// +build docker

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/VirgilSecurity/virgild/config"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"gopkg.in/virgil.v4"
	"gopkg.in/virgil.v4/transport/virgilhttp"
)

var host = "http://localhost:8080"

func TestCacheGetCard(t *testing.T) {
	eSnap := []byte(`{"public_key":"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUNvd0JRWURLMlZ3QXlFQVlSNTAxa1YxdFVuZTJ1T2RrdzRrRXJSUmJKcmMyU3lhejVWMWZ1RytyVnM9Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo=","identity":"com.virgilsecurity.cards","identity_type":"application","scope":"global","data":null}`)
	eCreated := "2016-09-30T10:22:17+0000"
	eVersion := "3.0"

	v, err := virgil.NewClient("", virgil.ClientTransport(virgilhttp.NewTransportClient(host, host, host, host)))
	assert.NoError(t, err, "Cannot create virgil client")
	c, err := v.GetCard("3e29d43373348cfb373b7eae189214dc01d7237765e572db685839b64adca853")
	assert.NoError(t, err, "Cannot get card(3e29d43373348cfb373b7eae189214dc01d7237765e572db685839b64adca853)")

	assert.EqualValues(t, eSnap, c.Snapshot, "Snapshots are not equal")
	assert.EqualValues(t, eCreated, c.CreatedAt, "Created dates are not equal")
	assert.EqualValues(t, eVersion, c.CardVersion, "Versions are not equal")
}

func TestCacheSearchCards(t *testing.T) {
	eSnap := []byte(`{"public_key":"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUNvd0JRWURLMlZ3QXlFQVlSNTAxa1YxdFVuZTJ1T2RrdzRrRXJSUmJKcmMyU3lhejVWMWZ1RytyVnM9Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo=","identity":"com.virgilsecurity.cards","identity_type":"application","scope":"global","data":null}`)
	eCreated := "2016-09-30T10:22:17+0000"
	eVersion := "3.0"

	v, err := virgil.NewClient("", virgil.ClientTransport(virgilhttp.NewTransportClient(host, host, host, host)))
	assert.NoError(t, err, "Cannot create virgil client")
	cs, err := v.SearchCards(virgil.SearchCriteriaByAppBundle("com.virgilsecurity.cards"))
	assert.NoError(t, err, "Cannot search cards by com.virgilsecurity.cards")
	assert.Len(t, cs, 1, "Number of found cards are not equal")
	c := cs[0]

	assert.EqualValues(t, eSnap, c.Snapshot, "Snapshots are not equall")
	assert.EqualValues(t, eCreated, c.CreatedAt, "Created dates are not equal")
	assert.EqualValues(t, eVersion, c.CardVersion, "Versions are not equal")
}

func TestSyncGetCard(t *testing.T) {
	appID := os.Getenv("SYNC_APP_ID")
	token := os.Getenv("SYNC_TOKEN")

	lc, err := virgil.NewClient("", virgil.ClientTransport(virgilhttp.NewTransportClient(host, host, host, host)))
	assert.NoError(t, err, "Cannot create virgil client")
	vc, err := virgil.NewClient(token)
	assert.NoError(t, err, "Cannot create virgil client")

	expected, err := vc.GetCard(appID)
	assert.NoError(t, err)

	actual, err := lc.GetCard(appID)
	assert.NoError(t, err)

	assert.EqualValues(t, expected.Snapshot, actual.Snapshot, "Snapshots are not equal")
	assert.EqualValues(t, expected.CreatedAt, actual.CreatedAt, "Created dates are not equal")
	assert.EqualValues(t, expected.CardVersion, actual.CardVersion, "Versions are not equal")
}

