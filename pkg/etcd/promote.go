package etcd

import (
	"context"
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func Promote(leaderEndpoint string, memberId uint64) (*clientv3.MemberPromoteResponse, error) {
	var endpoints = make([]string, 1)
	endpoints[0] = leaderEndpoint
	cli := NewEtcdClient(endpoints)
	defer cli.Close()

	resp, err := cli.MemberPromote(context.Background(), memberId)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return resp, nil
}
