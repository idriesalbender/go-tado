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
	DefaultBaseURL      = "https://my.tado.com/api/v2/"
	DefaultUserAgent    = "go-tado"
	TadoAPIClientID     = "public-api-preview"
	TadoAPIClientSecret = "4HJGRffVR8xb3XdEUQpjgZ1VplJi6Xgw"
	TadoAPIAuthURL      = "https://auth.tado.com/oauth/token"
)

var ErrNonNilContext = errors.New("context must not be nil")

var DefaultOAuth2Config = &oauth2.Config{
	ClientID:     TadoAPIClientID,
	ClientSecret: TadoAPIClientSecret,
	Endpoint: oauth2.Endpoint{
		TokenURL: TadoAPIAuthURL,
	},
}

// Client is the main client for interacting with the Tado API.
type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string
	common    service

	auth struct {
		config      *oauth2.Config
		token       *oauth2.Token
		credentials *credentials
	}

	User *UserService
	Home *HomeService
}

type service struct {
	client *Client
}

type credentials struct {
	username string
	password string
}

// NewClient returns a new Tado API client. If httpClient is nil, a default
// http.Client is used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	hc := *httpClient

	client := &Client{client: &hc}
	client.initialize()
	return client
}

// WithOAuthClient sets up the client to use an OAuth2 client with the provided
// oauth2.Config and token. If config is nil, the default Tado OAuth2
// configuration is used. Token must not be nil.
func (c *Client) WithOAuthClient(ctx context.Context, config *oauth2.Config, token *oauth2.Token) *Client {
	clone := c.clone()
	defer clone.initialize()

	if config == nil {
		config = DefaultOAuth2Config
	}

	tokenSource := config.TokenSource(ctx, token)
	clone.client = oauth2.NewClient(ctx, tokenSource)

	return clone
}

// initialize sets up the client with default values and initializes the
// services.
func (c *Client) initialize() {
	if c.client == nil {
		c.client = &http.Client{}
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
}

// clone returns a copy of the client. Must be initialized before use using
// Client.initialize.
func (c *Client) clone() *Client {
	clone := Client{
		client:    &http.Client{},
		BaseURL:   c.BaseURL,
		UserAgent: c.UserAgent,
		common:    c.common,
	}

	if c.client != nil {
		clone.client.Transport = c.client.Transport
		clone.client.CheckRedirect = c.client.CheckRedirect
		clone.client.Jar = c.client.Jar
		clone.client.Timeout = c.client.Timeout
	}

	return &clone
}

type RequestOption func(req *http.Request)

func (c *Client) NewRequest(method, path string, body interface{}, opts ...RequestOption) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	url, err := c.BaseURL.Parse(path)
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
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}
