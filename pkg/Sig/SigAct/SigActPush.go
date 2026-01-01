package SigAct

import (
	"net/http"
)

/** -------------------------------------------- EXT --------------------------------------------- */

type ActionPush struct {
	// 静态资源根目录
	_static string
}

func PushUrl() string {
	return "/rtc/push"
}

func PushNew(static string) *ActionPush {
	return &ActionPush{
		_static: static + "/rtcPush.html",
	}
}

func (act *ActionPush) Act(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// 响应静态文件
	http.ServeFile(w, r, act._static)
}
