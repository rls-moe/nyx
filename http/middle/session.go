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

func IsAdminSession(sess session.Session) bool {
	if sess == nil {
		return false
	}
	if sess.CAttr("mode") == "admin" {
		return true
	}
	return false
}

func IsModSession(sess session.Session) bool {
	if sess == nil {
		return false
	}
	if IsAdminSession(sess) {
		return true
	}
	if sess.CAttr("mode") == "mod" {
		return true
	}
	return false
}
