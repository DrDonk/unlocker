# The Patches

## VMware Code

VMware virtualises the CPU but importantly it emulates a lot of hardware, USB, SCSI, SATA, NIC, audio, etc. It also emulates an Apple SMC controller (vSMC). The vSMC works as expected in VMware Fusion and ESXi when it is running on Apple hardware. We need to enable this code to be functional on all platforms not just Apple hardware. This is the fundamental patch needed to get things working on Windows, Linux and ESXi.

There are two ancillary patches just to make life easier. The first is to enable the Apple macOS settings in the GUI when creating or modifying a VM, which is for Windows and Linux. The second is to allow vCenter server to start macOS guests on an ESXi host, which requires a small patch on the ESXi server.

## Patch 1 - Virtual SMC Controller Data

**Note: Needed for Workstation, Player and ESXi**

The UI spawns a new executable that actually runs the guest and provides the hardware virtualisation interface from the host to the guest. It also contains all the emulation code for the virtual devices, including the vSMC, used by the guest.

* vmware-vmx - Linux & macOS
* vmware-vmx.exe - Windows
* vmx - ESXi

(There are also debug and stats variants used for troubleshooting; e.g. on ESXi vmx, vmx-debug and vmx-stats files in the /bin folder.)

The code for the vSMC is compiled into vmware-vmx on every platform and consists of the emulated device, mapped either to an IO port or a memory address, and a table of static SMC key and values. These keys and values are read directly from the executable's compiled data.

There are 2 separate SMC implementations which can be changed via the VMX file:

`smc.version = "0"` - create an I/O port mapped vSMC used in earlier Intel Macs

`smc.version = "1"` - create a memory mapped vSMC used in current Intel Macs

The vSMC tables were found simply by looking for known key names such as OSK0.

Here's small section of the vSMC table:

```                           AppleSMCHeaderV0                                XREF[2]:     FUN_004b3ab0:004b3b5e(*), 
                                                                                          0166f990(*)  
        016a6fa0 c0 6f 6a        addr       AppleSMCDataV0                                   = "YEK#"
                 01 00 00 
                 00 00
        016a6fa8 f2 00 00 00     ddw        F2h
        016a6fac f0 00 00 00     ddw        F0h
        016a6fb0 00 00 00        ??[16]
                 00 00 00 
                 00 00 00 
                             AppleSMCDataV0                                  XREF[1]:     016a6fa0(*)  
        016a6fc0 59 45 4b 23     char[4]    "YEK#"
        016a6fc4 04              db         4h
        016a6fc5 32 33 69 75     char[4]    "23iu"
        016a6fc9 80              db         80h
        016a6fca 00 00 00        db[6]
                 00 00 00
        016a6fd0 70 3f 4b        addr       FUN_004b3f70
                 00 00 00 
                 00 00
        016a6fd8 00 00 00        db[48]
                 00 00 00 
                 00 00 00 
        016a7008 53 4b 4c 2b     char[4]    "SKL+"
        016a700c 01              db         1h
        016a700d 67 61 6c 66     char[4]    "galf"
        016a7011 90              db         90h
        016a7012 00 00 00        db[6]
                 00 00 00
        016a7018 90 3e 4b        addr       FUN_004b3e90
                 00 00 00 
                 00 00
        016a7020 07 00 00        db[48]
                 00 00 00 
                 00 00 00 
```
A little bit of reverse engineering, by looking at the repetitive structure of the table, gave us these data structures as documented in the Python based ESXi unlocker code.

```
vSMC Header Structure
=====================
Offset  Length  Struct Type Description
----------------------------------------
0x00/00 0x08/08 Q      ptr  Offset to key table
0x08/08 0x04/4  I      int  Number of private keys
0x0C/12 0x04/4  I      int  Number of public keys

vSMC Key Data Structure
Offset  Length  Struct Type Description
----------------------------------------
0x00/00 0x04/04 4s     int  Key name (byte reversed e.g. #KEY is YEK#)
0x04/04 0x01/01 B      byte Length of returned data
0x05/05 0x04/04 4s     int  Data type (byte reversed e.g. ui32 is 23iu)
0x09/09 0x01/01 B      byte Flag R/W
0x0A/10 0x06/06 6x     byte Padding
0x10/16 0x08/08 Q      ptr  Internal VMware routine
0x18/24 0x30/48 48B    byte Data

The internal VMware routines point to 4 variants:
AppleSMCHandleDefault
AppleSMCHandleNTOK
AppleSMCHandleNumKeys
AppleSMCHandleOSK
```

All the data structures are fixed length (in bytes):

* header = 16
* row = 72
   * key = 24
   * data = 48
* table size = row * header_number_of_private_keys

