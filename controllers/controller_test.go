package controllers

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/protocols"
	"testing"
)

type MockStorage struct {
	mock.Mock
}

func (s MockStorage) GetCard(id string) (*models.CardResponse, *models.ErrorResponse) {
	args := s.Called(id)
	var (
		cr *models.CardResponse
		er *models.ErrorResponse
		ok bool
	)
	if cr, ok = args.Get(0).(*models.CardResponse); !ok {
		cr = nil
	}
	if er, ok = args.Get(1).(*models.ErrorResponse); !ok {
		er = nil
	}
	return cr, er
}
func (s MockStorage) SearchCards(c models.Criteria) ([]models.CardResponse, *models.ErrorResponse) {
	args := s.Called(c)
	var (
		cr []models.CardResponse
		er *models.ErrorResponse
		ok bool
	)
	if cr, ok = args.Get(0).([]models.CardResponse); !ok {
		cr = nil
	}
	if er, ok = args.Get(1).(*models.ErrorResponse); !ok {
		er = nil
	}
	return cr, er
}

func (s MockStorage) CreateCard(c *models.CardResponse) (*models.CardResponse, *models.ErrorResponse) {
	args := s.Called(c)
	var (
		cr *models.CardResponse
		er *models.ErrorResponse
		ok bool
	)
	if cr, ok = args.Get(0).(*models.CardResponse); !ok {
		cr = nil
	}
	if er, ok = args.Get(1).(*models.ErrorResponse); !ok {
		er = nil
	}
	return cr, er
}

func (s MockStorage) RevokeCard(id string, c *models.CardResponse) *models.ErrorResponse {
	args := s.Called(id, c)
	var (
		er *models.ErrorResponse
		ok bool
	)
	if er, ok = args.Get(0).(*models.ErrorResponse); !ok {
		er = nil
	}
	return er
}

func MakeFakeCardResponse() *models.CardResponse {
	return MakeFakeCardResponseWith("test")
}

func MakeFakeCardResponseWith(text string) *models.CardResponse {
	return &models.CardResponse{
		Snapshot: []byte(text),
		Meta: models.ResponseMeta{
			CreatedAt:   text,
			CardVersion: "v4",
			Signatures: map[string][]byte{
				text: []byte(text),
			},
		},
	}
}

func AssertControllerRespose(t *testing.T, expected *models.CardResponse, r []byte, e protocols.CodeResponse) {
	assert.Equal(t, protocols.Ok, e)

	actial := new(models.CardResponse)
	assert.Nil(t, json.Unmarshal(r, actial), "Cannot restore object from response")
	assert.Equal(t, expected, actial)
}

func Test_GetCard_StorageReturnErr_ReturnErr(t *testing.T) {
	id := "test"
	testTable := map[protocols.CodeResponse]*models.ErrorResponse{
		protocols.ServerError:  models.MakeError(10000),
		protocols.RequestError: models.MakeError(12344), // any other code
	}

	for k, v := range testTable {
		actual := new(models.ErrorResponse)

		mStorage := MockStorage{}
		mStorage.On("GetCard", id).Return(nil, v)

		c := Controller{
			Storage: mStorage,
		}
		data, code := c.GetCard(id)

		assert.Equal(t, k, code)
		assert.Nil(t, json.Unmarshal(data, actual), "Cannot restore object from response")
		assert.Equal(t, v, actual)
	}

}

func Test_GetCard_StorageReturnVal_ReturnJsonByte(t *testing.T) {
	id := "test"
	expected := MakeFakeCardResponse()
	mStorage := MockStorage{}
	mStorage.On("GetCard", id).Return(expected, nil)

	c := Controller{
		Storage: mStorage,
	}
	r, code := c.GetCard(id)
	AssertControllerRespose(t, expected, r, code)
}

func Test_GetCard_StorageReturnNilValue_ReturnNilByte(t *testing.T) {
	id := "test"
	mStorage := MockStorage{}
	mStorage.On("GetCard", id).Return(nil, nil)

	c := Controller{
		Storage: mStorage,
	}
	r, code := c.GetCard(id)
	assert.Equal(t, protocols.NotFound, code)
	assert.Nil(t, r)
}

func Test_SearchCards_BrokenRequestData_ReturnErr(t *testing.T) {
	expected := models.MakeError(30000)
	actual := new(models.ErrorResponse)
	mStorage := MockStorage{}
	c := Controller{
		Storage: mStorage,
	}
	data, code := c.SearchCards([]byte("Test"))

	assert.Equal(t, protocols.RequestError, code)
	assert.Nil(t, json.Unmarshal(data, actual), "Cannot restore object from response")
	assert.Equal(t, expected, actual)
}

func Test_SearchCards_StorageReturnErr_ReturnErr(t *testing.T) {
	criteria := models.Criteria{
		Scope: "global",
	}
	data, _ := json.Marshal(&criteria)

	testTable := map[protocols.CodeResponse]*models.ErrorResponse{
		protocols.ServerError:  models.MakeError(10000),
		protocols.RequestError: models.MakeError(12344), // any other code
	}

	for k, v := range testTable {
		actual := new(models.ErrorResponse)

		mStorage := MockStorage{}
		mStorage.On("SearchCards", criteria).Return([]models.CardResponse{}, v)

		c := Controller{
			Storage: mStorage,
		}
		data, code := c.SearchCards(data)

		assert.Equal(t, k, code)
		assert.Nil(t, json.Unmarshal(data, actual), "Cannot restore object from response")
		assert.Equal(t, v, actual)
	}
}

