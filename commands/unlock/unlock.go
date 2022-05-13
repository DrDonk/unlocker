// SPDX-FileCopyrightText: © 2014-2022 David Parsons
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
	fmt.Printf("Unlocker %s for VMware Workstation/Player\n", vmwpatch.VERSION)
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

	fmt.Println()
	fmt.Println("Installing unlock patches")

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
				p, _, _ := vmwpatch.IsGOSPatched(gos.FileToPatch)
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

	// Copy ISOs
	fmt.Println()
	fmt.Println("Copying VMware Tools...")
	_, err := vmwpatch.CopyFile(v.SrcISOMacOSX, v.DstISOMacOSX)
	if err != nil {
		fmt.Println("Error copying darwinPre15.iso")
	}
	_, err = vmwpatch.CopyFile(v.SrcISOmacOS, v.DstISOmacOS)
	if err != nil {
		fmt.Println("Error copying darwin.iso")
	}

	// Start all VMW services and tasks on Windows
	vmwpatch.VMWStart(v)

	waitExit()
	return
}
