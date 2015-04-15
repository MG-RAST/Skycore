#!/bin/bash

# template requires discovery token and public ssh keys
# this hack below is needed to get IP address
# tested on magellan
 
until ! [[ -z $COREOS_PRIVATE_IPV4 ]]; do
 
 echo "COREOS_PUBLIC_IPV4=$(curl http://169.254.169.254/latest/meta-data/public-ipv4)" > /etc/environment
 echo "COREOS_PRIVATE_IPV4=$(curl http://169.254.169.254/latest/meta-data/local-ipv4)" >> /etc/environment
 source /etc/environment
 
done


until ! [[ -z $INSTANCE_TYPE ]]; do
 INSTANCE_TYPE=$(curl http://169.254.169.254/latest/meta-data/instance-type)
done

MOUNT_POINT="/mnt"
mkdir -p ${MOUNT_POINT}


EPHEMERAL="`curl -s -f -m 20 http://169.254.169.254/latest/meta-data/block-device-mapping/ephemeral0`"
if [ -z "${EPHEMERAL}" ]; then
	# workaround for a bug in EEE 2
	EPHEMERAL="`curl -s -f -m 20 http://169.254.169.254/latest/meta-data/block-device-mapping/ephemeral`"
fi
if [ ! -z "${EPHEMERAL}" ]; then
MOUNTUNIT=$(cat <<EOF
   - name: media-ephemeral.mount
     command: start
     content: |
       [Mount]
       What=${EPHEMERAL}
       Where=/media/ephemeral
       Type=ext3
EOF
)
fi


 
cat << 'EOF' > /tmp/user_data.yml 
#cloud-config


coreos: 
 etcd:
   discovery: https://discovery.etcd.io/<token>
   addr: \$private_ipv4:4001
   peer-addr: \$private_ipv4:7001
 fleet:
   public-ip: \$private_ipv4
EOF
cat << EOF >> /tmp/user_data.yml
   metadata: "instance_type=${INSTANCE_TYPE}"
 units:
   - name: etcd.service
     command: start
   - name: fleet.service
     command: start
   - name: settimezone.service
     command: start
     content: |
       [Unit]
       Description=Set the timezone

       [Service]
       ExecStart=/usr/bin/timedatectl set-timezone America/Chicago
       RemainAfterExit=yes
       Type=oneshot
${MOUNTUNIT}
EOF
cat << 'EOF' >> /tmp/user_data.yml
ssh_authorized_keys:
  # include one or more SSH public keys
  - <public ssh key>
EOF
 
sudo coreos-cloudinit --from-file=/tmp/user_data.yml