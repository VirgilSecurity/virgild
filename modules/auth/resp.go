package auth

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/valyala/fasthttp"
)

type logger interface {
	Printf(format string, args ...interface{})
}

type scopes func(t string) ([]string, error)

func wrap(l logger, tokenType string, getScopes scopes) func(string, fasthttp.RequestHandler) fasthttp.RequestHandler {
	ttSpace := append([]byte(tokenType), ' ')
	errWrap := errWrapFunc(l)
	return func(scope string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {

			hauth := ctx.Request.Header.Peek("Authorization")
			if !bytes.HasPrefix(hauth, ttSpace) {
				errWrap(errTokenInvalid, ctx)
				return
			}

			scopes, err := getScopes(string(hauth[len(ttSpace):]))
			if err != nil {
				errWrap(errTokenInvalid, ctx)
				return
			}
			ctx.Response.Header.Add("X-OAuth-Scopes", strings.Join(scopes, ", "))
			ctx.Response.Header.Add("X-Accepted-OAuth-Scopes", scope)

			for _, s := range scopes {
				if s == scope || s == "*" {
					next(ctx)
					return
				}
			}
			errWrap(errAuthServiceDenny, ctx)
		}
	}
}

type respErr struct {
	Code    errResponse `json:"code"`
	Message string      `json:"message"`
}

func errWrapFunc(l logger) func(err error, ctx *fasthttp.RequestCtx) {
	return func(err error, ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		switch e := err.(type) {
		case errResponse:
			msg, ok := errMap[e]
			if !ok {
				msg = "Unknow response error"
			}
			json.NewEncoder(ctx).Encode(respErr{e, msg})
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		default:
			l.Printf("Auth card request: %v", err)
			ctx.Error("", fasthttp.StatusInternalServerError)
		}
	}
}

var errMap = map[errResponse]string{
	20300: `The Virgil access token or token header was not specified or is invalid`,
	20301: `The Virgil authenticator service responded with an error`,
	20302: `The Virgil access token validation has failed on the Virgil Authenticator service`,
	20303: `The application was not found for the access token`,
	20500: `The Virgil Card is not available`,
}
