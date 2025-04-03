package tado

import (
	"context"
	"net/http"
)

// UserService handles communication with the user-related methods of the Tado
// API.
type UserService service

// User represents a Tado user.
type User struct {
	Name          string         `json:"name,omitempty"`
	Email         string         `json:"email,omitempty"`
	Username      string         `json:"username,omitempty"`
	ID            string         `json:"id,omitempty"`
	Homes         []BareHome     `json:"homes,omitempty"`
	Locale        string         `json:"locale,omitempty"`
	MobileDevices []MobileDevice `json:"mobileDevices,omitempty"`
}

type BareHome struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Get returns the authenticated user.
func (s *UserService) Get() (*User, error) {
	req, err := s.client.NewRequest(http.MethodGet, "me", nil)
	if err != nil {
		return nil, err
	}

	var user *User
	_, err = s.client.Do(context.Background(), req, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
