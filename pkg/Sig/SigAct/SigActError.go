package SigAct

import (
	"fmt"
	"net/http"
	"rtcServer/pkg/Log"
)

type ActionError struct {
	_code int
	_info string
}

func (act *ActionError) Execute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(act._code)
	fmt.Fprintf(w, "%d - %s", act._code, act._info)
}

func ActErrNotfound(w http.ResponseWriter, r *http.Request) {
	err := ActionError{
		_code: http.StatusNotFound,
		_info: "Not Found",
	}

	Log.Log().Errorf("Action not found error. request: %s", DumpAction(r))
	err.Execute(w, r)
}

func ActErrInternalError(w http.ResponseWriter, r *http.Request) {
	err := ActionError{
		_code: http.StatusInternalServerError,
		_info: "Server Internal Error",
	}

	Log.Log().Errorf("Action server internal error. request: %s", DumpAction(r))
	err.Execute(w, r)
}

func ActErrOther(w http.ResponseWriter, r *http.Request, code int, info string) {
	err := ActionError{
		_code: code,
		_info: info,
	}

	Log.Log().Errorf("Action other error. code: %d, info: %s, request: %s", code, info, DumpAction(r))
	err.Execute(w, r)
}
