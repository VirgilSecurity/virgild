package symmetric

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type SymmetricRepo interface {
	Create(k SymmetricKey) error
	Remove(keyID, userID string) error
	Get(keyID, userID string) (k *SymmetricKey, err error)
	KeysByUser(userID string) (ks []SymmetricKey, err error)
	UsersByKey(keyID string) (ks []SymmetricKey, err error)
}

func createKey(repo SymmetricRepo) Response {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		var k SymmetricKey
		err := json.Unmarshal(ctx.PostBody(), &k)
		if err != nil {
			return nil, ErrorJSONIsInvalid
		}
		err = repo.Create(k)
		if err != nil {
			return nil, errors.Wrap(err, "Cannnot store symmetric key")
		}
		return nil, nil
	}
}

func getKey(repo SymmetricRepo) Response {
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

type keyUserModel struct {
	KeyID  string `json:"key_id"`
	UserID string `json:"user_id"`
}

func symmetricKyes2KeyUserModels(s []SymmetricKey) (kps []keyUserModel) {
	for _, v := range s {
		kps = append(kps, keyUserModel{v.KeyID, v.UserID})
	}
	return
}

func getUsersByKey(repo SymmetricRepo) Response {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		keyID := ctx.UserValue("key_id").(string)
		k, err := repo.UsersByKey(keyID)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Cannnot get symmetric key (key id: %v )", keyID))
		}
		return symmetricKyes2KeyUserModels(k), nil
	}
}

func getKeysByUser(repo SymmetricRepo) Response {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		userID := ctx.UserValue("user_id").(string)
		k, err := repo.KeysByUser(userID)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Cannnot get symmetric key (user id: %v)", userID))
		}
		return symmetricKyes2KeyUserModels(k), nil
	}
}
