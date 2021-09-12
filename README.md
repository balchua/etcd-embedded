# Embedded etcd

This project is to simulate embedding etcd with high availability.

## How to

Before you start, you to satisfy these pre-requisites
* etcdctl - https://github.com/etcd-io/etcd/releases
* lxd - This will be used to simulate multi node etcd
* Must have these available ports `2379`, `2380`

## Build

Simply run `go build .`

## Setting up the lxd

We will start a 3 node LXD instance.  Launch an LXD instance

```shell
lxc launch -p default -p microk8s ubuntu:20.04 mk8s-0
lxc launch -p default -p microk8s ubuntu:20.04 mk8s-1
lxc launch -p default -p microk8s ubuntu:20.04 mk8s-2
```

### Copy the default config and  binary lxd instance

```shell
lxc file push default-config.yaml mk8s-0/home/ubuntu/config.yaml
lxc file push default-config.yaml mk8s-1/home/ubuntu/config.yaml
lxc file push default-config.yaml mk8s-2/home/ubuntu/config.yaml

lxc file push etcd-embedded mk8s-0/home/ubuntu/etcd-embedded
lxc file push etcd-embedded mk8s-1/home/ubuntu/etcd-embedded
lxc file push etcd-embedded mk8s-2/home/ubuntu/etcd-embedded

```

### Get all lxd instance IPs

```shell
lxc list

+-----------------------+---------+-----------------------+------+-----------+-----------+
|         NAME          |  STATE  |         IPV4          | IPV6 |   TYPE    | SNAPSHOTS |
+-----------------------+---------+-----------------------+------+-----------+-----------+
| mk8s-0                | RUNNING | 10.124.129.137 (eth0) |      | CONTAINER | 0         |
+-----------------------+---------+-----------------------+------+-----------+-----------+
| mk8s-1                | RUNNING | 10.124.129.151 (eth0) |      | CONTAINER | 0         |
+-----------------------+---------+-----------------------+------+-----------+-----------+
| mk8s-2                | RUNNING | 10.124.129.8 (eth0)   |      | CONTAINER | 0         |
+-----------------------+---------+-----------------------+------+-----------+-----------+

```
## First start (`mk8s-0`) the main server

```shell
lxc exec mk8s-0 -- /home/ubuntu/etcd-embedded server /home/ubuntu/config.yaml
```

## Add the second node (`mk8s-1`) as a learner.

```shell
lxc exec mk8s-1 -- /home/ubuntu/etcd-embedded join http://10.124.129.137:2380 /home/ubuntu/config.yaml
```

Before starting the second node (`mk8s-1`), check  that the node is added as a learner.

```shell

./etcdctl --endpoints=10.124.129.137:2379,10.124.129.151:2379,10.124.129.8:2379 member list -w table
+------------------+-----------+--------+----------------------------+----------------------------+------------+
|        ID        |  STATUS   |  NAME  |         PEER ADDRS         |        CLIENT ADDRS        | IS LEARNER |
+------------------+-----------+--------+----------------------------+----------------------------+------------+
| 58a3d8673d099781 |   started | mk8s-0 | http://10.124.129.137:2380 | http://10.124.129.137:2379 |      false |
| dd94fc49fe0983b8 | unstarted |        | http://10.124.129.151:2379 |                            |       true |
+------------------+-----------+--------+----------------------------+----------------------------+------------+

```

Then start the second node (`mk8s-1`)

```shell
lxc exec mk8s-1 -- /home/ubuntu/etcd-embedded server /home/ubuntu/config.yaml
```

Wait for a few seconds, the promotion to a voter node is automatic.

```shell
./etcdctl --endpoints=10.124.129.137:2379,10.124.129.151:2379,10.124.129.8:2379 member list -w table
+------------------+---------+--------+----------------------------+----------------------------+------------+
|        ID        | STATUS  |  NAME  |         PEER ADDRS         |        CLIENT ADDRS        | IS LEARNER |
+------------------+---------+--------+----------------------------+----------------------------+------------+
| 58a3d8673d099781 | started | mk8s-0 | http://10.124.129.137:2380 | http://10.124.129.137:2379 |      false |
| dd94fc49fe0983b8 | started | mk8s-1 | http://10.124.129.151:2379 | http://10.124.129.151:2379 |      false |
+------------------+---------+--------+----------------------------+----------------------------+------------+

```

## Add the third node (`mk8s-2`) as a learner.

