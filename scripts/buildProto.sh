#!/bin/bash 

cd ../proto

if ls *.go >/dev/null 2>&1; then
    rm *.go
fi

protoc --go_out=. --go-grpc_out=. *.proto