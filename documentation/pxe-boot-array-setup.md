

This are instructions to create an mdadm RAID1 (mirror) with swap and btrfs partitions. First step is removal of any existing mdadm array or LVM ond the disks.


## 1) Clean disks

```bash
umount /media/ephemeral/
swapoff /dev/md0p1
```

remove existing RAID:
```bash
mdadm --stop /dev/md0
mdadm --remove /dev/md0

mdadm --zero-superblock /dev/sda
mdadm --zero-superblock /dev/sdb
```

also wipe MBR:
```bash
for device in /dev/sda /dev/sdb ; do 
 dd if=/dev/zero of=${device} bs=1M count=1 ;
 # wipe last megabyte to get rid of RAID
 # 2048 is 1M/512bytes (getsz returns nuber of 512blocks)
 dd if=/dev/zero of=${device} bs=512 count=2048 seek=$((`blockdev --getsz ${device}` - 2048)) ;
done
```

remove LVM:
```bash
lvm vgremove --force vg01
# lvm lvdisplay
# sudo partprobe
```


## 2) Create RAID1 (mirror) with swap and btrfs partition:

```bash
Only necessary early in boot process:
while [ ! -e /dev/sda ] ; do sleep 3; echo 'Waiting for /dev/sda'; done
while [ ! -e /dev/sdb ] ; do sleep 3; echo 'Waiting for /dev/sdb'; done
```

RAID1 with mdadm:
```bash
mdadm --create --metadata=0.90 --verbose /dev/md0 --level=mirror --raid-devices=2 /dev/sda /dev/sdb
```

On top of this RAID1 array we create two partitions. Create first partition (swap,200G):
```bash
echo -e -n "g\nn\n1\n2048\n+200G\np\nt\n14\np\nw" | fdisk /dev/md0
sleep 3 # wait before you create the next one, issue in scripts
```

Create second partition, the ephemeral disk:
```bash
echo -e -n "n\n2\n\n\np\nw" | fdisk /dev/md0
```

Create filesystems:
```bash
/usr/sbin/wipefs -f /dev/md0p1
/usr/sbin/wipefs -f /dev/md0p2

/usr/sbin/mkswap /dev/md0p1
/usr/sbin/mkfs.btrfs -f /dev/md0p2
```

Done.

Following steps will executed by cloud-config, but you can use them to test your setup.
```bash
/usr/sbin/swapon /dev/md0p1
/usr/bin/mkdir -p /media/ephemeral/
/usr/bin/mount -t btrfs /dev/md0p2 /media/ephemeral/
```

Example procedure (remove LVM, create RAID 1, create swap+data partitions):
```bash
lvm vgremove --force vg01
sleep 3

for device in /dev/sda /dev/sdb ; do 
 dd if=/dev/zero of=${device} bs=1M count=1 ;
 # wipe last megabyte to get rid of RAID
 # 2048 is 1M/512bytes (getsz returns nuber of 512blocks)
 dd if=/dev/zero of=${device} bs=512 count=2048 seek=$((`blockdev --getsz ${device}` - 2048)) ;
done
sleep 2

mdadm --create --metadata=0.90 --verbose /dev/md0 --level=mirror --raid-devices=2 /dev/sda /dev/sdb
sleep 3

echo -e -n "g\nn\n1\n2048\n+200G\np\nt\n14\np\nw" | fdisk /dev/md0
sleep 3 # wait before you create the next one, issue in scripts

echo -e -n "n\n2\n\n\np\nw" | fdisk /dev/md0
sleep 3

/usr/sbin/wipefs -f /dev/md0p1
/usr/sbin/wipefs -f /dev/md0p2
```
Filesystem will be created by cloud-config.
