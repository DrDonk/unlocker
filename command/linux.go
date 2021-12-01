// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

//go:build linux
// +build linux

package main

import (
	"bufio"
	"fmt"
	"github.com/djherbis/times"
	"github.com/mitchellh/go-ps"
	"golocker/vmwpatch"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type VMwareInfo struct {
	BuildNumber    string
	ProductVersion string
	InstallDir     string
	InstallDir64   string
	Workstation    string
	Player         string
	KVM            string
	REST           string
	Tray           string
	AuthD          string
	HostD          string
	USBD           string
	VMXDefault     string
	VMXDebug       string
	VMXStats       string
	VMwareBase     string
	PathVMXDefault string
	PathVMXDebug   string
	PathVMXStats   string
	PathVMwareBase string
}

func amAdmin() bool {
	if os.Geteuid() == 0 {
		return true
	}
	return false
}

func backupExists(v *VMwareInfo) bool {
	currentFolder, _ := os.Getwd()
	backupFolder := filepath.Join(currentFolder, "backup", v.ProductVersion)
	if _, err := os.Stat(backupFolder); !os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

//goland:noinspection GoUnhandledErrorResult
func copyFile(src, dst string) (int64, error) {
	fmt.Printf("%s -> %s\n", src, dst)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)

	// Ensure timestamops are correct
	srcTimes, _ := times.Stat(src)
	_ = os.Chtimes(dst, srcTimes.AccessTime(), srcTimes.ModTime())

	return nBytes, err
}

func delFile(src, dst string) error {
	// Get file mode RW/RO
	fi, _ := os.Stat(dst)
	err := os.Chmod(dst, 666)
	if err != nil {
		return err
	}

	// Copy file back
	_, err = copyFile(src, dst)
	if err != nil {
		return err
	}

	// Remove the backup
	err = os.Remove(src)
	if err != nil {
		return err
	}

	// Restore file mode
	err = os.Chmod(dst, fi.Mode())
	if err != nil {
		return err
	}

	return nil
}

func printHelp() {
	fmt.Printf("usage: unlocker.exe <install | uninstall>\n")
	fmt.Printf("\tinstall - install patches\n")
	fmt.Printf("\tuninstall - uninstall patches")
}

func taskRunning(name string) int {
	pid := 0
	tasks, _ := ps.Processes()

	for i := range tasks {
		if tasks[i].Executable() == name {
			pid = tasks[i].Pid()
		}
	}
	return pid
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' {
			s = s[1:]
		}
		if i := len(s) - 1; s[i] == '"' {
			s = s[:i]
		}
	}
	return s
}

func vmwBackup(v *VMwareInfo) {
	currentFolder, _ := os.Getwd()
	backupFolder := filepath.Join(currentFolder, "backup", v.ProductVersion)
	err := os.MkdirAll(backupFolder, os.ModePerm)
	if err != nil {
		panic(err)
	}
	_, err = copyFile(v.PathVMwareBase, filepath.Join(backupFolder, v.VMwareBase))
	if err != nil {
		panic(err)
	}
	_, err = copyFile(v.PathVMXDefault, filepath.Join(backupFolder, v.VMXDefault))
	if err != nil {
		panic(err)
	}
	_, err = copyFile(v.PathVMXDebug, filepath.Join(backupFolder, v.VMXDebug))
	if err != nil {
		panic(err)
	}
	_, err = copyFile(v.PathVMXStats, filepath.Join(backupFolder, v.VMXStats))
	if err != nil {
		panic(err)
	}
	return
}

func vmwRestore(v *VMwareInfo) {
	currentFolder, _ := os.Getwd()
	backupFolder := filepath.Join(currentFolder, "backup", v.ProductVersion)
	err := delFile(filepath.Join(backupFolder, v.VMwareBase), v.PathVMwareBase)
	if err != nil {
		panic(err)
	}
	err = delFile(filepath.Join(backupFolder, v.VMXDefault), v.PathVMXDefault)
	if err != nil {
		panic(err)
	}
	err = delFile(filepath.Join(backupFolder, v.VMXDebug), v.PathVMXDebug)
	if err != nil {
		panic(err)
	}
	err = delFile(filepath.Join(backupFolder, v.VMXStats), v.PathVMXStats)
	if err != nil {
		panic(err)
	}

	err = os.RemoveAll(backupFolder)
	return
}

