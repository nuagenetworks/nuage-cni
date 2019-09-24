#!/bin/bash

set -x

USER_ID=$(id -u)
GROUP_ID=$(id -g)
DOCKERFILE="Dockerfile.build"
PWD=$(pwd)

MAKE_TARGET=${1}

docker run --privileged \
    --rm \
    -e NUAGE_PROJECT=${NUAGE_PROJECT} \
    -e NUAGE_BUILD_RELEASE=${NUAGE_BUILD_RELEASE} \
    -e NUAGE_BUILD_NUMBER=${NUAGE_BUILD_NUMBER} \
    -e version=${NUAGE_PROJECT}.${NUAGE_BUILD_RELEASE} \
    -e GOPATH=/BUILD/go \
    -e USER_ID=${USER_ID} \
    -e GROUP_ID=${GROUP_ID} \
    -v ${PWD}:/BUILD/go/src/github.com/nuagenetworks/nuage-cni \
    -v /usr/global:/usr/global \
    -v /root:/root \
    -v /var/run/docker.sock:/var/run/docker.sock \
    registry.mv.nuagenetworks.net:5000/build/nuage-cni \
    make ${MAKE_TARGET}
