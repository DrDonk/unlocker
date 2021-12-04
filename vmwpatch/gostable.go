// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package vmwpatch

import (
	"fmt"
	binarypack "github.com/canhlinh/go-binary-pack"
	"github.com/edsrzf/mmap-go"
	"regexp"
)

func findGOSTable(contents mmap.MMap) [][]int {
	// Regexp pattern for GOS table Darwin entries
	pattern := "\x10\x00\x00\x00[\x10|\x20]\x00\x00\x00[\x01|\x02]\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"

	// Find all occurrences
	re, _ := regexp.Compile(pattern)
	indices := re.FindAllStringIndex(string(contents), -1)
	return indices
}

func setBit(n int, pos uint) int {
	n |= 1 << pos
	return n
}

func PatchGOS(filename string) (string, string) {

	// MMap the file
	f, contents := mapFile(filename)

	unpatched := sha256File(contents)

	// Find GOS table entries for Darwin
	indices := findGOSTable(contents)

	// Setup struct pack string
	var flagPack = []string{"b"}

	// Create BinaryPack object
	bp := new(binarypack.BinaryPack)

	for _, index := range indices {

		// Unpack binary key data
		offset := index[0] + 32
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

		// Print details
		fmt.Printf("Flag patched @ offset: 0x%08x  Flag: 0x%01x -> 0x%01x\n", offset, oldFlag, newFlag)

	}

	// Flush to disk
	flushFile(contents)
	patched := sha256File(contents)
	unmapFile(f, contents)

	return unpatched, patched

}

func IsGOSPatched(filename string) (int, string) {

	// MMap the file
	f, contents := mapFile(filename)

	// Check if the file is already patched
	indices := findGOSTable(contents)
	patched := 0
	if indices == nil {
		patched = 1
	}
	hash256 := sha256File(contents)
	unmapFile(f, contents)
	return patched, hash256
}
