#!/bin/bash

me=$(cd ${0%/*}; pwd)

set -ex

docker run --rm -it -v "${me}":/app -w /app golang:1.4.2 sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o okgo'
docker build -t rgmccaw/okgo ${me}/
