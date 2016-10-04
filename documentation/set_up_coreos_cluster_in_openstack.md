# Setting up CoreOS cluster
Instructions for OpenStack

## Create CoreOS image

Search for latest image in http://stable.release.core-os.net/amd64-usr/ and download:
```bash
export COREOS_VERSION=$(curl http://stable.release.core-os.net/amd64-usr/current/version.txt | grep "COREOS_VERSION_ID" | cut -d '=' -f 2)
echo ${COREOS_VERSION}

wget http://stable.release.core-os.net/amd64-usr/${COREOS_VERSION}/coreos_production_openstack_image.img.bz2
bunzip2 coreos_production_openstack_image.img.bz2
```

install glance
```bash
pip install python-glanceclient
```

import image in OpenStack:
```bash
glance image-create --name CoreOS_${COREOS_VERSION}\
  --container-format bare \
  --disk-format qcow2 \
  --file coreos_production_openstack_image.img
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

Other docs moved to:
https://github.com/MG-RAST/MG-RAST-infrastructure/tree/master/docs
