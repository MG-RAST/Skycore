#!/bin/bash
set -e
set -x

PUBLIC_KEYS=$(cat keys.yaml| sed ':a;N;$!ba;s/\n/\\n/g')
DISCOVERY_TOKEN=$(cat discovery_token.txt)
NETWORK_INTERFACE=enp2s0f0


# use template from local git repo or download if not exist
if [ ! -e cloud-config-pxe.yml.template ] 
then
	wget --no-check-certificate https://raw.githubusercontent.com/wgerlach/Skycore/master/cloud-config-pxe.yml.template 
fi

sed -e "s;%ssh_authorized_keys%;${PUBLIC_KEYS};g" -e "s;%network_interface%;${NETWORK_INTERFACE};g" -e "s;%discovery_token%;${DISCOVERY_TOKEN};g" cloud-config-pxe.yml.template > cloud-config-pxe.yml