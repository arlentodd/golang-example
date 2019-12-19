package main

import (
	"fmt"
	"github.com/piaohao/godis"
	"log"
	"time"
)

func main() {
	option := &godis.Option{
		Host: "localhost",
		Port: 6379,
		Db:   0,
	}
	pool := godis.NewPool(&godis.PoolConfig{}, option)
	var sub = func(channel string) {
		redis, _ := pool.GetResource()
		defer redis.Close()
		pubsub := &godis.RedisPubSub{
			OnMessage: func(channel, message string) {
				// 接受消息函数
				log.Printf("OnMessage: channel=%s message=%s", channel, message)
			},
			OnSubscribe: func(channel string, subscribedChannels int) {
				log.Printf("OnSubscribe: channel=%s subscribedChannels=%d", channel, subscribedChannels)
			},
			OnPong: func(channel string) {
				log.Print("receive pong")
			},
		}

		redis.Subscribe(pubsub, channel) //第一个监听
	}
	go sub("godis")
	go sub("godis")

	forever := make(chan int, 1)
	{
		for i := 0; i < 20; i++ {
			time.Sleep(time.Millisecond * 50)
			redis, _ := pool.GetResource()
			redis.Publish("godis", fmt.Sprintf("godis pubsub %d", i))
			redis.Close()
		}
	}
	{
		redis, _ := pool.GetResource()
		defer redis.Close()
		value, err := redis.Get("godis")
		if err != nil {
			panic(err)
		}
		log.Printf("Get [godis] value : %s", value)
	}
	log.Printf(" [*] To exit press CTRL+C")
	<-forever
}
