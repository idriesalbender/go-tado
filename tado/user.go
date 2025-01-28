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
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	ID       string `json:"id"`
	Homes    []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"homes"`
	Locale string `json:"locale"`
	// MobileDevices []MobileDevice
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
