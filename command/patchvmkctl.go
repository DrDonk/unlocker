package main

import (
	"golocker/vmwpatch"
)

func main() {
	// Print titles
	println("PatchVMKCTL")
	println("-----------")
	println("© 2014-2021 David Parsons\n")
	vmwpatch.PatchVMKCTL()
}
