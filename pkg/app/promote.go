package app

import (
	"context"
	"log"
	"time"

	eetcd "github.com/balchua/etcd-embedded/pkg/etcd"
	"go.etcd.io/etcd/server/v3/embed"
)

func promote(etcdConfig string) {
	cfg, err := embed.ConfigFromFile(etcdConfig)

	currentUrl := cfg.ACUrls[0].String()
	if !eetcd.IsLeader(currentUrl) {
		return
	}

	etcdMembers, err := eetcd.ShowMembers(currentUrl)

	if err != nil {
		log.Fatal(err)
	}
	for _, member := range etcdMembers {
		if member.IsLearner {
			resp, err := eetcd.Promote(currentUrl, member.ID)
			if err != nil {
				log.Printf("Member ID: %d is not ready for promotion", member.ID)
			} else {
				log.Printf("Promoted member %d, ", resp.Header.MemberId)
			}

		}
	}

	return
}

func DoPromote(ctx context.Context, etcdConfig string) {

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			promote(etcdConfig)
		}
	}
}
