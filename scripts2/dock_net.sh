#!/bin/bash

#  Created : 2023-Mar-26
# Modified : 2024-Apr-04

# This shell script can be used to create a bridge network for a docker image.

# WARNING! Before trying make sure it does not exist!

# This is just a memo..
# Class A: 10.0.0.0 - 10.255.255.255
# Class B: 172.16.0.0 - 172.31.255.255
# Class C: 192.168.0.0 - 192.168.255.255

NETWORK=net_17216
IP_RANGE=172.16.0.0/16

docker network create --subnet=$IP_RANGE $NETWORK

# --- END ---
