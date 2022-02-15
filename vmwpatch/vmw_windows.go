// SPDX-FileCopyrightText: © 2014-2022 David Parsons
// SPDX-License-Identifier: MIT

//go:build windows
// +build windows

package vmwpatch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

var manager *mgr.Mgr

func IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		exe, _ := os.Executable()
		cwd, _ := os.Getwd()

		verbPtr, _ := syscall.UTF16PtrFromString("runas")
		exePtr, _ := syscall.UTF16PtrFromString(exe)
		cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
		argPtr, _ := syscall.UTF16PtrFromString(strings.Join(os.Args[1:], " "))

		err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, 1)
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	return true
}

func VMWStart(v *VMwareInfo) {
	fmt.Println()
	fmt.Println("Starting VMware services and tasks...")
	svcStart(v.AuthD)
	svcStart(v.USBD)
	svcStart(v.HostD)
	err := manager.Disconnect()
	if err != nil {
		fmt.Println("Disconnect from SCM failed")
		// Not stopping the process over this one
	}
	taskStart(filepath.Join(v.InstallDir, v.Tray))
}

func VMWStop(v *VMwareInfo) {
	fmt.Println()
	fmt.Println("Stopping VMware services and tasks...")

	var err error
	manager, err = mgr.Connect()
	if err != nil {
		panic("SCM connection failed")
	}

	svcStop(v.HostD)
	svcStop(v.USBD)
	svcStop(v.AuthD)
	taskStop(v.ShellExt) // No Need to re-exec this, it's part of a registered shell extension
	taskStop(v.Tray)
}

func VMWInfo() *VMwareInfo {
	v := &VMwareInfo{}

	// Store known service names
	v.AuthD = "VMAuthdService"
	v.HostD = "VMwareHostd"
	v.USBD = "VMUSBArbService"

	// Access registry for version, build and installation path
	var access uint32
	access = registry.QUERY_VALUE
	if runtime.GOARCH == "amd64" {
		access = access | registry.WOW64_32KEY
	}
	regKey, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\VMware, Inc.\VMware Player`, access)
	if err != nil {
		panic("Failed to open registry")
	}
	//goland:noinspection GoUnhandledErrorResult
	defer regKey.Close()

	v.ProductVersion, _, err = regKey.GetStringValue("ProductVersion")
	if err != nil {
		panic("Failed to locate registry key ProductVersion")
	}

	v.BuildNumber, _, err = regKey.GetStringValue("BuildNumber")
	if err != nil {
		panic("Failed to locate registry key BuildNumber")
	}

	v.InstallDir, _, err = regKey.GetStringValue("InstallPath")
	if err != nil {
		panic("Failed to locate registry key InstallPath")
	}

	// Construct needed filenames from reg settings
	v.InstallDir64 = filepath.Join(v.InstallDir, "x64")
	v.Player = "vmplayer.exe"
	v.Workstation = "vmware.exe"
	v.KVM = "vmware-kvm.exe"
	v.REST = "vmrest.exe"
	v.ShellExt = "vmware-shell-ext-thunker.exe"
	v.Tray = "vmware-tray.exe"
	v.VMXDefault = "vmware-vmx.exe"
	v.VMXDebug = "vmware-vmx-debug.exe"
	v.VMXStats = "vmware-vmx-stats.exe"
	v.VMwareBase = "vmwarebase.dll"
	v.PathVMXDefault = filepath.Join(v.InstallDir64, v.VMXDefault)
	v.PathVMXDebug = filepath.Join(v.InstallDir64, v.VMXDebug)
	v.PathVMXStats = filepath.Join(v.InstallDir64, v.VMXStats)
	v.PathVMwareBase = filepath.Join(v.InstallDir, v.VMwareBase)
	currentFolder, _ := os.Getwd()
	v.BackDir = filepath.Join(currentFolder, "backup", v.ProductVersion)
	v.BackVMXDefault = filepath.Join(v.BackDir, v.VMXDefault)
	v.BackVMXDebug = filepath.Join(v.BackDir, v.VMXDebug)
	v.BackVMXStats = filepath.Join(v.BackDir, v.VMXStats)
	v.BackVMwareBase = filepath.Join(v.BackDir, v.VMwareBase)
	v.PathISOMacOSX = filepath.Join(v.InstallDir, "darwinPre15.iso")
	v.PathISOmacOS = filepath.Join(v.InstallDir, "darwin.iso")
	return v
}

func setCTime(path string, ctime time.Time) error {
	//setCTime will set the create time on a file. On Windows, this requires
	//calling SetFileTime and explicitly including the create time.
	ctimespec := windows.NsecToTimespec(ctime.UnixNano())
	pathp, e := windows.UTF16PtrFromString(path)
	if e != nil {
		return e
	}
	h, e := windows.CreateFile(pathp,
		windows.FILE_WRITE_ATTRIBUTES, windows.FILE_SHARE_WRITE, nil,
		windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if e != nil {
		return e
	}
	//goland:noinspection GoUnhandledErrorResult
	defer windows.CloseHandle(h)
	c := windows.NsecToFiletime(windows.TimespecToNsec(ctimespec))
	return windows.SetFileTime(h, &c, nil, nil)
}

func svcState(s *mgr.Service) svc.State {
	status, err := s.Query()
	if err != nil {
		panic(fmt.Sprintf("Query(%s) failed: %s", s.Name, err))
	}
	return status.State
}

func svcWaitState(s *mgr.Service, want svc.State) {
	state := make(chan svc.State, 1)
	defer close(state)
	t := time.NewTimer(3 * time.Second)
	for {
		select {
		case <-t.C:
			panic(fmt.Sprintf("%s state change timeout", s.Name))
		case currentState := <-state:
			if currentState == want {
				t.Stop()
				return
			}
		case <-time.After(300 * time.Millisecond):
			state <- svcState(s)
		}
	}
}

func svcStart(name string) {
	s, err := manager.OpenService(name)
	if err != nil {
		return
	}

	fmt.Println("Starting service ", name)

	//goland:noinspection ALL
	defer s.Close()

	if svcState(s) == svc.Stopped {
		err = s.Start()
		if err != nil {
			panic(fmt.Sprintf("Control(%s) failed: %s", name, err))
		}
		svcWaitState(s, svc.Running)
	}
}

func svcStop(name string) {
	s, err := manager.OpenService(name)
	if err != nil {
		return
	}

	fmt.Println("Stopping service ", name)

	//goland:noinspection ALL
	defer s.Close()

	if svcState(s) == svc.Running {
		_, err = s.Control(svc.Stop)
		if err != nil {
			panic(fmt.Sprintf("Control(%s) failed: %s", name, err))
		}
		svcWaitState(s, svc.Stopped)
	}
}

func taskStart(filename string) {
	fmt.Println("Starting task ", filename)
	c := exec.Command(filename)
	_ = c.Start()
}

func taskStop(name string) {
	if TaskRunning(name) != 0 {
		fmt.Println("Stopping task ", name)
		c := exec.Command("taskkill.exe", "/F", "/IM", name)
		_ = c.Run()
	}
}
