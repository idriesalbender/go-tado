package tado

import (
	"context"
	"fmt"
	"time"
)

// Presence represents a Tado presence.
type Presence string

// HomeService handles communication with the home-related methods of the Tado
// API.
type HomeService service

const (
	PresenceHome Presence = "HOME"
	PresenceAway Presence = "AWAY"
)

// Home represents a Tado home.
type Home struct {
	ID                         int       `json:"id"`
	Name                       string    `json:"name"`
	DateTimeZone               string    `json:"dateTimeZone"`
	DateCreated                time.Time `json:"dateCreated"`
	TemperatureUnit            string    `json:"temperatureUnit"`
	Partner                    string    `json:"partner"`
	SimpleSmartScheduleEnabled bool      `json:"simpleSmartScheduleEnabled"`
	AwayRadiusInMeters         float64   `json:"awayRadiusInMeters"`
	InstallationCompleted      bool      `json:"installationCompleted"`
	IncidentDetection          struct {
		Supported bool `json:"supported"`
		Enabled   bool `json:"enabled"`
	} `json:"incidentDetection"`
	Generation              string        `json:"generation"`
	ZonesCount              int           `json:"zonesCount"`
	Language                string        `json:"language"`
	PreventFromSubscribing  bool          `json:"preventFromSubscribing"`
	Skills                  []interface{} `json:"skills"`
	ChristmasModeEnabled    bool          `json:"christmasModeEnabled"`
	ShowAutoAssistReminders bool          `json:"showAutoAssistReminders"`
	ContactDetails          struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"contactDetails"`
	Address struct {
		AddressLine1 string      `json:"addressLine1"`
		AddressLine2 interface{} `json:"addressLine2"`
		ZipCode      string      `json:"zipCode"`
		City         string      `json:"city"`
		State        interface{} `json:"state"`
		Country      string      `json:"country"`
	} `json:"address"`
	Geolocation struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"geolocation"`
	ConsentGrantSkippable bool     `json:"consentGrantSkippable"`
	EnabledFeatures       []string `json:"enabledFeatures"`
	IsAirComfortEligible  bool     `json:"isAirComfortEligible"`
	IsBalanceAcEligible   bool     `json:"isBalanceAcEligible"`
	IsEnergyIqEligible    bool     `json:"isEnergyIqEligible"`
	IsHeatSourceInstalled bool     `json:"isHeatSourceInstalled"`
	IsHeatPumpInstalled   bool     `json:"isHeatPumpInstalled"`
}

// State represents the state of a Tado home.
type State struct {
	Presence       Presence `json:"presence"`
	PresenceLocked bool     `json:"presenceLocked"`
}

// AirComfort represents the air comfort of a Tado home.
type AirComfort struct {
	Freshness struct {
		Value          string    `json:"value"`
		LastOpenWindow time.Time `json:"lastOpenWindow"`
	} `json:"freshness"`
	Comfort []struct {
		RoomID           int    `json:"roomId"`
		TemperatureLevel string `json:"temperatureLevel"`
		HumidityLevel    string `json:"humidityLevel"`
		Coordinate       struct {
			Radial  float64 `json:"radial"`
			Angular int     `json:"angular"`
		} `json:"coordinate"`
	} `json:"comfort"`
}

// HeatingSystem represents the various heating systems in a Tado home.
type HeatingSystem struct {
	Boiler struct {
		Present bool `json:"present"`
		ID      int  `json:"id"`
		Found   bool `json:"found"`
	} `json:"boiler"`
	UnderfloorHeating struct {
		Present bool `json:"present"`
	} `json:"underfloorHeating"`
}

