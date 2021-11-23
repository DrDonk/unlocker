// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"golocker/vmwpatch"
)

func main() {
	// Print titles
	println("PatchSMC")
	println("--------")
	println("© 2014-2021 David Parsons & Sam Bingner\n")
	vmwpatch.PatchSMC()
}
