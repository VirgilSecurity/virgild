package mode

import (
	"testing"

	"gopkg.in/virgil.v4"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/stretchr/testify/assert"
)

func TestDummyGet_ReturnNotFound(t *testing.T) {
	m := DummyCardsMiddleware{}
	_, err := m.Get("134")
	assert.Equal(t, core.ErrorEntityNotFound, err)
}

func TestDummySearch_ReturnEmptyVal(t *testing.T) {
	m := DummyCardsMiddleware{}
	cards, _ := m.Search(&virgil.Criteria{})
	assert.Len(t, cards, 0)
}

func TestDummyCreate_ReturnNotFound(t *testing.T) {
	m := DummyCardsMiddleware{}
	card, _ := m.Create(&core.CreateCardRequest{})

	assert.NotEmpty(t, card.ID)
	assert.Equal(t, "v4", card.Meta.CardVersion)
	assert.NotEmpty(t, card.Meta.CreatedAt)
}

func TestDummyRevoke_ReturnNil(t *testing.T) {
	m := DummyCardsMiddleware{}
	err := m.Revoke(&core.RevokeCardRequest{})
	assert.NoError(t, err)
}
