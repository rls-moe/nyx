package admin

import (
	"errors"
	"github.com/tidwall/buntdb"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"go.rls.moe/nyx/resources"
	"net/http"
)

func handleNewBoard(w http.ResponseWriter, r *http.Request) {
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

	var board = &resources.Board{}

	board.ShortName = r.FormValue("shortname")
	board.LongName = r.FormValue("longname")

	if board.ShortName == "" {
		errw.ErrorWriter(errors.New("Need shortname"), w, r)
		return
	}

	if board.ShortName == "admin" || board.ShortName == "@" || board.ShortName == "mod"{
		errw.ErrorWriter(errors.New("No"), w, r)
	}

	if board.LongName == "" && len(board.LongName) < 5 {
		errw.ErrorWriter(errors.New("Need 5 characters for long name"), w, r)
		return
	}

	if err = db.Update(func(tx *buntdb.Tx) error {
		return resources.NewBoard(tx, r.Host, board)
	}); err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	http.Redirect(w, r, "/admin/panel.html", http.StatusSeeOther)
}
