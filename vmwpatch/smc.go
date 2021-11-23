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
	"encoding/hex"
	"fmt"
	binarypack "github.com/canhlinh/go-binary-pack"
	"github.com/edsrzf/mmap-go"
)

const hdrLength = 16
const keyLength = 24
const dataLength = 48
const rowLength = keyLength + dataLength

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

func fourCCToString(s string) (result string) {
	for _, v := range s {
		if v != 0 {
			result = string(v) + result
		} else {
			result = " " + result
		}
	}
	return
}

//goland:noinspection GoUnusedExportedFunction
func stringToFourCC(s string) (result string) {
	for _, v := range s {
		if v != 32 {
			result = string(v) + result
		} else {
			result = "\x00" + result
		}
	}
	return
}

func printHdr(version string, offset int, vmxHdr smcHdr) {
	println(fmt.Sprintf("appleSMCTableV%s (smc.version = '%s')", version, version))
	println(fmt.Sprintf("File Offset  : 0x%08x", offset))
	println(fmt.Sprintf("Keys Address : 0x%08x", vmxHdr.address))
	println(fmt.Sprintf("Private Key #: 0x%04x/%04d", vmxHdr.cntPublic, vmxHdr.cntPublic))
	println(fmt.Sprintf("Public Key  #: 0x%04x/%04d", vmxHdr.cntPrivate, vmxHdr.cntPrivate))
	println("")
}

func printKey(offset int, vmxKey smcKey) {
	// Convert binary string to hexdump
	data := hex.EncodeToString([]byte(vmxKey.data)[0:vmxKey.length])

	// Print the key
	println(fmt.Sprintf("0x%08x %04s %02d  %-04s 0x%02x 0x%08x %s",
		offset,
		vmxKey.key,
		vmxKey.length,
		vmxKey.dataType,
		vmxKey.flag,
		vmxKey.ptrFunc,
		data))
	return
}

func getHdr(contents mmap.MMap, offset int) smcHdr {
	// Setup struct pack string
	var hdrPack = []string{"Q", "I", "I"}

	// Create BinaryPack object
	bp := new(binarypack.BinaryPack)

	// Unpack binary key data
	hdr, err := bp.UnPack(hdrPack, contents[offset:offset+hdrLength])
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

func getKey(contents mmap.MMap, offset int) smcKey {
	// Setup struct pack string
	var keyPack = []string{"4s", "B", "4s", "B", "B", "B", "B", "B", "B", "B", "Q", "48s"}

	// Create BinaryPack object
	bp := new(binarypack.BinaryPack)

	// Unpack binary key data
	keyRow, err := bp.UnPack(keyPack, contents[offset:offset+rowLength])
	if err != nil {
		println(err)
	}

	// Return the smcKey as a struct
	var vmxKey smcKey
	vmxKey.key = fourCCToString(keyRow[0].(string))
	vmxKey.length = byte(keyRow[1].(int))
	vmxKey.dataType = fourCCToString(keyRow[2].(string))
	vmxKey.flag = byte(keyRow[3].(int))
	vmxKey.ptrFunc = uintptr(keyRow[10].(int))
	vmxKey.data = keyRow[11].(string)
	return vmxKey
}
