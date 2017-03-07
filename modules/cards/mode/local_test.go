package mode

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"gopkg.in/virgil.v4"

	"github.com/VirgilSecurity/virgild/modules/cards/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fakeRepo struct {
	mock.Mock
}

func (f *fakeRepo) Get(id string) (card *core.SqlCard, err error) {
	args := f.Called(id)
	card, _ = args.Get(0).(*core.SqlCard)
	err = args.Error(1)
	return
}
func (f *fakeRepo) Find(identitis []string, identityType string, scope string) (cards []core.SqlCard, err error) {
	args := f.Called(identitis, identityType, scope)
	cards, _ = args.Get(0).([]core.SqlCard)
	err = args.Error(1)
	return
}
func (f *fakeRepo) Add(cs core.SqlCard) error {
	args := f.Called(cs)
	return args.Error(0)
}
func (f *fakeRepo) MarkDeletedById(id string) error {
	args := f.Called(id)
	return args.Error(0)
}
func (f *fakeRepo) DeleteById(id string) error {
	args := f.Called(id)
	return args.Error(0)
}
func (f *fakeRepo) DeleteBySearch(identitis []string, identityType string, scope string) error {
	args := f.Called(identitis, identityType, scope)
	return args.Error(0)
}

func makeCardAndSqlCard(t *testing.T) (*core.Card, *core.SqlCard) {
	bi := make([]byte, 10)
	_, err := rand.Read(bi)
	assert.NoError(t, err)
	identity := hex.EncodeToString(bi)

	snapshot, err := json.Marshal(virgil.CardModel{
		Identity:     identity,
		IdentityType: "temp",
		PublicKey:    []byte("public key"),
		Scope:        virgil.CardScope.Application,
		Data: map[string]string{
			"test": "data",
		},
		DeviceInfo: virgil.DeviceInfo{
			Device:     "temp IoT",
			DeviceName: "temp name",
		},
	})
	assert.NoError(t, err)

	card := &core.Card{
		ID:       hex.EncodeToString(virgil.Crypto().CalculateFingerprint(snapshot)),
		Snapshot: snapshot,
		Meta: core.CardMeta{
			Signatures: map[string][]byte{
				"test": []byte("sign"),
			},
			CardVersion: "v4",
			CreatedAt:   "today",
		},
	}

	jcard, err := json.Marshal(card)
	assert.NoError(t, err)

	scard := &core.SqlCard{
		CardID:       card.ID,
		Identity:     identity,
		IdentityType: "temp",
		Scope:        "application",
		Card:         jcard,
	}
	return card, scard
}

func TestLocalGet_LocalReturnErr_ReturnErr(t *testing.T) {
	r := new(fakeRepo)
	r.On("Get", mock.Anything).Return(nil, fmt.Errorf("Error"))
	lcm := LocalCardsMiddleware{r}
	f := lcm.Get(func(id string) (*core.Card, error) {
		t.FailNow()
		return nil, nil
	})

	_, err := f("id")
	assert.Error(t, err)
}

func TestLocalGet_LocalReturnVal_ReturnVal(t *testing.T) {
	const id = "id"
	expected, scard := makeCardAndSqlCard(t)
	scard.ExpireAt = time.Now().UTC().Add(time.Hour).Unix()

	r := new(fakeRepo)
	r.On("Get", id).Return(scard, nil)
	lcm := LocalCardsMiddleware{r}
	f := lcm.Get(func(id string) (*core.Card, error) {
		t.FailNow()
		return nil, nil
	})

	actual, _ := f(id)

	assert.Equal(t, expected, actual)
}

func TestLocalGet_LocalReturnValMarkedDeleted_ReturnNotFound(t *testing.T) {
	const id = "id"
	_, scard := makeCardAndSqlCard(t)
	scard.ExpireAt = time.Now().UTC().Add(time.Hour).Unix()
	scard.Deleted = true

	r := new(fakeRepo)
	r.On("Get", id).Return(scard, nil)
	lcm := LocalCardsMiddleware{r}
	f := lcm.Get(func(id string) (*core.Card, error) {
		t.FailNow()
		return nil, nil
	})

	_, err := f(id)

	assert.Equal(t, core.ErrorEntityNotFound, err)
}

func TestLocalGet_NextReturnErr_ReturnErr(t *testing.T) {
	r := new(fakeRepo)
	r.On("Get", mock.Anything).Return(nil, core.ErrorEntityNotFound)
	lcm := LocalCardsMiddleware{r}
	f := lcm.Get(func(id string) (*core.Card, error) {
		return nil, fmt.Errorf("Error")
	})

	_, err := f("id")

	assert.Error(t, err)
}

