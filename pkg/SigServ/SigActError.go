package SigServ

import (
	"fmt"
	"net/http"
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

	err.Execute(w, r)
}

func ActErrInternalError(w http.ResponseWriter, r *http.Request) {
	err := ActionError{
		_code: http.StatusInternalServerError,
		_info: "Server Internal Error",
	}

	err.Execute(w, r)
}

func ActErrOther(w http.ResponseWriter, r *http.Request, code int, info string) {
	err := ActionError{
		_code: code,
		_info: info,
	}

	err.Execute(w, r)
}
