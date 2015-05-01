#!/bin/bash


set -x



export ETCD_ENDPOINT=$(route|grep default|awk '{print $2}'):4001

export TOML_FILE="/Skycore/services/nginx/confd/nginx.toml"

#or
#export ETCD_ENDPOINT=172.17.42.1:4001



confd -onetime -node ${ETCD_ENDPOINT} -config-file ${TOML_FILE} 

# start nginx
service nginx start &
set -e

confd -watch -node ${ETCD_ENDPOINT} -config-file ${TOML_FILE} 

