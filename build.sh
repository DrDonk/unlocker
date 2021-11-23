#!/usr/bin/env bash
env GOOS=windows GOARCH=amd64 go build -o ./build/dumpsmc.exe ./command/dumpsmc.go
env GOOS=linux GOARCH=amd64 go build -o ./build/dumpsmc.linux ./command/dumpsmc.go
env GOOS=darwin GOARCH=amd64 go build -o ./build/dumpsmc-amd64.macos ./command/dumpsmc.go
env GOOS=darwin GOARCH=arm64 go build -o ./build/dumpsmc-arm64.macos ./command/dumpsmc.go
env GOOS=windows GOARCH=amd64 go build -o ./build/patchsmc.exe ./command/patchsmc.go
env GOOS=linux GOARCH=amd64 go build -o ./build/patchsmc.linux ./command/patchsmc.go
env GOOS=darwin GOARCH=amd64 go build -o ./build/patchsmc-amd64.macos ./command/patchsmc.go
env GOOS=darwin GOARCH=arm64 go build -o ./build/patchsmc-arm64.macos ./command/patchsmc.go
env GOOS=windows GOARCH=amd64 go build -o ./build/patchgos.exe ./command/patchgos.go
env GOOS=linux GOARCH=amd64 go build -o ./build/patchgos.linux ./command/patchgos.go
env GOOS=darwin GOARCH=amd64 go build -o ./build/patchgos-amd64.macos ./command/patchgos.go
env GOOS=darwin GOARCH=arm64 go build -o ./build/patchgos-arm64.macos ./command/patchgos.go
env GOOS=windows GOARCH=amd64 go build -o ./build/patchvmkctl.exe ./command/patchvmkctl.go
env GOOS=linux GOARCH=amd64 go build -o ./build/patchvmkctl.linux ./command/patchvmkctl.go
env GOOS=darwin GOARCH=amd64 go build -o ./build/patchvmkctl-amd64.macos ./command/patchvmkctl.go
env GOOS=darwin GOARCH=arm64 go build -o ./build/patchvmkctl-arm64.macos ./command/patchvmkctl.go
