// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package vmwpatch

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/edsrzf/mmap-go"
)

func mapFile(filename string, flag int) (file *os.File, contents mmap.MMap) {

	// Check valid flag of R0 or RDWR
	var mapFlag int
	switch flag {
	case os.O_RDONLY:
		mapFlag = mmap.RDONLY
	case os.O_RDWR:
		mapFlag = mmap.RDWR
	default:
		flag = os.O_RDONLY
		mapFlag = mmap.RDONLY
	}

	// Print filename
	fmt.Printf("File: %s Mode: %d\n", filename, flag)

	// Open the file
	f, err := os.OpenFile(filename, flag, 0444)

	if err != nil {
		panic("Cannot find file")
	}

	// Memory map file
	contents, err = mmap.Map(f, mapFlag, 0)
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
	_ = contents.Unmap()
	err := file.Close()
	if err != nil {
		panic("Error closing file")
	}
	return
}

func sha256File(contents mmap.MMap) string {
	return fmt.Sprintf("%x", sha256.Sum256(contents))
}
