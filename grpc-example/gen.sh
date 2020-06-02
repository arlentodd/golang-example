#!/usr/bin/env bash

cd $(dirname $0)

pwd

go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go

protoc --proto_path=. --go_out=plugins=grpc:. protoc/greeter/greeter.proto

protoc --proto_path=. -I/usr/local/include -I. -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.6 \
  --grpc-gateway_out=logtostderr=true:. \
  --go_out=plugins=grpc:. protoc/http_greeter/http_greeter.proto

echo "=========================================================================="
mkdir -p certs

openssl genrsa -out ./certs/ca.key

# 创建认证请求信息
openssl req -new -key ./certs/ca.key -out ./certs/ca.csr -subj "/C=CN/ST=HB/L=WH/O=matosiki/OU=dev/CN=matosiki.localhost/emailAddress=wx11055@163.com"

# 使用x509自签名
openssl req -new -x509 -days 365 -key ./certs/ca.key -out ./certs/cert.crt -subj "/C=CN/ST=HB/L=WH/O=matosiki/OU=dev/CN=matosiki.localhost/emailAddress=wx11055@163.com"

# 使用 CA 证书及CA密钥 对请求签发证书进行签发，生成x509证书
openssl x509 -req -days 3650 -in ./certs/ca.csr -CA ./certs/cert.crt -CAkey ./certs/ca.key -CAcreateserial -out ./certs/ca.crt

# 查看认证请求文件信息,包含了基本信息，并没有签名。
# openssl req -in ./certs/ca.csr  -noout -text

# 查看证书
#openssl x509 -in ./certs/ca.crt -noout -text

# 使用x509查看证书
#openssl x509 -in ./certs/ca.crt -noout -text

# 直接生成证书，跳过生成认证请求文件，证书自己颁发自己生成签名。
# openssl req -new -x509 -days -nodes  365 -key ./certs/ca.key -out ./certs/ca.crt -subj "/C=CN/ST=HB/L=WH/O=matosiki/OU=dev/CN=matosiki.localhost/emailAddress=wx11055@163.com"
