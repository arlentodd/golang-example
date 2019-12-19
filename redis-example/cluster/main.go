package main

import (
	"github.com/piaohao/godis"
	"log"
)

func main() {
	cluster := godis.NewRedisCluster(&godis.ClusterOption{
		Nodes:             []string{"localhost:7000", "localhost:7001", "localhost:7002", "localhost:7003", "localhost:7004", "localhost:7005"},
		ConnectionTimeout: 0,
		SoTimeout:         0,
		MaxAttempts:       0,
		Password:          "",
		PoolConfig:        &godis.PoolConfig{},
	})
	cluster.Set("cluster", "godis cluster")
	reply, _ := cluster.Get("cluster")
	log.Print(reply)
}
