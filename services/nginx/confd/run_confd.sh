#!/bin/bash


set -x

if [ ! -e /usr/bin/docker ] ; then
  curl -O https://get.docker.com/builds/Linux/x86_64/docker-1.6.0 && mv docker-1.6.0 /usr/bin/docker && chmod +x /usr/bin/docker
fi

export ETCD_ENDPOINT=$(route|grep default|awk '{print $2}'):4001
# usually => export ETCD_ENDPOINT=172.17.42.1:4001
export CONF_DIR="/Skycore/services/nginx/confd"
export TOML_FILE="/Skycore/services/nginx/confd/conf.d/nginx.toml"
export CONFD_ARGS="-node ${ETCD_ENDPOINT} -confdir=${CONF_DIR} -config-file=${TOML_FILE}"

mkdir -p /etc/nginx/sites-enabled/

# the first call often fails
confd -onetime ${CONFD_ARGS}
sleep 1
# start nginx

set -e

confd -watch=false -interval=5 ${CONFD_ARGS}


# rund in another container now
#nginx -c /Skycore/services/nginx/nginx.conf
