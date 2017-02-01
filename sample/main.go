package main

import (
	"fmt"

	virgil "gopkg.in/virgil.v4"

	"gopkg.in/virgil.v4/errors"
	"gopkg.in/virgil.v4/transport/virgilhttp"
)

var (
	cardsServicePublicKey = []byte(`MCowBQYDK2VwAyEA6Pij81JVf2ewrvaHHd9MUjq38yrPUZ9aSeuuAHKCOIo=`)
	cardServiceID         = "e26f8b6a4a5919c5db1710710421298e956f9bf836feaa765a6f04d93f8595ae"
)

func main() {
	crypto := virgil.Crypto()

	//set up custom validator with stg key  & our app key
	customValidator := virgil.NewCardsValidator()

	// cardsServicePublic, _ := crypto.ImportPublicKey(cardsServicePublicKey)

	// customValidator.AddVerifier(cardServiceID, cardsServicePublic)

	client, err := virgil.NewClient("123",
		virgil.ClientTransport(
			virgilhttp.NewTransportClient(
				"http://localhost:8080",
				"http://localhost:8080",
				"http://localhost:8080",
				"http://localhost:8080")),
		virgil.ClientCardsValidator(customValidator))

	//generate device key & create card
	deviceKeypair, _ := crypto.GenerateKeypair()

	req, err := virgil.NewCreateCardRequest("Device #1", "Smart Iot Device", deviceKeypair.PublicKey(), virgil.CardParams{
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
	signer.SelfSign(req, deviceKeypair.PrivateKey())

	//.....
	// signer.AuthoritySign(req, appCardID, appPrivateKey)

	card, err := client.CreateCard(req)
	if err != nil {

		e, ok := errors.ToSdkError(err)
		if ok {
			fmt.Println(e.ServiceErrorCode())
		}
		fmt.Printf("Error:%+v\n", err)
		panic(err)
	}
	fmt.Println(card.ID)
	fmt.Println(card.Identity)
	fmt.Println(card.CreatedAt)

	// appCard, err := client.GetCard(cardServiceID)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("App card:", appCard.Identity)

	cards, err := client.SearchCards(virgil.SearchCriteriaByAppBundle("com.gibsonmic.ed255app"))
	if err != nil {
		fmt.Println(err)
	}
	if len(cards) != 0 {
		fmt.Println("Find global card")
		fmt.Println("appCard:", cards[0].ID)
	}

	cards, err = client.SearchCards(virgil.Criteria{
		IdentityType: "Smart Iot Device",
		Identities: []string{
			"Device #1",
		},
	})
	if err != nil {
		panic(err)
	}

	for _, c := range cards {
		fmt.Println(c.Identity)
		gotCard, err := client.GetCard(c.ID)
		if err != nil {
			panic(err)
		}
		revreq, err := virgil.NewRevokeCardRequest(gotCard.ID, virgil.RevocationReason.Unspecified)
		if err != nil {
			panic(err)
		}
		// signer.AuthoritySign(revreq, appCard.ID, appPrivateKey)
		err = client.RevokeCard(revreq)
		if err != nil {
			panic(err)
		}
		fmt.Println(gotCard.ID)
	}
}
