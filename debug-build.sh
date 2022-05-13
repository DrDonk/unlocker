#!/usr/bin/env zsh
set -x

echo Building debug executables
mkdir -p ./dist/macos

pushd ./commands/dumpsmc
go-winres make --arch amd64 --product-version 9.9.9 --file-version 9.9.9
env GOOS=windows GOARCH=amd64 go build -o ../../dist/windows
env GOOS=linux GOARCH=amd64 go build -o ../../dist/linux
env GOOS=darwin GOARCH=amd64 go build -o ../../dist/macos
rm -v rsrc_windows_amd64.syso
popd

pushd ./commands/relock
go-winres make --arch amd64 --product-version 9.9.9 --file-version 9.9.9
env GOOS=windows GOARCH=amd64 go build -o ../../dist/windows
env GOOS=linux GOARCH=amd64 go build -o ../../dist/linux
env GOOS=darwin GOARCH=amd64 go build -o ../../dist/macos
rm -v rsrc_windows_amd64.syso
popd

pushd ./commands/unlock
go-winres make --arch amd64 --product-version 9.9.9 --file-version 9.9.9
env GOOS=windows GOARCH=amd64 go build -o ../../dist/windows
env GOOS=linux GOARCH=amd64 go build -o ../../dist/linux
env GOOS=darwin GOARCH=amd64 go build -o ../../dist/macos
rm -v rsrc_windows_amd64.syso
popd

pushd ./commands/patchgos
go-winres make --arch amd64 --product-version 9.9.9 --file-version 9.9.9
env GOOS=windows GOARCH=amd64 go build -o ../../dist/windows
env GOOS=linux GOARCH=amd64 go build -o ../../dist/linux
env GOOS=darwin GOARCH=amd64 go build -o ../../dist/macos
rm -v rsrc_windows_amd64.syso
popd

pushd ./commands/patchsmc
go-winres make --arch amd64 --product-version 9.9.9 --file-version 9.9.9
env GOOS=windows GOARCH=amd64 go build -o ../../dist/windows
env GOOS=linux GOARCH=amd64 go build -o ../../dist/linux
env GOOS=darwin GOARCH=amd64 go build -o ../../dist/macos
rm -v rsrc_windows_amd64.syso
popd

pushd ./commands/patchvmkctl
go-winres make --arch amd64 --product-version 9.9.9 --file-version 9.9.9
env GOOS=windows GOARCH=amd64 go build -o ../../dist/windows
env GOOS=linux GOARCH=amd64 go build -o ../../dist/linux
env GOOS=darwin GOARCH=amd64 go build -o ../../dist/macos
rm -v rsrc_windows_amd64.syso
popd

cp -v LICENSE ./dist
cp -v *.md ./dist
cp -vr ./iso ./dist
cp -vr ./tools ./dist
