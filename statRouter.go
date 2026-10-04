package goter

import (
	e "github.com/adagit94/err"
)

type statRoutes[H any] map[string]H
type statTree[H any] map[string]statRoutes[H]

// Creates static, faster router without support for dynamic segments and other dynamic functionalities that uses plain map based string paths matching.
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

// Create new route for which respective http methods handlers can be defined.
func (router *statRouter[H]) Route(path string) IRoute[H] {
	route := &route[H]{path: path, registerHandler: router.registerHandler}
	return route
}

// Select is mean't to be wrapped inside of specific http implementation entry point for incoming requests. It returns selected handler for a found route. It panics in case no route gets matched or method doesn't exists in register, so it's important to recover and handle such potential errors gracefully in outer code.
func (router *statRouter[H]) Select(path string, method string) H {
	routes, methodKey := router.tree[method]

	if !methodKey {
		panic(&e.Err{Code: e.MethodNotRegisteredCode, Message: "Method not registered."})
	}

	handler, routeKey := routes[path]

	if !routeKey {
		panic(&e.Err{Code: e.RouteNotRegisteredCode, Message: "Route not registered."})
	}

	return handler
}
