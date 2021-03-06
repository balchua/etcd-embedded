package app

import (
	"log"

	eetcd "github.com/balchua/etcd-embedded/pkg/etcd"
	"github.com/gofiber/fiber/v2"
)

type ClusterMember struct {
	ClientURLs []string `json:"client_urls"`
	IsLearner  bool     `json:"learner"`
	PeerURLs   []string `json:"peer_urls"`
	Id         uint64   `json:"id"`
	Name       string   `json:"name"`
}

func members(c *fiber.Ctx) error {
	var response []ClusterMember
	e, etcdErr := eetcd.NewEtcd(etcdConfig)

	if etcdErr != nil {
		return etcdErr
	}
	etcdMembers, err := e.ShowMembers()

	if err != nil {
		log.Fatal(err)
	}

	for _, member := range etcdMembers {
		membr := ClusterMember{}
		membr.ClientURLs = member.ClientURLs
		membr.PeerURLs = member.PeerURLs
		membr.Id = member.ID
		membr.Name = member.Name
		response = append(response, membr)
	}

	return c.JSON(response)
}
