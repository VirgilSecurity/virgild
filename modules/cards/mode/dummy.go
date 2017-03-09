package mode

import (
	"encoding/hex"
	"time"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	virgil "gopkg.in/virgil.v4"
)

type DummyCardsMiddleware struct {
}

func (dcm *DummyCardsMiddleware) Get(id string) (card *core.Card, err error) {
	return nil, core.ErrorEntityNotFound
}

func (dcm *DummyCardsMiddleware) Search(criteria *virgil.Criteria) ([]core.Card, error) {
	return make([]core.Card, 0), nil
}

func (dcm *DummyCardsMiddleware) Create(req *core.CreateCardRequest) (*core.Card, error) {
	id := hex.EncodeToString(virgil.Crypto().CalculateFingerprint(req.Request.Snapshot))

	return &core.Card{
		ID:       id,
		Snapshot: req.Request.Snapshot,
		Meta: core.CardMeta{
			CreatedAt:   time.Now().Format(time.RFC3339),
			CardVersion: "v4",
			Signatures:  req.Request.Meta.Signatures,
			Relations:   make(map[string][]byte),
		},
	}, nil
}

func (dcm *DummyCardsMiddleware) Revoke(req *core.RevokeCardRequest) error {
	return nil
}
