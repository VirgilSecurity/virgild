package card

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"gopkg.in/virgil.v4"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fakeHttpClient struct {
	mock.Mock
}

func (f *fakeHttpClient) Do(req *http.Request) (resp *http.Response, err error) {
	args := f.Called(req)
	resp, _ = args.Get(0).(*http.Response)
	err = args.Error(1)
	return
}

func TestCloudTable_ClientReturnErr_ReturnErr(t *testing.T) {
	table := map[string]func(cloud *cloudCard) error{
		"Get card": func(cloud *cloudCard) (err error) {
			_, err = cloud.getCard(context.Background(), "1234")
			return
		},
		"Search cards": func(cloud *cloudCard) (err error) {
			_, err = cloud.searchCards(context.Background(), &virgil.Criteria{})
			return
		},
		"Create card": func(cloud *cloudCard) (err error) {
			_, err = cloud.createCard(context.Background(), &core.CreateCardRequest{Request: virgil.SignableRequest{Snapshot: []byte(`Snapshot`)}})
			return
		},
		"Revoke card": func(cloud *cloudCard) (err error) {
			return cloud.revokeCard(context.Background(), &core.RevokeCardRequest{Request: virgil.SignableRequest{Snapshot: []byte(`Snapshot`)}})
		},
		"Create relation": func(cloud *cloudCard) (err error) {
			_, err = cloud.createRelation(context.Background(), &core.CreateRelationRequest{Request: virgil.SignableRequest{Snapshot: []byte(`Snapshot`)}})
			return
		},
		"Revoke relation": func(cloud *cloudCard) (err error) {
			_, err = cloud.revokeRelation(context.Background(), &core.RevokeRelationRequest{Request: virgil.SignableRequest{Snapshot: []byte(`Snapshot`)}})
			return
		},
	}

	f := new(fakeHttpClient)
	f.On("Do", mock.Anything).Return(nil, fmt.Errorf("ERROR"))
	cloud := cloudCard{RAService: "ra-service", CardsService: "cards-service", Client: f}
	for name, function := range table {
		err := function(&cloud)
		assert.Error(t, err, "ERROR", name)
	}
}

func TestCloudTable_ClientReturnOtherErr_ReturnVirgilErr(t *testing.T) {
	table := map[string]func(cloud *cloudCard) error{
		"Get card": func(cloud *cloudCard) (err error) {
			_, err = cloud.getCard(context.Background(), "1234")
			return
		},
		"Search cards": func(cloud *cloudCard) (err error) {
			_, err = cloud.searchCards(context.Background(), &virgil.Criteria{})
			return
		},
		"Create card": func(cloud *cloudCard) (err error) {
			_, err = cloud.createCard(context.Background(), &core.CreateCardRequest{Request: virgil.SignableRequest{Snapshot: []byte(`Snapshot`)}})
			return
		},
		"Revoke card": func(cloud *cloudCard) (err error) {
			return cloud.revokeCard(context.Background(), &core.RevokeCardRequest{Request: virgil.SignableRequest{Snapshot: []byte(`Snapshot`)}})
		},
		"Create relation": func(cloud *cloudCard) (err error) {
			_, err = cloud.createRelation(context.Background(), &core.CreateRelationRequest{Request: virgil.SignableRequest{Snapshot: []byte(`Snapshot`)}})
			return
		},
		"Revoke relation": func(cloud *cloudCard) (err error) {
			_, err = cloud.revokeRelation(context.Background(), &core.RevokeRelationRequest{Request: virgil.SignableRequest{Snapshot: []byte(`Snapshot`)}})
			return
		},
	}

	f := new(fakeHttpClient)
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(`{"code":1234}`)),
	}
	f.On("Do", mock.Anything).Return(resp, nil)
	cloud := cloudCard{RAService: "ra-service", CardsService: "cards-service", Client: f}

	for name, function := range table {
		err := function(&cloud)
		assert.Error(t, err, coreapi.APIError{Code: 1234, StatusCode: http.StatusBadRequest}, name)
	}
}

