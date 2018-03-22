package auth

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

func noAuth(permission string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if !ctx.RemoteIP().IsLoopback() &&
			(permission == PermissionCreateCard ||
				permission == PermissionRevokeCard) {

			json.NewEncoder(ctx).Encode(respErr{errForbidden, errMap[errForbidden]})
			ctx.SetStatusCode(fasthttp.StatusForbidden)
			return
		}
		next(ctx)
	}
}
