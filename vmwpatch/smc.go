// SPDX-FileCopyrightText: © 2014-2023 David Parsons
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
	"os"
	"unsafe"

	binarypack "github.com/canhlinh/go-binary-pack"
)

const hdrLength = 16
const keyLength = 24
const dataLength = 48
const rowLength = keyLength + dataLength

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
			result = " " + result
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

func patchELF(contents []byte, AppleSMCHandleOSK uintptr, AppleSMCHandleDefault uintptr) {
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
			}
		}
	}
}

func printHdr(version string, offset int, vmxHdr smcHdr) {
	fmt.Printf("appleSMCTableV%s (smc.version = '%s')\n", version, version)
	fmt.Printf("File Offset  : 0x%08x\n", offset)
	fmt.Printf("Keys Address : 0x%08x\n", vmxHdr.address)
	fmt.Printf("Public Keys  : 0x%04x/%04d\n", vmxHdr.cntPublic, vmxHdr.cntPublic)
	fmt.Printf("Private Keys : 0x%04x/%04d\n", vmxHdr.cntPrivate, vmxHdr.cntPrivate)
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

func findHdrs(contents []byte) (int, int) {

	// Find the vSMC headers
	var smcHeaderV0 = []byte{0xF2, 0x00, 0x00, 0x00, 0xF0, 0x00, 0x00, 0x00}
	var smcHeaderV1 = []byte{0xB4, 0x01, 0x00, 0x00, 0xB0, 0x01, 0x00, 0x00}
	smcHeaderV0Offset := bytes.Index(contents, smcHeaderV0) - 8
	smcHeaderV1Offset := bytes.Index(contents, smcHeaderV1) - 8
	return smcHeaderV0Offset, smcHeaderV1Offset
}

func findKey(contents []byte) (int, int) {
	// Find '#KEY' keys
	var keyKey = []byte{0x59, 0x45, 0x4B, 0x23, 0x04, 0x32, 0x33, 0x69, 0x75}
	smcKey0 := bytes.Index(contents, keyKey)
	smcKey1 := bytes.LastIndex(contents, keyKey)
	return smcKey0, smcKey1
}

func checkPatch(contents []byte) (int, string) {
	// Check if the file is already patched
	osk0 := bytes.Index(contents, []byte(osk0Data))
	osk1 := bytes.Index(contents, []byte(osk1Data))

	patchFlag := 0
	patchStatus := ""
	if osk0 == -1 && osk1 == -1 {
		patchFlag = 0
		patchStatus = "Unpatched"
	} else if osk0 != -1 && osk1 != -1 {
		patchFlag = 1
		patchStatus = "Patched"
	} else {
		patchFlag = 2
		patchStatus = "Unknown"
	}
	return patchFlag, patchStatus
}

func getHdr(contents []byte, offset int) smcHdr {
	// Setup struct pack string
	var hdrPack = []string{"Q", "I", "I"}

	// Unpack binary key data
	hdr, err := binarypack.New().UnPack(hdrPack, contents[offset:offset+hdrLength])
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

func getKey(contents []byte, offset int) smcKey {
	// Setup struct pack string
	var keyPack = []string{"4s", "B", "4s", "B", "B", "B", "B", "B", "B", "B", "Q", "48s"}

	// Unpack binary key data
	keyRow, err := binarypack.New().UnPack(keyPack, contents[offset:offset+rowLength])
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

func dumpKeys(contents []byte, offset int, count int) {
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

func putKey(contents []byte, offset int, vmxKey smcKey) {
	// Setup struct pack string
	var keyPack = []string{"4s", "B", "4s", "B", "B", "B", "B", "B", "B", "B", "Q", "48s"}

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
	keyPacked, err := binarypack.New().Pack(keyPack, keyRow)
	if err != nil {
		panic(err)
	}

	// Copy data to file
	copy(contents[offset:offset+rowLength], keyPacked)

	return
}

func patchKeys(contents []byte, offset int, count int) (uintptr, uintptr) {
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

	// Read the file
	contents := loadFile(filename)

	// Find the vSMC headers
	smcHeaderV0Offset, smcHeaderV1Offset := findHdrs(contents)

	// Find '#KEY' keys
	smcKey0, smcKey1 := findKey(contents)

	// Print vSMC0 tables and keys
	vmxhdr0 := getHdr(contents, smcHeaderV0Offset)
	printHdr("0", smcHeaderV0Offset, vmxhdr0)
	dumpKeys(contents, smcKey0, int(vmxhdr0.cntPrivate))

	fmt.Printf("\n")

	// Print vSMC1 tables and keys
	vmxhdr1 := getHdr(contents, smcHeaderV1Offset)
	printHdr("1", smcHeaderV1Offset, vmxhdr1)
	dumpKeys(contents, smcKey1, int(vmxhdr1.cntPrivate))

	return
}

func PatchSMC(filename string) (string, string) {

	// Read the file
	contents := loadFile(filename)
	unpatched256 := sha256File(contents)

	// Internal patch checker
	patchFlag, _ := checkPatch(contents)

	// Guard checks
	if patchFlag == 1 {
		fmt.Printf("File %s is already patched\n", filename)
		fmt.Printf("Patch Status: %d\nSHA256: %s\n", patchFlag, unpatched256)
		os.Exit(1)
	}
	if patchFlag == 2 {
		fmt.Printf("File %s is in an indeterminate state\n", filename)
		os.Exit(2)
	}
	if patchFlag > 2 {
		fmt.Printf("Unknown issue with filename %s\n", filename)
		os.Exit(3)
	}

	// Find the vSMC headers
	smcHeaderV0Offset, smcHeaderV1Offset := findHdrs(contents)

	// Find '#KEY' keys
	smcKey0, smcKey1 := findKey(contents)

	// Patch vSMC0 tables and keys
	vmxhdr0 := getHdr(contents, smcHeaderV0Offset)
	printHdr("0", smcHeaderV0Offset, vmxhdr0)
	_, _ = patchKeys(contents, smcKey0, int(vmxhdr0.cntPrivate))

	fmt.Println()

	// Patch vSMC1 tables and keys
	vmxhdr1 := getHdr(contents, smcHeaderV1Offset)
	printHdr("1", smcHeaderV1Offset, vmxhdr1)
	AppleSMCHandleOSK, AppleSMCHandleDefault := patchKeys(contents, smcKey1, int(vmxhdr1.cntPrivate))
	patchELF(contents, AppleSMCHandleOSK, AppleSMCHandleDefault)

	// Flush to disk
	patched256 := sha256File(contents)
	saveFile(filename, contents)

	return unpatched256, patched256
}

func IsSMCPatched(filename string) (int, string, string) {

	// Read the file
	contents := loadFile(filename)

	// Internal patch checker
	patchFlag, patchStatus := checkPatch(contents)

	// Calc sha256
	hash256 := sha256File(contents)
	return patchFlag, patchStatus, hash256
}
