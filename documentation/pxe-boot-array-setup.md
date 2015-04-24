

This are instructions to create an mdadm RAID1 (mirror) with swap and btrfs partitions. First step is removal of any existing mdadm array or LVM ond the disks.


## Remove existing RAID

raid_wipe.sh
```bash
#!/bin/bash
set -x

umount /media/ephemeral/
swapoff /dev/md0p1
set -e

mdadm --stop /dev/md0
mdadm --remove /dev/md0

mdadm --zero-superblock /dev/sda
mdadm --zero-superblock /dev/sdb
```


## Remove LVM
Example procedure (remove LVM, create RAID 1, create swap+data partitions):

lvm_wipe.sh
```bash
#!/bin/bash
set -x
set -e

lvm vgremove --force vg01
# lvm lvdisplay
# sudo partprobe
sleep 3

for device in /dev/sda /dev/sdb ; do 
 dd if=/dev/zero of=${device} bs=1M count=1 ;
 # wipe last megabyte to get rid of RAID
 # 2048 is 1M/512bytes (getsz returns nuber of 512blocks)
 dd if=/dev/zero of=${device} bs=512 count=2048 seek=$((`blockdev --getsz ${device}` - 2048)) ;
done
sleep 2
# seems to require reboot here, because resource is busy. not sure where that comes from
```

## Create RAID1 with swap+data partitions
raid1.sh
```bash
#!/bin/bash
set -x
set -e

#create RAID1
echo y | mdadm --create --metadata=0.90 --verbose /dev/md0 --level=mirror --raid-devices=2 /dev/sda /dev/sdb
sleep 3

#remove secondary GPT header
echo -e -n "o\\ny\\nw\\ny\\n" | gdisk /dev/md0
sleep 3

# create swap partition
echo -e -n "g\nn\n1\n2048\n+200G\np\nt\n14\np\nw" | fdisk /dev/md0
sleep 3 # wait before you create the next one, issue in scripts

#create data partition
echo -e -n "n\n2\n\n\np\nw" | fdisk /dev/md0
sleep 3

#remove secondary GPT header
echo -e -n "o\\ny\\nw\\ny\\n" | gdisk /dev/sda
sleep 1
echo -e -n "o\\ny\\nw\\ny\\n" | gdisk /dev/sdb
sleep 3

/usr/sbin/wipefs -f /dev/md0p1
/usr/sbin/wipefs -f /dev/md0p2
```

Filesystem will be created by cloud-config. but you can manually test you setup:
```bash
#/usr/sbin/mkswap /dev/md0p1
#/usr/sbin/mkfs.btrfs -f /dev/md0p2

/usr/sbin/swapon /dev/md0p1
/usr/bin/mkdir -p /media/ephemeral/
/usr/bin/mount -t btrfs /dev/md0p2 /media/ephemeral/
```


## Multiple machines example:
```bash
export MACHINES=`eval echo "{1..8} {10..11}"` ; echo ${MACHINES}
# test ssh
for i in ${MACHINES} ; do echo "$i: " ; ssh -o ConnectTimeout=1 -i ~/.ssh/wo_magellan_private_key.pem core@bio-worker${i} grep PRETTY /etc/os-release ; done
# copy lvm_wipe.sh
for i in ${MACHINES} ; do echo "$i: " ; scp -i ~/.ssh/wo_magellan_private_key.pem lvm_wipe.sh core@bio-worker${i}: ; done
# execute lvm_wipe.sh
for i in ${MACHINES} ; do echo "$i: " ; ssh -i ~/.ssh/wo_magellan_private_key.pem core@bio-worker${i} sudo ./lvm_wipe.sh ; done
#reboot
for i in ${MACHINES} ; do echo "$i: " ; ssh -i ~/.ssh/wo_magellan_private_key.pem core@bio-worker${i} sudo reboot ; done
#remove keys:
for i in ${MACHINES} ; do echo "$i: " ; ssh-keygen -f "/homes/wgerlach/.ssh/known_hosts" -R bio-worker${i} ; done
#test again ssh
for i in ${MACHINES} ; do echo "$i: " ; ssh -o ConnectTimeout=1 -i ~/.ssh/wo_magellan_private_key.pem core@bio-worker${i} grep PRETTY /etc/os-release ; done
#copy raid1.sh
for i in ${MACHINES} ; do echo "$i: " ; scp -i ~/.ssh/wo_magellan_private_key.pem raid1.sh core@bio-worker${i}: ; done
# execute raid1.sh
for i in ${MACHINES} ; do echo "$i: " ; ssh -i ~/.ssh/wo_magellan_private_key.pem core@bio-worker${i} sudo ./raid1.sh ; done
# reboot last time
for i in ${MACHINES} ; do echo "$i: " ; ssh -i ~/.ssh/wo_magellan_private_key.pem core@bio-worker${i} sudo reboot ; done
```
