# Troubleshooting Notes

## macOS Ventura Guests
There is a known issue that Ventura guests boot loop in VMware Workstation. There are different reasons for this
dependent on the CPU you are using on the host machine:

1. All CPUS - the Intel E1000E NIC has been dropped from Ventura and the default used by VMware is that NIC. It is
   a known issue and can be fixed by editing the guest's VMX file and ensuring the etehrnet device is set to:
   ``

Please do not upograde to Ventura at this time.

| CPU                | macOS         | Notes                          |
|--------------------|---------------|--------------------------------|
| Intel pre-Haswell  | macOS 10.6-12 | ESXi use 7.0 U2 hardware       |
| Intel post-Haswell | macOS 10.6-13 | ESXi use 7.0 U2 hardware       |
| AMD Ryzen          | macOS 11-12   | These are the only ones tested |

## 8. Important News

VMware will stop supporting new macOS versions as guests starting with Fusion 13 and the next version of ESXi.

This is documented at the VMware KB for Fusion:
https://kb.vmware.com/s/article/88697

And for ESXi:
https://kb.vmware.com/s/article/88698

This now seems to have affected ESXi 8, but Fusion 13 supports macos from 10.6 to 13 (Ventura).


## False positive from anti-virus software

Some Windows AV programs mistakenly flag the unlocker as having a virus. This is a false positive and it mainly is seen with Windows Defender. The Go Langauge web site documents the problem.

https://go.dev/doc/faq#virus

Always download from the releases section here, and if you have any doubts you can compile the source code and will see the false positive from Defender.


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


## Set a specific macOS Guest resolution

Installing the VMWare Tools should allow different video modes to be selected. If you have already installed them and
still does not change resolution you can try this. Open Terminal and run:

`sudo /Library/Application Support/VMware Tools/vmware-resolutionSet <width> <height>`

where width and height are the pixel size you want. For example to get 1440x900:

`sudo /Library/Application Support/VMware Tools/vmware-resolutionSet 1440 900`

### AMD CPUs
A patched macOS AMD kernel must be used to run on older AMD systems, but there is a workaround if you have a modern
AMD Ryzen CPU. The unlocker cannot patch this but we can recommend settings for the VMX file that allows macOS to
run on recent AMD CPUs. The tests are being recorded in this issue, and it would be useful if more can report
success or failures in that issue.

https://github.com/DrDonk/unlocker/issues/33

You must have Hyper-V disabled on Windows or VMware falls back to a mode called ULM. CPUID masking is not available
in ULM mode as it is pushed to Windows for the low level VMX/SVM operations. You can see if VMware is in ULM mode on
Hyper-V by searching the guest's vmware.log file for these 2 lines:

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
vvtd.enable = "FALSE"
```
3. Make sure there are no duplicate lines in the VMX file or the guest will not start and a dictionary error will
   be displayed by VMware.
4. You can now install and run macOS as a guest.
