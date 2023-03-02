# OC4VM  - macOS Recovery VMDK Maker

## Introduction
This is a utility  to download the recovery image for macOS from Apple's servers and create a bootable 
VMDK file that can be used to start an Internet installtion of macOS.

It also allows SIP and other recovery only mode commands to be used on an existing macOS guest. 
This is needed due to a bug in VMware EFI implementation which stops the Recovery environment being 
booted directly by using the CMD-R combination.

## Pre-requisites

###  Python 3
The tool is written in Python 3 and so Python must be installed on the host computer. Python 3.10 was used for 
the development of the tool.

* Linux   - Use your distro package manager
* macOS   - `brew install python3`
* Windows - Recommended way is to install Python 3.10 from the Microsoft App Store, or use Chocolatey or Scoop.

### qemu-img

You will need to have qemu-img utility, from QEMU, on the path.

* Linux - These can be installed from Linux repos, for example Debian based distros:

    `sudo apt install -y qemu-utils`


* macOS - The [brew](https://brew.sh) package manager is the best way to install on macOS, 
but the full Qemu package must be installed:

    `brew install qemu`


* Windows - use [Chocolatey](https://chocolatey.org) or [Scoop](https://scoop.sh) to install the programs:

    `choco/scoop install qemu`
   

## Instructions
1. Unzip the archive maintaining the folder structure
2. Open a console/shell in the folder with the tool
3. Run the tool: `python3 recovery_vmdk.py`
4. The menu will be displayed and just select the macOS version you want using the number on the menu.
```
OC4VM macOS Recovery VMDK Maker
===============================
(c) David Parsons 2022-23

Create a VMware VMDK Recovery Image
1. Catalina
2. Big Sur
3. Monterey
4. Ventura
Input menu number: 
```
The tool will  download the BaseSystem.dmg for the macOS version you selected and convert it to a VMDK.

Aftr the tool has finished there will be 3 files present in the folder. For example if downloading Ventura
there will be:

1. ventura.dmg
2. ventura.chunklist
3. ventura.vmdk

The .vmdk file is the only one needed for VMware. The .dmg and .chunklist files are the orginal files downloaded 
from Apple and can be removed if not needed.

## Acknowledgements
This tool wraps several other great open source commands. Thanks to the authors of those tools.

* macrecovery.py - https://github.com/acidanthera/OpenCorePkg
* qemu - https://www.qemu.org
* qemu-img for Windows - https://cloudbase.it/qemu-img-windows