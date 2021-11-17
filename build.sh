#!/usr/bin/env bash
env GOOS=windows GOARCH=amd64 go build -o dumpsmc.exe dumpsmc.go
env GOOS=linux GOARCH=amd64 go build -o dumpsmc.linux dumpsmc.go
env GOOS=darwin GOARCH=amd64 go build -o dumpsmc-amd64.macos dumpsmc.go
env GOOS=darwin GOARCH=arm64 go build -o dumpsmc-arm64.macos dumpsmc.go
