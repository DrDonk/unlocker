#!/usr/bin/env bash
env GOOS=windows GOARCH=amd64 go build -o patchsmc.exe patchsmc.go
env GOOS=linux GOARCH=amd64 go build -o patchsmc.linux patchsmc.go
env GOOS=darwin GOARCH=amd64 go build -o patchsmc-amd64.macos patchsmc.go
env GOOS=darwin GOARCH=arm64 go build -o patchsmc-arm64.macos patchsmc.go
