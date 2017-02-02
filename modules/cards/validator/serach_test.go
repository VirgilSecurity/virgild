package validator

import (
	"testing"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/stretchr/testify/assert"
	"gopkg.in/virgil.v4"
)

func TestSearchCards_ScopeIncorrect_ReturnErr(t *testing.T) {
	c := virgil.Criteria{}
	s := SearchCards(func(crit *virgil.Criteria) ([]core.Card, error) {
		return nil, nil
	})
	_, err := s(&c)

	assert.Equal(t, core.ErrorScopeMustBeGlobalOrApplication, err)
}

func TestSearchCards_IdentityEmpty_ReturnErr(t *testing.T) {
	c := virgil.Criteria{Scope: virgil.CardScope.Application}
	s := SearchCards(func(crit *virgil.Criteria) ([]core.Card, error) {
		return nil, nil
	})
	_, err := s(&c)

	assert.Equal(t, core.ErrorSearchIdentitesEmpty, err)
}

func TestSearchCards_CriteriaCorrect_NextExecuted(t *testing.T) {
	c := virgil.Criteria{Scope: virgil.CardScope.Application, Identities: []string{"test"}}
	var executed bool
	s := SearchCards(func(crit *virgil.Criteria) ([]core.Card, error) {
		executed = true
		return nil, nil
	})
	s(&c)
	assert.True(t, executed)
}
