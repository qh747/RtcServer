package SigAct

import (
	"net/http"
	"path/filepath"
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
		} else {
			// 默认使用http.ServeFile自动检测Content-Type
		}

		http.ServeFile(w, r, filePath)
	} else {
		http.NotFound(w, r)
	}
}
