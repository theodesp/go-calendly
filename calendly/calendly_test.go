package calendly

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type CalendlyClientTestSuite struct {
	suite.Suite

	// client is the Calendly client being tested.
	client *Client

	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
}

func TestCalendlyClientTestSuite(t *testing.T) {
	suite.Run(t, new(CalendlyClientTestSuite))
}

func (suite *CalendlyClientTestSuite) SetupTest() {
	suite.client = NewClient(nil)
	suite.mux = http.NewServeMux()
	suite.server = httptest.NewServer(suite.mux)

	// calendly client configured to use test server
	url, _ := url.Parse(suite.server.URL)
	suite.client.BaseURL = url
}

func (suite *CalendlyClientTestSuite) TearDownTest() {
	suite.server.Close()
}

func (suite *CalendlyClientTestSuite) TestClient_TestNewRequest() {
	assert := assert.New(suite.T())

	inURL, outURL := "echo", suite.server.URL+"/echo"
	inBody, outBody := &Echo{Email: "echo@echo.com"}, `{"email":"echo@echo.com"}`+"\n"

	req, err := suite.client.NewRequest("GET", inURL, inBody)
	assert.Nil(err)

	assert.Equal(req.URL.String(), outURL)

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	assert.Equal(outBody, string(body))

	// test that default user-agent is attached to the request
	assert.Equal(req.Header.Get("User-Agent"), suite.client.UserAgent)
}

func (suite *CalendlyClientTestSuite) TestNewRequest_invalidJSON() {
	assert := assert.New(suite.T())
	type T struct {
		A map[interface{}]interface{}
	}
	_, err := suite.client.NewRequest("GET", "/", &T{})
	assert.NotNil(err)
	assert.IsType(&json.UnsupportedTypeError{}, err)
}

func (suite *CalendlyClientTestSuite) TestDo() {
	assert := assert.New(suite.T())
	type foo struct {
		A string
	}

	suite.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := suite.client.Get(".")
	body := new(foo)
	suite.client.Do(context.Background(), req, body)

	want := &foo{"a"}
	assert.Equal(want, body)
}

func (suite *CalendlyClientTestSuite) TestDo_httpError() {
	assert := assert.New(suite.T())

	suite.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := suite.client.Get(".")
	resp, _ := suite.client.Do(context.Background(), req, nil)

	assert.Equal(400, resp.StatusCode)
}

func (suite *CalendlyClientTestSuite) TestClient_SetBaseUrl() {
	assert := assert.New(suite.T())
	expectedBaseUrl, _ := url.Parse("https://calendly.com/api/v2")

	err := suite.client.SetBaseURL("https://calendly.com/api/v2")

	assert.Nil(err)
	assert.Equal(suite.client.BaseURL, expectedBaseUrl)
}

func (suite *CalendlyClientTestSuite) TestClient_SetUserAgent() {
	assert := assert.New(suite.T())
	expectedUserAgent := "go-calendly-0.2.0"

	suite.client.SetUserAgent(expectedUserAgent)

	assert.Equal(suite.client.UserAgent, expectedUserAgent)
}


func (suite *CalendlyClientTestSuite) TestAddUrlOptions() {
	assert := assert.New(suite.T())
	type opts struct {
		HasName bool `url:"has_name"`
	}
	testCases := []struct {
		name      string
		s         string
		opt       interface{}
		want      string
		wantError bool
	}{
		{"TestNilOpts", "123", nil, "123", false},
		{"TestInvalidUrl", "http://[::1]a", nil, "http://[::1]a", true},
		{"TestEmptyOpts", "http://[::1]:80", struct{}{}, "http://[::1]:80", false},
		{"TestOptsWithNoTags", "http://[::1]:80", struct{ name string }{"hello"}, "http://[::1]:80", false},
		{"TestOptsWithTags", "http://[::1]:80", opts{true}, "http://[::1]:80?has_name=true", false},
		{"TestInvalidOpts", "http://[::1]:80", "::", "http://[::1]:80", true},
	}
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			v, err := addUrlOptions(tc.s, tc.opt)
			assert.Equal(tc.want, v)
			if tc.wantError {
				assert.NotNil(err)
				return
			}
			assert.Nil(err)
		})

	}
}