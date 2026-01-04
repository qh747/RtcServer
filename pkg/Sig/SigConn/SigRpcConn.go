package SigConn

import (
	"context"
	rpc "rtcServer/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SigRpcConn struct {
	_addr string
	_msg  string
}

func NewSigRpcConn(addr string, msg string) *SigRpcConn {
	return &SigRpcConn{_addr: addr, _msg: msg}
}

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

func (c *SigRpcConn) ReqAsync(cb func(resp string, err error)) {
	go func() {
		resp, err := c.Req()
		cb(resp, err)
	}()
}
