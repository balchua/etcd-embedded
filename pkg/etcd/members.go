package etcd

import (
	"context"
	"log"
)

type EtcdMember struct {
	ClientURLs []string
	IsLearner  bool
	PeerURLs   []string
	ID         uint64
}

func ShowMembers(endpoint string) ([]EtcdMember, error) {
	var endpoints []string
	endpoints = make([]string, 1)
	endpoints[0] = endpoint
	cli := NewEtcdClient(endpoints)
	defer cli.Close()
	resp, err := cli.MemberList(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var members []EtcdMember
	for _, member := range resp.Members {
		etcdMember := EtcdMember{
			ClientURLs: member.ClientURLs,
			IsLearner:  member.IsLearner,
			PeerURLs:   member.PeerURLs,
			ID:         member.ID,
		}
		members = append(members, etcdMember)

	}

	return members, nil
}
