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

rm -rf ~/rpmbuild/BUILD/nuage*
rm -rf ~/rpmbuild/SOURCES/nuage*
rm -rf ~/rpmbuild/RPMS/x86_64/nuage*
rm -rf ~/rpmbuild/SRPMS/nuage*
rm -rf /tmp/nuage-cni*

cd $GOPATH/src/github.com/nuagenetworks/nuage-cni
go build -o nuage-cni.go
go build -o nuage-cni-k8s nuage-cni.go
go build -o nuage-cni-openshift nuage-cni.go

cd /tmp
cp -r $GOPATH/src/github.com/nuagenetworks/nuage-cni nuage-cni-k8s-${version}
cp -r $GOPATH/src/github.com/nuagenetworks/nuage-cni nuage-cni-openshift-${version}
tar -czvf $HOME/rpmbuild/SOURCES/nuage-cni-k8s-${version}.tar.gz nuage-cni-k8s-${version}
tar -czvf $HOME/rpmbuild/SOURCES/nuage-cni-openshift-${version}.tar.gz nuage-cni-openshift-${version}
rpmbuild -ba $GOPATH/src/github.com/nuagenetworks/nuage-cni/rpmbuild/nuage-cni-k8s.spec
rpmbuild -ba $GOPATH/src/github.com/nuagenetworks/nuage-cni/rpmbuild/nuage-cni-openshift.spec
