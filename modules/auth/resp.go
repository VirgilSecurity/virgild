package auth

import (
	"bytes"
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type Logger interface {
	Printf(format string, args ...interface{})
}

type cards struct {
	Code    errResponse `json:"code"`
	Message string      `json:"message"`
}

func getToken(t []byte, auth func(token string) error) authHandler {
	typeSpace := append(t, ' ')
	return func(ctx *fasthttp.RequestCtx) error {
		hauth := ctx.Request.Header.Peek("Authorization")
		if !bytes.HasPrefix(hauth, typeSpace) {
			return errTokenInvalid
		}
		return auth(string(hauth[len(typeSpace):]))
	}
}

func cardsWrap(logger Logger, validator authHandler, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		err := validator(ctx)
		if err != nil {
			switch e := err.(type) {
			case errResponse:
				msg, ok := errMap[e]
				if !ok {
					msg = "Unknow response error"
				}
				b, err := json.Marshal(cards{e, msg})
				if err != nil {
					logger.Printf("Auth card request. Cannot marshal response: %v", e)
					ctx.Error("", fasthttp.StatusInternalServerError)
					return
				}
				ctx.SetContentType("application/json")
				ctx.Write(b)
				ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			default:
				logger.Printf("Auth card request: %v", err)
				ctx.Error("", fasthttp.StatusInternalServerError)
				return
			}
		}
		next(ctx)
	}
}

var errMap = map[errResponse]string{
	20300: `The Virgil access token or token header was not specified or is invalid`,
	20301: `The Virgil authenticator service responded with an error`,
	20302: `The Virgil access token validation has failed on the Virgil Authenticator service`,
	20303: `The application was not found for the access token`,
}
