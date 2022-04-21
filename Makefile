.PHONY: all
all: test
.SILENT:local-libjq
PWD=$(shell pwd)
export CGO_ENABLED=1
export CGO_CFLAGS="-I$(PWD)/libjq/include"
export CGO_LDFLAGS="-L$(PWD)/libjq/lib"

local-libjq:
	if [ ! -d "./libjq" ]; then \
		wget -q https://github.com/flant/libjq-go/releases/download/jq-b6be13d5-0/libjq-glibc-amd64.tgz; \
		tar zxf libjq-glibc-amd64.tgz; \
		rm -f libjq-glibc-amd64.tgz; \
	fi
test: local-libjq
	go test ./...
build: local-libjq
	go build y2c.go
