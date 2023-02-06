#!/usr/bin/env bash
#set -x

# Read current version
VERSION=$(<VERSION)

echo Building release executables - $VERSION

mkdir -p ./build/iso
mkdir -p ./build/linux
mkdir -p ./build/windows
mkdir -p ./build/templates

pushd ./commands/check
echo "Building check"
go-winres make --arch amd64 --product-version $VERSION --file-version $VERSION
env GOOS=windows GOARCH=amd64 go build -o ../../build/windows/check.exe
env GOOS=linux GOARCH=amd64 go build -o ../../build/linux/check
rm rsrc_windows_amd64.syso
popd

pushd ./commands/relock
echo "Building relock"
go-winres make --arch amd64 --product-version $VERSION --file-version $VERSION
env GOOS=windows GOARCH=amd64 go build -o ../../build/windows/relock.exe
env GOOS=linux GOARCH=amd64 go build -o ../../build/linux/relock
rm rsrc_windows_amd64.syso
popd

pushd ./commands/unlock
echo "Building unlock"
go-winres make --arch amd64 --product-version $VERSION --file-version $VERSION
env GOOS=windows GOARCH=amd64 go build -o ../../build/windows/unlock.exe
env GOOS=linux GOARCH=amd64 go build -o ../../build/linux/unlock
rm rsrc_windows_amd64.syso
popd

pushd ./commands/dumpsmc
echo "Building dumpsmc"
go-winres make --arch amd64 --product-version $VERSION --file-version $VERSION
env GOOS=windows GOARCH=amd64 go build -o ../../build/windows/dumpsmc.exe
env GOOS=linux GOARCH=amd64 go build -o ../../build/linux/dumpsmc
rm rsrc_windows_amd64.syso
popd

pushd ./commands/hostcaps
echo "Building hostcaps"
go-winres make --arch amd64 --product-version $VERSION --file-version $VERSION
env GOOS=windows GOARCH=amd64 go build -o ../../build/windows/hostcaps.exe
env GOOS=linux GOARCH=amd64 go build -o ../../build/linux/hostcaps
rm rsrc_windows_amd64.syso
popd

cp -v LICENSE ./build
cp -v *.md ./build
cp -v ./cpuid/linux/cpuid ./buildlinux/cpuid
cp -v ./cpuid/windows/cpuid.exe ./build/windows/cpuid.exe
cp -vr ./iso ./build
cp -vr ./templates ./build
