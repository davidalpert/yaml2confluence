#!/bin/bash

CGO_ENABLED=1 \
CGO_CFLAGS="-I`pwd`/libjq/include" \
CGO_LDFLAGS="-L`pwd`/libjq/lib" \
go run cmd/y2c/main.go $@