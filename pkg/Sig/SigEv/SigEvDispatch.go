package SigEv

import (
	"errors"
	"sync"
)

// const
// @param EvTopicPush 推送事件
const (
	EvTopicPush string = "EvTopicPush"
)

type (
	// EvTopic 事件主题
	EvTopic string

	// EvHandler 事件处理函数
	EvHandler func(args ...any)

	// EvDispatch 事件分发器
	EvDispatch struct {
		_id       string
		_lock     sync.RWMutex
		_handlers map[EvTopic]map[string]EvHandler
	}
)

// newDispatch 创建事件分发器
// @param id 事件分发器id
// @return *EvDispatch
// @return error 创建是否存在错误
func newDispatch(id string) (*EvDispatch, error) {
	if "" == id {
		return nil, errors.New("Dispatch id empty")
	}
	return &EvDispatch{
		_id:       id,
		_handlers: make(map[EvTopic]map[string]EvHandler),
	}, nil
}

// Subscribe 订阅事件
// @receiver dispatch 事件分发器
// @param topic 事件主题
// @param id 事件监听者id
// @param handler 事件处理函数
func (dispatch *EvDispatch) Subscribe(topic EvTopic, id string, handler EvHandler) {
	dispatch._lock.Lock()
	defer dispatch._lock.Unlock()

	if nil == dispatch._handlers[topic] {
		dispatch._handlers[topic] = make(map[string]EvHandler)
	}

	dispatch._handlers[topic][id] = handler
}

// Unsubscribe 取消订阅事件
// @receiver dispatch 事件分发器
// @param topic 事件主题
// @param id 事件监听者id
func (dispatch *EvDispatch) Unsubscribe(topic EvTopic, id string) {
	dispatch._lock.Lock()
	defer dispatch._lock.Unlock()

	if nil != dispatch._handlers[topic] {
		delete(dispatch._handlers[topic], id)
	}
}

// Publish 发布事件
// @receiver dispatch 事件分发器
// @param topic 事件主题
// @param args 事件参数
func (dispatch *EvDispatch) Publish(topic EvTopic, args ...any) {
	dispatch._lock.RLock()
	defer dispatch._lock.RUnlock()

	if nil != dispatch._handlers[topic] {
		for _, handler := range dispatch._handlers[topic] {
			go handler(args...)
		}
	}
}
