package validator

import (
	"context"
	"testing"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/stretchr/testify/assert"
	"gopkg.in/virgil.v4"
)

func TestSearchCards_ScopeIncorrect_ReturnErr(t *testing.T) {
	crit := &virgil.Criteria{}
	search := SearchCards(func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
		return nil, nil
	})
	_, err := search(context.Background(), crit)

	assert.Equal(t, core.ScopeMustBeGlobalOrApplicationErr, err)
}

func TestSearchCards_IdentityEmpty_ReturnErr(t *testing.T) {
	crit := &virgil.Criteria{Scope: virgil.CardScope.Application}
	search := SearchCards(func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
		return nil, nil
	})
	_, err := search(context.Background(), crit)

	assert.Equal(t, core.SearchIdentitesEmptyErr, err)
}

func TestSearchCards_CriteriaCorrect_NextExecuted(t *testing.T) {
	crit := &virgil.Criteria{Scope: virgil.CardScope.Application, Identities: []string{"test"}}
	var executed bool
	search := SearchCards(func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
		executed = true
		return nil, nil
	})
	search(context.Background(), crit)
	assert.True(t, executed)
}
