package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"rtcServer/pkg/Com/Conf"
	"rtcServer/pkg/Com/Log"
	"rtcServer/pkg/Sig/SigConn"
	"rtcServer/pkg/Sig/SigServ"
	"sync"
	"syscall"
	"time"
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
	if err := Log.InitLog(Log.LogParam{
		LogDir:     Conf.LogConf.LogDir,
		LogPrefix:  Conf.LogConf.LogPrefix,
		LogLevel:   Conf.LogConf.LogLevel,
		LogMaxSize: Conf.LogConf.LogMaxSize,
	}); nil != err {
		fmt.Fprintf(os.Stderr, "Init log error. err: %v\n", err)
		os.Exit(1)
	}

	// 加载媒体服务连接配置
	if err := SigConn.InitSigConnSelect(Conf.SigConf.SigConnAddr); nil != err {
		fmt.Fprintf(os.Stderr, "Init media conn error. err: %v\n", err)
		os.Exit(1)
	}
}

// 启动环境
func startEnvir() {
	// 创建context用于优雅关闭
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigServ := SigServ.NewSigServ(Conf.SigConf.GetAddr(), Conf.SigConf.SigStatic)
	sigSslServ := SigServ.NewSigSslServ(Conf.SigConf.GetSslAddr(), Conf.SigConf.SigStatic, Conf.SigConf.SigSslKey, Conf.SigConf.SigSslCert)

	// 创建等待组用于等待所有服务关闭
	var wg sync.WaitGroup

	// 子协程启动服务
	wg.Go(func() {
		sigServ.Start()
	})

	wg.Go(func() {
		sigSslServ.Start()
	})

	// 等待退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sigChan
	Log.Log().Info("Received shutdown signal, shutting down servers...")

	cancel()

	// 关闭服务
	sigServ.Stop()
	sigSslServ.Stop()

	// 等待所有服务协程结束
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	// 等待服务关闭或超时
	select {
	case <-done:
		Log.Log().Info("All servers stopped gracefully")
	case <-time.After(5 * time.Second):
		Log.Log().Warn("Timeout waiting for servers to stop")
	}
}
