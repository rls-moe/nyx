package admin

import (
	"github.com/tidwall/buntdb"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"go.rls.moe/nyx/resources"
	"net/http"
)

func handleDelAdmin(w http.ResponseWriter, r *http.Request) {
	sess := middle.GetSession(r)
	if !middle.IsAdminSession(sess) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	db := middle.GetDB(r)

	adminID := r.FormValue("adminid")
	if len(adminID) > 255 {
		errw.ErrorWriter(errw.MakeErrorWithTitle("Too long", "The ID of the administrator is too long"), w, r)
		return
	}
	if len(adminID) < 4 {
		errw.ErrorWriter(errw.MakeErrorWithTitle("Too short", "The ID of the administrator is too short"), w, r)
		return
	}

	if err = db.Update(func(tx *buntdb.Tx) error {
		return resources.DelAdmin(tx, adminID)
	}); err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	http.Redirect(w, r, "/admin/panel.html", http.StatusSeeOther)
}

func handleNewAdmin(w http.ResponseWriter, r *http.Request) {
	sess := middle.GetSession(r)
	if !middle.IsAdminSession(sess) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	db := middle.GetDB(r)

	var admin = &resources.AdminPass{}

	admin.ID = r.FormValue("adminid")
	if len(admin.ID) > 255 {
		errw.ErrorWriter(errw.MakeErrorWithTitle("Too long", "The ID of the administrator is too long"), w, r)
		return
	}
	if len(admin.ID) < 4 {
		errw.ErrorWriter(errw.MakeErrorWithTitle("Too short", "The ID of the administrator is too short"), w, r)
		return
	}
	if len(r.FormValue("adminpass")) > 255 {
		errw.ErrorWriter(errw.MakeErrorWithTitle("Too long", "The Password of the administrator is too long"), w, r)
		return
	}
	if len(r.FormValue("adminpass")) < 12 {
		errw.ErrorWriter(errw.MakeErrorWithTitle("Too short", "The Password of the administrator is too short"), w, r)
		return
	}
	admin.HashLogin(r.FormValue("adminpass"))

	if err = db.Update(func(tx *buntdb.Tx) error {
		return resources.NewAdmin(tx, admin)
	}); err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	http.Redirect(w, r, "/admin/panel.html", http.StatusSeeOther)
}
