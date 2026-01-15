#!/bin/bash 

rm ../bin/sigServ/signalServer
go build -o ../bin/sigServ/signalServer ../cmd/sigServ/main.go

if [ ! -f "../bin/sigServ/sigConfig.ini" ]; then
    cp ../conf/ini/sigConfig.ini ../bin/sigServ/sigConfig.ini
fi