#!/bin/bash 

rm ../bin/signalServer
go build -o ../bin/signalServer ../cmd/sigServ/main.go