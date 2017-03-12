package middle

import (
	"github.com/icza/session"
	"net/http"
)

func init() {
	session.Global.Close()
	session.Global = session.NewCookieManager(session.NewInMemStore())
}

func GetSession(r *http.Request) session.Session {
	return session.Get(r)
}
