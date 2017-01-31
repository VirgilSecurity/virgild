package validator

import (
	"github.com/VirgilSecurity/virgild/modules/cards/core"
	virgil "gopkg.in/virgil.v4"
)

var searchValidator = []func(c *virgil.Criteria) (bool, error){
	scopeMustGlobalOrApplication,
	searchIdentitiesNotEmpty,
}

func SearchCards(next core.SearchCards) core.SearchCards {
	return func(c *virgil.Criteria) ([]core.Card, error) {
		for _, v := range searchValidator {
			if ok, err := v(c); !ok {
				return nil, err
			}
		}
		return next(c)
	}
}

func scopeMustGlobalOrApplication(criteria *virgil.Criteria) (bool, error) {
	if criteria.Scope == virgil.CardScope.Application || criteria.Scope == virgil.CardScope.Global {
		return true, nil
	}
	return false, core.ErrorScopeMustBeGlobalOrApplication
}

func searchIdentitiesNotEmpty(crit *virgil.Criteria) (bool, error) {
	if len(crit.Identities) == 0 {
		return false, core.ErrorSearchIdentitesEmpty
	}
	return true, nil
}
