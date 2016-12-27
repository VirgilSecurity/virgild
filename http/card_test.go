package http

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
	"github.com/virgilsecurity/virgild/core"
)

type FakeCardHandler struct {
	mock.Mock
}

func (h *FakeCardHandler) Get(id string, resp core.Response) {
	h.Called(id, resp)
}
func (h *FakeCardHandler) Search(criteria core.Criteria, resp core.Response) {
	h.Called(criteria, resp)
}
func (h *FakeCardHandler) Create(req *core.Request, resp core.Response) {
	h.Called(req, resp)
}
func (h *FakeCardHandler) Revoke(req *core.Request, resp core.Response) {
	h.Called(req, resp)
}

func Test_JSONInvalid_ReturnErr(t *testing.T) {
	var c CardController

	table := []fasthttp.RequestHandler{
		c.Create,
		c.Revoke,
		c.Search,
	}
	for _, f := range table {
		ctx := makeRequestCtx("adf\"")
		f(ctx)
		assertResponse(t, core.ErrorJSONIsInvalid, ctx)
	}
}

func TestGet_HandlerInvoke(t *testing.T) {
	const id = "123"
	ctx := makeRequestCtx(nil)
	ctx.SetUserValue("id", id)
	h := new(FakeCardHandler)
	h.On("Get", id, mock.Anything).Once()
	c := CardController{Card: h}
	c.Get(ctx)

	h.AssertExpectations(t)
}

func TestCreate_HandlerInvoke(t *testing.T) {
	expected := &core.Request{
		Snapshot: []byte(`snapshot`),
		Meta: core.RequestMeta{
			Signatures: map[string][]byte{
				"id": []byte(`sign`),
			},
		},
	}
	ctx := makeRequestCtx(expected)

	h := new(FakeCardHandler)
	h.On("Create", expected, mock.Anything).Once()
	c := CardController{Card: h}
	c.Create(ctx)

	h.AssertExpectations(t)
}

func TestRevoke_HandlerInvoke(t *testing.T) {
	expected := &core.Request{
		Snapshot: []byte(`snapshot`),
		Meta: core.RequestMeta{
			Signatures: map[string][]byte{
				"id": []byte(`sign`),
			},
		},
	}
	ctx := makeRequestCtx(expected)

	h := new(FakeCardHandler)
	h.On("Revoke", expected, mock.Anything).Once()
	c := CardController{Card: h}
	c.Revoke(ctx)

	h.AssertExpectations(t)
}

func TestSearch_HandlerInvoke(t *testing.T) {
	expected := core.Criteria{
		Identities:   []string{"id1", "id2"},
		Scope:        "test",
		IdentityType: "type",
	}
	ctx := makeRequestCtx(expected)

	h := new(FakeCardHandler)
	h.On("Search", expected, mock.Anything).Once()
	c := CardController{Card: h}
	c.Search(ctx)

	h.AssertExpectations(t)
}
