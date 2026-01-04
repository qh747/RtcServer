package Conf

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/unknwon/goconfig"
)

// LogConfParam 日志配置参数
type LogConfParam struct {
	LogDir     string
	LogPrefix  string
	LogLevel   int
	LogMaxSize int64
}

// LogConf 日志配置参数全局变量
var LogConf LogConfParam = LogConfParam{
	LogDir:     "./log",
	LogPrefix:  "unknown",
	LogLevel:   0,
	LogMaxSize: 5,
}

// SigConfParam 信令服务配置参数
type SigConfParam struct {
	SigAddr     string
	SigPort     uint16
	SigSslPort  uint16
	SigStatic   string
	SigSslKey   string
	SigSslCert  string
	SigConnAddr string
}

// SigConf 信令配置参数全局变量
var SigConf SigConfParam = SigConfParam{
	SigAddr:     "0.0.0.0",
	SigPort:     8083,
	SigSslPort:  8443,
	SigStatic:   "./web/static",
	SigSslKey:   "./conf/cert/key.pem",
	SigSslCert:  "./conf/cert/cert.pem",
	SigConnAddr: "http://127.0.0.1:9083",
}

// InitConf      初始化配置
// @param file   配置文件
// @return error 初始化是否成功
func InitConf(file string) error {
	conf, err := goconfig.LoadConfigFile(file)
	if nil != err {
		fmt.Printf("Init config file error. err: %e, path: %s.\n", err, file)
		return err
	}

	logParam, err := conf.GetSection("LOG")
	if nil == err {
		LogConf.loadFrom(logParam)
	}

	sigParam, err := conf.GetSection("SIG")
	if nil == err {
		SigConf.loadFrom(sigParam)
	}

	return nil
}

// loadFrom      加载日志配置
// @receiver l   日志配置
// @param conf   配置参数
// @return error 加载是否成功
func (l *LogConfParam) loadFrom(conf map[string]string) error {
	var logConf LogConfParam
	dir, ok := conf["log_dir"]
	if !ok || "" == dir {
		return errors.New("Log config dir empty.")
	}
	logConf.LogDir = dir

	prefix, ok := conf["log_name"]
	if !ok || "" == prefix {
		return errors.New("Log config name empty.")
	}
	logConf.LogPrefix = prefix

	level, ok := conf["log_level"]
	if !ok || "" == level {
		return errors.New("Log config level empty.")
	}

	levelVal, err := strconv.Atoi(level)
	if nil != err {
		return errors.New("Log config level invalid.")
	}
	logConf.LogLevel = levelVal

	size, ok := conf["log_size"]
	if !ok {
		return errors.New("Log config size empty.")
	}

	sizeVal, err := strconv.Atoi(size)
	if nil != err {
		return errors.New("Log config size invalid.")
	}
	logConf.LogMaxSize = int64(sizeVal)

	LogConf = logConf
	return nil
}

// GetAddr        获取信令服务地址
// @receiver s    信令配置
// @return string 信令服务地址
func (s *SigConfParam) GetAddr() string {
	return s.SigAddr + ":" + strconv.Itoa(int(s.SigPort))
}

// GetSslAddr     获取加密信令服务地址
// @receiver s    信令配置
// @return string 加密信令服务地址
func (s *SigConfParam) GetSslAddr() string {
	return s.SigAddr + ":" + strconv.Itoa(int(s.SigSslPort))
}

// loadFrom       加载信令配置
// @receiver s	  信令配置
// @param conf	  配置参数
// @return error  加载是否成功
func (s *SigConfParam) loadFrom(conf map[string]string) error {
	var sigConf SigConfParam
	addr, ok := conf["sig_addr"]
	if !ok {
		return errors.New("Sig config addr empty.")
	}
	sigConf.SigAddr = addr

	port, ok := conf["sig_port"]
	if !ok {
		return errors.New("Sig config port empty.")
	}

	portVal, err := strconv.Atoi(port)
	if nil != err {
		return errors.New("Sig config port invalid.")
	}
	sigConf.SigPort = uint16(portVal)

	sslPort, ok := conf["sig_ssl_port"]
	if !ok {
		return errors.New("Sig config ssl port empty.")
	}

	sslPortVal, err := strconv.Atoi(sslPort)
	if nil != err {
		return errors.New("Sig config ssl port invalid.")
	}
	sigConf.SigSslPort = uint16(sslPortVal)

	static, ok := conf["sig_static"]
	if !ok {
		return errors.New("Sig config static empty.")
	}
	sigConf.SigStatic = static

	key, ok := conf["sig_ssl_key"]
	if !ok {
		return errors.New("Sig config ssl key empty.")
	}
	sigConf.SigSslKey = key

	cert, ok := conf["sig_ssl_cert"]
	if !ok {
		return errors.New("Sig config ssl cert empty.")
	}
	sigConf.SigSslCert = cert

	connAddr, ok := conf["sig_conn_addr"]
	if !ok {
		return errors.New("Sig config conn addr empty.")
	}
	sigConf.SigConnAddr = connAddr

	SigConf = sigConf
	return nil
}
