package auth

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type authHandler func(ctx *fasthttp.RequestCtx) (bool, error)

type errResponse int

func (e errResponse) Error() string {
	return fmt.Sprintf("code: %v", e)
}

var (
	errTokenInvalid         errResponse = 20300
	errAuthServiceReturnErr errResponse = 20301
	errAuthServiceDenny     errResponse = 20302
)
