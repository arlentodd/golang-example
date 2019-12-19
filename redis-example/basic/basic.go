package main

import (
	"github.com/piaohao/godis"
	"github.com/rs/zerolog/log"
)

func main() {
	redis := godis.NewRedis(&godis.Option{
		Host: "localhost",
		Port: 6379,
		Db:   0,
	})
	defer redis.Close()
	redis.Set("godis", "1")
	arr, _ := redis.Get("godis")
	log.Print(arr)
}
