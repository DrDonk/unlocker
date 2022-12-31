// SPDX-FileCopyrightText: © 2014-2023 David Parsons
// SPDX-License-Identifier: MIT

//go:build darwin

package vmwpatch

import (
	"os"
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
	// Dummy function on macOS
	return
}

//goland:noinspection GoUnusedParameter
func VMWStop(v *VMwareInfo) {
	// Dummy function on macOS
	return
}

func VMWInfo() *VMwareInfo {
	// TODO: Parse macOS Fusion specific files
	// /Library/Preferences/VMware Fusion/lastLocationUsed --> App Bundle
	// /Applications/VMware Fusion.app/Contents/Info.plist
	//	<key>CFBundleShortVersionString</key>
	//	<string>12.2.1</string>
	//	<key>CFBundleVersion</key>
	//	<string>18811640</string>

	v := &VMwareInfo{}
	return v
}

//goland:noinspection GoUnusedParameter
func setCTime(path string, ctime time.Time) error {
	// Dummy function on macOS
	return nil
}
