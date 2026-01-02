package SigAct

import (
	"fmt"
	"net/http"
	"rtcServer/pkg/Com/Json"
	"rtcServer/pkg/Com/Log"
	"strconv"
)

/** -------------------------------------------- EXT --------------------------------------------- */

type ActionError struct {
	_code int
	_info string
}

func (act *ActionError) Act(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(act._code)
	fmt.Fprintf(w, "%d - %s", act._code, act._info)
}

func (act *ActionError) ActJson(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(act._code)
	fmt.Fprintf(w, "%s", act._info)
}

func ActErrNotfound(w http.ResponseWriter, r *http.Request) {
	err := ActionError{
		_code: http.StatusNotFound,
		_info: "Not Found",
	}

	Log.Log().Errorf("Action not found error. request: %s", DumpAction(r))
	err.Act(w, r)
}

func ActErrInternalError(w http.ResponseWriter, r *http.Request) {
	err := ActionError{
		_code: http.StatusInternalServerError,
		_info: "Server Internal Error",
	}

	Log.Log().Errorf("Action server internal error. request: %s", DumpAction(r))
	err.Act(w, r)
}

func ActErrInvalidPushRequest(w http.ResponseWriter, r *http.Request, code int, info string) {
	body, _ := Json.NewPushResq(strconv.Itoa(code), info)

	err := ActionError{
		_code: http.StatusForbidden,
		_info: (*body).ToString(),
	}

	Log.Log().Errorf("Action invalid push request error. request: %s", DumpAction(r))
	err.ActJson(w, r)
}

func ActErrOther(w http.ResponseWriter, r *http.Request, code int, info string) {
	err := ActionError{
		_code: code,
		_info: info,
	}

	Log.Log().Errorf("Action other error. code: %d, info: %s, request: %s", code, info, DumpAction(r))
	err.Act(w, r)
}
