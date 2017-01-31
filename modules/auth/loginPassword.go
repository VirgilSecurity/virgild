package auth

import (
	"encoding/base64"
	"strings"

	"github.com/valyala/fasthttp"
)

func loginPasswrod(login, password string) authHandler {
	return func(ctx *fasthttp.RequestCtx) (bool, error) {
		auth := ctx.Request.Header.Peek("Authorization")
		if auth == nil {
			return false, errTokenInvalid
		}
		u, p, ok := parseBasicAuth(string(auth))
		if !ok {
			return false, errTokenInvalid
		}
		if u == login && p == password {
			return true, nil
		}
		return false, errTokenInvalid
	}
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}
