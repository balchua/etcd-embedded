package etcd

import (
	"fmt"
	"os"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewEtcdClient(endpoints []string) *clientv3.Client {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return cli
}
