package main

import (
	"github.com/piaohao/godis"
	"log"
)

func main() {
	option := &godis.Option{
		Host: "localhost",
		Port: 6379,
		Db:   0,
	}
	pool := godis.NewPool(nil, option)
	redis, _ := pool.GetResource()
	defer redis.Close()
	p, _ := redis.Multi()
	infoResp, _ := p.Info()
	timeResp, _ := p.Time()
	p.Exec()
	timeList, _ := timeResp.Get()
	log.Print(timeList)
	info, _ := infoResp.Get()
	log.Print(info)
}
