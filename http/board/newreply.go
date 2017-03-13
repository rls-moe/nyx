package board

import (
	"fmt"
	"github.com/pressly/chi"
	"github.com/tidwall/buntdb"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"go.rls.moe/nyx/resources"
	"net/http"
	"strconv"
)

func handleNewReply(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	if !resources.VerifyCaptcha(r) {
		http.Redirect(w, r,
			fmt.Sprintf("/%s/%s/thread.html?err=wrong_captcha",
				chi.URLParam(r, "board"), chi.URLParam(r, "thread")),
			http.StatusSeeOther)
		return
	}

	var reply = &resources.Reply{}

	reply.Board = chi.URLParam(r, "board")
	tid, err := strconv.Atoi(chi.URLParam(r, "thread"))
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	reply.Thread = tid
	reply.Text = r.FormValue("text")
	if len(reply.Text) > 10000 {
		errw.ErrorWriter(errw.MakeErrorWithTitle("I'm sorry but I can't do that", "These are too many characters"), w, r)
		return
	}
	if len(reply.Text) < 10 {
		errw.ErrorWriter(errw.MakeErrorWithTitle("I'm sorry but I can't do that", "These are not enough characters"), w, r)
		return
	}

	if score, err := resources.SpamScore(reply.Text); err != nil || !resources.CaptchaPass(score) {
		http.Redirect(w, r,
			fmt.Sprintf("/%s/%s/thread.html?err=trollthrottle",
				chi.URLParam(r, "board"), chi.URLParam(r, "thread")),
			http.StatusSeeOther)
		return
	}

	reply.Metadata = map[string]string{}
	if r.FormValue("tripcode") != "" {
		reply.Metadata["trip"] = resources.CalcTripCode(r.FormValue("tripcode"))
	} else {
		reply.Metadata["trip"] = "Anonymous"
	}

	db := middle.GetDB(r)
	if err = db.Update(func(tx *buntdb.Tx) error {
		thread, err := resources.GetThread(tx, r.Host, reply.Board, reply.Thread)
		if err != nil {
			return err
		}
		return resources.NewReply(tx, r.Host, reply.Board, thread, reply, false)
	}); err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%d/thread.html", chi.URLParam(r, "board"), reply.Thread), http.StatusSeeOther)
}
