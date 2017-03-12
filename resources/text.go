package resources

import (
	"html/template"
	"strings"
)

func OperateReplyText(unsafe string) template.HTML {
	unsafe = template.HTMLEscapeString(unsafe)
	unsafe = strings.Replace(unsafe, "\n", "<br />", -1)
	return template.HTML(unsafe)
}
