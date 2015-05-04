#!/bin/bash


set -x

export ETCD_ENDPOINT=$(route|grep default|awk '{print $2}'):4001
# usually => export ETCD_ENDPOINT=172.17.42.1:4001
export CONF_DIR="/Skycore/services/nginx/confd"
# does not work: export TOML_FILE="/Skycore/services/nginx/confd/nginx.toml"
export CONFD_ARGS="-node ${ETCD_ENDPOINT} -confdir=${CONF_DIR}"

mkdir -p /etc/nginx/sites-enabled/

# the first call often fails
confd -onetime ${CONFD_ARGS}
sleep 1
# start nginx

set -e

confd -watch=false ${CONFD_ARGS} &

sleep 1

nginx -c /Skycore/services/nginx/nginx.conf
