package middleware

import (
	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"gopkg.in/virgil.v4"
)

func SetApplicationScopForSearch(next core.SearchCards) core.SearchCards {
	return func(crit *virgil.Criteria) ([]core.Card, error) {
		if crit.Scope != virgil.CardScope.Global {
			crit.Scope = virgil.CardScope.Application
		}
		return next(crit)
	}
}