func TestSyncSearchCards(t *testing.T) {
	appID := os.Getenv("SYNC_APP_ID")
	token := os.Getenv("SYNC_TOKEN")
	appPrivateKey := os.Getenv("SYNC_APP_KEY")
	appPrivateKeyPass := os.Getenv("SYNC_APP_KEY_PASS")
	priv, err := virgil.Crypto().ImportPrivateKey([]byte(appPrivateKey), appPrivateKeyPass)
	assert.NoError(t, err, "Cannot import private key")

	deviceKeypair, err := virgil.Crypto().GenerateKeypair()
	assert.NoError(t, err, "Cannot generate key pair")

	uid := uuid.Must(uuid.NewV4()).String()
	req, err := virgil.NewCreateCardRequest(uid, "temp", deviceKeypair.PublicKey(), virgil.CardParams{
		Scope: virgil.CardScope.Application,
		Data: map[string]string{
			"os": "macOS",
		},
		DeviceInfo: virgil.DeviceInfo{
			Device:     "iphone7",
			DeviceName: "my iphone",
		},
	})
	assert.NoError(t, err, "Cannot create card request")

	signer := virgil.RequestSigner{}
	err = signer.SelfSign(req, deviceKeypair.PrivateKey())
	assert.NoError(t, err, "Cannot self sign")
	signer.AuthoritySign(req, appID, priv)

	vc, err := virgil.NewClient(token)
	assert.NoError(t, err, "Cannot create virgil client")

	card, err := vc.CreateCard(req)
	assert.NoError(t, err, "Cannot create card in the cloud")

	// HOT FIX
	if card.Relations == nil {
		card.Relations = make(map[string][]byte)
	}

	time.Sleep(5 * time.Second)

	defer func() {
		req, _ := virgil.NewRevokeCardRequest(card.ID, virgil.RevocationReason.Unspecified)
		signer.AuthoritySign(req, appID, priv)
		vc.RevokeCard(req)
	}()

	lc, err := virgil.NewClient("", virgil.ClientTransport(virgilhttp.NewTransportClient(host, host, host, host)))
	assert.NoError(t, err, "Cannot create virgil client")
	cs, err := lc.SearchCards(virgil.SearchCriteriaByIdentities(uid))
	assert.NoError(t, err, "Cannot search cards by temp name (%v)", uid)
	assert.Len(t, cs, 1, "Number of found cards are not equal")
	c := cs[0]

	assert.EqualValues(t, card, c)
}

func TestSyncCreateCard(t *testing.T) {
	appID := os.Getenv("SYNC_APP_ID")
	token := os.Getenv("SYNC_TOKEN")
	appPrivateKey := os.Getenv("SYNC_APP_KEY")
	appPrivateKeyPass := os.Getenv("SYNC_APP_KEY_PASS")

	priv, err := virgil.Crypto().ImportPrivateKey([]byte(appPrivateKey), appPrivateKeyPass)
	assert.NoError(t, err, "Cannot import private key")

	deviceKeypair, err := virgil.Crypto().GenerateKeypair()
	assert.NoError(t, err, "Cannot generate key pair")

	uid := uuid.Must(uuid.NewV4()).String()
	req, err := virgil.NewCreateCardRequest(uid, "temp", deviceKeypair.PublicKey(), virgil.CardParams{
		Scope: virgil.CardScope.Application,
		Data: map[string]string{
			"os": "macOS",
		},
		DeviceInfo: virgil.DeviceInfo{
			Device:     "iphone7",
			DeviceName: "my iphone",
		},
	})
	assert.NoError(t, err, "Cannot create card request")

	signer := virgil.RequestSigner{}
	err = signer.SelfSign(req, deviceKeypair.PrivateKey())
	assert.NoError(t, err, "Cannot self sign")

	lc, err := virgil.NewClient("", virgil.ClientTransport(virgilhttp.NewTransportClient(host, host, host, host)))
	assert.NoError(t, err, "Cannot create virgil client")

	vc, err := virgil.NewClient(token)
	assert.NoError(t, err, "Cannot create virgil client")

	card, err := lc.CreateCard(req)
	assert.NoError(t, err, "Cannot create create card in the virgild")

	time.Sleep(5 * time.Second)

	defer func() {
		req, _ := virgil.NewRevokeCardRequest(card.ID, virgil.RevocationReason.Unspecified)
		signer.AuthoritySign(req, appID, priv)
		vc.RevokeCard(req)
	}()

	c, err := vc.GetCard(card.ID)
	assert.NoError(t, err, "Cannot get card by id(%v)", card.ID)

	assert.EqualValues(t, card, c)
}

