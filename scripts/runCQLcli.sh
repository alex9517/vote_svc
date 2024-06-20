#!/bin/bash

#  Created : 2023-Nov-14
# Modified : 2023-Nov-23

# This is a shell script to run the Apache Cassandra CQLSH client
# as a Docker image. Do not forget to add that server name to your
# computer's (server's,..) /etc/hosts (or whatever DNS you prefer):
#   172.25.70.27   s7027

IMAGE=cassandra
# VERSION=latest
NETWORK=net_17216
# IP_ADDR=172.16.70.31
HOSTNAME=s7031
DATA_PATH=/scripts/data.cql
# NAME=cassandra-cqlsh
# echo $NETWORK
# echo $IP_ADDR
# echo $HOSTNAME

# -d, --detach : Run container in background and print container ID;
# --name : Assign a name to the container;
#  -e ALLOW_ORIGINS=https://ws4,http://localhost,https://localhost,http://localhost:3000 \

# docker run -it --network some-network --rm cassandra cqlsh some-cassandra
# docker run --rm --network cassandra -v "$(pwd)/data.cql:/scripts/data.cql" -e CQLSH_HOST=cassandra -e CQLSH_PORT=9042 -e CQLVERSION=3.4.6 nuvo/docker-cqlsh


# docker run -it --net $NETWORK --rm $IMAGE cqlsh $HOSTNAME
docker run -it --net $NETWORK -v "$(pwd)/data.cql:$DATA_PATH" --rm $IMAGE cqlsh $HOSTNAME

# --- END ---
