package sync

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/virgilsecurity/virgild/models"
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

func Test_GetCard_LocalHasValue_ReturnVal(t *testing.T) {
	var (
		local, remote MockStorage
	)
	id := "test"
	expected := MakeFakeCardResponseWith("test")
	local.On("GetCard", id).Return(expected, nil)

	sync := Sync{
		Local:  local,
		Remote: remote,
	}

	actual, err := sync.GetCard(id)
	assert.Nil(t, err, "Error should be nil when we found a card in local storage")
	assert.Equal(t, expected, actual)
}

func Test_GetCard_LocalReturnErr_ReturnRemoteVal(t *testing.T) {
	var (
		local, remote MockStorage
	)
	id := "test"
	expected := MakeFakeCardResponseWith("test")
	local.On("GetCard", id).Return(nil, models.MakeError(12312))
	remote.On("GetCard", id).Return(expected, nil)

	sync := Sync{
		Local:  local,
		Remote: remote,
	}
	actual, err := sync.GetCard(id)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_GetCard_LocalNotFoundCardRemoteReturnVal_ReturnVal(t *testing.T) {
	var (
		local, remote MockStorage
	)
	id := "test"
	expected := MakeFakeCardResponseWith("test")

	local.On("GetCard", id).Return(nil, nil)
	local.On("CreateCard", mock.Anything).Return(expected, nil)

	remote.On("GetCard", id).Return(expected, nil)

	sync := Sync{
		Local:  local,
		Remote: remote,
	}
	actual, err := sync.GetCard(id)
	assert.Nil(t, err, "Error should be nil when we found a card in local storage")
	assert.Equal(t, expected, actual)
}

func Test_GetCard_LocalNotFoundCardRemoteReturnVal_AddToLocal(t *testing.T) {
	var (
		local, remote MockStorage
	)
	id := "test"
	expected := MakeFakeCardResponseWith("test")

	local.On("GetCard", id).Return(nil, nil)
	local.On("CreateCard", expected).Return(expected, nil)

	remote.On("GetCard", id).Return(expected, nil).Once()

	sync := Sync{
		Local:  local,
		Remote: remote,
	}
	sync.GetCard(id)
	local.AssertExpectations(t)
}

func Test_GetCard_LocalNotFoundCardRemoteNotFoundCard_ReturnNil(t *testing.T) {
	var (
		local, remote MockStorage
	)
	id := "test"
	local.On("GetCard", id).Return(nil, nil)
	remote.On("GetCard", id).Return(nil, nil).Once()

	sync := Sync{
		Local:  local,
		Remote: remote,
	}
	card, err := sync.GetCard(id)
	assert.Nil(t, card)
	assert.Nil(t, err)
}

func Test_GetCard_LocalNotFoundCardRemoteReturnErr_ReturnErr(t *testing.T) {
	var (
		local, remote MockStorage
	)
	id := "test"
	expected := models.MakeError(123)
	local.On("GetCard", id).Return(nil, nil)
	remote.On("GetCard", id).Return(nil, expected)

	sync := Sync{
		Local:  local,
		Remote: remote,
	}
	card, err := sync.GetCard(id)
	assert.Nil(t, card)
	assert.NotNil(t, err)
	assert.Equal(t, expected, err)
}

func Test_SearchCards_LocalReturnErr_ReturnRemoteVal(t *testing.T) {
	var (
		local, remote MockStorage
		c             models.Criteria
	)
	expected := []models.CardResponse{
		*MakeFakeCardResponseWith("test1"),
		*MakeFakeCardResponseWith("test2"),
	}
	local.On("SearchCards", c).Return(nil, models.MakeError(1234))
	remote.On("SearchCards", mock.Anything).Return(expected, nil)
	sync := Sync{
		Local:  local,
		Remote: remote,
	}
	actual, _ := sync.SearchCards(c)
	assert.Equal(t, expected, actual)
}

func Test_SearchCards_LocalReturnValCountEqualOfIdentitiesInCriteria_ReturnVal(t *testing.T) {
	var (
		local, remote MockStorage
		c             models.Criteria
		expected      []models.CardResponse
	)
	c.Identities = append(c.Identities, "test")
	expected = append(expected, *MakeFakeCardResponseWith("test"))

	local.On("SearchCards", c).Return(expected, nil)

	sync := Sync{
		Local:  local,
		Remote: remote,
	}
	actual, err := sync.SearchCards(c)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_SearchCards_LocalReturnValCountNotEqualOfIdentitiesInCriteriaRemoteReturnVal_ReturnVal(t *testing.T) {
	var (
		local, remote MockStorage
		c             models.Criteria
		locResult     []models.CardResponse
		remResult     []models.CardResponse
	)
	c.Identities = append(c.Identities, "test", "new")

	cr1 := MakeFakeCardResponseWith("test1")
	cr2 := MakeFakeCardResponseWith("test2")
	remResult = append(remResult, *cr1, *cr2)

	local.On("SearchCards", c).Return(locResult, nil)
	local.On("CreateCard", mock.Anything).Return(nil, nil)
	remote.On("SearchCards", c).Return(remResult, nil)

	sync := Sync{
		Local:  local,
		Remote: remote,
	}
	actual, err := sync.SearchCards(c)

	assert.Nil(t, err)
	assert.Equal(t, remResult, actual)
}

func Test_SearchCards_LocalReturnValCountNotEqualOfIdentitiesInCriteriaRemoteReturnVal_AddToLocal(t *testing.T) {
	var (
		local, remote MockStorage
		c             models.Criteria
		locResult     []models.CardResponse
		remResult     []models.CardResponse
	)
	c.Identities = append(c.Identities, "test", "new")

	cr1 := MakeFakeCardResponseWith("test1")
	cr2 := MakeFakeCardResponseWith("test2")
	locResult = append(locResult, *cr1)
	remResult = append(remResult, *cr1, *cr2)

	local.On("SearchCards", c).Return(locResult, nil)
	local.On("CreateCard", mock.Anything).Return(nil, nil).Once()
	remote.On("SearchCards", c).Return(remResult, nil)

	sync := Sync{
		Local:  local,
		Remote: remote,
	}
	sync.SearchCards(c)

	// TODO Check passed parameter
	local.AssertExpectations(t)
}

func Test_SearchCards_LocalReturnValCountNotEqualOfIdentitiesInCriteriaRemoteReturnErr_ReturnErr(t *testing.T) {
	var (
		local, remote MockStorage
		c             models.Criteria
	)
	expected := models.MakeError(1234)
	c.Identities = append(c.Identities, "test")
	local.On("SearchCards", c).Return(nil, nil)
	remote.On("SearchCards", c).Return(nil, expected)

	sync := Sync{
		Local:  local,
		Remote: remote,
	}
	_, err := sync.SearchCards(c)
	assert.Equal(t, expected, err)
}

func Test_CreateCard_RemoteReturnErr_ReturnErr(t *testing.T) {
	var (
		local, remote MockStorage
		r             models.CardResponse
	)
	expected := models.MakeError(1234)
	remote.On("CreateCard", &r).Return(nil, expected)

	sync := Sync{
		Local:  local,
		Remote: remote,
	}

	_, err := sync.CreateCard(&r)
	assert.Equal(t, expected, err)
}

func Test_CreateCard_RemoteReturnVal_ReturnVal(t *testing.T) {
	var (
		local, remote MockStorage
		r             models.CardResponse
	)
	expected := MakeFakeCardResponseWith("text")
	local.On("CreateCard", expected).Return(nil, nil)
	remote.On("CreateCard", &r).Return(expected, nil)

	sync := Sync{
		Local:  local,
		Remote: remote,
	}

	actual, err := sync.CreateCard(&r)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_CreateCard_RemoteReturnVal_AddToLocal(t *testing.T) {
	var (
		local, remote MockStorage
		r             models.CardResponse
	)
	expected := MakeFakeCardResponseWith("text")
	local.On("CreateCard", expected).Return(nil, nil).Once()
	remote.On("CreateCard", &r).Return(expected, nil)

	sync := Sync{
		Local:  local,
		Remote: remote,
	}

	sync.CreateCard(&r)

	// TODO Check passed parameter
	local.AssertExpectations(t)
}

func Test_RevokeCard_RemoteReturnErr_ReturnErr(t *testing.T) {
	var (
		local, remote MockStorage
		r             models.CardResponse
	)
	id := "test"
	expected := models.MakeError(1234)
	remote.On("RevokeCard", id, &r).Return(expected)

	sync := Sync{
		Local:  local,
		Remote: remote,
	}

	err := sync.RevokeCard(id, &r)
	assert.NotNil(t, err)
	assert.Equal(t, expected, err)
}

func Test_RevokeCard_RemoteReturnNil_ReturnNil(t *testing.T) {
	var (
		local, remote MockStorage
		r             models.CardResponse
	)
	id := "test"
	remote.On("RevokeCard", id, &r).Return(nil)
	local.On("RevokeCard", id, &r).Return(nil)

	sync := Sync{
		Local:  local,
		Remote: remote,
	}

	err := sync.RevokeCard(id, &r)
	assert.Nil(t, err)
}

func Test_RevokeCard_RemoteReturnNil_DeleteInLocal(t *testing.T) {
	var (
		local, remote MockStorage
		r             models.CardResponse
	)
	id := "test"
	remote.On("RevokeCard", id, &r).Return(nil)
	local.On("RevokeCard", id, &r).Return(nil).Once()

	sync := Sync{
		Local:  local,
		Remote: remote,
	}

	sync.RevokeCard(id, &r)

	// TODO Check passed parameter
	local.AssertExpectations(t)
}
