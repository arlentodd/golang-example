package main

//优雅的关闭服务

import (
	"context"
	"encoding/json"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	port := 5000
	address := ""
	server := NewServer(fmt.Sprintf("%s:%d", address, port), NewServeMux())
	mainContext, mainFunc := context.WithCancel(context.Background())
	exitContext, exitFunc := context.WithCancel(mainContext)
	defer mainFunc()
	go func() {
		log.Println("启动服务成功,端口号", server.Addr)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
		log.Println("主进程阻塞,等待退出信号.")
	}()
	go func() {
		client, err := NewConsulClient()
		if err != nil {
			panic(err)
		}
		id := "fitmgr-bp-abcd123"
		err = client.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
			ID:      id,
			Name:    "fitmgr-bp",
			Port:    port,
			Tags:    []string{"master"},
			Address: address,
			Check: &consulapi.AgentServiceCheck{
				HTTP:                           fmt.Sprintf("http://%s:%d%s", address, port, "/health"),
				Timeout:                        "3s",
				Interval:                       "5s",  // 健康检查间隔
				DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务，注销时间，相当于过期时间
			},
		})
		defer client.Agent().ServiceDeregister(id)
		if err != nil {
			panic(err)
		}

		<-exitContext.Done()
	}()
	go func() {
		quitSign := make(chan os.Signal, 1)
		signal.Notify(quitSign, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quitSign
		log.Println("退出信号类型:", sig, ",优雅关闭中,关闭服务的事情....")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.SetKeepAlivesEnabled(false)
		err := server.Shutdown(ctx)
		if err != nil {
			log.Fatalf("不能够优雅的关闭服务: %v \n", err)
		}
		time.Sleep(time.Second) //睡眠一秒,表示处理相关关闭服务的事情
		log.Println("关闭服务相关事情处理结束")
		exitFunc()
		mainFunc()
	}()
	<-mainContext.Done()
	log.Println("收到退出信号,主进程结束,退出.")
}

func NewServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second}
}

var ConsulClient *consulapi.Client

func init() {
	var err error
	ConsulClient, err = NewConsulClient()
	if err != nil {
		panic(err)
	}
}

func NewConsulClient() (*consulapi.Client, error) {
	return consulapi.NewClient(consulapi.DefaultConfig())
}

var count int64

func NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		s := "consulCheck" + fmt.Sprint(count) + "\tremote:" + r.RemoteAddr + "\t" + r.URL.String()
		fmt.Println(s)
		fmt.Fprintln(w, s)
		count++
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m := map[string]interface{}{}
		m["method"] = r.Method
		m["path"] = r.RequestURI
		m["header"] = r.Header
		data := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&data)
		for i := range r.URL.Query() {
			data[i] = r.URL.Query() [i]
		}
		for i := range r.Form {
			data[i] = r.Form[i]
		}
		for i := range r.PostForm {
			data[i] = r.PostForm[i]
		}

		v, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(v)
		if err != nil {
			panic(err)
		}

	})

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	mux.HandleFunc("/tohello", func(w http.ResponseWriter, r *http.Request) {

	})
	return mux
}
