package middleware

import (
	"testing"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/stretchr/testify/assert"
	"gopkg.in/virgil.v4"
)

func TestSetApplicationScopForSearch_ScopeGlobal_Skeep(t *testing.T) {
	c := virgil.Criteria{
		Scope: virgil.CardScope.Global,
	}

	s := SetApplicationScopForSearch(func(c *virgil.Criteria) ([]core.Card, error) {
		assert.Equal(t, virgil.CardScope.Global, c.Scope)
		return nil, nil
	})
	s(&c)
}
func TestSetApplicationScopForSearch_ScopeOther_Skeep(t *testing.T) {
	c := virgil.Criteria{
		Scope: virgil.Enum("test"),
	}

	s := SetApplicationScopForSearch(func(c *virgil.Criteria) ([]core.Card, error) {
		assert.Equal(t, virgil.CardScope.Application, c.Scope)
		return nil, nil
	})
	s(&c)
}