// FlowTemperaturOptimization represents the flow temperature optimization of a Tado home.
type FlowTemperaturOptimization struct {
	HasMultipleBoilerControlDevices bool `json:"hasMultipleBoilerControlDevices"`
	MaxFlowTemperature              int  `json:"maxFlowTemperature"`
	MaxFlowTemperatureConstraints   struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"maxFlowTemperatureConstraints"`
	AutoAdaptation struct {
		Enabled            bool `json:"enabled"`
		MaxFlowTemperature int  `json:"maxFlowTemperature"`
	} `json:"autoAdaptation"`
	OpenThermDeviceSerialNumber string `json:"openThermDeviceSerialNumber"`
}

// Weather represents the weather of a Tado home.
type Weather struct {
	SolarIntensity struct {
		Type       string    `json:"type"`
		Percentage int       `json:"percentage"`
		Timestamp  time.Time `json:"timestamp"`
	} `json:"solarIntensity"`
	OutsideTemperature struct {
		Celsius    float64   `json:"celsius"`
		Fahrenheit float64   `json:"fahrenheit"`
		Timestamp  time.Time `json:"timestamp"`
		Type       string    `json:"type"`
		Precision  struct {
			Celsius    float64 `json:"celsius"`
			Fahrenheit float64 `json:"fahrenheit"`
		} `json:"precision"`
	} `json:"outsideTemperature"`
	WeatherState struct {
		Type      string    `json:"type"`
		Value     string    `json:"value"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"weatherState"`
}

// Get returns the home with the given ID.
func (s *HomeService) Get(id int) (*Home, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("homes/%d", id), nil)
	if err != nil {
		return nil, err
	}

	var home *Home
	_, err = s.client.Do(context.Background(), req, &home)
	if err != nil {
		return nil, err
	}

	return home, nil
}

// GetAirComfort returns the air comfort of the home with the given ID.
func (s *HomeService) GetAirComfort(id int) (*AirComfort, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("homes/%d/airComfort", id), nil)
	if err != nil {
		return nil, err
	}

	var airComfort *AirComfort
	_, err = s.client.Do(context.Background(), req, &airComfort)
	if err != nil {
		return nil, err
	}

	return airComfort, nil
}

// GetHeatSystem returns the heating system of the home with the given ID.
func (s *HomeService) GetHeatingSystem(id int) (*HeatingSystem, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("homes/%d/heatingSystem", id), nil)
	if err != nil {
		return nil, err
	}

	var heatingSystem *HeatingSystem
	_, err = s.client.Do(context.Background(), req, &heatingSystem)
	if err != nil {
		return nil, err
	}

	return heatingSystem, nil
}

// GetFlowTemperatureOptimization returns the flow temperature optimization of the home with the given ID.
func (s *HomeService) GetFlowTemperatureOptimization(id int) (*FlowTemperaturOptimization, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("homes/%d/flowTemperatureOptimization", id), nil)
	if err != nil {
		return nil, err
	}

	var flowTemperatureOptimization *FlowTemperaturOptimization
	_, err = s.client.Do(context.Background(), req, &flowTemperatureOptimization)
	if err != nil {
		return nil, err
	}

	return flowTemperatureOptimization, nil
}

// GetWeather returns the weather of the home with the given ID.
func (s *HomeService) GetWeather(id int) (*Weather, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("homes/%d/weather", id), nil)
	if err != nil {
		return nil, err
	}

	var weather *Weather
	_, err = s.client.Do(context.Background(), req, &weather)
	if err != nil {
		return nil, err
	}

	return weather, nil
}

// GetState returns the state of the home with the given ID.
func (s *HomeService) GetState(id int) (*State, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("homes/%d/state", id), nil)
	if err != nil {
		return nil, err
	}

	var state *State
	_, err = s.client.Do(context.Background(), req, &state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

// SetState sets the state of the home with the given ID.
func (s *HomeService) SetState(id int, presence Presence) error {
	req, err := s.client.NewRequest("PUT", fmt.Sprintf("homes/%d/presenceLock", id), &map[string]string{"homePresence": string(presence)})
	if err != nil {
		return err
	}

	_, err = s.client.Do(context.Background(), req, nil)
	if err != nil {
		return err
	}

	return nil
}
