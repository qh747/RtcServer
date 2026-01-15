#!/bin/bash 

rm ../bin/medServ/mediaServer
go build -o ../bin/medServ/mediaServer ../cmd/medServ/main.go

if [ ! -f "../bin/medServ/medConfig.ini" ]; then
    cp ../conf/ini/medConfig.ini ../bin/medServ/medConfig.ini
fi