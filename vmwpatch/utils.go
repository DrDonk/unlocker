// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package vmwpatch

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func loadFile(filename string) (contents []byte) {

	// Print filename
	fmt.Printf("File: %s\n", filename)

	// Open the file
	contents, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	// Return the contents
	return contents
}

func saveFile(filename string, contents []byte) {
	err := os.WriteFile(filename, contents, 0444)
	if err != nil {
		panic(err)
	}
	return
}

func sha256File(contents []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(contents))
}
