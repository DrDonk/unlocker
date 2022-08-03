# macOS Unlocker V4 for VMware Workstation

## 1. Introduction
Unlocker 4 is designed for VMware Workstation and Player 16 and has been tested against:

* Workstation Pro 16 on Windows and Linux
* Workstation Player 16 on Windows and Linux

The Unlocker enables certain flags and data tables that are required to see the macOS type when setting
the guest OS type, and modify the implmentation of the virtual SMC controller device. These capabiltiites are normally 
exposed in Fusion and ESXi when running on Apple hardware.

The patch code carries out the following modifications dependent on the product being patched:

* Patch vmware-vmx and derivatives to allow macOS to boot
* Patch vmwarebase.dll or libvmwarebase.so to allow Apple guests to be selected during VM creation
* Provide a copy of the macOS VMware Tools for the guest

It is important to understand that the Unlocker cannot add any new capabilities to VMware Workstation and Player
but enables support for macOS that is disabled in the VMware products that do not run on Apple hardware.

The Unlocker cannot:

* add support for new versions of macOS
* add paravirtualized Apple GPU support 
* add older (non-Ryzen) AMD CPU support

or any other features that are not already in the VMware compiled code. 

## 2. Running the patcher
### 2.1 Download Release
The code is written in Go and has no pre-requisites and should run directly from the release zip download.

* Download a binary release from https://github.com/DrDonk/unlocker/releases
* Unzip the archive to extract the executables and VMware Tools images
* Navigate to the folder with the extracted files

### 2.2 Windows
On Windows you can double-click the executables to launch them from the "windows" folder. You may be prompted to run 
with Administrator rights which are required to patch VMware.

* unlock.exe - apply patches to VMware
* relock.exe - remove patches from VMware
* check.exe  - check the patch status of VMware 

### 2.3 Linux
On Linux navigate to the folder where the unlocker has been unpacked, and run the required command from the 
"linux" folder. You will need to be either root or use sudo to run the scripts in the terminal.

* unlock - apply patches to VMware
* relock - remove patches from VMware
* check  - check the patch status of your VMware installation

## 3. Upgrading VMware
If a new VMware version is installed the patches will be lost as the update will overwrite the patched files.
You will need to re-run the unlock program to patch the newly installed VMware files.

## 4. Compiling from source code
The release ZIP file contains pre-compiled executables but if you want to compile the code please follow 
these instructions.

The Go language allows cross-compilation for multiple platforms and architectures. It is possible to
compile the unlocker on Windows, Linux and macOS for VMware Workstation on Linux and Windows.

### 4.1 Building on Windows
To build from source, do the following:
1. Install go and git using either [scoop](https://scoop.sh) or [chocolatey](https://chocolatey.org)
2. Open a Windows command prompt
3. `go install github.com/tc-hib/go-winres@latest`
4. `git clone https://github.com/DrDonk/unlocker.git`
5. `cd unlocker`
6. `build.cmd x.y.z` where x.y.z is the desiered Unlocker version number; e.g. 4.2.1.
7. `cd dist/windows`
8. Run the binaries as needed.

### 4.2 Building on Linux
To build from source, do the following:
1. Install go from your distro's repos
2. `go install github.com/tc-hib/go-winres@latest`
3. `export PATH="$PATH:$HOME/go/bin"`
4. `git clone https://github.com/DrDonk/unlocker.git`
5. `cd unlocker`
6. `sh build.sh x.y.z` where x.y.z is the desired Unlocker version number; e.g. 4.2.1.
7. `cd dist/linux`
8. Run the binaries as needed. 

### 4.3 Build a specific release
If you want to build from the source of a specific release and not of the latest commit, 
add `--branch <tag_name>` to the `git clone` command. For example, if you want to build from
the v4.2.1 tag, type in:

`git clone --branch v4.2.1 https://github.com/DrDonk/unlocker.git`.

and then build the code on your platform.

## 5. FAQS
### 5.1 False positive from antivirus software**
Some Windows antivirus programs mistakenly flag the unlocker as having a virus. This is a false positive
and it is mainly seen with Windows Defender. The Go Langauge website documents the problem.

https://go.dev/doc/faq#virus

If you have any doubts you can compile the source code yourself.

### 5.2 AMD CPUs
A patched macOS AMD kernel must be used to run on older AMD systems, but there is a workaround if you have a modern
AMD Ryzen CPU. The unlocker cannot patch this but we can recommend settings for the VMX file that allows macOS to
run on recent AMD CPUs. The tests are being recorded in this issue, and it would be useful if more can report
success or failures in that issue.

https://github.com/DrDonk/unlocker/issues/33

1. Read this KB article to learn how to edit a guest's VMX file safely https://kb.vmware.com/s/article/2057902
2. Add the following lines were to the VMX file:
```
cpuid.0.eax = "0000:0000:0000:0000:0000:0000:0000:1011"
cpuid.0.ebx = "0111:0101:0110:1110:0110:0101:0100:0111"
cpuid.0.ecx = "0110:1100:0110:0101:0111:0100:0110:1110"
cpuid.0.edx = "0100:1001:0110:0101:0110:1110:0110:1001"
cpuid.1.eax = "0000:0000:0000:0001:0000:0110:0111:0001"
cpuid.1.ebx = "0000:0010:0000:0001:0000:1000:0000:0000"
cpuid.1.ecx = "1000:0010:1001:1000:0010:0010:0000:0011"
cpuid.1.edx = "0000:0111:1000:1011:1111:1011:1111:1111"
vhv.enable = "FALSE"
vpmc.enable = "FALSE"
```
3. Make sure there are no duplicate lines in the VMX file or the guest will not start and a dictionary error will
   be displayed by VMware.
4. You can now install and run macOS as a guest.

## 6. VMware Downloads
These URLs will link to the latest versions of VMware's hosted products:

* VMware Fusion https://vmware.com/go/getfusion
* VMware Workstation for Windows https://www.vmware.com/go/getworkstation-win
* VMware Workstation for Linux https://www.vmware.com/go/getworkstation-linux
* VMware Player for Windows https://www.vmware.com/go/getplayer-win
* VMware Player for Linux https://www.vmware.com/go/getplayer-linux

## 7. VMware Tools
The Unlocker provides the VMware tools ISO images. There can be newer releases available which can be downloaded
from these URLs if the script has not yet been updated:

* Mac OS X 10.5 - 10.10 https://customerconnect.vmware.com/en/downloads/details?downloadGroup=VMTOOLS10012&productId=491
* macOS 10.11+ https://customerconnect.vmware.com/downloads/info/slug/datacenter_cloud_infrastructure/vmware_tools/12_x

_These URLs require a VMware login to download._

Version 16 of Workstation Pro recognises the darwin.iso files and the tools can be installed in the usual way by
using the "Install VMware Tools" menu item. The Player version does not automatically pick up the ISO images and so
the ISO must be maually attached to the VM via the guest's settings.

## 8. Thanks
Thanks to Zenith432 for originally building the C++ Unlocker and Mac Son of Knife
(MSoK) for all the testing and support.

Thanks also to Sam B for finding the solution for ESXi 6 and helping me with
debugging expertise. Sam also wrote the code for patching ESXi ELF files and
modified the Unlocker code to run on Python 3 in the ESXi 6.5 environment.

(c) 2011-2022 David Parsons
