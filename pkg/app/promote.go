package app

import (
	"context"
	"log"
	"time"

	eetcd "github.com/balchua/etcd-embedded/pkg/etcd"
	"go.uber.org/zap"
)

func promote(etcdConfig *eetcd.EtcdConfig) {
	var lg *zap.Logger
	lg, err := zap.NewProduction()

	currentUrl := etcdConfig.AdvertiseClientUrls
	if !eetcd.IsLeader(currentUrl) {
		return
	}

	etcdMembers, err := eetcd.ShowMembers(etcdConfig)

	if err != nil {
		log.Fatal(err)
	}
	for _, member := range etcdMembers {
		if member.IsLearner {
			response, err := eetcd.Promote(currentUrl, member.ID)
			if err != nil {
				lg.Info("Member is not ready for promotion.", zap.Uint64("memberId", member.ID))
				continue
			}
			lg.Info("Member promoted.", zap.Uint64("memberId", response.Header.MemberId))
		}
	}
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
