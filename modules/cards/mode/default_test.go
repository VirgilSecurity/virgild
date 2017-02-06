package mode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"gopkg.in/virgil.v4"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fakeCardRepository struct {
	mock.Mock
}

func (f *fakeCardRepository) Get(id string) (c *core.SqlCard, err error) {
	args := f.Called(id)
	c, _ = args.Get(0).(*core.SqlCard)
	err = args.Error(1)
	return
}

func (f *fakeCardRepository) Find(identitis []string, identityType string, scope string) (c []core.SqlCard, err error) {
	args := f.Called(identitis, identityType, scope)
	c, _ = args.Get(0).([]core.SqlCard)
	err = args.Error(1)
	return
}

func (f *fakeCardRepository) Add(cs core.SqlCard) error {
	args := f.Called(cs)
	return args.Error(0)
}

func (f *fakeCardRepository) DeleteById(id string) error {
	args := f.Called(id)
	return args.Error(0)
}

func (f *fakeCardRepository) DeleteBySearch(identitis []string, identityType string, scope string) error {
	args := f.Called(identitis, identityType, scope)
	return args.Error(0)
}

type fakeFingerprint struct {
	mock.Mock
}

func (f *fakeFingerprint) Calculate(data []byte) string {
	args := f.Called(data)
	return args.String(0)
}

func TestDefault_Get_RepoReturnErr_ReturnErr(t *testing.T) {
	repo := new(fakeCardRepository)
	repo.On("Get", mock.Anything).Return(nil, fmt.Errorf("Error"))
	h := DefaultModeCardHandler{repo, nil}
	_, err := h.Get("id")
	assert.NotNil(t, err)
}

func TestDefault_Get_RepoReturnValErrorCodeNotZero_ReturnErr(t *testing.T) {
	repo := new(fakeCardRepository)
	repo.On("Get", mock.Anything).Return(&core.SqlCard{ErrorCode: -1}, nil)
	h := DefaultModeCardHandler{repo, nil}
	_, err := h.Get("id")
	assert.Equal(t, core.ResponseErrorCode(-1), err)
}

func TestDefault_Get_RepoReturnValCardInvalid_ReturnErr(t *testing.T) {
	repo := new(fakeCardRepository)
	repo.On("Get", mock.Anything).Return(&core.SqlCard{Card: []byte("it's broken data")}, nil)
	h := DefaultModeCardHandler{repo, nil}
	_, err := h.Get("id")
	assert.NotNil(t, err)
}

func TestDefault_Get_RepoReturnValCard_ReturnCard(t *testing.T) {
	expected := &core.Card{Snapshot: []byte("test"), Meta: core.CardMeta{Signatures: map[string][]byte{
		"test": []byte("sign"),
	}}}
	jcard, _ := json.Marshal(expected)
	repo := new(fakeCardRepository)
	repo.On("Get", "id").Return(&core.SqlCard{Card: jcard}, nil)
	h := DefaultModeCardHandler{repo, nil}

	actual, _ := h.Get("id")

	assert.Equal(t, expected, actual)
}

func TestDefault_Search_RepoReturnErr_ReturnErr(t *testing.T) {
	repo := new(fakeCardRepository)
	repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("Error"))
	h := DefaultModeCardHandler{repo, nil}
	_, err := h.Search(&virgil.Criteria{})
	assert.NotNil(t, err)
}

func TestDefault_Search_RepoReturnValErrorCodeNotZero_ReturnErr(t *testing.T) {
	repo := new(fakeCardRepository)
	repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return([]core.SqlCard{core.SqlCard{ErrorCode: -1}}, nil)
	h := DefaultModeCardHandler{repo, nil}
	cs, err := h.Search(&virgil.Criteria{})
	assert.Len(t, cs, 0)
	assert.Nil(t, err)
}

func TestDefault_Search_RepoReturnValCardInvalid_ReturnErr(t *testing.T) {
	repo := new(fakeCardRepository)
	repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return([]core.SqlCard{core.SqlCard{Card: []byte("it's broken data")}}, nil)
	h := DefaultModeCardHandler{repo, nil}
	cs, err := h.Search(&virgil.Criteria{})

	assert.Len(t, cs, 0)
	assert.Nil(t, err)
}

