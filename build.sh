#!/usr/bin/env bash
env GOOS=windows GOARCH=amd64 go build -o ./windows/dumpsmc.exe ./tools/dumpsmc.go
env GOOS=windows GOARCH=amd64 go build -o ./windows/unlocker.exe ./command/unlocker.go
env GOOS=linux GOARCH=amd64 go build -o ./linux/dumpsmc ./tools/dumpsmc.go
env GOOS=linux GOARCH=amd64 go build -o ./linux/unlocker ./command/unlocker.go
env GOOS=darwin GOARCH=amd64 go build -o ./macos/dumpsmc ./tools/dumpsmc.go
