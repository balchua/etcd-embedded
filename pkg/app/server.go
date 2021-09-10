package app

import (
	"fmt"
	"log"
	"net"

	"github.com/gofiber/fiber/v2"
	"go.etcd.io/etcd/server/v3/embed"
)

var (
	cfg                *embed.Config
	advertiseClientUrl string
)

func Start(config string) {
	port := "3000"
	cfg, err := embed.ConfigFromFile(config)
	advertiseClientUrl = cfg.ACUrls[0].String()

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
		log.Printf("Unable to start the web server, this is fine %v", err)
	}
}

func isPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Print("Port 3000 is already in use.")
		return false
	}

	ln.Close()
	return true

}
