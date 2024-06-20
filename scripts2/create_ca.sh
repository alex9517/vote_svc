#!/bin/bash
#  Created : 2022-Oct-20
# Modified : 2022-Oct-20

# This script creates a self-signed
#  (or call it "fake") Certification Authority (CA).

# Expiration time ~5 years
#  (it's too much .. but I'm lazy, .. and it's just for development).

EXPIRE_AFTER=1827

# Create Root signing Key:

openssl genrsa -out ca.key 4096

# Generate self-signed Root certificate

openssl req -new -x509 -key ca.key -sha256 -subj "/C=EU/ST=Europe/O=Inter Soft Dev" -days $EXPIRE_AFTER -out ca.cert

# -END-
