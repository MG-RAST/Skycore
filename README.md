# Skycore

Tool to push Docker images into Shock and pull from Shock. Preserves some metadata and uses etcd configuration to deploy Docker images.

## Get Skycore binary
Either use the Dockerfile in this repository to statically compile skycore (The Dockerfile contains some more comments), or download pre-compiled binary (amd64):

```bash
wget https://github.com/wgerlach/Skycore/releases/download/latest/skycore
```
## Skycore from source
```bash
go get github.com/wgerlach/Skycore/skycore
```
## update skycore vendors
```bash
go get -v github.com/mjibson/party
cd $GOPATH/src/github.com/wgerlach/
git clone --recursive https://github.com/wgerlach/Skycore.git
cd Skycore/skycore
party -d vendor -c -u
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


## Deploy the skycore binary on your CoreOS machines

get IP addresses: Either a) from fleetctl (if installed)
```bash
export MACHINES=`fleetctl list-machines --full --no-legend | cut -f 2 | tr '\n' ' '` ; echo ${MACHINES}
```
or b) from nova (I admit, this is ugly.)
```bash
export MACHINES=`nova list --name <my_coreos> | grep -E -o "([0-9]{1,3}[\.]){3}[0-9]{1,3}" | tr '\n' ' '` ; echo ${MACHINES}
```
To get rid of the ssh warning "WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED", you can run:
```bash
for i in ${MACHINES} ; do ssh-keygen -f "/home/ubuntu/.ssh/known_hosts" -R ${i} ; done
```

Do some testing (read coreos version or openstack uuid):
```bash
for i in ${MACHINES} ; do echo -n "$i: " ; ssh -i ~/.ssh/coreos.pem -o StrictHostKeyChecking=no core@${i} grep PRETTY /etc/os-release  ; done
for i in ${MACHINES} ; do echo -n "$i: " ; ssh -i ~/.ssh/coreos.pem -o StrictHostKeyChecking=no core@${i} curl -s http://169.254.169.254/openstack/2013-10-17/meta_data.json | json_xs | grep uuid  ; done
```

Finally, copy the binary:
```bash
rm -f skycore ; wget https://github.com/wgerlach/Skycore/releases/download/latest/skycore
chmod +x skycore
for i in ${MACHINES} ; do scp -i ~/.ssh/coreos.pem -o StrictHostKeyChecking=no ./skycore core@${i}: ; done
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
