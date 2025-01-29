package tado

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestNewClient(t *testing.T) {
	var tests = []struct {
		name  string
		input *http.Client
	}{
		{"with nil http client", nil},
		{"with non-nil http client", &http.Client{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewClient(tt.input)
			assert.NotNil(t, got)
			assert.NotNil(t, got.client)

			if tt.input != nil {
				assert.Equal(t, tt.input, got.client)
			}
		})
	}
}

func TestClient_WithOAuthClient(t *testing.T) {
	mockOAuthConfig := &oauth2.Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://example.com/oauth/token",
		},
	}

	mockToken := &oauth2.Token{
		AccessToken: "test-access-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	var tests = []struct {
		name             string
		inputOAuthConfig *oauth2.Config
		inputToken       *oauth2.Token
	}{
		{"with non-nil oauth config and non-nil token", mockOAuthConfig, mockToken},
		{"with nil oauth config and non-nil token", nil, mockToken},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(nil).WithOAuthClient(context.Background(), tt.inputOAuthConfig, tt.inputToken)

			assert.NotNil(t, client.client)
			assert.Equal(t, reflect.TypeOf(client.client.Transport), reflect.TypeOf(&oauth2.Transport{}))
		})
	}
}

func TestClient_NewRequest(t *testing.T) {
	client := NewClient(nil)

	type T struct {
		A map[interface{}]interface{}
	}

	var tests = []struct {
		name      string
		inputPath string
		inputBody interface{}
	}{
		{"with relative path", "foo", nil},
		{"with absolute path", "/foo", nil},
		{"with trailing slash", "foo/", nil},
		{"with body", "foo", User{ID: "test-id"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := client.NewRequest("GET", tt.inputPath, tt.inputBody)
			assert.Equal(t, DefaultBaseURL+strings.TrimPrefix(tt.inputPath, "/"), req.URL.String())

			if tt.inputBody != nil {
				assert.NotNil(t, req.Body)

				body, _ := io.ReadAll(req.Body)
				assert.Equal(t, `{"id":"test-id"}`+"\n", string(body))
			}
		})
	}
}

func TestClient_NewRequest_invalidJSON(t *testing.T) {
	client := NewClient(nil)

	type T struct {
		A map[interface{}]interface{}
	}

	_, err := client.NewRequest("GET", "foo", &T{})

	assert.Error(t, err)
	assert.ErrorIs(t, err, err.(*json.UnsupportedTypeError))
}

func TestClient_NewRequest_badURL(t *testing.T) {
	client := NewClient(nil)

	_, err := client.NewRequest("GET", ":", nil)

	assert.Error(t, err)
	assert.ErrorIs(t, err, err.(*url.Error))
}

func TestClient_NewRequest_badMethod(t *testing.T) {
	client := NewClient(nil)

	_, err := client.NewRequest("BAD METHOD", "foo", nil)

	assert.Error(t, err)
}
