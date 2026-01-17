// SPDX-FileCopyrightText: © 2014-2026 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"

	"github.com/drdonk/golocker/vmwpatch"
)

func main() {
	// Print titles
	fmt.Printf("PatchSMC %s\n", vmwpatch.VERSION)
	fmt.Printf("===============\n")

	// Get and check file passed as parameter
	var filename string
	if len(os.Args) < 2 {
		fmt.Printf("Please pass a file name!\n")
		return
	} else {
		filename = os.Args[1]
		fmt.Printf("Filename: %s\n\n", filename)
	}

	unpatched, patched := vmwpatch.PatchSMC(filename)
	fmt.Printf("\nSHA256\nunpatched: %s\npatched:   %s\n", unpatched, patched)
}
