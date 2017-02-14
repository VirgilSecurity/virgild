package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/virgil.v4"
)

func TestInitDefault_ConfigEmpty_SetDefaultValue(t *testing.T) {
	expected := Config{
		Admin: AdminConfig{
			Login:    "admin",
			Password: "8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918",
		},
		DB:      "sqlite3:virgild.db",
		LogFile: "console",
	}

	actual := initDefault(Config{})
	assert.Equal(t, expected.Admin, actual.Admin)
	assert.Equal(t, expected.DB, actual.DB)
	assert.Equal(t, expected.LogFile, actual.LogFile)
	_, err := virgil.Crypto().ImportPrivateKey([]byte(actual.Cards.Signer.PrivateKey), actual.Cards.Signer.PrivateKeyPassword)
	assert.Nil(t, err)
}

func TestInitDefault_ConfigSetRemote_SetDefaultValue(t *testing.T) {
	expected := Config{
		Cards: CardsConfig{
			Remote: RemoteConfig{
				Services: ServicesConfig{
					Cards:    "https://cards.virgilsecurity.com",
					CardsRO:  "https://cards-ro.virgilsecurity.com",
					Identity: "https://identity.virgilsecurity.com",
					VRA:      "https://ra.virgilsecurity.com",
				},
				Authority: AuthorityConfig{
					CardID:    "3e29d43373348cfb373b7eae189214dc01d7237765e572db685839b64adca853",
					PublicKey: "MCowBQYDK2VwAyEAYR501kV1tUne2uOdkw4kErRRbJrc2Syaz5V1fuG+rVs=",
				},
				Cache: 3600,
			},
		},
	}

	actual := initDefault(Config{Cards: CardsConfig{Remote: RemoteConfig{}}})
	assert.Equal(t, expected.Cards.Remote, actual.Cards.Remote)
}
