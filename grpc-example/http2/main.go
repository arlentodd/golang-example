package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"grpc-example/protoc/http_greeter"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	rt "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

const Addr = "127.0.0.1:50051"

type server struct {
}

func (s *server) Say(ctx context.Context, request *http_greeter.HttpRequest) (*http_greeter.HttpResponse, error) {
	resp := new(http_greeter.HttpResponse)
	resp.Message = "Hello " + request.Name + "."
	return resp, nil
}

func main() {
	conn, err := net.Listen("tcp", Addr)
	if err != nil {
		grpclog.Fatalf("TCP Listen err:%v\n", err)
	}
	tlsConfig, err := getTLSConfig("./certs/ca.crt", "./certs/ca.key")
	if err != nil {
		grpclog.Fatalf("Failed to get TLSConfig %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))

	http_greeter.RegisterGreeterHTTPAPIServer(grpcServer, &server{})

	gwmux := rt.NewServeMux()

	//dares, err := credentials.NewClientTLSFromFile("./certs/ca.crt", "")
	//if err != nil {
	//	grpclog.Fatalf("Failed to create client TLS credentials %v", err)
	//}
	//
	//http_greeter.RegisterGreeterHTTPAPIHandlerFromEndpoint(context.Background(), gwmux, Addr, []grpc.DialOption{grpc.WithTransportCredentials(dares)})

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	srv := &http.Server{
		Addr:      Addr,
		Handler:   grpcHandlerFunc(grpcServer, mux),
		TLSConfig: tlsConfig,
	}

	grpclog.Infof("gRPC and https listen on: %s\n", Addr)

	if err := srv.Serve(tls.NewListener(conn, srv.TLSConfig)); err != nil {
		grpclog.Fatal("ListenAndServe: ", err)
	}
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	if otherHandler == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			grpcServer.ServeHTTP(w, r)
		})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}
func getTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	cert, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		grpclog.Fatalf("TLS KeyPair err: %v\n", err)
		return nil, err
	}
	certPool := x509.NewCertPool()

	cert2, err := ioutil.ReadFile("./certs/cert.crt")
	if err != nil {
		return nil, err
	}
	if !certPool.AppendCertsFromPEM(cert2) {
		panic("fail to append test ca")
	}
	return &tls.Config{
		Certificates: []tls.Certificate{pair},
		NextProtos:   []string{http2.NextProtoTLS}, // HTTP2 TLS支持
		ClientCAs:    certPool,
	}, err
}
