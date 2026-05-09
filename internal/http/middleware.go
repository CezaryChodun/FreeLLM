package http

import nethttp "net/http"

type Middleware func(nethttp.Handler) nethttp.Handler

func Chain(handler nethttp.Handler, middlewares ...Middleware) nethttp.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
