package handler

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/VirgilSecurity/virgild/modules/symmetric/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

type fakeSymmetricRepo struct {
	mock.Mock
}

func (f *fakeSymmetricRepo) Create(k core.SymmetricKey) error {
	args := f.Called(k)
	return args.Error(0)
}
func (f *fakeSymmetricRepo) Remove(keyID, userID string) error {
	args := f.Called(keyID, userID)
	return args.Error(0)
}
func (f *fakeSymmetricRepo) Get(keyID, userID string) (k *core.SymmetricKey, err error) {
	args := f.Called(keyID, userID)
	k, _ = args.Get(0).(*core.SymmetricKey)
	err = args.Error(1)
	return
}
func (f *fakeSymmetricRepo) KeysByUser(userID string) (ks []core.KeyUserPair, err error) {
	args := f.Called(userID)
	ks, _ = args.Get(0).([]core.KeyUserPair)
	err = args.Error(1)
	return
}
func (f *fakeSymmetricRepo) UsersByKey(keyID string) (ks []core.KeyUserPair, err error) {
	args := f.Called(keyID)
	ks, _ = args.Get(0).([]core.KeyUserPair)
	err = args.Error(1)
	return
}

func makeRequestCtx(body interface{}) *fasthttp.RequestCtx {
	res := &fasthttp.RequestCtx{
		Request: fasthttp.Request{
			Header: fasthttp.RequestHeader{},
		},
		Response: fasthttp.Response{},
	}

	switch body.(type) {
	case []byte:
		res.Request.AppendBody(body.([]byte))
	case string:
		res.Request.AppendBodyString(body.(string))
	case nil:

	default:
		b, _ := json.Marshal(body)
		res.Request.AppendBody(b)

	}

	return res
}

func TestGetKey_RepoReturnErr_ReturnErr(t *testing.T) {
	repo := new(fakeSymmetricRepo)
	repo.On("Get", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("ERROR"))

	ctx := fasthttp.RequestCtx{}
	ctx.SetUserValue("key_id", "1")
	ctx.SetUserValue("user_id", "1")

	f := GetKey(repo)
	_, err := f(&ctx)

	assert.NotNil(t, err)
}

func TestGetKey_RepoReturnVal_ReturnVal(t *testing.T) {
	expected := &core.SymmetricKey{
		UserID:       "user_id",
		KeyID:        "key_id",
		EncryptedKey: []byte("test"),
	}
	repo := new(fakeSymmetricRepo)
	repo.On("Get", "key_id", "user_id").Return(expected, nil)

	ctx := fasthttp.RequestCtx{}
	ctx.SetUserValue("key_id", "key_id")
	ctx.SetUserValue("user_id", "user_id")

	f := GetKey(repo)
	actual, _ := f(&ctx)

	assert.Equal(t, expected, actual)
}

func TestCreateKey_JSONInvalid_ReturnErr(t *testing.T) {
	ctx := makeRequestCtx([]byte("broken JSON,."))

	f := CreateKey(nil)
	_, err := f(ctx)

	assert.Equal(t, core.ErrorJSONIsInvalid, err)
}

func TestCreateKey_RepoReturnErr_ReturnErr(t *testing.T) {
	repo := new(fakeSymmetricRepo)
	repo.On("Create", mock.Anything).Return(fmt.Errorf("ERROR"))

	ctx := makeRequestCtx(core.SymmetricKey{KeyID: "1234", UserID: "1234", EncryptedKey: []byte("test")})

	f := CreateKey(repo)
	_, err := f(ctx)

	assert.NotNil(t, err)
}

func TestCreateKey_RepoReturnVal_ReturnVal(t *testing.T) {
	k := core.SymmetricKey{
		UserID:       "user_id",
		KeyID:        "key_id",
		EncryptedKey: []byte("test"),
	}
	repo := new(fakeSymmetricRepo)
	repo.On("Create", k).Return(nil)

	ctx := makeRequestCtx(k)

	f := CreateKey(repo)
	result, err := f(ctx)

	assert.Nil(t, err)
	assert.Nil(t, result)
}

func TestGetUsersByKey_RepoReturnErr_ReturnErr(t *testing.T) {
	repo := new(fakeSymmetricRepo)
	repo.On("UsersByKey", mock.Anything).Return(nil, fmt.Errorf("ERROR"))

	ctx := &fasthttp.RequestCtx{}
	ctx.SetUserValue("key_id", "id")

	f := GetUsersByKey(repo)
	_, err := f(ctx)

	assert.NotNil(t, err)
}

func TestGetUsersByKey_RepoReturnVal_ReturnVal(t *testing.T) {
	expected := []core.KeyUserPair{core.KeyUserPair{
		UserID: "user_id",
		KeyID:  "key_id",
	}}
	repo := new(fakeSymmetricRepo)
	repo.On("UsersByKey", "id").Return(expected, nil)

	ctx := &fasthttp.RequestCtx{}
	ctx.SetUserValue("key_id", "id")

	f := GetUsersByKey(repo)
	actual, err := f(ctx)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestGetKeysByUser_RepoReturnErr_ReturnErr(t *testing.T) {
	repo := new(fakeSymmetricRepo)
	repo.On("KeysByUser", mock.Anything).Return(nil, fmt.Errorf("ERROR"))

	ctx := &fasthttp.RequestCtx{}
	ctx.SetUserValue("user_id", "id")

	f := GetKeysByUser(repo)
	_, err := f(ctx)

	assert.NotNil(t, err)
}

func TestGetKeysByUser_RepoReturnVal_ReturnVal(t *testing.T) {
	expected := []core.KeyUserPair{core.KeyUserPair{
		UserID: "user_id",
		KeyID:  "key_id",
	}}
	repo := new(fakeSymmetricRepo)
	repo.On("KeysByUser", "id").Return(expected, nil)

	ctx := &fasthttp.RequestCtx{}
	ctx.SetUserValue("user_id", "id")

	f := GetKeysByUser(repo)
	actual, err := f(ctx)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
