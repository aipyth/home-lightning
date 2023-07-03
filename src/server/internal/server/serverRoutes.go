package server

import (
	"github.com/valyala/fasthttp"
)

func (s *server) Routes() {
	s.Router.ServeFiles("/static/{filepath:*}", "./static")
	s.Router.GET("/", logger(s.serveHome))
	s.Router.GET("/socket/web-interface", s.serveWebInterface)
}

func (s *server) serveHome(ctx *fasthttp.RequestCtx) {
	fasthttp.ServeFile(ctx, "templates/home.html")
}

