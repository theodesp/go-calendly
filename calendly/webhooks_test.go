package calendly

import (
	"github.com/stretchr/testify/assert"
	"fmt"
	"net/http"
	"context"
	"io/ioutil"
	"net/url"
)

func (suite *CalendlyClientTestSuite) TestWebhooksService_Create() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf("/%s", webhooksPath)

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodPost)

		body, _ := ioutil.ReadAll(r.Body)
		expected := `%22url%253Dhttp%253A%252F%252Fwebhook%2526events%255B%255D%253Dinvitee.cancelled%22%0A`
		assert.Equal(expected, url.QueryEscape(string(body)))

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":123}`)
	})
	opts := &WebhooksOpts{
		Url: "http://webhook",
		Events: []EventHookType{InviteeCancelledHookType},
	}
	v, resp, err := suite.client.Webhooks.Create(context.Background(), opts)
	want := &Webhook{ID: int64(123)}

	assert.Nil(err)
	assert.Equal(want, v)
	assert.Equal(http.StatusCreated, resp.StatusCode)
}


func (suite *CalendlyClientTestSuite) TestWebhooksService_GetByID() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf("/%s", fmt.Sprintf(getWebhookpath, "1"))

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodGet)
		fmt.Fprint(w, `{"data":{"id":123}}`)
	})

	webHook, _, err := suite.client.Webhooks.GetByID(context.Background(), int64(1))
	assert.Nil(err)

	want := &Webhook{
		ID: int64(123),
	}
	assert.Equal(want, webHook)
}

func (suite *CalendlyClientTestSuite) TestEventTypesService_ListWebhooks() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf("/%s", webhooksPath)

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodGet)
		fmt.Fprint(w, `{"data":[{"id":123},{"id":456}]}`)
	})

	webhooks, _, err := suite.client.Webhooks.List(context.Background())
	assert.Nil(err)

	want := []*Webhook{
		{ID: int64(123)},
		{ID: int64(456)},
	}
	assert.Equal(want, webhooks)
}

func (suite *CalendlyClientTestSuite) TestEventTypesService_DeleteWebhook() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf("/%s", fmt.Sprintf(getWebhookpath, "1"))

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodDelete)
		w.WriteHeader(http.StatusOK)
	})

	resp, err := suite.client.Webhooks.Delete(context.Background(), int64(1))
	assert.Nil(err)

	assert.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *CalendlyClientTestSuite) TestWebhooksService_InvalidParams() {
	assert := assert.New(suite.T())

	_, _, err := suite.client.Webhooks.Create(context.Background(), nil)
	assert.NotNil(err)

	_, _, err = suite.client.Webhooks.Create(context.Background(), &WebhooksOpts{})
	assert.NotNil(err)

	_, _, err = suite.client.Webhooks.Create(context.Background(), &WebhooksOpts{
		Url: "http://192.168.0.%31/",
	})
	assert.NotNil(err)
}
