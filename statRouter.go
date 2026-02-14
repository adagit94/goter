package goter

import (
	e "github.com/adagit94/gono/gotils/errors"
)

type statRoutes[H any] map[string]H
type statTree[H any] map[string]statRoutes[H]

func CreateStaticRouter[H any]() IStatRouter[H] {
	router := &statRouter[H]{tree: make(statTree[H])}
	return router
}

type statRouter[H any] struct {
	tree statTree[H]
}

func (router *statRouter[H]) registerHandler(path string, method string, handler H) {
	_, keyExists := router.tree[method]

	if !keyExists {
		router.tree[method] = make(statRoutes[H])
	}

	router.tree[method][path] = handler
}

type IStatRouter[H any] interface {
	Route(path string) IRoute[H]
	Select(path string, method string) H
}

func (router *statRouter[H]) Route(path string) IRoute[H] {
	route := &route[H]{path: path, registerHandler: router.registerHandler}
	return route
}

func (router *statRouter[H]) Select(path string, method string) H {
	routes, methodKey := router.tree[method]

	if !methodKey {
		panic(&e.CodeError{Code: e.MethodNotRegisteredCode, Message: "Method not registered."})
	}

	handler, routeKey := routes[path]

	if !routeKey {
		panic(&e.CodeError{Code: e.RouteNotRegisteredCode, Message: "Route not registered."})
	}

	return handler
}
