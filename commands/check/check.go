//
// SPDX-FileCopyrightText: © 2014-2023 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"github.com/drdonk/golocker/vmwpatch"
	"os"
)

func waitExit() {
	quiet := os.Getenv("UNLOCK_QUIET")
	if quiet != "1" {
		fmt.Println()
		fmt.Println("\nPress Enter key to continue...")

		_, _ = fmt.Scanln()
	}
}

func main() {
	// Titles
	fmt.Printf("Checker %s for VMware Workstation/Player\n", vmwpatch.VERSION)
	fmt.Println("============================================")
	fmt.Println()

	// Check admin rights
	if !vmwpatch.IsAdmin() {
		fmt.Println("Re-run with admin/root privileges")
		waitExit()
		return
	}

	// Get VMware product details from registry and file system
	v := vmwpatch.VMWInfo()
	if v.ProductVersion != "" {
		fmt.Println("VMware is installed at: ", v.InstallDir)
		fmt.Println("VMware version: ", v.BuildNumber)
	} else {
		fmt.Println("VMware is not installed on this computer")
		waitExit()
		return
	}
	patchFlag, patchStatus, hash256 := vmwpatch.IsSMCPatched(v.PathVMXDefault)
	fmt.Printf("Patch Status: %s (%d)\nSHA256: %s\n", patchStatus, patchFlag, hash256)

	patchFlag, patchStatus, hash256 = vmwpatch.IsSMCPatched(v.PathVMXDebug)
	fmt.Printf("Patch Status: %s (%d)\nSHA256: %s\n", patchStatus, patchFlag, hash256)

	_, err := os.Stat(v.PathVMXStats)
	if err == nil {
		patchFlag, patchStatus, hash256 = vmwpatch.IsSMCPatched(v.PathVMXStats)
		fmt.Printf("Patch Status: %s (%d)\nSHA256: %s\n", patchStatus, patchFlag, hash256)
	}

	patchFlag, patchStatus, hash256 = vmwpatch.IsGOSPatched(v.PathVMwareBase)
	fmt.Printf("Patch Status: %s (%d)\nSHA256: %s\n", patchStatus, patchFlag, hash256)

	waitExit()
	return
}
