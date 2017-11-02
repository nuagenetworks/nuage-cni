# nuage-cni

The goal of `nuage-cni` is to provide Nuage Container Network Interface (CNI) plugin and allow launched containers/pods to be attached to and detached from Nuage overlay networks. For more information on CNI refer to https://github.com/containernetworking/cni


# Nuage CNI plugin modes

The Nuage CNI plugin functions in two modes:

- CNI mode

- Audit Daemon mode

## CNI Mode

In the CNI mode, the Nuage CNI plugin gets invoked by network/cni isolator when a container gets launched or deleted.

 - When container/pod gets launched, the Nuage CNI Plugin gets invoked. It creates a container port in VRS and resolves that container with an IP address from Nuage defined VSD overlay network.

 - When container/pod gets deleted, the Nuage CNI Plugin gets invoked. It deletes the container port entry from VRS thereby detaching the container from Nuage defined VSD overlay network.

## Audit Daemon Mode

In Audit Daemon mode, the Nuage CNI plugin also operates as a background systemd service (nuage-cni) on each agent VRS node and periodically audits agent VRS nodes to make sure the ports in VRS correspond to the currently functional containers/pods. If there are any stale VRS ports which do not correspond to any currently running containers/pods, the nuage-cni service deletes those ports from VRS. nuage-cni service will be started by default on all agent VRS nodes as a part of the CNI plugin installation. To stop the audit daemon, execute `systemctl stop nuage-cni` on the agent VRS node.


# Build Nuage CNI plugin

## Steps to generate CNI plugin binaries

- Clone the nuage-cni repository under src folder in your GOPATH

- Execute "make" from "nuage-cni" folder. Following binaries will be generated under "nuage-cni" folder:

 - nuage-cni-mesos
 - nuage-cni-k8s
 - nuage-cni-openshift

## Steps to generate CNI plugin rpm packages

- Clone the nuage-cni repository under src folder in your GOPATH

- Set the version required for the rpm: export version=`desired rpm version`

- Update the `desired rpm version` in rpmbuild spec files for Mesos, k8s and Openshift

- Then run ./scripts/buildRPM.sh which will generate RPMs for Nuage CNI plugin for Mesos/K8S/OSE under rpmbuild directory on your host


## Steps to generate Nuage CNI docker image for CNI daemon sets install

- Clone https://github.com/nuagenetworks/nuage-cni.git to your $GOPATH/src folder on your host machine

- Change directory to `nuage-cni` folder

- Set desired image version for your CNI docker image using `export version=<image-version>` eg: export version=v5.1.2

- Then run the build script `./scripts/build-nuage-cni-docker.sh`

- At the end of this script execution `nuage-cni-docker-<image-version>.tar` will be generated under nuage-cni folder

- To load the CNI docker image on your slave nodes, copy the `nuage-cni-docker-<image-version>.tar` file generated above to your slave nodes and do `docker load -i nuage-cni-docker-<image-version>.tar`


# Automated Nuage CNI plugin installation

- Please refer to the Nuage VSP Mesos Integration Guide for detailed steps on Nuage Mesos CNI plugin frameworks installation on Mesos clusters.

- Please refer to the Nuage VSP OpenShift Integration Guide for detailed steps on Nuage Openshift CNI plugin Ansible installation on Openshift clusters (HA and single master)

- Please refer to the Nuage VSP Kubernetes Integration Guide for detailed steps on Nuage Kubernetes CNI plugin Ansible installation on Kubernetes setups