You can dump the keys from the vmx executable using a program I created called dumpsmc program. Here is an extract:
```
DumpSMC 4.2.5
==============
© 2014-2023 David Parsons 

Filename: ./samples/esxi8/vmx

appleSMCTableV0 (smc.version = '0')
File Offset  : 0x0137a6e0
Keys Address : 0x0137c700
Public Keys  : 0x00f0/0240
Private Keys : 0x00f2/0242

Table Offset : 0x0137a700
Offset     Name Len Type Flag FuncPtr    Data
-------    ---- --- ---- ---- -------    ----
0x0137a700 #KEY 04  ui32 0x80 0x003a1c30 00000000
0x0137a748 +LKS 01  flag 0x90 0x003a1bf0 07
0x0137a790 ACEN 01  ui8  0x10 0x003a1bf0 00
0x0137a7d8 ACFP 01  flag 0x80 0x003a1bf0 01
….
0x0137e9f0 zDBG 01  ui8  0x50 0x003a1bf0 00
0x0137ea38 zSCI 01  ui8  0xd0 0x003a1bf0 00
0x0137ea80 OSK0 32  ch8* 0x80 0x003a1c50 0000000000000000000000000000000000000000000000000000000000000000
0x0137eac8 OSK1 32  ch8* 0x80 0x003a1c50 0000000000000000000000000000000000000000000000000000000000000000
```

When macOS requests a value from the SMC, VMware looks it up from this table and returns it to the guest. The one exception is that the static table of data does not contain the Apple haiku and vmware-vmx tries to read the OSK0/1 keys from the real host hardware. This will work a real Mac and return the required data from the real SMC by storing it at runtime in the table held in memory.

Linux and Windows do have drivers for real SMC via Linux kernel module and Windows device driver (it's in the Bootcamp download), but these drivers are never called from code not even when running Linux or Windows on real Apple hardware.

**So patch 1's goal is to statically add the haiku to the data table, and patch the function to read it directly from the static data and not try to call the host's SMC.**

From the analysis of the data tables and the functions it was found a patch could be quite simply written to override the calls to the host.

So the logic to patch the vSMC tables:

* Locate the vSMC table in the file by looking for the #KEY which is always the first entry (including real Macs)
* Loop through the table until we find OSK0 or OSK1 and add the haiku into the data area of the vSMC record.
* Replace the call to the real SMC (AppleSMCHandleOSK) with the default function (AppleSMCHandleDefault) for OSK0/1 pointer field.

There is one more patch needed for ELF type executables used by Linux and ESXi. ELF files have a table which shows where each function is being used and patches the in-memory address when loaded. As we have altered the OSK function pointers we need to find them in the RELA entries table and point them to the new functions.

Windows PE binaries are simpler than Linux/ESXi ELF binaries and do not need that fix up.

## Patch 2 - Guest OS Visibility

**Note: Workstation and Player only**

VMware supports many guest OS types and displays them in the UI when creating or modifying the a guest. This is driven by data tables which are held in a shared object loaded by the UI.

It's called by different names on the different platforms:

* libvmwarebase.dylib - macOS
* vmwarebase.dll - Windows
* libvmwarebase.so - Linux

So the interesting thing here is that the table exists on all platforms but Fusion displays Apple OS whereas Workstation does not. Reverse engineering of the Fusion file versus the Workstation files allowed the bit that exposes the guest to the UI to be found.

The GOSTable stores a lot of different information, and much of it seems to be in bit fields. It looks like it includes the fields that specify what the initial virtual hardware for the guest is when it is created. For example a guest should have SATA rather than SCSI virtual drives.

To patch this a pattern was found from the table to run a regular expression against the binary file to find all instances for macOS GOS entries. It was simpler to do it this way as the actual full structure of the table was not important, just finding the bit to expose the guest type to the UI was all that was needed.

Regexp pattern for GOS table Darwin entries:

`"\x10\x00\x00\x00[\x10|\x20]\x00\x00\x00[\x01|\x02]\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"`

Once we find the entries we loop though them in the file and set the exposed bit.

## Patch 3 - ESXi SMC Present

**Note: ESXi only**

The final patch is exclusively for ESXi when attached to a vCenter server. If the ESXi server is standalone then macOS can started from the ESXi web UI. However to start the VM from a vCenter server the ESXi host has to respond to an API request that returns whether the host has an SMC controller present. i.e. a real Apple machine. Of course we don't have one on non-Apple hardware so need to fix that by patching the API vCenter calls to find out if the SMC is present.

The vCenter and ESXi APIs are big and complex but it was possible to see which call was being used, via Wireshark and checking logs, and it's one that retrieves the host hardware capabilities. This can be done via the ESXi web API or directly on the ESXi host by logging in to the console or using SSH and running:

**`vim-cmd hostsvc/hosthardware`**

This will return output and you should see an smcPresent field whose value is returned by another shared object called libvmkctl.so. This shared object looks to see if the applesmc module has been loaded into the vmkernel.  The patch is a very simple switch from "applesmc" to "vmkernel" which is always present and as you can guess it's the name of the core of ESXi.

