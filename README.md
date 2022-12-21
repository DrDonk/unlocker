# macOS Unlocker V4 for VMware Workstation
***

Please checkout the wiki. https://github.com/DrDonk/unlocker/wiki
***
 
## 1. Introduction
Unlocker 4 is designed for VMware Workstation and Player and has been tested against:

* Workstation Pro 16/17 on Windows and Linux
* Workstation Player 16/17 on Windows and Linux

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

## 2. Running the Unlocker
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

## 2.4 Upgrading VMware
If a new VMware version is installed the patches will be lost as the update will overwrite the patched files.
You will need to re-run the unlock program to patch the newly installed VMware files. This does not apply to the Arch Linux AUR packages due to the provided ALPM hooks.

## 2.5 VMware Tools
The Unlocker provides the VMware tools ISO images.Version 16/17 of Workstation Pro recognises the darwin.iso files and the tools can be installed in the usual way by using the "Install VMware Tools" menu item. The Player version does not automatically pick up the ISO images and so
the ISO must be maually attached to the VM via the guest's settings.

## 3. VMware Downloads
These URLs will link to the latest versions of VMware's hosted products:

* VMware Fusion https://vmware.com/go/getfusion
* VMware Workstation for Windows https://www.vmware.com/go/getworkstation-win
* VMware Workstation for Linux https://www.vmware.com/go/getworkstation-linux
* VMware Player for Windows https://www.vmware.com/go/getplayer-win
* VMware Player for Linux https://www.vmware.com/go/getplayer-linux
* VMware Guest Tools https://vmware.com/go/tools

## 4. Thanks
Thanks to Zenith432 for originally building the C++ Unlocker and Mac Son of Knife
(MSoK) for all the testing and support.

Thanks also to Sam B for finding the solution for ESXi 6 and helping me with
debugging expertise. Sam also wrote the code for patching ESXi ELF files and
modified the Unlocker code to run on Python 3 in the ESXi 6.5 environment.

There have also been many new testers via GitHub, some who want to remain anoymous, but a big thanks to them.

(c) 2011-2022 David Parsons
