# Fleet service deployment

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
systemctl status -l <service>
journalctl -b -u <service>  # only locally
fleetctl journal <service>  # creates a ssh tunnel to remote machine
```

## Starting services
AWE server, uses MachineID as argument
```bash
fleetctl start awe-server{,-mongodb,-discovery}@1dc3558aa345483292f2f858de0e23e1.service
```
AWE clients (use ranges)
```bash
for i in {2..4} ; do fleetctl start awe-client\@$i.service ; done
```
