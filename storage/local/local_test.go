package local

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	. "github.com/virgilsecurity/virgil-apps-cards-cacher/database/sqlmodels"
	"github.com/virgilsecurity/virgil-apps-cards-cacher/models"
	"gopkg.in/virgilsecurity/virgil-sdk-go.v4"
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

type MockLogger struct {
	mock.Mock
}

func (l *MockLogger) Println(v ...interface{}) {
	l.Called(v...)
}
func (l *MockLogger) Printf(format string, v ...interface{}) {
	l.Called()
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
	l := MockLogger{}
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	actual, err := local.GetCard(id)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_GetCard_ReturnErr_ReturnErr(t *testing.T) {
	id := "Test"
	expected := models.MakeError(10000)
	r := MockCardRepository{}
	r.On("Get", id).Return(nil, errors.New("Some error"))
	l := MockLogger{}
	l.On("Printf")
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	_, err := local.GetCard(id)
	assert.Equal(t, expected, err)
}

func Test_GetCard_RepoReturnErr_LogErr(t *testing.T) {
	id := "Test"
	r := MockCardRepository{}
	r.On("Get", id).Return(nil, errors.New("Some error"))
	l := MockLogger{}
	l.On("Printf").Once()
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	local.GetCard(id)
	l.AssertExpectations(t)
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
	l := MockLogger{}
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	actual, err := local.SearchCards(c)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_SearchCards_ReturnEmpty_ReturnEmpty(t *testing.T) {
	c := models.Criteria{}
	r := MockCardRepository{}
	r.On("Find", c).Return([]CardSql{}, nil)
	l := MockLogger{}
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	actual, _ := local.SearchCards(c)
	assert.Len(t, actual, 0)
}

func Test_SearchCards_ReturnErr_ReturnErr(t *testing.T) {
	c := models.Criteria{}
	expected := models.MakeError(10000)
	r := MockCardRepository{}
	r.On("Find", c).Return([]CardSql{}, errors.New("Some error"))
	l := MockLogger{}
	l.On("Printf")
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	_, err := local.SearchCards(c)
	assert.Equal(t, expected, err)
}

func Test_SearchCards_ReturnErr_LogErr(t *testing.T) {
	c := models.Criteria{}
	r := MockCardRepository{}
	r.On("Find", c).Return([]CardSql{}, errors.New("Some error"))
	l := MockLogger{}
	l.On("Printf").Once()
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	local.SearchCards(c)
	l.AssertExpectations(t)
}

func Test_CreateCard_BrokenSnapshout_ReturnErr(t *testing.T) {
	cr := MakeFakeCardResponseWith("test")
	expected := models.MakeError(30107)
	r := MockCardRepository{}
	l := MockLogger{}
	l.On("Printf")
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	_, err := local.CreateCard(cr)
	assert.Equal(t, expected, err)
}

func Test_CreateCard_BrokenSnapshout_LogErr(t *testing.T) {
	cr := MakeFakeCardResponseWith("test")
	r := MockCardRepository{}
	l := MockLogger{}
	l.On("Printf").Once()
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	local.CreateCard(cr)
	l.AssertExpectations(t)
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
	l := MockLogger{}
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	local.CreateCard(cr)
	r.AssertCalled(t, "Add", expected)
}

func Test_CreateCard_CardIdEmpty_CalcId(t *testing.T) {
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
		ID:       "",
		Snapshot: jReq,
	}

	crypto := virgil.Crypto()
	fp := crypto.CalculateFingerprint(jReq)
	id := hex.EncodeToString(fp)

	ecr := &models.CardResponse{
		ID:       id,
		Snapshot: jReq,
	}

	jCr, _ := json.Marshal(ecr)
	expected := CardSql{
		Id:           id,
		Identity:     req.Identity,
		IdentityType: req.IdentityType,
		Scope:        req.Scope,
		Card:         string(jCr[:]),
	}
	r := MockCardRepository{}
	r.On("Add", expected).Return(nil).Once()
	l := MockLogger{}
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	local.CreateCard(cr)
	r.AssertCalled(t, "Add", expected)
}

func Test_CreateCard_RepoReturnErr_ReturnErr(t *testing.T) {
	expected := models.MakeError(10000)
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
	r := MockCardRepository{}
	r.On("Add", mock.Anything).Return(errors.New("Some error"))
	l := MockLogger{}
	l.On("Printf")
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	_, err := local.CreateCard(cr)
	assert.Equal(t, expected, err)
}

func Test_CreateCard_RepoReturnErr_LogErr(t *testing.T) {
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
	r := MockCardRepository{}
	r.On("Add", mock.Anything).Return(errors.New("Some error"))
	l := MockLogger{}
	l.On("Printf").Once()
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	local.CreateCard(cr)
	l.AssertExpectations(t)
}

func Test_RevokeCard_RepoReturnNil_ReturnNil(t *testing.T) {
	id := "id"
	rc := MakeFakeCardResponseWith("test")
	r := MockCardRepository{}
	r.On("Delete", id).Return(nil)
	l := MockLogger{}
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	err := local.RevokeCard(id, rc)
	assert.Nil(t, err)
}

func Test_RevokeCard_RepoReturnErr_ReturnErr(t *testing.T) {
	expected := models.MakeError(10000)
	id := "id"
	rc := MakeFakeCardResponseWith("test")
	r := MockCardRepository{}
	r.On("Delete", id).Return(errors.New("Some error"))
	l := MockLogger{}
	l.On("Printf")
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	err := local.RevokeCard(id, rc)
	assert.Equal(t, expected, err)
}

func Test_RevokeCard_RepoReturnErr_LogErr(t *testing.T) {
	id := "id"
	rc := MakeFakeCardResponseWith("test")
	r := MockCardRepository{}
	r.On("Delete", id).Return(errors.New("Some error"))
	l := MockLogger{}
	l.On("Printf").Once()
	local := Local{
		Repo:   &r,
		Logger: &l,
	}
	local.RevokeCard(id, rc)
	l.AssertExpectations(t)
}
