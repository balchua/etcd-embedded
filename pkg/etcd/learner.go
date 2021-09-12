package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func AddMemberAsLearner(leaderEndpoint string, cfg *EtcdConfig) (*clientv3.MemberAddResponse, error) {
	var lg *zap.Logger
	lg, err := zap.NewProduction()
	peerUrls := make([]string, 1)
	peerUrls[0] = cfg.ListenPeerUrls

	lg.Info("AddMemberAsLearner", zap.String("peerUrl", peerUrls[0]), zap.String("LeaderEndpoint", leaderEndpoint))

	var endpoints = make([]string, 1)
	endpoints[0] = leaderEndpoint
	cli := NewEtcdClient(endpoints)
	defer cli.Close()

	resp, err := cli.MemberAddAsLearner(context.Background(), peerUrls)
	if err != nil {
		lg.Warn("", zap.Error(err))
		return nil, err
	}
	return resp, nil
}
