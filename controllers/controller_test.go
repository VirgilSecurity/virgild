package controllers

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
	"testing"
)

type MockStorage struct {
	mock.Mock
}

func (s MockStorage) GetCard(id string) (*models.CardResponse, error) {
	args := s.Called(id)
	v, ok := args.Get(0).(*models.CardResponse)
	if ok {
		return v, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}
func (s MockStorage) SearchCards(c models.Criteria) ([]models.CardResponse, error) {
	args := s.Called(c)
	return args.Get(0).([]models.CardResponse), args.Error(1)
}

func (s MockStorage) CreateCard(c models.CardResponse) (*models.CardResponse, error) {
	args := s.Called(c)
	v, ok := args.Get(0).(*models.CardResponse)
	if ok {
		return v, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (s MockStorage) RevokeCard(id string, c models.CardResponse) error {
	args := s.Called(id, c)
	return args.Error(0)
}

func MakeFakeCardResponse() models.CardResponse {
	return MakeFakeCardResponseWith("test")
}

func MakeFakeCardResponseWith(text string) models.CardResponse {
	return models.CardResponse{
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

func AssertControllerRespose(t *testing.T, expected models.CardResponse, r []byte, e error) {
	assert.Nil(t, e)

	var actial models.CardResponse
	e = json.Unmarshal(r, &actial)
	assert.Nil(t, e, "Cannot restore object from response")
	assert.Equal(t, expected, actial)
}

func Test_GetCard_StorageReturnErr_ReturnErr(t *testing.T) {
	id := "test"
	errText := "Some error"
	mStorage := MockStorage{}
	mStorage.On("GetCard", id).Return(new(models.CardResponse), errors.New(errText))

	c := Controller{
		Storage: mStorage,
	}
	_, err := c.GetCard(id)
	assert.NotNil(t, err)
	assert.Equal(t, errText, err.Error())
}

func Test_GetCard_StorageReturnVal_ReturnJsonByte(t *testing.T) {
	id := "test"
	expected := MakeFakeCardResponse()
	mStorage := MockStorage{}
	mStorage.On("GetCard", id).Return(&expected, nil)

	c := Controller{
		Storage: mStorage,
	}
	r, err := c.GetCard(id)
	AssertControllerRespose(t, expected, r, err)
}

func Test_GetCard_StorageReturnNilValue_ReturnNilByte(t *testing.T) {
	id := "test"
	mStorage := MockStorage{}
	mStorage.On("GetCard", id).Return(nil, nil)

	c := Controller{
		Storage: mStorage,
	}
	r, err := c.GetCard(id)
	assert.Nil(t, err)
	assert.Nil(t, r)
}

func Test_SearchCards_BrokenRequestData_ReturnErr(t *testing.T) {
	errText := "Data has incorrect format"
	mStorage := MockStorage{}
	c := Controller{
		Storage: mStorage,
	}
	_, err := c.SearchCards([]byte("Test"))
	assert.NotNil(t, err)
	assert.Equal(t, errText, err.Error())
}

func Test_SearchCards_StorageReturnErr_ReturnErr(t *testing.T) {
	criteria := models.Criteria{
		Scope: "global",
	}
	data, _ := json.Marshal(&criteria)

	errText := "Some error"
	mStorage := MockStorage{}
	mStorage.On("SearchCards", criteria).Return([]models.CardResponse{}, errors.New(errText))

	c := Controller{
		Storage: mStorage,
	}
	_, err := c.SearchCards(data)
	assert.NotNil(t, err)
	assert.Equal(t, errText, err.Error())
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
		MakeFakeCardResponse(),
		MakeFakeCardResponse(),
	}
	mStorage := MockStorage{}
	mStorage.On("SearchCards", restoredCriteria).Return(expected, nil)

	c := Controller{
		Storage: mStorage,
	}
	r, err := c.SearchCards(criteria)

	assert.Nil(t, err)

	var actial []models.CardResponse
	err = json.Unmarshal(r, &actial)
	assert.Nil(t, err, "Cannot restore object from response")
	assert.Equal(t, expected, actial)
}

func Test_CreateCard_BrokenRequestData_ReturnErr(t *testing.T) {
	errText := "Data has incorrect format"
	mStorage := MockStorage{}
	c := Controller{
		Storage: mStorage,
	}
	_, err := c.CreateCard([]byte("Test"))
	assert.NotNil(t, err)
	assert.Equal(t, errText, err.Error())
}

func Test_CreateCard_StorageReturnErr_ReturnErr(t *testing.T) {
	param := MakeFakeCardResponse()
	data, _ := json.Marshal(&param)

	errText := "Some error"
	mStorage := MockStorage{}
	mStorage.On("CreateCard", param).Return(nil, errors.New(errText))

	c := Controller{
		Storage: mStorage,
	}
	_, err := c.CreateCard(data)
	assert.NotNil(t, err)
	assert.Equal(t, errText, err.Error())
}

func Test_CreateCard_StorageReturnVal_ReturnJsonByte(t *testing.T) {
	param := MakeFakeCardResponse()
	data, _ := json.Marshal(&param)

	expected := MakeFakeCardResponseWith("expected")
	mStorage := MockStorage{}
	mStorage.On("CreateCard", param).Return(&expected, nil)

	c := Controller{
		Storage: mStorage,
	}
	r, err := c.CreateCard(data)

	AssertControllerRespose(t, expected, r, err)
}

func Test_RevokeCard_BrokenRequestData_ReturnErr(t *testing.T) {
	id := "test"
	errText := "Data has incorrect format"
	mStorage := MockStorage{}
	c := Controller{
		Storage: mStorage,
	}
	err := c.RevokeCard(id, []byte("Test"))
	assert.NotNil(t, err)
	assert.Equal(t, errText, err.Error())
}

func Test_RevokeCard_StorageReturnErr_ReturnErr(t *testing.T) {
	id := "test"
	param := MakeFakeCardResponse()
	data, _ := json.Marshal(&param)

	errText := "Some error"
	mStorage := MockStorage{}
	mStorage.On("RevokeCard", id, param).Return(errors.New(errText))

	c := Controller{
		Storage: mStorage,
	}
	err := c.RevokeCard(id, data)
	assert.NotNil(t, err)
	assert.Equal(t, errText, err.Error())
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
	err := c.RevokeCard(id, data)
	assert.Nil(t, err)
}
