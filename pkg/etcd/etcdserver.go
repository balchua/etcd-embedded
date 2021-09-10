package etcd

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"go.etcd.io/etcd/server/v3/embed"
)

func StartEtcd(ctx context.Context, config string) {
	cfg, err := embed.ConfigFromFile(config)

	if err != nil {
		log.Fatalf("%v", err)
	}

	existingCluster, err := isInitialized(ctx, cfg.Dir)

	if !existingCluster {
		// err := initDir(ctx, cfg.Dir)
		// if err != nil {
		// 	log.Print(err)
		// }
	}

	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Print(err)
	}
	defer e.Close()
	select {
	case <-e.Server.ReadyNotify():
		log.Printf("Server is ready!")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		log.Printf("Server took too long to start!")
	}
	log.Fatal(<-e.Err())

}

func isInitialized(ctx context.Context, etcdDir string) (bool, error) {
	dir := filepath.Join(etcdDir, "member", "wal")
	if s, err := os.Stat(dir); err == nil && s.IsDir() {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, errors.Wrapf(err, "invalid state for wal directory %s", dir)
	}
}
