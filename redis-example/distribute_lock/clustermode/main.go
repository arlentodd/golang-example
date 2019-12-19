package main

import (
	"github.com/piaohao/godis"
	"time"
)

func main() {
	locker := godis.NewClusterLocker(&godis.ClusterOption{
		Nodes:             []string{"localhost:7000", "localhost:7001", "localhost:7002", "localhost:7003", "localhost:7004", "localhost:7005"},
		ConnectionTimeout: 0,
		SoTimeout:         0,
		MaxAttempts:       0,
		Password:          "",
		PoolConfig:        &godis.PoolConfig{},
	}, &godis.LockOption{
		Timeout: 5 * time.Second,
	})
	lock, err := locker.TryLock("lock")
	if err == nil && lock != nil {
		//do something
		locker.UnLock(lock)
	}
}
