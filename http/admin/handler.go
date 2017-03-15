package admin

import (
	"bytes"
	"github.com/GeertJohan/go.rice"
	"github.com/pressly/chi"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"html/template"
	"net/http"
	"time"
)

var (
	panelTmpl  = template.New("admin/panel")
	loginTmpl  = template.New("admin/login")
	statusTmpl = template.New("admin/status")
)

func LoadTemplates() error {
	var err error
	box, err := rice.FindBox("res/")
	if err != nil {
		panic(err)
	}
	panelTmpl, err = panelTmpl.Parse(box.MustString("panel.html"))
	if err != nil {
		panic(err)
	}
	loginTmpl, err = loginTmpl.Parse(box.MustString("index.html"))
	if err != nil {
		panic(err)
	}
	statusTmpl, err = statusTmpl.Parse(box.MustString("status.html"))
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
	r.Post("/cleanup.sh", handleCleanup)
	r.Post("/login.sh", handleLogin)
	r.Post("/logout.sh", handleLogout)
	r.Post("/new_admin.sh", handleNewAdmin)
	r.Post("/del_admin.sh", handleDelAdmin)
	r.Get("/status.sh", serveStatus)
	r.Post("/set_rules.sh", handleSetRules)
}

// Router sets up moderation functions
// It **must** be setup on the /mod/ basepath
func ModRouter(r chi.Router) {
	r.Post("/del_reply.sh", handleDelPost)
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
