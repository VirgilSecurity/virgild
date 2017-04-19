package validator

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"gopkg.in/virgil.v4"

	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/stretchr/testify/assert"
)

func fakeCreateCardReuest(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
	return nil, nil
}

func TestCreateCard_IdentityEmpty_ReturnErr(t *testing.T) {
	req := new(core.CreateCardRequest)
	c := CreateCard(fakeCreateCardReuest)
	_, err := c(context.Background(), req)

	assert.Equal(t, core.CardIdentityEmptyErr, err)
}

func TestCreateCard_GlobalCardIdentityTypeIsNotEmail_ReturnErr(t *testing.T) {
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Identity:     "test",
			IdentityType: "test",
			Scope:        virgil.CardScope.Global,
		},
	}
	c := CreateCard(fakeCreateCardReuest)
	_, err := c(context.Background(), req)

	assert.Equal(t, core.GlobalCardIdentityTypeMustBeEmailErr, err)
}

func TestCreateCard_IdentityTypeEmailIdentityIsNot_ReturnErr(t *testing.T) {
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Identity:     "test",
			IdentityType: "email",
		},
	}
	c := CreateCard(fakeCreateCardReuest)
	_, err := c(context.Background(), req)

	assert.Equal(t, core.EmailIdentityIvalidErr, err)
}

func TestCreateCard_ApplicationCardIdentityTypeIsNotEmail_Skeep(t *testing.T) {
	kp, _ := virgil.Crypto().GenerateKeypair()
	snapshot := []byte("test")
	id := virgil.Crypto().CalculateFingerprint(snapshot)
	sign, _ := virgil.Crypto().Sign(id, kp.PrivateKey())

	epub, _ := kp.PublicKey().Encode()
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Identity:     "test",
			IdentityType: "test",
			Scope:        virgil.CardScope.Application,
			PublicKey:    epub,
		},
		Request: virgil.SignableRequest{
			Snapshot: snapshot,
			Meta: virgil.RequestMeta{
				Signatures: map[string][]byte{
					hex.EncodeToString(id): sign,
				},
			},
		},
	}
	c := CreateCard(fakeCreateCardReuest)
	_, err := c(context.Background(), req)

	assert.Nil(t, err)
}

func TestCreateCard_KeyLengthIvalid_ReturnErr(t *testing.T) {
	table := []int{15, 4000, 2049}

	for _, l := range table {
		req := &core.CreateCardRequest{
			Info: virgil.CardModel{
				Identity:     "test",
				IdentityType: "test",
				Scope:        virgil.CardScope.Application,
				PublicKey:    make([]byte, l),
			},
		}
		c := CreateCard(fakeCreateCardReuest)
		_, err := c(context.Background(), req)

		assert.Equal(t, core.PublicKeyLentghInvalidErr, err, fmt.Sprintf("Len:%v", l))
	}
}

func TestCreateCard_DataEntiesGreaterThan16_ReturnErr(t *testing.T) {
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Identity:     "test",
			IdentityType: "test",
			Scope:        virgil.CardScope.Application,
			PublicKey:    make([]byte, 1024),
			Data:         make(map[string]string),
		},
	}

	for i := 0; i < 20; i++ {
		req.Info.Data[strconv.Itoa(i)] = strconv.Itoa(i)
	}
	c := CreateCard(fakeCreateCardReuest)
	_, err := c(context.Background(), req)

	assert.Equal(t, core.CardDataCannotContainsMoreThan16EntriesErr, err)
}

func TestCreateCard_DataEntiesExceed256_ReturnErr(t *testing.T) {
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Identity:     "test",
			IdentityType: "test",
			Scope:        virgil.CardScope.Application,
			PublicKey:    make([]byte, 1024),
			Data:         make(map[string]string),
		},
	}
	var data [300]byte
	rand.Read(data[:])
	req.Info.Data["test"] = hex.EncodeToString(data[:])

	c := CreateCard(fakeCreateCardReuest)
	_, err := c(context.Background(), req)

	assert.Equal(t, core.DataValueExceed256Err, err)
}

func TestCreateCard_InfoDeviceExceed256_ReturnErr(t *testing.T) {
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Identity:     "test",
			IdentityType: "test",
			Scope:        virgil.CardScope.Application,
			PublicKey:    make([]byte, 1024),
			Data:         make(map[string]string),
		},
	}
	var data [300]byte
	rand.Read(data[:])
	req.Info.DeviceInfo.Device = hex.EncodeToString(data[:])

	c := CreateCard(fakeCreateCardReuest)
	_, err := c(context.Background(), req)

	assert.Equal(t, core.InfoValueExceed256Err, err)
}

func TestCreateCard_InfoDeviceNameExceed256_ReturnErr(t *testing.T) {
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Identity:     "test",
			IdentityType: "test",
			Scope:        virgil.CardScope.Application,
			PublicKey:    make([]byte, 1024),
			Data:         make(map[string]string),
		},
	}
	var data [300]byte
	rand.Read(data[:])
	req.Info.DeviceInfo.DeviceName = hex.EncodeToString(data[:])

	c := CreateCard(fakeCreateCardReuest)
	_, err := c(context.Background(), req)

	assert.Equal(t, core.InfoValueExceed256Err, err)
}

