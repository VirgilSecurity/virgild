package mode

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	virgil "gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/errors"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fakeVirgilClient struct {
	mock.Mock
}

func (f *fakeVirgilClient) GetCard(id string) (c *virgil.Card, err error) {
	args := f.Called(id)
	c, _ = args.Get(0).(*virgil.Card)
	err = args.Error(1)
	return
}

func (f *fakeVirgilClient) SearchCards(crit *virgil.Criteria) (cs []*virgil.Card, err error) {
	args := f.Called(crit)
	cs, _ = args.Get(0).([]*virgil.Card)
	err = args.Error(1)
	return
}

func (f *fakeVirgilClient) CreateCard(req *virgil.SignableRequest) (c *virgil.Card, err error) {
	args := f.Called(req)
	c, _ = args.Get(0).(*virgil.Card)
	err = args.Error(1)
	return
}

func (f *fakeVirgilClient) RevokeCard(req *virgil.SignableRequest) error {
	args := f.Called(req)
	return args.Error(0)
}

func makeVCard() *virgil.Card {
	var r [32]byte
	rand.Read(r[:])

	identity := hex.EncodeToString(r[:])
	kp, _ := virgil.Crypto().GenerateKeypair()
	req, _ := virgil.NewCreateCardRequest(identity, "app", kp.PublicKey(), virgil.CardParams{
		Scope: virgil.CardScope.Application,
		Data: map[string]string{
			identity: identity,
		},
		DeviceInfo: virgil.DeviceInfo{
			Device:     "iphone",
			DeviceName: "iphone7",
		},
	})

	return &virgil.Card{
		ID:           hex.EncodeToString(virgil.Crypto().CalculateFingerprint(req.Snapshot)),
		Identity:     identity,
		IdentityType: "app",
		Scope:        virgil.CardScope.Application,
		Data: map[string]string{
			identity: identity,
		},
		DeviceInfo: virgil.DeviceInfo{
			Device:     "iphone",
			DeviceName: "iphone7",
		},
		Snapshot: req.Snapshot,
		Signatures: map[string][]byte{
			identity: []byte("sign"),
		},
	}
}

func TestApp_GetCard_RepoHasValue_ReturnVal(t *testing.T) {
	vcard := makeVCard()
	expected := vcard2Card(vcard)
	sCard, _ := vcard2SqlCard(vcard)
	sCard.ExpireAt = time.Now().Add(time.Hour).Unix()

	repo := new(fakeCardRepository)
	repo.On("Get", "id").Return(sCard, nil)

	a := AppModeCardHandler{repo, nil}

	actual, _ := a.Get("id")
	assert.Equal(t, expected, actual)
}

func TestApp_GetCard_TestApp_GetCard_RepoReturnErr_ReturnErr(t *testing.T) {
	repo := new(fakeCardRepository)
	repo.On("Get", mock.Anything).Return(nil, fmt.Errorf("ERROR"))

	a := AppModeCardHandler{repo, nil}
	_, err := a.Get("id")

	assert.NotNil(t, err)
}

func TestApp_GetCard_LocalEntityNotFoundRemoteReturnVal_ReturnVal(t *testing.T) {
	vcard := makeVCard()
	expected := vcard2Card(vcard)

	repo := new(fakeCardRepository)
	repo.On("Get", "id").Return(nil, core.ErrorEntityNotFound)
	repo.On("Add", mock.Anything).Return(nil)
	vc := new(fakeVirgilClient)
	vc.On("GetCard", "id").Return(vcard, nil)

	a := AppModeCardHandler{repo, vc}

	actual, _ := a.Get("id")
	assert.Equal(t, expected, actual)
}

func TestApp_GetCard_LocalEntityNotFoundRemoteReturnVal_AddToLocal(t *testing.T) {
	vcard := makeVCard()
	expected, _ := vcard2SqlCard(vcard)

	repo := new(fakeCardRepository)
	repo.On("Get", mock.Anything).Return(nil, core.ErrorEntityNotFound)
	repo.On("Add", *expected).Return(nil).Once()
	vc := new(fakeVirgilClient)
	vc.On("GetCard", mock.Anything).Return(vcard, nil)

	a := AppModeCardHandler{repo, vc}

	a.Get("id")

	repo.AssertExpectations(t)
}

