package board

import (
	"bytes"
	"github.com/GeertJohan/go.rice"
	"github.com/pressly/chi"
	"github.com/tidwall/buntdb"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"go.rls.moe/nyx/resources"
	"html/template"
	"net/http"
	"time"
)

var riceConf = rice.Config{
	LocateOrder: []rice.LocateMethod{
		rice.LocateWorkingDirectory,
		rice.LocateEmbedded,
		rice.LocateAppended,
	},
}

var box = riceConf.MustFindBox("http/board/res/")

var (
	dirTmpl    = template.New("board/dir")
	boardTmpl  = template.New("board/board")
	threadTmpl = template.New("board/thread")

	hdlFMap = template.FuncMap{
		"renderText": resources.OperateReplyText,
	}
)

func init() {
	var err error
	dirTmpl, err = dirTmpl.Parse(box.MustString("dir.html"))
	if err != nil {
		panic(err)
	}
	boardTmpl, err = boardTmpl.Funcs(hdlFMap).Parse(box.MustString("board.html"))
	if err != nil {
		panic(err)
	}
	threadTmpl, err = threadTmpl.Funcs(hdlFMap).Parse(box.MustString("thread.html"))
	if err != nil {
		panic(err)
	}
}

func Router(r chi.Router) {
	r.Get("/", serveDir)
	r.Get("/dir.html", serveDir)
	r.Get("/:board/board.html", serveBoard)
	r.Post("/:board/new_thread.sh", handleNewThread)
	r.Get("/:board/:thread/thread.html", serveThread)
	r.Get("/:board/:thread/:post/post.html", servePost)
	r.Post("/:board/:thread/reply.sh", handleNewReply)
}

func servePost(w http.ResponseWriter, r *http.Request) {
	return
}

func serveDir(w http.ResponseWriter, r *http.Request) {
	dat := bytes.NewBuffer([]byte{})
	db := middle.GetDB(r)
	ctx := middle.GetBaseCtx(r)
	err := db.View(func(tx *buntdb.Tx) error {
		bList, err := resources.ListBoards(tx, r.Host)
		if err != nil {
			return err
		}
		ctx["Boards"] = bList
		return nil
	})
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	err = dirTmpl.Execute(dat, ctx)
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	http.ServeContent(w, r, "dir.html", time.Now(), bytes.NewReader(dat.Bytes()))
}
