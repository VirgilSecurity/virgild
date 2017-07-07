package card

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	virgil "gopkg.in/virgil.v4"

	"github.com/VirgilSecurity/virgild/coreapi"
	"github.com/VirgilSecurity/virgild/modules/card/core"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

var cloudDurationMetrics = prometheus.NewSummaryVec(prometheus.SummaryOpts{
	Name:      "cards_service",
	Subsystem: "cards",
	Namespace: "virgild",
}, []string{"type"})

type client interface {
	Do(req *http.Request) (*http.Response, error)
}

type cloudCard struct {
	RAService    string
	CardsService string
	Client       client
}

func (c *cloudCard) getCard(ctx context.Context, id string) (*virgil.CardResponse, error) {
	var body []byte
	var err error

	timer := prometheus.NewTimer(cloudDurationMetrics.WithLabelValues("get_card"))
	body, err = c.send(ctx, http.MethodGet, c.CardsService+"/v4/card/"+id, nil)
	timer.ObserveDuration()

	if err != nil {
		return nil, err
	}
	card := new(virgil.CardResponse)
	err = json.Unmarshal(body, card)
	return card, err
}

func (c *cloudCard) searchCards(ctx context.Context, crit *virgil.Criteria) ([]virgil.CardResponse, error) {
	var body []byte
	var err error

	timer := prometheus.NewTimer(cloudDurationMetrics.WithLabelValues("search"))
	body, err = c.send(ctx, http.MethodPost, c.CardsService+"/v4/card/actions/search", crit)
	timer.ObserveDuration()

	if err != nil {
		return nil, err
	}
	var cards []virgil.CardResponse
	err = json.Unmarshal(body, &cards)
	return cards, err
}

func (c *cloudCard) createCard(ctx context.Context, req *core.CreateCardRequest) (*virgil.CardResponse, error) {
	var body []byte
	var err error

	timer := prometheus.NewTimer(cloudDurationMetrics.WithLabelValues("create_card"))
	body, err = c.send(ctx, http.MethodPost, c.RAService+"/v1/card", req.Request)
	timer.ObserveDuration()

	if err != nil {
		return nil, err
	}
	card := new(virgil.CardResponse)
	err = json.Unmarshal(body, card)
	return card, err
}

func (c *cloudCard) revokeCard(ctx context.Context, req *core.RevokeCardRequest) error {
	var err error

	timer := prometheus.NewTimer(cloudDurationMetrics.WithLabelValues("revoke_card"))
	_, err = c.send(ctx, http.MethodDelete, c.RAService+"/v1/card/"+req.Info.ID, req.Request)
	timer.ObserveDuration()

	return err
}

func (c *cloudCard) createRelation(ctx context.Context, req *core.CreateRelationRequest) (*virgil.CardResponse, error) {
	var body []byte
	var err error

	timer := prometheus.NewTimer(cloudDurationMetrics.WithLabelValues("create_relation"))
	body, err = c.send(ctx, http.MethodPost, c.CardsService+"/v4/card/"+req.ID+"/collections/relations", req.Request)
	timer.ObserveDuration()

	if err != nil {
		return nil, err
	}
	card := new(virgil.CardResponse)
	err = json.Unmarshal(body, card)
	return card, err
}

func (c *cloudCard) revokeRelation(ctx context.Context, req *core.RevokeRelationRequest) (*virgil.CardResponse, error) {
	var body []byte
	var err error

	timer := prometheus.NewTimer(cloudDurationMetrics.WithLabelValues("revoke_relation"))
	body, err = c.send(ctx, http.MethodDelete, c.CardsService+"/v4/card/"+req.ID+"/collections/relations", req.Request)
	timer.ObserveDuration()

	if err != nil {
		return nil, err
	}
	card := new(virgil.CardResponse)
	err = json.Unmarshal(body, card)
	return card, err
}

func (c *cloudCard) send(ctx context.Context, method string, urlStr string, payload interface{}) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		bp, err := json.Marshal(payload)
		if err != nil {
			return nil, errors.Wrapf(err, "Cloud.Send(cannot marshal payload [payload: %v])", payload)
		}
		body = bytes.NewReader(bp)
	}
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, errors.Wrap(err, "Cloud.Send(cannot create request)")
	}
	auth := core.GetAuthHeader(ctx)
	req.Header.Set("Authorization", auth)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Cloud.Send(default client send req)")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Cloud.Send(read reasponse)")
	}
	if resp.StatusCode == http.StatusOK {
		return respBody, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, coreapi.EntityNotFoundErr
	}

	verr := new(virgilError)
	err = json.Unmarshal(respBody, verr)
	if err != nil {
		return nil, errors.Wrapf(err, "Cloud.Send(unmarshal error [body: %s])", respBody)
	}

	return nil, coreapi.APIError{
		Code:       verr.Code,
		StatusCode: resp.StatusCode,
	}
}

type virgilError struct {
	Code int
}
