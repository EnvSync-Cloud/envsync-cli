package services

import (
	"fmt"

	"github.com/EnvSync-Cloud/envsync-cli/internal/config"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services/responses"
	"resty.dev/v3"
)

type AuthService interface {
	LoginDeviceCode() (responses.DeviceCodeResponse, error)
	LoginToken(deviceCode, clientID, authDomain string) (string, error)
}

type auth struct {
	client *resty.Client
	cfg    config.AppConfig
}

func NewAuthService() AuthService {
	cfg := config.New()
	client := resty.New().SetBaseURL(cfg.BackendURL)

	return &auth{
		client: client,
		cfg:    cfg,
	}
}

func (s *auth) LoginDeviceCode() (responses.DeviceCodeResponse, error) {
	var resBody responses.DeviceCodeResponse

	res, err := s.client.
		R().
		SetResult(&resBody).
		Get("/access/cli")

	if err != nil {
		return responses.DeviceCodeResponse{}, fmt.Errorf("failed to get login URL: %w", err)
	}

	if res.StatusCode() != 201 {
		return responses.DeviceCodeResponse{}, fmt.Errorf("unexpected status code while fetching login URL: %d", res.StatusCode())
	}

	return resBody, nil
}

func (s *auth) LoginToken(deviceCode, clientID, authDomain string) (string, error) {
	var resBody responses.LoginTokenResponse

	res, err := s.client.
		SetBaseURL("https://" + authDomain).
		R().
		SetResult(&resBody).
		SetFormData(map[string]string{
			"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
			"device_code": deviceCode,
			"client_id":   clientID,
		}).
		Post("/oauth/token")

	if err != nil {
		return "", fmt.Errorf("failed to get login token: %w", err)
	}

	if res.StatusCode() != 200 {
		return "", fmt.Errorf("unexpected status code while fetching login token: %d", res.StatusCode())
	}

	return resBody.AccessToken, nil
}
