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
for i in ${MACHINES} ; do scp -i <your_private_ssh_key> -o StrictHostKeyChecking=no ./skycore core@${i}: ; done
```

## Docker image registration for services with etcd
Once you have built and uploaded a new Docker image for a particular service to Shock, you need to update the etcd configuration to point to the new Shock node. To get access to etcd you probably have to log into one of the machines. The service name has to match the unit name, for example "mg-rast-v4-web":
```bash
curl -L http://127.0.0.1:4001/v2/keys/service_images/<servicename>/shock -XPUT -d value="shock.metagenomics.anl.gov/node/<node>"
```

You can read the current configuration with the same url:
```bash
curl -L http://127.0.0.1:4001/v2/keys/service_images/<servicename>/shock
```

You can also use the etcdctl command to modify values and to browse the etcd tree. For example "etcdctl ls /service_images/" will show for which services Docker images are registered.


## Fleet service deployment
The unit files in this example are using skycore, which needs to be installed on all machines. This also means that docker images have to be registered with etcd.

Log into a machine and confirm:
```bash
fleetctl list-machines
```

Download unit files from git repo. Then deploy unit files for a service and its sidekick (called discovery): 
```bash
fleetctl submit mg-rast-v4-web\@.service mg-rast-v4-web-discovery\@.service
fleetctl list-unit-files
```

Start 2 instances:
```bash
fleetctl start fleetctl start mg-rast-v4-web{,-discovery}@{1..2}.service
fleetctl list-units
```

Destroy units and delete unit files. Delete unit files only when you need to make changes to them:
```bash
fleetctl destroy fleetctl start mg-rast-v4-web{,-discovery}@{1..2}.service
fleetctl destroy mg-rast-v4-web\@.service mg-rast-v4-web-discovery\@.service
```

Monitoring:
```bash
fleetctl list-machines
fleetctl list-unit-files 
fleetctl list-units
```
