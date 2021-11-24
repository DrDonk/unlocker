// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package vmwpatch

import (
	"bytes"
	"fmt"
)

//goland:noinspection GoUnhandledErrorResult
func PatchVMKCTL(filename string) {

	// MMap the file
	contents := checkFile(filename)

	// Replace applesmc with variable always found on ESXi
	var APPLESMC = []byte("applesmc")
	var VMKERNEL = []byte("vmkernel")

	// Find and replace string
	offset := bytes.Index(contents, APPLESMC)
	before := string(contents[offset : offset+8])
	copy(contents[offset:offset+8], VMKERNEL)
	err := contents.Flush()
	if err != nil {
		panic(err)
	}
	after := string(contents[offset : offset+8])
	println(fmt.Sprintf("Patching %s -> %s", before, after))

}
