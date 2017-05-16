package middleware

import (
	"net/http"
	"testing"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fakeAppStore struct {
	mock.Mock
}

func (f *fakeAppStore) GetById(id string) (app *core.Application, err error) {
	args := f.Called(id)
	app, _ = args.Get(0).(*core.Application)
	err = args.Error(1)
	return
}

type fakeTokenStore struct {
	mock.Mock
}

func (f *fakeTokenStore) GetByValue(val string) (token *core.Token, err error) {
	args := f.Called(val)
	token, _ = args.Get(0).(*core.Token)
	err = args.Error(1)
	return
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

	assert.Nil(t, s, "Seccessfull result is nil")
	assert.Nil(t, err, "Error result is nil")
}
