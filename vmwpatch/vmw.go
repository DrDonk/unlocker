// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package vmwpatch

import (
	"fmt"
	"github.com/djherbis/times"
	"github.com/mitchellh/go-ps"
	"io"
	"os"
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
	PathISOmacOS   string
	PathISOMacOSX  string
	BackDir        string
	BackVMXDefault string
	BackVMXDebug   string
	BackVMXStats   string
	BackVMwareBase string
}

func CheckStatus(v *VMwareInfo) int {
	// TODO: Find a better way to check combinations
	var status = -1
	vmxBase, _ := IsGOSPatched(v.PathVMwareBase)
	vmxDefault, _ := IsSMCPatched(v.PathVMXDefault)
	vmxDebug, _ := IsSMCPatched(v.PathVMXDebug)
	if v.PathVMXStats != "" {
		vmxStats, _ := IsSMCPatched(v.PathVMXStats)
		status = vmxBase + vmxDefault + vmxDebug + vmxStats
		switch status {
		case 0:
			status = 0
		case 4:
			status = 1
		default:
			status = 2
		}
	} else {
		status = vmxBase + vmxDefault + vmxDebug
		switch status {
		case 0:
			status = 0
		case 3:
			status = 1
		default:
			status = 2
		}
	}
	return status
}

func Backup(v *VMwareInfo) {
	err := os.MkdirAll(v.BackDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	_, err = CopyFile(v.PathVMwareBase, v.BackVMwareBase)
	if err != nil {
		panic(err)
	}
	_, err = CopyFile(v.PathVMXDefault, v.BackVMXDefault)
	if err != nil {
		panic(err)
	}
	_, err = CopyFile(v.PathVMXDebug, v.BackVMXDebug)
	if err != nil {
		panic(err)
	}
	_, err = CopyFile(v.PathVMXStats, v.BackVMXStats)
	if err != nil {
		panic(err)
	}
	return
}

func Restore(v *VMwareInfo) {
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
	err = DelFile(v.BackVMXStats, v.PathVMXStats)
	if err != nil {
		panic(err)
	}

	err = os.RemoveAll(v.BackDir)
	return
}

func BackupExists(v *VMwareInfo) bool {
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
	nBytes, err := io.Copy(destination, source)

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
	pid := 0
	tasks, _ := ps.Processes()

	for i := range tasks {
		if tasks[i].Executable() == name {
			pid = tasks[i].Pid()
		}
	}
	return pid
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
	return
}
