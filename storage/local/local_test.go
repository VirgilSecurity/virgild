package local

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	. "github.com/virgilsecurity/virgil-apps-cards-cacher/database/sqlmodels"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
	"testing"
)

type MockCardRepository struct {
	mock.Mock
}

func (r *MockCardRepository) Get(id string) (*CardSql, error) {
	args := r.Called(id)
	if v, ok := args.Get(0).(*CardSql); ok {
		return v, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (r *MockCardRepository) Find(c models.Criteria) ([]CardSql, error) {
	args := r.Called(c)
	return args.Get(0).([]CardSql), args.Error(1)
}

func (r *MockCardRepository) Add(c CardSql) error {
	args := r.Called(c)
	return args.Error(0)
}

func (r *MockCardRepository) Delete(id string) error {
	args := r.Called(id)
	return args.Error(0)
}

func MakeFakeCardResponseWith(text string) *models.CardResponse {
	return &models.CardResponse{
		ID:       text,
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

func MakeFakePairCardSqlCardResponseWith(text string) (*CardSql, *models.CardResponse) {
	r := MakeFakeCardResponseWith(text)
	jr, _ := json.Marshal(r)
	c := &CardSql{
		Card:         string(jr[:]),
		Id:           r.ID,
		Identity:     text,
		IdentityType: text,
		Scope:        text,
	}
	return c, r
}

func Test_GetCard_ResultEmpty_ReturnNil(t *testing.T) {
	id := "Test"
	r := MockCardRepository{}
	r.On("Get", id).Return(nil, nil)
	local := Local{
		Repo: &r,
	}
	c, err := local.GetCard(id)
	assert.Nil(t, c)
	assert.Nil(t, err)
}

func Test_GetCard_ResultVal_ReturnVal(t *testing.T) {
	id := "Test"
	cs, expected := MakeFakePairCardSqlCardResponseWith("test")
	r := MockCardRepository{}
	r.On("Get", id).Return(cs, nil)
	local := Local{
		Repo: &r,
	}
	actual, err := local.GetCard(id)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_GetCard_ReturnNil_ReturnNil(t *testing.T) {
	id := "Test"
	r := MockCardRepository{}
	r.On("Get", id).Return(nil, errors.New("Some error"))
	local := Local{
		Repo: &r,
	}
	_, err := local.GetCard(id)
	assert.EqualError(t, err, "Some error")
}

func Test_SearchCards_ReturnVal_ReturnVal(t *testing.T) {
	cs1, r1 := MakeFakePairCardSqlCardResponseWith("test1")
	cs2, r2 := MakeFakePairCardSqlCardResponseWith("test2")
	expected := []models.CardResponse{
		*r1,
		*r2,
	}
	cs := []CardSql{
		*cs1,
		*cs2,
	}
	c := models.Criteria{}
	r := MockCardRepository{}
	r.On("Find", c).Return(cs, nil)
	local := Local{
		Repo: &r,
	}
	actual, err := local.SearchCards(c)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_SearchCards_ReturnEmpty_ReturnEmpty(t *testing.T) {
	c := models.Criteria{}
	r := MockCardRepository{}
	r.On("Find", c).Return([]CardSql{}, nil)
	local := Local{
		Repo: &r,
	}
	actual, _ := local.SearchCards(c)
	assert.Len(t, actual, 0)
}

func Test_SearchCards_ReturnErr_ReturnErr(t *testing.T) {
	c := models.Criteria{}
	r := MockCardRepository{}
	r.On("Find", c).Return([]CardSql{}, errors.New("Some error"))
	local := Local{
		Repo: &r,
	}
	_, err := local.SearchCards(c)
	assert.EqualError(t, err, "Some error")
}

func Test_CreateCard_BrokenSnapshout_ReturnErr(t *testing.T) {
	cr := MakeFakeCardResponseWith("test")
	r := MockCardRepository{}
	local := Local{
		Repo: &r,
	}
	_, err := local.CreateCard(cr)
	assert.IsType(t, models.ErrorResponse{}, err)
	assert.Equal(t, 30107, err.(models.ErrorResponse).Code)
}

func Test_CreateCard_AddVal_ActionInvoked(t *testing.T) {
	req := CardRequest{
		Identity:     "identity",
		IdentityType: "application",
		PublicKey:    []byte(`some value`),
		Scope:        "global",
		Data: map[string]string{
			"test": "test",
		},
		DeviceInfo: DeviceInfo{
			Device:     "iphone7",
			DeviceName: "my",
		},
	}
	jReq, _ := json.Marshal(req)
	cr := &models.CardResponse{
		ID:       "id",
		Snapshot: jReq,
	}
	jCr, _ := json.Marshal(cr)
	expected := CardSql{
		Id:           cr.ID,
		Identity:     req.Identity,
		IdentityType: req.IdentityType,
		Scope:        req.Scope,
		Card:         string(jCr[:]),
	}
	r := MockCardRepository{}
	r.On("Add", expected).Return(nil).Once()
	local := Local{
		Repo: &r,
	}
	local.CreateCard(cr)
	r.AssertCalled(t, "Add", expected)
}

func Test_CreateCard_RepoReturnNil_ReturnErr(t *testing.T) {
	req := CardRequest{
		Identity:     "identity",
		IdentityType: "application",
		PublicKey:    []byte(`some value`),
		Scope:        "global",
		Data: map[string]string{
			"test": "test",
		},
		DeviceInfo: DeviceInfo{
			Device:     "iphone7",
			DeviceName: "my",
		},
	}
	jReq, _ := json.Marshal(req)
	cr := &models.CardResponse{
		ID:       "id",
		Snapshot: jReq,
	}
	jCr, _ := json.Marshal(cr)
	expected := CardSql{
		Id:           cr.ID,
		Identity:     req.Identity,
		IdentityType: req.IdentityType,
		Scope:        req.Scope,
		Card:         string(jCr[:]),
	}
	r := MockCardRepository{}
	r.On("Add", expected).Return(errors.New("Some error")).Once()
	local := Local{
		Repo: &r,
	}
	_, err := local.CreateCard(cr)
	assert.EqualError(t, err, "Some error")
}

func Test_RevokeCard_RepoReturnNil_ReturnNil(t *testing.T) {
	id := "id"
	rc := MakeFakeCardResponseWith("test")
	r := MockCardRepository{}
	r.On("Delete", id).Return(nil)
	local := Local{
		Repo: &r,
	}
	err := local.RevokeCard(id, rc)
	assert.NoError(t, err)
}

func Test_RevokeCard_RepoReturnErr_ReturnErr(t *testing.T) {
	id := "id"
	rc := MakeFakeCardResponseWith("test")
	r := MockCardRepository{}
	r.On("Delete", id).Return(errors.New("Some error"))
	local := Local{
		Repo: &r,
	}
	err := local.RevokeCard(id, rc)
	assert.EqualError(t, err, "Some error")
}
