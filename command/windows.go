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

func main() {
	println(fmt.Sprintf("Unlocker %s for VMware Workstation/Player", vmwpatch.VERSION))
	println("============================================")
	println(fmt.Sprintf("%s \n", vmwpatch.COPYRIGHT))

	v := vmwInfo()
	println(fmt.Sprintf("Patching version %s", v.ProductVersion))
	println(v.InstallDir)

	tray, _ := processID(v.Tray)
	println(tray)
	c := exec.Command(filepath.Join(v.InstallDir, "vmware-task.exe"))
	println(c)
	_ = c.Start()
	svcStop(v.AuthD)
	svcStop(v.HostD)
	svcStop(v.USBD)
	svcStart(v.AuthD)
	svcStart(v.HostD)
	svcStart(v.USBD)
	return
}
