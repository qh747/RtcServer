package SigAct

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rtcServer/pkg/Com/Log"
)

// ActionError http请求处理错误
type ActionError struct {
	_code int
	_info string
}

// ErrInvalidReuqest 无效请求错误
type ErrInvalidReuqest struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Act 				执行响应
// @receiver act 	http请求处理错误
// @param w 		http响应
// @param r 		http请求
func (act *ActionError) Act(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(act._code)
	fmt.Fprintf(w, "%s", act._info)
}

// ActErrNotfound	http请求404错误
// @param w 		http响应
// @param r 		http请求
func ActErrNotfound(w http.ResponseWriter, r *http.Request) {
	err := ActionError{
		_code: http.StatusNotFound,
		_info: http.StatusText(http.StatusNotFound),
	}

	Log.Log().Errorf("Action not found error. request: %s", DumpAction(r))
	err.Act(w, r)
}

// ActErrInternalError	http请求500错误
// @param w 			http响应
// @param r 			http请求
func ActErrInternalError(w http.ResponseWriter, r *http.Request) {
	err := ActionError{
		_code: http.StatusInternalServerError,
		_info: http.StatusText(http.StatusInternalServerError),
	}

	Log.Log().Errorf("Action server internal error. request: %s", DumpAction(r))
	err.Act(w, r)
}

// ActErrInvalidRequest	http请求无效错误
// @param w 			http响应
// @param r 			http请求
// @param msg 			错误信息
func ActErrInvalidRequest(w http.ResponseWriter, r *http.Request, msg string) {
	errBody, _ := json.Marshal(ErrInvalidReuqest{
		Code: -1,
		Msg:  "Reuqest invalid",
	})

	err := ActionError{
		_code: http.StatusBadRequest,
		_info: string(errBody),
	}

	Log.Log().Errorf("Action invalid request error. request: %s", DumpAction(r))
	err.Act(w, r)
}

// ActErrOther	http请求其他错误
// @param w 	http响应
// @param r 	http请求
// @param code 	错误码
// @param info 	错误信息
func ActErrOther(w http.ResponseWriter, r *http.Request, code int, info string) {
	err := ActionError{
		_code: code,
		_info: http.StatusText(code),
	}

	Log.Log().Errorf("Action other error. code: %d, info: %s, request: %s", code, info, DumpAction(r))
	err.Act(w, r)
}
