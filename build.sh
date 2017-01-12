#!/bin/bash

me=$(cd ${0%/*}; pwd)

set -ex

docker run --rm -it -v "${me}":/app -w /app golang:1.7-alpine sh -c 'apk update && apk add git && go get gopkg.in/gin-gonic/gin.v1 && CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o okgo'
docker build -t rgmccaw/okgo ${me}/
