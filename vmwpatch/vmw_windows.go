// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

//go:build windows
// +build windows

package vmwpatch

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

func IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func VMWStart(v *VMwareInfo) {
	fmt.Printf("\nStarting VMware services and tasks...\n")
	svcStart(v.AuthD)
	svcStart(v.HostD)
	svcStart(v.USBD)
	taskStart(filepath.Join(v.InstallDir, v.Tray))
}

func VMWStop(v *VMwareInfo) {
	fmt.Printf("\nStopping VMware services and tasks...\n")
	svcStop(v.AuthD)
	svcStop(v.HostD)
	svcStop(v.USBD)
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
	ctimespec := syscall.NsecToTimespec(ctime.UnixNano())
	pathp, e := syscall.UTF16PtrFromString(path)
	if e != nil {
		return e
	}
	h, e := syscall.CreateFile(pathp,
		syscall.FILE_WRITE_ATTRIBUTES, syscall.FILE_SHARE_WRITE, nil,
		syscall.OPEN_EXISTING, syscall.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if e != nil {
		return e
	}
	//goland:noinspection GoUnhandledErrorResult
	defer syscall.Close(h)
	c := syscall.NsecToFiletime(syscall.TimespecToNsec(ctimespec))
	return syscall.SetFileTime(h, &c, nil, nil)
}

func svcState(s *mgr.Service) svc.State {
	status, err := s.Query()
	if err != nil {
		panic(fmt.Sprintf("Query(%s) failed: %s", s.Name, err))
	}
	return status.State
}

func svcWaitState(s *mgr.Service, want svc.State) {
	for i := 0; ; i++ {
		have := svcState(s)
		if have == want {
			return
		}
		if i > 10 {
			panic(fmt.Sprintf("%s state is=%d, waiting timeout", s.Name, have))
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func svcStart(name string) {
	m, err := mgr.Connect()
	if err != nil {
		panic("SCM connection failed")
	}

	//goland:noinspection ALL
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		//fmt.Printf("Invalid service %s", name))
		return
	} else {
		fmt.Printf("Starting service %s\n", name)
	}

	//goland:noinspection ALL
	defer s.Close()

	if svcState(s) == svc.Stopped {
		err = s.Start()
		if err != nil {
			panic(fmt.Sprintf("Control(%s) failed: %s", name, err))
		}
		svcWaitState(s, svc.Running)
	}

	err = m.Disconnect()

}

func svcStop(name string) {
	m, err := mgr.Connect()
	if err != nil {
		panic("SCM connection failed")

	}

	//goland:noinspection ALL
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		//fmt.Printf("Invalid service %s", name))
		return
	} else {
		fmt.Printf("Stopping service %s\n", name)
	}

	//goland:noinspection ALL
	defer s.Close()

	if svcState(s) == svc.Running {
		_, err = s.Control(svc.Stop)
		if err != nil {
			panic(fmt.Sprintf("Control(%s) failed: %s", name, err))
		}
		svcWaitState(s, svc.Stopped)
	}

	err = m.Disconnect()

}

func taskStart(filename string) {
	fmt.Printf("Starting task %s\n", filename)
	c := exec.Command(filename)
	_ = c.Start()
	return
}

func taskStop(name string) {
	if TaskRunning(name) != 0 {
		fmt.Printf("Stopping task %s\n", name)
		c := exec.Command("taskkill.exe", "/F", "/IM", name)
		_ = c.Run()
	}
	return
}
