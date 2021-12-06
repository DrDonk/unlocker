// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"golocker/vmwpatch"
	"os"
)

func waitExit() {
	quiet := os.Getenv("UNLOCK_QUIET")
	if quiet != "1" {
		fmt.Printf("\nPress any key to continue...")
		//goland:noinspection GoUnhandledErrorResult
		fmt.Scanln()
	}
}

func printHelp() {
	fmt.Printf("usage: unlocker.exe <install | uninstall>\n")
	fmt.Printf("\tinstall - install patches\n")
	fmt.Printf("\tuninstall - uninstall patches\n")
	waitExit()
}

func main() {
	// Titles
	fmt.Printf("Unlocker %s for VMware Workstation/Player\n", vmwpatch.VERSION)
	fmt.Printf("============================================\n")
	fmt.Printf("%s\n\n", vmwpatch.COPYRIGHT)

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
		waitExit()
		return
	}

	// Get VMware product details from registry and file system
	v := vmwpatch.VMWInfo()
	fmt.Printf("VMware is installed at: %s\n", v.InstallDir)
	fmt.Printf("VMware version: %s\n", v.BuildNumber)

	// Check no VMs running
	if vmwpatch.IsRunning(v) {
		fmt.Printf("Aborting patching!\n")
		waitExit()
		return
	}

	// Stop all VMW services and tasks on Windows
	// Dummy calls on Linux
	vmwpatch.VMWStop(v)

	if install {
		fmt.Printf("\nInstalling unlocker\n")

		// Check patch status
		fmt.Printf("Checking patch status of files...\n")
		if vmwpatch.CheckStatus(v) != 0 {
			fmt.Printf("Aborting install as files already patched!\n")
			waitExit()
			return
		}

		// Check backup status
		if vmwpatch.BackupExists(v) {
			fmt.Printf("Aborting install as backup folder already exists!\n")
			waitExit()
			return
		}

		// Backup files
		fmt.Printf("\nBacking up files...\n")
		vmwpatch.Backup(v)

		// Patch files
		fmt.Printf("\nPatching...\n")
		unpatched, patched := vmwpatch.PatchSMC(v.PathVMXDefault)
		vmwpatch.WriteHashes(v.BackVMXDefault, unpatched, patched)
		fmt.Printf("\n")
		unpatched, patched = vmwpatch.PatchSMC(v.PathVMXDebug)
		vmwpatch.WriteHashes(v.BackVMXDebug, unpatched, patched)
		fmt.Printf("\n")
		_, err := os.Stat(v.PathVMXStats)
		if err == nil {
			unpatched, patched = vmwpatch.PatchSMC(v.PathVMXStats)
			vmwpatch.WriteHashes(v.BackVMXStats, unpatched, patched)
			fmt.Printf("\n")
		}
		unpatched, patched = vmwpatch.PatchGOS(v.PathVMwareBase)
		vmwpatch.WriteHashes(v.BackVMwareBase, unpatched, patched)

		// Copy iso ISOs
		fmt.Printf("\nCopying VMware Tools...\n")
		_, err = vmwpatch.CopyFile("./iso/darwinPre15.iso", v.PathISOMacOSX)
		if err != nil {
			fmt.Printf("Error copying darwinPre15.iso")
		}
		_, err = vmwpatch.CopyFile("./iso/darwin.iso", v.PathISOmacOS)
		if err != nil {
			fmt.Printf("Error copying darwin.iso")
		}

	} else {
		fmt.Printf("\nUninstalling unlocker\n")

		// Check patch status
		fmt.Printf("Checking patch status of files...\n")
		if vmwpatch.CheckStatus(v) != 1 {
			fmt.Printf("Aborting install as files already unpatched!\n")
			waitExit()
			return
		}

		// Check backup status
		if !vmwpatch.BackupExists(v) {
			fmt.Printf("Aborting uninstall as backup folder does not exist!\n")
			waitExit()
			return
		}

		// Restore files
		fmt.Printf("\nRestoring files...\n")
		vmwpatch.Restore(v)

		// Removing iso ISOs
		fmt.Printf("\nRemoving VMware Tools...\n")
		fmt.Printf("%s\n", v.PathISOMacOSX)
		_ = os.Remove(v.PathISOMacOSX)
		fmt.Printf("%s\n", v.PathISOmacOS)
		_ = os.Remove(v.PathISOmacOS)
	}

	// Start all VMW services and tasks on Windows
	// Dummy calls on Linux
	vmwpatch.VMWStart(v)

	waitExit()
	return
}
