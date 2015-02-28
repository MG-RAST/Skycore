# Skycore

Tool to push Docker images into Shock and pull from Shock. Preserves some metadata and uses etcd configuration to deploy Docker images.

## Get Skycore binary
Either use the Dockerfile in this repository to statically compile skycore (The Dockerfile contains some more comments), or download pre-compiled binary:

```bash
wget http://dunkirk.mcs.anl.gov/~wgerlach/skycore
```

## Example deployment process for a fleet service using skycore
Build image (requires docker):
```bash
docker build --tag=mgrast/m5nr-solr:20150223_1700 --force-rm --no-cache https://raw.githubusercontent.com/MG-RAST/myM5NR/master/solr/docker/Dockerfile
```
Upload image to Shock:
```bash
skycore push --shock=http://shock.metagenomics.anl.gov --token=$TOKEN mgrast/m5nr-solr:20150223_1700
```
Register shock node (of the new image) with etcd (requires etcd access):
```bash
curl -L http://127.0.0.1:4001/v2/keys/service_images/m5nr-solr/shock -XPUT -d value="shock.metagenomics.anl.gov/node/<node_id>"
```
Please update/add the corresponding line register_docker_image_for_service_all.sh .

And restart fleet service... either with fleetctl or fleet api..

# Setting up CoreOS cluster
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
chmod +x skycore
for i in ${MACHINES} ; do scp -i <your_private_ssh_key> -o StrictHostKeyChecking=no ./skycore core@${i}: ; done
```

## Log in to your CoreOS cluster

Login with forwarding your ssh user agent. Run these commands on your client outside of the CoreOS cluster:
```bash
eval $(ssh-agent)
ssh-add <path to your private key>
ssh -A core@<instance>
```
You may want to assign a public IP address to one of you CoreOS instances.

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
fleetctl submit mg-rast-v4-web{,-discovery}\@.service
fleetctl list-unit-files
```

Start 2 instances of each of mg-rast-v4-web and mg-rast-v4-web-discovery:
```bash
fleetctl start mg-rast-v4-web{,-discovery}\@{1..2}.service
fleetctl list-units
```
The mg-rast-v4-web-discovery sidekicks provide service discovery via the etcd keys /services/mg-rast-v4-web/mg-rast-v4-web@1 and /services/mg-rast-v4-web/mg-rast-v4-web@2 . The example below shows the service information stored by a sidekick:

```bash
etcdctl get /services/mg-rast-v4-web/mg-rast-v4-web@1
{ "host":"coreos-wolfgang-139c22c0-4fbc-4de1-9d94-81507ccf323f.novalocal","port": 80,"COREOS_PRIVATE_IPV4":"10.1.12.67","COREOS_PUBLIC_IPV4":""}
```

Destroy units and delete unit files. Delete unit files only when you need to make changes to them:
```bash
fleetctl destroy mg-rast-v4-web{,-discovery}\@{1..2}.service
fleetctl destroy mg-rast-v4-web{,-discovery}\@.service
```

Monitoring:
```bash
fleetctl list-machines
fleetctl list-unit-files 
fleetctl list-units
```

Debugging:
```bash
systemctl status -l service
journalctl -b -u service
```
