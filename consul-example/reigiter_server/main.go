package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"log"
	"net"
	"net/http"
)

var count int64

func main() {
	registerServer()
}

// consul 服务端会自己发送请求，来进行健康检查
func consulCheck(w http.ResponseWriter, r *http.Request) {
	s := "consulCheck" + fmt.Sprint(count) + "\tremote:" + r.RemoteAddr + "\t" + r.URL.String()
	fmt.Println(s)
	fmt.Fprintln(w, s)
	count++
}
func registerServer() {

	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "service-echo-a2f155" // 服务节点的名称
	registration.Name = "service-echo"      // 服务名称
	registration.Port = 9527                // 服务端口
	registration.Tags = []string{"v1000"}   // tag，可以为空
	registration.Address = localIP()        // 服务 IP
	registration.Address = "127.0.0.1"        // 服务 IP
	fmt.Println(registration.Address)
	serverPort := 8080
	registration.Check = &consulapi.AgentServiceCheck{ // 健康检查
		HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, serverPort, "/check"),
		Timeout:                        "3s",
		Interval:                       "5s",  // 健康检查间隔
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务，注销时间，相当于过期时间
		// GRPC:     fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service),// grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatal("register server error : ", err)
	}

	http.HandleFunc("/check", consulCheck)
	log.Printf("server start up success , url=http://localhost:%d",serverPort)
	http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil)

}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
