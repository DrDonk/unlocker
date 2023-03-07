# Tokamak macOS Recovery VMDK Maker

## Introduction
This is port of the Tokamak utilities for VMware Fusion adapted to run on Linux and Windows. 
It is designed to download the recovery image for macOS from Apple's servers and create a bootable 
VMDK file that can be used to start an Internet installtion of macOS.

It also allows SIP and other recovery only mode commands to be used on an existing macOS guest. 
This is needed due to a bug in VMware EFI implementation which stops the Recovery environment being 
booted directly by using the CMD-R combination.

## Pre-requisites

###  Python 3
The tool is written in Python 3 and so Python must be installed on the host computer. Python 3.10 was used for 
the development of the tool.

* macOS   - `brew install python3`
* Linux   - Use your distro package manager
* Windows - Recommended way is to install Python 3.10 from the Microosft App Store

### dmg2img & qemu-img

macOS and Linux users will need to have two utilities on the path:
1. dmg2img
2. qemu-img

These can be installed from Linux repos, for example Debian based distros:

`sudo apt install -y dmg2img`

`sudo apt install -y qemu-utils`

The brew package manager is the best way to install on macOS, but the full Qemu package must be installed:

`brew install dmg2img`

`brew install qemu`

Windows binaries are shipped as part of the zip file and so no additional downloads are required.

## Instructions
Download the zip file from: [recovery_vmdk.zip](https://github.com/DrDonk/unlocker/files/9992470/recovery_vmdk.zip)

1. Unzip the archive maintaining the folder strucutre
2. Open a console/shell in the folder with the tool
3. Run the tool: `python3 recovery_vmdk.py`
4. The menu will be displayed and just select the macOS version you want using the number on the menu.
```
Tokamak macOS Recovery VMDK Maker
=================================
(c) David Parsons 2022

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
* dmg2img - http://vu1tur.eu.org/tools
* qemu - https://www.qemu.org
* qemu-img for Windows - https://cloudbase.it/qemu-img-windows
