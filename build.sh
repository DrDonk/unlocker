#!/usr/bin/env zsh
set -x

echo Building release executables

pushd ./commands/dumpsmc
go-winres make --arch amd64 --product-version 9.9.9 --file-version 9.9.9
env GOOS=windows GOARCH=amd64 go build -o ../../dist/windows
env GOOS=linux GOARCH=amd64 go build -o ../../dist/linux
rm -v rsrc_windows_amd64.syso
popd

pushd ./commands/relock
go-winres make --arch amd64 --product-version 9.9.9 --file-version 9.9.9
env GOOS=windows GOARCH=amd64 go build -o ../../dist/windows
env GOOS=linux GOARCH=amd64 go build -o ../../dist/linux
rm -v rsrc_windows_amd64.syso
popd

pushd ./commands/unlock
go-winres make --arch amd64 --product-version 9.9.9 --file-version 9.9.9
env GOOS=windows GOARCH=amd64 go build -o ../../dist/windows
env GOOS=linux GOARCH=amd64 go build -o ../../dist/linux
rm -v rsrc_windows_amd64.syso
popd


cp -v LICENSE ./dist
cp -v *.md ./dist
cp -vr ./iso ./dist
cp -vr ./tools ./dist
