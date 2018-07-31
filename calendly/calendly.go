package calendly

import (
	"net/http"
	"net/url"
	"reflect"

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"golang.org/x/net/context/ctxhttp"
	"io"
	"io/ioutil"
)

const (
	libraryVersion = "0.1.0"
	defaultBaseURL = "https://calendly.com/api/v1/"
	userAgent      = "go-calendly-" + libraryVersion
	mediaType      = "application/json"
	testRoute      = "echo"
)

type Client struct {
	// HTTP client used to communicate with the DO API.
	client *http.Client

	common ApiService

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// Event Types Service
	EventTypes EventTypesService
}

// Response is a Calendly response. This wraps the standard http.Response returned
// from Calendly.
type Response struct {
	*http.Response
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Message string `json:"message"`

	// RequestID returned from the API, useful to contact support.
	RequestID string `json:"request_id"`
}

type ApiService struct {
	client *Client
}

// NewClient returns a new DigitalOcean API client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.common.client = c
	c.EventTypes = EventTypesService{c}

	return c
}

// SetUserAgent is a client option for setting the user agent.
func (c *Client) SetUserAgent(ua string) {
	c.UserAgent = ua
}

// SetBaseURL is a client option for setting the base URL.
func (c *Client) SetBaseURL(bUrl string) error {
	u, err := url.Parse(bUrl)
	if err != nil {
		return err
	}

	c.BaseURL = u
	return nil
}

// NewRequest creates an API request.
// A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", mediaType)
	}
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := ctxhttp.Do(ctx, c.client, req)

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		} else {
			io.CopyN(ioutil.Discard, resp.Body, 512)
		}
	}()

	response := newResponse(resp)
	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, err
}

// Convenient shorthand for GET requests
func (c *Client) Get(urlStr string) (*http.Request, error) {
	return c.NewRequest(http.MethodGet, urlStr, nil)
}

// Convenient shorthand for POST requests
func (c *Client) Post(urlStr string, body interface{}) (*http.Request, error) {
	return c.NewRequest(http.MethodPost, urlStr, body)
}

// Convenient shorthand for PUT requests
func (c *Client) Put(urlStr string, body interface{}) (*http.Request, error) {
	return c.NewRequest(http.MethodPut, urlStr, body)
}

// Convenient shorthand for DELETE requests
func (c *Client) Delete(urlStr string) (*http.Request, error) {
	return c.NewRequest(http.MethodDelete, urlStr, nil)
}

// CheckResponse checks the API response for errors, and returns them if present.
// A response is considered an
// error if it has a status code outside the 200 range.
// API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.
// Any other response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}

func (r *ErrorResponse) Error() string {
	if r.RequestID != "" {
		return fmt.Sprintf("%v %v: %d (request %q) %v",
			r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.RequestID, r.Message)
	}
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

func (c *Client) Echo(ctx context.Context) (*Echo, *Response, error) {
	req, err := c.NewRequest(http.MethodGet, testRoute, nil)
	if err != nil {
		return nil, nil, err
	}

	e := &Echo{}
	resp, err := c.Do(ctx, req, e)
	if err != nil {
		return nil, resp, err
	}

	return e, resp, nil
}

// Test Authentication Token
// Use this endpoint to test your Authentication Token.
type Echo struct {
	Email string `json:"email"`
}

func (r *Echo) String() string {
	return fmt.Sprintf("EchoResponse: email=%v", r.Email)
}

func addUrlOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}

// newResponse creates a new Response for the provided http.Response
func newResponse(r *http.Response) *Response {
	response := Response{Response: r}
	return &response
}
