package Log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志初始化
// param 日志初始化参数
func InitLog(param LogParam) error {
	// 创建日志句柄
	logHandle = logrus.New()
	logHandle.SetReportCaller(true)
	logHandle.SetLevel(getLevel(param.LogLevel))
	logHandle.SetFormatter(&logFormat{_fmt: "2006-01-02 15:04:05.000"})

	// 创建日志文件目录
	param.LogDir += "/" + time.Now().Format("2006-01-02")
	err := os.MkdirAll(param.LogDir, 0755)
	if err != nil {
		// 目录创建失败，回退到标准输出
		logHandle.SetOutput(os.Stdout)
		logHandle.Errorf("Create log directory error. err: %v", err)
		return err
	}

	// 创建日志文件
	lumberjackLogger := &lumberjack.Logger{
		// 使用参数中的日志目录和前缀
		Filename: param.LogDir + "/" + param.LogPrefix + ".log",
		// 日志文件最大大小(MB)
		MaxSize: int(param.LogMaxSize),
		// 保留旧文件的最大数量
		MaxBackups: 10,
		// 保留旧文件的最大天数
		MaxAge: 30,
		// 使用本地时间格式备份文件名
		LocalTime: true,
	}

	// 同时输出到文件和命令行
	multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)
	logHandle.SetOutput(multiWriter)
	return nil
}

// 获取日志句柄
// return 日志句柄
func Log() *logrus.Logger {
	return logHandle
}

// 日志句柄
var logHandle *logrus.Logger

// 日志类型
const (
	LFatal = iota
	LError
	LWarn
	LInfo
	LDebug
	LTrace
)

// 日志参数
type LogParam struct {
	LogDir     string
	LogPrefix  string
	LogLevel   int
	LogMaxSize int64
}

// 日志格式
type logFormat struct {
	_fmt string
}

// 日志内容格式化
// return 写入日志长度, 是否存在错误
// entry  日志信息
func (f *logFormat) Format(entry *logrus.Entry) ([]byte, error) {
	// 如果 entry 的 buffer 为空，创建一个新的
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 1. 时间戳：[年-月-日 时:分:秒.毫秒]
	timestamp := entry.Time
	fmt.Fprintf(b, "[%s]", timestamp.Format(f._fmt))

	// 2. 日志级别：[级别]
	level := strings.ToUpper(entry.Level.String())
	fmt.Fprintf(b, "[%s]", level)

	// 3. 文件和行号：[文件：行数]
	if entry.Caller != nil {
		filename := filepath.Base(entry.Caller.File)
		fmt.Fprintf(b, "[%s:%d]", filename, entry.Caller.Line)
	} else {
		b.WriteString("[unknown:0]")
	}

	// 4. 函数名称：[函数名称]
	if entry.HasCaller() && entry.Caller != nil {
		parts := strings.Split(entry.Caller.Function, ".")
		if len(parts) > 0 {
			funcName := parts[len(parts)-1]
			fmt.Fprintf(b, "[%s]", funcName)
		}
	} else {
		b.WriteString("[unknown]")
	}

	// 5. 日志内容
	b.WriteString(" ")

	// 如果有格式化的消息，先输出
	if entry.Message != "" {
		b.WriteString(entry.Message)
	}

	// 如果有字段，输出字段
	if len(entry.Data) > 0 {
		b.WriteString(" ")
		for k, v := range entry.Data {
			fmt.Fprintf(b, "%s=%v ", k, v)
		}
	}

	b.WriteString("\n")
	return b.Bytes(), nil
}

// 获取日志级别
// return logrus日志级别
// l      当前包指定的日志级别
func getLevel(l int) logrus.Level {
	switch l {
	case LFatal:
		return logrus.FatalLevel
	case LError:
		return logrus.ErrorLevel
	case LWarn:
		return logrus.WarnLevel
	case LInfo:
		return logrus.InfoLevel
	case LDebug:
		return logrus.DebugLevel
	case LTrace:
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}
