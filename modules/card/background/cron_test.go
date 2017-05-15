package background

import (
	"fmt"
	"testing"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fakeDevPortalClient struct {
	mock.Mock
}

func (f *fakeDevPortalClient) GetApplications() (apps []core.Application, err error) {
	args := f.Called()
	apps, _ = args.Get(0).([]core.Application)
	err = args.Error(1)

	return
}

func (f *fakeDevPortalClient) GetTokens() (apps []core.Token, err error) {
	args := f.Called()
	apps, _ = args.Get(0).([]core.Token)
	err = args.Error(1)

	return
}

type fakeAppStore struct {
	mock.Mock
}

func (f *fakeAppStore) Add(app core.Application) error {
	args := f.Called(app)
	return args.Error(0)
}

func (f *fakeAppStore) Delete(id string) error {
	args := f.Called(id)
	return args.Error(0)
}

func (f *fakeAppStore) GetAll() (apps []core.Application, err error) {
	args := f.Called()
	apps, _ = args.Get(0).([]core.Application)
	err = args.Error(1)
	return
}

type fakeTokenStore struct {
	mock.Mock
}

func (f *fakeTokenStore) Add(token core.Token) error {
	args := f.Called(token)
	return args.Error(0)
}
func (f *fakeTokenStore) Delete(id string) error {
	args := f.Called(id)
	return args.Error(0)
}
func (f *fakeTokenStore) GetAll() (tokens []core.Token, err error) {
	args := f.Called()
	tokens, _ = args.Get(0).([]core.Token)
	err = args.Error(1)
	return
}

func TestUpdateAppsCronJob_ClientReturnErr_ReturnErr(t *testing.T) {
	c := new(fakeDevPortalClient)
	c.On("GetApplications").Return(nil, fmt.Errorf("Error"))
	s := new(fakeAppStore)
	job := UpdateAppsCronJob(s, c)
	err := job()

	assert.Error(t, err, "Error")
}

func TestUpdateAppsCronJob_StoreReturnErr_ReturnErr(t *testing.T) {
	c := new(fakeDevPortalClient)
	c.On("GetApplications").Return([]core.Application{}, nil)
	s := new(fakeAppStore)
	s.On("GetAll").Return(nil, fmt.Errorf("Error"))
	job := UpdateAppsCronJob(s, c)
	err := job()

	assert.Error(t, err, "Error")
}

func TestUpdateAppsCronJob_AppsWasNotChanged_DoNothing(t *testing.T) {
	app := core.Application{
		ID:        "test",
		UpdatedAt: "today",
	}
	c := new(fakeDevPortalClient)
	c.On("GetApplications").Return([]core.Application{app}, nil)
	s := new(fakeAppStore)
	s.On("GetAll").Return([]core.Application{app}, nil)
	job := UpdateAppsCronJob(s, c)
	job()

	s.AssertNotCalled(t, "Delete", mock.Anything)
	s.AssertNotCalled(t, "Add", mock.Anything)
}

func TestUpdateAppsCronJob_AppDeletedStoreDeleteReturnErr_ReturnErr(t *testing.T) {

	c := new(fakeDevPortalClient)
	c.On("GetApplications").Return([]core.Application{}, nil)
	s := new(fakeAppStore)
	s.On("GetAll").Return([]core.Application{core.Application{
		ID:        "test",
		UpdatedAt: "today",
	}}, nil)
	s.On("Delete", mock.Anything).Return(fmt.Errorf("Error"))
	job := UpdateAppsCronJob(s, c)
	err := job()

	assert.Error(t, err, "Error")
}

func TestUpdateAppsCronJob_AppDeleted_StoreDelete(t *testing.T) {
	expected := core.Application{
		ID:        "test",
		UpdatedAt: "today",
	}

	c := new(fakeDevPortalClient)
	c.On("GetApplications").Return([]core.Application{}, nil)
	s := new(fakeAppStore)
	s.On("GetAll").Return([]core.Application{expected}, nil)
	s.On("Delete", expected.ID).Return(nil).Once()
	job := UpdateAppsCronJob(s, c)
	job()

	s.AssertExpectations(t)
}

func TestUpdateAppsCronJob_AppUpdated_StoreDelete(t *testing.T) {
	c := new(fakeDevPortalClient)
	c.On("GetApplications").Return([]core.Application{core.Application{
		ID:        "test",
		UpdatedAt: "today",
	}}, nil)
	s := new(fakeAppStore)
	s.On("GetAll").Return([]core.Application{core.Application{
		ID:        "test",
		UpdatedAt: "yesterday",
	}}, nil)
	s.On("Delete", "test").Return(nil).Once()
	s.On("Add", mock.Anything).Return(nil)
	job := UpdateAppsCronJob(s, c)
	job()

	s.AssertExpectations(t)
}

func TestUpdateAppsCronJob_AppAddedStoreAddReturnErr_ReturnErr(t *testing.T) {

	c := new(fakeDevPortalClient)
	c.On("GetApplications").Return([]core.Application{core.Application{
		ID:        "test",
		UpdatedAt: "today",
	}}, nil)
	s := new(fakeAppStore)
	s.On("GetAll").Return([]core.Application{}, nil)
	s.On("Add", mock.Anything).Return(fmt.Errorf("Error"))
	job := UpdateAppsCronJob(s, c)
	err := job()

	assert.Error(t, err, "Error")
}

func TestUpdateAppsCronJob_AppAdded_StoreAdd(t *testing.T) {
	expected := core.Application{
		ID:        "test",
		UpdatedAt: "today",
	}

	c := new(fakeDevPortalClient)
	c.On("GetApplications").Return([]core.Application{expected}, nil)
	s := new(fakeAppStore)
	s.On("GetAll").Return([]core.Application{}, nil)
	s.On("Add", expected).Return(nil).Once()
	job := UpdateAppsCronJob(s, c)
	job()

	s.AssertExpectations(t)
}

func TestUpdateAppsCronJob_AppUpdated_StoreAdd(t *testing.T) {
	expected := core.Application{
		ID:        "test",
		UpdatedAt: "today",
	}
	c := new(fakeDevPortalClient)
	c.On("GetApplications").Return([]core.Application{expected}, nil)
	s := new(fakeAppStore)
	s.On("GetAll").Return([]core.Application{core.Application{
		ID:        "test",
		UpdatedAt: "yesterday",
	}}, nil)
	s.On("Delete", mock.Anything).Return(nil)
	s.On("Add", expected).Return(nil).Once()
	job := UpdateAppsCronJob(s, c)
	job()

	s.AssertExpectations(t)
}

func TestUpdateTokensCronJob_ClientReturnErr_ReturnErr(t *testing.T) {
	c := new(fakeDevPortalClient)
	c.On("GetTokens").Return(nil, fmt.Errorf("Error"))
	s := new(fakeTokenStore)
	job := UpdateTokensCronJob(s, c)
	err := job()

	assert.Error(t, err, "Error")
}

func TestUpdateTokensCronJob_StoreReturnErr_ReturnErr(t *testing.T) {
	c := new(fakeDevPortalClient)
	c.On("GetTokens").Return([]core.Token{}, nil)
	s := new(fakeTokenStore)
	s.On("GetAll").Return(nil, fmt.Errorf("Error"))
	job := UpdateTokensCronJob(s, c)
	err := job()

	assert.Error(t, err, "Error")
}

func TestUpdateTokensCronJob_TokensWasNotChanged_DoNothing(t *testing.T) {
	token := core.Token{
		ID:        "test",
		UpdatedAt: "today",
	}
	c := new(fakeDevPortalClient)
	c.On("GetTokens").Return([]core.Token{token}, nil)
	s := new(fakeTokenStore)
	s.On("GetAll").Return([]core.Token{token}, nil)
	job := UpdateTokensCronJob(s, c)
	job()

	s.AssertNotCalled(t, "Delete", mock.Anything)
	s.AssertNotCalled(t, "Add", mock.Anything)
}

func TestUpdateTokensCronJob_TokenDeletedStoreDeleteReturnErr_ReturnErr(t *testing.T) {

	c := new(fakeDevPortalClient)
	c.On("GetTokens").Return([]core.Token{}, nil)
	s := new(fakeTokenStore)
	s.On("GetAll").Return([]core.Token{core.Token{
		ID:        "test",
		UpdatedAt: "today",
	}}, nil)
	s.On("Delete", mock.Anything).Return(fmt.Errorf("Error"))
	job := UpdateTokensCronJob(s, c)
	err := job()

	assert.Error(t, err, "Error")
}

func TestUpdateTokensCronJob_TokenDeleted_StoreDelete(t *testing.T) {
	expected := core.Token{
		ID:        "test",
		UpdatedAt: "today",
	}

	c := new(fakeDevPortalClient)
	c.On("GetTokens").Return([]core.Token{}, nil)
	s := new(fakeTokenStore)
	s.On("GetAll").Return([]core.Token{expected}, nil)
	s.On("Delete", expected.ID).Return(nil).Once()
	job := UpdateTokensCronJob(s, c)
	job()

	s.AssertExpectations(t)
}

func TestUpdateTokensCronJob_TokenUpdated_StoreDelete(t *testing.T) {
	c := new(fakeDevPortalClient)
	c.On("GetTokens").Return([]core.Token{core.Token{
		ID:        "test",
		UpdatedAt: "today",
	}}, nil)
	s := new(fakeTokenStore)
	s.On("GetAll").Return([]core.Token{core.Token{
		ID:        "test",
		UpdatedAt: "yesterday",
	}}, nil)
	s.On("Delete", "test").Return(nil).Once()
	s.On("Add", mock.Anything).Return(nil)
	job := UpdateTokensCronJob(s, c)
	job()

	s.AssertExpectations(t)
}

func TestUpdateTokensCronJob_TokenAddedStoreAddReturnErr_ReturnErr(t *testing.T) {

	c := new(fakeDevPortalClient)
	c.On("GetTokens").Return([]core.Token{core.Token{
		ID:        "test",
		UpdatedAt: "today",
	}}, nil)
	s := new(fakeTokenStore)
	s.On("GetAll").Return([]core.Token{}, nil)
	s.On("Add", mock.Anything).Return(fmt.Errorf("Error"))
	job := UpdateTokensCronJob(s, c)
	err := job()

	assert.Error(t, err, "Error")
}

func TestUpdateTokensCronJob_TokenAdded_StoreAdd(t *testing.T) {
	expected := core.Token{
		ID:        "test",
		UpdatedAt: "today",
	}

	c := new(fakeDevPortalClient)
	c.On("GetTokens").Return([]core.Token{expected}, nil)
	s := new(fakeTokenStore)
	s.On("GetAll").Return([]core.Token{}, nil)
	s.On("Add", expected).Return(nil).Once()
	job := UpdateTokensCronJob(s, c)
	job()

	s.AssertExpectations(t)
}

func TestUpdateTokensCronJob_TokenUpdated_StoreAdd(t *testing.T) {
	expected := core.Token{
		ID:        "test",
		UpdatedAt: "today",
	}
	c := new(fakeDevPortalClient)
	c.On("GetTokens").Return([]core.Token{expected}, nil)
	s := new(fakeTokenStore)
	s.On("GetAll").Return([]core.Token{core.Token{
		ID:        "test",
		UpdatedAt: "yesterday",
	}}, nil)
	s.On("Delete", mock.Anything).Return(nil)
	s.On("Add", expected).Return(nil).Once()
	job := UpdateTokensCronJob(s, c)
	job()

	s.AssertExpectations(t)
}
