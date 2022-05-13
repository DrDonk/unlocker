#!/usr/bin/env zsh
#set -x
echo Zip distribution
if ! [ $# -eq 1 ] ; then
  echo "Product version not found" >&2
  exit 1
fi

rm -vf ./build/unlocker$1.zip
rm -vrf ./build/unlocker$1
7z a ./build/unlocker$1.zip ./dist/*
7z x -o./build/unlocker$1 ./build/unlocker$1.zip
