package main

import (
	"context"
	"net/http"
)

type Handler interface{}

type Momentum struct {
	mws []Middleware
	*Router
}

type momentumContextKey struct{}

func (m *Momentum) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	m.Router.ServeHTTP(rw, req)
}

func (m *Momentum) UseMiddleware(mw Middleware) {
	m.mws = append(m.mws, mw)
}

func (m *Momentum) createContext(rw http.ResponseWriter, req *http.Request) *Context {
	c := &Context{
		mws:  m.mws,
		Resp: rw,
	}
	c.Req = req.WithContext(context.WithValue(req.Context(), momentumContextKey{}, c))
	return c
}
