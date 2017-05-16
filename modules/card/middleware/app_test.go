package middleware

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fakeTokenStore struct {
	mock.Mock
}

func (f *fakeTokenStore) GetByValue(val string) (token *core.Token, err error) {
	args := f.Called(val)
	token, _ = args.Get(0).(*core.Token)
	err = args.Error(1)
	return
}

type fakeCache struct {
	mock.Mock
}

func (f *fakeCache) Get(key string, val interface{}) bool {
	args := f.Called(key)
	has := args.Bool(0)
	if has {
		if appID, ok := val.(*string); ok {
			*appID = args.Get(1).(string)
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

func TestAppMiddlewareRequestApp_OwnerEmpty_FuncExecute(t *testing.T) {
	funcExecuted := false
	am := AppMiddleware{}
	m := am.RequestApp(func(req *http.Request) (interface{}, error) {
		funcExecuted = true
		return nil, nil
	})
	req, _ := http.NewRequest(http.MethodPost, "http://localhost/action", nil)
	s, err := m(req)

	assert.Nil(t, s, "Seccessfull result isn't nil")
	assert.Nil(t, err, "Error result isn't nil")
}

func TestAppMiddlewareRequestApp_AppIDCachePass_FuncExecute(t *testing.T) {
	fCache := new(fakeCache)
	fCache.On("Get", "owner").Return(true, "appID")
	am := AppMiddleware{Cache: fCache}
	m := am.RequestApp(func(req *http.Request) (interface{}, error) {
		owner := core.GetOwnerRequest(req.Context())

		assert.Equal(t, "appID", owner, "Owner was not updated")

		return nil, nil
	})
	req, _ := http.NewRequest(http.MethodPost, "http://localhost/action", nil)
	req = req.WithContext(core.SetOwnerRequest(context.Background(), "owner"))
	s, err := m(req)

	assert.Nil(t, s, "Seccessfull result is not nil")
	assert.Nil(t, err, "Error result is not nil")
}

func TestAppMiddlewareRequestApp_AppIDCacheMissTokenStoreReturnErr_ReturnErr(t *testing.T) {
	fCache := new(fakeCache)
	fCache.On("Get", "owner").Return(false, "")
	ts := new(fakeTokenStore)
	ts.On("GetByValue", mock.Anything).Return(nil, fmt.Errorf("ERROR"))

	am := AppMiddleware{Cache: fCache, TokenStore: ts}
	m := am.RequestApp(func(req *http.Request) (interface{}, error) {
		assert.FailNow(t, "Function was executed")

		return nil, nil
	})

	req, _ := http.NewRequest(http.MethodPost, "http://localhost/action", nil)
	req = req.WithContext(core.SetOwnerRequest(context.Background(), "owner"))
	s, err := m(req)

	assert.Nil(t, s, "Seccessfull result is not nil")
	assert.Error(t, err)
}

func TestAppMiddlewareRequestApp_AppIDCacheMissTokenStoreEntityNotFound_ReturnErr(t *testing.T) {
	funcExecuted := false
	expectedOwner := "owner"

	fCache := new(fakeCache)
	fCache.On("Get", "owner").Return(false, "")
	ts := new(fakeTokenStore)
	ts.On("GetByValue", mock.Anything).Return(nil, core.EntityNotFoundErr)

	am := AppMiddleware{Cache: fCache, TokenStore: ts}
	m := am.RequestApp(func(req *http.Request) (interface{}, error) {
		funcExecuted = true

		owner := core.GetOwnerRequest(req.Context())
		assert.Equal(t, expectedOwner, owner)
		return nil, nil
	})

	req, _ := http.NewRequest(http.MethodPost, "http://localhost/action", nil)
	req = req.WithContext(core.SetOwnerRequest(context.Background(), expectedOwner))
	s, err := m(req)

	assert.Nil(t, s, "Seccessfull result is not nil")
	assert.Nil(t, err, "Error result is not nil")
	assert.True(t, funcExecuted)
}

func TestAppMiddlewareRequestApp_AppIDFound_SetCache(t *testing.T) {
	appID := "appID"
	owner := "owner"

	fCache := new(fakeCache)
	fCache.On("Get", owner).Return(false, "")
	fCache.On("Set", owner, appID).Return(nil).Once()

	ts := new(fakeTokenStore)
	ts.On("GetByValue", mock.Anything).Return(&core.Token{Application: appID}, nil)

	am := AppMiddleware{Cache: fCache, TokenStore: ts}
	m := am.RequestApp(func(req *http.Request) (interface{}, error) {
		return nil, nil
	})

	req, _ := http.NewRequest(http.MethodPost, "http://localhost/action", nil)
	req = req.WithContext(core.SetOwnerRequest(context.Background(), owner))
	m(req)

	fCache.AssertExpectations(t)
}

func TestAppMiddlewareRequestApp_AppIDFound_UpdateOwnerContext(t *testing.T) {
	funcExecuted := false
	appID := "appID"
	owner := "owner"

	fCache := new(fakeCache)
	fCache.On("Get", owner).Return(false, "")
	fCache.On("Set", owner, appID).Return(nil).Once()

	ts := new(fakeTokenStore)
	ts.On("GetByValue", owner).Return(&core.Token{Application: appID}, nil)

	am := AppMiddleware{Cache: fCache, TokenStore: ts}
	m := am.RequestApp(func(req *http.Request) (interface{}, error) {
		funcExecuted = true

		actualOwner := core.GetOwnerRequest(req.Context())
		assert.Equal(t, appID, actualOwner)
		return nil, nil
	})

	req, _ := http.NewRequest(http.MethodPost, "http://localhost/action", nil)
	req = req.WithContext(core.SetOwnerRequest(context.Background(), owner))
	m(req)
}
