package admin

import (
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"strconv"
	"go.rls.moe/nyx/resources"
	"fmt"
	"net/http"
	"github.com/tidwall/buntdb"
)

func handleDelPost(w http.ResponseWriter, r *http.Request) {
	sess := middle.GetSession(r)
	if sess == nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}
	if sess.CAttr("mode") != "admin" && sess.CAttr("mode") != "mod" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	rid, err := strconv.Atoi(r.FormValue("reply_id"))
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	trid, err := strconv.Atoi(r.FormValue("thread_id"))
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	board := r.FormValue("board")

	if sess.CAttr("mode") == "mod" && sess.CAttr("board") != board {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Not on this board"))
		return
	}

	db := middle.GetDB(r)

	err = db.Update(func(tx *buntdb.Tx) error {
		reply, err := resources.GetReply(tx, r.Host, board, trid, rid)
		if err != nil {
			return err
		}
		reply.Text = "[deleted]"
		reply.Metadata["deleted"] = "yes"
		reply.Image = nil
		reply.Thumbnail = nil
		err = resources.UpdateReply(tx, r.Host, board, reply)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%d/thread.html", board, trid), http.StatusSeeOther)
}
