package SigAct

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rtcServer/pkg/Com/Log"
)

type ActionError struct {
	_code int
	_info string
}

type ErrInvalidReuqest struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (act *ActionError) Act(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(act._code)
	fmt.Fprintf(w, "%s", act._info)
}

func ActErrNotfound(w http.ResponseWriter, r *http.Request) {
	err := ActionError{
		_code: http.StatusNotFound,
		_info: http.StatusText(http.StatusNotFound),
	}

	Log.Log().Errorf("Action not found error. request: %s", DumpAction(r))
	err.Act(w, r)
}

func ActErrInternalError(w http.ResponseWriter, r *http.Request) {
	err := ActionError{
		_code: http.StatusInternalServerError,
		_info: http.StatusText(http.StatusInternalServerError),
	}

	Log.Log().Errorf("Action server internal error. request: %s", DumpAction(r))
	err.Act(w, r)
}

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

func ActErrOther(w http.ResponseWriter, r *http.Request, code int, info string) {
	err := ActionError{
		_code: code,
		_info: http.StatusText(code),
	}

	Log.Log().Errorf("Action other error. code: %d, info: %s, request: %s", code, info, DumpAction(r))
	err.Act(w, r)
}
