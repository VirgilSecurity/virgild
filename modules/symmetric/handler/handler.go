package handler

import (
	"encoding/json"
	"fmt"

	"github.com/VirgilSecurity/virgild/modules/symmetric/core"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func CreateKey(repo core.SymmetricRepo) core.Response {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		var k core.SymmetricKey
		err := json.Unmarshal(ctx.PostBody(), &k)
		if err != nil {
			return nil, core.ErrorJSONIsInvalid
		}
		err = repo.Create(k)
		if err != nil {
			return nil, errors.Wrap(err, "Cannnot store symmetric key")
		}
		return nil, nil
	}
}

func GetKey(repo core.SymmetricRepo) core.Response {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		keyID, _ := ctx.UserValue("key_id").(string)
		userID, _ := ctx.UserValue("user_id").(string)
		k, err := repo.Get(keyID, userID)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Cannnot get symmetric key (key id: %v user id: %v)", keyID, userID))
		}
		return k, nil
	}
}

func GetUsersByKey(repo core.ListSymmetricRepo) core.Response {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		keyID := ctx.UserValue("key_id").(string)
		k, err := repo.UsersByKey(keyID)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Cannnot get symmetric key (key id: %v )", keyID))
		}
		return k, nil
	}
}

func GetKeysByUser(repo core.ListSymmetricRepo) core.Response {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		userID := ctx.UserValue("user_id").(string)
		k, err := repo.KeysByUser(userID)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Cannnot get symmetric key (user id: %v)", userID))
		}
		return k, nil
	}
}
