package main

import (
	"github.com/piaohao/godis"
	"time"
)

func main() {
	locker := godis.NewLocker(&godis.Option{
		Host: "localhost",
		Port: 6379,
		Db:   0,
	}, &godis.LockOption{
		Timeout: 5 * time.Second,
	})
	lock, err := locker.TryLock("lock")
	if err == nil && lock != nil {
		locker.UnLock(lock)
	}

}
