package middleware

import (
	"net/http"
	"strings"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/VirgilSecurity/virgild/modules/card/core"
)

var tokenType = "VIRGIL "

func RequestOwner(next coreapi.APIHandler) coreapi.APIHandler {
	return func(req *http.Request) (interface{}, error) {
		authHeader := req.Header.Get("Authorization")

		// maybe it's global request
		if len(authHeader) == 0 {
			return next(req)
		}

		if !strings.HasPrefix(authHeader, tokenType) {
			return nil, core.UnsupportedAuthTypeErr
		}

		token := string(authHeader[len(tokenType):])
		ctx := core.SetAuthHeader(req.Context(), authHeader)
		ctx = core.SetOwnerRequest(ctx, token)
		req = req.WithContext(ctx)

		return next(req)
	}
}
