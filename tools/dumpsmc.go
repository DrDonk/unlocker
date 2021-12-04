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
	fmt.Printf("DumpSMC %s\n", vmwpatch.VERSION)
	fmt.Printf("==============\n")
	fmt.Printf("%s \n\n", vmwpatch.COPYRIGHT)

	// Get and check file passed as parameter
	var filename string
	if len(os.Args) < 2 {
		fmt.Printf("Please pass a file name!\n")
		return
	} else {
		filename = os.Args[1]
	}

	//vmwpatch.DumpSMC(filename)
	state, hash256 := vmwpatch.IsPatched(filename)
	fmt.Printf("Patch Status: %d\nSHA256: %s\n", state, hash256)

}
