package main

import (
	"golocker/vmwpatch"
)

func main() {
	// Print titles
	println("PatchGOS")
	println("-----------")
	println("© 2014-2021 David Parsons\n")
	vmwpatch.PatchGOS()
}
