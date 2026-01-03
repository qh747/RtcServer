package SigConn

import (
	"context"
	"errors"
	"rtcServer/pkg/Com/Json"
	rpc "rtcServer/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewSigRpcPushConn(addr string, room string, user string, typeStr string, msg string) *SigRpcConn {
	return &SigRpcConn{
		_addr:   addr,
		_method: MehtodPushReq,
		_data: map[string]string{
			"room": room,
			"user": user,
			"type": typeStr,
			"msg":  msg,
		},
	}
}

const (
	MehtodPushReq = "PushReq"
)

type SigRpcConn struct {
	_addr   string
	_method string
	_data   map[string]string
}

func (c *SigRpcConn) Req() (string, error) {
	if MehtodPushReq == c._method {
		return c.pushReq()
	} else {
		return "", errors.New("Method invalid")
	}
}

func (c *SigRpcConn) ReqAsync(cb func(resp string, err error)) {
	go func() {
		resp, err := c.Req()
		cb(resp, err)
	}()
}

func (c *SigRpcConn) pushReq() (string, error) {
	room, ok := c._data["room"]
	if !ok || "" == room {
		return "", errors.New("Room invalid")
	}

	user, ok := c._data["user"]
	if !ok || "" == user {
		return "", errors.New("User invalid")
	}

	typeStr, ok := c._data["type"]
	if !ok || "" == typeStr {
		return "", errors.New("Type invalid")
	}

	msg, ok := c._data["msg"]
	if !ok || "" == msg {
		return "", errors.New("Msg invalid")
	}

	conn, err := grpc.NewClient(c._addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if nil != err {
		return "", err
	}

	defer conn.Close()

	cli := rpc.NewRpcConnClient(conn)
	resp, err := cli.RtcPush(context.Background(), &rpc.RtcPushReqArgs{
		Room: room,
		User: user,
		Type: typeStr,
		Msg:  msg,
	})

	if nil != err {
		return "", err
	}

	respBody, err := Json.NewPushResq(resp.GetCode(), resp.GetMsg())
	if nil != err {
		return "", err
	}
	return respBody.ToString(), nil
}
