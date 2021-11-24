// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"golocker/vmwpatch"
	"os"
)

func main() {
	// Print titles
	println("PatchGOS")
	println("-----------")
	println("© 2014-2021 David Parsons\n")

	// Get and check file passed as parameter
	var filename string
	if len(os.Args) < 2 {
		println("Please pass a file name!")
		return
	} else {
		filename = os.Args[1]
	}

	vmwpatch.PatchGOS(filename)
}
