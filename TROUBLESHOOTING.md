# Troubleshooting Notes

## False positive from anti-virus software

Some Windows AV programs mistakenly flag the unlocker as having a virus. This is a false positive and it mainly is seen with Windows Defender. The Go Langauge web site documents the problem. 

https://go.dev/doc/faq#virus

Always download from the releases section here, and if you have any doubts you can compile the source code and will see the false positive from Defender.

## Set a specific macOS Guest resolutiom

Installing the VMWare Tools should allow different video modes. to be selected. If you have already installed them and still does not change resolution you can try this. Open Terminal and run:

`sudo /Library/Application Support/VMware Tools/vmware-resolutionSet <width> <height>`

where width and height are the pixel size you want. For example to get 1440x900:

`sudo /Library/Application Support/VMware Tools/vmware-resolutionSet 1440 900`
