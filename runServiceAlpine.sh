#!/bin/bash

#  Created : 2024-Apr-05
# Modified : 2024-Apr-05

# This is a shell script to run the 'vote-svc' microservice
# as a Docker image. Do not forget to add that server's name to your
# computer's (server's,..) /etc/hosts (or whatever DNS you prefer):
#   172.25.70.26   s7026

# This is just a memo..
# Class A: 10.0.0.0 - 10.255.255.255
# Class B: 172.16.0.0 - 172.31.255.255
# Class C: 192.168.0.0 - 192.168.255.255

# NOTE!
# IP_ADDR and HOSTNAME were used to create a self-signed certificate.
# If you change IP or hostname, you have to create new 'service.pem' and
# 'service.key'.

IMAGE=ms-vote-svc-alpine
VERSION=latest
NETWORK=net_17216
IP_ADDR=172.16.70.26
HOSTNAME=s7026
CONTAINER=vote-service

# echo $NETWORK
# echo $IP_ADDR
# echo $HOSTNAME
# echo $BACKUP

# -d, --detach : Run container in background and print container ID;
# --name : Assign a name to the container;

docker run --rm --net $NETWORK --ip $IP_ADDR --hostname $HOSTNAME \
  --add-host=ws4:192.168.10.94 \
  --name $CONTAINER \
  -e GOMEMLIMIT=2147483648 \
  -e ALLOW_ORIGINS=* \
  -v /home/alex/tmp/proj2/go/vote_svc/scripts2/service.pem:/etc/x509/https/service.pem \
  -v /home/alex/tmp/proj2/go/vote_svc/scripts2/service.key:/etc/x509/https/service.key \
  -p 49177:8080 \
  -p 49178:8081 \
  -p 49179:8443 \
  -d $IMAGE:$VERSION

# --- END ---
