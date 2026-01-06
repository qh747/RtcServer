package SigEv

import (
	"rtcServer/pkg/Com/Log"
	"rtcServer/pkg/Sig/SigConn"
)

// OnPushEvent 推流事件
// @param args 事件参数
func OnPushEvent(args ...any) {
	// 获取服务连接
	id := args[0].(string)
	addr := SigConn.Selector.GetAddr(id)

	switch addr.GetType() {
	case SigConn.Rpc:
		onRpcPushEvent(addr, args...)
	case SigConn.Http:
		onHttpPushEvent(addr, args...)
	case SigConn.Https:
		onHttpsPushEvent(addr, args...)
	default:
		Log.Logger.Errorf("Push event warning. conn type invalid. addr: %s.\n", addr.ToString())
	}
}

// onRpcPushEvent 使用rpc方式推流
// @param addr 服务连接
// @param args 事件参数
func onRpcPushEvent(addr SigConn.Addr, args ...any) {
	if len(args) < 3 {
		Log.Logger.Errorf("On rpc push event error. arguments insufficient. expected at least 3, got %d", len(args))
		return
	}

	reqMsg, ok := args[1].(string)
	if !ok {
		Log.Logger.Errorf("On rpc push event error. argument[1] type invalid. expected string. got %T", args[1])
		return
	}

	cb, ok := args[2].(func(resp string, err error))
	if !ok {
		Log.Logger.Errorf("On rpc push event error. argument[1] type invalid. expected func(resp string, err error). got %T", args[2])
		return
	}

	conn := SigConn.NewRpcConn(addr.ToString(), reqMsg)
	conn.ReqAsync(cb)
}

// onHttpPushEvent 使用http方式推流
// @param addr 服务连接
// @param args 事件参数
func onHttpPushEvent(addr SigConn.Addr, args ...any) {
	if len(args) < 3 {
		Log.Logger.Errorf("On http push event error. arguments insufficient. expected at least 3. got %d", len(args))
		return
	}

	reqMsg, ok := args[1].(string)
	if !ok {
		Log.Logger.Errorf("On http push event error. argument[1] type invalid. expected string. got %T", args[1])
		return
	}

	cb, ok := args[2].(func(resp string, err error))
	if !ok {
		Log.Logger.Errorf("On http push event error. argument[1] type invalid. expected func(resp string, err error). got %T", args[2])
		return
	}

	url := addr.ToString() + "/push/start"

	head := map[string]string{
		"Content-Type": "application/json",
	}

	conn := SigConn.NewHttpConn(url, "POST", head, reqMsg)
	conn.ReqAsync(cb)
}

// onHttpsPushEvent 使用https方式推流
// @param addr 服务连接
// @param args 事件参数
func onHttpsPushEvent(addr SigConn.Addr, args ...any) {
	if len(args) < 3 {
		Log.Logger.Errorf("On https push event error. arguments insufficient. expected at least 3, got %d", len(args))
		return
	}

	reqMsg, ok := args[1].(string)
	if !ok {
		Log.Logger.Errorf("On https push event error. argument[1] type invalid. expected string. got %T", args[1])
		return
	}

	cb, ok := args[2].(func(resp string, err error))
	if !ok {
		Log.Logger.Errorf("On https push event error. argument[1] type invalid. expected func(resp string, err error). got %T", args[2])
		return
	}

	url := addr.ToString() + "/push/start"

	head := map[string]string{
		"Content-Type": "application/json",
	}

	conn := SigConn.NewHttpsConn(url, "POST", head, reqMsg)
	conn.ReqAsync(cb)
}
