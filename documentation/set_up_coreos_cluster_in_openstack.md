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

Other service ports:
```bash
nova secgroup-create solr-m5nr "Solr m5nr 8983"
nova secgroup-add-rule solr-m5nr tcp 8983 8983 0.0.0.0/0

nova secgroup-create mgrast-v4-web "MGRAST v4 web 80"
nova secgroup-add-rule mgrast-v4-web tcp 80 80 0.0.0.0/0
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
## Log in to your CoreOS cluster

Login with forwarding your ssh user agent. Run these commands on your client outside of the CoreOS cluster:
```bash
cd ~/.ssh 
ln -s <your private key> coreos.pem
eval $(ssh-agent)
ssh-add ~/.ssh/coreos.pem
ssh -A core@<instance>
```
You may want to assign a public IP address to one of you CoreOS instances.

## Optional: Set up fleetctl locally to talk to cluster
```bash
wget https://github.com/coreos/fleet/releases/download/v0.9.1/fleet-v0.9.1-linux-amd64.tar.gz
tar xvzf fleet-v0.9.1-linux-amd64.tar.gz
cp fleet-v0.9.1-linux-amd64/fleetctl /usr/local/bin/

#in your .bashrc
export FLEETCTL_TUNNEL=<ip address of one coreos instance>
```
