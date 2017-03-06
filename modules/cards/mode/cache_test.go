package mode

import (
	"fmt"
	"testing"

	virgil "gopkg.in/virgil.v4"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fakeCacheManager struct {
	mock.Mock
}

func (f *fakeCacheManager) Get(key string, val interface{}) bool {
	args := f.Called(key)
	has := args.Bool(0)
	if has {

		if c, ok := val.(*core.Card); ok {
			setcard := args.Get(1).(*core.Card)
			*c = *setcard
		} else if c, ok := val.(**core.Card); ok {
			*c = args.Get(1).(*core.Card)
		} else if ids, ok := val.(*[]string); ok {
			setids := args.Get(1).([]string)
			*ids = setids
		}
	}
	return has
}
func (f *fakeCacheManager) Set(key string, val interface{}) {
	f.Called(key, val)
}
func (f *fakeCacheManager) Del(key string) {
	f.Called(key)
}

func TestCacheGet_CacheExist_ReturnVal(t *testing.T) {
	m := new(fakeCacheManager)
	const id = "1234"
	expected := &core.Card{ID: "1234", Snapshot: []byte("asdf")}
	m.On("Get", id, mock.Anything).Return(true, expected)
	ccm := CacheCardsMiddleware{m}
	f := ccm.Get(func(id string) (*core.Card, error) {
		t.FailNow()
		return nil, nil
	})

	actual, err := f(id)
	assert.Exactly(t, expected, actual, "Cards are not equals")
	assert.NoError(t, err)
}

func TestCacheGet_MissCache_ReturnNextVal(t *testing.T) {
	m := new(fakeCacheManager)
	m.On("Get", mock.Anything, mock.Anything).Return(false)
	m.On("Set", mock.Anything, mock.Anything)

	const expectedId = "1234"
	expected := &core.Card{ID: "1234", Snapshot: []byte("asdf")}
	ccm := CacheCardsMiddleware{m}
	f := ccm.Get(func(id string) (*core.Card, error) {
		assert.Equal(t, expectedId, id)
		return expected, nil
	})

	actual, err := f(expectedId)
	assert.Exactly(t, expected, actual, "Cards are not equals")
	assert.NoError(t, err)
}

func TestCacheGet_MissCache_UpdateCache(t *testing.T) {
	const id = "1234"
	expected := &core.Card{ID: "1234", Snapshot: []byte("asdf")}

	m := new(fakeCacheManager)
	m.On("Get", mock.Anything, mock.Anything).Return(false)
	m.On("Set", id, expected).Once()

	ccm := CacheCardsMiddleware{m}
	f := ccm.Get(func(id string) (*core.Card, error) {
		return expected, nil
	})

	f(id)

	m.AssertExpectations(t)
}

func TestCacheGet_MissCacheNextReturnErr_ReturnErr(t *testing.T) {
	const id = "1234"

	m := new(fakeCacheManager)
	m.On("Get", mock.Anything, mock.Anything).Return(false)

	ccm := CacheCardsMiddleware{m}
	f := ccm.Get(func(id string) (*core.Card, error) {
		return nil, fmt.Errorf("Error")
	})

	_, err := f(id)

	assert.Error(t, err)
}

func TestCacheSearch_CacheExist_ReturnVal(t *testing.T) {
	crit := virgil.Criteria{
		Identities:   []string{"bob"},
		Scope:        virgil.CardScope.Global,
		IdentityType: "nick",
	}
	const id = "1234"
	expected := core.Card{ID: "1234", Snapshot: []byte("asdf")}

	m := new(fakeCacheManager)
	m.On("Get", id, mock.Anything).Return(true, &expected)
	m.On("Get", fmt.Sprint(crit.IdentityType, crit.Scope, "bob"), mock.Anything).Return(true, []string{id})

	ccm := CacheCardsMiddleware{m}
	f := ccm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		t.FailNow()
		return nil, nil
	})

	actual, err := f(&crit)
	assert.Equal(t, []core.Card{expected}, actual, "Cards are not equals")
	assert.NoError(t, err)
}

func TestCacheSearch_CachePartiallyExist_UpdateCache(t *testing.T) {
	const id = "1234"
	expected := core.Card{ID: id, Snapshot: []byte("asdf")}
	crit := virgil.Criteria{
		Identities:   []string{"bob"},
		Scope:        virgil.CardScope.Global,
		IdentityType: "nick",
	}

	m := new(fakeCacheManager)
	m.On("Get", id, mock.Anything).Return(false)
	m.On("Get", fmt.Sprint(crit.IdentityType, crit.Scope, "bob"), mock.Anything).Return(true, []string{id})

	m.On("Set", id, &expected).Once()
	m.On("Set", fmt.Sprint(crit.IdentityType, crit.Scope, "bob"), []string{id}).Once()

	ccm := CacheCardsMiddleware{m}
	f := ccm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		return []core.Card{expected}, nil
	})

	f(&crit)

	m.AssertExpectations(t)
}