```shell
lxc exec mk8s-2 -- /home/ubuntu/etcd-embedded join http://10.124.129.137:2380 /home/ubuntu/config.yaml
{"level":"info","ts":1631414981.2671242,"caller":"util/network.go:21","msg":"Interface is a loopback","name":"lo"}
{"level":"info","ts":1631414981.2682645,"caller":"util/network.go:30","msg":"Interface details","name":"eth0","isUp":true}
{"level":"info","ts":1631414981.2686086,"caller":"util/network.go:34","msg":"Interface name"}
{"level":"info","ts":1631414981.2686648,"caller":"util/network.go:37","msg":"IP found","ip":"10.124.129.8"}
{"level":"info","ts":1631414981.2689445,"caller":"etcd/learner.go:16","msg":"AddMemberAsLearner","peerUrl":"http://10.124.129.8:2379","LeaderEndpoint":"http://10.124.129.137:2380"}
{"level":"info","ts":1631414981.2927167,"caller":"cmd/join.go:41","msg":"Join Successful.","PeerURL":["http://10.124.129.8:2379"],"MemberId":2288430810285401478,"ClientURL":[],"IsLearner":true}

```

Before starting the third node (`mk8s-2`), check that the node is added as a learner.

```shell
./etcdctl --endpoints=10.124.129.137:2379,10.124.129.151:2379,10.124.129.8:2379 member list -w table
+------------------+-----------+--------+----------------------------+----------------------------+------------+
|        ID        |  STATUS   |  NAME  |         PEER ADDRS         |        CLIENT ADDRS        | IS LEARNER |
+------------------+-----------+--------+----------------------------+----------------------------+------------+
| 1fc223b2841ee586 | unstarted |        |   http://10.124.129.8:2379 |                            |       true |
| 58a3d8673d099781 |   started | mk8s-0 | http://10.124.129.137:2380 | http://10.124.129.137:2379 |      false |
| dd94fc49fe0983b8 |   started | mk8s-1 | http://10.124.129.151:2379 | http://10.124.129.151:2379 |      false |
+------------------+-----------+--------+----------------------------+----------------------------+------------+

```

Then start the third node (`mk8s-2`)

```shell
lxc exec mk8s-2 -- /home/ubuntu/etcd-embedded server /home/ubuntu/config.yaml
```

Check that the third node (`mk8s-2`) is all part of the cluster

```shell
./etcdctl --endpoints=10.124.129.137:2379,10.124.129.151:2379,10.124.129.8:2379 member list -w table
+------------------+---------+--------+----------------------------+----------------------------+------------+
|        ID        | STATUS  |  NAME  |         PEER ADDRS         |        CLIENT ADDRS        | IS LEARNER |
+------------------+---------+--------+----------------------------+----------------------------+------------+
| 1fc223b2841ee586 | started | mk8s-2 |   http://10.124.129.8:2379 |   http://10.124.129.8:2379 |      false |
| 58a3d8673d099781 | started | mk8s-0 | http://10.124.129.137:2380 | http://10.124.129.137:2379 |      false |
| dd94fc49fe0983b8 | started | mk8s-1 | http://10.124.129.151:2379 | http://10.124.129.151:2379 |      false |
+------------------+---------+--------+----------------------------+----------------------------+------------+

```

Now you have a working etcd HA cluster.


## Endpoint

The app tries to add simple endpoint to show the etcd member nodes.

`http://localhost:3000/members`


## etcdctl helpful commands

```shell
./etcdctl --endpoints=10.124.129.137:2379,10.124.129.151:2379,10.124.129.8:2379 member list -w table
+------------------+---------+--------+----------------------------+----------------------------+------------+
|        ID        | STATUS  |  NAME  |         PEER ADDRS         |        CLIENT ADDRS        | IS LEARNER |
+------------------+---------+--------+----------------------------+----------------------------+------------+
| 1fc223b2841ee586 | started | mk8s-2 |   http://10.124.129.8:2379 |   http://10.124.129.8:2379 |      false |
| 58a3d8673d099781 | started | mk8s-0 | http://10.124.129.137:2380 | http://10.124.129.137:2379 |      false |
| dd94fc49fe0983b8 | started | mk8s-1 | http://10.124.129.151:2379 | http://10.124.129.151:2379 |      false |
+------------------+---------+--------+----------------------------+----------------------------+------------+
```

Check the status of the endpoints, shows who is the leader

```shell
./etcdctl --endpoints=10.124.129.137:2379,10.124.129.151:2379,10.124.129.8:2379 endpoint status -w table
+---------------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------+
|      ENDPOINT       |        ID        | VERSION | DB SIZE | IS LEADER | IS LEARNER | RAFT TERM | RAFT INDEX | RAFT APPLIED INDEX | ERRORS |
+---------------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------+
| 10.124.129.137:2379 | 58a3d8673d099781 |   3.5.0 |   20 kB |      true |      false |         2 |         10 |                 10 |        |
| 10.124.129.151:2379 | dd94fc49fe0983b8 |   3.5.0 |   20 kB |     false |      false |         2 |         10 |                 10 |        |
|   10.124.129.8:2379 | 1fc223b2841ee586 |   3.5.0 |   20 kB |     false |      false |         2 |         10 |                 10 |        |
+---------------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------+

```