func TestCloudTable_ClientReturnOkBodyInvalid_ReturnErr(t *testing.T) {
	table := map[string]func(cloud *cloudCard) error{
		"Get card": func(cloud *cloudCard) (err error) {
			_, err = cloud.getCard(context.Background(), "1234")
			return
		},
		"Search cards": func(cloud *cloudCard) (err error) {
			_, err = cloud.searchCards(context.Background(), &virgil.Criteria{})
			return
		},
		"Create card": func(cloud *cloudCard) (err error) {
			_, err = cloud.createCard(context.Background(), &core.CreateCardRequest{Request: virgil.SignableRequest{Snapshot: []byte(`Snapshot`)}})
			return
		},
		"Revoke card": func(cloud *cloudCard) (err error) {
			return cloud.revokeCard(context.Background(), &core.RevokeCardRequest{Request: virgil.SignableRequest{Snapshot: []byte(`Snapshot`)}})
		},
		"Create relation": func(cloud *cloudCard) (err error) {
			_, err = cloud.createRelation(context.Background(), &core.CreateRelationRequest{Request: virgil.SignableRequest{Snapshot: []byte(`Snapshot`)}})
			return
		},
		"Revoke relation": func(cloud *cloudCard) (err error) {
			_, err = cloud.revokeRelation(context.Background(), &core.RevokeRelationRequest{Request: virgil.SignableRequest{Snapshot: []byte(`Snapshot`)}})
			return
		},
	}

	f := new(fakeHttpClient)
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(`asdf: fasd`)),
	}
	f.On("Do", mock.Anything).Return(resp, nil)
	cloud := cloudCard{RAService: "ra-service", CardsService: "cards-service", Client: f}

	for name, function := range table {
		err := function(&cloud)
		assert.NotNil(t, err, name)
	}
}

func TestCloudGetCard_ClientReturn404_ReturnVirgilErr(t *testing.T) {
	f := new(fakeHttpClient)
	resp := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	f.On("Do", mock.Anything).Return(resp, nil)
	cloud := cloudCard{RAService: "ra-service", CardsService: "cards-service", Client: f}
	_, err := cloud.getCard(context.Background(), "1234")
	assert.Error(t, err, coreapi.EntityNotFoundErr)
}

func TestCloudGetCard_ClientReturnOk_ReturnVal(t *testing.T) {
	const authHeader = "header"
	expectedReq, _ := http.NewRequest(http.MethodGet, "cards-service/v4/card/1234", nil)
	expectedReq.Header.Set("Authorization", authHeader)

	expectedCard := &virgil.CardResponse{
		Snapshot: []byte(`snapshot`),
		ID:       "1234",
	}

	respBody, _ := json.Marshal(expectedCard)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}
	f := new(fakeHttpClient)
	f.On("Do", expectedReq).Return(resp, nil)
	cloud := cloudCard{RAService: "ra-service", CardsService: "cards-service", Client: f}

	ctx := core.SetAuthHeader(context.Background(), authHeader)
	card, err := cloud.getCard(ctx, "1234")

	assert.NoError(t, err)
	assert.Equal(t, expectedCard, card)
}

func TestCloudSearchCards_ClientReturnOk_ReturnVal(t *testing.T) {
	const authHeader = "header"
	crit := &virgil.Criteria{
		Identities: []string{"alice", "bob"},
	}
	bCrit, _ := json.Marshal(crit)
	matchReq := func(req *http.Request) bool {
		body, _ := ioutil.ReadAll(req.Body)
		return req.Method == http.MethodPost &&
			req.URL.String() == "cards-service/v4/card/actions/search" &&
			req.Header.Get("Authorization") == authHeader &&
			bytes.Equal(body, bCrit)
	}

	expectedCard := []virgil.CardResponse{
		virgil.CardResponse{
			Snapshot: []byte(`snapshot`),
			ID:       "1234",
		},
	}
	respBody, _ := json.Marshal(expectedCard)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}

	f := new(fakeHttpClient)
	f.On("Do", mock.MatchedBy(matchReq)).Return(resp, nil)
	cloud := cloudCard{RAService: "ra-service", CardsService: "cards-service", Client: f}

	ctx := core.SetAuthHeader(context.Background(), authHeader)
	card, err := cloud.searchCards(ctx, crit)

	assert.NoError(t, err)
	assert.Equal(t, expectedCard, card)
}

func TestCloudCreateCard_ClientReturnOk_ReturnVal(t *testing.T) {
	const authHeader = "header"
	signableRequest := virgil.SignableRequest{
		Snapshot: []byte(`Snapshot`),
		Meta: virgil.RequestMeta{
			Signatures: map[string][]byte{
				"1234": []byte(`sign`),
			},
		},
	}
	bodyReqest, _ := json.Marshal(signableRequest)
	createCard := &core.CreateCardRequest{Request: signableRequest}

	matchReq := func(req *http.Request) bool {
		body, _ := ioutil.ReadAll(req.Body)
		return req.Method == http.MethodPost &&
			req.URL.String() == "ra-service/v1/card" &&
			req.Header.Get("Authorization") == authHeader &&
			bytes.Equal(body, bodyReqest)
	}

	expectedCard := &virgil.CardResponse{
		Snapshot: []byte(`snapshot`),
		ID:       "1234",
	}

	respBody, _ := json.Marshal(expectedCard)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}
	f := new(fakeHttpClient)
	f.On("Do", mock.MatchedBy(matchReq)).Return(resp, nil)
	cloud := cloudCard{RAService: "ra-service", CardsService: "cards-service", Client: f}

	ctx := core.SetAuthHeader(context.Background(), authHeader)
	card, err := cloud.createCard(ctx, createCard)

	assert.NoError(t, err)
	assert.Equal(t, expectedCard, card)
}

