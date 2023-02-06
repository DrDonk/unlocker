#!/bin/sh
#set -x

# Read current version
VERSION=$(<VERSION)
VERSION=${VERSION//.}

echo Create distribution files - $VERSION

rm -vf ./dist/unlocker$VERSION.zip
rm -vrf ./dist/unlocker$VERSION
7z a ./dist/unlocker$VERSION.zip ./build/*
tar czvf ./dist/unlocker$VERSION.tgz ./build/*
7z x -o./dist/unlocker$VERSION ./dist/unlocker$VERSION.zip

shasum -a 256  ./dist/unlocker$VERSION.tgz
shasum -a 256  ./dist/unlocker$VERSION.zip
shasum -a 512  ./dist/unlocker$VERSION.tgz
shasum -a 512  ./dist/unlocker$VERSION.zip
