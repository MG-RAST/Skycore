# Skycore

Tool to push Docker images into Shock and pull from Shock. Preserves some metadata and uses etcd configuration to deploy Docker images.

## Get Skycore binary
Either use the Dockerfile in this repository to statically compile skycore (The Dockerfile contains some more comments), or download pre-compiled binary (amd64):

```bash
wget https://github.com/wgerlach/Skycore/releases/download/latest/skycore
```
## Skycore from source (requires golang)
```bash
go get github.com/wgerlach/Skycore/skycore
```

### Build container
```bash
docker build --tag skycore/factory:latest .
```

### Use container to compile and get binary:
```bash
mkdir -p ~/skycore_bin
docker run -t -i --name sky_fac -v ~/skycore_bin:/gopath/bin skycore/factory:latest bash
cd /gopath/src/github.com/wgerlach/Skycore && git pull && /compile.sh
```

### skycore execution within container
```bash
# container will need access to docker socket
mkdir -p ~/skycore_bin
docker run -t -i -v /var/run/docker.sock:/var/run/docker.sock --name sky_fac -v ~/skycore_bin:/gopath/bin skycore/factory:latest bash
```

## update skycore vendors
```bash
go get -v github.com/mjibson/party
cd $GOPATH/src/github.com/wgerlach/
git clone --recursive git@github.com:wgerlach/Skycore.git
cd Skycore/skycore
party -d vendor -c -u
```



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

