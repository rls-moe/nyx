package middle

import (
	"github.com/icza/session"
	"go.rls.moe/nyx/config"
	"net/http"
)

func SetupSessionManager(c *config.Config) {
	session.Global.Close()
	session.Global = session.NewCookieManagerOptions(session.NewInMemStore(),
		&session.CookieMngrOptions{
			AllowHTTP: c.DisableSecurity,
		})
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
