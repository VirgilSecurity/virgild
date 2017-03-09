package mode

import (
	"fmt"
	"testing"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	virgil "gopkg.in/virgil.v4"
)

type fakeClient struct {
	mock.Mock
}

func (f *fakeClient) GetCard(id string) (card *virgil.Card, err error) {
	args := f.Called(id)
	card, _ = args.Get(0).(*virgil.Card)
	err = args.Error(1)
	return
}
func (f *fakeClient) SearchCards(crit *virgil.Criteria) (cards []*virgil.Card, err error) {
	args := f.Called(crit)
	cards, _ = args.Get(0).([]*virgil.Card)
	err = args.Error(1)
	return
}
func (f *fakeClient) CreateCard(req *virgil.SignableRequest) (cards *virgil.Card, err error) {
	args := f.Called(req)
	cards, _ = args.Get(0).(*virgil.Card)
	err = args.Error(1)
	return
}
func (f *fakeClient) RevokeCard(req *virgil.SignableRequest) error {
	args := f.Called(req)
	return args.Error(0)
}

func TestRemoteGet_RemoteReturnErr_ReturnErr(t *testing.T) {
	c := new(fakeClient)
	c.On("GetCard", mock.Anything).Return(nil, fmt.Errorf("Error"))
	rcm := RemoteCardsMiddleware{c}

	_, err := rcm.Get("id")

	assert.Error(t, err)
}

func TestRemoteGet_RemoteReturnVal_ReturnVal(t *testing.T) {
	const id = "123"
	vcard := &virgil.Card{
		ID:          id,
		Snapshot:    []byte("snapshot"),
		CardVersion: "v4",
		Signatures: map[string][]byte{
			"123": []byte("sign"),
		},
		CreatedAt: "today",
	}
	expected := &core.Card{
		ID:       id,
		Snapshot: vcard.Snapshot,
		Meta: core.CardMeta{
			CardVersion: vcard.CardVersion,
			CreatedAt:   vcard.CreatedAt,
			Signatures:  vcard.Signatures,
		},
	}

	c := new(fakeClient)
	c.On("GetCard", id).Return(vcard, nil)
	rcm := RemoteCardsMiddleware{c}

	actual, _ := rcm.Get(id)

	assert.Equal(t, expected, actual)
}

func TestRemoteSearch_RemoteReturnErr_ReturnErr(t *testing.T) {
	c := new(fakeClient)
	c.On("SearchCards", mock.Anything).Return(nil, fmt.Errorf("Error"))
	rcm := RemoteCardsMiddleware{c}

	_, err := rcm.Search(&virgil.Criteria{})

	assert.Error(t, err)
}

func TestRemoteSearch_RemoteReturnVal_ReturnVal(t *testing.T) {
	vcard := []*virgil.Card{
		&virgil.Card{
			ID:          "1124",
			Snapshot:    []byte("snapshot"),
			CardVersion: "v4",
			Signatures: map[string][]byte{
				"123": []byte("sign"),
			},
			CreatedAt: "today",
		},
		&virgil.Card{
			ID:          "4321",
			Snapshot:    []byte("snapshot 2"),
			CardVersion: "v4",
			Signatures: map[string][]byte{
				"123": []byte("new sign"),
			},
			CreatedAt: "today",
		},
	}
	expected := []core.Card{
		core.Card{
			ID:       vcard[0].ID,
			Snapshot: vcard[0].Snapshot,
			Meta: core.CardMeta{
				CardVersion: vcard[0].CardVersion,
				CreatedAt:   vcard[0].CreatedAt,
				Signatures:  vcard[0].Signatures,
			},
		},
		core.Card{
			ID:       vcard[1].ID,
			Snapshot: vcard[1].Snapshot,
			Meta: core.CardMeta{
				CardVersion: vcard[1].CardVersion,
				CreatedAt:   vcard[1].CreatedAt,
				Signatures:  vcard[1].Signatures,
			},
		},
	}

	crit := virgil.SearchCriteriaByAppBundle("asdf")
	c := new(fakeClient)
	c.On("SearchCards", crit).Return(vcard, nil)
	rcm := RemoteCardsMiddleware{c}

	actual, _ := rcm.Search(crit)

	assert.Equal(t, expected, actual)
}

func TestRemoteCreate_RemoteReturnErr_ReturnErr(t *testing.T) {
	c := new(fakeClient)
	c.On("CreateCard", mock.Anything).Return(nil, fmt.Errorf("Error"))
	rcm := RemoteCardsMiddleware{c}

	_, err := rcm.Create(&core.CreateCardRequest{})

	assert.Error(t, err)
}

func TestRemoteCreate_RemoteReturnVal_ReturnVal(t *testing.T) {
	vcard := &virgil.Card{
		ID:          "1124",
		Snapshot:    []byte("snapshot"),
		CardVersion: "v4",
		Signatures: map[string][]byte{
			"123": []byte("sign"),
		},
		CreatedAt: "today",
	}
	expected := &core.Card{
		ID:       vcard.ID,
		Snapshot: vcard.Snapshot,
		Meta: core.CardMeta{
			CardVersion: vcard.CardVersion,
			CreatedAt:   vcard.CreatedAt,
			Signatures:  vcard.Signatures,
		},
	}

	req := virgil.SignableRequest{
		Snapshot: vcard.Snapshot,
		Meta: virgil.RequestMeta{
			Signatures: vcard.Signatures,
		},
	}

	c := new(fakeClient)
	c.On("CreateCard", &req).Return(vcard, nil)
	rcm := RemoteCardsMiddleware{c}

	actual, _ := rcm.Create(&core.CreateCardRequest{
		Info:    virgil.CardModel{},
		Request: req,
	})

	assert.Equal(t, expected, actual)
}

func TestRemoteRevoke_RemoteReturnErr_ReturnErr(t *testing.T) {
	c := new(fakeClient)
	c.On("RevokeCard", mock.Anything).Return(fmt.Errorf("Error"))
	rcm := RemoteCardsMiddleware{c}

	err := rcm.Revoke(&core.RevokeCardRequest{})

	assert.Error(t, err)
}

func TestRemoteRevoke_RemoteReturnNil_ReturnNil(t *testing.T) {
	req := virgil.SignableRequest{
		Snapshot: []byte("snapshot"),
		Meta: virgil.RequestMeta{
			Signatures: map[string][]byte{
				"test": []byte("sign"),
			},
		},
	}

	c := new(fakeClient)
	c.On("RevokeCard", &req).Return(nil)
	rcm := RemoteCardsMiddleware{c}

	err := rcm.Revoke(&core.RevokeCardRequest{
		Info:    virgil.RevokeCardRequest{},
		Request: req,
	})

	assert.NoError(t, err)
}
