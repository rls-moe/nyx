package board

import (
	"bytes"
	"github.com/pressly/chi"
	"github.com/tidwall/buntdb"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"go.rls.moe/nyx/resources"
	"log"
	"net/http"
	"time"
)

func serveBoard(w http.ResponseWriter, r *http.Request) {
	dat := bytes.NewBuffer([]byte{})
	db := middle.GetDB(r)
	ctx := middle.GetBaseCtx(r)
	err := db.View(func(tx *buntdb.Tx) error {
		bName := chi.URLParam(r, "board")
		b, err := resources.GetBoard(tx, r.Host, bName)
		if err != nil {
			return err
		}
		ctx["Board"] = b

		threads, err := resources.ListThreads(tx, r.Host, bName)
		if err != nil {
			return err
		}
		log.Println("Number of Thread on board: ", len(threads))

		for k := range threads {
			err := resources.FillReplies(tx, r.Host, threads[k])
			if err != nil {
				return err
			}
		}
		ctx["Threads"] = threads
		return nil
	})
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	err = boardTmpl.Execute(dat, ctx)
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	http.ServeContent(w, r, "board.html", time.Now(), bytes.NewReader(dat.Bytes()))
}
