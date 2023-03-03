# Changes

All dates are UK DD/MM/YY format.


## 03/03/23 4.2.6
### This is the last release of Unlocker.

_drdonk:_
* Tidy up of code
* Removed template VMs
* Added macOS recovery VMDK utility

## 06/02/23 4.2.5
_drdonk_
* Added dumpsmc to release build as useful for diagnostic info
* SMC text labels for private and public key count incorrectly swapped
* Add wait for key press to hostcaps
* Updated cpuid binaries from http://github.com/drdonk/cpuid
* Build process updated to make it more automatedbuilding
* Added TAR Gzipped distribution file for Linux 

## 21/12/22 4.2.4
_drdonk_
* Changes to templates to disable webcam in VM
* Removed files from templates that are generated when VM is first opened to ensure uniqueness of VMs.
* New templates with good defaults for different versions of VMware and macOS guests
* Remove some duplicate messages in executables
* Removed a marker in patches as never used (KPPW/KPST)
* Wiki now source of more dynamic informtion such as Ventura supoort

_samhocevar_
* Check for vmware-vmx-stats file exists before trying to patch

_SimPilotAdamT_
* Clarification on using AUR with Unlocker

## 21/09/22 4.2.3
_drdonk_
* VMware Tools 12.1.0 - last VMware tools for macOS to be released
* Template VMs for Intel and AMD CPUs with sensible defaults and preformatted HFS+J virtual disk

## 30/07/22 4.2.2
_drdonk_
* Re-worked README.md file
* Updated darwin.iso to version 12.0.6
* relock restarts Windows services unpatching fails
* Initial changes for 22H2 tech previews. New Windows service for auto-starting VMs needs to be stopped and started
* Ensure using bash and fix dist folder creation in the source code build scripts fro Linux/macOS

## 03/06/22 4.2.1
_drdonk_:
* Missing message when removing ISOs in relock
* Updated darwin.iso to version 12.0.5

## 13/05/22 4.2.0
_drdonk:_
* Split unlocker.exe into 2 new files unlock.exe and relock.exe so can be directly executed on Windows without 
having to pass a parameter to install or uninstall
* Added Windows manifest to allow the executables to directly run as Administrator without using re-spawn code 
* Minimize false AV detection especially Windows Defender
* Added a new command to check the status of the patching called check.exe
* Fix to restarting vmware-tray.exe on Windows which had the wrong working directory

## 03/03/22 4.1.3
_drdonk:_
* Update VMware Tools to version 12.0.0

## 15/02/22 4.1.2

_drdonk:_
* Fix stopping/starting services on Windows due to service dependencies

## 26/01/22 4.1.1

_drdonk:_
* Fix copying of VMware tools ISO to the VMware installation folder
* Added a new troubleshooting page
* Updated copyright to 2022

## 09/12/21 4.1.0

Changed version to 4.1.0 due to the large number of changes to code submitted by boostchicken. 
Thanks to boostchicken for the fixes and improvements to the code.

_boostchicken:_
* Refactor to be more go like
* Fix vmware-shell-ext-thunker.exe preventing patching on windows
* Fixed header offset bug and binpack allocs
* Cleaned up logic regarding restoring backups or not
* Elevate permissions on Windows automatically

_drdonk:_
* Key press message changed to correct input of enter key
* Build macOS dumpsmc executable

## 06/12/21 4.0.0
* Initial release