func TestLocalGet_NextReturnVal_ReturnVal(t *testing.T) {
	expected, _ := makeCardAndSqlCard(t)

	r := new(fakeRepo)
	r.On("Get", mock.Anything).Return(nil, core.ErrorEntityNotFound)
	r.On("Add", mock.Anything).Return(nil)

	lcm := LocalCardsMiddleware{r}
	f := lcm.Get(func(id string) (*core.Card, error) {
		return expected, nil
	})

	actual, _ := f("id")

	assert.Equal(t, expected, actual)
}

func TestLocalGet_NextReturnVal_AddLocalStore(t *testing.T) {
	expected, scard := makeCardAndSqlCard(t)

	r := new(fakeRepo)
	r.On("Get", mock.Anything).Return(nil, core.ErrorEntityNotFound)
	r.On("Add", *scard).Return(nil).Once()

	lcm := LocalCardsMiddleware{r}
	f := lcm.Get(func(id string) (*core.Card, error) {
		return expected, nil
	})

	f("id")

	r.AssertExpectations(t)
}

func TestLocalGet_NextReturnValConver2SqlCardReturnErr_ReturnErr(t *testing.T) {
	card := &core.Card{
		Snapshot: []byte("snapshot"),
	}
	r := new(fakeRepo)
	r.On("Get", mock.Anything).Return(nil, core.ErrorEntityNotFound)

	lcm := LocalCardsMiddleware{r}
	f := lcm.Get(func(id string) (*core.Card, error) {
		return card, nil
	})

	_, err := f("id")

	assert.Error(t, err)
	r.AssertNotCalled(t, "Add", mock.Anything)
}

func TestLocalGet_ExpiredNextReturnErr_ReturnErr(t *testing.T) {
	_, scard := makeCardAndSqlCard(t)
	scard.ExpireAt = time.Now().UTC().Add(-time.Hour).Unix()

	r := new(fakeRepo)
	r.On("Get", mock.Anything).Return(scard, nil)
	r.On("DeleteById", mock.Anything).Return(nil)
	lcm := LocalCardsMiddleware{r}
	f := lcm.Get(func(id string) (*core.Card, error) {
		return nil, fmt.Errorf("Error")
	})

	_, err := f("id")

	assert.Error(t, err)
}

func TestLocalGet_ExpiredNextReturnVal_ReturnVal(t *testing.T) {
	expected, scard := makeCardAndSqlCard(t)
	scard.ExpireAt = time.Now().UTC().Add(-time.Hour).Unix()

	r := new(fakeRepo)
	r.On("Get", mock.Anything).Return(scard, nil)
	r.On("DeleteById", mock.Anything).Return(nil)
	r.On("Add", mock.Anything).Return(nil)

	lcm := LocalCardsMiddleware{r}
	f := lcm.Get(func(id string) (*core.Card, error) {
		return expected, nil
	})

	actual, _ := f("id")

	assert.Equal(t, expected, actual)
}

func TestLocalGet_ExpiredNextReturnVal_AddLocalStore(t *testing.T) {
	expected, scard := makeCardAndSqlCard(t)
	uscard := *scard
	scard.ExpireAt = time.Now().UTC().Add(-time.Hour).Unix()

	r := new(fakeRepo)
	r.On("Get", mock.Anything).Return(scard, nil)
	r.On("Add", uscard).Return(nil).Once()
	r.On("DeleteById", expected.ID).Return(nil).Once()

	lcm := LocalCardsMiddleware{r}
	f := lcm.Get(func(id string) (*core.Card, error) {
		return expected, nil
	})

	f(expected.ID)

	r.AssertExpectations(t)
}

func TestLocalGet_ExpiredNextReturnValConver2SqlCardReturnErr_ReturnErr(t *testing.T) {
	_, scard := makeCardAndSqlCard(t)
	scard.ExpireAt = time.Now().UTC().Add(-time.Hour).Unix()
	card := &core.Card{
		Snapshot: []byte("snapshot"),
	}
	r := new(fakeRepo)
	r.On("Get", mock.Anything).Return(scard, nil)
	r.On("DeleteById", mock.Anything).Return(nil)

	lcm := LocalCardsMiddleware{r}
	f := lcm.Get(func(id string) (*core.Card, error) {
		return card, nil
	})

	_, err := f("id")

	assert.Error(t, err)
	r.AssertNotCalled(t, "Add", mock.Anything)
}

func TestLocalSearch_LocalReturnErr_ReturnErr(t *testing.T) {
	r := new(fakeRepo)
	r.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("ERROR"))

	lcm := LocalCardsMiddleware{r}
	f := lcm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		t.FailNow()
		return nil, nil
	})

	_, err := f(&virgil.Criteria{})

	assert.Error(t, err)
}

