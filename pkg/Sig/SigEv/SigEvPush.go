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
	reqMsg := args[1].(string)
	conn := SigConn.NewRpcConn(addr.ToString(), reqMsg)

	cb := args[2].(func(resp string, err error))
	conn.ReqAsync(cb)
}

// onHttpPushEvent 使用http方式推流
// @param addr 服务连接
// @param args 事件参数
func onHttpPushEvent(addr SigConn.Addr, args ...any) {
	url := addr.ToString() + "/push/start"

	head := map[string]string{
		"Content-Type": "application/json",
	}

	reqMsg := args[1].(string)
	conn := SigConn.NewHttpConn(url, "POST", head, reqMsg)

	cb := args[2].(func(resp string, err error))
	conn.ReqAsync(cb)
}

// onHttpsPushEvent 使用https方式推流
// @param addr 服务连接
// @param args 事件参数
func onHttpsPushEvent(addr SigConn.Addr, args ...any) {
	url := addr.ToString() + "/push/start"

	head := map[string]string{
		"Content-Type": "application/json",
	}

	reqMsg := args[1].(string)
	conn := SigConn.NewHttpsConn(url, "POST", head, reqMsg)

	cb := args[2].(func(resp string, err error))
	conn.ReqAsync(cb)
}
