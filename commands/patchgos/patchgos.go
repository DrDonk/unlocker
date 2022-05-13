// SPDX-FileCopyrightText: © 2014-2022 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"

	"github.com/drdonk/golocker/vmwpatch"
)

func main() {
	// Print titles
	fmt.Printf("PatchGOS %s\n", vmwpatch.VERSION)
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

	patchFlag, _, hash256 := vmwpatch.IsGOSPatched(filename)

	switch patchFlag {
	case 0:
		unpatched, patched := vmwpatch.PatchGOS(filename)
		fmt.Printf("\nSHA256\nunpatched: %s\npatched:   %s\n", unpatched, patched)
	case 1:
		fmt.Printf("File %s is already patched\n", filename)
		fmt.Printf("Patch Status: %d\nSHA256: %s\n", patchFlag, hash256)
		return
	case 2:
		fmt.Printf("File %s is in an indeterminate state\n", filename)
		return
	default:
		fmt.Printf("Unknown issue with filename %s\n", filename)
		return
	}

}
