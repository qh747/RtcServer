package SigServ

import (
	"net/http"
	"rtcServer/pkg/Com/Log"
	"rtcServer/pkg/Sig/SigAct"
	"strings"
)

// SignalServer	信令服务
type SignalServer struct {
	// 静态资源根目录
	_static string

	// HTTP服务器实例
	_impl *http.Server

	// 请求处理回调函数map
	_acts SigAct.ActionMap
}

// SignalSslServer 信令加密服务
type SignalSslServer struct {
	// 密钥
	_key string

	// 证书
	_cert string

	// 信令服务器实例
	_serv *SignalServer
}

// NewSigServ 				创建信令服务
// @param addr 				服务地址
// @param static			静态资源根目录
// @return *SignalServer	信令服务
func NewSigServ(addr string, static string) *SignalServer {
	// 创建服务
	serv := SignalServer{
		_static: static,
		_acts: map[string]SigAct.Action{
			SigAct.PushUrl(): SigAct.PushNew(static),
		},
	}

	serv._impl = &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for url, act := range serv._acts {
				if url == r.RequestURI || strings.HasPrefix(r.RequestURI, url+"/") {
					// 找到响应
					if nil == act {
						// 响应无效
						SigAct.ActErrInternalError(w, r)
					} else {
						// 处理响应
						act.Act(w, r)
					}
					return
				}
			}

			// 未找到
			SigAct.ActErrNotfound(w, r)
		}),
	}

	return &serv
}

// Start 			启动服务
// @receiver serv 	信令服务
func (serv *SignalServer) Start() {
	Log.Log().Infof("Start signal server. listen: %s", serv._impl.Addr)
	if err := serv._impl.ListenAndServe(); nil != err && http.ErrServerClosed != err {
		Log.Log().Errorf("Start signal server error. addr: %s, err: %v", serv._impl.Addr, err)
		return
	}
}

// Stop 			停止服务
// @receiver serv 	信令服务
func (serv *SignalServer) Stop() {
	if serv._impl != nil {
		serv._impl.Close()
	}
}

// NewSigSslServ 				创建信令加密服务
// @param addr 					服务地址
// @param static 				静态资源根目录
// @param key 					密钥
// @param cert 					证书
// @return *SignalSslServer		信令加密服务
func NewSigSslServ(addr string, static string, key string, cert string) *SignalSslServer {
	// 创建服务
	serv := SignalSslServer{
		_key:  key,
		_cert: cert,
		_serv: NewSigServ(addr, static),
	}

	return &serv
}

// Start 			启动加密服务
// @receiver serv 	信令加密服务
func (serv *SignalSslServer) Start() {
	Log.Log().Infof("Start signal ssl server. listen: %s", serv._serv._impl.Addr)
	if err := serv._serv._impl.ListenAndServeTLS(serv._cert, serv._key); nil != err && http.ErrServerClosed != err {
		Log.Log().Errorf("Start signal ssl server error. addr: %s, err: %v", serv._serv._impl.Addr, err)
		return
	}
}

// Stop 			停止加密服务
// @receiver serv 	信令加密服务
func (serv *SignalSslServer) Stop() {
	serv._serv.Stop()
}
