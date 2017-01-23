package main

import (
	"encoding/json"
	"time"

	virgil "gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/errors"
)

func vcard2SqlCard(vcard *virgil.Card) (*cardSql, error) {
	card := vcard2Card(vcard)
	jcard, err := json.Marshal(card)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &cardSql{
		CardID:       vcard.ID,
		Identity:     vcard.Identity,
		IdentityType: vcard.IdentityType,
		Scope:        string(vcard.Scope),
		Card:         jcard,
	}, nil
}

func vcard2Card(vcard *virgil.Card) *Card {
	return &Card{
		ID:       vcard.ID,
		Snapshot: vcard.Snapshot,
		Meta: CardMeta{
			CreatedAt:   vcard.CreatedAt,
			CardVersion: vcard.CardVersion,
			Signatures:  vcard.Signatures,
		},
	}
}

func sqlCard2Card(sql *cardSql) (*Card, error) {
	if sql.ErrorCode != 0 {
		return nil, ResponseErrorCode(sql.ErrorCode)
	}
	card := new(Card)
	err := json.Unmarshal(sql.Card, card)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return card, nil
}

func sqlCards2Cards(sql []cardSql) ([]*Card, error) {
	cards := make([]*Card, 0)
	for _, v := range sql {
		c, err := sqlCard2Card(&v)
		if err != nil {
			if _, ok := err.(ResponseErrorCode); !ok {
				return nil, err
			}
			continue
		}
		cards = append(cards, c)
	}
	return cards, nil
}

type AppModeCardHandler struct {
	Repo      CardRepository
	Signer    RequestSigner
	Validator Validator
	Remote    VirgilClient
}

func (h *AppModeCardHandler) remoteGet(id string) (interface{}, error) {
	vc, err := h.Remote.GetCard(id)
	if err != nil {
		verr, ok := errors.ToSdkError(err)
		if ok {
			code := verr.ServiceErrorCode()
			if verr.HTTPErrorCode() == 404 {
				code = int(ErrorEntityNotFound)
			}
			h.Repo.Add(cardSql{
				CardID:    id,
				ErrorCode: code,
			})
		}
		return nil, errors.Wrap(err, "")
	}
	sqlCard, err := vcard2SqlCard(vc)
	if err != nil {
		return nil, err
	}
	h.Repo.Add(*sqlCard)
	return vcard2Card(vc), nil
}

func (h *AppModeCardHandler) Get(id string) (interface{}, error) {
	c, err := h.Repo.Get(id)
	if err == ErrorEntityNotFound {
		return h.remoteGet(id)
	}
	if err != nil {
		return nil, err
	}
	if c.ExpireAt.After(time.Now()) {
		h.Repo.DeleteById(id)
		return h.remoteGet(id)
	}

	return sqlCard2Card(c)
}

func (h *AppModeCardHandler) remoteSearch(criteria *virgil.Criteria) (interface{}, error) {
	vcards, err := h.Remote.SearchCards(*criteria)
	if err != nil {
		verr, ok := errors.ToSdkError(err)
		if ok {
			for k, _ := range criteria.Identities {
				h.Repo.Add(cardSql{
					Identity:     criteria.Identities[k],
					IdentityType: criteria.IdentityType,
					Scope:        string(criteria.Scope),
					ErrorCode:    verr.ServiceErrorCode(),
				})
			}
		}
		return nil, errors.Wrap(err, "")
	}

	cards := make([]*Card, 0)
	for _, vc := range vcards {

		sqlCard, err := vcard2SqlCard(vc)
		if err != nil {
			return nil, err
		}
		h.Repo.Add(*sqlCard)

		cards = append(cards, vcard2Card(vc))
	}
	return cards, nil
}

func (h *AppModeCardHandler) Search(criteria *virgil.Criteria) (interface{}, error) {
	ok, err := h.Validator.IsValidSearchCriteria(criteria)
	if !ok {
		return nil, err
	}
	cards, err := h.Repo.Find(criteria.Identities, criteria.IdentityType, string(criteria.Scope))
	if err != nil {
		return nil, err
	}
	if len(cards) == 0 {
		return h.remoteSearch(criteria)
	}

	for _, v := range cards {
		if v.ExpireAt.After(time.Now()) {
			h.Repo.DeleteBySearch(criteria.Identities, criteria.IdentityType, string(criteria.Scope))
			return h.remoteSearch(criteria)
		}
	}
	return sqlCards2Cards(cards)
}

func (h *AppModeCardHandler) Create(req *CreateCardRequest) (interface{}, error) {
	ok, err := h.Validator.IsValidCreateCardRequest(req)
	if !ok {
		return nil, err
	}
	err = h.Signer.Sign(&req.Request)
	if err != nil {
		return nil, err
	}
	vcard, err := h.Remote.CreateCard(&req.Request)
	if err != nil {
		return nil, err
	}
	sqlCard, err := vcard2SqlCard(vcard)
	if err != nil {
		return nil, err
	}

	h.Repo.Add(*sqlCard)
	return vcard2Card(vcard), nil
}

func (h *AppModeCardHandler) Revoke(req *RevokeCardRequest) (interface{}, error) {
	ok, err := h.Validator.IsValidRevokeCardRequest(req)
	if !ok {
		return nil, err
	}
	err = h.Signer.Sign(&req.Request)
	if err != nil {
		return nil, err
	}
	err = h.Remote.RevokeCard(&req.Request)
	if err != nil {
		return nil, err
	}

	h.Repo.DeleteById(req.Info.ID)
	return nil, nil
}
