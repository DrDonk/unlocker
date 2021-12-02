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
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	binarypack "github.com/canhlinh/go-binary-pack"
	"github.com/edsrzf/mmap-go"
	"unsafe"
)

const hdrLength = 16
const keyLength = 24
const dataLength = 48
const rowLength = keyLength + dataLength

//goland:noinspection GoUnusedConst
const kppwData = "\x53\x70\x65\x63\x69\x61\x6c\x69\x73\x52\x65\x76\x65\x6c\x69\x6f\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"
const kpstData = "\x01"
const osk0Data = "\x6f\x75\x72\x68\x61\x72\x64\x77\x6f\x72\x6b\x62\x79\x74\x68\x65\x73\x65\x77\x6f\x72\x64\x73\x67\x75\x61\x72\x64\x65\x64\x70\x6c"
const osk1Data = "\x65\x61\x73\x65\x64\x6f\x6e\x74\x73\x74\x65\x61\x6c\x28\x63\x29\x41\x70\x70\x6c\x65\x43\x6f\x6d\x70\x75\x74\x65\x72\x49\x6e\x63"

const elfMagic = "7f454c46"

type smcHdr struct {
	address    uintptr
	cntPrivate uint16
	cntPublic  uint16
}

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

func patchELF(contents mmap.MMap, AppleSMCHandleOSK uintptr, AppleSMCHandleDefault uintptr) {
	// Process ELF RELA records
	progType := hex.EncodeToString(contents[0:4])
	if progType == elfMagic {
		fmt.Printf("\nModifying ELF RELA records from 0x%08x -> 0x%08x\n", AppleSMCHandleOSK, AppleSMCHandleDefault)
		defPtr := ptrToBytes(AppleSMCHandleDefault)
		oskPtr := ptrToBytes(AppleSMCHandleOSK)
		var relaPtr int
		for i := 0; i < 4; i++ {
			relaPtr = bytes.Index(contents, oskPtr)
			if relaPtr != -1 {
				// Replace the function pointer
				fmt.Printf("Relocation modified at: 0x%08x\n", relaPtr)
				copy(contents[relaPtr:relaPtr+8], defPtr)

				// Flush to disk
				flushFile(contents)
			}
		}
	}
}

func printHdr(version string, offset int, vmxHdr smcHdr) {
	fmt.Printf("appleSMCTableV%s (smc.version = '%s')\n", version, version)
	fmt.Printf("File Offset  : 0x%08x\n", offset)
	fmt.Printf("Keys Address : 0x%08x\n", vmxHdr.address)
	fmt.Printf("Private Key #: 0x%04x/%04d\n", vmxHdr.cntPublic, vmxHdr.cntPublic)
	fmt.Printf("Public Key  #: 0x%04x/%04d\n", vmxHdr.cntPrivate, vmxHdr.cntPrivate)
	fmt.Printf("\n")
}

func printKey(offset int, vmxKey smcKey) {
	// Convert binary string to hexdump
	data := hex.EncodeToString([]byte(vmxKey.data)[0:vmxKey.length])

	// Print the key
	fmt.Printf("0x%08x %04s %02d  %-04s 0x%02x 0x%08x %s\n",
		offset,
		vmxKey.key,
		vmxKey.length,
		vmxKey.dataType,
		vmxKey.flag,
		vmxKey.ptrFunc,
		data)
	return
}

func findHdrs(contents mmap.MMap) (int, int) {

	// Find the vSMC headers
	var smcHeaderV0 = []byte{0xF2, 0x00, 0x00, 0x00, 0xF0, 0x00, 0x00, 0x00}
	var smcHeaderV1 = []byte{0xB4, 0x01, 0x00, 0x00, 0xB0, 0x01, 0x00, 0x00}
	smcHeaderV0Offset := bytes.Index(contents, smcHeaderV0) - 8
	smcHeaderV1Offset := bytes.Index(contents, smcHeaderV1) - 8
	return smcHeaderV0Offset, smcHeaderV1Offset
}

func findKEY(contents mmap.MMap) (int, int) {
	// Find '#KEY' keys
	var keyKey = []byte{0x59, 0x45, 0x4B, 0x23, 0x04, 0x32, 0x33, 0x69, 0x75}
	smcKey0 := bytes.Index(contents, keyKey)
	smcKey1 := bytes.LastIndex(contents, keyKey)
	return smcKey0, smcKey1
}

