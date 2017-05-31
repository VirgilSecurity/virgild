package card

import (
	"context"
	"fmt"
	"testing"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/virgil.v4"
)

type fakeCache struct {
	mock.Mock
}

func (f *fakeCache) Get(key string, val interface{}) bool {
	args := f.Called(key)
	has := args.Bool(0)
	if has {
		if card, ok := val.(**virgil.CardResponse); ok {
			*card = args.Get(1).(*virgil.CardResponse)
		} else if ids, ok := val.(*[]string); ok {
			*ids = args.Get(1).([]string)
		}
	}
	return has
}

func (f *fakeCache) Set(key string, val interface{}) {
	f.Called(key, val)
}

func (f *fakeCache) Del(key string) {
	f.Called(key)
}

func TestCacheGetCard_KeyExist_ReturnVal(t *testing.T) {
	owner := "owner"
	id := "card_id"
	expected := &virgil.CardResponse{
		Snapshot: []byte(`snapshot`),
	}

	cache := new(fakeCache)
	cache.On("Get", owner+"_"+id).Return(true, expected)

	ctx := core.SetOwnerRequest(context.Background(), owner)

	cacheCard := cacheCardMiddleware{cache}
	actual, err := cacheCard.GetCard(func(ctx context.Context, id string) (*virgil.CardResponse, error) {
		t.Fatal("Function executed")
		return nil, nil
	})(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestCacheGetCard_KeyNotExist_FuncExecuted(t *testing.T) {
	owner := "owner"
	id := "card_id"
	expected := &virgil.CardResponse{
		Snapshot: []byte(`snapshot`),
	}

	cache := new(fakeCache)
	cache.On("Get", mock.Anything).Return(false)

	ctx := core.SetOwnerRequest(context.Background(), owner)

	cacheCard := cacheCardMiddleware{cache}
	actual, err := cacheCard.GetCard(func(ctx context.Context, id string) (*virgil.CardResponse, error) {
		return expected, fmt.Errorf("ERROR")
	})(ctx, id)

	assert.Error(t, err, "ERROR")
	assert.Equal(t, expected, actual)
}

func TestCacheGetCard_FuncExecuted_SetCache(t *testing.T) {
	owner := "owner"
	id := "card_id"
	expected := &virgil.CardResponse{
		Snapshot: []byte(`snapshot`),
	}

	cache := new(fakeCache)
	cache.On("Get", mock.Anything).Return(false)
	cache.On("Set", owner+"_"+id, expected).Once()

	ctx := core.SetOwnerRequest(context.Background(), owner)

	cacheCard := cacheCardMiddleware{cache}
	cacheCard.GetCard(func(ctx context.Context, id string) (*virgil.CardResponse, error) {
		return expected, nil
	})(ctx, id)

	cache.AssertExpectations(t)
}

func TestCacheSearchCards_CacheMiss_FuncExecuted(t *testing.T) {
	expected := []virgil.CardResponse{
		virgil.CardResponse{},
		virgil.CardResponse{},
	}
	cache := new(fakeCache)
	cache.On("Get", mock.Anything).Return(false)
	cache.On("Set", mock.Anything, mock.Anything)

	cacheCard := cacheCardMiddleware{cache}
	cards, err := cacheCard.SearchCards(func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
		return expected, nil
	})(context.Background(), &virgil.Criteria{})

	assert.NoError(t, err)
	assert.Equal(t, expected, cards)
}

func TestCacheSearchCards_FuncExecuted_CacheSet(t *testing.T) {
	owner := "owner"
	card1 := virgil.CardResponse{ID: "1"}
	card2 := virgil.CardResponse{ID: "2"}
	expected := []virgil.CardResponse{
		card1,
		card2,
	}
	crit := &virgil.Criteria{
		Identities:   []string{"bob", "alice"},
		Scope:        virgil.CardScope.Application,
		IdentityType: "test",
	}
	searchKey := owner + "_test_application_alice_bob"
	cache := new(fakeCache)
	cache.On("Get", mock.Anything).Return(false)
	cache.On("Set", owner+"_"+card1.ID, card1).Once()
	cache.On("Set", owner+"_"+card2.ID, card2).Once()
	cache.On("Set", searchKey, []string{"1", "2"}).Once()

	cacheCard := cacheCardMiddleware{cache}
	cacheCard.SearchCards(func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
		return expected, nil
	})(core.SetOwnerRequest(context.Background(), owner), crit)

	cache.AssertExpectations(t)
}

