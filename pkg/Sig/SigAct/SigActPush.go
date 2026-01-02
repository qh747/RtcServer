package SigAct

import (
	"fmt"
	"net/http"
	"path/filepath"
	"rtcServer/pkg/Com/Json"
	"rtcServer/pkg/Com/Log"
	"strings"
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
		_static: static,
	}
}

func (act *ActionPush) Act(w http.ResponseWriter, r *http.Request) {
	if "GET" == r.Method {
		act.actGet(w, r)
	} else if "POST" == r.Method {
		act.actPost(w, r)
	} else {
		Log.Log().Errorf("Action push error. request method invalid. request: %s", DumpAction(r))
		ActErrNotfound(w, r)
	}
}

func (act *ActionPush) actGet(w http.ResponseWriter, r *http.Request) {
	if PushUrl() == r.RequestURI {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, act._static+"/push.html")
	} else if strings.HasPrefix(r.RequestURI, PushUrl()+"/") {
		// 获取请求路径中除前缀外的部分
		relativePath := r.RequestURI[len(PushUrl()):]

		// 移除开头的斜杠
		if len(relativePath) > 0 && relativePath[0] == '/' {
			relativePath = relativePath[1:]
		}

		// 构造实际文件路径
		filePath := filepath.Join(act._static, relativePath)

		// 检查路径安全性，防止目录遍历攻击
		if !strings.HasPrefix(filePath, act._static) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// 根据文件扩展名设置正确的Content-Type
		if strings.HasSuffix(filePath, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(filePath, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(filePath, ".html") || strings.HasSuffix(filePath, ".htm") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		} else if strings.HasSuffix(filePath, ".png") {
			w.Header().Set("Content-Type", "image/png")
		} else if strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg") {
			w.Header().Set("Content-Type", "image/jpeg")
		} else if strings.HasSuffix(filePath, ".gif") {
			w.Header().Set("Content-Type", "image/gif")
		}

		http.ServeFile(w, r, filePath)
	} else {
		Log.Log().Errorf("Action push get error. request url invalid. request: %s", DumpAction(r))
		http.NotFound(w, r)
	}
}

func (act *ActionPush) actPost(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.RequestURI, PushUrl()+"/start") {
		if _, err := Json.NewPushReq(r); nil != err {
			ActErrInvalidPushRequest(w, r, -1, "Load request error")
			return
		}

		// 发送成功响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		pushResp, _ := Json.NewPushResq("0", "success")
		fmt.Fprintf(w, "%s", pushResp.ToString())
	} else {
		Log.Log().Errorf("Action push post error. request url invalid. request: %s", DumpAction(r))
		ActErrNotfound(w, r)
	}
}

/** -------------------------------------------- IN --------------------------------------------- */
