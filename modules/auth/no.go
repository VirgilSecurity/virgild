package auth

import "github.com/valyala/fasthttp"

func noAuth(permission string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return next
}
