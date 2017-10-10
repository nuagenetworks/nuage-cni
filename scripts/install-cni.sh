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
rm -f /host/opt/cni/bin/$1
rm -f /host/etc/cni/net.d/nuage-net.conf

# Choose which default cni binaries should be copied
SKIP_CNI_BINARIES=${SKIP_CNI_BINARIES:-""}
SKIP_CNI_BINARIES=",$SKIP_CNI_BINARIES,"

if [ "$1" = "nuage-cni-k8s" ]; then
    TMP_CONF='/tmp/vsp-k8s.yaml'
    NUAGE_CONF='/usr/share/vsp-k8s/vsp-k8s.yaml'
    CONFIG_DIR='vsp-k8s'
fi

if [ "$1" = "nuage-cni-openshift" ]; then
    TMP_CONF='/tmp/vsp-openshift.yaml'
    NUAGE_CONF='/usr/share/vsp-openshift/vsp-openshift.yaml'
    CONFIG_DIR='vsp-openshift'
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

# Deleting any pre-existing iptable mangle forward rule to
# mark underlay packets to support node port
iptables -t mangle -D FORWARD -o svc-pat-tap -p tcp -m mark ! --mark 0x2 -j MARK --set-xmark 0xd/0xd
while ! test -f "/var/run/openvswitch/db.sock"; do
  sleep 10
  echo "Waiting for VRS container to come up"
done

# Adding iptable mangle rule to mark the underlay packets
# to support node port functionality
sleep 2
iptables -t mangle -A FORWARD -o svc-pat-tap -p tcp -m mark ! --mark 0x2 -j MARK --set-xmark 0xd/0xd

# Start Nuage CNI audit daemon to run infinitely here.
# This prevents Kubernetes from restarting the pod repeatedly.
/opt/cni/bin/$1 -daemon
