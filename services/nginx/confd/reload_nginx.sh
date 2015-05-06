#!/bin/bash
set -e
set -x


export DOCKER_VERSION="1.5.0"
if [ ! -e /usr/bin/docker ] ; then
  curl -O https://get.docker.com/builds/Linux/x86_64/docker-${DOCKER_VERSION}
  mv docker-${DOCKER_VERSION} /usr/bin/docker 
  chmod +x /usr/bin/docker
fi

/usr/bin/docker exec mgrast_nginx /usr/sbin/nginx -s reload -c /Skycore/services/nginx/nginx.conf