func getHdr(contents mmap.MMap, offset int) smcHdr {
	// Setup struct pack string
	var hdrPack = []string{"Q", "I", "I"}

	// Create BinaryPack object
	bp := new(binarypack.BinaryPack)

	// Unpack binary key data
	hdr, err := bp.UnPack(hdrPack, contents[offset:offset+hdrLength])
	if err != nil {
		panic(err)
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
		panic(err)
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

func dumpKeys(contents mmap.MMap, offset int, count int) {
	fmt.Printf("Table Offset : 0x%08x\n", offset)
	fmt.Printf("Offset     Name Len Type Flag FuncPtr    Data\n")
	fmt.Printf("-------    ---- --- ---- ---- -------    ----\n")

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
	flushFile(contents)
	return
}

func patchKeys(contents mmap.MMap, offset int, count int) (uintptr, uintptr) {
	// Loop for each count and print key
	// Last key should be OSK1
	var vmxKey smcKey
	var AppleSMCHandleDefault uintptr
	var AppleSMCHandleOSK uintptr

	fmt.Printf("Table Offset : 0x%08x\n", offset)
	for i := 0; i < count; i++ {
		// Unpack binary key data
		ptrCurrent := offset + (i * rowLength)
		vmxKey = getKey(contents, ptrCurrent)
		switch vmxKey.key {
		case "+LKS":
			fmt.Printf("Getting +LKS:\n")
			printKey(ptrCurrent, vmxKey)
			AppleSMCHandleDefault = vmxKey.ptrFunc
		case "KPPW":
			fmt.Printf("Patching KPPW:\n")
			printKey(ptrCurrent, vmxKey)
			//vmxKey.data = kppwData
			hash := sha256.Sum256(contents)
			vmxKey.length = 32
			vmxKey.data = string(hash[:])
			putKey(contents, ptrCurrent, vmxKey)
			vmxKey = getKey(contents, ptrCurrent)
			printKey(ptrCurrent, vmxKey)
		case "KPST":
			fmt.Printf("Patching KPST:\n")
			printKey(ptrCurrent, vmxKey)
			vmxKey.data = kpstData
			putKey(contents, ptrCurrent, vmxKey)
			vmxKey = getKey(contents, ptrCurrent)
			printKey(ptrCurrent, vmxKey)
		case "OSK0":
			fmt.Printf("Patching OSK0:\n")
			printKey(ptrCurrent, vmxKey)
			AppleSMCHandleOSK = vmxKey.ptrFunc
			vmxKey.ptrFunc = AppleSMCHandleDefault
			vmxKey.data = osk0Data
			putKey(contents, ptrCurrent, vmxKey)
			vmxKey = getKey(contents, ptrCurrent)
			printKey(ptrCurrent, vmxKey)
		case "OSK1":
			fmt.Printf("Patching OSK1:\n")
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

func DumpSMC(filename string) {

	// MMap the file
	f, contents := mapFile(filename)

	// Find the vSMC headers
	smcHeaderV0Offset, smcHeaderV1Offset := findHdrs(contents)

	// Find '#KEY' keys
	smcKey0, smcKey1 := findKEY(contents)

	// Print vSMC0 tables and keys
	vmxhdr0 := getHdr(contents, smcHeaderV0Offset)
	printHdr("0", smcHeaderV0Offset, vmxhdr0)
	dumpKeys(contents, smcKey0, int(vmxhdr0.cntPrivate))

	fmt.Printf("\n")

	// Print vSMC1 tables and keys
	vmxhdr1 := getHdr(contents, smcHeaderV1Offset)
	printHdr("1", smcHeaderV1Offset, vmxhdr1)
	dumpKeys(contents, smcKey1, int(vmxhdr1.cntPrivate))

	// Unmap file
	unmapFile(f, contents)

}

func PatchSMC(filename string) {

	// MMap the file
	f, contents := mapFile(filename)

	// Check if the file is already patched
	patched := bytes.Index(contents, []byte(osk0Data))
	if patched != -1 {
		fmt.Printf("Aborting as file %s already patched\n", filename)
		unmapFile(f, contents)
		return
	}

	// Find the vSMC headers
	smcHeaderV0Offset, smcHeaderV1Offset := findHdrs(contents)

	// Find '#KEY' keys
	smcKey0, smcKey1 := findKEY(contents)

	// Patch vSMC0 tables and keys
	vmxhdr0 := getHdr(contents, smcHeaderV0Offset)
	printHdr("0", smcHeaderV0Offset, vmxhdr0)
	patchKeys(contents, smcKey0, int(vmxhdr0.cntPrivate))

	fmt.Printf("\n")

	// Patch vSMC1 tables and keys
	vmxhdr1 := getHdr(contents, smcHeaderV1Offset)
	printHdr("1", smcHeaderV1Offset, vmxhdr1)
	AppleSMCHandleOSK, AppleSMCHandleDefault := patchKeys(contents, smcKey1, int(vmxhdr1.cntPrivate))
	patchELF(contents, AppleSMCHandleOSK, AppleSMCHandleDefault)

	// Flush to disk
	flushFile(contents)
	unmapFile(f, contents)
}
