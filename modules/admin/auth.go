package admin

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

var authPrefix = []byte("Basic ")

type authValidator func(t string) error

func adminWrap(auth authValidator) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			hauth := ctx.Request.Header.Peek("Authorization")
			if !bytes.HasPrefix(hauth, authPrefix) {
				ctx.Response.Header.Add("WWW-Authenticate", "Basic")
				ctx.SetStatusCode(fasthttp.StatusUnauthorized)
				return
			}

			err := auth(string(hauth[len(authPrefix):]))
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusUnauthorized)
				return
			}
			next(ctx)
		}
	}
}

func staticLogin(login, password string) authValidator {
	return func(token string) error {
		u, p, ok := parseLoginPassword(token)
		if !ok {
			return fmt.Errorf("Cannot pars Base auth")
		}
		hp := sha256.Sum256([]byte(p))
		if login == u && password == hex.EncodeToString(hp[:]) {
			return nil
		}
		return fmt.Errorf("User name or password is not equal")
	}
}

func parseLoginPassword(token string) (string, string, bool) {
	c, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", "", false
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return "", "", false
	}
	return cs[:s], cs[s+1:], true
}
