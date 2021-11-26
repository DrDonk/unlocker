// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"golocker/vmwpatch"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type VMwareInfo struct {
	BuildNumber    string
	ProductVersion string
	InstallDir     string
	InstallDir64   string
	Workstation    string
	Player         string
	KVM            string
	REST           string
	Tray           string
	AuthD          string
	HostD          string
	USBD           string
	VMXDefault     string
	VMXDebug       string
	VMXStats       string
	VMwareBase     string
}

// unsafe.Sizeof(windows.ProcessEntry32{})
const processEntrySize = 568

func processID(name string) (uint32, error) {
	h, e := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if e != nil {
		return 0, e
	}
	p := windows.ProcessEntry32{Size: processEntrySize}
	for {
		e := windows.Process32Next(h, &p)
		if e != nil {
			return 0, e
		}
		if windows.UTF16ToString(p.ExeFile[:]) == name {
			return p.ProcessID, nil
		}
	}
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
		//println(fmt.Sprintf("Invalid service %s", name))
		return
	} else {
		println(fmt.Sprintf("Starting service %s", name))
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
		//println(fmt.Sprintf("Invalid service %s", name))
		return
	} else {
		println(fmt.Sprintf("Stopping service %s", name))
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
	c := exec.Command(filename)
	_ = c.Start()
	return
}

func taskRunning(name string) bool {
	pid, err := processID(name)
	if (pid != 0) && (err == nil) {
		return true
	} else {
		return false
	}
}

func taskStop(name string) {
	if taskRunning(name) {
		c := exec.Command("taskkill.exe", "/F", "/IM", name)
		_ = c.Run()
	}
	return
}

func vmwInfo() *VMwareInfo {
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
	v.VMXDefault = filepath.Join(v.InstallDir, "x64", "vmware-vmx.exe")
	v.VMXDebug = filepath.Join(v.InstallDir, "x64", "vmware-vmx-debug.exe")
	v.VMXStats = filepath.Join(v.InstallDir, "x64", "vmware-vmx-stats.exe")
	v.VMwareBase = filepath.Join(v.InstallDir, "vmwarebase.dll")

	return v
}

func vmwRunning(v *VMwareInfo) bool {
	if taskRunning(v.Workstation) {
		println("VMware Workstation is running")
		return true
	}
	if taskRunning(v.Player) {
		println("VMware Player is running")
		return true
	}
	if taskRunning(v.KVM) {
		println("VMware KVM is running")
		return true
	}
	if taskRunning(v.REST) {
		println("VMware REST API is running")
		return true
	}
	if taskRunning(v.VMXDefault) {
		println("VMware VM (vmware-vmx) is running")
		return true
	}
	if taskRunning(v.VMXDebug) {
		println("VMware VM (vmware-vmx-debug) is running")
		return true
	}
	if taskRunning(v.VMXStats) {
		println("VMware VM (vmware-vmx-stats) is running")
		return true
	}
	return false
}

func main() {
	// Titles
	println(fmt.Sprintf("Unlocker %s for VMware Workstation/Player", vmwpatch.VERSION))
	println("============================================")
	println(fmt.Sprintf("%s \n", vmwpatch.COPYRIGHT))

	// Get VMware product details from registry and file system
	v := vmwInfo()
	println(fmt.Sprintf("VMware is installed at: %s", v.InstallDir))
	println(fmt.Sprintf("Patching VMware product version %s", v.ProductVersion))

	// Check no VMs running
	if vmwRunning(v) {
		println("Aborting patching!")
		return
	}

	// Stop all services and task with open handles to the exe/dll files
	println("Stopping VMware services and tasks...")
	svcStop(v.AuthD)
	svcStop(v.HostD)
	svcStop(v.USBD)
	taskStop(v.Tray)

	//Backing up files...

	//Patching...

	//Getting VMware Tools...

	println("Starting VMware services and tasks...")
	svcStart(v.AuthD)
	svcStart(v.HostD)
	svcStart(v.USBD)
	taskStart(filepath.Join(v.InstallDir, v.Tray))

	println("Finished!")
	return
}
