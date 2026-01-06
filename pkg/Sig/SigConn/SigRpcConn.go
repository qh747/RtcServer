package SigConn

import (
	"context"
	rpc "rtcServer/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RpcConn rpc连接
type RpcConn struct {
	_addr string
	_msg  string
}

// NewRpcConn 创建rpc连接
// @param addr rpc服务地址
// @param msg rpc请求消息
// @return *RpcConn
func NewRpcConn(addr string, msg string) *RpcConn {
	return &RpcConn{_addr: addr, _msg: msg}
}

// Req rpc请求
// @receiver c rpc连接
// @return string rpc响应
// @return error rpc请求是否存在错误
func (c *RpcConn) Req() (string, error) {
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

// ReqAsync 异步rpc请求
// @receiver c rpc连接
// @param cb rpc请求回调函数
func (c *RpcConn) ReqAsync(cb func(resp string, err error)) {
	go func() {
		resp, err := c.Req()
		cb(resp, err)
	}()
}
