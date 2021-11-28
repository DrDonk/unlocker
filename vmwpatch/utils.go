// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package vmwpatch

import (
	"fmt"
	"github.com/edsrzf/mmap-go"
	"os"
)

func mapFile(filename string) (file *os.File, contents mmap.MMap) {

	// Print filename
	fmt.Printf("File: %s\n", filename)

	// Open the file
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		panic("Cannot find file")
	}

	// Memory map file
	contents, err = mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		panic("Error mapping file")
	}
	return f, contents
}

func flushFile(contents mmap.MMap) {
	err := contents.Flush()
	if err != nil {
		panic("Error flushing file")
	}
	return
}

func unmapFile(file *os.File, contents mmap.MMap) {
	err := contents.Unmap()
	if err != nil {
		panic("Error unmapping file")
	}

	err = file.Close()
	if err != nil {
		panic("Error closing file")
	}
	return
}
