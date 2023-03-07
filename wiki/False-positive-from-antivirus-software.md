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
