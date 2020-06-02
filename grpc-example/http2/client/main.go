package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"grpc-example/protoc/http_greeter"
	"log"
)

const Addr = "matosiki.localhost:50051"

func main() {

	// TLS连接
	creds, err := credentials.NewClientTLSFromFile("./certs/ca.crt", "")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}
	conn, err := grpc.Dial(Addr, grpc.WithTransportCredentials(creds))

	if err != nil {
		grpclog.Fatalln(err)
	}

	defer conn.Close()

	// 初始化客户端
	c := http_greeter.NewGreeterHTTPAPIClient(conn)

	// 调用方法
	reqBody := new(http_greeter.HttpRequest)
	reqBody.Name = "gRPC"
	r, err := c.Say(context.Background(), reqBody)

	if err != nil {
		grpclog.Fatalln(err)
	}

	log.Println(r.Message)
}