func TestLocalSearch_LocalReturnVal_ReturnVal(t *testing.T) {
	crit := &virgil.Criteria{
		Identities:   []string{"bob", "alice"},
		IdentityType: "type",
		Scope:        virgil.CardScope.Application,
	}
	card, scard := makeCardAndSqlCard(t)
	scard.ExpireAt = time.Now().UTC().Add(time.Hour).Unix()

	r := new(fakeRepo)
	r.On("Find", crit.Identities, crit.IdentityType, string(crit.Scope)).Return([]core.SqlCard{*scard}, nil)

	lcm := LocalCardsMiddleware{r}
	f := lcm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		t.FailNow()
		return nil, nil
	})

	actual, _ := f(crit)

	assert.Equal(t, []core.Card{*card}, actual)
}

func TestLocalSearch_CardRevoked_ReturnEmpty(t *testing.T) {
	_, scard := makeCardAndSqlCard(t)
	scard.Deleted = true

	r := new(fakeRepo)
	r.On("Find", mock.Anything, mock.Anything, mock.Anything).Return([]core.SqlCard{*scard}, nil)

	lcm := LocalCardsMiddleware{r}
	f := lcm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		t.FailNow()
		return nil, nil
	})

	actual, _ := f(&virgil.Criteria{})

	assert.Empty(t, actual)
}

func TestLocalSearch_NextReturnErr_ReturnErr(t *testing.T) {
	r := new(fakeRepo)
	r.On("Find", mock.Anything, mock.Anything, mock.Anything).Return([]core.SqlCard{}, nil)

	lcm := LocalCardsMiddleware{r}
	f := lcm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		return nil, fmt.Errorf("ERROR")
	})

	_, err := f(&virgil.Criteria{})

	assert.Error(t, err)
}

func TestLocalSearch_NextReturnVal_ReturnVal(t *testing.T) {
	expcrit := &virgil.Criteria{
		Identities:   []string{"bob", "alice"},
		IdentityType: "type",
		Scope:        virgil.CardScope.Application,
	}
	expected, _ := makeCardAndSqlCard(t)

	r := new(fakeRepo)
	r.On("Find", mock.Anything, mock.Anything, mock.Anything).Return([]core.SqlCard{}, nil)
	r.On("Add", mock.Anything).Return(nil)

	lcm := LocalCardsMiddleware{r}
	f := lcm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		assert.Equal(t, expcrit, crit)
		return []core.Card{*expected}, nil
	})

	actual, _ := f(expcrit)

	assert.Equal(t, []core.Card{*expected}, actual)
}

func TestLocalSearch_NextReturnVal_AddToLocalStore(t *testing.T) {

	expected, scard := makeCardAndSqlCard(t)

	r := new(fakeRepo)
	r.On("Find", mock.Anything, mock.Anything, mock.Anything).Return([]core.SqlCard{}, nil)
	r.On("Add", *scard).Return(nil).Once()

	lcm := LocalCardsMiddleware{r}
	f := lcm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		return []core.Card{*expected}, nil
	})

	f(&virgil.Criteria{})

	r.AssertExpectations(t)
}

func TestLocalSearch_ExpiredNextReturnErr_ReturnErr(t *testing.T) {
	_, scard := makeCardAndSqlCard(t)
	scard.ExpireAt = time.Now().UTC().Add(-time.Hour).Unix()

	r := new(fakeRepo)
	r.On("Find", mock.Anything, mock.Anything, mock.Anything).Return([]core.SqlCard{*scard}, nil)
	r.On("DeleteBySearch", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	lcm := LocalCardsMiddleware{r}
	f := lcm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		return nil, fmt.Errorf("ERROR")
	})

	_, err := f(&virgil.Criteria{})

	assert.Error(t, err)
}

func TestLocalSearch_ExpiredRemoveExpiredData_ExecuteDeleteBySearch(t *testing.T) {
	crit := &virgil.Criteria{
		Identities:   []string{"bob", "alice"},
		IdentityType: "type",
		Scope:        virgil.CardScope.Application,
	}
	_, scard := makeCardAndSqlCard(t)
	scard.ExpireAt = time.Now().UTC().Add(-time.Hour).Unix()

	r := new(fakeRepo)
	r.On("Find", mock.Anything, mock.Anything, mock.Anything).Return([]core.SqlCard{*scard}, nil)
	r.On("DeleteBySearch", crit.Identities, crit.IdentityType, string(crit.Scope)).Return(nil).Once()

	lcm := LocalCardsMiddleware{r}
	f := lcm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		return nil, fmt.Errorf("ERROR")
	})

	f(crit)

	r.AssertExpectations(t)
}

