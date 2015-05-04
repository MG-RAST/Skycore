#!/bin/bash


set -x

export ETCD_ENDPOINT=$(route|grep default|awk '{print $2}'):4001
# usually => export ETCD_ENDPOINT=172.17.42.1:4001
export CONF_DIR="/Skycore/services/nginx/confd"
export TOML_FILE="/Skycore/services/nginx/confd/nginx.toml"
export CONFD_ARGS="-node ${ETCD_ENDPOINT} -config-file ${TOML_FILE} -confdir=${CONF_DIR}"


# the first call often fails
confd -onetime ${CONFD_ARGS}

# start nginx

set -e

confd -watch=false ${CONFD_ARGS} &

nginx -c /Skycore/services/nginx/nginx.conf
