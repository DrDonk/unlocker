// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package vmwpatch

import (
	"fmt"
	"github.com/edsrzf/mmap-go"
	"os"
)

func checkFile(filename string) (contents mmap.MMap) {

	// Print filename
	println(fmt.Sprintf("File: %s", filename))

	// Open the file
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		panic(fmt.Sprintf("Cannot find file %s", filename))
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	// Memory map file
	contents, err = mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		panic(fmt.Sprintf("Error mapping file %s", filename))
	}
	defer func(contents *mmap.MMap) {
		err := contents.Unmap()
		if err != nil {

		}
	}(&contents)

	println(fmt.Sprintf("File: %s", filename))
	return contents
}
