// SPDX-FileCopyrightText: © 2014-2023 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"

	"github.com/drdonk/golocker/vmwpatch"
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
	fmt.Printf("Relocker %s for VMware Workstation/Player\n", vmwpatch.VERSION)
	fmt.Println("============================================")
	fmt.Println(vmwpatch.COPYRIGHT)
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

	// Check no VMs running
	if vmwpatch.IsRunning(v) {
		fmt.Println("Aborting patching!")
		waitExit()
		return
	}

	// Stop all VMW services and tasks on Windows
	vmwpatch.VMWStop(v)

	fmt.Println("Uninstalling unlock patches")

	// Check backup status
	if !v.BackupExists() {
		fmt.Println("Aborting uninstall as backup folder does not exist!")
		// Ensure Services restart even with failed relock
		vmwpatch.VMWStart(v)
		waitExit()
		return
	}

	// Restore files
	fmt.Println()
	fmt.Println("Restoring files...")
	v.Restore()

	// Removing ISOs
	fmt.Println()
	fmt.Println("Removing VMware Tools...")
	fmt.Println(v.DstISOMacOSX)
	_ = os.Remove(v.DstISOMacOSX)
	fmt.Println(v.DstISOmacOS)
	_ = os.Remove(v.DstISOmacOS)

	// Start all VMW services and tasks on Windows
	vmwpatch.VMWStart(v)

	waitExit()
	return
}
