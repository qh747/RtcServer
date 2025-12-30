package main

import (
	"rtcServer/pkg/Log"
	SigServ "rtcServer/pkg/Sig/SigServ"
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
	sigServ := SigServ.New(":8083")
	sigServ.Start()
}