func TestCreateCard_SignEmpty_ReturnErr(t *testing.T) {
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Identity:     "test",
			IdentityType: "test",
			Scope:        virgil.CardScope.Application,
			PublicKey:    make([]byte, 1024),
			Data:         make(map[string]string),
		},
	}

	c := CreateCard(fakeCreateCardReuest)
	_, err := c(context.Background(), req)

	assert.Equal(t, core.SignsIsEmptyErr, err)
}

func TestCreateCard_SelfSignMissing_ReturnErr(t *testing.T) {
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Identity:     "test",
			IdentityType: "test",
			Scope:        virgil.CardScope.Application,
			PublicKey:    make([]byte, 1024),
			Data:         make(map[string]string),
		},
		Request: virgil.SignableRequest{
			Snapshot: []byte("test"),
			Meta: virgil.RequestMeta{
				Signatures: map[string][]byte{
					"test": []byte("test"),
				},
			},
		},
	}

	c := CreateCard(fakeCreateCardReuest)
	_, err := c(context.Background(), req)

	assert.Equal(t, core.SignItemInvalidForClientErr, err)
}

func TestCreateCard_SelfSignPublicKeyIncorrect_ReturnErr(t *testing.T) {
	snapshot := []byte("test")
	id := hex.EncodeToString(virgil.Crypto().CalculateFingerprint(snapshot))
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Identity:     "test",
			IdentityType: "test",
			Scope:        virgil.CardScope.Application,
			PublicKey:    make([]byte, 1024),
			Data:         make(map[string]string),
		},
		Request: virgil.SignableRequest{
			Snapshot: snapshot,
			Meta: virgil.RequestMeta{
				Signatures: map[string][]byte{
					id: []byte("test"),
				},
			},
		},
	}

	c := CreateCard(fakeCreateCardReuest)
	_, err := c(context.Background(), req)

	assert.Equal(t, core.SnapshotIncorrectErr, err)
}

func TestCreateCard_SelfSignSignInvalid_ReturnErr(t *testing.T) {
	snapshot := []byte("test")
	id := hex.EncodeToString(virgil.Crypto().CalculateFingerprint(snapshot))
	kp, _ := virgil.Crypto().GenerateKeypair()
	epub, _ := kp.PublicKey().Encode()
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Identity:     "test",
			IdentityType: "test",
			Scope:        virgil.CardScope.Application,
			PublicKey:    epub,
			Data:         make(map[string]string),
		},
		Request: virgil.SignableRequest{
			Snapshot: snapshot,
			Meta: virgil.RequestMeta{
				Signatures: map[string][]byte{
					id: []byte("test"),
				},
			},
		},
	}

	c := CreateCard(fakeCreateCardReuest)
	_, err := c(context.Background(), req)

	assert.Equal(t, core.SignItemInvalidForClientErr, err)
}

func TestCreateCard_CustomValidator_Executed(t *testing.T) {
	kp, _ := virgil.Crypto().GenerateKeypair()
	snapshot := []byte("test")
	id := virgil.Crypto().CalculateFingerprint(snapshot)
	sign, _ := virgil.Crypto().Sign(id, kp.PrivateKey())

	epub, _ := kp.PublicKey().Encode()
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Identity:     "test",
			IdentityType: "test",
			Scope:        virgil.CardScope.Application,
			PublicKey:    epub,
		},
		Request: virgil.SignableRequest{
			Snapshot: snapshot,
			Meta: virgil.RequestMeta{
				Signatures: map[string][]byte{
					hex.EncodeToString(id): sign,
				},
			},
		},
	}
	var executed bool
	c := CreateCard(fakeCreateCardReuest, func(ctx context.Context, req *core.CreateCardRequest) (bool, error) {
		executed = true
		return false, nil
	})
	c(context.Background(), req)

	assert.True(t, executed)
}

func TestCreateCard_RequestValid_NextFuncExecuted(t *testing.T) {
	kp, _ := virgil.Crypto().GenerateKeypair()
	snapshot := []byte("test")
	id := virgil.Crypto().CalculateFingerprint(snapshot)
	sign, _ := virgil.Crypto().Sign(id, kp.PrivateKey())

	epub, _ := kp.PublicKey().Encode()
	req := &core.CreateCardRequest{
		Info: virgil.CardModel{
			Identity:     "test",
			IdentityType: "test",
			Scope:        virgil.CardScope.Application,
			PublicKey:    epub,
		},
		Request: virgil.SignableRequest{
			Snapshot: snapshot,
			Meta: virgil.RequestMeta{
				Signatures: map[string][]byte{
					hex.EncodeToString(id): sign,
				},
			},
		},
	}
	var executed bool
	c := CreateCard(func(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
		executed = true
		return nil, nil
	})
	c(context.Background(), req)

	assert.True(t, executed)
}

func TestWrapCreateValidateVRASign_FuncExecuted(t *testing.T) {
	var executed bool
	wrap := WrapCreateValidateVRASign(func(req *virgil.SignableRequest) (bool, error) {
		executed = true
		return true, fmt.Errorf("Error")
	})
	ok, err := wrap(context.Background(), new(core.CreateCardRequest))

	assert.True(t, executed)
	assert.True(t, ok)
	assert.NotNil(t, err)
}