func TestLocalSearch_ExpiredNextReturnVal_ReturnVal(t *testing.T) {
	expcrit := &virgil.Criteria{
		Identities:   []string{"bob", "alice"},
		IdentityType: "type",
		Scope:        virgil.CardScope.Application,
	}
	expected, scard := makeCardAndSqlCard(t)
	scard.ExpireAt = time.Now().UTC().Add(-time.Hour).Unix()

	r := new(fakeRepo)
	r.On("Find", mock.Anything, mock.Anything, mock.Anything).Return([]core.SqlCard{*scard}, nil)
	r.On("DeleteBySearch", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	r.On("Add", mock.Anything).Return(nil)

	lcm := LocalCardsMiddleware{r}
	f := lcm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		assert.Equal(t, expcrit, crit)
		return []core.Card{*expected}, nil
	})

	actual, _ := f(expcrit)

	assert.Equal(t, []core.Card{*expected}, actual)
}

func TestLocalSearch_ExpiredNextReturnVal_AddToLocalStore(t *testing.T) {

	expected, scard := makeCardAndSqlCard(t)
	uscard := *scard
	scard.ExpireAt = time.Now().UTC().Add(-time.Hour).Unix()

	r := new(fakeRepo)
	r.On("Find", mock.Anything, mock.Anything, mock.Anything).Return([]core.SqlCard{*scard}, nil)
	r.On("DeleteBySearch", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	r.On("Add", uscard).Return(nil).Once()

	lcm := LocalCardsMiddleware{r}
	f := lcm.Search(func(crit *virgil.Criteria) ([]core.Card, error) {
		return []core.Card{*expected}, nil
	})

	f(&virgil.Criteria{})

	r.AssertExpectations(t)
}

func TestLocalCreate_NextReturnErr_ReturnErr(t *testing.T) {
	lcm := LocalCardsMiddleware{}
	f := lcm.Create(func(req *core.CreateCardRequest) (*core.Card, error) {
		return nil, fmt.Errorf("Error")
	})
	_, err := f(&core.CreateCardRequest{})

	assert.Error(t, err)
}

func TestLocalCreate_NextReturnValConver2SqlCardReturnErr_ReturnErr(t *testing.T) {
	card := &core.Card{
		Snapshot: []byte("snapshot"),
	}
	lcm := LocalCardsMiddleware{}
	f := lcm.Create(func(req *core.CreateCardRequest) (*core.Card, error) {
		return card, nil
	})
	_, err := f(&core.CreateCardRequest{})

	assert.Error(t, err)
}

func TestLocalCreate_NextReturnVal_ReturnVal(t *testing.T) {
	expected, _ := makeCardAndSqlCard(t)

	m := new(fakeRepo)
	m.On("Add", mock.Anything).Return(nil)
	lcm := LocalCardsMiddleware{m}
	f := lcm.Create(func(req *core.CreateCardRequest) (*core.Card, error) {
		return expected, nil
	})
	actual, _ := f(&core.CreateCardRequest{})

	assert.Equal(t, expected, actual)
}

func TestLocalCreate_NextReturnVal_AddToLocalStore(t *testing.T) {
	expected, scard := makeCardAndSqlCard(t)

	m := new(fakeRepo)
	m.On("Add", *scard).Return(nil).Once()
	lcm := LocalCardsMiddleware{m}
	f := lcm.Create(func(req *core.CreateCardRequest) (*core.Card, error) {
		return expected, nil
	})
	f(&core.CreateCardRequest{})

	m.AssertExpectations(t)
}

func TestLocalRevoke_NextReturnErr_ReturnErr(t *testing.T) {
	lcm := LocalCardsMiddleware{}
	f := lcm.Revoke(func(req *core.RevokeCardRequest) error {
		return fmt.Errorf("Error")
	})
	err := f(&core.RevokeCardRequest{})

	assert.Error(t, err)
}

func TestLocalRevoke_NextReturnNil_ReturnNil(t *testing.T) {
	m := new(fakeRepo)
	m.On("MarkDeletedById", mock.Anything).Return(nil)
	lcm := LocalCardsMiddleware{m}
	f := lcm.Revoke(func(req *core.RevokeCardRequest) error {
		return nil
	})
	err := f(&core.RevokeCardRequest{})

	assert.NoError(t, err)
}

func TestLocalRevoke_NextReturnNil_MarkCardDeleted(t *testing.T) {
	const id = "id"
	m := new(fakeRepo)
	m.On("MarkDeletedById", id).Return(nil).Once()
	lcm := LocalCardsMiddleware{m}
	f := lcm.Revoke(func(req *core.RevokeCardRequest) error {
		return nil
	})
	f(&core.RevokeCardRequest{Info: virgil.RevokeCardRequest{
		ID: id,
	}})

	m.AssertExpectations(t)
}
