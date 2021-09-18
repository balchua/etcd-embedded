package app

import (
	"context"
	"log"
	"time"

	eetcd "github.com/balchua/etcd-embedded/pkg/etcd"
	"go.uber.org/zap"
)

func promote(etcdConfig *eetcd.EtcdConfig) error {
	var lg *zap.Logger
	lg, err := zap.NewProduction()
	e, etcdErr := eetcd.NewEtcd(etcdConfig)
	if etcdErr != nil {
		lg.Error("Failed to initialize etcd client", zap.Error(etcdErr))
	}

	currentUrl := etcdConfig.AdvertiseClientUrls
	// Do not promote when the node is not the leader
	if !e.IsLeader(currentUrl) {
		return nil
	}
	etcdMembers, err := e.ShowMembers()

	if err != nil {
		log.Fatal(err)
	}
	for _, member := range etcdMembers {
		if member.IsLearner {
			response, err := e.Promote(currentUrl, member.ID)
			if err != nil {
				lg.Info("Member is not ready for promotion.", zap.Uint64("memberId", member.ID))
				continue
			}
			lg.Info("Member promoted.", zap.Uint64("memberId", response.Header.MemberId))
		}
	}
	return nil
}

func DoPromote(ctx context.Context, etcdConfig *eetcd.EtcdConfig) {

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			promote(etcdConfig)
		}
	}
}
