// SPDX-FileCopyrightText: © 2023 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"github.com/drdonk/cpuid-go"
	"os"
)

func waitExit() {
	quiet := os.Getenv("UNLOCK_QUIET")
	if quiet != "1" {
		fmt.Println()
		fmt.Println("\nPress Enter key to continue...")

		_, _ = fmt.Scanln()
	}
}

func main() {
	fmt.Println("macOS Unlocker Host Capabilities")
	fmt.Println("================================")
	fmt.Println("")
	fmt.Println("Host CPU")
	fmt.Println("--------")
	fmt.Println("Name:", cpuid.CPU.BrandName)
	fmt.Println("Vendor ID:", cpuid.CPU.VendorID)
	fmt.Println("CPU Family:", cpuid.CPU.Family, "Model:", cpuid.CPU.Model, "Stepping:", cpuid.CPU.Stepping)
	fmt.Println("")
	fmt.Println("Host Report")
	fmt.Println("-----------")
	fmt.Println("AVX2 support:")
	fmt.Println("\tmacOS 13 Ventura: ", cpuid.CPU.Has(cpuid.AVX2))
	fmt.Println("Hypervisor mode: ")
	fmt.Println("\tHypervisor Present:", cpuid.CPU.Has(cpuid.HYPERVISOR))
	if cpuid.CPU.Has(cpuid.HYPERVISOR) {
		fmt.Println("\tHypervisor String:", cpuid.CPU.HypervisorString)
		fmt.Println("\tHypervisor ID:", cpuid.CPU.HypervisorID)
	}

	if cpuid.CPU.HypervisorString != "VMwareVMware" {
		fmt.Println("\tHypervisor found:")
		fmt.Println("\t\t- CPUID masking may not work correctly")
		fmt.Println("\t\t- Graphics speed may be affected")
	}
	waitExit()
	return
}
