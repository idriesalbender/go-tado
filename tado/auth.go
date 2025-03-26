package tado

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

// Authenticator defines the interface for Tado API authentication mechanisms
type Authenticator interface {
	TokenSource(context.Context) (oauth2.TokenSource, error)
}

var TadoDeviceAuthClientID = "1bb50063-6b0c-4d11-bd99-387f4a91cc46"
var TadoDeviceAuthURL = "https://login.tado.com/oauth2/device_authorize"
var TadoDeviceAuthTokenURL = "https://login.tado.com/oauth2/token"

var TadoDeviceAuthDefaultOAuth2Config = &oauth2.Config{
	ClientID: TadoDeviceAuthClientID,
	Endpoint: oauth2.Endpoint{
		DeviceAuthURL: TadoDeviceAuthURL,
		TokenURL:      TadoDeviceAuthTokenURL,
	},
	Scopes: []string{"offline-access"},
}

// DeviceAuthenticator provides an authentication mechanism using the OAuth2
// device authorization flow for the Tado API.
//
// This authenticator uses an oauth2.Config to obtain a device code and prompt
// the user to visit a verification URI. Once the user authorizes the device, it
// provides a TokenSource for accessing the Tado API.
//
// The DeviceAuthenticator can be initialized with a custom oauth2.Config, or it
// defaults to TadoDeviceAuthDefaultOAuth2Config if none is provided.
type DeviceAuthenticator struct {
	config *oauth2.Config
}

// NewDeviceAuthenticator creates a new DeviceAuthenticator.
//
// If the provided config is nil, it defaults to
// TadoDeviceAuthDefaultOAuth2Config.
func NewDeviceAuthenticator(config *oauth2.Config) *DeviceAuthenticator {
	c := config

	if c == nil {
		c = TadoDeviceAuthDefaultOAuth2Config
	}

	return &DeviceAuthenticator{
		config: c,
	}
}

// TokenSource implements the Authenticator interface.
//
// It is a blocking call that asks the user to visit the verification URI and
// enter the user code. Once the user has done so, it returns a TokenSource for
// the authenticated user.
func (a *DeviceAuthenticator) TokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	deviceCode, err := a.config.DeviceAuth(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Visit %s to log in.\n", deviceCode.VerificationURIComplete)
	fmt.Printf("Enter the code: %s\n", deviceCode.UserCode)

	token, err := a.config.DeviceAccessToken(ctx, deviceCode)
	if err != nil {
		return nil, err
	}

	return a.config.TokenSource(ctx, token), nil
}
