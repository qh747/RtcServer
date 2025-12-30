package SigServ

import (
	"net/http"
	"rtcServer/pkg/Log"
	"rtcServer/pkg/Sig/SigAct"
)

// 信令服务
type SignalServer struct {
	// 服务地址
	_addr string

	// 请求处理回调函数map
	_acts SigAct.ActionMap
}

/** -------------------------------------------- EXT --------------------------------------------- */

// 创建服务
// return 信令服务
// addr   服务地址
func New(addr string) *SignalServer {
	serv := new(SignalServer)

	// 设置服务地址
	serv._addr = addr

	// 注册请求处理回调函数
	serv.registAction()
	return serv
}

// 启动服务
func (serv *SignalServer) Start() {
	// 启动http服务
	Log.Log().Infof("Start signal server. listen: %s", serv._addr)
	err := http.ListenAndServe(serv._addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serv.execute(w, r)
	}))

	if nil != err {
		Log.Log().Errorf("Start signal server error. addr: %s, err: %v", serv._addr, err)
		return
	}
}

/** -------------------------------------------- IN --------------------------------------------- */

// 注册请求处理回调函数
func (serv *SignalServer) registAction() {
	serv._acts = make(SigAct.ActionMap)

	// 注册推流url
	serv._acts[SigAct.ActionPushUrl()] = SigAct.NewActionPush()
}

// 处理请求
// w 响应句柄
// r 请求内容
func (serv *SignalServer) execute(w http.ResponseWriter, r *http.Request) {
	act, ok := serv._acts[r.RequestURI]
	switch {
	case !ok:
		// 未找到
		SigAct.ActErrNotfound(w, r)
	case nil == act:
		// 响应无效
		SigAct.ActErrInternalError(w, r)
	default:
		// 处理响应
		act.Execute(w, r)
	}
}
