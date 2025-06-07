package repository

import (
	"fmt"

	"github.com/EnvSync-Cloud/envsync-cli/internal/config"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/requests"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/responses"
	"resty.dev/v3"
)

type ApplicationRepository interface {
	Create(app requests.ApplicationRequest) error
	GetAll() ([]responses.AppResponse, error)
	Delete(id string) error
}

type appRepo struct {
	client *resty.Client
}

func NewApplicationRepository() ApplicationRepository {
	cfg := config.New()
	c := resty.New().
		SetBaseURL(cfg.BackendURL).
		SetAuthToken(cfg.AccessToken)

	return &appRepo{
		client: c,
	}
}

func (a *appRepo) Create(app requests.ApplicationRequest) error {
	var response responses.AppResponse

	resp, err := a.client.R().
		SetBody(app).
		SetResult(&response).
		Post("/app")

	if err != nil {
		return err
	}

	if resp.StatusCode() != 201 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}

func (a *appRepo) GetAll() ([]responses.AppResponse, error) {
	var apps []responses.AppResponse

	resp, err := a.client.R().
		SetResult(&apps).
		Get("/app")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return apps, nil
}

func (a *appRepo) Delete(id string) error {
	resp, err := a.client.R().
		SetPathParam("id", id).
		Delete("/app/{id}")

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}
