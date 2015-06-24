


## Introduction to CoreOS/Fleet

http://www.slideshare.net/geekle/coreos-control-your-fleet


## From Dockerfile to Fleet service


Dockerfile:

https://github.com/MG-RAST/MG-RASTv4/blob/master/docker/Dockerfile

Building:

https://github.com/MG-RAST/MG-RASTv4

Upload image to Shock and register in etcd (this example uses another service):

https://github.com/wgerlach/Skycore#example-deployment-process-for-a-fleet-service-using-skycore

Script to register all services in MG-RAST:

https://github.com/MG-RAST/MG-RAST-infrastructure/blob/master/register_docker_image_for_service_all.sh

Fleet unit:

https://github.com/MG-RAST/MG-RAST-infrastructure/blob/master/fleet-units/mg-rast-v4-web%40.service

Fleet unit discovery:

https://github.com/MG-RAST/MG-RAST-infrastructure/blob/master/fleet-units/mg-rast-v4-web-discovery%40.service

Starting a service:

https://github.com/wgerlach/Skycore/blob/master/documentation/fleet_service_deployment.md


## Create CoreOS cluster in OpenStack

https://github.com/wgerlach/Skycore/blob/master/documentation/set_up_coreos_cluster_in_openstack.md


## Create CoreOS cluster with PXE-boot

Create RAID1 with btrfs partition:

https://github.com/wgerlach/Skycore/blob/master/documentation/pxe-boot-array-setup.md

Script to create cloud-config from template:

https://github.com/MG-RAST/MG-RAST-infrastructure/blob/master/cloud-config/update_cloud_config_pxe.sh

Cloud-config template file. Use update_cloud_config_pxe.sh to create actual cloud-config file. This is a two-step cloud-config file:

https://github.com/MG-RAST/MG-RAST-infrastructure/blob/master/cloud-config/cloud-config-pxe.yml.template




