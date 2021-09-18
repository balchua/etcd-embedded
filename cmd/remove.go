package cmd

import (
	"log"

	eetcd "github.com/balchua/etcd-embedded/pkg/etcd"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a node from the etcd cluster",
	Long:  `Remove a node from the etcd cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatal("Invalid arguments")
		}
		name := args[1]
		config := args[2]
		leaderEndpoint := args[0]
		remove(leaderEndpoint, name, config)
	},
}

func remove(leaderEndpoint string, nodeName string, configFile string) {
	lg, _ := zap.NewProduction()
	etcdConfig := eetcd.LoadEtcdConfig(configFile)
	e, etcdErr := eetcd.NewEtcd(etcdConfig)
	if etcdErr != nil {
		lg.Error("Failed to initialize etcd client", zap.Error(etcdErr))
	}
	remove_err := e.RemoveMember(leaderEndpoint, nodeName)

	if remove_err != nil {
		log.Print(remove_err)
	}
	lg.Info("Remove Successful.")
}
