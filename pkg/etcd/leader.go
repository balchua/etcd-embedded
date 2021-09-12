package etcd

import (
	"context"

	"go.uber.org/zap"
)

func IsLeader(advertiseClientUrl string) bool {
	var lg *zap.Logger
	lg, err := zap.NewProduction()

	var endpoints = make([]string, 1)
	endpoints[0] = advertiseClientUrl
	cli := NewEtcdClient(endpoints)
	defer cli.Close()
	resp, err := cli.Status(context.Background(), advertiseClientUrl)
	if err != nil {
		lg.Warn("Leader check error.", zap.Error(err))
		return false
	}
	if resp.Header.MemberId == resp.Leader {
		return true
	}
	return false

}
