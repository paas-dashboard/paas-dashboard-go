package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"net/http"
	"path/filepath"
)

type MainController struct {
	web.Controller
}

func (m *MainController) RedirectToIndex() {
	root := filepath.Join(".", "static")
	m.Ctx.ResponseWriter.Header().Set("Content-Type", "text/html")
	http.ServeFile(m.Ctx.ResponseWriter, m.Ctx.Request, filepath.Join(root, "index.html"))
}
