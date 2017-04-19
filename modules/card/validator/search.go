package validator

import (
	"context"

	"github.com/VirgilSecurity/virgild/modules/card/core"

	virgil "gopkg.in/virgil.v4"
)

type searchCardsValidatorHandler func(ctx context.Context, c *virgil.Criteria) (bool, error)

var searchValidator = []searchCardsValidatorHandler{
	scopeMustGlobalOrApplication,
	searchIdentitiesNotEmpty,
}

func SearchCards(f core.SearchCardsHandler) core.SearchCardsHandler {
	return func(ctx context.Context, c *virgil.Criteria) ([]virgil.CardResponse, error) {
		for _, v := range searchValidator {
			if ok, err := v(ctx, c); !ok {
				return nil, err
			}
		}
		return f(ctx, c)
	}
}

func scopeMustGlobalOrApplication(ctx context.Context, criteria *virgil.Criteria) (bool, error) {
	if criteria.Scope == virgil.CardScope.Application || criteria.Scope == virgil.CardScope.Global {
		return true, nil
	}
	return false, core.ScopeMustBeGlobalOrApplicationErr
}

func searchIdentitiesNotEmpty(ctx context.Context, crit *virgil.Criteria) (bool, error) {
	if len(crit.Identities) == 0 {
		return false, core.SearchIdentitesEmptyErr
	}
	return true, nil
}
