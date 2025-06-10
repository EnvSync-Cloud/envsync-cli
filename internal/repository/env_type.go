package repository

import (
	"fmt"

	"resty.dev/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/config"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/responses"
)

type EnvTypeRepository interface {
	GetAll() ([]responses.EnvTypeResponse, error)
}

type envTypeRepo struct {
	client *resty.Client
}

func NewEnvTypeRepository() EnvTypeRepository {
	cfg := config.New()
	c := resty.New().
		SetBaseURL(cfg.BackendURL).
		SetDisableWarn(true).
		SetAuthToken(cfg.AccessToken)

	return &envTypeRepo{
		client: c,
	}
}

func (e *envTypeRepo) GetAll() ([]responses.EnvTypeResponse, error) {
	var response []responses.EnvTypeResponse

	resp, err := e.client.R().
		SetResult(&response).
		Get("/env_type")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return response, nil
}
