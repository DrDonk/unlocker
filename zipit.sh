#!/bin/sh
#set -x
echo Zip distribution
if ! [ $# -eq 1 ] ; then
  echo "Product version not found: xyz (e.g. 123)" >&2
  exit 1
fi

rm -vf ./dist/unlocker$1.zip
rm -vrf ./dist/unlocker$1
7z a ./dist/unlocker$1.zip ./build/*
7z x -o./dist/unlocker$1 ./dist/unlocker$1.zip
