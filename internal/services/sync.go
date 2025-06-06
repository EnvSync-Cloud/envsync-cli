package services

import (
	"fmt"

	"github.com/EnvSync-Cloud/envsync-cli/internal/config"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services/requests"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services/responses"
	"resty.dev/v3"
)

type SyncService interface {
	PullEnv(appId, envType string) (map[string]string, error)
}

type sync struct {
	client *resty.Client
}

func NewSyncService() SyncService {
	cfg := config.New()
	client := resty.New().
		SetBaseURL(cfg.BackendURL).
		SetHeader("Content-Type", "application/json").
		SetAuthToken(cfg.AccessToken)

	return &sync{
		client: client,
	}
}

func (s *sync) PullEnv(appId, envType string) (map[string]string, error) {
	var env responses.EnvVariableList

	res, err := s.client.
		R().
		SetResult(&env).
		SetBody(requests.EnvVariableRequest{
			AppID:   appId,
			EnvType: envType,
		}).
		Post("/env")

	if err != nil {
		return nil, err
	}

	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}

	return env.ToMap(), nil
}
