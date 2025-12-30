package testlog

import (
	"rtcServer/pkg/Common/Log"
	"strings"
	"sync"
	"testing"
	"time"
)

// 测试日志初始化
func TestGlobalLoggerInitialization(t *testing.T) {
	Log.InitLog(Log.LogParam{
		LogDir:     "./logInit",
		LogPrefix:  "test",
		LogLevel:   Log.LTrace,
		LogMaxSize: 5,
	})

	Log.Log().Info("Application started with global logger")
	Log.Log().Error("这是一个错误级别的日志")
	Log.Log().Debug("这是一个调试级别的日志")

	// 模拟一些操作
	time.Sleep(time.Millisecond * 100)

	Log.Log().Info("Application running with global logger")
	Log.Log().Info("All goroutines can safely use the same logger instance")
}

// 测试日志轮转功能
func TestLogRotation(t *testing.T) {
	Log.InitLog(Log.LogParam{
		LogDir:     "./logRotate",
		LogPrefix:  "test",
		LogLevel:   Log.LTrace,
		LogMaxSize: 1,
	})

	// 写入大量日志来触发轮转
	longMessage := strings.Repeat("这是一条很长的日志消息用于测试日志文件轮转功能。", 100)
	for i := range 10000 {
		Log.Log().Infof("日志消息 #%d: %s", i+1, longMessage)
	}

	Log.Log().Info("日志轮转测试完成")
}

// 测试多协程写日志
func TestLogMultiGoRoutine(t *testing.T) {
	Log.InitLog(Log.LogParam{
		LogDir:     "./logRoutine",
		LogPrefix:  "test",
		LogLevel:   Log.LTrace,
		LogMaxSize: 5,
	})

	// 协程数量
	var count = 10

	// 创建协程等待组
	var wait sync.WaitGroup
	wait.Add(count)

	for idx := range count {
		go func() {
			for inIdx := range 100 {
				Log.Log().Infof("日志消息 协程id:%d: 计数:%d 内容:%s", idx, inIdx, "this is a test.")
			}
			wait.Done()
		}()
	}

	wait.Wait()
}

// TestMultiOutput 测试日志同时输出到文件和命令行
func TestMultiOutput(t *testing.T) {
	Log.InitLog(Log.LogParam{
		LogDir:     "./logMulti",
		LogPrefix:  "test",
		LogLevel:   Log.LTrace,
		LogMaxSize: 5,
	})

	t.Log("=== 测试开始：验证日志同时输出到文件和命令行 ===")

	// 测试不同级别的日志
	Log.Log().Debug("这是一条调试级别的日志 - 应该同时显示在命令行和文件中")
	Log.Log().Info("这是一条信息级别的日志 - 应该同时显示在命令行和文件中")
	Log.Log().Warn("这是一条警告级别的日志 - 应该同时显示在命令行和文件中")
	Log.Log().Error("这是一条错误级别的日志 - 应该同时显示在命令行和文件中")

	// 测试带格式的日志
	Log.Log().Infof("当前时间: %s, 测试ID: %d", time.Now().Format("15:04:05"), 12345)

	// 模拟一些业务日志
	Log.Log().Info("用户登录成功: user_id=1001, ip=192.168.1.100")
	Log.Log().Warn("磁盘空间不足: current=85%, threshold=90%")
	Log.Log().Error("数据库连接失败: timeout=5s, retries=3")

	t.Log("=== 测试完成：请检查命令行输出和 ./log/ 目录下的日志文件 ===")
}
