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
			Info: creq,
			Request: virgil.SignableRequest{
				Snapshot: req.Snapshot,
				Meta: virgil.RequestMeta{
					Signatures: req.Meta.Signatures,
				},
			}})
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
			Info: creq,
			Request: virgil.SignableRequest{
				Snapshot: req.Snapshot,
				Meta: virgil.RequestMeta{
					Signatures: req.Meta.Signatures,
				},
			}})
	}
}

type CardsRepository interface {
	Count() (int64, error)
}

func GetCountCards(repo CardsRepository) Response {
	type CardsCountModel struct {
		Count int64 `json:"count"`
	}
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {

		c, err := repo.Count()
		if err != nil {
			return nil, err
		}
		return CardsCountModel{c}, nil
	}
}

type RequestStatisticRepository interface {
	Search(from int64, to int64, token string) ([]core.RequestStatistics, error)
}

func GetRequestStatistic(repo RequestStatisticRepository) Response {
	type RequestStatistic struct {
		Data     int64  `json:"data"`
		Token    string `json:"token"`
		Method   string `json:"method"`
		Resource string `Json:"resource"`
	}
	type FilterRequestStatistic struct {
		From  int64  `json:"from,omitempty"`
		To    int64  `json:"to,omitempty"`
		Toekn string `json:"token,omitempty"`
	}
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		filter := new(FilterRequestStatistic)
		err := json.Unmarshal(ctx.PostBody(), filter)
		if err != nil {
			return nil, err
		}
		return repo.Search(filter.From, filter.To, filter.Toekn)
	}
}
