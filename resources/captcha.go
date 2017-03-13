package resources

import (
	"github.com/dchest/captcha"
	"net/http"
)

func MakeCaptcha() string {
	return captcha.NewLen(5)
}

func VerifyCaptcha(r *http.Request) bool {
	return captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution"))
}

var ServeCaptcha = captcha.Server(captcha.StdWidth, captcha.StdHeight)