func TestSyncRevokeCard(t *testing.T) {
	appID := os.Getenv("SYNC_APP_ID")
	token := os.Getenv("SYNC_TOKEN")
	appPrivateKey := os.Getenv("SYNC_APP_KEY")
	appPrivateKeyPass := os.Getenv("SYNC_APP_KEY_PASS")

	priv, err := virgil.Crypto().ImportPrivateKey([]byte(appPrivateKey), appPrivateKeyPass)
	assert.NoError(t, err, "Cannot import private key")

	deviceKeypair, err := virgil.Crypto().GenerateKeypair()
	assert.NoError(t, err, "Cannot generate key pair")

	uid := uuid.Must(uuid.NewV4()).String()
	req, err := virgil.NewCreateCardRequest(uid, "temp", deviceKeypair.PublicKey(), virgil.CardParams{
		Scope: virgil.CardScope.Application,
		Data: map[string]string{
			"os": "macOS",
		},
		DeviceInfo: virgil.DeviceInfo{
			Device:     "iphone7",
			DeviceName: "my iphone",
		},
	})
	assert.NoError(t, err, "Cannot create card request")

	signer := virgil.RequestSigner{}
	err = signer.SelfSign(req, deviceKeypair.PrivateKey())
	assert.NoError(t, err, "Cannot self sign")
	err = signer.AuthoritySign(req, appID, priv)
	assert.NoError(t, err, "Cannot add app sing")

	lc, err := virgil.NewClient("", virgil.ClientTransport(virgilhttp.NewTransportClient(host, host, host, host)))
	assert.NoError(t, err, "Cannot create virgil client")

	vc, err := virgil.NewClient(token)
	assert.NoError(t, err, "Cannot create virgil client")

	card, err := vc.CreateCard(req)
	assert.NoError(t, err, "Cannot create create card in the cloud")

	time.Sleep(5 * time.Second)

	req, _ = virgil.NewRevokeCardRequest(card.ID, virgil.RevocationReason.Unspecified)
	err = lc.RevokeCard(req)

	assert.NoError(t, err, "Cannot revoke card")

	_, err = vc.GetCard(card.ID)
	assert.Error(t, err, "We expected that card not found")
}

type vd struct {
	CardID    string `json:"card_id"`
	PublicKey string `json:"public_key"`
}

func TestLocalGetCard(t *testing.T) {
	resp, err := http.Get(host + "/api/card")
	if err != nil {
		t.Fatalf("Cannot get public card info: %+v", err)
	}
	defer resp.Body.Close()

	var vd config.VirgilDCard
	err = json.NewDecoder(resp.Body).Decode(&vd)
	assert.NoError(t, err, "Cannot decode VirgilD Card Info")

	pub, err := virgil.Crypto().ImportPublicKey([]byte(vd.PublicKey))
	assert.NoError(t, err, "Cannot import public key")

	cv := virgil.NewCardsValidator()
	cv.AddVerifier(vd.CardID, pub)
	client, err := virgil.NewClient("", virgil.ClientTransport(virgilhttp.NewTransportClient(host, host, host, host)),
		virgil.ClientCardsValidator(cv))
	assert.NoError(t, err, "Cannot create virgil client")

	_, err = client.GetCard(vd.CardID)
	assert.NoError(t, err)
}

func TestLocalSearchCards(t *testing.T) {
	resp, err := http.Get(host + "/api/card")
	if err != nil {
		t.Fatalf("Cannot get public card info: %+v", err)
	}
	defer resp.Body.Close()

	var vd config.VirgilDCard
	err = json.NewDecoder(resp.Body).Decode(&vd)
	assert.NoError(t, err, "Cannot decode VirgilD Card Info")

	pub, err := virgil.Crypto().ImportPublicKey([]byte(vd.PublicKey))
	assert.NoError(t, err, "Cannot import public key")

	cv := virgil.NewCardsValidator()
	cv.AddVerifier(vd.CardID, pub)
	client, err := virgil.NewClient("", virgil.ClientTransport(virgilhttp.NewTransportClient(host, host, host, host)),
		virgil.ClientCardsValidator(cv))
	assert.NoError(t, err, "Cannot create virgil client")

	deviceKeypair, err := virgil.Crypto().GenerateKeypair()
	assert.NoError(t, err, "Cannot generate key pair")

	uid := uuid.Must(uuid.NewV4()).String()
	req, err := virgil.NewCreateCardRequest(uid, "temp", deviceKeypair.PublicKey(), virgil.CardParams{
		Scope: virgil.CardScope.Application,
		Data: map[string]string{
			"os": "macOS",
		},
		DeviceInfo: virgil.DeviceInfo{
			Device:     "iphone7",
			DeviceName: "my iphone",
		},
	})
	assert.NoError(t, err, "Cannot create card request")

	signer := virgil.RequestSigner{}
	err = signer.SelfSign(req, deviceKeypair.PrivateKey())
	assert.NoError(t, err, "Cannot self sign")

	card, err := client.CreateCard(req)
	assert.NoError(t, err, "Cannot create create card in the cloud")

	defer func() {
		req, _ := virgil.NewRevokeCardRequest(card.ID, virgil.RevocationReason.Unspecified)
		client.RevokeCard(req)
	}()

	cs, err := client.SearchCards(virgil.SearchCriteriaByIdentities(uid))
	assert.NoError(t, err, "Cannot search cards by temp name (%v)", uid)
	assert.Len(t, cs, 1, "Number of found cards are not equal")
	c := cs[0]

	assert.EqualValues(t, card, c)
}

