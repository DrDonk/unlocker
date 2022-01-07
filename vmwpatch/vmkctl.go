// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package vmwpatch

import (
	"bytes"
	"fmt"
)

//goland:noinspection GoUnhandledErrorResult
func PatchVMKCTL(filename string) {

	// Read the file
	contents := loadFile(filename)

	// Replace applesmc with variable always found on ESXi
	var APPLESMC = []byte("applesmc")
	var VMKERNEL = []byte("vmkernel")

	// Find and replace string
	offset := bytes.Index(contents, APPLESMC)
	before := string(contents[offset : offset+8])
	copy(contents[offset:offset+8], VMKERNEL)

	after := string(contents[offset : offset+8])
	fmt.Printf("Patching %s -> %s\n", before, after)

	// Flush to disk
	saveFile(filename, contents)
}
