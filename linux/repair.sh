#!/bin/bash
set -e

echo "Repair Vmware using"

# Ensure we only use unmodified commands
export PATH=/bin:/sbin:/usr/bin:/usr/sbin

# Make sure only root can run our script
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

version=$(grep -i player.product.version /etc/vmware/config | cut -d'"' -f2- | rev | cut -c 2- | rev)
build=$(grep -i product.buildnumber /etc/vmware/config | cut -d'"' -f2- | rev | cut -c 2- | rev)
IFS='.' read -r -a product <<< "$version"

printf "VMware product version: %s.%s\n\n" "$version" "$build"
#printf "Major:    ${product[0]}\n"
#printf "Minor:    ${product[1]}\n"
#printf "Revision: ${product[2]}\n"
#printf "Build:    ${build}\n"


echo Finished!

