package sync

import (
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
	v, ok := args.Get(0).([]models.CardResponse)
	if ok {
		return v, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}
func (s MockStorage) CreateCard(r *models.CardResponse) (*models.CardResponse, error) {
	args := s.Called(r)
	v, ok := args.Get(0).(*models.CardResponse)
	if ok {
		return v, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}
func (s MockStorage) RevokeCard(id string, r *models.CardResponse) error {
	args := s.Called(id, r)
	return args.Error(0)
}

type MockLogger struct {
	mock.Mock
}

func (l MockLogger) Println(v ...interface{}) {
	l.Called()
}
func (l MockLogger) Printf(format string, v ...interface{}) {
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
		Logger: MockLogger{},
	}

	actual, err := sync.GetCard(id)
	assert.Nil(t, err, "Error should be nil when we found a card in local storage")
	assert.Equal(t, expected, actual)
}

func Test_GetCard_LocalReturnErr_LogErr(t *testing.T) {
	var (
		local, remote MockStorage
	)
	id := "test"
	l := MockLogger{}
	l.On("Println").Once()
	err := errors.New("Some error")
	local.On("GetCard", id).Return(nil, err)
	remote.On("GetCard", id).Return(nil, nil)

	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: l,
	}
	sync.GetCard(id)
	l.AssertExpectations(t)
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
		Logger: MockLogger{},
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
		Logger: MockLogger{},
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
		Logger: MockLogger{},
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
	local.On("GetCard", id).Return(nil, nil)
	remote.On("GetCard", id).Return(nil, errors.New("Some error"))

	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: MockLogger{},
	}
	card, err := sync.GetCard(id)
	assert.Nil(t, card)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Some error")
}

func Test_SearchCards_LocalReturnErr_LogErr(t *testing.T) {
	var (
		local, remote MockStorage
		c             models.Criteria
		l             MockLogger
	)
	local.On("SearchCards", c).Return(nil, errors.New("Some error"))
	remote.On("SearchCards", mock.Anything).Return(nil, nil)
	l.On("Println").Once()
	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: l,
	}
	sync.SearchCards(c)
	l.AssertExpectations(t)
}

func Test_SearchCards_LocalReturnValCountEqualOfIdentitiesInCriteria_ReturnVal(t *testing.T) {
	var (
		local, remote MockStorage
		c             models.Criteria
		l             MockLogger
		expected      []models.CardResponse
	)
	c.Identities = append(c.Identities, "test")
	expected = append(expected, *MakeFakeCardResponseWith("test"))

	local.On("SearchCards", c).Return(expected, nil)
	l.On("Println")

	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: l,
	}
	actual, err := sync.SearchCards(c)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_SearchCards_LocalReturnValCountNotEqualOfIdentitiesInCriteriaRemoteReturnVal_ReturnVal(t *testing.T) {
	var (
		local, remote MockStorage
		c             models.Criteria
		l             MockLogger
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
	l.On("Println")

	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: l,
	}
	actual, err := sync.SearchCards(c)

	assert.Nil(t, err)
	assert.Equal(t, remResult, actual)
}

func Test_SearchCards_LocalReturnValCountNotEqualOfIdentitiesInCriteriaRemoteReturnVal_AddToLocal(t *testing.T) {
	var (
		local, remote MockStorage
		c             models.Criteria
		l             MockLogger
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
	l.On("Println")

	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: l,
	}
	sync.SearchCards(c)

	// TODO Check passed parameter
	local.AssertExpectations(t)
}

func Test_SearchCards_LocalReturnValCountNotEqualOfIdentitiesInCriteriaRemoteReturnErr_ReturnErr(t *testing.T) {
	var (
		local, remote MockStorage
		l             MockLogger
		c             models.Criteria
	)
	c.Identities = append(c.Identities, "test")
	local.On("SearchCards", c).Return(nil, nil)
	remote.On("SearchCards", c).Return(nil, errors.New("Some error"))

	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: l,
	}
	_, err := sync.SearchCards(c)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Some error")
}

func Test_CreateCard_RemoteReturnErr_ReturnErr(t *testing.T) {
	var (
		local, remote MockStorage
		l             MockLogger
		r             models.CardResponse
	)
	remote.On("CreateCard", &r).Return(nil, errors.New("Some error"))

	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: l,
	}

	_, err := sync.CreateCard(&r)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Some error")
}

func Test_CreateCard_RemoteReturnVal_ReturnVal(t *testing.T) {
	var (
		local, remote MockStorage
		l             MockLogger
		r             models.CardResponse
	)
	expected := MakeFakeCardResponseWith("text")
	local.On("CreateCard", expected).Return(nil, nil)
	remote.On("CreateCard", &r).Return(expected, nil)

	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: l,
	}

	actual, err := sync.CreateCard(&r)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_CreateCard_RemoteReturnVal_AddToLocal(t *testing.T) {
	var (
		local, remote MockStorage
		l             MockLogger
		r             models.CardResponse
	)
	expected := MakeFakeCardResponseWith("text")
	local.On("CreateCard", expected).Return(nil, nil).Once()
	remote.On("CreateCard", &r).Return(expected, nil)

	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: l,
	}

	sync.CreateCard(&r)

	// TODO Check passed parameter
	local.AssertExpectations(t)
}

func Test_CreateCard_RemoteReturnValLocalReturnErr_LogErr(t *testing.T) {
	var (
		local, remote MockStorage
		r             models.CardResponse
		l             MockLogger
	)
	l.On("Println").Once()

	expected := MakeFakeCardResponseWith("text")
	local.On("CreateCard", expected).Return(nil, errors.New("Some error"))
	remote.On("CreateCard", &r).Return(expected, nil)

	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: l,
	}

	sync.CreateCard(&r)
	l.AssertExpectations(t)
}

func Test_RevokeCard_RemoteReturnErr_ReturnErr(t *testing.T) {
	var (
		local, remote MockStorage
		l             MockLogger
		r             models.CardResponse
	)
	id := "test"
	remote.On("RevokeCard", id, &r).Return(errors.New("Some error"))

	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: l,
	}

	err := sync.RevokeCard(id, &r)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Some error")
}

func Test_RevokeCard_RemoteReturnNil_ReturnNil(t *testing.T) {
	var (
		local, remote MockStorage
		l             MockLogger
		r             models.CardResponse
	)
	id := "test"
	remote.On("RevokeCard", id, &r).Return(nil)
	local.On("RevokeCard", id, &r).Return(nil)

	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: l,
	}

	err := sync.RevokeCard(id, &r)
	assert.Nil(t, err)
}

func Test_RevokeCard_RemoteReturnNil_DeleteToLocal(t *testing.T) {
	var (
		local, remote MockStorage
		l             MockLogger
		r             models.CardResponse
	)
	id := "test"
	remote.On("RevokeCard", id, &r).Return(nil)
	local.On("RevokeCard", id, &r).Return(nil).Once()

	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: l,
	}

	sync.RevokeCard(id, &r)

	// TODO Check passed parameter
	local.AssertExpectations(t)
}

func Test_RevokeCard_RemoteReturnValLocalReturnErr_LogErr(t *testing.T) {
	var (
		local, remote MockStorage
		r             models.CardResponse
		l             MockLogger
	)
	l.On("Println").Once()

	id := "test"
	remote.On("RevokeCard", id, &r).Return(nil)
	local.On("RevokeCard", id, &r).Return(errors.New("Some error")).Once()

	sync := Sync{
		Local:  local,
		Remote: remote,
		Logger: l,
	}

	sync.RevokeCard(id, &r)
	l.AssertExpectations(t)
}