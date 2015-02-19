# Skycore

Tool to push Docker images into Shock and pull from Shock. Preserves some metadata and uses etcd configuration to deploy Docker images.

Use the Dockerfile in this repository to statically compile the skycore. The Dockerfile contains some more comments.

# CoreOS stuff
Instructions for OpenStack

## Create CoreOS image

Search for latest image in http://stable.release.core-os.net/amd64-usr/ and download:
```bash
wget http://stable.release.core-os.net/amd64-usr/557.2.0/coreos_production_openstack_image.img.bz2
bunzip2 coreos_production_openstack_image.img.bz2
```

import image in OpenStack:
```bash
export COREOS="CoreOS_VERSION"
glance image-create --name ${COREOS}\
  --container-format bare \
  --disk-format qcow2 \
  --file coreos_production_openstack_image.img \
  --is-public False
```

## Prepare cloud-config.sh
Get etcd discovery token:
```bash
curl -w "\n" https://discovery.etcd.io/new
```

Download cloud-config-openstack.sh from this repository and add your new etcd token and add public ssh keys.

## Start CoreOS VM

If you do not already have an security group for CoreOS, please create one to open the required ports:
```bash
nova secgroup-create coreos "CoreOS ports 4001 and 7001"
nova secgroup-add-rule coreos tcp 4001 4001 10.1.0.0/16
nova secgroup-add-rule coreos tcp 7001 7001 10.1.0.0/16
```

Use nova boot to start your instances:
```bash
nova boot \
  --user-data ./cloud-config.sh \
  --image CoreOS_557.2.0 \
  --key-name <your_openstack_public_ssh_key_name> \
  --flavor i2.medium.sd \
  --num-instances 3 \
  --security-groups default,coreos \
  my_coreos
```

## Deploy the skycore binary on your CoreOS machines

get IP addresses: (I admit, this is ugly.)
```bash
export MACHINES=`nova list --name my_coreos | grep -E -o "([0-9]{1,3}[\.]){3}[0-9]{1,3}" | tr '\n' ' '` ; echo ${MACHINES}
```
and copy the binary
```bash
for i in ${MACHINES} ; do scp -i ~/.ssh/wo_magellan_pubkey.pem -o StrictHostKeyChecking=no ./skycore core@${i}: ; done
```

## Fleet
Deploy and start service using fleet. 
```bash
fleetctl submit mg-rast-v4-web\@.service mg-rast-v4-web-discovery\@.service
fleetctl start mg-rast-v4-web\@1
(running)
fleetctl stop mg-rast-v4-web\@1
fleetctl unload mg-rast-v4-web\@1
fleetctl destroy mg-rast-v4-web\@.service
```

Monitoring:
```bash
fleetctl list-machines
fleetctl list-unit-files 
fleetctl list-units
```
