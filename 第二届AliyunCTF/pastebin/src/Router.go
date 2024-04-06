package main

import (
	"net/http"
)

type Router struct {
	routes []RouteEntry
	m      *Momentum
}

type RouteEntry struct {
	Path        string
	Method      string
	HandlerFunc http.HandlerFunc
}
type Key struct{}

type reqContextKey = Key

func (ent *RouteEntry) Match(method string, pattern string) bool {
	if method != ent.Method {
		return false // Method mismatch
	}
	if pattern != ent.Path {
		return false // Path mismatch
	}
	return true
}

func (rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, e := range rtr.routes {
		match := e.Match(r.Method, r.URL.Path)
		if !match {
			continue
		}

		e.HandlerFunc.ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}

func (rtr *Router) Handle(method string, pattern string, handlers []Handler) {
	rtr.handle(method, pattern, func(resp http.ResponseWriter, req *http.Request) {
		c := rtr.m.createContext(resp, req)
		for _, h := range handlers {
			c.mws = append(c.mws, getMWFromHandler(h))
		}
		c.run()
	})
}

func (rtr *Router) handle(method, pattern string, handle http.HandlerFunc) {
	e := RouteEntry{
		Method:      method,
		Path:        pattern,
		HandlerFunc: handle,
	}
	rtr.routes = append(rtr.routes, e)
}

func (rtr *Router) Get(pattern string, h ...Handler)  { rtr.Handle("GET", pattern, h) }
func (rtr *Router) Post(pattern string, h ...Handler) { rtr.Handle("POST", pattern, h) }
