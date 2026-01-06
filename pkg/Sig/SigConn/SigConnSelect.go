package SigConn

import (
	"errors"
	"rtcServer/pkg/Com/Log"
	"strings"
	"sync"
)

// connSelector	连接选择器
type connSelector struct {
	_idx      int
	_addr     []*Addr
	_lock     sync.RWMutex
	_bindAddr map[string]*Addr
}

// selector 全局连接选择器变量
var Selector *connSelector

// InitConnSelect 初始化连接选择器
// @param param 初始化参数
// @return error 初始化是否存在错误
func InitConnSelect(param string) error {
	if nil == Selector {
		Selector = &connSelector{
			_idx: 0,
		}

		if err := Selector.setAddr(param); nil != err {
			Selector = nil
			return err
		}
	}
	return nil
}

// setAddr 设置连接地址
// @receiver s 连接选择器
// @param param 连接地址
// @return error 设置是否存在错误
func (s *connSelector) setAddr(param string) error {
	s._lock.RLock()
	defer s._lock.RUnlock()

	for val := range strings.SplitSeq(param, ",") {
		addr := new(Addr)
		if err := addr.LoadFrom(val); nil != err {
			Log.Logger.Warnf("Load addr warning. addr: %s. err: %e\n", val, err)
			continue
		}
		s._addr = append(s._addr, addr)
	}

	if 0 == len(s._addr) {
		return errors.New("addr list empty")
	}
	return nil
}

// GetAddr 获取连接地址
// @receiver s 连接选择器
// @param id 连接id
// @return Addr 连接地址
func (s *connSelector) GetAddr(id string) Addr {
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
