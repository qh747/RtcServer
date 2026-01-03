#!/bin/bash 

rm ../bin/signalServer
go build -o ../bin/signalServer ../cmd/sigServ/main.go

if [ ! -f "../bin/config.ini" ]; then
    cp ../conf/ini/config.ini ../bin/config.ini
fi