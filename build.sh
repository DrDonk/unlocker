#!/usr/bin/env bash
echo Building release executables

echo Building Windows executables...
env GOOS=windows GOARCH=amd64 go build -o ./windows/dumpsmc.exe ./tools/dumpsmc.go
env GOOS=windows GOARCH=amd64 go build -o ./windows/unlocker.exe ./command/unlocker.go

echo Building Linux executables...
env GOOS=linux GOARCH=amd64 go build -o ./linux/dumpsmc ./tools/dumpsmc.go
env GOOS=linux GOARCH=amd64 go build -o ./linux/unlocker ./command/unlocker.go
