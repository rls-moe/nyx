package admin

import (
	"github.com/icza/session"
	"github.com/tidwall/buntdb"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"go.rls.moe/nyx/resources"
	"net/http"
)

func handleLogout(w http.ResponseWriter, r *http.Request) {
	sess := middle.GetSession(r)
	if sess == nil {
		http.Redirect(w, r, "/admin/index.html", http.StatusSeeOther)
	}
	session.Remove(sess, w)
	http.Redirect(w, r, "/admin/index.html", http.StatusSeeOther)
}
func handleLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		errw.ErrorWriter(err, w, r)
	}
	db := middle.GetDB(r)

	var admin = &resources.AdminPass{}
	err = db.View(func(tx *buntdb.Tx) error {
		var err error
		admin, err = resources.GetAdmin(tx, r.FormValue("id"))
		return err
	})
	if err != nil {
		err = errw.MakeErrorWithTitle("Access Denied", "User or Password Invalid")
		errw.ErrorWriter(err, w, r)
		return
	}
	err = admin.VerifyLogin(r.FormValue("pass"))
	if err != nil {
		err = errw.MakeErrorWithTitle("Access Denied", "User or Password Invalid")
		errw.ErrorWriter(err, w, r)
		return
	}
	sess := session.NewSessionOptions(&session.SessOptions{
		CAttrs: map[string]interface{}{"mode": "admin"},
	})
	session.Add(sess, w)

	http.Redirect(w, r, "/admin/panel.html", http.StatusSeeOther)
}