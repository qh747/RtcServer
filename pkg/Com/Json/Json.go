package Json

import (
	"encoding/json"
	"fmt"
)

func NewPushReq(str []byte) (*PushReq, error) {
	var req PushReq
	if err := json.Unmarshal(str, &req); err != nil {
		return nil, fmt.Errorf("Request body invalid. body: %s", string(str))
	}
	return &req, nil
}

func NewPushResq(code string, msg string) (*PushResp, error) {
	return &PushResp{
		Code: code,
		Msg:  msg,
	}, nil
}

type PushReq struct {
	// 房间id
	Room string `json:"room"`
	// 用户id
	User string `json:"user"`
	// 请求类型
	Type string `json:"type"`
	// 消息
	Msg string `json:"msg"`
}

type PushResp struct {
	// 响应码
	Code string `json:"code"`
	// 消息
	Msg string `json:"msg"`
}

func (resp *PushResp) ToString() string {
	str, err := json.Marshal(resp)
	if err != nil {
		return ""
	}
	return string(str)
}
