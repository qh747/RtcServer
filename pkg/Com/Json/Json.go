package Json

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

/** -------------------------------------------- EXT --------------------------------------------- */

type PushReq struct {
	// 房间id
	_roomId string
	// 用户id
	_userId string
	// 请求类型
	_type string
	// sdp
	_sdp string
}

type PushResp struct {
	// 响应码
	_code string
	// 消息
	_msg string
}

func NewPushReq(r *http.Request) (*PushReq, error) {
	roomId := r.URL.Query().Get("room")
	if "" == roomId {
		return nil, fmt.Errorf("Room id invalid")
	}

	userId := r.URL.Query().Get("user")
	if "" == userId {
		return nil, fmt.Errorf("User id invalid")
	}

	str, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Request body invalid")
	}

	var body pushReqBody
	if err := json.Unmarshal(str, &body); err != nil {
		return nil, fmt.Errorf("Request body invalid. body: %s", string(str))
	}

	return &PushReq{
		_roomId: roomId,
		_userId: userId,
		_type:   body.Type,
		_sdp:    body.Sdp,
	}, nil
}

func NewPushResq(code string, msg string) (*PushResp, error) {
	if "" == code {
		return nil, fmt.Errorf("Code invalid")
	}

	if "" == msg {
		return nil, fmt.Errorf("Msg invalid")
	}

	return &PushResp{
		_code: code,
		_msg:  msg,
	}, nil
}

func (resp *PushResp) ToString() string {
	if "0" == resp._code {
		succResp := pushRespBody{
			Code: resp._code,
			Sdp:  resp._msg,
		}

		str, err := json.Marshal(succResp)
		if err != nil {
			return ""
		}
		return string(str)
	} else {
		errResp := pushErrRespBody{
			Code: resp._code,
			Msg:  resp._msg,
		}

		str, err := json.Marshal(errResp)
		if err != nil {
			return ""
		}
		return string(str)
	}
}

/** -------------------------------------------- IN --------------------------------------------- */

type pushReqBody struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
}

type pushRespBody struct {
	Code string `json:"code"`
	Sdp  string `json:"sdp"`
}

type pushErrRespBody struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}
