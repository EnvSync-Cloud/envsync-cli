package repository

import (
	"fmt"

	"resty.dev/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/config"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/responses"
)

type EnvTypeRepository interface {
	GetAll() ([]responses.EnvTypeResponse, error)
	GetByID(id string) (responses.EnvTypeResponse, error)
	GetByAppID(appID string) ([]responses.EnvTypeResponse, error)
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

func (e *envTypeRepo) GetByID(id string) (responses.EnvTypeResponse, error) {
	var response responses.EnvTypeResponse

	resp, err := e.client.R().
		SetResult(&response).
		Get(fmt.Sprintf("/env_type/%s", id))

	if err != nil {
		return responses.EnvTypeResponse{}, err
	}

	if resp.StatusCode() != 200 {
		return responses.EnvTypeResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return response, nil
}

func (e *envTypeRepo) GetByAppID(appID string) ([]responses.EnvTypeResponse, error) {
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

	// Filter the response by appID
	var filteredResponse []responses.EnvTypeResponse
	for _, envType := range response {
		if envType.AppID == appID {
			filteredResponse = append(filteredResponse, envType)
		}
	}

	return filteredResponse, nil
}
