package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDB_MissColon_ReturnErr(t *testing.T) {
	_, err := initDB("virgild.db")
	assert.NotNil(t, err)
}

func TestInitDB_DriverNotFound_ReturnErr(t *testing.T) {
	_, err := initDB("test:virgild.db")
	assert.NotNil(t, err)
}

func TestInitDB_ReturnOrm(t *testing.T) {
	orm, err := initDB("sqlite3:")
	assert.NotNil(t, orm)
	assert.Nil(t, err)
}

func TestInitVRA_Nil_ReturnNil(t *testing.T) {
	table := []AuthorityConfig{
		AuthorityConfig{CardID: "ad"},
		AuthorityConfig{PublicKey: "ad"},
		AuthorityConfig{},
	}

	for _, v := range table {
		vra, err := initVRA(v)
		assert.Nil(t, vra)
		assert.Nil(t, err)
	}
}

func TestInitVRA_PublicKeyInvalid_ReturnErr(t *testing.T) {
	_, err := initVRA(AuthorityConfig{CardID: "asdf", PublicKey: "asdf"})
	assert.NotNil(t, err)
}

func TestInitVRA_ReturnVRA(t *testing.T) {
	vra, err := initVRA(AuthorityConfig{CardID: "3e29d43373348cfb373b7eae189214dc01d7237765e572db685839b64adca853", PublicKey: "MCowBQYDK2VwAyEAYR501kV1tUne2uOdkw4kErRRbJrc2Syaz5V1fuG+rVs="})
	assert.Nil(t, err)
	assert.NotNil(t, vra)
}

func TestInitSigner_PrivateKeyInvalid_ReturnErr(t *testing.T) {
	_, err := initSigner(&SignerConfig{})
	assert.Nil(t, err)
}

func TestInitSigner_CardIdEmpty_ReturnCard(t *testing.T) {
	signer, err := initSigner(&SignerConfig{PrivateKey: "MC4CAQAwBQYDK2VwBCIEICRiZ4/yefgZHLqsUULmmvv798Jd+7P9kaKG512LbMjc"})
	assert.Nil(t, err)
	assert.NotNil(t, signer.Card)
	assert.Equal(t, signer.CardID, signer.Card.ID)
}

func TestInitSigner_ReturnSigner(t *testing.T) {
	signer, err := initSigner(&SignerConfig{CardID: "123", PrivateKey: "MC4CAQAwBQYDK2VwBCIEICRiZ4/yefgZHLqsUULmmvv798Jd+7P9kaKG512LbMjc"})
	assert.Nil(t, err)
	assert.Nil(t, signer.Card)
	assert.Equal(t, "123", signer.CardID)
}

func TestInitRemote_PublicKeyInvalid_ReturnErr(t *testing.T) {
	_, err := initRemote(RemoteConfig{Authority: AuthorityConfig{PublicKey: "+rVs="}})
	assert.NotNil(t, err)
}

func TestInitRemote_ReturnRemote(t *testing.T) {
	remote, err := initRemote(RemoteConfig{Authority: AuthorityConfig{PublicKey: "MCowBQYDK2VwAyEAYR501kV1tUne2uOdkw4kErRRbJrc2Syaz5V1fuG+rVs="}})
	assert.Nil(t, err)
	assert.NotNil(t, remote)
}

func TestSetSiteVirgilD_Return(t *testing.T) {
	expected := VirgilDCard{
		CardID:    "123",
		PublicKey: "MCowBQYDK2VwAyEAZCptOxevuISeK4rlqSkQF9gCi7jXicA5Rtl2pwKLmO0=",
	}
	signer, err := initSigner(&SignerConfig{
		CardID:     "123",
		PrivateKey: "MC4CAQAwBQYDK2VwBCIEICRiZ4/yefgZHLqsUULmmvv798Jd+7P9kaKG512LbMjc",
	})
	assert.Nil(t, err)
	actual, err := setSiteVirgilD(signer)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
