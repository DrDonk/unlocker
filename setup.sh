#!/bin/bash
set -e
rm -rfv ./tests/*
cp -prv ./samples/* ./tests/
