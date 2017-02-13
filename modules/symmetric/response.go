package symmetric

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type Logger interface {
	Printf(format string, args ...interface{})
}

type Response func(ctx *fasthttp.RequestCtx) (interface{}, error)

func MakeResponseWrapper(logger Logger) func(f Response) fasthttp.RequestHandler {
	return func(f Response) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			s, err := f(ctx)

			ctx.SetContentType("application/json")
			ctx.ResetBody()

			if err != nil {
				responseError(err, ctx, logger)
			} else {
				responseSeccess(s, ctx, logger)
			}
		}
	}
}

type responseErrorModel struct {
	Code    ResponseErrorCode `json:"code"`
	Message string            `json:"message"`
}

func responseError(err error, ctx *fasthttp.RequestCtx, logger Logger) {
	initErr := errors.Cause(err)
	switch e := initErr.(type) {
	case ResponseErrorCode:
		if e == ErrorEntityNotFound {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusBadRequest)

		json.NewEncoder(ctx).Encode(responseErrorModel{
			Code:    e,
			Message: mapCode2Msg(e),
		})
	default:
		logger.Printf("Intrenal error: %+v", err)

		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		json.NewEncoder(ctx).Encode(responseErrorModel{
			Code:    ResponseErrorCode(10000),
			Message: mapCode2Msg(10000),
		})
	}
}

func responseSeccess(model interface{}, ctx *fasthttp.RequestCtx, logger Logger) {
	err := json.NewEncoder(ctx).Encode(model)
	if err != nil {
		responseError(err, ctx, logger)
	}
}

var code2Resp = map[ResponseErrorCode]string{
	10000: `Internal application error. You know, shit happens, so do internal server errors. Just take a deep breath and try harder.`,
	// 	30000: `JSON specified as a request is invalid`,
	// 	30010: `A data inconsistency error`,
	// 	30100: `Global Virgil Card identity type is invalid, because it can be only an 'email'`,
	// 	30101: `Virgil Card scope must be either 'global' or 'application'`,
	// 	30102: `Virgil Card id validation failed`,
	// 	30103: `Virgil Card data parameter cannot contain more than 16 entries`,
	// 	30104: `Virgil Card info parameter cannot be empty if specified and must contain 'device' and/or 'device_name' key`,
	// 	30105: `Virgil Card info parameters length validation failed. The value must be a string and mustn't exceed 256 characters`,
	// 	30106: `Virgil Card data parameter must be an associative array (https://en.wikipedia.org/wiki/Associative_array)`,
	// 	30107: `A CSR parameter (content_snapshot) parameter is missing or is incorrect`,
	// 	30111: `Virgil Card identities passed to search endpoint must be a list of non-empty strings`,
	// 	30113: `Virgil Card identity type is invalid`,
	// 	30114: `Segregated Virgil Card custom identity value must be a not empty string`,
	// 	30115: `Virgil Card identity email is invalid`,
	// 	30116: `Virgil Card identity application is invalid`,
	// 	30117: `Public key length is invalid. It goes from 16 to 2048 bytes`,
	// 	30118: `Public key must be base64-encoded string`,
	// 	30119: `Virgil Card data parameter must be a key/value list of strings`,
	// 	30120: `Virgil Card data parameters must be strings`,
	// 	30121: `Virgil Card custom data entry value length validation failed. It mustn't exceed 256 characters`,
	// 	30123: `SCR signs list parameter is missing or is invalid`,
	// 	30128: `SCR sign item is invalid or missing for the application`,
	// 	30131: `Virgil Card id specified in the request body must match with the one passed in the URL`,
	// 	30134: `Virgil Card data parameters key must be aplphanumerical`,
	// 	30137: `Global Virigl Card cannot be created unconfirmed (which means that Virgil Identity service sign is mandatory)`,
	// 	30138: `Virigl Card with the same fingerprint exists already`,
	// 	30139: `Virigl Card revocation reason isn't specified or is invalid`,
	// 	30140: `SCR sign validation failed`,
	// 	30141: `SCR one of signers Virgil Cards is not found `,
	// 	30142: `SCR sign item is invalid or missing for the Client`,
	// 	30143: `SCR sign item is invalid or missing for the Virgil Registration Authority service`,
}

func mapCode2Msg(code ResponseErrorCode) string {
	if msg, ok := code2Resp[code]; ok {
		return msg
	}
	return "Unknow response error"
}
