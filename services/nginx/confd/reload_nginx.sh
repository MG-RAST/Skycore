#!/bin/bash
set -e
set -x


if [ ! -e /usr/bin/docker ] ; then
  curl -O https://get.docker.com/builds/Linux/x86_64/docker-1.6.0 && mv docker-1.6.0 /usr/bin/docker && chmod +x /usr/bin/docker
fi

/usr/bin/docker exec mgrast_nginx /usr/sbin/nginx -s reload -c /Skycore/services/nginx/nginx.conf
