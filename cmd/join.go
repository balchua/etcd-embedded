package cmd

import (
	"log"

	eetcd "github.com/balchua/etcd-embedded/pkg/etcd"
	"github.com/spf13/cobra"
	"go.etcd.io/etcd/server/v3/embed"
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

func join(leaderEndpoint string, config string) {
	cfg, err := embed.ConfigFromFile(config)

	memberResponse, err := eetcd.AddMemberAsLearner(leaderEndpoint, cfg.LPUrls[0].String())
	if err != nil {
		log.Print(err)
	}
	log.Printf("PeerURL: %s, MemberId %d, ClientURL: %s, IsLearner: %t", memberResponse.Member.GetPeerURLs(), memberResponse.Member.ID, memberResponse.Member.ClientURLs, memberResponse.Member.IsLearner)

}
