package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"context"
	"grpc-example/protoc/greeter"
	"log"
	"time"

	_ "google.golang.org/grpc/balancer/grpclb"
)

const (
	ADDRESS = "matosiki.localhost:50051"
)

func main() {
	cred, err := credentials.NewClientTLSFromFile("./certs/ca.crt", "")
	if err != nil {
		log.Fatalln(err)
	}

	//通过grpc 库 建立一个连接
	//conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure())
	conn, err := grpc.Dial(ADDRESS,  grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	//通过刚刚的连接 生成一个client对象。
	c := greeter.NewGreeterClient(conn)
	say, err := c.Say(context.TODO(), &greeter.Request{
		Name: "john",
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("result:", say.Message)

	//调用服务端推送流
	reqstreamData := &greeter.StreamReqData{Data: "client request stream"}
	res, _ := c.GetStream(context.Background(), reqstreamData)
	for {
		aa, err := res.Recv()
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(aa)
	}
	//客户端 推送 流
	putRes, _ := c.PutStream(context.Background())
	i := 1
	for {
		i++
		putRes.Send(&greeter.StreamReqData{Data: "client request putStream"})
		time.Sleep(time.Second)
		if i > 5 {
			break
		}
	}
	//服务端 客户端 双向流
	allStr, _ := c.AllStream(context.Background())
	go func() {
		var i = 0
		for {
			i++
			if data, err := allStr.Recv(); err == nil {
				log.Println(data)
			} else {
				allStr.CloseSend()
				break
			}
			if i == 10 {
				allStr.CloseSend()
				break
			}
		}
	}()

	go func() {
		var i = 0
		for {
			i++
			allStr.Send(&greeter.StreamReqData{Data: fmt.Sprintf("client request allStream %d", i)})
			time.Sleep(time.Second)
		}
	}()

	select {}

}
