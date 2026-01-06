package SigConn

import (
	"errors"
	"strconv"
	"strings"
)

// const 连接类型
const (
	Rpc   = "rpc"
	Http  = "http"
	Https = "https"
)

type (
	// Type 连接类型
	Type string

	// Conn 媒体服务连接接口
	Conn interface {
		// 发起请求
		Req() (string, error)
		// 发起异步请求
		ReqAsync(func(string, error))
	}

	// Addr 连接地址
	Addr struct {
		_type string
		_addr string
		_port uint16
	}
)

// LoadFrom 加载连接地址
// @receiver c 连接地址
// @param addr 连接地址
// @return error 加载是否存在错误
func (c *Addr) LoadFrom(addr string) error {
	if "" == addr {
		return errors.New("Param empty")
	}

	// 解析连接类型
	var ipPort string
	if strings.HasPrefix(addr, "rpc://") {
		c._type = Rpc
		ipPort = strings.TrimPrefix(addr, "rpc://")
	} else if strings.HasPrefix(addr, "http://") {
		c._type = Http
		ipPort = strings.TrimPrefix(addr, "http://")
	} else if strings.HasPrefix(addr, "https://") {
		c._type = Https
		ipPort = strings.TrimPrefix(addr, "https://")
	} else {
		return errors.New("Param type invalid: " + addr)
	}

	// 解析ip和端口，格式为 ip:port
	parts := strings.Split(ipPort, ":")
	if len(parts) != 2 {
		return errors.New("Param addr invalid. expected ip:port: " + addr)
	}
	c._addr = parts[0]

	portStr := parts[1]
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return errors.New("Param port invalid. : " + portStr)
	}
	c._port = uint16(port)
	return nil
}

// ToString 连接地址转换为字符串
// @receiver c 连接地址
// @return string 连接地址字符串
func (c *Addr) ToString() string {
	return c._type + "://" + c._addr + ":" + strconv.FormatUint(uint64(c._port), 10)
}

func (c *Addr) GetType() Type {
	return Type(c._type)
}
