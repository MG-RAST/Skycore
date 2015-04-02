#!/bin/bash

set -e
set -x



export ETCD_ENDPOINT=$(route|grep default|awk '{print $2}'):4001

#or
#export ETCD_ENDPOINT=172.17.42.1:4001



confd -onetime -node ${ETCD_ENDPOINT} -config-file /etc/confd/conf.d/nginx.toml 

# start nginx
service nginx start &

confd -watch -node ${ETCD_ENDPOINT} -config-file /etc/confd/conf.d/nginx.toml 

