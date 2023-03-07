There are a known issues that Ventura guests boot loop in VMware Workstation. The situation after extensive testing is:

1. Intel pre-Haswell CPUs dropped by Apple in Ventura - no way to patch need to look at OCLP
2. AMD CPUs no longer work in Ventura possibly due to how CPUID leaf 7 is read - no way to patch need to look at OpenCore and other solutions
3. Intel Haswell+ CPUS, which are supported, make sure you set 

`ethernet0.virtualDev = "e1000e"`

to

`ethernet0.virtualDev = "vmxnet3"`



| CPU                | macOS         | Notes                          |
|--------------------|---------------|--------------------------------|
| Intel pre-Haswell  | macOS 10.6-12 | ESXi use 7.0 U2 hardware       |
| Intel post-Haswell | macOS 10.6-13 | ESXi use 7.0 U2 hardware       |
| AMD Ryzen          | macOS 11-12   | These are the only ones tested |

Please do not upgrade to Ventura at this time if you have AMD or Intel pre-Haswell CPUs.
