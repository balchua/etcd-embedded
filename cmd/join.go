package cmd

import (
	"fmt"
	"log"
	"strings"

	eetcd "github.com/balchua/etcd-embedded/pkg/etcd"
	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(joinCmd)
}

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join to an existing etcd cluster",
	Long:  `Join to an existing etcd cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatal("Invalid arguments")
		}
		config := args[1]
		leaderEndpoint := args[0]
		join(leaderEndpoint, config)
	},
}

func join(leaderEndpoint string, configFile string) {
	lg, err := zap.NewProduction()
	etcdConfig := eetcd.LoadEtcdConfig(configFile)
	e, etcdErr := eetcd.NewEtcd(etcdConfig)
	if etcdErr != nil {
		lg.Error("Failed to initialize etcd client", zap.Error(etcdErr))
	}
	memberResponse, err := e.AddMemberAsLearner(leaderEndpoint)

	if err != nil {
		log.Print(err)
	}
	lg.Info("Join Successful.",
		zap.Strings("PeerURL", memberResponse.Member.PeerURLs),
		zap.Uint64("MemberId", memberResponse.Member.ID),
		zap.Strings("ClientURL", memberResponse.Member.ClientURLs),
		zap.Bool("IsLearner", memberResponse.Member.IsLearner))

	etcdConfig.InitialClusterState = "existing"
	setInitialCluster(memberResponse, etcdConfig, configFile)

}

func setInitialCluster(memberResponse *clientv3.MemberAddResponse, etcdConfig *eetcd.EtcdConfig, configFile string) {
	conf := []string{}
	newMemberId := memberResponse.Member.ID
	for _, member := range memberResponse.Members {

		for _, u := range member.PeerURLs {
			n := member.Name
			if member.ID == newMemberId {
				n = etcdConfig.Name
			}
			conf = append(conf, fmt.Sprintf("%s=%s", n, u))
		}
	}

	etcdConfig.InitialCluster = strings.Join(conf[:], ",")
	etcdConfig.ToFile(configFile)

}