func TestApp_GetCard_LocalEntityNotFoundRemoteEntityNotFoun_ReturnEntityNotFound(t *testing.T) {
	repo := new(fakeCardRepository)
	repo.On("Get", mock.Anything).Return(nil, core.ErrorEntityNotFound)
	repo.On("Add", mock.Anything).Return(nil)
	vc := new(fakeVirgilClient)
	vc.On("GetCard", mock.Anything).Return(nil, errors.NewHttpError(404, "test"))

	a := AppModeCardHandler{repo, vc}

	_, err := a.Get("id")
	assert.Equal(t, core.ErrorEntityNotFound, errors.Cause(err))
}

func TestApp_GetCard_LocalNotFoundRemoteReturnErr_ReturnErr(t *testing.T) {
	repo := new(fakeCardRepository)
	repo.On("Get", mock.Anything).Return(nil, core.ErrorEntityNotFound)
	repo.On("Add", mock.Anything).Return(nil)

	expected := errors.NewHttpError(500, "test")
	vc := new(fakeVirgilClient)
	vc.On("GetCard", mock.Anything).Return(nil, expected)

	a := AppModeCardHandler{repo, vc}

	_, err := a.Get("id")
	assert.Equal(t, expected, errors.Cause(err))
}

func TestApp_GetCard_LocalNotFoundRemoteReturnErr_CahceRespErr(t *testing.T) {
	repo := new(fakeCardRepository)
	repo.On("Get", mock.Anything).Return(nil, core.ErrorEntityNotFound)
	repo.On("Add", core.SqlCard{CardID: "id", ErrorCode: 100000}).Return(nil).Once()

	vc := new(fakeVirgilClient)
	vc.On("GetCard", mock.Anything).Return(nil, errors.NewServiceError(100000, 300, ""))

	a := AppModeCardHandler{repo, vc}

	a.Get("id")
	repo.AssertExpectations(t)
}

func TestApp_Search_RepoReturnErr_ReturnErr(t *testing.T) {
	repo := new(fakeCardRepository)
	repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("Error"))

	a := AppModeCardHandler{repo, nil}

	_, err := a.Search(&virgil.Criteria{})

	assert.NotNil(t, err)
}

func TestApp_Serarch_RepoReturnEmptyRemoteReturnErr_ReturnErr(t *testing.T) {
	vc := new(fakeVirgilClient)
	vc.On("SearchCards", mock.Anything).Return(nil, fmt.Errorf("Error"))
	repo := new(fakeCardRepository)
	repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(make([]core.SqlCard, 0), nil)

	a := AppModeCardHandler{repo, vc}

	_, err := a.Search(&virgil.Criteria{})

	assert.NotNil(t, err)
}

func TestApp_Serarch_RepoReturnEmptyRemoteReturnSDKErr_AddCache(t *testing.T) {
	c := virgil.Criteria{
		Identities:   []string{"test1", "test2"},
		IdentityType: "app",
		Scope:        virgil.CardScope.Application,
	}
	expected1 := core.SqlCard{
		Identity:     c.Identities[0],
		IdentityType: c.IdentityType,
		Scope:        string(c.Scope),
		ErrorCode:    300,
	}
	expected2 := core.SqlCard{
		Identity:     c.Identities[1],
		IdentityType: c.IdentityType,
		Scope:        string(c.Scope),
		ErrorCode:    300,
	}
	vc := new(fakeVirgilClient)
	vc.On("SearchCards", mock.Anything).Return(nil, errors.NewServiceError(300, 404, "msg"))
	repo := new(fakeCardRepository)
	repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(make([]core.SqlCard, 0), nil)
	repo.On("Add", expected1).Return(nil).Once()
	repo.On("Add", expected2).Return(nil).Once()

	a := AppModeCardHandler{repo, vc}

	a.Search(&c)

	repo.AssertExpectations(t)
}

func TestApp_Search_RepoReturnEmptyRemoteReturnCards_AddCardsToRepo(t *testing.T) {
	c := virgil.Criteria{
		Identities:   []string{"test1", "test2"},
		IdentityType: "app",
		Scope:        virgil.CardScope.Application,
	}
	vcard := makeVCard()
	expected, _ := vcard2SqlCard(vcard)
	vc := new(fakeVirgilClient)
	vc.On("SearchCards", mock.Anything).Return([]*virgil.Card{vcard}, nil)
	repo := new(fakeCardRepository)
	repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(make([]core.SqlCard, 0), nil)
	repo.On("Add", *expected).Return(nil).Once()

	a := AppModeCardHandler{repo, vc}

	a.Search(&c)

	repo.AssertExpectations(t)
}

