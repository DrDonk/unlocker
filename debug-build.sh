#!/usr/bin/env bash
echo Building debug executables

echo Building Windows executables...
env GOOS=windows GOARCH=amd64 go build -o ./windows/dumpsmc.exe ./tools/dumpsmc.go
env GOOS=windows GOARCH=amd64 go build -o ./windows/patchsmc.exe ./tools/patchsmc.go
env GOOS=windows GOARCH=amd64 go build -o ./windows/patchgos.exe ./tools/patchgos.go
env GOOS=windows GOARCH=amd64 go build -o ./windows/patchvmkctl.exe ./tools/patchvmkctl.go
env GOOS=windows GOARCH=amd64 go build -o ./windows/unlocker.exe ./command/unlocker.go

echo Building Linux executables...
env GOOS=linux GOARCH=amd64 go build -o ./linux/dumpsmc ./tools/dumpsmc.go
env GOOS=linux GOARCH=amd64 go build -o ./linux/patchsmc ./tools/patchsmc.go
env GOOS=linux GOARCH=amd64 go build -o ./linux/patchgos ./tools/patchgos.go
env GOOS=linux GOARCH=amd64 go build -o ./linux/patchvmkctl ./tools/patchvmkctl.go
env GOOS=linux GOARCH=amd64 go build -o ./linux/unlocker ./command/unlocker.go

echo Building macOS executables...
env GOOS=darwin GOARCH=amd64 go build -o ./darwin/dumpsmc ./tools/dumpsmc.go
env GOOS=darwin GOARCH=amd64 go build -o ./darwin/patchsmc ./tools/patchsmc.go
env GOOS=darwin GOARCH=amd64 go build -o ./darwin/patchgos ./tools/patchgos.go
env GOOS=darwin GOARCH=amd64 go build -o ./darwin/patchvmkctl ./tools/patchvmkctl.go
env GOOS=darwin GOARCH=amd64 go build -o ./darwin/unlocker ./command/unlocker.go

