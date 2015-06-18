#!/bin/bash
set -x

export DEVICES=`echo /dev/sd{a,b}`

umount /media/ephemeral/
swapoff /dev/md0p1

mdadm --stop /dev/md0
mdadm --remove /dev/md0

mdadm --zero-superblock /dev/sd{a,b}1
mdadm --zero-superblock ${DEVICES}

for device in ${DEVICES} ; do 
  echo -e -n "o\\ny\\nw\ny\\n" | gdisk ${device}
  sleep 2
done
