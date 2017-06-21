#!/bin/bash

set -e

if [ -z ${GOPATH} ]; then
    echo "\"GOPATH\" environmental variable is not set";
    exit 1
fi

if [ -z ${version} ]; then
    echo "\"version\" environmental variable is not set";
    exit 1
fi

cd $GOPATH/src/nuage-cni
make
sudo docker build -t nuage/cni:${version} .
docker save nuage/cni:${version} > nuage-cni-docker.tar
