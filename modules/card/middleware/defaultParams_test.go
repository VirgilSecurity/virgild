package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/virgil.v4"
)

func TestSetApplicationScopForSearch_ScopeGlobal_Skeep(t *testing.T) {
	crit := &virgil.Criteria{
		Scope: virgil.CardScope.Global,
	}

	search := SetApplicationScopForSearch(func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
		assert.Equal(t, virgil.CardScope.Global, crit.Scope)
		return nil, nil
	})
	search(context.Background(), crit)
}
func TestSetApplicationScopForSearch_ScopeOther_Skeep(t *testing.T) {
	crit := &virgil.Criteria{
		Scope: virgil.Enum("test"),
	}

	search := SetApplicationScopForSearch(func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
		assert.Equal(t, virgil.CardScope.Application, crit.Scope)
		return nil, nil
	})
	search(context.Background(), crit)
}
