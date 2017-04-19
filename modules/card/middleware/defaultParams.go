package middleware

import (
	"context"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"gopkg.in/virgil.v4"
)

func SetApplicationScopForSearch(f core.SearchCardsHandler) core.SearchCardsHandler {
	return func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
		if crit.Scope != virgil.CardScope.Global {
			crit.Scope = virgil.CardScope.Application
		}
		return f(ctx, crit)
	}
}
