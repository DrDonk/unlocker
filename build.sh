#!/usr/bin/env bash
env GOOS=windows GOARCH=amd64 go build ./tools/dumpsmc.go
env GOOS=linux GOARCH=amd64 go build ./tools/dumpsmc.go
env GOOS=windows GOARCH=amd64 go build ./command/unlocker.go
env GOOS=linux GOARCH=amd64 go build ./command/unlocker.go
zip ./build/release.zip dumpsmc dumpsmc.exe unlocker unlocker.exe unlocker.exe.manifest iso/*
