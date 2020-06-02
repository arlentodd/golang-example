package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials"

	"grpc-example/protoc/greeter"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

const Addr = "127.0.0.1:50051"

type server struct {
}

//
//func (s *server) Say(req *greeter.Request, rps greeter.Response) error {
//}

func (s *server) Say(ctx context.Context, req *greeter.Request) (*greeter.Response, error) {
	return &greeter.Response{Message: "Hello " + req.Name + "!"}, nil
}

//服务端 单向流
func (s *server) GetStream(req *greeter.StreamReqData, res greeter.Greeter_GetStreamServer) error {
	i := 0
	for {
		i++
		res.Send(&greeter.StreamResData{Data: fmt.Sprintf("server return getStream:%d:%v", i, time.Now().Unix())})
		time.Sleep(1 * time.Second)
		if i >= 5 {
			break
		}
	}
	return nil
}

//客户端 单向流
func (s *server) PutStream(cliStr greeter.Greeter_PutStreamServer) error {

	for {
		if tem, err := cliStr.Recv(); err == nil {
			log.Println("server get putStream", tem)
		} else {
			log.Println("break, err :", err)
			break
		}
	}

	return nil
}

//客户端服务端 双向流
func (s *server) AllStream(allStr greeter.Greeter_AllStreamServer) error {

	wg := sync.WaitGroup{}
	wg.Add(2)
	var i = 0
	go func() {
		for {
			if data, err := allStr.Recv(); err == nil && i <= 10 {
				log.Println("server get allStream", data)
				i++
			} else {
				if err != nil {
					log.Println("break, err :", err)
				}
				break
			}
		}
		wg.Done()
	}()

	go func() {
		for {
			err := allStr.Send(&greeter.StreamResData{Data: "server return allStream"})
			if err != nil {
				break
			}
			time.Sleep(time.Second)
		}
		wg.Done()
	}()

	wg.Wait()
	return nil
}

func main() {

	cred, err := credentials.NewServerTLSFromFile("./certs/ca.crt", "./certs/ca.key")
	if err != nil {
		log.Fatalln(err)
	}

	//监听端口
	lis, err := net.Listen("tcp",Addr)
	if err != nil {
		return
	}
	//创建一个grpc 服务器
	s := grpc.NewServer(grpc.Creds(cred))
	//注册事件
	greeter.RegisterGreeterServer(s, &server{})
	//处理链接
	s.Serve(lis)
}
