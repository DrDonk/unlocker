// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

/*
vSMC Header Structure
=====================
Offset  Length  Struct Type Description
----------------------------------------
0x00/00 0x08/08 Q      ptr  Offset to key table
0x08/08 0x04/4  I      int  Number of private keys
0x0C/12 0x04/4  I      int  Number of public keys

vSMC Key Data Structure
Offset  Length  Struct Type Description
----------------------------------------
0x00/00 0x04/04 4s     int  Key name (byte reversed e.g. #KEY is YEK#)
0x04/04 0x01/01 B      byte Length of returned data
0x05/05 0x04/04 4s     int  Data type (byte reversed e.g. ui32 is 23iu)
0x09/09 0x01/01 B      byte Flag R/W
0x0A/10 0x06/06 6x     byte Padding
0x10/16 0x08/08 Q      ptr  Internal VMware routine
0x18/24 0x30/48 48B    byte Data

The internal VMware routines point to 4 variants:
AppleSMCHandleDefault
AppleSMCHandleNTOK
AppleSMCHandleNumKeys
AppleSMCHandleOSK
*/

package vmwpatch

import (
	"bytes"
	"fmt"
	"github.com/edsrzf/mmap-go"
	"os"
)

func dumpKeys(contents mmap.MMap, offset int, count int) {
	println(fmt.Sprintf("Table Offset : 0x%08x", offset))
	println("Offset     Name Len Type Flag FuncPtr    Data")
	println("-------    ---- --- ---- ---- -------    ----")

	// Loop for each count and print key
	// Last key should be OSK1
	var vmxKey smcKey
	for i := 0; i < count; i++ {
		// Unpack binary key data
		ptrCurrent := offset + (i * rowLength)
		vmxKey = getKey(contents, ptrCurrent)
		printKey(ptrCurrent, vmxKey)
	}
}

//goland:noinspection GoUnhandledErrorResult
func DumpSMC() {

	//Get and check file passed as parameter
	var filename string
	if len(os.Args) < 2 {
		println("Please pass a file name!")
	} else {
		filename = os.Args[1]
	}

	// Open the file
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
	//goland:noinspection Annotator
	defer contents.Unmap()

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

	// Print vSMC0 tables and keys
	vmxhdr0 := getHdr(contents, smcHeaderV0Offset)
	printHdr("0", smcHeaderV0Offset, vmxhdr0)
	dumpKeys(contents, smcKey0, int(vmxhdr0.cntPrivate))

	println("\n")

	// Print vSMC1 tables and keys
	vmxhdr1 := getHdr(contents, smcHeaderV1Offset)
	printHdr("1", smcHeaderV1Offset, vmxhdr1)
	dumpKeys(contents, smcKey1, int(vmxhdr1.cntPrivate))

}
