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
		return
	} else {
		filename = os.Args[1]
	}

	//	Open the file
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		println(fmt.Sprintf("Cannot find file %s", filename))
		panic(err)
	}
	defer f.Close()

	// Memory map file
	contents, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		println("error mapping: %s", err)
		panic(err)
	}
	defer contents.Unmap()

	//Print titles
	println("PatchVMKCTL")
	println("-----------")
	println("© 2014-2021 David Parsons\n")
	println(fmt.Sprintf("File: %s", filename))
	println()

	offset := bytes.Index(contents, APPLESMC)
	before := string(contents[offset : offset+8])
	copy(contents[offset:offset+8], VMKERNEL)
	err = contents.Flush()
	if err != nil {
		panic(err)
	}
	after := string(contents[offset : offset+8])
	println(fmt.Sprintf("Patching %s -> %s", before, after))

}

func main() {
	vmkctl()
}
