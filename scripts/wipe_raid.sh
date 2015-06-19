#!/bin/bash
set -x

export DEVICES_STR="/dev/sd{a,b}"
export DEVICES=$(eval echo ${DEVICES_STR})
export PARTITIONS=$(eval echo ${DEVICES_STR}1)

umount /media/ephemeral/
swapoff /dev/md0p1

mdadm --stop /dev/md0
mdadm --remove /dev/md0

mdadm --zero-superblock ${PARTITIONS}
mdadm --zero-superblock ${DEVICES}

for device in ${DEVICES} ; do 
  echo -e -n "o\\ny\\nw\ny\\n" | gdisk ${device}
  sleep 2
done
