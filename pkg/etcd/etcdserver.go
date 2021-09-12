package etcd

import (
	"context"
	"log"
	"time"

	"go.etcd.io/etcd/server/v3/embed"
	"go.uber.org/zap"
)

func StartEtcd(ctx context.Context, config *EtcdConfig) {
	var lg *zap.Logger
	lg, err := zap.NewProduction()

	e, err := embed.StartEtcd(config.ToEmbedEtcdConfig())
	if err != nil {
		lg.Warn("Unable to start etcd.", zap.Error(err))
	}
	defer e.Close()
	select {
	case <-e.Server.ReadyNotify():
		lg.Info("etcd Server is ready!")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		lg.Warn("etcd didn't start on time.")
	}
	log.Fatal(<-e.Err())

}
