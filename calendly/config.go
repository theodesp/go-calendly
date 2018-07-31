package calendly

import (
	"fmt"
	"net/http"
)

const DefaultHeaderTokenKey = "X-Token"

// Config represents an OAuth1 consumer's (client's) key and secret, the
// callback URL, and the provider Endpoint to which the consumer corresponds.
type Config struct {
	// API Key (Client Identifier)
	ApiKey string

	// Header identifier for passing the API key
	HeaderKey string
}

// NewTokenAuthClient returns a new http Client which signs requests via header Token.
func NewTokenAuthClient(config *Config) *http.Client {
	return &http.Client{Transport: &Transport{Base: http.DefaultTransport, config: config}}
}

// Transport is an http.RoundTripper which makes Authenticated HTTP requests. It
// wraps a base RoundTripper and adds an API header using the
// token from the config.
//
// Transport is a low-level component, most users should use NewClient to create
// an http.Client instead.
type Transport struct {
	// Base is the base RoundTripper used to make HTTP requests. If nil, then
	// http.DefaultTransport is used
	Base http.RoundTripper

	// Config that is used for this transport
	config *Config
}

// RoundTrip authorizes the request by passing the API token to the request header
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.config == nil {
		return nil, fmt.Errorf("go-calendly: Transport's config is nil")
	}

	if t.config.HeaderKey == "" {
		t.config.HeaderKey = DefaultHeaderTokenKey
	}

	if t.config.ApiKey == "" {
		return nil, fmt.Errorf("go-calendly: API Key token is missing")
	}

	req.Header.Set(t.config.HeaderKey, t.config.ApiKey)
	return t.Base.RoundTrip(req)
}
