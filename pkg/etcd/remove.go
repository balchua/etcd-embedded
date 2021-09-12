package etcd

import (
	"context"

	"go.uber.org/zap"
)

func RemoveMember(leaderEndpoint string, nodeName string, etcdConfig *EtcdConfig) error {
	var lg *zap.Logger
	lg, err := zap.NewProduction()
	peerUrls := make([]string, 1)
	peerUrls[0] = etcdConfig.ListenClientUrls

	var endpoints = make([]string, 1)
	endpoints[0] = leaderEndpoint
	cli := NewEtcdClient(endpoints)
	defer cli.Close()

	members, err := ShowMembers(etcdConfig)
	if err != nil {
		lg.Warn("", zap.Error(err))
		return err
	}
	for _, member := range members {
		if member.Name == nodeName {
			resp, err := cli.MemberRemove(context.Background(), member.ID)
			if err != nil {
				lg.Error("Unable to remove the member", zap.String("nodeNameToRemove", nodeName))
			} else {
				lg.Info("RemoveMember", zap.Uint64("memberId", resp.Header.MemberId))
				return nil
			}
		}
	}

	return nil
}
