package repository

import (
	"fmt"

	"github.com/EnvSync-Cloud/envsync-cli/internal/config"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/requests"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/responses"
	"resty.dev/v3"
)

type SyncRepository interface {
	GetAllEnv() ([]responses.EnvironmentVariables, error)
	BatchCreateEnv(env requests.BatchSyncEnvRequest) error
	BatchUpdateEnv(env requests.BatchSyncEnvRequest) error
	BatchDeleteEnv(env requests.BatchSyncEnvRequest) error
}

type syncRepo struct {
	client    *resty.Client
	appID     string
	envTypeID string
}

func NewSyncRepository(appID, envTypeID string) SyncRepository {
	cfg := config.New()
	client := resty.New().
		SetDisableWarn(true).
		SetBaseURL(cfg.BackendURL).
		SetHeader("Content-Type", "application/json").
		SetAuthToken(cfg.AccessToken)

	return &syncRepo{
		client:    client,
		appID:     appID,
		envTypeID: envTypeID,
	}
}

func (s *syncRepo) GetAllEnv() ([]responses.EnvironmentVariables, error) {
	var env []responses.EnvironmentVariables

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

func (s *syncRepo) BatchCreateEnv(env requests.BatchSyncEnvRequest) error {
	res, err := s.client.
		R().
		SetBody(env).
		Put("/env/batch")

	if err != nil {
		return err
	}

	if res.StatusCode() != 200 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}

	return nil
}

func (s *syncRepo) BatchUpdateEnv(env requests.BatchSyncEnvRequest) error {
	res, err := s.client.
		R().
		SetBody(env).
		Patch("/env/batch")

	if err != nil {
		return err
	}

	if res.StatusCode() != 200 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}

	return nil
}

func (s *syncRepo) BatchDeleteEnv(env requests.BatchSyncEnvRequest) error {
	res, err := s.client.
		R().
		SetBody(env).
		Delete("/env/batch")

	if err != nil {
		return err
	}

	if res.StatusCode() != 200 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}

	return nil
}
