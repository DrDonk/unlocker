#!/bin/bash
set -e

# Ensure we only use unmodified commands
export PATH=/bin:/sbin:/usr/bin:/usr/sbin

# Info
echo "Repair VMware using $1"

# Make sure only root can run our script
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

# Get a temp folder for extract
$1 -x ./repair


echo "Restoring unpatched files..."
cp -fpv ./repair/vmware-vmx/lib/bin/vmware-vmx /usr/lib/vmware/bin/vmware-vmx
cp -fpv ./repair/vmware-vmx/lib/bin/vmware-vmx-debug /usr/lib/vmware/bin/vmware-vmx-debug
cp -fpv ./repair/vmware-vmx/lib/bin/vmware-vmx-stats /usr/lib/vmware/bin/vmware-vmx-stats
cp -fpv ./repair/vmware-vmx/lib/lib/libvmwarebase.so/libvmwarebase.so /usr/lib/vmware/lib/libvmwarebase.so/libvmwarebase.so

echo Removing extracted setup...
rm -rf ./repair
echo Finished!

