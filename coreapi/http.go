package coreapi

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func wrapAPIHandler(logger Logger) func(fun APIHandler) http.Handler {
	return func(handler APIHandler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ok bool
			w.Header().Set("Content-Type", "application/json")

			seccess, err := handler(r)
			if err != nil {
				var apiErr APIError

				innerErr := errors.Cause(err)
				if apiErr, ok = innerErr.(APIError); !ok {
					logger.Err("API wrapper: %+v", err)
					apiErr = InternalServerErr
				}
				w.WriteHeader(apiErr.StatusCode)

				if apiErr == EntityNotFoundErr {
					return
				}

				b, _ := json.Marshal(apiErr)
				w.Write(b)
				return
			}

			w.WriteHeader(http.StatusOK)

			if seccess == nil {
				return
			}

			var body []byte

			if body, ok = seccess.([]byte); !ok {
				body, _ = json.Marshal(seccess)
			}
			w.Write(body)
		})
	}
}
