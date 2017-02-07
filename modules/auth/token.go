package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type permission int

const (
	get    permission = 1
	search permission = 2
	create permission = 4
	revoke permission = 8
)

type token struct {
	Token      string `xorm:"PK"`
	Permission permission
}

type tokenRepo interface {
	All() ([]token, error)
	Remove(token string) error
	Get(token string) (*token, error)
	Create(token token) error
}

type tokenHandler struct {
	repo tokenRepo
}

type tokenModel struct {
	Token string          `json:"token,omitempty"`
	Perm  map[string]bool `json:"permissions"`
}

func (h *tokenHandler) All(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	t, err := h.repo.All()
	if err != nil {
		ctx.Error("", fasthttp.StatusInternalServerError)
		return
	}
	resp := make([]tokenModel, 0)
	for _, v := range t {
		resp = append(resp, tokenModel{v.Token, perm2Map(v.Permission)})
	}

	b, err := json.Marshal(resp)
	if err != nil {
		ctx.Error("", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Write(b)
}

func (h *tokenHandler) Remove(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	id, ok := ctx.UserValue("id").(string)
	if !ok {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}
	err := h.repo.Remove(id)
	if err != nil {
		ctx.Error("", fasthttp.StatusInternalServerError)
		return
	}
}

func (h *tokenHandler) Update(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	id, ok := ctx.UserValue("id").(string)
	if !ok {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}
	var perm tokenModel
	err := json.Unmarshal(ctx.PostBody(), &perm)
	if err != nil {
		ctx.Error(`{"message":"JSON body invalid"}`, fasthttp.StatusBadRequest)
		return
	}

	err = h.repo.Remove(id)
	if err != nil {
		ctx.Error("", fasthttp.StatusInternalServerError)
		return
	}
	err = h.repo.Create(token{id, map2Perm(perm.Perm)})
	if err != nil {
		ctx.Error("", fasthttp.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(tokenModel{id, perm.Perm})
	if err != nil {
		ctx.Error("", fasthttp.StatusInternalServerError)
		return
	}
	ctx.Write(b)
}
func (h *tokenHandler) Create(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	var perm tokenModel
	err := json.Unmarshal(ctx.PostBody(), &perm)
	if err != nil {
		ctx.Error(`{"message":"JSON body invalid"}`, fasthttp.StatusBadRequest)
		return
	}

	id := make([]byte, 32)
	rand.Read(id)
	t := token{hex.EncodeToString(id), map2Perm(perm.Perm)}
	err = h.repo.Create(t)
	if err != nil {
		ctx.Error("", fasthttp.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(tokenModel{t.Token, perm2Map(t.Permission)})
	if err != nil {
		ctx.Error("", fasthttp.StatusInternalServerError)
		return
	}
	ctx.Write(b)
}

func auth(repo tokenRepo, permName string) func(token string) error {
	perm := map2Perm(map[string]bool{permName: true})
	return func(t string) error {
		tt, err := repo.Get(t)
		if err != nil {
			return errTokenInvalid
		}
		if tt.Permission&perm == perm {
			return nil
		}
		return errAuthServiceDenny
	}
}

func perm2Map(perm permission) map[string]bool {
	return map[string]bool{
		PermissionGetCard:     perm&get == get,
		PermissionSearchCards: perm&search == search,
		PermissionCreateCard:  perm&create == create,
		PermissionRevokeCard:  perm&revoke == revoke,
	}
}

func map2Perm(m map[string]bool) permission {
	r := Btoi(m[PermissionGetCard])*int(get) +
		Btoi(m[PermissionSearchCards])*int(search) +
		Btoi(m[PermissionCreateCard])*int(create) +
		Btoi(m[PermissionRevokeCard])*int(revoke)
	return permission(r)
}

func Btoi(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
