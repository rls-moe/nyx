package board

import (
	"github.com/GeertJohan/go.rice/embedded"
	"time"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    `board.html`,
		FileModTime: time.Unix(1489412682, 0),
		Content:     string("<!doctype html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"utf-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>{{.Config.Site.Title}} - /{{.Board.ShortName}}/</title>\n    <link rel=\"stylesheet\" href=\"/@/style.css\">\n    <link rel=\"stylesheet\" href=\"/@/custom.css\">\n</head>\n<body>\n<div class=\"banner logo\">\n    <div class=\"site title\"><h1><span class=\"reflink\"><a href=\"/{{.Board.ShortName}}/board.html\">/{{.Board.ShortName}}/</a></span></h1></div>\n    <div class=\"site description\"><h2>{{.Board.LongName}}</h2></div>\n</div>\n{{ $boardlink := .Board.ShortName }}\n{{ if .Session }}\n{{ if eq (.Session.CAttr \"mode\") \"admin\" }}\nLogged in as Admin\n{{ end }}\n{{ if eq (.Session.CAttr \"mode\") \"mod\" }}\nLogged in as Mod for {{ .Session.CAttr \"board\" }}\n{{ end }}\n{{ end }}\n<hr />\n{{ template \"thread/post\" . }}\n<div class=\"postlists\">\n    {{ $board := .Board }}\n    {{ $csrf := .CSRFToken }}\n    {{ $session := .Session }}\n    {{range .Threads}}\n        {{ template \"thread/postlists\" dict \"Thread\" . \"Board\" $board \"CSRFToken\" $csrf \"Session\" $session }}\n    {{end}}\n</div>\n</body>\n</html>"),
	}
	file3 := &embedded.EmbeddedFile{
		Filename:    `dir.html`,
		FileModTime: time.Unix(1489315156, 0),
		Content:     string("<!doctype html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"utf-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>{{.Config.Site.Title}} Boards</title>\n    <link rel=\"stylesheet\" href=\"/@/style.css\">\n    <link rel=\"stylesheet\" href=\"/@/custom.css\">\n</head>\n<body>\n    <div class=\"banner logo\">\n        <div class=\"site title\"><h1>{{.Config.Site.Title}}</h1></div>\n        <div class=\"site description\"><h2>{{.Config.Site.Description}}</h2></div>\n    </div>\n    <div class=\"boardlist\">\n        <div class=\"boardtitle\">\n            <h3>Boards</h3>\n        </div>\n        <div class=\"boardlist\">\n            <ul>\n                {{range .Boards}}\n                    <li>\n                        <a class=\"boardlink\" href=\"/{{ .ShortName}}/board.html\">{{.ShortName}}: {{.LongName}}</a>\n                    </li>\n                {{end}}\n            </ul>\n        </div>\n    </div>\n</body>\n</html>"),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    `thread.html`,
		FileModTime: time.Unix(1489412660, 0),
		Content:     string("<!doctype html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"utf-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>{{.Config.Site.Title}} - /{{.Board.ShortName}}/</title>\n    <link rel=\"stylesheet\" href=\"/@/style.css\">\n    <link rel=\"stylesheet\" href=\"/@/custom.css\">\n</head>\n<body>\n<div class=\"banner logo\">\n    <div class=\"site title\"><h1><span class=\"reflink\"><a href=\"/{{.Board.ShortName}}/board.html\">/{{.Board.ShortName}}/</a></span></h1></div>\n    <div class=\"site description\"><h2>{{.Board.LongName}}</h2></div>\n    <div class=\"site thread\"><h3>{{.Thread.ID}}</h3></div>\n</div>\n{{ $boardlink := .Board.ShortName }}\n<hr />\n{{ template \"thread/post\" . }}\n{{ template \"thread/postlists\" . }}\n</body>\n</html>"),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    `thread.tmpl.html`,
		FileModTime: time.Unix(1489566591, 0),
		Content:     string("{{ define \"thread/post\" }}\n<div class=\"postarea\">\n    {{ if .Thread }}\n    <form id=\"postform\"\n          action=\"/{{.Board.ShortName}}/{{.Thread.ID}}/reply.sh\"\n          method=\"POST\"\n          enctype=\"multipart/form-data\">\n    {{ else }}\n    <form id=\"postform\"\n          action=\"/{{.Board.ShortName}}/new_thread.sh\"\n          method=\"POST\"\n          enctype=\"multipart/form-data\">\n    {{ end }}\n        <table>\n            <tbody>\n            {{ if .PreviousError }}\n            <tr>\n                <td class=\"postblock\">\n                    Error\n                </td>\n                <td>\n                    {{.PreviousError}}\n                </td>\n            </tr>\n            {{ end }}\n            <tr>\n                <td class=\"postblock\">\n                    TripCode\n                </td>\n                <td>\n                    <input type=\"text\" name=\"tripcode\" size=48 placeholder=\"Anonymous\"/>\n                    <input\n                            type=\"hidden\"\n                            name=\"csrf_token\"\n                            value=\"{{ .CSRFToken }}\" />\n                </td>\n            </tr>\n            <tr>\n                <td class=\"postblock\">\n                    Comment\n                </td>\n                <td>\n                        <textarea\n                                name=\"text\"\n                                placeholder=\"your comment\"\n                                rows=\"4\"\n                                cols=\"48\"\n                                minlength=\"5\"\n                                required\n                        ></textarea>\n                </td>\n            </tr>\n            <tr>\n                <td class=\"postblock\">\n                    Image File\n                </td>\n                <td>\n                    <input type=\"file\" name=\"image\" />\n                </td>\n            </tr>\n            {{ if ne .Config.Captcha.Mode \"disabled\" }}\n            <tr>\n                <td class=\"postblock\">\n                    Captcha\n                </td>\n                <td>\n                    {{ $captchaId := makeCaptcha }}\n                    <img id=\"image\" src=\"/captcha/{{$captchaId}}.png\" alt=\"Captcha Image\"/>\n                    <audio id=audio controls style=\"display:none\" src=\"/captcha/{{$captchaId}}.wav\" preload=none>\n                        You browser doesn't support audio.\n                        <a href=\"/captcha/download/{{$captchaId}}.wav\">Download file</a> to play it in the external player.\n                    </audio>\n                    <br>\n                    <input type=\"text\" name=\"captchaSolution\" size=48 />\n                    <input type=\"hidden\"\n                           name=\"captchaId\"\n                           value=\"{{$captchaId}}\"/>\n                </td>\n            </tr>\n            {{ end }}\n            {{ if (isModSession .Session) }}\n            <tr>\n                <td class=\"postblock\">\n                    Mod Post\n                </td>\n                <td>\n                    <label>\n                        <input type=\"checkbox\" name=\"modpost\"/>Mark as Mod Post\n                    </label>\n                    {{ if (isAdminSession .Session) }}\n                    <label>\n                        <input type=\"checkbox\" name=\"adminpost\"/>Mark as Admin Post\n                    </label>\n                    {{ end }}\n                </td>\n            </tr>\n            {{ end }}\n            <tr>\n                <td class=\"postblock\">\n\n                </td>\n                <td>\n                    <input type=\"submit\" value=\"Post\" />\n                </td>\n            </tr>\n            {{ if .Board.Metadata.rules }}\n            <tr>\n                <td class=\"postblock\">\n                    Rules\n                </td>\n                <td class=\"rules\">\n                    {{ renderText .Board.Metadata.rules }}\n                </td>\n            </tr>\n            {{ end }}\n            </tbody>\n        </table>\n    </form>\n</div>\n<hr />\n{{ end }}\n\n{{ define \"thread/reply\" }}\n    <label><span class=\"postertrip\">\n        {{ if .Reply.Metadata.trip }}\n            {{ .Reply.Metadata.trip}}\n        {{ else }}\n            Anonymous\n        {{ end }}\n        {{ if .Reply.Metadata.modpost }}\n            (Mod)\n        {{ end }}\n        {{ if .Reply.Metadata.adminpost }}\n            [Admin]\n        {{ end }}\n    </span></label>\n    <span class=\"date\">{{dateFromID .Reply.ID | formatDate}}</span>\n    {{ if .Session }}\n    {{ if eq (.Session.CAttr \"mode\") \"admin\" }}\n    <form class=\"delform\" action=\"/mod/del_reply.sh\" method=\"POST\">\n        <input\n                type=\"hidden\"\n                name=\"csrf_token\"\n                value=\"{{ .CSRF }}\" />\n        <input\n                type=\"hidden\"\n                name=\"reply_id\"\n                value=\"{{ .Reply.ID }}\" />\n        <input\n                type=\"hidden\"\n                name=\"thread_id\"\n                value=\"{{ .ThreadID }}\" />\n        <input\n                type=\"hidden\"\n                name=\"board\"\n                value=\"{{ .Boardlink }}\" />\n        <input type=\"submit\" value=\"delete\" />\n    </form>\n    {{ end }}\n    {{ end }}\n    <span>\n        {{ if not .Reply.Metadata.spamscore }}\n        {{ $score := (rateSpam .Reply.Text) }}\n            {{printf \"[SpamScore: %f]\" $score }}\n            {{printf \"[Captcha: %.3f%%]\" (percentFloat (captchaProb $score)) }}\n            {{printf \"[OLD]\"}}\n        {{ else }}\n            {{ printf \"[SpamScore: %s]\" .Reply.Metadata.spamscore }}\n            {{ printf \"[Captcha: %s %%]\" .Reply.Metadata.captchaprob }}\n        {{ end }}\n    </span>\n    <span class=\"reflink\">\n        <a href=\"/{{.Boardlink}}/{{.ThreadID}}/thread.html\">No.{{.Reply.ID}}</a>\n    </span>\n    {{ if .Reply.Thumbnail }}\n    <br />\n    <a target=\"_blank\" href=\"/{{.Boardlink}}/{{.ThreadID}}/{{.Reply.ID}}/{{.Reply.ID}}.png\">\n    <img\n            src=\"/{{.Boardlink}}/{{.ThreadID}}/{{.Reply.ID}}/thumb.png\"\n            class=\"thumb\"\n    />\n    </a>\n    {{ end }}\n    {{ if .Reply.Metadata.deleted }}\n    <blockquote><blockquote class=\"deleted\">\n        {{ renderText .Reply.Text }}\n    </blockquote></blockquote>\n    {{ else }}\n    <blockquote><blockquote>\n        {{ renderText .Reply.Text}}\n    </blockquote></blockquote>\n    {{ end }}\n{{ end }}\n\n{{ define \"thread/main\" }}\n<div class=\"postlists\">\n    {{ $boardlink := .Board.ShortName }}\n    {{ $threadrid := .Thread.GetReply.ID }}\n    {{ $threadid := .Thread.ID }}\n    {{ $csrf := .CSRFToken }}\n    {{ $session := .Session }}\n    {{ with .Thread }}\n        {{ with .GetReply }}\n        {{ with dict \"Reply\" . \"Boardlink\" $boardlink \"CSRF\" $csrf \"ThreadID\" $threadid \"Session\" $session }}\n            {{ template \"thread/reply\" . }}\n        {{ end }}\n        {{ end }}\n    {{range .GetReplies}}\n    {{ if ne .ID $threadrid }}\n    <table class=\"reply-table\"><tbody><tr><td class=\"doubledash\">&gt;&gt;</td>\n        <td class=\"reply\" id=\"reply{{.ID}}\">\n            {{ with dict \"Reply\" . \"Boardlink\" $boardlink \"CSRF\" $csrf \"ThreadID\" $threadid \"Session\" $session }}\n            {{ template \"thread/reply\" . }}\n            {{ end }}\n        </td>\n    </tr></tbody></table>\n    {{end}}\n    {{end}}\n    {{end}}\n    <br clear=\"left\" /><hr />\n</div>\n{{ end }}\n\n{{ template \"thread/main\" . }}"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   ``,
		DirModTime: time.Unix(1489566591, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // board.html
			file3, // dir.html
			file4, // thread.html
			file5, // thread.tmpl.html

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`res/`, &embedded.EmbeddedBox{
		Name: `res/`,
		Time: time.Unix(1489566591, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"board.html":       file2,
			"dir.html":         file3,
			"thread.html":      file4,
			"thread.tmpl.html": file5,
		},
	})
}
