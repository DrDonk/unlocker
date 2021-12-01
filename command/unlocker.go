// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"golocker/vmwpatch"
	"os"
	"path/filepath"
)

func printHelp() {
	fmt.Printf("usage: unlocker.exe <install | uninstall>\n")
	fmt.Printf("\tinstall - install patches\n")
	fmt.Printf("\tuninstall - uninstall patches\n")
}

func main() {
	// Titles
	fmt.Printf("Unlocker %s for VMware Workstation/Player\n", vmwpatch.VERSION)
	fmt.Printf("============================================\n")
	fmt.Printf("%s \n\n", vmwpatch.COPYRIGHT)

	// Simple arg parser
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	var install bool
	switch os.Args[1] {
	case "install":
		install = true
	case "uninstall":
		install = false
	default:
		printHelp()
		return
	}

	// Check admin rights
	if !vmwpatch.IsAdmin() {
		fmt.Printf("Re-run with admin/root privileges\n")
		return
	}

	// Get VMware product details from registry and file system
	v := vmwpatch.VMWInfo()
	fmt.Printf("VMware is installed at: %s\n", v.InstallDir)
	fmt.Printf("VMware version: %s\n\n", v.BuildNumber)

	// Check no VMs running
	if vmwpatch.IsRunning(v) {
		fmt.Printf("Aborting patching!\n")
		return
	}

	// Abort if installing and backup is present
	if install {
		fmt.Printf("Installing unlocker\n")
		if vmwpatch.BackupExists(v) {
			fmt.Printf("Aborting install as backup folder already exists!\n")
			return
		}
	}

	// Abort if uninstalling and backup is missing
	if !install {
		fmt.Printf("Uninstalling unlocker\n")
		if !vmwpatch.BackupExists(v) {
			fmt.Printf("Aborting uninstall as backup folder does not exist!\n")
			return
		}
	}

	// Stop all VMW services and tasks on Windows
	// Dummy calls on Linux
	vmwpatch.VMWStop(v)

	if install {
		// Backup files
		fmt.Printf("\nBacking up files...\n")
		vmwpatch.Backup(v)

		// Patch files
		fmt.Printf("\nPatching...\n")
		vmwpatch.PatchSMC(v.PathVMXDefault)
		fmt.Printf("\n")
		vmwpatch.PatchSMC(v.PathVMXDebug)
		fmt.Printf("\n")
		vmwpatch.PatchSMC(v.PathVMXStats)
		fmt.Printf("\n")
		vmwpatch.PatchGOS(v.PathVMwareBase)
		fmt.Printf("\n")

		// Copy iso ISOs
		fmt.Printf("\nCopying VMware Tools...\n")
		_, _ = vmwpatch.CopyFile("./iso/darwinPre15.iso", filepath.Join(v.InstallDir, "darwinPre15.iso"))
		_, _ = vmwpatch.CopyFile("./iso/darwin.iso", filepath.Join(v.InstallDir, "darwin.iso"))
	} else {
		// Restore files
		fmt.Printf("\nRestoring files...\n")
		vmwpatch.Restore(v)

		// Removing iso ISOs
		fmt.Printf("\nRemoving VMware Tools...\n")
		isoPath := filepath.Join(v.InstallDir, "darwinPre15.iso")
		fmt.Printf("%s\n", isoPath)
		_ = os.Remove(isoPath)
		isoPath = filepath.Join(v.InstallDir, "darwin.iso")
		fmt.Printf("%s\n", isoPath)
		_ = os.Remove(isoPath)
	}

	// Start all VMW services and tasks on Windows
	// Dummy calls on Linux
	vmwpatch.VMWStart(v)

	fmt.Printf("\nPress any key to finish...")
	//goland:noinspection GoUnhandledErrorResult
	fmt.Scanln()
	return
}
