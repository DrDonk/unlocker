// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"golocker/vmwpatch"
	"os"
)

func main() {
	// Print titles
	fmt.Printf("PatchSMC %s\n", vmwpatch.VERSION)
	fmt.Printf("===============\n")
	fmt.Printf("%s \n\n", vmwpatch.COPYRIGHT)

	// Get and check file passed as parameter
	var filename string
	if len(os.Args) < 2 {
		fmt.Printf("Please pass a file name!\n")
		return
	} else {
		filename = os.Args[1]
	}

	state, hash256 := vmwpatch.IsSMCPatched(filename)

	switch state {
	case 0:
		unpatched, patched := vmwpatch.PatchSMC(filename)
		fmt.Printf("\nSHA256\nunpatched: %s\npatched:   %s\n", unpatched, patched)
	case 1:
		fmt.Printf("File %s is already patched\n", filename)
		fmt.Printf("Patch Status: %d\nSHA256: %s\n", state, hash256)
		return
	case 2:
		fmt.Printf("File %s is in an indeterminate state\n", filename)
		return
	default:
		fmt.Printf("Unknown issue with filename %s\n", filename)
		return
	}

}
