#!/bin/bash
set -x

umount /media/ephemeral/
swapoff /dev/md0p1

mdadm --stop /dev/md0
mdadm --remove /dev/md0

mdadm --zero-superblock /dev/sda1 /dev/sdb1
mdadm --zero-superblock /dev/sda /dev/sdb

echo -e -n "o\\ny\\nw\ny\\n" | gdisk /dev/sda
sleep 2
echo -e -n "o\\ny\\nw\ny\\n" | gdisk /dev/sdb
sleep 2
