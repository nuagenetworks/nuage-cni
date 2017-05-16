# nuage-mesos-cni
Repository for maintaining Nuage Integration with Mesos/Kubernetes/Openshift using CNI plugin

To build CNI plugin binaries, execute "make" from "nuage-cni" folder. Make sure the GOPATH is set correctly. Following binaries will be generated under "nuage-cni" folder:

- nuage-cni-mesos
- nuage-cni-k8s
- nuage-cni-openshift

To generate Nuage CNI plugin rpm builds:

- Clone the nuage-cni repository under src folder in your GOPATH

- Set the version required for the rpm: export version=0.0

- Update the version in all rpmbuild spec files

- Then run ./scripts/buildRPM.sh which will generate RPMs for Nuage CNI plugin for Mesos/K8S/OSE
