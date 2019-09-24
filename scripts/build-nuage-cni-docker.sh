#!/bin/bash

set -xe

if [ -z ${GOPATH} ]; then
    echo "\"GOPATH\" environmental variable is not set";
    exit 1
fi

if [ -z ${version} ]; then
    echo "\"version\" environmental variable is not set";
    exit 1
fi

export DOCKER_API_VERSION=1.22
NUAGE_BUILD_NUMBER=${NUAGE_BUILD_NUMBER:-0}
version=$version-$NUAGE_BUILD_NUMBER
cd $GOPATH/src/github.com/nuagenetworks/nuage-cni
make -f Makefile
sudo docker build -t nuage/cni:${version} .
docker save nuage/cni:${version} > nuage-cni-docker-${version}.tar
docker rmi nuage/cni:${version}
