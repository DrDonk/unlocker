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
	vmwpatch.IsAdmin()

	// Print titles
	fmt.Printf("DumpSMC %s\n", vmwpatch.VERSION)
	fmt.Printf("==============\n")
	fmt.Printf("%s \n\n", vmwpatch.COPYRIGHT)
	// Get and check file passed as parameter
	var filename string
	if len(os.Args) < 2 {
		fmt.Printf("Please pass a file name!\n")
		waitExit()
		return
	} else {
		filename = os.Args[1]
		fmt.Printf("Filename: %s\n\n", filename)
	}

	vmwpatch.DumpSMC(filename)
	patchFlag, patchStatus, hash256 := vmwpatch.IsSMCPatched(filename)
	fmt.Printf("\nPatch Status: %s (%d)\nSHA256: %s\n", patchStatus, patchFlag, hash256)
	waitExit()

}
