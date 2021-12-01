// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

//go:build darwin
// +build darwin

package vmwpatch

import (
	"os"
	"path/filepath"
	"time"
)

func IsAdmin() bool {
	if os.Geteuid() == 0 {
		return true
	}
	return false
}

//goland:noinspection GoUnusedParameter
func VMWStart(v *VMwareInfo) {
	// Dummy function on Linux
	return
}

//goland:noinspection GoUnusedParameter
func VMWStop(v *VMwareInfo) {
	// Dummy function on Linux
	return
}

func VMWInfo() *VMwareInfo {
	v := &VMwareInfo{}

	// Store known service names
	// Not used on macOS
	v.AuthD = ""
	v.HostD = ""
	v.USBD = ""

	// TODO: Parse macOS Fusion specific files
	// /Library/Preferences/VMware Fusion/lastLocationUsed --> App Bundle
	// /Applications/VMware Fusion.app/Contents/Info.plist
	//	<key>CFBundleShortVersionString</key>
	//	<string>12.2.1</string>
	//	<key>CFBundleVersion</key>
	//	<string>18811640</string>

	// Access /etc/vmware/config for version, build and installation path
	file, err := os.Open("/Library/Preferences/VMware Fusion/lastLocationUsed")
	if err != nil {
		panic(err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	// Basic product settings
	v.ProductVersion = ""
	v.BuildNumber = ""
	v.InstallDir = ""

	// Construct needed filenames from reg settings
	v.InstallDir64 = ""
	v.Player = ""
	v.Workstation = "/Applications/VMware Fusion.app/Contents/MacOS/VMware Fusion"
	v.KVM = ""
	v.REST = "vmrest"
	v.Tray = ""
	v.VMXDefault = "vmware-vmx"
	v.VMXDebug = "vmware-vmx-debug"
	v.VMXStats = "vmware-vmx-stats"
	v.VMwareBase = "libvmwarebase.dylib"
	v.PathVMXDefault = filepath.Join(v.InstallDir, "Contents", "Library", "vmware-vmx")
	v.PathVMXDebug = filepath.Join(v.InstallDir, "Contents", "Library", "vmware-vmx-debug")
	v.PathVMXStats = filepath.Join(v.InstallDir, "Contents", "Library", "vmware-vmx-stats")
	v.PathVMwareBase = filepath.Join(v.InstallDir, "Contents", "Frameworks", "libvmwarebase.dylib")
	return v
}

//goland:noinspection GoUnusedParameter
func setCTime(path string, ctime time.Time) error {
	// Dummy function on Linux
	return nil
}
