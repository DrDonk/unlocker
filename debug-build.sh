#!/usr/bin/env bash
#set -x
# Read current version
VERSION=$(<VERSION)

echo Building debug executables - $VERSION
echo package vmwpatch > ./vmwpatch/version.go
echo const VERSION = \"$VERSION\" >> ./vmwpatch/version.go

mkdir -p ./build/iso
mkdir -p ./build/linux
mkdir -p ./build/macos
mkdir -p ./build/windows

pushd ./commands/check
echo "Building check"
go-winres make --arch amd64 --product-version $VERSION --file-version $VERSION
env GOOS=windows GOARCH=amd64 go build -o ../../build/windows/check.exe
env GOOS=linux GOARCH=amd64 go build -o ../../build/linux/check
env GOOS=darwin GOARCH=amd64 go build -o ../../build/macos/check
rm rsrc_windows_amd64.syso
popd

pushd ./commands/relock
echo "Building relock"
go-winres make --arch amd64 --product-version $VERSION --file-version $VERSION
env GOOS=windows GOARCH=amd64 go build -o ../../build/windows/relock.exe
env GOOS=linux GOARCH=amd64 go build -o ../../build/linux/relock
env GOOS=darwin GOARCH=amd64 go build -o ../../build/macos/relock
rm rsrc_windows_amd64.syso
popd

pushd ./commands/unlock
echo "Building unlock"
go-winres make --arch amd64 --product-version $VERSION --file-version $VERSION
env GOOS=windows GOARCH=amd64 go build -o ../../build/windows/unlock.exe
env GOOS=linux GOARCH=amd64 go build -o ../../build/linux/unlock
env GOOS=darwin GOARCH=amd64 go build -o ../../build/macos/unlock
rm rsrc_windows_amd64.syso
popd

pushd ./commands/dumpsmc
echo "Building dumpsmc"
go-winres make --arch amd64 --product-version $VERSION --file-version $VERSION
env GOOS=windows GOARCH=amd64 go build -o ../../build/windows/dumpsmc.exe
env GOOS=linux GOARCH=amd64 go build -o ../../build/linux/dumpsmc
env GOOS=darwin GOARCH=amd64 go build -o ../../build/macos/dumpsmc
rm rsrc_windows_amd64.syso
popd

pushd ./commands/patchgos
echo "Building patchgos"
go-winres make --arch amd64 --product-version $VERSION --file-version $VERSION
env GOOS=windows GOARCH=amd64 go build -o ../../build/windows/patchgos.exe
env GOOS=linux GOARCH=amd64 go build -o ../../build/linux/patchgos
env GOOS=darwin GOARCH=amd64 go build -o ../../build/macos/patchgos
rm rsrc_windows_amd64.syso
popd

pushd ./commands/patchsmc
echo "Building patchsmc"
go-winres make --arch amd64 --product-version $VERSION --file-version $VERSION
env GOOS=windows GOARCH=amd64 go build -o ../../build/windows/patchsmc.exe
env GOOS=linux GOARCH=amd64 go build -o ../../build/linux/patchsmc
env GOOS=darwin GOARCH=amd64 go build -o ../../build/macos/patchsmc
rm rsrc_windows_amd64.syso
popd

pushd ./commands/patchvmkctl
echo "Building patchvmkctl"
go-winres make --arch amd64 --product-version $VERSION --file-version $VERSION
env GOOS=windows GOARCH=amd64 go build -o ../../build/windows/patchvmkctl.exe
env GOOS=linux GOARCH=amd64 go build -o ../../build/linux/patchvmkctl
env GOOS=darwin GOARCH=amd64 go build -o ../../build/macos/patchvmkctl
rm rsrc_windows_amd64.syso
popd

pushd ./commands/hostcaps
echo "Building hostcaps"
go-winres make --arch amd64 --product-version $VERSION --file-version $VERSION
env GOOS=windows GOARCH=amd64 go build -o ../../build/windows/hostcaps.exe
env GOOS=linux GOARCH=amd64 go build -o ../../build/linux/hostcaps
env GOOS=darwin GOARCH=amd64 go build -o ../../build/macos/hostcaps
rm rsrc_windows_amd64.syso
popd

cp -v LICENSE ./build
cp -v *.md ./build
cp -vr ./cpuid/* ./build
cp -vr ./iso ./build
