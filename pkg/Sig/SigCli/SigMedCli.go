package SigCli

import (
	"errors"
	"rtcServer/pkg/Com/Log"
	"strconv"
	"strings"
	"sync"
)

// 初始化媒体服务连接
// return 初始化结果
// param  param 初始化参数
func InitMedCli(param string) error {
	if nil == selector {
		var err error
		if selector, err = newSelector(param); nil != err {
			return err
		}
	}

	return nil
}

// 连接类型
const (
	connRpc   = "rpc"
	connHttp  = "http"
	connHttps = "https"
)

// 连接地址
type connAddr struct {
	_type string
	_addr string
	_port uint16
}

// 解析连接地址
// return 解析是否存在错误
// param  addr 待解析连接地址
func (c *connAddr) loadFrom(addr string) error {
	if "" == addr {
		return errors.New("Param empty")
	}

	// 解析连接类型
	var ipPort string
	if strings.HasPrefix(addr, "rpc://") {
		c._type = connRpc
		ipPort = strings.TrimPrefix(addr, "rpc://")
	} else if strings.HasPrefix(addr, "http://") {
		c._type = connHttp
		ipPort = strings.TrimPrefix(addr, "http://")
	} else if strings.HasPrefix(addr, "https://") {
		c._type = connHttps
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

// 格式化连接地址
// return 连接地址
func (c *connAddr) toString() string {
	return c._type + "://" + c._addr + ":" + strconv.FormatUint(uint64(c._port), 10)
}

// 创建连接选择器
// return 连接选择器地址，创建是否存在错误
// param  param 创建参数
func newSelector(param string) (*connSelector, error) {
	s := new(connSelector)
	s._idx = 0

	if err := s.setAddr(param); nil != err {
		return nil, err
	}

	return s, nil
}

// 连接选择器
type connSelector struct {
	_idx      int
	_addr     []*connAddr
	_lock     sync.RWMutex
	_bindAddr map[string]*connAddr
}

// 设置连接地址
// return 设置是否存在错误
// param  param 连接地址
func (s *connSelector) setAddr(param string) error {
	s._lock.RLock()
	defer s._lock.RUnlock()

	for val := range strings.SplitSeq(param, ",") {
		addr := new(connAddr)
		if err := addr.loadFrom(val); nil != err {
			Log.Log().Warnf("Load addr warning. addr: %s. err: %e\n", val, err)
			continue
		}
		s._addr = append(s._addr, addr)
	}

	if 0 == len(s._addr) {
		return errors.New("addr list empty")
	}
	return nil
}

// 获取连接地址
// return 连接地址
// param  id 连接地址关联id
func (s *connSelector) getAddr(id string) string {
	s._lock.RLock()
	defer s._lock.RUnlock()

	// 优先从已绑定列表中查找
	addr, ok := s._bindAddr[id]
	if ok && nil != addr {
		return addr.toString()
	}

	// 如果未绑定则从新增绑定关系
	s._idx++
	if s._idx >= len(s._addr) {
		s._idx = 0
	}

	addr = s._addr[s._idx]
	s._bindAddr[id] = addr

	return addr.toString()
}

// 全局连接选择器变量
var selector *connSelector
