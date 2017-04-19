package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/VirgilSecurity/virgild/modules/card/core"

	virgil "gopkg.in/virgil.v4"
)

func GetCard(f core.GetCardHandler) coreapi.APIHandler {
	return func(req *http.Request) (interface{}, error) {
		id := req.URL.Query().Get(":id")
		return f(req.Context(), id)
	}
}

func SearchCards(f core.SearchCardsHandler) coreapi.APIHandler {
	return func(req *http.Request) (interface{}, error) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, core.JSONInvalidErr
		}

		crit := new(virgil.Criteria)
		err = json.Unmarshal(body, crit)
		if err != nil {
			return nil, core.JSONInvalidErr
		}
		return f(req.Context(), crit)
	}
}

func CreateCard(f core.CreateCardHandler) coreapi.APIHandler {
	return func(req *http.Request) (interface{}, error) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, core.JSONInvalidErr
		}

		var createReq virgil.SignableRequest
		err = json.Unmarshal(body, &createReq)
		if err != nil {
			return nil, core.JSONInvalidErr
		}
		var info virgil.CardModel
		err = json.Unmarshal(createReq.Snapshot, &info)
		if err != nil {
			return nil, core.SnapshotIncorrectErr
		}

		ccr := &core.CreateCardRequest{
			Info:    info,
			Request: createReq,
		}

		return f(req.Context(), ccr)
	}
}

func RevokeCard(f core.RevokeCardHandler) coreapi.APIHandler {
	return func(req *http.Request) (interface{}, error) {
		id := req.URL.Query().Get(":id")
		ctx := core.SetURLCardID(req.Context(), id)
		req = req.WithContext(ctx)

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, core.JSONInvalidErr
		}

		var revokeReq virgil.SignableRequest
		err = json.Unmarshal(body, &revokeReq)
		if err != nil {
			return nil, core.JSONInvalidErr
		}

		var info virgil.RevokeCardRequest
		err = json.Unmarshal(revokeReq.Snapshot, &info)
		if err != nil {
			return nil, core.SnapshotIncorrectErr
		}

		rcr := &core.RevokeCardRequest{
			Info:    info,
			Request: revokeReq,
		}

		err = f(req.Context(), rcr)
		if err != nil {
			return nil, err
		}
		return []byte("{}"), nil
	}
}

func CreateRelation(f core.CreateRelationHandler) coreapi.APIHandler {
	return func(req *http.Request) (interface{}, error) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, core.JSONInvalidErr
		}

		var createReq virgil.SignableRequest
		err = json.Unmarshal(body, &createReq)
		if err != nil {
			return nil, core.JSONInvalidErr
		}
		return f(req.Context(), &core.CreateRelationRequest{
			ID:      req.URL.Query().Get(":id"),
			Request: createReq,
		})
	}
}
func RevokeRelation(f core.RevokeRelationHandler) coreapi.APIHandler {
	return func(req *http.Request) (interface{}, error) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, core.JSONInvalidErr
		}

		var revokeReq virgil.SignableRequest
		err = json.Unmarshal(body, &revokeReq)
		if err != nil {
			return nil, core.JSONInvalidErr
		}

		var info virgil.RevokeCardRequest
		err = json.Unmarshal(revokeReq.Snapshot, &info)
		if err != nil {
			return nil, core.SnapshotIncorrectErr
		}

		return f(req.Context(), &core.RevokeRelationRequest{
			ID:      req.URL.Query().Get(":id"),
			Info:    info,
			Request: revokeReq,
		})
	}
}
