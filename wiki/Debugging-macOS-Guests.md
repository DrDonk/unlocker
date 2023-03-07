Sometimes it is neccessary to get debugging info from both VMware and the macOS guest to diagnose problems. Please follow these steps to add logging and debug information to a macOS guest. Please note this is best run with a single virtual drive holding macOS to make it easier to find the relevant boot files.

1. Shutdown the guest in VMware
2. Close VMware Fusion, Workstation or Player
3. Open the guest's VMX file in a text editor and add the following lines, making sure there are no duplicate lines:
```
answer.msg.serial.file.open = "Replace"
serial0.fileName = "serial.log"
serial0.fileType = "file"
serial0.present = "TRUE"
serial0.tryNoRxLoss = "FALSE"
serial0.yieldOnMsrRead = "TRUE"
vmx.buildType = "debug"
```


4. Power on to the firmware using the menus:
VM > Power > Power On to Firmware.

5. Select the EFI Shell option:

![Screenshot 2022-11-08 at 12 51 50](https://user-images.githubusercontent.com/869796/200572162-b5c1dfb1-a22a-4e82-919f-039ce9e4dfb6.jpg)

6. Use the `ls` command to find the boot.efi file. In this case it is in fs1:, but you may ahve to look in other partitions.
 
![Screenshot 2022-11-08 at 12 53 05](https://user-images.githubusercontent.com/869796/200572378-5fcfa232-bc18-4ff8-9440-bd833f20d8a2.jpg)

7. Run the booter using the following parameters:

![Screenshot 2022-11-08 at 12 54 39](https://user-images.githubusercontent.com/869796/200572595-a937fb8b-bd07-482e-8153-72ff2be9704b.jpg)

8. The guest should start booting and you will get this screen, but after that all output goes to serial.log.

![Screenshot 2022-11-08 at 12 54 48](https://user-images.githubusercontent.com/869796/200572757-60b30293-8180-4a74-b015-df474f0735ef.jpg)

After the test has been completed shut the guest down and zip up the serial.log and vmware.log files in the guest folder. They can then be attached to any issues or dicussions.

Please note that these errors are normal(?!) even when running macOS in Fusion on real Apple hardware.

```
ACPI Error: ACPI Error: [PSTA][PSTA] Namespace lookup failure, AE_NOT_FOUND Namespace lookup failure, AE_NOT_FOUND (20160930/psargs-463)
 (20160930/psargs-463)
[_STA] @00001 #002D:  [_STA] @00001 #002D:  
No Local Variables are initialized for method [_STA]
No Local Variables are initialized for method [_STA]
No Arguments are initialized for method [_STA]
No Arguments are initialized for method [_STA]
ACPI Error: ACPI Error: Method parse/execution failed Method parse/execution failed [\_SB.PCI0.P2P0.S1F0._STA] (Node ffffffa2292e7630)[\_SB.PCI0.P2P0.S1F0._STA] (Node ffffffa2292e7630), AE_NOT_FOUND, AE_NOT_FOUND (20160930/psparse-632)
 (20160930/psparse-632)
```
