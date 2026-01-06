package SigAct

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type (
	// Action http请求处理接口
	Action interface {
		// 执行响应
		Act(w http.ResponseWriter, r *http.Request)
	}

	// ActionMap http请求处理接口map
	ActionMap = map[string]Action
	// ActionMapPtr http请求处理接口map指针
	ActionMapPtr = *ActionMap
)

// DumpAction 输出http请求
// @param r http请求
// @return string http请求字符串
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
