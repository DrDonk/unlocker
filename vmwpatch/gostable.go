// SPDX-FileCopyrightText: © 2014-2022 David Parsons
// SPDX-License-Identifier: MIT

package vmwpatch

import (
	"fmt"
	"regexp"

	binarypack "github.com/canhlinh/go-binary-pack"
)

func findGOSTable(contents []byte) [][]int {
	// Regexp pattern for GOS table Darwin entries
	pattern := "\x10\x00\x00\x00[\x10|\x20]\x00\x00\x00[\x01|\x02]\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"

	// Find all occurrences
	re, _ := regexp.Compile(pattern)
	indices := re.FindAllStringIndex(string(contents), -1)
	return indices
}

//goland:noinspection GoUnusedFunction
func clearBit(n int, pos uint) int {
	mask := ^(1 << pos)
	n &= mask
	return n
}

func hasBit(n int, pos uint) bool {
	val := n & (1 << pos)
	return val > 0
}

func setBit(n int, pos uint) int {
	n |= 1 << pos
	return n
}

func PatchGOS(filename string) (string, string) {

	// Read the file
	contents := loadFile(filename)

	unpatched := sha256File(contents)

	// Find GOS table entries for Darwin
	indices := findGOSTable(contents)

	// Setup struct pack string
	var flagPack = []string{"b"}

	// Create BinaryPack object
	bp := new(binarypack.BinaryPack)

	count := 0
	for _, index := range indices {
		count++
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

		// Copy data to file
		copy(contents[offset:offset+1], flagPacked)

		// Print details
		fmt.Printf("Flag patched @ offset: 0x%08x  Flag: 0x%01x -> 0x%01x\n", offset, oldFlag, newFlag)
	}
	fmt.Printf("Patched %d flags\n", count)

	// Flush to disk
	patched := sha256File(contents)
	saveFile(filename, contents)

	return unpatched, patched

}

func IsGOSPatched(filename string) (int, string) {

	// MMap the file
	contents := loadFile(filename)

	// Check if the file is already patched
	indices := findGOSTable(contents)

	// Setup struct pack string
	var flagPack = []string{"b"}

	// Create BinaryPack object
	patched := 0
	count := 0
	for _, index := range indices {

		// Keep a counter for checks
		count++

		// Unpack binary key data
		offset := index[0] + 32
		unpackFlag, err := binarypack.New().UnPack(flagPack, contents[offset:offset+32])
		if err != nil {
			panic(err)
		}

		// Loop through each entry and test top bit
		// 0xBF (WKS 12/13)
		// 0x3F (WKS 14+)
		oldFlag := unpackFlag[0].(int)
		if hasBit(oldFlag, 0) {
			patched++
		}
	}

	// Check patched byte count
	switch patched {
	case 0:
		patched = 0
	case count:
		patched = 1
	default:
		patched = 2
	}

	hash256 := sha256File(contents)
	saveFile(filename, contents)
	return patched, hash256
}