func TestLocalCreateCard(t *testing.T) {
	resp, err := http.Get(host + "/api/card")
	if err != nil {
		t.Fatalf("Cannot get public card info: %+v", err)
	}
	defer resp.Body.Close()

	var vd config.VirgilDCard
	err = json.NewDecoder(resp.Body).Decode(&vd)
	assert.NoError(t, err, "Cannot decode VirgilD Card Info")

	pub, err := virgil.Crypto().ImportPublicKey([]byte(vd.PublicKey))
	assert.NoError(t, err, "Cannot import public key")

	cv := virgil.NewCardsValidator()
	cv.AddVerifier(vd.CardID, pub)
	client, err := virgil.NewClient("", virgil.ClientTransport(virgilhttp.NewTransportClient(host, host, host, host)),
		virgil.ClientCardsValidator(cv))
	assert.NoError(t, err, "Cannot create virgil client")

	deviceKeypair, err := virgil.Crypto().GenerateKeypair()
	assert.NoError(t, err, "Cannot generate key pair")

	uid := uuid.Must(uuid.NewV4()).String()
	req, err := virgil.NewCreateCardRequest(uid, "temp", deviceKeypair.PublicKey(), virgil.CardParams{
		Scope: virgil.CardScope.Application,
		Data: map[string]string{
			"os": "macOS",
		},
		DeviceInfo: virgil.DeviceInfo{
			Device:     "iphone7",
			DeviceName: "my iphone",
		},
	})
	assert.NoError(t, err, "Cannot create card request")

	signer := virgil.RequestSigner{}
	err = signer.SelfSign(req, deviceKeypair.PrivateKey())
	assert.NoError(t, err, "Cannot self sign")

	card, err := client.CreateCard(req)
	assert.NoError(t, err, "Cannot create create card in the virgild")

	defer func() {
		req, _ := virgil.NewRevokeCardRequest(card.ID, virgil.RevocationReason.Unspecified)
		client.RevokeCard(req)
	}()

	c, err := client.GetCard(card.ID)
	assert.NoError(t, err, "Cannot get card by id(%v)", card.ID)

	assert.EqualValues(t, card, c)
}

func TestLocalRevokeCard(t *testing.T) {
	resp, err := http.Get(host + "/api/card")
	if err != nil {
		t.Fatalf("Cannot get public card info: %+v", err)
	}
	defer resp.Body.Close()

	var vd config.VirgilDCard
	err = json.NewDecoder(resp.Body).Decode(&vd)
	assert.NoError(t, err, "Cannot decode VirgilD Card Info")

	pub, err := virgil.Crypto().ImportPublicKey([]byte(vd.PublicKey))
	assert.NoError(t, err, "Cannot import public key")

	cv := virgil.NewCardsValidator()
	cv.AddVerifier(vd.CardID, pub)
	client, err := virgil.NewClient("", virgil.ClientTransport(virgilhttp.NewTransportClient(host, host, host, host)),
		virgil.ClientCardsValidator(cv))
	assert.NoError(t, err, "Cannot create virgil client")

	deviceKeypair, err := virgil.Crypto().GenerateKeypair()
	assert.NoError(t, err, "Cannot generate key pair")

	uid := uuid.Must(uuid.NewV4()).String()
	req, err := virgil.NewCreateCardRequest(uid, "temp", deviceKeypair.PublicKey(), virgil.CardParams{
		Scope: virgil.CardScope.Application,
		Data: map[string]string{
			"os": "macOS",
		},
		DeviceInfo: virgil.DeviceInfo{
			Device:     "iphone7",
			DeviceName: "my iphone",
		},
	})
	assert.NoError(t, err, "Cannot create card request")

	signer := virgil.RequestSigner{}
	err = signer.SelfSign(req, deviceKeypair.PrivateKey())
	assert.NoError(t, err, "Cannot self sign")

	card, err := client.CreateCard(req)
	assert.NoError(t, err, "Cannot create card in the local storage")

	req, _ = virgil.NewRevokeCardRequest(card.ID, virgil.RevocationReason.Unspecified)
	err = client.RevokeCard(req)

	assert.NoError(t, err, "Cannot revoke card")

	_, err = client.GetCard(card.ID)
	assert.Error(t, err, "We expected that card not found")
}
