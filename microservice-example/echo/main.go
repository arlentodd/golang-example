package main

//优雅的关闭服务

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m := map[string]interface{}{}
		m["method"] = r.Method
		m["path"] = r.RequestURI
		m["header"] = r.Header
		data := map[string]interface{}{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			panic(err)
		}
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
	server := NewServer(":5000", mux)
	done := make(chan struct{}, 1)
	quitSign := make(chan os.Signal, 1)
	signal.Notify(quitSign, syscall.SIGINT, syscall.SIGTERM)
	go func() {
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
		close(done)
	}()
	log.Println("启动服务成功,端口号", server.Addr)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
	log.Println("主进程阻塞,等待退出信号.")
	<-done
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
