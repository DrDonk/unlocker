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
	fmt.Printf("PatchVMKCTL\n")
	fmt.Printf("-----------\n")
	fmt.Printf("© 2014-2021 David Parsons\n\n")

	// Get and check file passed as parameter
	var filename string
	if len(os.Args) < 2 {
		fmt.Printf("Please pass a file name!\n")
		return
	} else {
		filename = os.Args[1]
	}

	vmwpatch.PatchVMKCTL(filename)
}