func TestCloudRevokeCard_ClientReturnOk_ReturnVal(t *testing.T) {
	const authHeader = "header"
	signableRequest := virgil.SignableRequest{
		Snapshot: []byte(`Snapshot`),
		Meta: virgil.RequestMeta{
			Signatures: map[string][]byte{
				"1234": []byte(`sign`),
			},
		},
	}
	bodyReqest, _ := json.Marshal(signableRequest)
	revokeCard := &core.RevokeCardRequest{Info: virgil.RevokeCardRequest{ID: "1234"}, Request: signableRequest}

	matchReq := func(req *http.Request) bool {
		body, _ := ioutil.ReadAll(req.Body)
		return req.Method == http.MethodDelete &&
			req.URL.String() == "ra-service/v1/card/1234" &&
			req.Header.Get("Authorization") == authHeader &&
			bytes.Equal(body, bodyReqest)
	}

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(nil)),
	}
	f := new(fakeHttpClient)
	f.On("Do", mock.MatchedBy(matchReq)).Return(resp, nil)
	cloud := cloudCard{RAService: "ra-service", CardsService: "cards-service", Client: f}

	ctx := core.SetAuthHeader(context.Background(), authHeader)
	err := cloud.revokeCard(ctx, revokeCard)

	assert.NoError(t, err)
}

func TestCloudCreateRelations_ClientReturnOk_ReturnVal(t *testing.T) {
	const authHeader = "header"
	signableRequest := virgil.SignableRequest{
		Snapshot: []byte(`Snapshot`),
		Meta: virgil.RequestMeta{
			Signatures: map[string][]byte{
				"1234": []byte(`sign`),
			},
		},
	}
	bodyReqest, _ := json.Marshal(signableRequest)
	createRelation := &core.CreateRelationRequest{ID: "1234", Request: signableRequest}

	matchReq := func(req *http.Request) bool {
		body, _ := ioutil.ReadAll(req.Body)
		return req.Method == http.MethodPost &&
			req.URL.String() == "cards-service/v4/card/1234/collections/relations" &&
			req.Header.Get("Authorization") == authHeader &&
			bytes.Equal(body, bodyReqest)
	}

	expectedCard := &virgil.CardResponse{
		Snapshot: []byte(`snapshot`),
		ID:       "1234",
	}

	respBody, _ := json.Marshal(expectedCard)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}
	f := new(fakeHttpClient)
	f.On("Do", mock.MatchedBy(matchReq)).Return(resp, nil)
	cloud := cloudCard{RAService: "ra-service", CardsService: "cards-service", Client: f}

	ctx := core.SetAuthHeader(context.Background(), authHeader)
	card, err := cloud.createRelation(ctx, createRelation)

	assert.NoError(t, err)
	assert.Equal(t, expectedCard, card)
}

func TestCloudRevokeRelations_ClientReturnOk_ReturnVal(t *testing.T) {
	const authHeader = "header"
	signableRequest := virgil.SignableRequest{
		Snapshot: []byte(`Snapshot`),
		Meta: virgil.RequestMeta{
			Signatures: map[string][]byte{
				"1234": []byte(`sign`),
			},
		},
	}
	bodyReqest, _ := json.Marshal(signableRequest)
	revokeRelation := &core.RevokeRelationRequest{ID: "1234", Request: signableRequest}

	matchReq := func(req *http.Request) bool {
		body, _ := ioutil.ReadAll(req.Body)
		return req.Method == http.MethodDelete &&
			req.URL.String() == "cards-service/v4/card/1234/collections/relations" &&
			req.Header.Get("Authorization") == authHeader &&
			bytes.Equal(body, bodyReqest)
	}

	expectedCard := &virgil.CardResponse{
		Snapshot: []byte(`snapshot`),
		ID:       "1234",
	}

	respBody, _ := json.Marshal(expectedCard)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}
	f := new(fakeHttpClient)
	f.On("Do", mock.MatchedBy(matchReq)).Return(resp, nil)
	cloud := cloudCard{RAService: "ra-service", CardsService: "cards-service", Client: f}

	ctx := core.SetAuthHeader(context.Background(), authHeader)
	card, err := cloud.revokeRelation(ctx, revokeRelation)

	assert.NoError(t, err)
	assert.Equal(t, expectedCard, card)
}
