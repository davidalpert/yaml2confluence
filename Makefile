.PHONY: all
all: test
.SILENT:local-libjq
PWD=$(shell pwd)
export CGO_ENABLED=1
export CGO_CFLAGS="-I$(PWD)/libjq/include"
export CGO_LDFLAGS="-L$(PWD)/libjq/lib"

build-libjq:
	./scripts/build-libjq-go.sh 
test: local-libjq
	go test ./...
build:
	go build -a -installsuffix cgo y2c.go
