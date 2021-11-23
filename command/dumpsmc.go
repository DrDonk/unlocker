package main

import (
	"golocker/vmwpatch"
)

func main() {
	// Print titles
	println("DumpSMC")
	println("-------")
	println("© 2014-2021 David Parsons\n")
	vmwpatch.DumpSMC()
}
