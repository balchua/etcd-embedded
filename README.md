## Simulate a clustered etcd

Pre-requisite
* etcdctl
* Must have these available ports `2379`, `2380`, `12379`, `12380`, `22379`, `22380`



In this example, we will start a 3 node etcd cluster.  Each node's configuration is in the following files:

* First node `config.yaml`
* Second node `config2.yaml`
* Third node `config3.yaml`

Check the content of each yaml and make sure that the ports do not interfere with your existing system.
Example, the path where etcd stores its data, currently hardcoded to the following directories

```yaml
# Path to the data directory.
data-dir: /home/thor/etcd-data/data

# Path to the dedicated wal directory.
wal-dir: /home/thor/etcd-data/data
```

## First start the main server

```shell
./etcd-embedded server ./config.yaml
```

## Add the second node as a learner.

```shell
./etcd-embedded join http://localhost:2379 ./config2.yaml
```

Before starting the second node, check  that the node is added as a learner.

```shell
./etcdctl member list --endpoints=localhost:2379,localhost:12379,localhost:22379 -w table
+------------------+-----------+------+------------------------+-----------------------+------------+
|        ID        |  STATUS   | NAME |       PEER ADDRS       |     CLIENT ADDRS      | IS LEARNER |
+------------------+-----------+------+------------------------+-----------------------+------------+
| 8e9e05c52164694d |   started |   n0 |  http://localhost:2380 | http://localhost:2379 |      false |
| 97bbfd74ce4e6ce4 | unstarted |      | http://localhost:12380 |                       |       true |
+------------------+-----------+------+------------------------+-----------------------+------------+

```

Then start the second node

```shell
./etcd-embedded server ./config2.yaml
```

Wait for a few seconds, the promotion to a voter node is automatic.

```shell
./etcdctl member list --endpoints=localhost:2379,localhost:12379,localhost:22379 -w table
+------------------+---------+------+------------------------+------------------------+------------+
|        ID        | STATUS  | NAME |       PEER ADDRS       |      CLIENT ADDRS      | IS LEARNER |
+------------------+---------+------+------------------------+------------------------+------------+
| 8e9e05c52164694d | started |   n0 |  http://localhost:2380 |  http://localhost:2379 |      false |
| 97bbfd74ce4e6ce4 | started |   n1 | http://localhost:12380 | http://localhost:12379 |      false |
+------------------+---------+------+------------------------+------------------------+------------+
```

## Add the third node as a learner.

```shell
./etcd-embedded join http://localhost:2379 ./config3.yaml
```

Before starting the third node, check  that the node is added as a learner.

```shell
./etcdctl member list --endpoints=localhost:2379,localhost:12379,localhost:22379 -w table
+------------------+-----------+------+------------------------+------------------------+------------+
|        ID        |  STATUS   | NAME |       PEER ADDRS       |      CLIENT ADDRS      | IS LEARNER |
+------------------+-----------+------+------------------------+------------------------+------------+
| 278f416e9d2430e9 | unstarted |      | http://localhost:22380 |                        |       true |
| 8e9e05c52164694d |   started |   n0 |  http://localhost:2380 |  http://localhost:2379 |      false |
| 97bbfd74ce4e6ce4 |   started |   n1 | http://localhost:12380 | http://localhost:12379 |      false |
+------------------+-----------+------+------------------------+------------------------+------------+
```

Then start the third node

```shell
./etcd-embedded server ./config3.yaml
```

```shell
./etcdctl member list --endpoints=localhost:2379,localhost:12379,localhost:22379 -w table
+------------------+---------+------+------------------------+------------------------+------------+
|        ID        | STATUS  | NAME |       PEER ADDRS       |      CLIENT ADDRS      | IS LEARNER |
+------------------+---------+------+------------------------+------------------------+------------+
| 278f416e9d2430e9 | started |   n2 | http://localhost:22380 | http://localhost:22379 |      false |
| 8e9e05c52164694d | started |   n0 |  http://localhost:2380 |  http://localhost:2379 |      false |
| 97bbfd74ce4e6ce4 | started |   n1 | http://localhost:12380 | http://localhost:12379 |      false |
+------------------+---------+------+------------------------+------------------------+------------+
```

Now you have a working etcd HA cluster.

## etcdctl helpful commands

```shell
./etcdctl member list --endpoints=localhost:2379,localhost:12379,localhost:22379
```

Check the status of the endpoints, shows who is the leader
```shell
./etcdctl --endpoints=localhost:2379,localhost:12379,localhost:22379 -w table endpoint status
```