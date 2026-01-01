package main

import (
	"flag"
	"fmt"
	"os"
	"rtcServer/pkg/Common/Conf"
	"rtcServer/pkg/Common/Log"
	"rtcServer/pkg/Sig/SigServ"
)

// 主函数
func main() {
	initEnvir()
	startEnvir()
}

// 初始化环境
func initEnvir() {
	// 读取命令行输入参数
	confParam := flag.String("c", "", "config file path")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  myapp -c=./config.ini")
		os.Exit(0)
	}

	// 解析命令行参数
	flag.Parse()
	if "" == *confParam {
		flag.Usage()
	}

	// 加载配置文件
	if err := Conf.InitConf(*confParam); nil != err {
		fmt.Fprintf(os.Stderr, "Init config error. err: %v\n", err)
		os.Exit(1)
	}

	// 加载日志配置
	logConf := Conf.GetLogConf()
	if err := Log.InitLog(Log.LogParam{
		LogDir:     logConf.LogDir,
		LogPrefix:  logConf.LogPrefix,
		LogLevel:   logConf.LogLevel,
		LogMaxSize: logConf.LogMaxSize,
	}); nil != err {
		fmt.Fprintf(os.Stderr, "Init log error. err: %v\n", err)
	}
}

// 启动环境
func startEnvir() {
	sigConf := Conf.GetSigConf()
	sigServ := SigServ.NewSigServ(sigConf.GetAddr(), sigConf.SigStatic)
	sigServ.Start()
}
