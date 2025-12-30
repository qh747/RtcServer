package SigAct

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type Action interface {
	// 执行响应
	Execute(w http.ResponseWriter, r *http.Request)
}

// key = url, value response action
type ActionMap = map[string]Action
type ActionMapPtr = *ActionMap

// 输出http请求
// return http请求字符串
// r      http请求
func DumpAction(r *http.Request) string {
	if r == nil {
		return "<nil request>"
	}

	// true 表示包含 body
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		return fmt.Sprintf("Dump request error: %v", err)
	}
	return string(dump)
}
