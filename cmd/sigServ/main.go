package main

import (
	"rtcServer/pkg/Log"
	"rtcServer/pkg/SigServ"
)

func init() {
	// 初始化日志模块
	Log.InitLog(Log.LogParam{
		LogDir:     "../log",
		LogPrefix:  "sigServ",
		LogLevel:   Log.LDebug,
		LogMaxSize: 5,
	})
}

func main() {
	sigServ := new(SigServ.SignalServer)

	sigServ.Init(":8083")
	sigServ.Start()
}
