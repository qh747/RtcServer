package SigServ

import "net/http"

type Action interface {
	// 执行响应
	Execute(w http.ResponseWriter, r *http.Request)
}

// key = url, value response action
type ActionMap = map[string]Action
type ActionMapPtr = *ActionMap
