package app

import (
	"fmt"
	"log"
	"strings"

	betcd "github.com/balchua/etcd-embedded/pkg/etcd"
	"github.com/gofiber/fiber/v2"
)

func members(c *fiber.Ctx) error {
	var b strings.Builder

	etcdMembers, err := betcd.ShowMembers(advertiseClientUrl)

	if err != nil {
		log.Fatal(err)
	}
	b.Grow(32)
	for _, member := range etcdMembers {
		fmt.Fprintf(&b, "%s||%t||%s||%d, ", member.ClientURLs, member.IsLearner, member.PeerURLs, member.ID)
	}

	return c.SendString(b.String())
}
