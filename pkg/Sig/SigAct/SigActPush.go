package SigAct

import (
	"net/http"
)

func ActionPushUrl() string {
	return "/rtc/push"
}

func NewActionPush() *ActionPush {
	return new(ActionPush)
}

type ActionPush struct {
}

func (act *ActionPush) Execute(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// 响应静态文件
	http.ServeFile(w, r, "../web/static/rtcPush.html")
}
