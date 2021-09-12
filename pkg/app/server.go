package app

import (
	"fmt"
	"net"

	"github.com/balchua/etcd-embedded/pkg/etcd"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var etcdConfig *etcd.EtcdConfig

func Start(cfg *etcd.EtcdConfig) {
	etcdConfig = cfg
	lg, _ := zap.NewProduction()
	port := "3000"

	if !isPortAvailable(port) {
		return
	}
	fmt.Printf("TCP Port %q is available", port)

	// Fiber instance
	app := fiber.New()

	// Routes
	app.Get("/members", members)

	// Start server
	appErr := app.Listen(":" + port)
	if appErr != nil {
		lg.Warn("Unable to start the web server, this is fine ")
	}
}

func isPortAvailable(port string) bool {
	lg, _ := zap.NewProduction()
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		lg.Info("Port is already in use.", zap.String("port", port))
		return false
	}

	ln.Close()
	return true

}
