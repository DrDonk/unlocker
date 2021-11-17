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

package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/canhlinh/go-binary-pack"
	"github.com/edsrzf/mmap-go"
	"os"
)

const HdrLength = 16
const KeyLength = 24
const DataLength = 48
const RowLength = KeyLength + DataLength

const kppwData = "5370656369616c6973526576656c696f"
const kpstData = "01"
const osk0Data = "6f757268617264776f726b62797468657365776f72647367756172646564706c"
const osk1Data = "65617365646f6e74737465616c2863294170706c65436f6d7075746572496e63"

//goland:noinspection GoUnusedType
type smcHdr struct {
	address    uintptr
	cntPrivate uint16
	cntPublic  uint16
}

//goland:noinspection GoUnusedType
type smcKey struct {
	key      string
	length   byte
	dataType string
	flag     byte
	_        [6]byte
	ptrFunc  uintptr
	data     string
}

func reverse(s string) (result string) {
	for _, v := range s {
		if v != 0 {
			result = string(v) + result
		} else {
			result = " " + result
		}
	}
	return
}

func printhdr(version string, offset int, vmxHdr smcHdr) {
	println(fmt.Sprintf("appleSMCTableV%s (smc.version = '%s')", version, version))
	println(fmt.Sprintf("File Offset  : 0x%08x", offset))
	println(fmt.Sprintf("Keys Address : 0x%08x", vmxHdr.address))
	println(fmt.Sprintf("Private Key #: 0x%04x/%04d", vmxHdr.cntPublic, vmxHdr.cntPublic))
	println(fmt.Sprintf("Public Key  #: 0x%04x/%04d", vmxHdr.cntPrivate, vmxHdr.cntPrivate))
	println("")
}

func printkey(offset int, vmxKey smcKey) {
	println(fmt.Sprintf("0x%08x %04s %02d  %-04s 0x%02x 0x%08x %s",
		offset,
		vmxKey.key,
		vmxKey.length,
		vmxKey.dataType,
		vmxKey.flag,
		vmxKey.ptrFunc,
		vmxKey.data))
	return
}

func gethdr(contents []byte, offset int) smcHdr {
	// Setup struct pack string
	var hdrPack = []string{"Q", "I", "I"}

	// Create BinaryPack object
	bp := new(binarypack.BinaryPack)

	// Unpack binary key data
	hdr, err := bp.UnPack(hdrPack, contents[offset:offset+HdrLength])
	if err != nil {
		println(err)
	}

	// Return the smcHdr as a struct
	var vmxHdr smcHdr
	vmxHdr.address = uintptr(hdr[0].(int))
	vmxHdr.cntPrivate = uint16(hdr[1].(int))
	vmxHdr.cntPublic = uint16(hdr[2].(int))
	return vmxHdr
}

func getkey(contents []byte, offset int) smcKey {
	// Setup struct pack string
	var keyPack = []string{"4s", "B", "4s", "B", "B", "B", "B", "B", "B", "B", "Q", "48s"}

	// Create BinaryPack object
	bp := new(binarypack.BinaryPack)

	// Unpack binary key data
	keyRow, err := bp.UnPack(keyPack, contents[offset:offset+RowLength])
	if err != nil {
		println(err)
	}

	// Return the smcKey as a struct
	var vmxKey smcKey
	vmxKey.key = reverse(keyRow[0].(string))
	vmxKey.length = byte(keyRow[1].(int))
	vmxKey.dataType = reverse(keyRow[2].(string))
	vmxKey.flag = byte(keyRow[3].(int))
	vmxKey.ptrFunc = uintptr(keyRow[10].(int))
	vmxKey.data = hex.EncodeToString([]byte(keyRow[11].(string))[0:vmxKey.length])

	return vmxKey
}

//func putkey(contents []byte, offset int, vmxKey smcKey) {
//	// Setup struct pack string
//	var keyPack = []string{"4s", "B", "4s", "B", "B", "B", "B", "B", "B", "B", "Q", "48s"}
//
//	// Create BinaryPack object
//	bp := new(binarypack.BinaryPack)
//
//	// Unpack binary key data
//	keyRow, err := bp.Pack(keyPack, abc)
//
//	.UnPack(keyPack, contents[offset:offset+RowLength])
//	if err != nil {
//		println(err)
//	}

//// Return the smcKey as a struct
//var vmxKey smcKey
//vmxKey.key = reverse(keyRow[0].(string))
//vmxKey.length = byte(keyRow[1].(int))
//vmxKey.dataType = reverse(keyRow[2].(string))
//vmxKey.flag = byte(keyRow[3].(int))
//vmxKey.ptrFunc = uintptr(keyRow[10].(int))
//vmxKey.data = hex.EncodeToString([]byte(keyRow[11].(string))[0:vmxKey.length])

//	return
//}

func patchkeys(contents []byte, offset int, count int, funcAddr uintptr) {
	println(fmt.Sprintf("Table Offset : 0x%08x", offset))
	println("Offset     Name Len Type Flag FuncPtr    Data")
	println("-------    ---- --- ---- ---- -------    ----")

	// Loop for each count and print key
	// Last key should be OSK1
	var vmxKey smcKey
	for i := 0; i < count; i++ {
		// Unpack binary key data
		ptrCurrent := offset + (i * RowLength)
		vmxKey = getkey(contents, ptrCurrent)
		switch vmxKey.key {
		case "KPPW":
			//data, _ := hex.DecodeString(kppwData)
			printkey(ptrCurrent, vmxKey)
			vmxKey.data = kppwData
			printkey(ptrCurrent, vmxKey)
		case "KPST":
			printkey(ptrCurrent, vmxKey)
			vmxKey.data = kpstData
			printkey(ptrCurrent, vmxKey)
		case "OSK0":
			//data, _ := hex.DecodeString(osk0Data)
			printkey(ptrCurrent, vmxKey)
			vmxKey.ptrFunc = funcAddr
			vmxKey.data = osk0Data
			printkey(ptrCurrent, vmxKey)
		case "OSK1":
			//data, _ := hex.DecodeString(osk1Data)
			printkey(ptrCurrent, vmxKey)
			vmxKey.ptrFunc = funcAddr
			vmxKey.data = osk1Data
			printkey(ptrCurrent, vmxKey)
		}
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
	println("patchsmc")
	println("--------")
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

	// Find '$Adr' key in V0 table and used to patch OSK0 & OSK1 key functions
	var adrKey = []byte{0x72, 0x64, 0x41, 0x24, 0x04, 0x32, 0x33, 0x69, 0x75}
	smcAdr := bytes.Index(contents, adrKey)
	vmxAdr := getkey(contents, smcAdr)
	println(fmt.Sprintf("0x%08x", vmxAdr.ptrFunc))

	// Print vSMC0 tables and keys
	vmxhdr0 := gethdr(contents, smcHeaderV0Offset)
	printhdr("0", smcHeaderV0Offset, vmxhdr0)
	patchkeys(contents, smcKey0, int(vmxhdr0.cntPrivate), vmxAdr.ptrFunc)

	println("\n")

	// Print vSMC1 tables and keys
	vmxhdr1 := gethdr(contents, smcHeaderV1Offset)
	printhdr("1", smcHeaderV1Offset, vmxhdr1)
	patchkeys(contents, smcKey1, int(vmxhdr1.cntPrivate), vmxAdr.ptrFunc)

}

func main() {
	vSMC()
}
