package main

import (
	"context"
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

type ReqContext struct {
	*Context
	DBCon *sql.DB
	Nonce string
}

type Middleware = func(next http.Handler) http.Handler

type Context struct {
	mws []Middleware

	Req  *http.Request
	Resp http.ResponseWriter
}

const (
	headerContentType = "Content-Type"
	contentTypeHTML   = "text/html; charset=UTF-8"
	pathPrefix        = "./view/"
)

func (ctx *Context) run() {
	h := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	for i := len(ctx.mws) - 1; i >= 0; i-- {
		h = ctx.mws[i](h)
	}
	h.ServeHTTP(ctx.Resp, ctx.Req)
}

func (ctx *Context) HTML(status int, name string, data interface{}) {
	ctx.Resp.Header().Set(headerContentType, contentTypeHTML)
	ctx.Resp.WriteHeader(status)
	t, _ := template.ParseFiles(pathPrefix + name)
	err := t.ExecuteTemplate(ctx.Resp, name, data)
	if err != nil {
		log.Print(err)
		ctx.Resp.Write([]byte("Internal Server Error"))
	}
}

func GetCtx(c context.Context) *Context {
	if mc, ok := c.Value(momentumContextKey{}).(*Context); ok {
		return mc
	}
	return nil
}

func GetReqCtx(c context.Context) *ReqContext {
	if reqCtx, ok := c.Value(reqContextKey{}).(*ReqContext); ok {
		return reqCtx
	}
	return nil
}
