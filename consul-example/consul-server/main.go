package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/rs/zerolog/log"
	"net"
	"strconv"
)

func main() {
	var lastIndex uint64
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500" //consul server

	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("api new client is failed, err:", err)
		return
	}
	services, metainfo, err := client.Health().Service("serverNode", "v1000", true, &consulapi.QueryOptions{
		WaitIndex: lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
	})
	if err != nil {
		log.Printf("error retrieving instances from Consul: %v", err)
	}
	lastIndex = metainfo.LastIndex
	addrs := map[string]struct{}{}
	for _, service := range services {
		fmt.Println("service.Service.Address:", service.Service.Address, "service.Service.Port:", service.Service.Port)
		addrs[net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))] = struct{}{}
	}
}
