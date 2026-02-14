package goter

import (
	errs "github.com/adagit94/gono/gotils/errors"
	uri "github.com/adagit94/gono/gotils/uri"
	strs "strings"
)

const (
	Post    = "POST"
	Get     = "GET"
	Put     = "PUT"
	Patch   = "PATCH"
	Delete  = "DELETE"
	Options = "OPTIONS"
	Connect = "CONNECT"
	Head    = "HEAD"
	Trace   = "TRACE"
)

type IRouter[H any] interface {
	Route(path string) IRoute[H]
	Select(path string, query string, method string) (H, IParams)
}

func CreateRouter[H any]() IRouter[H] {
	router := &router[H]{tree: make(routesTree[H])}
	return router
}

type segmentConf struct {
	segment string
	static  bool
}

type routeConf[H any] struct {
	segments []segmentConf
	handler  H
}

type routes[H any] map[int][]routeConf[H]

type routesTree[H any] map[string]routes[H]

type router[H any] struct {
	tree routesTree[H]
}

func (router *router[H]) registerHandler(path string, method string, handler H) {
	segs := strs.Split(path, "/")
	segsCount := len(segs)

	if _, methodKeyExists := router.tree[method]; !methodKeyExists {
		router.tree[method] = make(routes[H])
	}

	if _, segsCountKeyExists := router.tree[method][segsCount]; !segsCountKeyExists {
		router.tree[method][segsCount] = make([]routeConf[H], 0)
	}

	router.tree[method][segsCount] = append(router.tree[method][segsCount], routeConf[H]{segments: genSegConfs(segs), handler: handler})
	sortRoutes(router.tree[method][segsCount])
}

func (router *router[H]) Route(path string) IRoute[H] {
	route := &route[H]{path: path, registerHandler: router.registerHandler}
	return route
}

func (router *router[H]) Select(path string, query string, method string) (H, IParams) {
	segs := strs.Split(path, "/")
	segsCount := len(segs)
	segsCountsMap, methodKey := router.tree[method]

	if !methodKey {
		panic(&errs.CodeError{Code: errs.MethodNotRegisteredCode, Message: "Method not registered."})
	}

	routes, segsCountKey := segsCountsMap[segsCount]

	if !segsCountKey {
		panic(&errs.CodeError{Code: errs.RouteNotRegisteredCode, Message: "Route not registered."})
	}

	for _, routeConf := range routes {
		pathParams := make(paramsMap)
		take := true

		for i, seg := range routeConf.segments {
			if seg.static {
				if seg.segment != segs[i] {
					take = false
					break
				}
			} else {
				pathParams[seg.segment] = segs[i]
			}
		}

		if take {
			return routeConf.handler, &params{path: &pathParams, query: uri.ParseQueryStr(query)}
		}
	}

	panic(&errs.CodeError{Code: errs.HandlerNotFoundCode, Message: "Handler not found."})
}

type IRoute[H any] interface {
	Post(H) IRoute[H]
	Get(H) IRoute[H]
	Put(H) IRoute[H]
	Patch(H) IRoute[H]
	Delete(H) IRoute[H]
	Options(H) IRoute[H]
	Connect(H) IRoute[H]
	Head(H) IRoute[H]
	Trace(H) IRoute[H]
}

type route[H any] struct {
	path            string
	registerHandler func(path string, method string, handler H)
}

func (route *route[H]) Post(handler H) IRoute[H] {
	route.registerHandler(route.path, Post, handler)
	return route
}

func (route *route[H]) Get(handler H) IRoute[H] {
	route.registerHandler(route.path, Get, handler)
	return route
}

func (route *route[H]) Put(handler H) IRoute[H] {
	route.registerHandler(route.path, Put, handler)
	return route
}

func (route *route[H]) Patch(handler H) IRoute[H] {
	route.registerHandler(route.path, Patch, handler)
	return route
}

func (route *route[H]) Delete(handler H) IRoute[H] {
	route.registerHandler(route.path, Delete, handler)
	return route
}

func (route *route[H]) Options(handler H) IRoute[H] {
	route.registerHandler(route.path, Options, handler)
	return route
}

func (route *route[H]) Connect(handler H) IRoute[H] {
	route.registerHandler(route.path, Connect, handler)
	return route
}

func (route *route[H]) Head(handler H) IRoute[H] {
	route.registerHandler(route.path, Head, handler)
	return route
}

func (route *route[H]) Trace(handler H) IRoute[H] {
	route.registerHandler(route.path, Trace, handler)
	return route
}

type paramsMap = map[string]string

type IParams interface {
	Path(param string) string
	Query(param string) string
}

type params struct {
	path  *paramsMap
	query *paramsMap
}

func (params *params) Path(param string) string {
	return (*params.path)[param]
}

func (params *params) Query(param string) string {
	return (*params.query)[param]
}
