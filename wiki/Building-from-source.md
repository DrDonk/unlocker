## Compiling from source code
The release ZIP file contains pre-compiled executables but if you want to compile the code please follow 
these instructions.

The Go language allows cross-compilation for multiple platforms and architectures. It is possible to
compile the unlocker on Windows, Linux and macOS for VMware Workstation on Linux and Windows.

### Building on Windows
To build from source, do the following:
1. Install go and git using either [scoop](https://scoop.sh) or [chocolatey](https://chocolatey.org)
2. Open a Windows command prompt
3. `go install github.com/tc-hib/go-winres@latest`
4. `git clone https://github.com/DrDonk/unlocker.git`
5. `cd unlocker`
6. `build.cmd x.y.z` where x.y.z is the desiered Unlocker version number; e.g. 4.2.1.
7. `cd build/windows`
8. Run the binaries as needed.

### Building on Linux
To build from source, do the following:
1. Install go from your distro's repos
2. `go install github.com/tc-hib/go-winres@latest`
3. `export PATH="$PATH:$HOME/go/bin"`
4. `git clone https://github.com/DrDonk/unlocker.git`
5. `cd unlocker`
6. `sh build.sh x.y.z` where x.y.z is the desired Unlocker version number; e.g. 4.2.1.
7. `cd build/linux`
8. Run the binaries as needed. 

### Build a specific release
If you want to build from the source of a specific release and not of the latest commit, 
add `--branch <tag_name>` to the `git clone` command. For example, if you want to build from
the v4.2.1 tag, type in:

`git clone --branch v4.2.1 https://github.com/DrDonk/unlocker.git`.

and then build the code on your platform.