func TestCacheSearchCards_CacheExist(t *testing.T) {
	owner := "owner"
	card1 := virgil.CardResponse{ID: "1"}
	card2 := virgil.CardResponse{ID: "2"}
	expected := []virgil.CardResponse{
		card1,
		card2,
	}
	crit := &virgil.Criteria{
		Identities:   []string{"bob", "alice"},
		Scope:        virgil.CardScope.Application,
		IdentityType: "test",
	}
	searchKey := owner + "_test_application_alice_bob"
	cache := new(fakeCache)
	cache.On("Get", searchKey).Return(true, []string{"1", "2"})
	cache.On("Get", owner+"_"+card1.ID).Return(true, &card1)
	cache.On("Get", owner+"_"+card2.ID).Return(true, &card2)

	cacheCard := cacheCardMiddleware{cache}
	cards, err := cacheCard.SearchCards(func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
		t.Fatal("Function executed")
		return nil, nil
	})(core.SetOwnerRequest(context.Background(), owner), crit)

	assert.NoError(t, err)
	assert.Equal(t, expected, cards)
}

func TestCacheSearchCards_CacheExist_OneOfCardIsMiss(t *testing.T) {
	owner := "owner"
	card := virgil.CardResponse{ID: "1"}

	crit := &virgil.Criteria{
		Identities:   []string{"bob", "alice"},
		Scope:        virgil.CardScope.Application,
		IdentityType: "test",
	}
	searchKey := owner + "_test_application_alice_bob"
	cache := new(fakeCache)
	cache.On("Get", searchKey).Return(true, []string{"1", "2"})
	cache.On("Get", owner+"_"+card.ID).Return(true, &card)
	cache.On("Get", owner+"_2").Return(false)

	funcExecuted := false
	cacheCard := cacheCardMiddleware{cache}
	cacheCard.SearchCards(func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
		funcExecuted = true
		return nil, fmt.Errorf("ERROR")
	})(core.SetOwnerRequest(context.Background(), owner), crit)

	assert.True(t, funcExecuted)
}

func TestCacheSearchCards_CacheMiss_FuncExecutedReturnErr(t *testing.T) {
	cache := new(fakeCache)
	cache.On("Get", mock.Anything).Return(false)
	cache.On("Set", mock.Anything, mock.Anything)

	cacheCard := cacheCardMiddleware{cache}
	cards, err := cacheCard.SearchCards(func(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
		return nil, fmt.Errorf("ERROR")
	})(context.Background(), &virgil.Criteria{})

	assert.Error(t, err, "ERROR")
	assert.Nil(t, cards)
}

func TestCacheCreateCard_FuncReturnErr_ReturnErr(t *testing.T) {
	cacheCard := cacheCardMiddleware{}
	actual, err := cacheCard.CreateCard(func(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
		return nil, fmt.Errorf("ERROR")
	})(context.Background(), nil)

	assert.Error(t, err, "ERROR")
	assert.Nil(t, actual)
}

