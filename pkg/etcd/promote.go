package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func Promote(leaderEndpoint string, memberId uint64) (*clientv3.MemberPromoteResponse, error) {
	var lg *zap.Logger
	lg, err := zap.NewProduction()
	var endpoints = make([]string, 1)
	endpoints[0] = leaderEndpoint
	cli := NewEtcdClient(endpoints)
	defer cli.Close()

	resp, err := cli.MemberPromote(context.Background(), memberId)
	if err != nil {
		lg.Warn("", zap.Error(err))
		return nil, err
	}
	return resp, nil
}
