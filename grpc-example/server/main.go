package main

import (
	"fmt"
	"grpc/pro"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

const PORT = ":50051"

type server struct {
}

//服务端 单向流
func (s *server) GetStream(req *pro.StreamReqData, res pro.Greeter_GetStreamServer) error {
	i := 0
	for {
		i++
		res.Send(&pro.StreamResData{Data: fmt.Sprintf("%v", time.Now().Unix())})
		time.Sleep(1 * time.Second)
		if i > 10 {
			break
		}
	}
	return nil
}

//客户端 单向流
func (s *server) PutStream(cliStr pro.Greeter_PutStreamServer) error {

	for {
		if tem, err := cliStr.Recv(); err == nil {
			log.Println(tem)
		} else {
			log.Println("break, err :", err)
			break
		}
	}

	return nil
}

//客户端服务端 双向流
func (s *server) AllStream(allStr pro.Greeter_AllStreamServer) error {

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for {
			if data, err := allStr.Recv(); err == nil {
				log.Println(data)
			} else {
				log.Println("break, err :", err)
				break
			}
		}
		wg.Done()
	}()

	go func() {
		for {
			allStr.Send(&pro.StreamResData{Data: "ssss"})
			time.Sleep(time.Second)
		}
		wg.Done()
	}()

	wg.Wait()
	return nil
}

func main() {
	//监听端口
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		return
	}
	//创建一个grpc 服务器
	s := grpc.NewServer()
	//注册事件
	pro.RegisterGreeterServer(s, &server{})
	//处理链接
	s.Serve(lis)
}
