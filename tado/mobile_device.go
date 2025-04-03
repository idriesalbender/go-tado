package tado

import (
	"context"
	"fmt"
)

// MobileDeviceService handles communication with the mobile device-related
// methods of the Tado API.
type MobileDeviceService service

// MobileDevice represents a Tado mobile device.
type MobileDevice struct {
	Name     string               `json:"name,omitempty"`
	ID       int                  `json:"id,omitempty"`
	Settings MobileDeviceSettings `json:"settings,omitempty"`
	Location struct {
		Stale           bool `json:"stale,omitempty"`
		AtHome          bool `json:"atHome,omitempty"`
		BearingFromHome struct {
			Degrees float64 `json:"degrees,omitempty"`
			Radians float64 `json:"radians,omitempty"`
		} `json:"bearingFromHome,omitempty"`
		RelativeDistanceFromHomeFence float64 `json:"relativeDistanceFromHomeFence,omitempty"`
	} `json:"location,omitempty"`
	DeviceMetadata struct {
		Platform  string `json:"platform,omitempty"`
		OsVersion string `json:"osVersion,omitempty"`
		Model     string `json:"model,omitempty"`
		Locale    string `json:"locale,omitempty"`
	} `json:"deviceMetadata,omitempty"`
}

type MobileDeviceSettings struct {
	GeoTrackingEnabled          bool `json:"geoTrackingEnabled,omitempty"`
	SpecialOffersEnabled        bool `json:"specialOffersEnabled,omitempty"`
	OnDemandLogRetrievalEnabled bool `json:"onDemandLogRetrievalEnabled,omitempty"`
	PushNotifications           struct {
		LowBatteryReminder          bool `json:"lowBatteryReminder,omitempty"`
		AwayModeReminder            bool `json:"awayModeReminder,omitempty"`
		HomeModeReminder            bool `json:"homeModeReminder,omitempty"`
		OpenWindowReminder          bool `json:"openWindowReminder,omitempty"`
		EnergySavingsReportReminder bool `json:"energySavingsReportReminder,omitempty"`
		IncidentDetection           bool `json:"incidentDetection,omitempty"`
		EnergyIqReminder            bool `json:"energyIqReminder,omitempty"`
		TariffHighPriceAlert        bool `json:"tariffHighPriceAlert,omitempty"`
		TariffLowPriceAlert         bool `json:"tariffLowPriceAlert,omitempty"`
	} `json:"pushNotifications"`
}

// List returns a list of all mobile devices for the provided home ID.
func (s *MobileDeviceService) List(ctx context.Context, id int) (*[]MobileDevice, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("homes/%d/mobileDevices", id), nil)
	if err != nil {
		return nil, err
	}

	var mobileDevices *[]MobileDevice
	_, err = s.client.Do(ctx, req, &mobileDevices)
	if err != nil {
		return nil, err
	}

	return mobileDevices, nil
}

// Get returns the mobile device with the given ID for the provided home ID.
func (s *MobileDeviceService) Get(ctx context.Context, homeID, deviceID int) (*MobileDevice, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("homes/%d/mobileDevices/%d", homeID, deviceID), nil)
	if err != nil {
		return nil, err
	}

	var mobileDevice *MobileDevice
	_, err = s.client.Do(ctx, req, &mobileDevice)
	if err != nil {
		return nil, err
	}

	return mobileDevice, nil
}

// Delete deletes the relationship between the given mobile device and home.
func (s *MobileDeviceService) Delete(ctx context.Context, homeID, deviceID int) error {
	req, err := s.client.NewRequest("DELETE", fmt.Sprintf("homes/%d/mobileDevices/%d", homeID, deviceID), nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetSettings returns the settings of the mobile device with the given ID for the provided home ID.
func (s *MobileDeviceService) GetSettings(ctx context.Context, homeID, deviceID int) (*MobileDeviceSettings, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("homes/%d/mobileDevices/%d/settings", homeID, deviceID), nil)
	if err != nil {
		return nil, err
	}

	var settings *MobileDeviceSettings
	_, err = s.client.Do(ctx, req, &settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// UpdateSettings updates the settings of the mobile device with the given ID for the provided home ID.
func (s *MobileDeviceService) UpdateSettings(ctx context.Context, homeID, deviceID int, settings MobileDeviceSettings) (*MobileDeviceSettings, error) {
	req, err := s.client.NewRequest("PUT", fmt.Sprintf("homes/%d/mobileDevices/%d/settings", homeID, deviceID), settings)
	if err != nil {
		return nil, err
	}

	var settings2 *MobileDeviceSettings
	_, err = s.client.Do(ctx, req, &settings2)
	if err != nil {
		return nil, err
	}

	return settings2, nil
}
