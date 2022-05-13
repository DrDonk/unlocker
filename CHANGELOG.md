# Changes

All dates are UK format.

## 13/05/22 4.2.0
_drdonk:_
* Split unlocker.exe into 2 new files unlock.exe and relock.exe so can be directly executed without having to pass
a parameter to install or uninstall
* Added Windows manifest to allow the executables to directly run as Administrator without using re-spawn code. 
This seems to help with false AV detection especially Windows Defencder.


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
