#!/usr/bin/env bash

MASTER_IP=$1

sudo rpm -Uvh http://repos.mesosphere.io/el/7/noarch/RPMS/mesosphere-el-repo-7-1.noarch.rpm
sudo yum -y install mesos marathon
sudo yum -y install mesosphere-zookeeper
echo 1 > /var/lib/zookeeper/myid
echo $MASTER_IP > /etc/mesos-master/hostname
echo server.1=$MASTER_IP:2888:3888 >> /etc/zookeeper/conf/zoo.cfg
sudo systemctl start zookeeper
echo zk://$MASTER_IP:2181/mesos > /etc/mesos/zk
echo 1 > /etc/mesos-master/quorum
systemctl stop mesos-slave.service
sudo service mesos-master restart
sudo service marathon restart
