package auth

import (
	"strings"

	"github.com/valyala/fasthttp"
)

func static(t string) authHandler {
	return func(ctx *fasthttp.RequestCtx) (bool, error) {
		token := string(ctx.Request.Header.Peek("Authorization"))
		if !strings.HasSuffix(token, t) {
			return false, errTokenInvalid
		}
		return true, nil
	}
}
