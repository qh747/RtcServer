package SigConn

import (
	"context"
	rpc "rtcServer/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SigRpcConn rpc连接
type SigRpcConn struct {
	_addr string
	_msg  string
}

// NewSigRpcConn 		创建rpc连接
// @param addr 			rpc服务地址
// @param msg 			rpc请求消息
// @return *SigRpcConn	rpc连接
func NewSigRpcConn(addr string, msg string) *SigRpcConn {
	return &SigRpcConn{_addr: addr, _msg: msg}
}

// GetType 			获取连接类型
// @receiver c 		rpc连接
// @return connType rpc连接类型
func (c *SigRpcConn) GetType() connType {
	return connRpc
}

// Req 				rpc请求
// @receiver c 		rpc连接
// @return string 	rpc响应
// @return error 	rpc请求是否存在错误
func (c *SigRpcConn) Req() (string, error) {
	conn, err := grpc.NewClient(c._addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if nil != err {
		return "", err
	}
	defer conn.Close()

	if resp, err := rpc.NewRpcConnClient(conn).RtcPush(context.Background(), &rpc.RtcPushReqArgs{Msg: c._msg}); nil != err {
		return "", err
	} else {
		return resp.GetMsg(), nil
	}
}

// ReqAsync 	异步rpc请求
// @receiver c 	rpc连接
// @param cb 	rpc请求回调函数
func (c *SigRpcConn) ReqAsync(cb func(resp string, err error)) {
	go func() {
		resp, err := c.Req()
		cb(resp, err)
	}()
}
