package main

import (
	"fmt"

	virgil "gopkg.in/virgil.v4"

	"gopkg.in/virgil.v4/errors"
	"gopkg.in/virgil.v4/transport/virgilhttp"
)

var (
	cardsServicePublicKey = []byte(`MCowBQYDK2VwAyEAI26M2oj8M+6r20kPE5JhgbvoXGT2IZr73klehP9W9mg=`)
	cardServiceID         = "e26f8b6a4a5919c5db1710710421298e956f9bf836feaa765a6f04d93f8595ae"
	crypto                = virgil.Crypto()
	client                *virgil.Client
)

func main() {
	//set up custom validator with stg key  & our app key
	client = setupClient()
	id := createCard("Device #1", "Smart Iot Device")
	fmt.Println()

	getCard(id)
	fmt.Println()

	cards := searchCards("Device #1", "Smart Iot Device")
	fmt.Println("Found cards by identity: Device #1 and IdentityType: Smart Iot Device")
	for _, c := range cards {
		fmt.Println(PrintCard(c))
	}
	fmt.Println()

	fmt.Println("Deleted cards:")
	for _, c := range cards {
		revreq, err := virgil.NewRevokeCardRequest(c.ID, virgil.RevocationReason.Unspecified)
		if err != nil {
			panic(err)
		}
		err = client.RevokeCard(revreq)
		if err != nil {
			panic(err)
		}
		fmt.Println(PrintCard(c))
	}
}

func setupClient() *virgil.Client {
	customValidator := virgil.NewCardsValidator()

	cardsServicePublic, _ := crypto.ImportPublicKey(cardsServicePublicKey)

	customValidator.AddVerifier(cardServiceID, cardsServicePublic)

	client, _ := virgil.NewClient("",
		virgil.ClientTransport(
			virgilhttp.NewTransportClient(
				"http://localhost:8080",
				"http://localhost:8080",
				"http://localhost:8080",
				"http://localhost:8080")),
		virgil.ClientCardsValidator(customValidator))
	return client
}

func createCard(i, it string) string {
	deviceKeypair, _ := crypto.GenerateKeypair()

	req, err := virgil.NewCreateCardRequest(i, it, deviceKeypair.PublicKey(), virgil.CardParams{
		Scope: virgil.CardScope.Application,
		Data: map[string]string{
			"os": "macOS",
		},
		DeviceInfo: virgil.DeviceInfo{
			Device:     "iphone7",
			DeviceName: "my iphone",
		},
	})
	signer := virgil.RequestSigner{}
	err = signer.SelfSign(req, deviceKeypair.PrivateKey())
	if err != nil {
		panic(err)
	}

	card, err := client.CreateCard(req)
	if err != nil {
		e, ok := errors.ToSdkError(err)
		if ok {
			fmt.Println("Service error code:", e.ServiceErrorCode())
		}
		fmt.Printf("Error: %+v\n", err)
		panic(err)
	}

	fmt.Println("Created card {", PrintCard(card), "}")
	return card.ID
}

func getCard(id string) {
	appCard, err := client.GetCard(id)
	if err != nil {
		panic(err)
	}
	fmt.Println("Found card by id (", id, "){", PrintCard(appCard), "}")
}

func searchCards(i, it string) []*virgil.Card {
	cards, err := client.SearchCards(&virgil.Criteria{
		IdentityType: it,
		Identities: []string{
			i,
		},
	})
	if err != nil {
		panic(err)
	}
	return cards
}

func PrintCard(c *virgil.Card) string {
	return fmt.Sprintf("ID: %v Identity: %v IdentityType: %v Scope: %v Data: %v", c.ID, c.Identity, c.IdentityType, c.Scope, c.Data)
}
