#!/usr/bin/env bash
#set -x
echo Clean buildribution
rm -vfr ./build
mkdir -vp ./build/iso
mkdir -vp ./build/linux
mkdir -vp ./build/windows
mkdir -vp ./build/templates

rm -vfr ./commands/rsrc_windows_amd64.syso
