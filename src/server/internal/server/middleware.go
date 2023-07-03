package server

import (
	"log"
	"time"
	"strconv"
	"github.com/valyala/fasthttp"
)

func getHttp(ctx *fasthttp.RequestCtx) string {
	if ctx.Response.Header.IsHTTP11() {
		return "HTTP/1.1"
	}
	return "HTTP/1.0"
}

func logger(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		startTime := time.Now()
		h(ctx)
		duration := time.Since(startTime)

		statusCode := ctx.Response.Header.StatusCode()

		log.Printf("%s[%v]%s %s %s %v | %s %v",
			buildColor(boldBright, red),
			duration,
			buildColor(reset),

			wrapByStatus(statusCode, string(ctx.Method())),
			wrapByStatus(statusCode, "|" + strconv.Itoa(statusCode) + "|"),

			string(ctx.RequestURI()),

			getHttp(ctx),
			ctx.RemoteAddr(),
		)
	}
}

