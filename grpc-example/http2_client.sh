#!/usr/bin/env bash

cd $(dirname $0)
pwd
go run ./http2/client/main.go
