package http

import (
	"github.com/GeertJohan/go.rice/embedded"
	"time"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    `admin.css`,
		FileModTime: time.Unix(1489250860, 0),
		Content:     string("/* CUSTOM CSS */\ndiv.admin.login {\n    border: 1px solid black;\n    width: 500px;\n    margin: auto;\n    margin-top: 100px;\n}\n.admin.form.row {\n    margin: auto;\n    padding: 5px;\n    width: 90%;\n    height: 22px;\n    left: 0;\n    right: 0;\n    display: flex;\n}\n.admin.form.input {\n    font-family: \"monospace\";\n    width: 100%;\n    height: 100%;\n    padding: 2px;\n    display: inline;\n}\n.admin.form.input.halfsize {\n    width: 50%;\n}"),
	}
	file3 := &embedded.EmbeddedFile{
		Filename:    `custom.css`,
		FileModTime: time.Unix(1489426703, 0),
		Content:     string("h1 {\n    font-size: 32px;\n}\n\nh2 {\n    font-size: 24px;\n}\n\nh3 {\n    font-size: 16px;\n}\n\ndiv {\n    display: block;\n    margin: 0;\n    padding: 0;\n}\n\nblockquote blockquote {\n    word-wrap: break-word;\n    word-break: break-all;\n    white-space: normal;\n    padding: 2px;\n    margin-bottom: 1em;\n    margin-top: 1em;\n    margin-left: 40px;\n    margin-right: 40px;\n}\n\n\n.delform {\n    display: inline;\n    margin: 0;\n    padding: 0;\n}\n.delform input {\n    display: inline;\n}\n\n.deleted {\n    color: #707070;\n}\n\n.reply-table {\n    display: block;\n}\n\n.reply {\n    display: table;\n}"),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    `style.css`,
		FileModTime: time.Unix(1489426859, 0),
		Content:     string("/* The following CSS is mostly taken from Wakaba, big thanks for the devs there! <3 */\n\nhtml, body {\n    background:#FFFFEE;\n    color:#800000;\n}\na {\n    color:#0000EE;\n}\na:hover {\n    color:#DD0000;\n}\n.adminbar {\n    text-align:right;\n    clear:both;\n    float:right;\n}\n.logo {\n    clear:both;\n    text-align:center;\n    font-size:2em;\n    color:#800000;\n    width:100%;\n}\n.theader {\n    background:#E04000;\n    text-align:center;\n    padding:2px;\n    color:#FFFFFF;\n    width:100%;\n}\n.postarea {\n}\n.rules {\n    font-size:0.7em;\n}\n.postblock {\n    background:#EEAA88;\n    color:#800000;\n    font-weight:800;\n}\n.footer {\n    text-align:center;\n    font-size:12px;\n    font-family:serif;\n}\n.passvalid {\n    background:#EEAA88;\n    text-align:center;\n    width:100%;\n    color:#ffffff;\n}\n.dellist {\n    font-weight: bold;\n    text-align:center;\n}\n.delbuttons {\n    text-align:center;\n    padding-bottom:4px;\n\n}\n.managehead {\n    background:#AAAA66;\n    color:#400000;\n    padding:0px;\n}\n.postlists {\n    background:#FFFFFF;\n    width:100%;\n    padding:0px;\n    color:#800000;\n}\n.row1 {\n    background:#EEEECC;\n    color:#800000;\n}\n.row2 {\n    background:#DDDDAA;\n    color:#800000;\n}\n.unkfunc {\n    background:inert;\n    color:#789922;\n}\n.filesize {\n    text-decoration:none;\n}\n.filetitle {\n    background:inherit;\n    font-size:1.2em;\n    color:#CC1105;\n    font-weight:800;\n}\n.postername {\n    color:#117743;\n    font-weight:bold;\n}\n.postertrip {\n    color:#228854;\n}\n.oldpost {\n    color:#CC1105;\n    font-weight:800;\n}\n.omittedposts {\n    color:#707070;\n}\n.reply {\n    background:#F0E0D6;\n    color:#800000;\n}\n.doubledash {\n    vertical-align:top;\n    clear:both;\n    float:left;\n}\n.replytitle {\n    font-size: 1.2em;\n    color:#CC1105;\n    font-weight:800;\n}\n.commentpostername {\n    color:#117743;\n    font-weight:800;\n}\n.thumbnailmsg {\n    font-size: small;\n    color:#800000;\n}\n\n\n\n.abbrev {\n    color:#707070;\n}\n.highlight {\n    background:#F0E0D6;\n    color:#800000;\n    border: 2px dashed #EEAA88;\n}\n\n/* From pl files */\n\n/* futaba_style.pl */\n\nform { margin-bottom: 0px }\nform .trap { display:none }\n.postarea { text-align: center }\n.postarea table { margin: 0px auto; text-align: left }\n.thumb { border: none; float: left; margin: 2px 20px }\n.nothumb { float: left; background: #eee; border: 2px dashed #aaa; text-align: center; margin: 2px 20px; padding: 1em 0.5em 1em 0.5em; }\n\n.reflink a { color: inherit; text-decoration: none }\n.reply .filesize { margin-left: 20px }\n.userdelete { float: right; text-align: center; white-space: nowrap }\n.replypage .replylink { display: none }"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   ``,
		DirModTime: time.Unix(1489426859, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // admin.css
			file3, // custom.css
			file4, // style.css

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`res/`, &embedded.EmbeddedBox{
		Name: `res/`,
		Time: time.Unix(1489426859, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"admin.css":  file2,
			"custom.css": file3,
			"style.css":  file4,
		},
	})
}
