#!/bin/sh

# Script to install Nuage CNI on a Kubernetes/Openshift host.
# - Expects the host CNI binary path to be mounted at /host/opt/cni/bin.
# - Expects the host CNI network config path to be mounted at /host/etc/cni/net.d.
# - Expects the desired Nuage VSP config in the NUAGE_VSP_CONFIG env variable.

# The directory on the host where CNI networks are installed. Defaults to
# /etc/cni/net.d, but can be overridden by setting CNI_NET_DIR.  This is used
# for populating absolute paths in the CNI network config to assets
# which are installed in the CNI network config directory.
HOST_CNI_NET_DIR=${CNI_NET_DIR:-/etc/cni/net.d}

# Clean up any existing Nuage CNI binary, network config and yaml file
rm -irf /host/opt/cni
rm -irf /host/etc/cni

# Create directory structure to be created on the host for
# CNI plugin binary and CNI netconf file under /etc and /opt
mkdir -p /host/etc/cni/net.d
mkdir -p /host/opt/cni/bin
chmod 755 /host/etc/cni/net.d
chmod 755 /host/opt/cni/bin

# Choose which default cni binaries should be copied
SKIP_CNI_BINARIES=${SKIP_CNI_BINARIES:-""}
SKIP_CNI_BINARIES=",$SKIP_CNI_BINARIES,"

if [ "$1" = "nuage-cni-k8s" ]; then
    TMP_CONF='/tmp/vsp-k8s.yaml'
    CONFIG_DIR='vsp-k8s'
    rm -irf /usr/share/vsp-k8s
    mkdir -p /usr/share/vsp-k8s
    chmod 755 /usr/share/vsp-k8s
    NUAGE_CONF='/usr/share/vsp-k8s/vsp-k8s.yaml'
fi

if [ "$1" = "nuage-cni-openshift" ]; then
    TMP_CONF='/tmp/vsp-openshift.yaml'
    CONFIG_DIR='vsp-openshift'
    NUAGE_CONF='/usr/share/vsp-openshift/vsp-openshift.yaml'
fi

if [ "$2" = "is_atomic" ]; then
    DIR='/var/usr/share'
    NUAGE_CONF=$DIR/$CONFIG_DIR/$CONFIG_DIR.yaml
fi

# Create a temporary file for Nuage vsp-k8s yaml
cat > $TMP_CONF << EOF
EOF
chmod 777 $TMP_CONF

# If specified, overwrite the network configuration file.
if [ "${NUAGE_VSP_CONFIG:-}" != "" ]; then
cat >$TMP_CONF <<EOF
${NUAGE_VSP_CONFIG:-}
EOF
fi

# This is to create the VSP yaml file on RHEL masters
# on atomic setups
if [ "$2" = "is_atomic" ]; then
    cp $TMP_CONF /usr/share/vsp-openshift/
fi

# Copy the generated yaml file to usr/share Nuage folder
mv $TMP_CONF $NUAGE_CONF

# Place the new binaries if the directory is writeable.
if [ -w "/host/opt/cni/bin/" ]; then
	for path in /opt/cni/bin/*
	do
		filename=$(basename $path)
		tmp=",$filename,"
		if [ "${SKIP_CNI_BINARIES#*$tmp}" = "$SKIP_CNI_BINARIES" ] && [ ! -f /host/opt/cni/bin/$filename ]; then
			cp /opt/cni/bin/$filename /host/opt/cni/bin/
		fi
	done
	echo "Wrote Nuage CNI binaries to /host/opt/cni/bin/"
fi

if [ "$1" = "nuage-cni-k8s" ]; then
    rm -f /host/opt/cni/bin/nuage-cni-openshift
else
    rm -f /host/opt/cni/bin/nuage-cni-k8s
fi

CNI_YAML_CONF='/etc/default/nuage-cni.yaml'
# Configuring Nuage CNI yaml file using daemon sets
if [ "${NUAGE_CNI_YAML_CONFIG:-}" != "" ]; then
cat > $CNI_YAML_CONF << EOF
EOF
chmod 777 $CNI_YAML_CONF

cat > $CNI_YAML_CONF <<EOF
${NUAGE_CNI_YAML_CONFIG:-}
EOF
fi
cp $CNI_YAML_CONF /host/etc/default

TMP_CONF='/nuage-net.conf.k8s'
if [ "$1" = "nuage-cni-openshift" ]; then
    TMP_CONF='/nuage-net.conf.openshift'
fi

# Move the temporary CNI config into place.
FILENAME=${CNI_CONF_NAME:-nuage-net.conf}
mv $TMP_CONF /host/etc/cni/net.d/${FILENAME}
echo "Wrote CNI config: $(cat /host/etc/cni/net.d/${FILENAME})"
echo "Done configuring CNI"

# Add iptables rule for Nuage overlay to underlay traffic
# and vice versa
iptables -L | grep nuage-vxlan
if [ $? -eq 0 ]
then
echo "iptables rule to allow Nuage vxlan ports is present"
else
iptables -w -I INPUT 1 -p udp --dport 4789 -j ACCEPT -m comment --comment "nuage-vxlan"
fi

iptables -L | grep nuage-overlay-underlay
if [ $? -eq 0 ]
then
echo "iptables rule to allow Nuage overlay to underlay traffic is present"
else
iptables -w -I FORWARD 1 -s ${NUAGE_CLUSTER_NW_CIDR:-} -j ACCEPT -m comment --comment "nuage-overlay-underlay"
fi

iptables -L | grep nuage-underlay-overlay
if [ $? -eq 0 ]
then
echo "iptables rule to allow Nuage underlay to overlay traffic is present"
else
iptables -w -I FORWARD 1 -d ${NUAGE_CLUSTER_NW_CIDR:-} -j ACCEPT -m comment --comment "nuage-underlay-overlay"  
fi

if [ "$1" = "nuage-cni-k8s" ]; then
# Create Nuage kubeconfig file for api server communication
cat > /usr/share/vsp-k8s/nuage.kubeconfig <<EOF
apiVersion: v1
kind: Config
current-context: nuage-to-cluster.local
preferences: {}
clusters:
- cluster:
    insecure-skip-tls-verify: true
    server: ${MASTER_API_SERVER_URL:-}
  name: cluster.local
contexts:
- context:
    cluster: cluster.local
    user: nuage
  name: nuage-to-cluster.local
users:
- name: nuage
  user:
    token: ${NUAGE_TOKEN:-}
EOF
fi

# Start Nuage CNI audit daemon to run infinitely here.
# This prevents Kubernetes from restarting the pod repeatedly.
/opt/cni/bin/$1 -daemon
