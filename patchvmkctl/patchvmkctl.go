// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"fmt"
	"github.com/edsrzf/mmap-go"
	"os"
)

//goland:noinspection GoUnhandledErrorResult
func vmkctl() {

	var APPLESMC = []byte("applesmc")
	var VMKERNEL = []byte("vmkernel")

	//Get and check file passed as parameter
	var filename string
	if len(os.Args) < 2 {
		println("Please pass a file name!")
	} else {
		filename = os.Args[1]
	}

	//	Open the file
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		println(fmt.Sprintf("Cannot find file %s", filename))
		println(err)
	}
	defer f.Close()

	// Memory map file
	contents, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		println("error mapping: %s", err)
	}
	defer contents.Unmap()

	//Print titles
	println("PatchVMKCTL")
	println("-----------")
	println("© 2014-2021 David Parsons\n")
	println(fmt.Sprintf("File: %s", filename))
	println()

	offset := bytes.Index(contents, APPLESMC)
	println(string(contents[offset : offset+8]))
	copy(contents[offset:offset+8], VMKERNEL)
	contents.Flush()

	println(string(contents[offset : offset+8]))

}

func main() {
	vmkctl()
}
