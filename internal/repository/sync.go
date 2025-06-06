package repository

import (
	"fmt"

	"github.com/EnvSync-Cloud/envsync-cli/internal/config"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/requests"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/responses"
	"resty.dev/v3"
)

type SyncRepository interface {
	PullEnv() (responses.EnvVariableList, error)
}

type syncRepo struct {
	client    *resty.Client
	appID     string
	envTypeID string
}

func NewSyncService(appID, envTypeID string) SyncRepository {
	cfg := config.New()
	client := resty.New().
		SetBaseURL(cfg.BackendURL).
		SetHeader("Content-Type", "application/json").
		SetAuthToken(cfg.AccessToken)

	return &syncRepo{
		client:    client,
		appID:     appID,
		envTypeID: envTypeID,
	}
}

func (s *syncRepo) PullEnv() (responses.EnvVariableList, error) {
	var env responses.EnvVariableList

	res, err := s.client.
		R().
		SetResult(&env).
		SetBody(requests.EnvVariableRequest{
			AppID:     s.appID,
			EnvTypeID: s.envTypeID,
		}).
		Post("/env")

	if err != nil {
		return nil, err
	}

	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}

	return env, nil
}
