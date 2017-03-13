package board

import (
	"bytes"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/pressly/chi"
	"github.com/tidwall/buntdb"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"go.rls.moe/nyx/resources"
	"image"
	"image/png"
	"net/http"
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

	if !resources.VerifyCaptcha(r) {
		http.Redirect(w, r,
			fmt.Sprintf("/%s/board.html?err=wrong_captcha",
				chi.URLParam(r, "board")),
			http.StatusSeeOther)
		return
	}

	var thread = &resources.Thread{}
	var mainReply = &resources.Reply{}

	mainReply.Board = chi.URLParam(r, "board")
	thread.Board = chi.URLParam(r, "board")
	mainReply.Text = r.FormValue("text")
	if len(mainReply.Text) > 10000 {
		errw.ErrorWriter(errw.MakeErrorWithTitle("I'm sorry but I can't do that", "These are too many characters"), w, r)
		return
	}
	if len(mainReply.Text) < 5 {
		errw.ErrorWriter(errw.MakeErrorWithTitle("I'm sorry but I can't do that", "These are not enough characters"), w, r)
		return
	}

	if score, err := resources.SpamScore(mainReply.Text); err != nil || !resources.CaptchaPass(score) {
		http.Redirect(w, r,
			fmt.Sprintf("/%s/board.html?err=trollthrottle",
				chi.URLParam(r, "board")),
			http.StatusSeeOther)
		return
	}

	{
		file, _, err := r.FormFile("image")
		if err != nil && err != http.ErrMissingFile {
			errw.ErrorWriter(err, w, r)
			return
		} else if err != http.ErrMissingFile {
			img, _, err := image.Decode(file)
			if err != nil {
				errw.ErrorWriter(err, w, r)
				return
			}
			thumb := resize.Thumbnail(128, 128, img, resize.Lanczos3)
			imgBuf := bytes.NewBuffer([]byte{})
			err = png.Encode(imgBuf, img)
			if err != nil {
				errw.ErrorWriter(err, w, r)
				return
			}
			fmt.Println("Image has size ", len(imgBuf.Bytes()))
			mainReply.Image = imgBuf.Bytes()
			imgBuf = bytes.NewBuffer([]byte{})
			err = png.Encode(imgBuf, thumb)
			if err != nil {
				errw.ErrorWriter(err, w, r)
				return
			}
			mainReply.Thumbnail = imgBuf.Bytes()
		}
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
