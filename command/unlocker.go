// SPDX-FileCopyrightText: © 2014-2021 David Parsons
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

func printHelp() {
	fmt.Println("usage: unlocker.py.exe <install | uninstall>")
	fmt.Println("\tinstall - install patches")
	fmt.Println("\tuninstall - uninstall patches")
	waitExit()
}

func main() {
	// Titles
	fmt.Printf("Unlocker %s for VMware Workstation/Player\n", vmwpatch.VERSION)
	fmt.Println("============================================")
	fmt.Println(vmwpatch.COPYRIGHT)
	fmt.Println()

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
		fmt.Println("Re-run with admin/root privileges")
		waitExit()
		return
	}

	// Get VMware product details from registry and file system
	v := vmwpatch.VMWInfo()
	fmt.Println("VMware is installed at: ", v.InstallDir)
	fmt.Println("VMware version: ", v.BuildNumber)

	// Check no VMs running
	if vmwpatch.IsRunning(v) {
		fmt.Println("Aborting patching!")
		waitExit()
		return
	}

	// Stop all VMW services and tasks on Windows
	vmwpatch.VMWStop(v)

	if install {
		fmt.Println()
		fmt.Println("Installing unlocker.py")

		patchSmc := make(chan *vmwpatch.PatchOperation)
		patchGos := make(chan *vmwpatch.PatchOperation)
		done := make(chan int)
		go func() {
			for {
				select {
				case smc := <-patchSmc:
					p, _, _ := vmwpatch.IsSMCPatched(smc.FileToPatch)
					if p == 0 {
						fmt.Println("Patching", smc.FileToPatch)
						backupSuccessful := smc.Backup()
						if !backupSuccessful {
							fmt.Println(smc.BackupLocation, "already exists, skipping backup, still patching")
						}
						unpatched, patched := vmwpatch.PatchSMC(smc.FileToPatch)
						vmwpatch.WriteHashes(smc.BackupLocation, unpatched, patched)
						fmt.Println()
					} else {
						fmt.Println(smc.FileToPatch, "already patched, skipping")
					}
				case gos := <-patchGos:
					p, _ := vmwpatch.IsGOSPatched(gos.FileToPatch)
					if p == 0 {
						fmt.Println("Patching", gos.FileToPatch)
						backupSuccessful := gos.Backup()
						if !backupSuccessful {
							fmt.Println(gos.BackupLocation, "already exists, skipping backup, still patching")
						}
						unpatched, patched := vmwpatch.PatchGOS(gos.FileToPatch)
						vmwpatch.WriteHashes(gos.BackupLocation, unpatched, patched)
						fmt.Println()
					} else {
						fmt.Println(gos.FileToPatch, "already patched, skipping")
					}
				case <-done:
					fmt.Println("Patching Complete!")
					return
				}
			}
		}()

		// Patch files
		fmt.Println("Patching...")
		v.PatchFiles(patchGos, patchSmc, done)

		// Copy iso ISOs
		fmt.Println()
		fmt.Println("Copying VMware Tools...")
		_, err := vmwpatch.CopyFile("../iso/darwinPre15.iso", v.PathISOMacOSX)
		if err != nil {
			fmt.Println("Error copying darwinPre15.iso")
		}
		_, err = vmwpatch.CopyFile("../iso/darwin.iso", v.PathISOmacOS)
		if err != nil {
			fmt.Println("Error copying darwin.iso")
		}

	} else {
		fmt.Println("Uninstalling unlocker.py")

		// Check backup status
		if !v.BackupExists() {
			fmt.Println("Aborting uninstall as backup folder does not exist!")
			waitExit()
			return
		}

		// Restore files
		fmt.Println()
		fmt.Println("Restoring files...")
		v.Restore()

		// Removing iso ISOs
		fmt.Println()
		fmt.Println("Removing VMware Tools...")
		fmt.Printf(v.PathISOMacOSX)
		_ = os.Remove(v.PathISOMacOSX)
		fmt.Println(v.PathISOmacOS)
		_ = os.Remove(v.PathISOmacOS)
	}

	// Start all VMW services and tasks on Windows
	vmwpatch.VMWStart(v)

	waitExit()
	return
}
