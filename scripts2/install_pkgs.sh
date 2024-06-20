#!/bin/bash
#  Created : 2022-Dec-17
# Modified : 2024-Apr-04

# IMPORTANT!
# Run this in the root dir of the project, i.e. where 'go.mod' is located.

# IMPORTANT!
# Probably you don't need all this packages, check carefully and install only what you really need.

# This script installs the external dependencies (3rd party packages) required by this proj. Of course, you can do it manually,
# using this script as a source for 'go get ...' commands, or, you can edit this script to add/skip/remove whatever you need.

# Go kit projects usually don't need a complicated config! So, skip 'spf13/viper';
# go get github.com/spf13/viper

go get github.com/go-kit/kit/log

# go get google.golang.org/grpc

go get github.com/go-kit/kit/metrics/prometheus

go get github.com/lightstep/lightstep-tracer-go

go get github.com/oklog/oklog/pkg/group

go get github.com/opentracing/opentracing-go

go get github.com/openzipkin-contrib/zipkin-go-opentracing

go get github.com/openzipkin/zipkin-go

go get github.com/openzipkin/zipkin-go/reporter/http

go get github.com/prometheus/client_golang/prometheus

go get github.com/prometheus/client_golang/prometheus/promhttp

go get sourcegraph.com/sourcegraph/appdash

go get sourcegraph.com/sourcegraph/appdash/opentracing

# go get github.com/disintegration/imaging

go get github.com/rs/cors

go get github.com/hashicorp/consul/api

# go get github.com/go-kit/kit/sd/consul

# go get github.com/gorilla/mux

# go get github.com/oklog/run


# Update protoc Go bindings via
#  go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

go get github.com/prometheus/client_golang@v1.19.0

go get github.com/hashicorp/consul/api@v1.18.1

# --- END OF FILE ---
