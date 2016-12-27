package handler

import (
	"encoding/base64"
	"encoding/json"

	"github.com/virgilsecurity/virgild/core"
	"github.com/virgilsecurity/virgild/core/storage"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4"
)

type Logger interface {
	Printf(format string, args ...interface{})
}

type RequestSigner interface {
	Sign(*virgil.SignableRequest) error
}

type response struct {
	ID       string       `json:"id"`
	Snapshot []byte       `json:"content_snapshot"` // the raw serialized version of CardRequest
	Meta     responseMeta `json:"meta"`
}
type responseMeta struct {
	CreatedAt   string            `json:"created_at,omitempty"`
	CardVersion string            `json:"card_version,omitempty"`
	Signatures  map[string][]byte `json:"signs"`
}

type Card struct {
	Storage storage.CardStorage
	Signer  RequestSigner
	Logger  Logger
}

func (c *Card) Get(id string, resp core.Response) {
	card, err := c.Storage.GetCard(id)
	if err == storage.ErrorForbidden {
		resp.Error(core.ErrorForbidden)
		return
	}
	if err != nil {
		c.Logger.Printf("CardHandler.Get(%v): %v", id, err)
		resp.Error(core.ErrorInernalApplication)
		return
	}
	resp.Success(convert2Response(card))
}

func (c *Card) Search(criteria core.Criteria, resp core.Response) {
	crit := virgil.Criteria{
		Identities:   criteria.Identities,
		IdentityType: criteria.IdentityType,
	}
	if criteria.Scope == "global" {
		crit.Scope = virgil.CardScope.Global
	}
	cards, err := c.Storage.SearchCards(crit)
	if err != nil {
		c.Logger.Printf("CardHandler.Search: %v", err)
		resp.Error(core.ErrorInernalApplication)
		return
	}
	cs := make([]*response, 0, len(cards))
	for _, v := range cards {
		cs = append(cs, convert2Response(v))
	}
	resp.Success(cs)
}

func (c *Card) Create(req *core.Request, resp core.Response) {
	sr, err := convertCreateRequest2SignableRequest(req)
	if err != nil {
		c.Logger.Printf("CardHandler.Create[conver]: %v", err)
		resp.Error(core.ErrorInernalApplication)
		return
	}
	err = c.Signer.Sign(sr)
	if err != nil {
		c.Logger.Printf("CardHandler.Create[sign]: %v", err)
		resp.Error(core.ErrorInernalApplication)
		return
	}
	card, err := c.Storage.CreateCard(sr)
	if err != nil {
		c.Logger.Printf("CardHandler.Create[create]: %v", err)
		resp.Error(core.ErrorInernalApplication)
		return
	}
	resp.Success(convert2Response(card))
}

func (c *Card) Revoke(req *core.Request, resp core.Response) {
	sr, err := convertRevokeRequest2SignableRequest(req)
	if err != nil {
		c.Logger.Printf("CardHandler.Revoke[conver]: %v", err)
		resp.Error(core.ErrorInernalApplication)
		return
	}
	err = c.Signer.Sign(sr)
	if err != nil {
		c.Logger.Printf("CardHandler.Revoke[sign]: %v", err)
		resp.Error(core.ErrorInernalApplication)
		return
	}
	err = c.Storage.RevokeCard(sr)
	if err != nil {
		c.Logger.Printf("CardHandler.Revoke[revoke]: %v", err)
		resp.Error(core.ErrorInernalApplication)
		return
	}
	resp.Success(nil)
}

func convertCreateRequest2SignableRequest(req *core.Request) (*virgil.SignableRequest, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	str := base64.StdEncoding.EncodeToString(b)
	return virgil.ImportCreateCardRequest([]byte(str))
}

func convertRevokeRequest2SignableRequest(req *core.Request) (*virgil.SignableRequest, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	str := base64.StdEncoding.EncodeToString(b)
	return virgil.ImportRevokeCardRequest([]byte(str))
}

func convert2Response(c *virgil.Card) *response {
	return &response{
		ID:       c.ID,
		Snapshot: c.Snapshot,
		Meta: responseMeta{
			CardVersion: c.CardVersion,
			Signatures:  c.Signatures,
			CreatedAt:   c.CreatedAt,
		},
	}
}
