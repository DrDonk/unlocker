// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/canhlinh/go-binary-pack"
	"github.com/edsrzf/mmap-go"
	"golocker/utils"
	"os"
)


func dumpkeys(contents []byte, offset int, count int) {
	println(fmt.Sprintf("Table Offset : 0x%08x", offset))
	println("Offset     Name Len Type Flag FuncPtr    Data")
	println("-------    ---- --- ---- ---- -------    ----")

	// Loop for each count and print key
	// Last key should be OSK1
	var vmxKe
	for i := 0; i < count; i++ {
		// Unpack binary key data
		ptrCurrent := offset + (i * RowLength)
		vmxKey = getkey(contents, ptrCurrent)
		printkey(ptrCurrent, vmxKey)
	}
}

//goland:noinspection GoUnhandledErrorResult
func vSMC() {

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
	println("dumpsmc")
	println("-------")
	println(fmt.Sprintf("File: %s", filename))
	println()

	// Find the vSMC headers
	var smcHeaderV0 = []byte{0xF2, 0x00, 0x00, 0x00, 0xF0, 0x00, 0x00, 0x00}
	var smcHeaderV1 = []byte{0xB4, 0x01, 0x00, 0x00, 0xB0, 0x01, 0x00, 0x00}
	smcHeaderV0Offset := bytes.Index(contents, smcHeaderV0) - 8
	smcHeaderV1Offset := bytes.Index(contents, smcHeaderV1) - 8

	// Find '#KEY' keys
	var keyKey = []byte{0x59, 0x45, 0x4B, 0x23, 0x04, 0x32, 0x33, 0x69, 0x75}
	smcKey0 := bytes.Index(contents, keyKey)
	smcKey1 := bytes.LastIndex(contents, keyKey)

	//TODO: Remove this when unlocker is re-written
	// Find '$Adr' key in V0 table and used to patch OSK0 & OSK1 key functions
	//var adrKey = []byte{0x72, 0x64, 0x41, 0x24, 0x04, 0x32, 0x33, 0x69, 0x75}
	//smcAdr := bytes.Index(contents, adrKey)
	//println(fmt.Sprintf("0x%08x", smcAdr))

	// Print vSMC0 tables and keys
	vmxhdr0 := gethdr(contents, smcHeaderV0Offset)
	printhdr("0", smcHeaderV0Offset, vmxhdr0)
	dumpkeys(contents, smcKey0, int(vmxhdr0.cntPrivate))

	println("\n")

	// Print vSMC1 tables and keys
	vmxhdr1 := gethdr(contents, smcHeaderV1Offset)
	printhdr("1", smcHeaderV1Offset, vmxhdr1)
	dumpkeys(contents, smcKey1, int(vmxhdr1.cntPrivate))

}

func main() {
	vSMC()
}
