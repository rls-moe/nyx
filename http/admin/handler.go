package admin

import (
	"bytes"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/icza/session"
	"github.com/pressly/chi"
	"github.com/tidwall/buntdb"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"go.rls.moe/nyx/resources"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var riceConf = rice.Config{
	LocateOrder: []rice.LocateMethod{
		rice.LocateWorkingDirectory,
		rice.LocateEmbedded,
		rice.LocateAppended,
	},
}

var box = riceConf.MustFindBox("http/admin/res/")

var (
	panelTmpl = template.New("admin/panel")
	loginTmpl = template.New("admin/login")
)

func init() {
	var err error
	panelTmpl, err = panelTmpl.Parse(box.MustString("panel.html"))
	if err != nil {
		panic(err)
	}
	loginTmpl, err = loginTmpl.Parse(box.MustString("index.html"))
	if err != nil {
		panic(err)
	}
}

// Router sets up the Administration Panel
// It **must** be setup on the /admin/ basepath
func AdminRouter(r chi.Router) {
	r.Get("/", serveLogin)
	r.Get("/index.html", serveLogin)
	r.Get("/panel.html", servePanel)
	r.Post("/new_board.sh", handleNewBoard)
	r.Post("/login.sh", handleLogin)
	r.Post("/logout.sh", handleLogout)
}

// Router sets up moderation functions
// It **must** be setup on the /mod/ basepath
func ModRouter(r chi.Router) {
	r.Post("/del_reply.sh", handleDelPost)
}

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

func serveLogin(w http.ResponseWriter, r *http.Request) {
	dat := bytes.NewBuffer([]byte{})
	err := loginTmpl.Execute(dat, middle.GetBaseCtx(r))
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	http.ServeContent(w, r, "index.html", time.Now(),
		bytes.NewReader(dat.Bytes()))
}

func servePanel(w http.ResponseWriter, r *http.Request) {
	sess := middle.GetSession(r)
	if sess == nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}
	dat := bytes.NewBuffer([]byte{})
	err := panelTmpl.Execute(dat, middle.GetBaseCtx(r))
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	http.ServeContent(w, r, "panel.html", time.Now(),
		bytes.NewReader(dat.Bytes()))
}

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
	}
	err = admin.VerifyLogin(r.FormValue("pass"))
	if err != nil {
		err = errw.MakeErrorWithTitle("Access Denied", "User or Password Invalid")
		errw.ErrorWriter(err, w, r)
	}
	sess := session.NewSessionOptions(&session.SessOptions{
		CAttrs: map[string]interface{}{"mode": "admin"},
	})
	session.Add(sess, w)

	http.Redirect(w, r, "/admin/panel.html", http.StatusSeeOther)
}