//goland:noinspection GoUnhandledErrorResult
func vmwInfo() *VMwareInfo {
	v := &VMwareInfo{}

	// Store known service names
	// Not used on Linux
	v.AuthD = ""
	v.HostD = ""
	v.USBD = ""

	// Access /etc/vmware/config for version, build and installation path
	file, err := os.Open("/etc/vmware/config")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	config := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
					value = trimQuotes(value)
				}
				config[key] = value
			}
		}
	}

	// Basic product settings
	v.ProductVersion = config["product.version"]
	v.BuildNumber = config["product.buildNumber"]
	v.InstallDir = config["libdir"]

	// Construct needed filenames from reg settings
	v.InstallDir64 = ""
	v.Player = "vmplayer"
	v.Workstation = "vmware"
	v.KVM = "vmware-kvm"
	v.REST = "vmrest"
	v.Tray = "vmware-tray"
	v.VMXDefault = "vmware-vmx"
	v.VMXDebug = "vmware-vmx-debug"
	v.VMXStats = "vmware-vmx-stats"
	v.VMwareBase = "libvmwarebase.so"
	v.PathVMXDefault = filepath.Join(v.InstallDir, "bin", "vmware-vmx")
	v.PathVMXDebug = filepath.Join(v.InstallDir, "bin", "vmware-vmx-debug")
	v.PathVMXStats = filepath.Join(v.InstallDir, "bin", "vmware-vmx-stats")
	v.PathVMwareBase = filepath.Join(v.InstallDir, "lib", "libvmwarebase.so", "libvmwarebase.so")
	return v
}

func vmwRunning(v *VMwareInfo) bool {
	if taskRunning(v.Workstation) != 0 {
		fmt.Printf("VMware Workstation is running\n")
		return true
	}
	if taskRunning(v.Player) != 0 {
		fmt.Printf("VMware Player is running\n")
		return true
	}
	if taskRunning(v.KVM) != 0 {
		fmt.Printf("VMware KVM is running\n")
		return true
	}
	if taskRunning(v.REST) != 0 {
		fmt.Printf("VMware REST API is running\n")
		return true
	}
	if taskRunning(v.VMXDefault) != 0 {
		fmt.Printf("VMware VM (vmware-vmx) is running\n")
		return true
	}
	if taskRunning(v.VMXDebug) != 0 {
		fmt.Printf("VMware VM (vmware-vmx-debug) is running\n")
		return true
	}
	if taskRunning(v.VMXStats) != 0 {
		fmt.Printf("VMware VM (vmware-vmx-stats) is running\n")
		return true
	}
	return false
}

func main() {
	// Titles
	fmt.Printf("Unlocker %s for VMware Workstation/Player\n", vmwpatch.VERSION)
	fmt.Printf("============================================\n")
	fmt.Printf("%s \n\n", vmwpatch.COPYRIGHT)

	// Simple arg parser
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	var install bool
	switch os.Args[1] {
	case "install":
		install = true
	case "uninstall":
		install = false
	default:
		printHelp()
		return
	}

	// Check admin rights
	if !amAdmin() {
		fmt.Printf("Run as Administrator/root\n")
		return
	}

	// Get VMware product details from registry and file system
	v := vmwInfo()
	fmt.Printf("VMware is installed at: %s\n", v.InstallDir)
	fmt.Printf("VMware version %s\n\n", v.ProductVersion)

	// Check no VMs running
	if vmwRunning(v) {
		fmt.Printf("Aborting patching!\n")
		return
	}

	// Abort if installing and backup is present
	if install {
		fmt.Printf("Installing unlocker\n")
		if backupExists(v) {
			fmt.Printf("Aborting install as backup folder already exists!\n")
		}
	}

	// Abort if uninstalling and backup is missing
	if !install {
		fmt.Printf("Uninstalling unlocker\n")
		if !backupExists(v) {
			fmt.Printf("Aborting uninstall as backup folder does not exist!\n")
		}
	}

	if install {
		// Backup files
		fmt.Printf("\nBacking up files...\n")
		vmwBackup(v)

		// Patch files
		fmt.Printf("\nPatching...\n")
		vmwpatch.PatchSMC(v.PathVMXDefault)
		fmt.Printf("\n")
		vmwpatch.PatchSMC(v.PathVMXDebug)
		fmt.Printf("\n")
		vmwpatch.PatchSMC(v.PathVMXStats)
		fmt.Printf("\n")
		vmwpatch.PatchGOS(v.PathVMwareBase)
		fmt.Printf("\n")

		// Copy iso ISOs
		fmt.Printf("\nCopying VMware Tools...\n")
		_, _ = copyFile("./iso/darwinPre15.iso", filepath.Join(v.InstallDir, "darwinPre15.iso"))
		_, _ = copyFile("./iso/darwin.iso", filepath.Join(v.InstallDir, "darwin.iso"))
	} else {
		// Restore files
		fmt.Printf("\nRestoring files...\n")
		vmwRestore(v)

		// Removing iso ISOs
		fmt.Printf("\nRemoving VMware Tools...\n")
		isoPath := filepath.Join(v.InstallDir, "darwinPre15.iso")
		fmt.Printf("%s\n", isoPath)
		_ = os.Remove(isoPath)
		isoPath = filepath.Join(v.InstallDir, "darwin.iso")
		fmt.Printf("%s\n", isoPath)
		_ = os.Remove(isoPath)
	}

	fmt.Printf("\nPress any key to finish...")
	//goland:noinspection GoUnhandledErrorResult
	fmt.Scanln()
	return
}
