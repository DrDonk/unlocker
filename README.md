# macOS Unlocker V4 for VMware Workstation

---
## Important News

VMware will stop supporting new macOS versions as guests starting with Fusion 13 and the next version of ESXi.

This is documented at the VMware KB for Fusion:
https://kb.vmware.com/s/article/88697

And for ESXi:
https://kb.vmware.com/s/article/88698

This will likely impact Workstation Pro/Player as well. Currently the exact details of what may change are not clear, but
it is safe to say we should look for alternative approaches to the unlocker.

---

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
* Optionally check the sha256 checksum matches that published in the release
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

### 2.3.1 Arch Linux and its derivatives
If you are using Arch Linux or a derivative of Arch, you can install and use this like any normal package [from the AUR](https://aur.archlinux.org/packages/vmware-unlocker-bin). Precompiled binaries and PKGBUILDs to build from the [source of the latest release](https://aur.archlinux.org/packages/vmware-unlocker) as well as [from the main branch](https://aur.archlinux.org/packages/vmware-unlocker-git) are provided. Included in all of the packages are some Pacman hooks to help unlock and relock VMware whenever there is an update to Unlocker, VMware Worksation Pro/Player, or any of VMware Workstation Pro/Player's components. The binary files are installed to the `/usr/bin` directory so they are already added to your `$PATH`.

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
### 5.1 False positive from antivirus software
Some Windows antivirus programs mistakenly flag the unlocker as having a virus. The code is compiled on macOS  
which is checked for viruses. These are false positives and if you want to check you can compile the code youself using
instruction in the readme file.

**Please do not log an issue about this!**

This is a known problem with Windows executables written in the Go languauge which is documented at the Go website.
[https://go.dev/doc/faq#virus]()

From the Go web site:

>Why does my virus-scanning software think my Go distribution or compiled binary is infected?
>
>This is a common occurrence, especially on Windows machines, and is almost always a false positive. 
>Commercial virus scanning programs are often confused by the structure of Go binaries, which they don't see as 
>often as those compiled from other languages.
>
>If you've just installed the Go distribution and the system reports it is infected, that's certainly a mistake. 
>To be really thorough, you can verify the download by comparing the checksum with those on the downloads page.
>
>In any case, if you believe the report is in error, please report a bug to the supplier of your virus scanner. 
>Maybe in time virus scanners can learn to understand Go programs.

### 5.2 AMD CPUs
A patched macOS AMD kernel must be used to run on older AMD systems, but there is a workaround if you have a modern
AMD Ryzen CPU. The unlocker cannot patch this but we can recommend settings for the VMX file that allows macOS to
run on recent AMD CPUs. The tests are being recorded in this issue, and it would be useful if more can report
success or failures in that issue.

https://github.com/DrDonk/unlocker/issues/33

You must have Hyper-V disabled on Windows or VMware falls back to a mode called ULM. CPUID masking is not available in ULM mode as it is pushed to Windows for the low level VMX/SVM operations. You can see if VMware is in ULM mode on Hyper_V by searching the guest's vmware.log file for these 2 lines:

`vmx IOPL_Init: Hyper-V detected by CPUID`

`Monitor Mode: ULM`

Here is a link to a Microsoft artcile on disabling Hyper-V.

https://docs.microsoft.com/en-us/troubleshoot/windows-client/application-management/virtualization-apps-not-work-with-hyper-v

_Editing the VMX file_
1. Read this KB article to learn how to edit a guest's VMX file safely https://kb.vmware.com/s/article/2057902
2. Add the following lines to the VMX file:
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
4. You can now install and run macOS as a gu

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
