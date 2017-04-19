package healthcheck

import (
	"net/http"

	"github.com/VirgilSecurity/virgild/coreapi"
)

func Init(c coreapi.Core) {
	c.HTTP.Router.Get("/health/status", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}
