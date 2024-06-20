#!/bin/bash
#  Created : 2022-Oct-20
# Modified : 2023-Jan-29

# This script creates a
#
#  SELF-SIGNED (or call it "fake")
#
#  SERVER CERTIFICATE signed by a
#
#  ROOT CERTIFICATE of a self-signed (i.e. "fake") CA.

# Expiration time ~5 years
#  (it's too much .. but it's just for development).

EXPIRE_AFTER=1827

# Create a certificate key for the service.

openssl genrsa -out service.key 4096

# Create signing CSR (for local testing you can use `'/CN=localhost'`,
#  for online testing `CN` needs to be replaced with something like
#  `'/CN=grpc.example.com'`. Include this in a config file
#  ([certificate.conf](certificate.conf)).

openssl req -new -key service.key -out service.csr -config certificate.conf

# Generate a certificate for the service.

openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial -out service.pem -days $EXPIRE_AFTER -sha256 -extfile certificate.conf -extensions req_ext

# Verify (optional).

# openssl x509 -in service.pem -text -noout
# openssl verify -CAfile ca.cert server.pem

# -END-
