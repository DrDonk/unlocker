// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package vmwpatch

import (
	"fmt"
	"os"

	"github.com/djherbis/times"
	"github.com/mitchellh/go-ps"
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
	ShellExt       string
	VMXDefault     string
	VMXDebug       string
	VMXStats       string
	VMwareBase     string
	PathVMXDefault string
	PathVMXDebug   string
	PathVMXStats   string
	PathVMwareBase string
	PathISOmacOS   string
	PathISOMacOSX  string
	BackDir        string
	BackVMXDefault string
	BackVMXDebug   string
	BackVMXStats   string
	BackVMwareBase string
}

type PatchOperation struct {
	FileToPatch    string
	BackupLocation string
}

func (v *VMwareInfo) PatchFiles(gos chan *PatchOperation, smc chan *PatchOperation, done chan int) {
	err := os.MkdirAll(v.BackDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	gos <- &PatchOperation{FileToPatch: v.PathVMwareBase, BackupLocation: v.BackVMwareBase}
	smc <- &PatchOperation{FileToPatch: v.PathVMXDefault, BackupLocation: v.BackVMXDefault}
	smc <- &PatchOperation{FileToPatch: v.PathVMXDebug, BackupLocation: v.BackVMXDebug}
	_, err = os.Stat(v.PathVMXStats)
	if err == nil {
		smc <- &PatchOperation{FileToPatch: v.PathVMXStats, BackupLocation: v.BackVMXStats}
	}
	done <- 1
}

func (p *PatchOperation) Backup() bool {
	_, err := os.Stat(p.BackupLocation)
	if err != nil {
		_, err = CopyFile(p.FileToPatch, p.BackupLocation)
		if err != nil {
			panic(err)
		}
		return true
	}
	return false
}

func (v *VMwareInfo) Restore() {
	err := DelFile(v.BackVMwareBase, v.PathVMwareBase)
	if err != nil {
		panic(err)
	}
	err = DelFile(v.BackVMXDefault, v.PathVMXDefault)
	if err != nil {
		panic(err)
	}
	err = DelFile(v.BackVMXDebug, v.PathVMXDebug)
	if err != nil {
		panic(err)
	}
	_ = DelFile(v.BackVMXStats, v.PathVMXStats)
	_ = os.RemoveAll(v.BackDir)
	return
}

func (v *VMwareInfo) BackupExists() bool {
	if _, err := os.Stat(v.BackDir); !os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

//goland:noinspection GoUnhandledErrorResult
func CopyFile(src, dst string) (int64, error) {
	fmt.Printf("%s -> %s\n", src, dst)
	srcFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !srcFileStat.Mode().IsRegular() {
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
	var nBytes int64
	nBytes, err = destination.ReadFrom(source)
	// Ensure file mode and ownership is correct
	_ = os.Chmod(dst, srcFileStat.Mode())

	// Ensure timestamps are correct
	srcTimes, _ := times.Stat(src)
	_ = os.Chtimes(dst, srcTimes.AccessTime(), srcTimes.ModTime())
	if srcTimes.HasBirthTime() {
		_ = setCTime(dst, srcTimes.BirthTime())
	}

	return nBytes, err
}

func DelFile(src, dst string) error {
	// Get file mode RW/RO
	fi, _ := os.Stat(dst)
	err := os.Chmod(dst, 666)
	if err != nil {
		return err
	}

	// Copy file back
	_, err = CopyFile(src, dst)
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

func IsRunning(v *VMwareInfo) bool {
	if TaskRunning(v.Workstation) != 0 {
		fmt.Printf("VMware Workstation is running\n")
		return true
	}
	if TaskRunning(v.Player) != 0 {
		fmt.Printf("VMware Player is running\n")
		return true
	}
	if TaskRunning(v.KVM) != 0 {
		fmt.Printf("VMware KVM is running\n")
		return true
	}
	if TaskRunning(v.REST) != 0 {
		fmt.Printf("VMware REST API is running\n")
		return true
	}
	if TaskRunning(v.VMXDefault) != 0 {
		fmt.Printf("VMware VM (vmware-vmx) is running\n")
		return true
	}
	if TaskRunning(v.VMXDebug) != 0 {
		fmt.Printf("VMware VM (vmware-vmx-debug) is running\n")
		return true
	}
	if TaskRunning(v.VMXStats) != 0 {
		fmt.Printf("VMware VM (vmware-vmx-stats) is running\n")
		return true
	}
	return false
}

func TaskRunning(name string) int {
	tasks, _ := ps.Processes()

	for _, task := range tasks {
		if task.Executable() == name {
			return task.Pid()
		}
	}
	return 0
}

func WriteHashes(filename string, unpatched string, patched string) {
	shafilename := fmt.Sprintf("%s%s", filename, ".sha256")
	f, err := os.Create(shafilename)
	if err != nil {
		panic(err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%s\n%s", unpatched, patched))
	if err != nil {
		panic(err)
	}

	//goland:noinspection GoUnhandledErrorResult
	f.Sync()
}
