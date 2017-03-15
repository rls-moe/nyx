package board

import (
	"errors"
	"fmt"
	"github.com/pressly/chi"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"go.rls.moe/nyx/resources"
	"net/http"
	"strconv"
)

var trollThrottle = errors.New("Troll throttle")

func parseReply(r *http.Request, reply *resources.Reply) error {
	reply.Board = chi.URLParam(r, "board")
	reply.Text = r.FormValue("text")
	if tidStr := chi.URLParam(r, "thread"); tidStr != "" {
		var err error
		reply.Thread, err = strconv.Atoi(tidStr)
		if err != nil {
			return err
		}
	}
	if len(reply.Text) > 10000 {
		return errw.MakeErrorWithTitle(
			"I'm sorry but I can't do that",
			"There are too many characters")
	}
	if len(reply.Text) < 5 {
		return errw.MakeErrorWithTitle(
			"I'm sorry but I can't do that",
			"There are not enough characters")
	}

	reply.Metadata = map[string]string{}

	spamScore, err := resources.SpamScore(reply.Text)
	if err != nil {
		return err
	}

	reply.Metadata["spamscore"] = fmt.Sprintf("%.6f", spamScore)
	reply.Metadata["captchaprob"] = fmt.Sprintf("%.2f", resources.CaptchaProb(spamScore)*100)

	if !resources.CaptchaPass(spamScore) {
		return trollThrottle
	}

	file, hdr, err := r.FormFile("image")
	err = parseImage(reply, file, hdr, err)
	if err != nil {
		return err
	}

	if r.FormValue("tripcode") != "" {
		reply.Metadata["trip"] = resources.CalcTripCode(r.FormValue("tripcode"))
	}
	if middle.IsModSession(middle.GetSession(r)) {
		if r.FormValue("modpost") != "" {
			reply.Metadata["modpost"] = "yes"
		}
		if middle.IsAdminSession(middle.GetSession(r)) {
			if r.FormValue("adminpost") != "" {
				reply.Metadata["adminpost"] = "yes"
			}
		}
	}

	return nil
}
