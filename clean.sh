#!/usr/bin/env bash
#set -x
echo Clean distribution
rm -vfr ./dist
mkdir -vp ./dist/iso
mkdir -vp ./dist/linux
mkdir -vp ./dist/windows
mkdir -vp ./dist/templates

rm -vfr ./commands/rsrc_windows_amd64.syso
