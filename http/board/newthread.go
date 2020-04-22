package board

import (
	"fmt"
	"net/http"

	"github.com/pressly/chi"
	"github.com/tidwall/buntdb"
	"go.rls.moe/nyx/config"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"go.rls.moe/nyx/resources"
)

func handleNewThread(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	err = r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	if middle.GetConfig(r).Captcha.Mode != config.CaptchaDisabled {
		if !resources.VerifyCaptcha(r) {
			http.Redirect(w, r,
				fmt.Sprintf("/%s/board.html?err=wrong_captcha",
					chi.URLParam(r, "board")),
				http.StatusSeeOther)
			return
		}
	}

	var thread = &resources.Thread{}
	var mainReply = &resources.Reply{}

	thread.Board = chi.URLParam(r, "board")
	thread.Metadata = map[string]string{}

	err = parseReply(r, mainReply)
	if err == trollThrottle {
		http.Redirect(w, r,
			fmt.Sprintf("/%s/board.html?err=trollthrottle",
				chi.URLParam(r, "board")),
			http.StatusSeeOther)
		return
	} else if err != nil {
		errw.ErrorWriter(err, w, r)
		return
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
