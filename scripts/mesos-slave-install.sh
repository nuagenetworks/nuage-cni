#!/usr/bin/env bash

MASTER_IP=$1

sudo rpm -Uvh http://repos.mesosphere.io/el/7/noarch/RPMS/mesosphere-el-repo-7-1.noarch.rpm
sudo yum -y install mesos
echo [dockerrepo] >> /etc/yum.repos.d/docker.repo
echo 'name=Docker Repository' >> /etc/yum.repos.d/docker.repo
echo 'baseurl=https://yum.dockerproject.org/repo/main/centos/7/' >> /etc/yum.repos.d/docker.repo
echo 'enabled=1' >> /etc/yum.repos.d/docker.repo
echo 'gpgcheck=1' >> /etc/yum.repos.d/docker.repo
echo 'gpgkey=https://yum.dockerproject.org/gpg' >> /etc/yum.repos.d/docker.repo
sudo yum -y install docker-engine
sudo systemctl enable docker.service
sudo systemctl start docker
echo OPTIONS='--selinux-enabled --log-driver=journald' >> /etc/sysconfig/docker
echo DOCKER_CERT_PATH=/etc/docker >> /etc/sysconfig/docker
echo HTTP_PROXY=http://global.proxy.alcatel-lucent.com:8000 >> /etc/sysconfig/docker
echo HTTPS_PROXY=https://global.proxy.alcatel-lucent.com:8000 >> /etc/sysconfig/docker
mkdir -p /etc/systemd/system/docker.service.d
echo [Service] >> /etc/systemd/system/docker.service.d/http-proxy.conf
echo 'Environment="HTTP_PROXY=http://global.proxy.alcatel-lucent.com:8000" "HTTPS_PROXY=https://global.proxy.alcatel-lucent.com:8000"' >> /etc/systemd/system/docker.service.d/http-proxy.conf
echo 'docker,mesos' > /etc/mesos-slave/containerizers
echo '5mins' > /etc/mesos-slave/executor_registration_timeout
echo zk://$MASTER_IP:2181/mesos > /etc/mesos/zk
systemctl disable mesos-master
sudo service mesos-slave restart
sudo systemctl restart docker
systemctl daemon-reload
docker pull busybox
