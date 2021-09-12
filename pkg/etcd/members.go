package etcd

import (
	"context"

	"go.uber.org/zap"
)

type EtcdMember struct {
	ClientURLs []string
	IsLearner  bool
	PeerURLs   []string
	ID         uint64
	Name       string
}

func ShowMembers(etcdConfig *EtcdConfig) ([]EtcdMember, error) {
	var lg *zap.Logger
	lg, err := zap.NewProduction()

	var endpoints []string
	endpoints = make([]string, 1)
	endpoints[0] = etcdConfig.AdvertiseClientUrls

	cli := NewEtcdClient(endpoints)
	defer cli.Close()
	resp, err := cli.MemberList(context.Background())

	if err != nil {
		lg.Warn("", zap.Error(err))
		return nil, err
	}

	var members []EtcdMember
	for _, member := range resp.Members {
		etcdMember := EtcdMember{
			ClientURLs: member.ClientURLs,
			IsLearner:  member.IsLearner,
			PeerURLs:   member.PeerURLs,
			ID:         member.ID,
			Name:       member.Name,
		}
		members = append(members, etcdMember)
	}

	return members, nil
}
