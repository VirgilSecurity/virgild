package auth

import "github.com/valyala/fasthttp"

func noAuth(permission string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx fasthttp.RequestCtx) {
		if !ctx.RemoteIP().IsLoopback() &&
			(permission == PermissionCreateCard ||
				permission == PermissionRevokeCard) {

			errWrap(errForbidden, ctx)
			return
		}
		next(ctx)
	}
}
