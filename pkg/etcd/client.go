package etcd

import (
	"fmt"
	"os"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewEtcdClient(leaderEndpoint string) *clientv3.Client {
	var endpoints = make([]string, 1)
	endpoints[0] = leaderEndpoint

	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return cli
}
