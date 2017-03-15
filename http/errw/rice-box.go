package errw

import (
	"github.com/GeertJohan/go.rice/embedded"
	"time"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    `error.html`,
		FileModTime: time.Unix(1489238440, 0),
		Content:     string("<!doctype html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"utf-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>{{.Config.Site.Title}} Admin Login</title>\n    <style>\n        div.error {\n            border: 1px solid black;\n            width: 500px;\n            margin: auto;\n            margin-top: 100px;\n        }\n        div.error h1 {\n            margin-bottom: 0px;\n            text-align: center;\n        }\n        div.error h2 {\n            text-align: center;\n        }\n        div.error h3 {\n            margin-top: 0px;\n            text-align: center;\n            color: #888;\n        }\n    </style>\n</head>\n<body>\n<div class=\"error\">\n    <h1>{{.Error.Title}}</h1><br/>\n    <h3>{{.Error.Code}}</h3><br/>\n    <h2>{{.Error.Description}}</h2>\n</div>\n</body>\n</html>"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   ``,
		DirModTime: time.Unix(1489240168, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // error.html

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`res/`, &embedded.EmbeddedBox{
		Name: `res/`,
		Time: time.Unix(1489240168, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"error.html": file2,
		},
	})
}
