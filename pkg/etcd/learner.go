package etcd

import (
	"context"
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func AddMemberAsLearner(leaderEndpoint string, peerUrl string) (*clientv3.MemberAddResponse, error) {
	var endpoints = make([]string, 1)
	endpoints[0] = leaderEndpoint
	cli := NewEtcdClient(endpoints)
	defer cli.Close()
	var peerUrls []string

	peerUrls = make([]string, 1)
	peerUrls[0] = peerUrl
	log.Printf("Peer URL: %s, Leader: %s", peerUrl, leaderEndpoint)
	resp, err := cli.MemberAddAsLearner(context.Background(), peerUrls)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}
