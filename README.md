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

## Removing a node from the cluster

Go to any node and execute the following:

```shell
lxc exec mk8s-0 -- /home/ubuntu/etcd-embedded remove http://10.124.129.137:2379 mk8s-2 /home/ubuntu/config.yaml
```

Arguments are:
* the endpoint of the leader node
* the name of the node you wish to remove
* the configuration file

Inspect the member list using `etcdctl`

```shell
./etcdctl --endpoints=10.124.129.137:2379,10.124.129.151:2379,10.124.129.8:2379 endpoint status -w table
{"level":"warn","ts":"2021-09-12T11:34:38.345+0800","logger":"etcd-client","caller":"v3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"etcd-endpoints://0xc0000d2a80/#initially=[10.124.129.137:2379;10.124.129.151:2379;10.124.129.8:2379]","attempt":0,"error":"rpc error: code = DeadlineExceeded desc = latest balancer error: last connection error: connection error: desc = \"transport: Error while dialing dial tcp 10.124.129.8:2379: connect: connection refused\""}
Failed to get the status of endpoint 10.124.129.8:2379 (context deadline exceeded)
+---------------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------+
|      ENDPOINT       |        ID        | VERSION | DB SIZE | IS LEADER | IS LEARNER | RAFT TERM | RAFT INDEX | RAFT APPLIED INDEX | ERRORS |
+---------------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------+
| 10.124.129.137:2379 | 58a3d8673d099781 |   3.5.0 |   20 kB |      true |      false |         2 |         11 |                 11 |        |
| 10.124.129.151:2379 | 71c4859c30bb0e6e |   3.5.0 |   20 kB |     false |      false |         2 |         11 |                 11 |        |
+---------------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------+

```

The error you see above is normal, since the remove node is automatically goes down and it is no longer listening to the ports.


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