func TestApp_Search_RepoReturnEmptyRemoteReturnCards_ReturnVal(t *testing.T) {
	c := &virgil.Criteria{
		Identities:   []string{"test1", "test2"},
		IdentityType: "app",
		Scope:        virgil.CardScope.Application,
	}
	vcard := makeVCard()
	expected := []core.Card{*vcard2Card(vcard)}
	vc := new(fakeVirgilClient)
	vc.On("SearchCards", c).Return([]*virgil.Card{vcard}, nil)
	repo := new(fakeCardRepository)
	repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(make([]core.SqlCard, 0), nil)
	repo.On("Add", mock.Anything).Return(nil)

	a := AppModeCardHandler{repo, vc}

	actual, _ := a.Search(c)

	assert.Equal(t, expected, actual)
}

func TestApp_Serarch_RepoReturnExpireCardsRemoteReturnErr_ReturnErr(t *testing.T) {
	vc := new(fakeVirgilClient)
	vc.On("SearchCards", mock.Anything).Return(nil, fmt.Errorf("Error"))
	repo := new(fakeCardRepository)
	exp := []core.SqlCard{core.SqlCard{ExpireAt: time.Now().Add(-time.Hour).Unix()}}
	repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(exp, nil)
	repo.On("DeleteBySearch", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	a := AppModeCardHandler{repo, vc}

	_, err := a.Search(&virgil.Criteria{})

	assert.NotNil(t, err)
}

func TestApp_Serarch_RepoReturnExpireCards_RemoveExpiredCache(t *testing.T) {
	c := virgil.Criteria{
		Identities:   []string{"test1", "test2"},
		IdentityType: "app",
		Scope:        virgil.CardScope.Application,
	}

	vc := new(fakeVirgilClient)
	vc.On("SearchCards", mock.Anything).Return(nil, fmt.Errorf("Error"))
	repo := new(fakeCardRepository)
	exp := []core.SqlCard{core.SqlCard{ExpireAt: time.Now().Add(-time.Hour).Unix()}}
	repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(exp, nil)
	repo.On("DeleteBySearch", c.Identities, c.IdentityType, string(c.Scope)).Return(nil).Once()

	a := AppModeCardHandler{repo, vc}

	a.Search(&c)

	repo.AssertExpectations(t)
}

func TestApp_Serarch_RepoReturnExpireCardsRemoteReturnSDKErr_AddCache(t *testing.T) {
	c := virgil.Criteria{
		Identities:   []string{"test1", "test2"},
		IdentityType: "app",
		Scope:        virgil.CardScope.Application,
	}
	expected1 := core.SqlCard{
		Identity:     c.Identities[0],
		IdentityType: c.IdentityType,
		Scope:        string(c.Scope),
		ErrorCode:    300,
	}
	expected2 := core.SqlCard{
		Identity:     c.Identities[1],
		IdentityType: c.IdentityType,
		Scope:        string(c.Scope),
		ErrorCode:    300,
	}
	vc := new(fakeVirgilClient)
	vc.On("SearchCards", mock.Anything).Return(nil, errors.NewServiceError(300, 404, "msg"))
	repo := new(fakeCardRepository)
	exp := []core.SqlCard{core.SqlCard{ExpireAt: time.Now().Add(-time.Hour).Unix()}}
	repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(exp, nil)
	repo.On("DeleteBySearch", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	repo.On("Add", expected1).Return(nil).Once()
	repo.On("Add", expected2).Return(nil).Once()

	a := AppModeCardHandler{repo, vc}

	a.Search(&c)

	repo.AssertExpectations(t)
}

func TestApp_Search_RepoReturnExpireCardsRemoteReturnCards_AddCardsToRepo(t *testing.T) {
	c := virgil.Criteria{
		Identities:   []string{"test1", "test2"},
		IdentityType: "app",
		Scope:        virgil.CardScope.Application,
	}
	vcard := makeVCard()
	expected, _ := vcard2SqlCard(vcard)
	vc := new(fakeVirgilClient)
	vc.On("SearchCards", mock.Anything).Return([]*virgil.Card{vcard}, nil)
	repo := new(fakeCardRepository)
	exp := []core.SqlCard{core.SqlCard{ExpireAt: time.Now().Add(-time.Hour).Unix()}}
	repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(exp, nil)
	repo.On("DeleteBySearch", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	repo.On("Add", *expected).Return(nil).Once()

	a := AppModeCardHandler{repo, vc}

	a.Search(&c)

	repo.AssertExpectations(t)
}

func TestApp_Search_RepoReturnExpireCardsRemoteReturnCards_ReturnVal(t *testing.T) {
	c := &virgil.Criteria{
		Identities:   []string{"test1", "test2"},
		IdentityType: "app",
		Scope:        virgil.CardScope.Application,
	}
	vcard := makeVCard()
	expected := []core.Card{*vcard2Card(vcard)}
	vc := new(fakeVirgilClient)
	vc.On("SearchCards", c).Return([]*virgil.Card{vcard}, nil)
	repo := new(fakeCardRepository)
	exp := []core.SqlCard{core.SqlCard{ExpireAt: time.Now().Add(-time.Hour).Unix()}}
	repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(exp, nil)
	repo.On("DeleteBySearch", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	repo.On("Add", mock.Anything).Return(nil)

	a := AppModeCardHandler{repo, vc}

	actual, _ := a.Search(c)

	assert.Equal(t, expected, actual)
}

func TestApp_Search_RepoReturnCards_ReturnVal(t *testing.T) {
	c := virgil.Criteria{
		Identities:   []string{"test1", "test2"},
		IdentityType: "app",
		Scope:        virgil.CardScope.Application,
	}
	vcard := makeVCard()
	scard, _ := vcard2SqlCard(vcard)
	scard.ExpireAt = time.Now().Add(time.Hour).Unix()
	expected := []core.Card{*vcard2Card(vcard)}
	repo := new(fakeCardRepository)
	repo.On("Find", c.Identities, c.IdentityType, string(c.Scope)).Return([]core.SqlCard{*scard}, nil)

	a := AppModeCardHandler{repo, nil}

	actual, _ := a.Search(&c)

	assert.Equal(t, expected, actual)
}

func TestApp_CreateCard_RemoteReturnErr_ReturnErr(t *testing.T) {
	vc := new(fakeVirgilClient)
	vc.On("CreateCard", mock.Anything).Return(nil, fmt.Errorf("Error"))

	a := AppModeCardHandler{nil, vc}

	_, err := a.Create(&core.CreateCardRequest{})
	assert.NotNil(t, err)
}

func TestApp_CreateCard_RemoteReturnVal_ReturnVal(t *testing.T) {
	vcard := makeVCard()
	expected := vcard2Card(vcard)
	req := &virgil.SignableRequest{Snapshot: vcard.Snapshot, Meta: virgil.RequestMeta{Signatures: vcard.Signatures}}

	vc := new(fakeVirgilClient)
	vc.On("CreateCard", req).Return(vcard, nil)
	repo := new(fakeCardRepository)
	repo.On("Add", mock.Anything).Return(nil)

	a := AppModeCardHandler{repo, vc}

	actual, _ := a.Create(&core.CreateCardRequest{Request: *req})
	assert.Equal(t, expected, actual)
}

func TestApp_CreateCard_RemoteReturnVal_AddToLocal(t *testing.T) {
	vcard := makeVCard()
	expected, _ := vcard2SqlCard(vcard)

	vc := new(fakeVirgilClient)
	vc.On("CreateCard", mock.Anything).Return(vcard, nil)
	repo := new(fakeCardRepository)
	repo.On("Add", *expected).Return(nil).Once()

	a := AppModeCardHandler{repo, vc}

	a.Create(&core.CreateCardRequest{})

	repo.AssertExpectations(t)
}

func TestApp_RevokeCard_RemoteReturnErr_ReturnErr(t *testing.T) {
	vc := new(fakeVirgilClient)
	vc.On("RevokeCard", mock.Anything).Return(fmt.Errorf("Error"))

	a := AppModeCardHandler{nil, vc}

	err := a.Revoke(&core.RevokeCardRequest{})
	assert.NotNil(t, err)
}

func TestApp_RevokeCard_RemoteReturnNil_DeleteFromLocal(t *testing.T) {
	vc := new(fakeVirgilClient)
	vc.On("RevokeCard", mock.Anything).Return(nil)
	repo := new(fakeCardRepository)
	repo.On("DeleteById", "id").Return(nil).Once()

	a := AppModeCardHandler{repo, vc}

	a.Revoke(&core.RevokeCardRequest{Info: virgil.RevokeCardRequest{ID: "id"}})

	repo.AssertExpectations(t)
}

func TestApp_RevokeCard_RemoteReturnNil_ReturnNil(t *testing.T) {
	expected, _ := virgil.NewRevokeCardRequest("id", virgil.RevocationReason.Compromised)
	vc := new(fakeVirgilClient)
	vc.On("RevokeCard", expected).Return(nil)
	repo := new(fakeCardRepository)
	repo.On("DeleteById", mock.Anything).Return(nil)

	a := AppModeCardHandler{repo, vc}

	err := a.Revoke(&core.RevokeCardRequest{Request: *expected})

	assert.Nil(t, err)
}
