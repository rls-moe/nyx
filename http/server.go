package http

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"go.rls.moe/nyx/config"
	"go.rls.moe/nyx/http/admin"
	"go.rls.moe/nyx/http/board"
	"go.rls.moe/nyx/http/middle"
	"net/http"
	"time"
)

var riceConf = rice.Config{
	LocateOrder: []rice.LocateMethod{
		rice.LocateWorkingDirectory,
		rice.LocateEmbedded,
		rice.LocateAppended,
	},
}

func Start(config *config.Config) {
	r := chi.NewRouter()

	fmt.Println("Setting up Router")
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CloseNotify)
	r.Use(middle.LimitSize(config))
	r.Use(middleware.DefaultCompress)

	r.Use(middle.ConfigCtx(config))

	r.Use(middle.CSRFProtect)
	{
		mw, err := middle.Database(config)
		if err != nil {
			panic(err)
		}
		r.Use(mw)
	}

	r.Route("/admin/", admin.AdminRouter)
	r.Route("/mod/", admin.ModRouter)
	{
		box := riceConf.MustFindBox("http/res")
		atFileServer := http.StripPrefix("/@/", http.FileServer(box.HTTPBox()))
		r.Mount("/@/", atFileServer)
	}
	r.Group(board.Router)

	fmt.Println("Setup Complete, Starting Web Server...")
	srv := &http.Server{
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		Handler:        r,
		Addr:           config.ListenOn,
		MaxHeaderBytes: 1 * 1024 * 1024,
	}
	srv.ListenAndServe()
}
