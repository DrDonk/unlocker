### **At this time Ventura does not work please read [Ventura KB](https://github.com/DrDonk/unlocker/wiki/macOS-Ventura-Guests).**

A patched macOS AMD kernel must be used to run on older AMD systems, but there is a workaround if you have a modern
AMD Ryzen CPU. The unlocker cannot patch this but we can recommend settings for the VMX file that allows macOS to
run on recent AMD CPUs. 

You must have Hyper-V disabled on Windows or VMware falls back to a mode called ULM. CPUID masking is not available in ULM mode as it is pushed to Windows for the low level VMX/SVM operations. You can see if VMware is in ULM mode on Hyper-V by searching the guest's vmware.log file for these 2 lines:

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
4. You can now install and run macOS as a guest
