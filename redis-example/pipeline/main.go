package main

import (
	"github.com/piaohao/godis"
	"github.com/rs/zerolog/log"
)

func main() {
	option := &godis.Option{
		Host: "localhost",
		Port: 6379,
		Db:   0,
	}
	pool := godis.NewPool(&godis.PoolConfig{}, option)
	redis, _ := pool.GetResource()
	defer redis.Close()
	p := redis.Pipelined()
	infoResp, _ := p.Info()
	timeResp, _ := p.Time()
	p.Sync()
	timeList, _ := timeResp.Get()
	log.Print(timeList)
	info, _ := infoResp.Get()
	log.Print(info)
}
