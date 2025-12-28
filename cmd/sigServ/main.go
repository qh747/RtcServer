package main

import (
	"rtcServer/pkg/SigServ"
)

func main() {
	sigServ := new(SigServ.SignalServer)

	sigServ.Init(":8083")
	sigServ.Start()
}