func TestDefault_Search_RepoReturnValCard_ReturnCard(t *testing.T) {
	crit := &virgil.Criteria{
		Identities:   []string{"test", "test1"},
		IdentityType: "app",
		Scope:        virgil.CardScope.Global,
	}
	expected := []core.Card{core.Card{Snapshot: []byte("test"), Meta: core.CardMeta{Signatures: map[string][]byte{
		"test": []byte("sign"),
	}}}}
	jcard, _ := json.Marshal(expected[0])
	repo := new(fakeCardRepository)
	repo.On("Find", crit.Identities, crit.IdentityType, string(crit.Scope)).Return([]core.SqlCard{core.SqlCard{Card: jcard}}, nil)
	h := DefaultModeCardHandler{repo, nil}

	actual, _ := h.Search(crit)

	assert.Equal(t, expected, actual)
}

func TestDefault_Create_RepoReturnErr_ReturnErr(t *testing.T) {
	repo := new(fakeCardRepository)
	repo.On("Add", mock.Anything).Return(fmt.Errorf("Error"))
	f := new(fakeFingerprint)
	f.On("Calculate", mock.Anything).Return("id")
	h := DefaultModeCardHandler{repo, f}

	_, err := h.Create(&core.CreateCardRequest{})

	assert.NotNil(t, err)
}

func TestDefault_Create_StoreCard(t *testing.T) {
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Scope:        virgil.CardScope.Application,
			Data:         map[string]string{"test": "data"},
			DeviceInfo:   virgil.DeviceInfo{Device: "iphone", DeviceName: "device name"},
			Identity:     "identity",
			IdentityType: "app",
			PublicKey:    []byte("test pub key"),
		},
		Request: virgil.SignableRequest{
			Snapshot: []byte("snapshot"),
			Meta: virgil.RequestMeta{
				Signatures: map[string][]byte{"test": []byte("sign")},
			},
		},
	}
	f := new(fakeFingerprint)
	f.On("Calculate", []byte("snapshot")).Return("id")
	repo := new(fakeCardRepository)
	repo.On("Add", mock.MatchedBy(func(scard core.SqlCard) bool {
		actual, _ := sqlCard2Card(&scard)
		return actual.ID == "id" &&
			bytes.Equal(actual.Snapshot, req.Request.Snapshot) &&
			actual.Meta.CardVersion == "v4" &&
			bytes.Equal(actual.Meta.Signatures["test"], []byte("sign")) &&
			scard.CardID == "id" &&
			scard.Deleted == false &&
			scard.ErrorCode == 0 &&
			scard.Identity == req.Info.Identity &&
			scard.IdentityType == req.Info.IdentityType &&
			scard.Scope == string(req.Info.Scope)
	})).Return(nil).Once()

	h := DefaultModeCardHandler{repo, f}
	h.Create(req)

	repo.AssertExpectations(t)
}

func TestDefault_Create_ReturnCard(t *testing.T) {
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Scope:        virgil.CardScope.Application,
			Data:         map[string]string{"test": "data"},
			DeviceInfo:   virgil.DeviceInfo{Device: "iphone", DeviceName: "device name"},
			Identity:     "identity",
			IdentityType: "app",
			PublicKey:    []byte("test pub key"),
		},
		Request: virgil.SignableRequest{
			Snapshot: []byte("snapshot"),
			Meta: virgil.RequestMeta{
				Signatures: map[string][]byte{"test": []byte("sign")},
			},
		},
	}
	f := new(fakeFingerprint)
	f.On("Calculate", []byte("snapshot")).Return("id")
	repo := new(fakeCardRepository)
	repo.On("Add", mock.Anything).Return(nil)

	h := DefaultModeCardHandler{repo, f}
	c, _ := h.Create(req)

	assert.Equal(t, req.Request.Snapshot, c.Snapshot)
	assert.Equal(t, "v4", c.Meta.CardVersion)
	assert.Equal(t, req.Request.Meta.Signatures, c.Meta.Signatures)
	assert.Equal(t, "id", c.ID)
}

func TestDefault_Revoke_ReturnErr(t *testing.T) {
	const id = "id"
	repo := new(fakeCardRepository)
	repo.On("DeleteById", id).Return(fmt.Errorf("Error"))

	h := DefaultModeCardHandler{repo, nil}
	err := h.Revoke(&core.RevokeCardRequest{Info: virgil.RevokeCardRequest{ID: id}})

	assert.NotNil(t, err)
}
