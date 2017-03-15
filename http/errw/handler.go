package errw

import (
	"errors"
	"github.com/GeertJohan/go.rice"
	"github.com/pressly/chi/middleware"
	"go.rls.moe/nyx/http/middle"
	"html/template"
	"net/http"
)

var (
	errorTmpl = template.New("errw/error")
)

func LoadTemplates() error {
	box, err := rice.FindBox("errw_res/")
	if err != nil {
		return err
	}
	errorTmpl, err = errorTmpl.Parse(box.MustString("error.html"))
	if err != nil {
		return err
	}
	return nil
}

type ErrorWithTitle interface {
	error
	ErrorTitle() string
}

type errorWTInt struct {
	message, title string
}

func (e errorWTInt) Error() string {
	return e.message
}

func (e errorWTInt) ErrorTitle() string {
	return e.title
}

func MakeErrorWithTitle(title, message string) ErrorWithTitle {
	return errorWTInt{message, title}
}

func ErrorWriter(err error, w http.ResponseWriter, r *http.Request) {
	ctx := middle.GetBaseCtx(r)

	if err == nil {
		ErrorWriter(errors.New("Unknonw Error"), w, r)
	}

	if errWT, ok := err.(ErrorWithTitle); ok {
		ctx["Error"] = map[string]string{
			"Code":        middleware.GetReqID(r.Context()),
			"Description": errWT.Error(),
			"Title":       errWT.ErrorTitle(),
		}
	} else {
		ctx["Error"] = map[string]string{
			"Code":        middleware.GetReqID(r.Context()),
			"Description": err.Error(),
			"Title":       "Error",
		}
	}
	errorTmpl.Execute(w, ctx)
	return
}
