package etcd

import (
	"context"
	"time"

	"log"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
	"go.uber.org/zap"
)

type ETCD struct {
	Config *EtcdConfig
	Logger *zap.Logger
}

func NewEtcd(config *EtcdConfig) (*ETCD, error) {
	lg, err := zap.NewProduction()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &ETCD{
		Config: config,
		Logger: lg,
	}, nil
}

type EtcdMember struct {
	ClientURLs []string
	IsLearner  bool
	PeerURLs   []string
	ID         uint64
	Name       string
}

func (d *ETCD) StartEtcd(ctx context.Context, config *EtcdConfig) {

	e, err := embed.StartEtcd(config.ToEmbedEtcdConfig())
	if err != nil {
		d.Logger.Warn("Unable to start etcd.", zap.Error(err))
	}
	defer e.Close()
	select {
	case <-e.Server.ReadyNotify():
		d.Logger.Info("etcd Server is ready!")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		d.Logger.Warn("etcd didn't start on time.")
	}
	etcdErr := <-e.Err()

	d.Logger.Error("ETCD fatal mishap", zap.Error(etcdErr))

}

func (d *ETCD) IsLeader(advertiseClientUrl string) bool {

	cli := NewEtcdClient(advertiseClientUrl)
	defer cli.Close()
	resp, err := cli.Status(context.Background(), advertiseClientUrl)
	if err != nil {
		d.Logger.Warn("Leader check error.", zap.Error(err))
		return false
	}
	if resp.Header.MemberId == resp.Leader {
		return true
	}
	return false

}

func (d *ETCD) AddMemberAsLearner(leaderEndpoint string) (*clientv3.MemberAddResponse, error) {
	var lg *zap.Logger
	lg, err := zap.NewProduction()
	peerUrls := make([]string, 1)
	peerUrls[0] = d.Config.ListenPeerUrls

	lg.Info("AddMemberAsLearner", zap.String("peerUrl", peerUrls[0]), zap.String("LeaderEndpoint", leaderEndpoint))
	cli := NewEtcdClient(leaderEndpoint)
	defer cli.Close()

	resp, err := cli.MemberAddAsLearner(context.Background(), peerUrls)
	if err != nil {
		lg.Warn("", zap.Error(err))
		return nil, err
	}
	return resp, nil
}

func (d *ETCD) ShowMembers() ([]EtcdMember, error) {
	var lg *zap.Logger
	lg, err := zap.NewProduction()

	cli := NewEtcdClient(d.Config.AdvertiseClientUrls)
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

func (d *ETCD) Promote(leaderEndpoint string, memberId uint64) (*clientv3.MemberPromoteResponse, error) {

	cli := NewEtcdClient(leaderEndpoint)
	defer cli.Close()

	resp, err := cli.MemberPromote(context.Background(), memberId)
	if err != nil {
		d.Logger.Warn("", zap.Error(err))
		return nil, err
	}
	return resp, nil
}

func (d *ETCD) RemoveMember(leaderEndpoint string, nodeName string) error {
	peerUrls := make([]string, 1)
	peerUrls[0] = d.Config.ListenClientUrls

	cli := NewEtcdClient(leaderEndpoint)
	defer cli.Close()

	members, err := d.ShowMembers()
	if err != nil {
		d.Logger.Warn("", zap.Error(err))
		return err
	}
	for _, member := range members {
		if member.Name == nodeName {
			resp, err := cli.MemberRemove(context.Background(), member.ID)
			if err != nil {
				d.Logger.Error("Unable to remove the member", zap.String("nodeNameToRemove", nodeName))
			} else {
				d.Logger.Info("RemoveMember", zap.Uint64("memberId", resp.Header.MemberId))
				return nil
			}
		}
	}

	return nil
}
