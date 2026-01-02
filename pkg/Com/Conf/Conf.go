package Conf

import (
	"fmt"
	"strconv"

	"github.com/unknwon/goconfig"
)

/** -------------------------------------------- EXT --------------------------------------------- */

// 日志配置参数
type LogConfParam struct {
	LogDir     string
	LogPrefix  string
	LogLevel   int
	LogMaxSize int64
}

// 信令服务配置参数
type SigConfParam struct {
	SigAddr     string
	SigPort     uint16
	SigSslPort  uint16
	SigStatic   string
	SigSslKey   string
	SigSslCert  string
	SigConnAddr string
}

func (s *SigConfParam) GetAddr() string {
	return s.SigAddr + ":" + strconv.Itoa(int(s.SigPort))
}

func (s *SigConfParam) GetSslAddr() string {
	return s.SigAddr + ":" + strconv.Itoa(int(s.SigSslPort))
}

func InitConf(file string) error {
	conf, err := goconfig.LoadConfigFile(file)
	if nil != err {
		fmt.Printf("Init config file error. err: %e, path: %s.\n", err, file)
		return err
	}

	loadLogConf(conf.GetSection("LOG"))
	loadSigConf(conf.GetSection("SIG"))
	return nil
}

func GetLogConf() LogConfParam {
	return logConf
}

func GetSigConf() SigConfParam {
	return sigConf
}

/** -------------------------------------------- IN --------------------------------------------- */

var logConf LogConfParam = LogConfParam{
	LogDir:     "./log",
	LogPrefix:  "unknown",
	LogLevel:   0,
	LogMaxSize: 5,
}

var sigConf SigConfParam = SigConfParam{
	SigAddr:     "0.0.0.0",
	SigPort:     8083,
	SigSslPort:  8443,
	SigStatic:   "./web/static",
	SigSslKey:   "./conf/cert/key.pem",
	SigSslCert:  "./conf/cert/cert.pem",
	SigConnAddr: "127.0.0.1:0983",
}

func loadLogConf(conf map[string]string, err error) {
	// 加载配置错误，使用默认配置
	if nil != err {
		fmt.Println("Load log config error. not exists. use default config.")
		return
	}

	dir, ok := conf["log_dir"]
	if !ok {
		fmt.Println("Log config dir empty.")
	} else {
		logConf.LogDir = dir
	}

	prefix, ok := conf["log_name"]
	if !ok {
		fmt.Println("Log config name empty.")
	} else {
		logConf.LogPrefix = prefix
	}

	level, ok := conf["log_level"]
	if !ok {
		fmt.Println("Log config level empty.")
	} else {
		val, err := strconv.Atoi(level)
		if nil != err {
			fmt.Println("Log config level invalid.")
		} else {
			logConf.LogLevel = val
		}
	}

	size, ok := conf["log_size"]
	if !ok {
		fmt.Println("Log config size empty.")
	} else {
		val, err := strconv.Atoi(size)
		if nil != err {
			fmt.Println("Log config size invalid.")
		} else {
			logConf.LogMaxSize = int64(val)
		}
	}
}

func loadSigConf(conf map[string]string, err error) {
	// 加载配置错误，使用默认配置
	if nil != err {
		fmt.Println("Load sig config error. not exists. use default config.")
		return
	}

	addr, ok := conf["sig_addr"]
	if !ok {
		fmt.Println("Sig config addr empty.")
	} else {
		sigConf.SigAddr = addr
	}

	port, ok := conf["sig_port"]
	if !ok {
		fmt.Println("Sig config port empty.")
	} else {
		val, err := strconv.Atoi(port)
		if nil != err {
			fmt.Println("Sig config port invalid.")
		} else {
			sigConf.SigPort = uint16(val)
		}
	}

	sslPort, ok := conf["sig_ssl_port"]
	if !ok {
		fmt.Println("Sig config ssl port empty.")
	} else {
		val, err := strconv.Atoi(sslPort)
		if nil != err {
			fmt.Println("Sig config ssl port invalid.")
		} else {
			sigConf.SigSslPort = uint16(val)
		}
	}

	static, ok := conf["sig_static"]
	if !ok {
		fmt.Println("Sig config static empty.")
	} else {
		sigConf.SigStatic = static
	}

	key, ok := conf["sig_ssl_key"]
	if !ok {
		fmt.Println("Sig config ssl key empty.")
	} else {
		sigConf.SigSslKey = key
	}

	cert, ok := conf["sig_ssl_cert"]
	if !ok {
		fmt.Println("Sig config ssl cert empty.")
	} else {
		sigConf.SigSslCert = cert
	}

	connAddr, ok := conf["sig_conn_addr"]
	if !ok {
		fmt.Println("Sig config conn addr empty.")
	} else {
		sigConf.SigConnAddr = connAddr
	}
}
