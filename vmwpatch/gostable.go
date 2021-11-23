// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package vmwpatch

import (
	"fmt"
	binarypack "github.com/canhlinh/go-binary-pack"
	"github.com/edsrzf/mmap-go"
	"os"
	"regexp"
)

func setBit(n int, pos uint) int {
	n |= 1 << pos
	return n
}

//goland:noinspection GoUnhandledErrorResult
func PatchGOS() {

	// Get and check file passed as parameter
	var filename string
	if len(os.Args) < 2 {
		println("Please pass a file name!")
		return
	} else {
		filename = os.Args[1]
	}

	// Open the file
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		println(fmt.Sprintf("Cannot find file %s", filename))
		println(err)
		return
	}
	defer f.Close()

	// Memory map file
	contents, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		println("error mapping: %s", err)
	}
	defer contents.Unmap()

	println(fmt.Sprintf("File: %s", filename))

	pattern := "\x10\x00\x00\x00[\x10|\x20]\x00\x00\x00[\x01|\x02]\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"

	re, _ := regexp.Compile(pattern)
	indices := re.FindAllStringIndex(string(contents), -1)

	// Setup struct pack string
	var flagPack = []string{"b"}

	// Create BinaryPack object
	bp := new(binarypack.BinaryPack)

	for _, index := range indices {
		offset := index[0] + 32

		// Unpack binary key data
		unpackFlag, err := bp.UnPack(flagPack, contents[offset:offset+32])
		if err != nil {
			panic(err)
		}

		// Loop through each entry and set top bit
		// 0xBE --> 0xBF (WKS 12/13)
		// 0x3E --> 0x3F (WKS 14+)
		oldFlag := unpackFlag[0].(int)
		newFlag := setBit(oldFlag, 0)

		// Pack binary key data
		flagPacked, err := bp.Pack(flagPack, []interface{}{newFlag})
		if err != nil {
			panic(err)
		}

		// Copy data to mmap file
		copy(contents[offset:offset+1], flagPacked)

		//Flush to disk
		err = contents.Flush()
		if err != nil {
			panic(err)
		}

		println(fmt.Sprintf("Flag patched @ offset: 0x%08x  Flag: 0x%01x -> 0x%01x", offset, oldFlag, newFlag))

	}
}
