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
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/canhlinh/go-binary-pack"
	"github.com/edsrzf/mmap-go"
	"os"
	"unsafe"
)

const hdrLength = 16
const keyLength = 24
const dataLength = 48
const rowLength = keyLength + dataLength

const kpstData = "\x01"
const kppwData = "\x53\x70\x65\x63\x69\x61\x6c\x69\x73\x52\x65\x76\x65\x6c\x69\x6f\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"
const osk0Data = "\x6f\x75\x72\x68\x61\x72\x64\x77\x6f\x72\x6b\x62\x79\x74\x68\x65\x73\x65\x77\x6f\x72\x64\x73\x67\x75\x61\x72\x64\x65\x64\x70\x6c"
const osk1Data = "\x65\x61\x73\x65\x64\x6f\x6e\x74\x73\x74\x65\x61\x6c\x28\x63\x29\x41\x70\x70\x6c\x65\x43\x6f\x6d\x70\x75\x74\x65\x72\x49\x6e\x63"

const elfMagic = "7f454c46"

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

func ptrToBytes(ptr uintptr) []byte {
	size := unsafe.Sizeof(ptr)
	bytePtr := make([]byte, size)
	switch size {
	case 4:
		binary.LittleEndian.PutUint32(bytePtr, uint32(ptr))
	case 8:
		binary.LittleEndian.PutUint64(bytePtr, uint64(ptr))
	default:
		panic(fmt.Sprintf("unknown uintptr size: %v", size))
	}
	return bytePtr
}

//goland:noinspection Annotator
func patchELF(contents mmap.MMap, AppleSMCHandleOSK uintptr, AppleSMCHandleDefault uintptr) {
	// Process ELF RELA records
	progType := hex.EncodeToString(contents[0:4])
	if progType == elfMagic {
		println(fmt.Sprintf("\nModifying ELF RELA records from 0x%08x -> 0x%08x", AppleSMCHandleOSK, AppleSMCHandleDefault))
		defPtr := ptrToBytes(AppleSMCHandleDefault)
		oskPtr := ptrToBytes(AppleSMCHandleOSK)
		var relaPtr int
		for {
			relaPtr = bytes.Index(contents, oskPtr)
			if relaPtr != -1 {
				// Replace the function pointer
				println(fmt.Sprintf("Relocation modified at: 0x%08x", relaPtr))
				copy(contents[relaPtr:relaPtr+8], defPtr)

				// Flush to disk
				err := contents.Flush()
				if err != nil {
					panic(err)
				}
			} else {
				break
			}
		}
	}
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

//goland:noinspection Annotator
func putKey(contents mmap.MMap, offset int, vmxKey smcKey) {
	// Setup struct pack string
	var keyPack = []string{"4s", "B", "4s", "B", "B", "B", "B", "B", "B", "B", "Q", "48s"}

	// Create BinaryPack object
	bp := new(binarypack.BinaryPack)

	//keyData, _ := hex.DecodeString(vmxKey.data)
	keyRow := []interface{}{
		stringToFourCC(vmxKey.key),
		int(vmxKey.length),
		stringToFourCC(vmxKey.dataType),
		int(vmxKey.flag),
		0, // Padding
		0, // Padding
		0, // Padding
		0, // Padding
		0, // Padding
		0, // Padding
		int(vmxKey.ptrFunc),
		vmxKey.data,
	}

	// Pack binary key data
	keyPacked, err := bp.Pack(keyPack, keyRow)
	if err != nil {
		panic(err)
	}

	// Copy data to mmap file
	copy(contents[offset:offset+rowLength], keyPacked)

	// Flush to disk
	err = contents.Flush()
	if err != nil {
		panic(err)
	}
	return
}

func patchKeys(contents mmap.MMap, offset int, count int) (uintptr, uintptr) {
	// Loop for each count and print key
	// Last key should be OSK1
	var vmxKey smcKey
	var AppleSMCHandleDefault uintptr
	var AppleSMCHandleOSK uintptr

	println(fmt.Sprintf("Table Offset : 0x%08x", offset))
	for i := 0; i < count; i++ {
		// Unpack binary key data
		ptrCurrent := offset + (i * rowLength)
		vmxKey = getKey(contents, ptrCurrent)
		switch vmxKey.key {
		case "+LKS":
			println("Getting +LKS:")
			printKey(ptrCurrent, vmxKey)
			AppleSMCHandleDefault = vmxKey.ptrFunc
		case "KPPW":
			println("Patching KPPW:")
			printKey(ptrCurrent, vmxKey)
			vmxKey.data = kppwData
			putKey(contents, ptrCurrent, vmxKey)
			vmxKey = getKey(contents, ptrCurrent)
			printKey(ptrCurrent, vmxKey)
		case "KPST":
			println("Patching KPST:")
			printKey(ptrCurrent, vmxKey)
			vmxKey.data = kpstData
			putKey(contents, ptrCurrent, vmxKey)
			vmxKey = getKey(contents, ptrCurrent)
			printKey(ptrCurrent, vmxKey)
		case "OSK0":
			println("Patching OSK0:")
			printKey(ptrCurrent, vmxKey)
			AppleSMCHandleOSK = vmxKey.ptrFunc
			vmxKey.ptrFunc = AppleSMCHandleDefault
			vmxKey.data = osk0Data
			putKey(contents, ptrCurrent, vmxKey)
			vmxKey = getKey(contents, ptrCurrent)
			printKey(ptrCurrent, vmxKey)
		case "OSK1":
			println("Patching OSK1:")
			printKey(ptrCurrent, vmxKey)
			vmxKey.ptrFunc = AppleSMCHandleDefault
			vmxKey.data = osk1Data
			putKey(contents, ptrCurrent, vmxKey)
			vmxKey = getKey(contents, ptrCurrent)
			printKey(ptrCurrent, vmxKey)
		}
	}
	return AppleSMCHandleOSK, AppleSMCHandleDefault
}

//goland:noinspection GoUnhandledErrorResult,Annotator,Annotator
func PatchSMC() {

	// Get and check file passed as parameter
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

	// Patch vSMC0 tables and keys
	vmxhdr0 := getHdr(contents, smcHeaderV0Offset)
	printHdr("0", smcHeaderV0Offset, vmxhdr0)
	patchKeys(contents, smcKey0, int(vmxhdr0.cntPrivate))

	println("\n")

	// Patch vSMC1 tables and keys
	vmxhdr1 := getHdr(contents, smcHeaderV1Offset)
	printHdr("1", smcHeaderV1Offset, vmxhdr1)
	AppleSMCHandleOSK, AppleSMCHandleDefault := patchKeys(contents, smcKey1, int(vmxhdr1.cntPrivate))
	patchELF(contents, AppleSMCHandleOSK, AppleSMCHandleDefault)
	contents.Flush()
}
