#!/usr/bin/env python3
# coding=utf-8

# SPDX-FileCopyrightText: © 2022-23 David Parsons
# SPDX-License-Identifier: MIT

import subprocess
import sys

import macrecovery


def main():

    print('\nOC4VM macOS Recovery VMDK Maker')
    print('=================================')
    print('(c) David Parsons 2022-23\n')

    # Print the menu
    print('Create a VMware VMDK Recovery Image')
    print('1. Catalina')
    print('2. Big Sur')
    print('3. Monterey')
    print('4. Ventura')

    # And get the input
    while True:
        selection = input('Input menu number: ')

        if selection == '1':
            basename = 'catalina'
            boardid = 'Mac-6F01561E16C75D06'
            break
        if selection == '2':
            basename = 'bigsur'
            boardid = 'Mac-2BD1B31983FE1663'
            break

        if selection == '3':
            basename = 'monterey'
            boardid = 'Mac-A5C67F76ED83108C'
            break
        if selection == '4':
            basename = 'ventura'
            boardid = 'Mac-B4831CEBD52A0C4C'
            break

    print('Downloading DMG... \n')

    # Setup args for macrecovery and get the download
    sys.argv = ['macrecovery.py',
                'download',
                '-b', boardid,
                '-m', '00000000000000000',
                '--basename', basename,
                '-os', 'latest']

    macrecovery.main()

    # Convert DMG to IMG using dmg2img
    dmg = f'{basename}.dmg'
    vmdk = f'{basename}.vmdk'

    print('Convertng to VMDK: ')
    qemu_img = f'qemu_img convert -O vmdk {dmg} {vmdk} -p'
    subprocess.call(qemu_img.split())
    print(f'Created VMDK disk: {vmdk}')
    return


if __name__ == '__main__':
    main()