func TestCacheSearch_MissCache_ReturnNextVal(t *testing.T) {
	m := new(fakeCacheManager)
	m.On("Get", mock.Anything, mock.Anything).Return(false)
	m.On("Set", mock.Anything, mock.Anything)

	expected := []core.Card{core.Card{ID: "1234", Snapshot: []byte("asdf")}}
	ccm := CacheCardsMiddleware{m}
	f := ccm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		return expected, nil
	})

	actual, err := f(&virgil.Criteria{})
	assert.Exactly(t, expected, actual, "Cards are not equals")
	assert.NoError(t, err)
}

func TestCacheSearch_MissCache_UpdateCache(t *testing.T) {
	const id = "1234"
	expected := core.Card{ID: id, Snapshot: []byte("asdf")}
	crit := virgil.Criteria{
		Identities:   []string{"bob", "alice"},
		Scope:        virgil.CardScope.Global,
		IdentityType: "nick",
	}

	m := new(fakeCacheManager)
	m.On("Get", mock.Anything, mock.Anything).Return(false)
	m.On("Set", id, &expected).Once()
	m.On("Set", fmt.Sprint(crit.IdentityType, crit.Scope, "alice_bob"), []string{id}).Once()

	ccm := CacheCardsMiddleware{m}
	f := ccm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		return []core.Card{expected}, nil
	})

	f(&crit)

	m.AssertExpectations(t)
}

func TestCacheSearch_MissCacheNextReturnErr_ReturnErr(t *testing.T) {
	m := new(fakeCacheManager)
	m.On("Get", mock.Anything, mock.Anything).Return(false)

	ccm := CacheCardsMiddleware{m}
	f := ccm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		return nil, fmt.Errorf("Error")
	})

	_, err := f(&virgil.Criteria{})

	assert.Error(t, err)
}

func TestCacheCreate_NextReturnErr_ReturnErr(t *testing.T) {
	ccm := CacheCardsMiddleware{}
	f := ccm.Create(func(req *core.CreateCardRequest) (*core.Card, error) {
		return nil, fmt.Errorf("Error")
	})

	_, err := f(&core.CreateCardRequest{})

	assert.Error(t, err)
}

func TestCacheCreate_NextReturnVal_ReturnVal(t *testing.T) {
	expected := &core.Card{ID: "1234", Snapshot: []byte("asdf")}

	m := new(fakeCacheManager)
	m.On("Set", mock.Anything, mock.Anything)

	ccm := CacheCardsMiddleware{m}
	f := ccm.Create(func(req *core.CreateCardRequest) (*core.Card, error) {
		return expected, nil
	})

	actual, err := f(&core.CreateCardRequest{})

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestCacheCreate_CacheVal(t *testing.T) {
	expected := &core.Card{ID: "1234", Snapshot: []byte("asdf")}

	m := new(fakeCacheManager)
	m.On("Set", expected.ID, expected).Once()

	ccm := CacheCardsMiddleware{m}
	f := ccm.Create(func(req *core.CreateCardRequest) (*core.Card, error) {
		return expected, nil
	})

	f(&core.CreateCardRequest{})

	m.AssertExpectations(t)
}

func TestCacheRevoke_NextReturnErr_ReturnErr(t *testing.T) {
	ccm := CacheCardsMiddleware{}
	f := ccm.Revoke(func(req *core.RevokeCardRequest) error {
		return fmt.Errorf("Error")
	})

	err := f(&core.RevokeCardRequest{})

	assert.Error(t, err)
}

func TestCacheRevoke_NextReturnNil_ReturnNil(t *testing.T) {
	m := new(fakeCacheManager)
	m.On("Del", mock.Anything)

	ccm := CacheCardsMiddleware{m}
	f := ccm.Revoke(func(req *core.RevokeCardRequest) error {
		return nil
	})

	err := f(&core.RevokeCardRequest{})

	assert.NoError(t, err)
}

func TestCacheRevoke_RemoveCache(t *testing.T) {
	const id = "1234"

	m := new(fakeCacheManager)
	m.On("Del", id).Once()

	ccm := CacheCardsMiddleware{m}
	f := ccm.Revoke(func(req *core.RevokeCardRequest) error {
		return nil
	})

	f(&core.RevokeCardRequest{Info: virgil.RevokeCardRequest{ID: id}})

	m.AssertExpectations(t)
}