func Test_SearchCards_StorageReturnVal_ReturnJsonByte(t *testing.T) {
	criteria, _ := json.Marshal(models.Criteria{
		IdentityType: "test",
		Identities: []string{
			"test1",
			"test2",
		},
		Scope: "test",
	})

	restoredCriteria := models.Criteria{
		IdentityType: "test",
		Identities: []string{
			"test1",
			"test2",
		},
		Scope: "application",
	}

	expected := []models.CardResponse{
		*MakeFakeCardResponse(),
		*MakeFakeCardResponse(),
	}
	mStorage := MockStorage{}
	mStorage.On("SearchCards", restoredCriteria).Return(expected, nil)

	c := Controller{
		Storage: mStorage,
	}
	r, code := c.SearchCards(criteria)

	assert.Equal(t, protocols.Ok, code)

	var actial []models.CardResponse
	assert.Nil(t, json.Unmarshal(r, &actial), "Cannot restore object from response")
	assert.Equal(t, expected, actial)
}

func Test_CreateCard_BrokenRequestData_ReturnErr(t *testing.T) {
	expected := models.MakeError(30000)
	actual := new(models.ErrorResponse)
	mStorage := MockStorage{}
	c := Controller{
		Storage: mStorage,
	}
	data, code := c.CreateCard([]byte("Test"))

	assert.Equal(t, protocols.RequestError, code)
	assert.Nil(t, json.Unmarshal(data, actual), "Cannot restore object from response")
	assert.Equal(t, expected, actual)
}

func Test_CreateCard_StorageReturnErr_ReturnErr(t *testing.T) {
	param := MakeFakeCardResponse()
	data, _ := json.Marshal(param)

	testTable := map[protocols.CodeResponse]*models.ErrorResponse{
		protocols.ServerError:  models.MakeError(10000),
		protocols.RequestError: models.MakeError(12344), // any other code
	}

	for k, v := range testTable {
		actual := new(models.ErrorResponse)

		mStorage := MockStorage{}
		mStorage.On("CreateCard", param).Return(nil, v)

		c := Controller{
			Storage: mStorage,
		}
		data, code := c.CreateCard(data)

		assert.Equal(t, k, code)
		assert.Nil(t, json.Unmarshal(data, actual), "Cannot restore object from response")
		assert.Equal(t, v, actual)
	}
}

func Test_CreateCard_StorageReturnVal_ReturnJsonByte(t *testing.T) {
	param := MakeFakeCardResponse()
	data, _ := json.Marshal(param)

	expected := MakeFakeCardResponseWith("expected")
	mStorage := MockStorage{}
	mStorage.On("CreateCard", param).Return(expected, nil)

	c := Controller{
		Storage: mStorage,
	}
	r, code := c.CreateCard(data)

	AssertControllerRespose(t, expected, r, code)
}

func Test_RevokeCard_BrokenRequestData_ReturnErr(t *testing.T) {
	id := "test"
	expected := models.MakeError(30000)
	actual := new(models.ErrorResponse)
	mStorage := MockStorage{}
	c := Controller{
		Storage: mStorage,
	}
	data, code := c.RevokeCard(id, []byte("Test"))

	assert.Equal(t, protocols.RequestError, code)
	assert.Nil(t, json.Unmarshal(data, actual), "Cannot restore object from response")
	assert.Equal(t, expected, actual)
}

func Test_RevokeCard_StorageReturnErr_ReturnErr(t *testing.T) {
	id := "test"
	param := MakeFakeCardResponse()
	data, _ := json.Marshal(param)

	testTable := map[protocols.CodeResponse]*models.ErrorResponse{
		protocols.ServerError:  models.MakeError(10000),
		protocols.RequestError: models.MakeError(12344), // any other code
	}

	for k, v := range testTable {
		actual := new(models.ErrorResponse)

		mStorage := MockStorage{}
		mStorage.On("RevokeCard", id, param).Return(v)

		c := Controller{
			Storage: mStorage,
		}
		data, code := c.RevokeCard(id, data)

		assert.Equal(t, k, code)
		assert.Nil(t, json.Unmarshal(data, actual), "Cannot restore object from response")
		assert.Equal(t, v, actual)
	}
}

func Test_RevokeCard_StorageReturnNilErr_ReturnNil(t *testing.T) {
	id := "test"
	param := MakeFakeCardResponse()
	data, _ := json.Marshal(&param)

	mStorage := MockStorage{}
	mStorage.On("RevokeCard", id, param).Return(nil)

	c := Controller{
		Storage: mStorage,
	}
	data, code := c.RevokeCard(id, data)
	assert.Equal(t, protocols.Ok, code)
	assert.Nil(t, data)
}
