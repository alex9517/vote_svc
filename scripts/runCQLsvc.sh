#!/bin/bash

#  Created : 2023-Nov-14
# Modified : 2024-Apr-16

# This is a shell script to run the Apache Cassandra service
# as a Docker image. Do not forget to add that server name to your
# computer's (server's,..) /etc/hosts (or whatever DNS you prefer):
#   172.25.70.27   s7027

# This is just a memo:
# Class A: 10.0.0.0 - 10.255.255.255
# Class B: 172.16.0.0 - 172.31.255.255
# Class C: 192.168.0.0 - 192.168.255.255

IMAGE=cassandra
VERSION=latest
NETWORK=net_17216
IP_ADDR=172.16.70.31
HOSTNAME=s7031
NAME=cassandra01
DATA_DIR=/u01/data/cassandra
# echo $NETWORK
# echo $IP_ADDR
# echo $HOSTNAME

# -d, --detach : Run container in background and print container ID;
# --name : Assign a name to the container;
#  -e ALLOW_ORIGINS=https://ws4,http://localhost,https://localhost,http://localhost:3000 \

docker run --rm -it --name $NAME  --net $NETWORK --ip $IP_ADDR --hostname $HOSTNAME \
  --add-host=ws4:192.168.10.94 \
  -v $DATA_DIR:/var/lib/cassandra \
  -p 49181:9042 \
  -d $IMAGE:$VERSION

#  -e ALLOW_ORIGINS=* \
#  -e DATABASE_URL=postgres://appuser:pass@ws4:5432/db2 \
#  -v /home/alex/tmp/proj2/go/showbiz/scripts/service.pem:/etc/x509/https/service.pem \
#  -v /home/alex/tmp/proj2/go/showbiz/scripts/service.key:/etc/x509/https/service.key \

# --- END ---
