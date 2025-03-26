package tado

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

const (
	DefaultBaseURL   = "https://my.tado.com/api/v2/"
	DefaultUserAgent = "go-tado"
)

var ErrNonNilContext = errors.New("context must not be nil")

// Client is the main client for interacting with the Tado API.
type Client struct {
	authenticator Authenticator
	client        *http.Client
	BaseURL       *url.URL
	UserAgent     string
	common        service

	User         *UserService
	Home         *HomeService
	MobileDevice *MobileDeviceService
}

type service struct {
	client *Client
}

type ClientOption func(*Client)

func WithAuthenticator(auth Authenticator) ClientOption {
	return func(c *Client) {
		c.authenticator = auth
	}
}

// NewClient returns a new Client instance with the given options.
//
// The Client returned by NewClient is not initialized until the first call to
// a method that requires authentication. If no Authenticator is provided, a
// DeviceAuthenticator with the default OAuth2 configuration is used.
func NewClient(opts ...ClientOption) *Client {
	tc := &Client{}
	for _, opt := range opts {
		opt(tc)
	}

	if tc.authenticator == nil {
		tc.authenticator = NewDeviceAuthenticator(nil)
	}

	tc.initialize()
	return tc
}

// initialize sets up the client with default values and initializes the
// services.
func (c *Client) initialize() {
	if c.client == nil {
		token, err := c.authenticator.TokenSource(context.Background())
		if err != nil {
			panic(err)
		}

		c.client = oauth2.NewClient(context.Background(), token)
	}

	if c.BaseURL == nil {
		c.BaseURL, _ = url.Parse(DefaultBaseURL)
	}

	if c.UserAgent == "" {
		c.UserAgent = DefaultUserAgent
	}

	c.common.client = c

	c.User = (*UserService)(&c.common)
	c.Home = (*HomeService)(&c.common)
	c.MobileDevice = (*MobileDeviceService)(&c.common)
}

// clone returns a copy of the client. Must be initialized before use using
// Client.initialize.
// func (c *Client) clone() *Client {
// 	clone := Client{
// 		client:    &http.Client{},
// 		BaseURL:   c.BaseURL,
// 		UserAgent: c.UserAgent,
// 		common:    c.common,
// 	}

// 	if c.client != nil {
// 		clone.client.Transport = c.client.Transport
// 		clone.client.CheckRedirect = c.client.CheckRedirect
// 		clone.client.Jar = c.client.Jar
// 		clone.client.Timeout = c.client.Timeout
// 	}

// 	return &clone
// }

type RequestOption func(req *http.Request)

func (c *Client) NewRequest(method, path string, body interface{}, opts ...RequestOption) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	url, err := c.BaseURL.Parse(strings.TrimPrefix(path, "/")) // trim prefix to prevent absolute paths from overwriting the base URL
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

// Response is a Tado API response. It wraps the standard http.Response returned
// from Tado.
type Response struct {
	*http.Response
}

// newResponse returns a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// bareDo sends an API request using the provided http.Client (`caller`) and
// lets you handle the http.Response on your own.
//
// The provided ctx must not be nil. If it is, bareDo returns ErrNonNilContext.
func (c *Client) bareDo(ctx context.Context, caller *http.Client, req *http.Request) (*Response, error) {
	if ctx == nil {
		return nil, ErrNonNilContext
	}

	req = req.WithContext(ctx)

	res, err := caller.Do(req)
	var response *Response
	if res != nil {
		response = newResponse(res)
	}

	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = url.String()
				return response, e
			}
		}

		return response, err
	}

	return response, err
}

// BareDo sends an API request and lets you handle the http.Response on your
// own.
//
// The provided ctx must not be nil. If it is, BareDo returns ErrNonNilContext.
func (c *Client) BareDo(ctx context.Context, req *http.Request) (*Response, error) {
	return c.bareDo(ctx, c.client, req)
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an error
// if an API error has occurred. If v implements the io.Writer interface, the
// raw response body will be written to v, without attempting to decode it. If v
// is nil and no error occurs, the response is returned as is.
//
// The provided ctx must not be nil. If it is, Do returns ErrNonNilContext.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	res, err := c.BareDo(ctx, req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, res.Body)
	default:
		derr := json.NewDecoder(res.Body).Decode(v)
		if derr == io.EOF {
			derr = nil // ignore EOF errors caused by empty response body
		}
		if derr != nil {
			err = derr
		}
	}

	return res, err
}

// roundTripperFunc creates a RoundTripper (transport).
// type roundTripperFunc func(*http.Request) (*http.Response, error)

// func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
// 	return fn(r)
// }
