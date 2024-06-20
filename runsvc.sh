#!/bin/bash
#  Created : 2023-Jan-12
# Modified : 2024-Apr-27

# Soft memory limit (bytes)
# (0.5 GB / 1.0 GB / 1.5 GB / 2.0 GB)
# NOTE! It only works with Go v1.19 or newer
# ------------------------------------------
# export GOMEMLIMIT=536870912
# export GOMEMLIMIT=1073741824
# export GOMEMLIMIT=1610612736
export GOMEMLIMIT=2147483648

# This can be a comma-separated list of allowed origins.

export ALLOW_ORIGINS=https://ws4

./vote-svc -service-cert "./cert/service.pem" -service-key "./cert/service.key" -rate-limit 500

# --- END ---
