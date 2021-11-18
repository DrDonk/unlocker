// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

const APPLESMC = "applesmc"
const VMKERNEL = "vmkernel"

//goland:noinspection GoUnhandledErrorResult
func vmkctl() {

	//Get and check file passed as parameter
	var filename string
	if len(os.Args) < 2 {
		println("Please pass a file name!")
	} else {
		filename = os.Args[1]
	}

	//	Open the file
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		println(fmt.Sprintf("Cannot find file %s", filename))
		println(err)
	}

	//Print titles
	println("PatchVMKCTL")
	println("-----------")
	println("© 2014-2021 David Parsons\n")
	println(fmt.Sprintf("File: %s", filename))
	println()

	offset := bytes.Index(contents, []byte(APPLESMC))
	println(fmt.Sprintf("0x%08x", offset))

	//	Open the file to write
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		println(fmt.Sprintf("Cannot find file %s", filename))
		println(err)
	}
	defer f.Close()
	f.Seek(int64(offset), 0)
	f.Write([]byte(VMKERNEL))
	f.Sync()
	f.Close()
}

func main() {
	vmkctl()
}
