package middle

import (
	"github.com/justinas/nosurf"
	"github.com/pressly/chi/middleware"
	"net/http"
)

func GetBaseCtx(r *http.Request) map[string]interface{} {
	val := map[string]interface{}{
		"Config":    GetConfig(r),
		"ReqID":     middleware.GetReqID(r.Context()),
		"CSRFToken": nosurf.Token(r),
	}

	return val
}

func CSRFProtect(next http.Handler) http.Handler {
	return nosurf.New(next)
}
