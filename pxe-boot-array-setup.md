## Clean disks

```bash
umount /media/ephemeral/
swapoff /dev/md0p1

# remove RAID:
mdadm --stop /dev/md0
mdadm --remove /dev/md0

mdadm --zero-superblock /dev/sda
mdadm --zero-superblock /dev/sdb


# wipe MBR:
for device in /dev/sda /dev/sdb ; do 
 dd if=/dev/zero of=${device} bs=1M count=1 ;
 # wipe last megabyte to get rid of RAID
 # 2048 is 1M/512bytes (getsz returns nuber of 512blocks)
 dd if=/dev/zero of=${device} bs=512 count=2048 seek=$((`blockdev --getsz ${device}` - 2048)) ;
done

#remove LVM ?
lvm vgremove --force vg01
# lvm lvdisplay
# sudo partprobe
```


## Mirror disks with swap:

```bash
Prepare disks for boot:
while [ ! -e /dev/sda ] ; do sleep 3; echo 'Waiting for /dev/sda'; done
while [ ! -e /dev/sdb ] ; do sleep 3; echo 'Waiting for /dev/sdb'; done

 #mdadm with RAID1:
mdadm --create --metadata=0.90 --verbose /dev/md0 --level=mirror --raid-devices=2 /dev/sda /dev/sdb

#create first partition (swap,200G)
echo -e -n "g\nn\n1\n2048\n+200G\np\nt\n14\np\nw" | fdisk /dev/md0
sleep 3
#create second partition
echo -e -n "n\n2\n\n\np\nw" | fdisk /dev/md0


/usr/sbin/wipefs -f /dev/md0p1
/usr/sbin/wipefs -f /dev/md0p2
/usr/sbin/mkfs.btrfs -f /dev/md0p2



#this is done by cloud-config:
/usr/sbin/mkswap /dev/md0p1
/usr/sbin/swapon /dev/md0p1
/usr/bin/mkdir -p /media/ephemeral/
```