func TestCacheCreateCard_FuncReturnVal(t *testing.T) {
	expected := &virgil.CardResponse{
		ID:       "1234",
		Snapshot: []byte(`snapshot`),
	}
	cache := new(fakeCache)
	cache.On("Set", mock.Anything, mock.Anything).Once()

	cacheCard := cacheCardMiddleware{cache}
	actual, err := cacheCard.CreateCard(func(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
		return expected, nil
	})(context.Background(), nil)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestCacheCreateCard_SetKey(t *testing.T) {
	owner := "owner"
	expected := &virgil.CardResponse{
		ID:       "1234",
		Snapshot: []byte(`snapshot`),
	}
	cache := new(fakeCache)
	cache.On("Set", owner+"_"+expected.ID, expected).Once()

	cacheCard := cacheCardMiddleware{cache}
	cacheCard.CreateCard(func(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
		return expected, nil
	})(core.SetOwnerRequest(context.Background(), owner), nil)

	cache.AssertExpectations(t)
}

func TestCacheRevokeCard_FuncReturnErr_ReturnErr(t *testing.T) {
	cacheCard := cacheCardMiddleware{}
	err := cacheCard.RevokeCard(func(ctx context.Context, req *core.RevokeCardRequest) error {
		return fmt.Errorf("ERROR")
	})(context.Background(), nil)

	assert.Error(t, err, "ERROR")
}

func TestCacheRevokeCard_DelKey(t *testing.T) {
	owner := "owner"
	id := "1234"

	cache := new(fakeCache)
	cache.On("Del", owner+"_"+id).Once()

	cacheCard := cacheCardMiddleware{cache}
	cacheCard.RevokeCard(func(ctx context.Context, req *core.RevokeCardRequest) error {
		return nil
	})(core.SetOwnerRequest(context.Background(), owner), &core.RevokeCardRequest{
		Info: virgil.RevokeCardRequest{
			ID: id,
		},
	})

	cache.AssertExpectations(t)
}

func TestCacheCreateRelation_FuncReturnErr_ReturnErr(t *testing.T) {
	cacheCard := cacheCardMiddleware{}
	actual, err := cacheCard.CreateRelations(func(ctx context.Context, req *core.CreateRelationRequest) (*virgil.CardResponse, error) {
		return nil, fmt.Errorf("ERROR")
	})(context.Background(), nil)

	assert.Error(t, err, "ERROR")
	assert.Nil(t, actual)
}

func TestCacheCreateRelation_FuncReturnVal(t *testing.T) {
	expected := &virgil.CardResponse{
		ID:       "1234",
		Snapshot: []byte(`snapshot`),
	}
	cache := new(fakeCache)
	cache.On("Set", mock.Anything, mock.Anything).Once()

	cacheCard := cacheCardMiddleware{cache}
	actual, err := cacheCard.CreateRelations(func(ctx context.Context, req *core.CreateRelationRequest) (*virgil.CardResponse, error) {
		return expected, nil
	})(context.Background(), nil)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestCacheCreateRelation_SetKey(t *testing.T) {
	owner := "owner"
	expected := &virgil.CardResponse{
		ID:       "1234",
		Snapshot: []byte(`snapshot`),
	}
	cache := new(fakeCache)
	cache.On("Set", owner+"_"+expected.ID, expected).Once()

	cacheCard := cacheCardMiddleware{cache}
	cacheCard.CreateRelations(func(ctx context.Context, req *core.CreateRelationRequest) (*virgil.CardResponse, error) {
		return expected, nil
	})(core.SetOwnerRequest(context.Background(), owner), nil)

	cache.AssertExpectations(t)
}

func TestCacheRevokeRelation_FuncReturnErr_ReturnErr(t *testing.T) {
	cacheCard := cacheCardMiddleware{}
	actual, err := cacheCard.RevokeRelations(func(ctx context.Context, req *core.RevokeRelationRequest) (*virgil.CardResponse, error) {
		return nil, fmt.Errorf("ERROR")
	})(context.Background(), nil)

	assert.Error(t, err, "ERROR")
	assert.Nil(t, actual)
}

func TestCacheRevokeRelation_FuncReturnVal(t *testing.T) {
	expected := &virgil.CardResponse{
		ID:       "1234",
		Snapshot: []byte(`snapshot`),
	}
	cache := new(fakeCache)
	cache.On("Set", mock.Anything, mock.Anything).Once()

	cacheCard := cacheCardMiddleware{cache}
	actual, err := cacheCard.RevokeRelations(func(ctx context.Context, req *core.RevokeRelationRequest) (*virgil.CardResponse, error) {
		return expected, nil
	})(context.Background(), nil)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestCacheRevokeRelation_SetKey(t *testing.T) {
	owner := "owner"
	expected := &virgil.CardResponse{
		ID:       "1234",
		Snapshot: []byte(`snapshot`),
	}
	cache := new(fakeCache)
	cache.On("Set", owner+"_"+expected.ID, expected).Once()

	cacheCard := cacheCardMiddleware{cache}
	cacheCard.RevokeRelations(func(ctx context.Context, req *core.RevokeRelationRequest) (*virgil.CardResponse, error) {
		return expected, nil
	})(core.SetOwnerRequest(context.Background(), owner), nil)

	cache.AssertExpectations(t)
}
