package board

import (
	"bytes"
	"errors"
	"github.com/GeertJohan/go.rice"
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

var box = riceConf.MustFindBox("http/board/res/")

var (
	tmpls = template.New("base")
	//dirTmpl    = template.New("board/dir")
	//boardTmpl  = template.New("board/board")
	//threadTmpl = template.New("board/thread")

	hdlFMap = template.FuncMap{
		"renderText": resources.OperateReplyText,
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"rateSpam":    resources.SpamScore,
		"makeCaptcha": resources.MakeCaptcha,
		"dateFromID":  resources.DateFromId,
		"formatDate": func(date time.Time) string {
			return date.Format("02 Jan 06 15:04:05")
		},
	}
)

func init() {
	var err error
	tmpls = tmpls.Funcs(hdlFMap)
	tmpls, err = tmpls.New("thread/postlists").Parse(box.MustString("thread.tmpl.html"))
	if err != nil {
		panic(err)
	}
	_, err = tmpls.New("board/dir").Parse(box.MustString("dir.html"))
	if err != nil {
		panic(err)
	}
	_, err = tmpls.New("board/board").Parse(box.MustString("board.html"))
	if err != nil {
		panic(err)
	}
	_, err = tmpls.New("board/thread").Parse(box.MustString("thread.html"))
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
	r.Get("/:board/:thread/:reply/:unused.png", serveFullImage)
	r.Get("/:board/:thread/:reply/thumb.png", serveThumb)
	r.Post("/:board/:thread/reply.sh", handleNewReply)
	r.Handle("/captcha/:captchaId.png", resources.ServeCaptcha)
	r.Handle("/captcha/:captchaId.wav", resources.ServeCaptcha)
	r.Handle("/captcha/download/:captchaId.wav", resources.ServeCaptcha)
}

func serveThumb(w http.ResponseWriter, r *http.Request) {
	dat := bytes.NewBuffer([]byte{})
	db := middle.GetDB(r)
	err := db.View(func(tx *buntdb.Tx) error {
		bName := chi.URLParam(r, "board")
		tid, err := strconv.Atoi(chi.URLParam(r, "thread"))
		if err != nil {
			return err
		}
		rid, err := strconv.Atoi(chi.URLParam(r, "reply"))
		if err != nil {
			return err
		}

		reply, err := resources.GetReply(tx, r.Host, bName, tid, rid)
		if err != nil {
			return err
		}
		_, err = dat.Write(reply.Thumbnail)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	http.ServeContent(w, r, "thumb.png", time.Now(), bytes.NewReader(dat.Bytes()))
}

func serveFullImage(w http.ResponseWriter, r *http.Request) {
	dat := bytes.NewBuffer([]byte{})
	db := middle.GetDB(r)
	err := db.View(func(tx *buntdb.Tx) error {
		bName := chi.URLParam(r, "board")
		tid, err := strconv.Atoi(chi.URLParam(r, "thread"))
		if err != nil {
			return err
		}
		rid, err := strconv.Atoi(chi.URLParam(r, "reply"))
		if err != nil {
			return err
		}

		reply, err := resources.GetReply(tx, r.Host, bName, tid, rid)
		if err != nil {
			return err
		}
		_, err = dat.Write(reply.Image)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	http.ServeContent(w, r, "image.png", time.Now(), bytes.NewReader(dat.Bytes()))
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
	err = tmpls.ExecuteTemplate(dat, "board/dir", ctx)
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	http.ServeContent(w, r, "dir.html", time.Now(), bytes.NewReader(dat.Bytes()))
}
