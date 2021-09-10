package etcd

import (
	"context"
	"log"
)

func IsLeader(advertiseClientUrl string) bool {
	var endpoints = make([]string, 1)
	endpoints[0] = advertiseClientUrl
	cli := NewEtcdClient(endpoints)
	defer cli.Close()
	resp, err := cli.Status(context.Background(), advertiseClientUrl)
	if err != nil {
		log.Print(err)
		return false
	}
	if resp.Header.MemberId == resp.Leader {
		return true
	}
	return false

}
