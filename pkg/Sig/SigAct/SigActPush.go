package SigAct

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"rtcServer/pkg/Com/Log"
	"rtcServer/pkg/Sig/SigEv"
	"strings"
)

type (
	// ActionPush 推流请求处理接口
	ActionPush struct {
		// 静态资源根目录
		_static string
	}

	// PushReuqest 推流请求
	PushReuqest struct {
		Room string `json:"room"`
		User string `json:"user"`
		Type string `json:"type"`
		Msg  string `json:"msg"`
	}
)

// PushUrl 推流请求路径
// @return string url
func PushUrl() string {
	return "/rtc/push"
}

// PushNew 创建推流请求处理接口
// @param static 静态资源根目录
// @return *ActionPush 推流请求处理接口
func PushNew(static string) *ActionPush {
	return &ActionPush{
		_static: static,
	}
}

// Act 执行响应
// @receiver act 推流请求处理接口
// @param w http响应
// @param r http请求
func (act *ActionPush) Act(w http.ResponseWriter, r *http.Request) {
	if "GET" == r.Method {
		act.actGet(w, r)
	} else if "POST" == r.Method {
		act.actPost(w, r)
	} else {
		Log.Logger.Errorf("Action push error. request method invalid. request: %s", DumpAction(r))
		ActErrNotfound(w, r)
	}
}

// actGet GET请求响应
// @receiver act 推流请求处理接口
// @param w http响应
// @param r http请求
func (act *ActionPush) actGet(w http.ResponseWriter, r *http.Request) {
	if PushUrl() == r.RequestURI {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, act._static+"/rtcPush.html")
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
		Log.Logger.Errorf("Action push get error. request url invalid. request: %s", DumpAction(r))
		http.NotFound(w, r)
	}
}

// actPost POST请求响应
// @receiver act 推流请求处理接口
// @param w http响应
// @param r http请求
func (act *ActionPush) actPost(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.RequestURI, PushUrl()+"/start") {
		// 读取请求
		body, err := io.ReadAll(r.Body)
		if err != nil {
			ActErrInvalidRequest(w, r, "Load request body failed")
			return
		}

		// 解析请求
		var pushReq PushReuqest
		if err := json.Unmarshal([]byte(body), &pushReq); nil != err {
			ActErrInvalidRequest(w, r, "Request body format invalid")
			return
		}
		Log.Logger.Infof("Receive push request. room: %s. user: %s. type: %s. msg: %s", pushReq.Room, pushReq.User, pushReq.Type, pushReq.Msg)

		// 转发请求给媒体服务
		SigEv.Dispatcher.Publish(SigEv.EvTopic(SigEv.EvTopicPush), pushReq.Room, body, func(resp string, err error) {
			if nil != err {
				ActErrInvalidRequest(w, r, err.Error())
			} else {
				Log.Logger.Infof("Receive push response. msg: %s", resp)
			}
		})
	} else {
		Log.Logger.Errorf("Action push post error. request url invalid. request: %s", DumpAction(r))
		ActErrNotfound(w, r)
	}
}
