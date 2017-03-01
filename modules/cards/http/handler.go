package http

import (
	"encoding/json"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/valyala/fasthttp"
	virgil "gopkg.in/virgil.v4"
)

type Response func(ctx *fasthttp.RequestCtx) (seccess interface{}, err error)

func GetCard(f core.GetCard) Response {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		id := ctx.UserValue("id").(string)
		return f(id)
	}
}

func SearchCards(f core.SearchCards) Response {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		c := new(virgil.Criteria)
		err := json.Unmarshal(ctx.PostBody(), c)
		if err != nil {
			return nil, core.ErrorJSONIsInvalid
		}
		return f(c)
	}
}

func CreateCard(f core.CreateCard) Response {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		req := new(virgil.SignableRequest)
		err := json.Unmarshal(ctx.PostBody(), req)
		if err != nil {
			return nil, core.ErrorJSONIsInvalid
		}

		creq := virgil.CardModel{}
		err = json.Unmarshal(req.Snapshot, &creq)
		if err != nil {
			return nil, core.ErrorSnapshotIncorrect
		}

		return f(&core.CreateCardRequest{
			Info:    creq,
			Request: *req})
	}
}

func RevokeCard(f core.RevokeCard) Response {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		req := new(virgil.SignableRequest)
		err := json.Unmarshal(ctx.PostBody(), req)
		if err != nil {
			return nil, core.ErrorJSONIsInvalid
		}

		creq := virgil.RevokeCardRequest{}
		err = json.Unmarshal(req.Snapshot, &creq)
		if err != nil {
			return nil, core.ErrorSnapshotIncorrect
		}

		return nil, f(&core.RevokeCardRequest{
			Info:    creq,
			Request: *req})
	}
}

type CardsRepository interface {
	Count() (int64, error)
}

type cardsCountModel struct {
	Count int64 `json:"count"`
}

func GetCountCards(repo CardsRepository) Response {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {

		c, err := repo.Count()
		if err != nil {
			return nil, err
		}
		return cardsCountModel{c}, nil
	}
}
