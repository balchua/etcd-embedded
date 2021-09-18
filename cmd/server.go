package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/balchua/etcd-embedded/pkg/app"
	eetcd "github.com/balchua/etcd-embedded/pkg/etcd"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts an embedded etcd server as the seed etcd",
	Long:  `Starts an embedded etcd server as the seed etcd`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Invalid arguments")
		}
		etcdConfig := eetcd.LoadEtcdConfig(args[0])
		etcdConfig.ToFile(args[0])
		e, _ := eetcd.NewEtcd(etcdConfig)

		go e.StartEtcd(context.Background(), etcdConfig)
		go app.DoPromote(context.Background(), etcdConfig)
		go app.Start(etcdConfig)

		ch := make(chan os.Signal)
		signal.Notify(ch, unix.SIGPWR)
		signal.Notify(ch, unix.SIGINT)
		signal.Notify(ch, unix.SIGQUIT)
		signal.Notify(ch, unix.SIGTERM)
		<-ch
	},
}
