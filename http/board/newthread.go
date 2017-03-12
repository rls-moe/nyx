package board

import (
	"fmt"
	"github.com/pressly/chi"
	"github.com/tidwall/buntdb"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"go.rls.moe/nyx/resources"
	"net/http"
)

func handleNewThread(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	var thread = &resources.Thread{}
	var mainReply = &resources.Reply{}

	mainReply.Board = chi.URLParam(r, "board")
	thread.Board = chi.URLParam(r, "board")
	mainReply.Text = r.FormValue("text")
	if len(mainReply.Text) > 1000 {
		errw.ErrorWriter(errw.MakeErrorWithTitle("I'm sorry but I can't do that", "These are too many characters"), w, r)
		return
	}
	if len(mainReply.Text) < 10 {
		errw.ErrorWriter(errw.MakeErrorWithTitle("I'm sorry but I can't do that", "These are not enough characters"), w, r)
		return
	}
	mainReply.Metadata = map[string]string{}
	if r.FormValue("tripcode") != "" {
		mainReply.Metadata["trip"] = resources.CalcTripCode(r.FormValue("tripcode"))
	}

	db := middle.GetDB(r)
	if err = db.Update(func(tx *buntdb.Tx) error {
		return resources.NewThread(tx, r.Host, mainReply.Board, thread, mainReply)
	}); err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%d/thread.html", chi.URLParam(r, "board"), thread.ID), http.StatusSeeOther)
}
