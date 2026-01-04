package SigConn

import (
	"errors"
	"rtcServer/pkg/Com/Log"
	"strconv"
	"strings"
	"sync"
)

// connAddr 连接地址
type connAddr struct {
	_type string
	_addr string
	_port uint16
}

// connSelector	连接选择器
type connSelector struct {
	_idx      int
	_addr     []*connAddr
	_lock     sync.RWMutex
	_bindAddr map[string]*connAddr
}

// selector 全局连接选择器变量
var selector *connSelector

// InitSigConnSelect	初始化连接选择器
// @param param 		初始化参数
// @return error 		初始化是否存在错误
func InitSigConnSelect(param string) error {
	if nil == selector {
		var err error
		if selector, err = newSelector(param); nil != err {
			return err
		}
	}
	return nil
}

// GetSelector 				获取连接选择器
// @return *connSelector 	连接选择器
func GetSelector() *connSelector {
	return selector
}

// loadFrom 		加载连接地址
// @receiver c		连接地址
// @param addr		连接地址
// @return error	加载是否存在错误
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

// toString 		连接地址转换为字符串
// @receiver c 		连接地址
// @return string 	连接地址字符串
func (c *connAddr) toString() string {
	return c._type + "://" + c._addr + ":" + strconv.FormatUint(uint64(c._port), 10)
}

// newSelector 				创建连接选择器
// @param param 			连接地址
// @return *connSelector 	连接选择器
// @return error 			创建是否存在错误
func newSelector(param string) (*connSelector, error) {
	s := new(connSelector)
	s._idx = 0

	if err := s.setAddr(param); nil != err {
		return nil, err
	}

	return s, nil
}

// setAddr 			设置连接地址
// @receiver s 		连接选择器
// @param param		连接地址
// @return error	设置是否存在错误
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

// getAddr 			获取连接地址
// @receiver s 		连接选择器
// @param id 		连接id
// @return connAddr	连接地址
func (s *connSelector) GetAddr(id string) connAddr {
	s._lock.RLock()
	defer s._lock.RUnlock()

	// 优先从已绑定列表中查找
	addr, ok := s._bindAddr[id]
	if ok && nil != addr {
		return *addr
	}

	// 如果未绑定则从新增绑定关系
	s._idx++
	if s._idx >= len(s._addr) {
		s._idx = 0
	}

	addr = s._addr[s._idx]
	s._bindAddr[id] = addr

	return *addr
}
