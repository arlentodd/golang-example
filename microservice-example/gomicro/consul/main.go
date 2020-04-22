package main

import (
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"log"
	"net/http"
	"time"
)

func main() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	service := micro.NewService(micro.Registry(reg), micro.Name("greeter.client"))
	service.Init()
	service.Client()

	port := 5000
	address := ""
	server := NewServer(fmt.Sprintf("%s:%d", address, port), NewServeMux())
	log.Println("启动服务成功,端口号", server.Addr)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
	log.Println("主进程阻塞,等待退出信号.")
}

func NewServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second}
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
