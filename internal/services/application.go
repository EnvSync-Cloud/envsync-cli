package services

import (
	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/mappers"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository"
)

type ApplicationService interface {
	CreateApp(app *domain.Application) (domain.Application, error)
	GetAppByID(id string) (domain.Application, error)
	GetAllApps() ([]domain.Application, error)
	DeleteApp(app domain.Application) error
	ReadAppEnvTypes() ([]domain.EnvironmentType, error)
}

type app struct {
	appRepo     repository.ApplicationRepository
	envTypeRepo repository.EnvTypeRepository
}

func NewAppService() ApplicationService {
	appRepo := repository.NewApplicationRepository()
	envTypeRepo := repository.NewEnvTypeRepository()

	return &app{
		appRepo:     appRepo,
		envTypeRepo: envTypeRepo,
	}
}

func (a *app) CreateApp(app *domain.Application) (domain.Application, error) {
	req := mappers.DomainToAppRequest(app)

	var appRes domain.Application
	if res, err := a.appRepo.Create(req); err != nil {
		return domain.Application{}, err
	} else {
		appRes = mappers.AppResponseToDomain(res)
	}

	return appRes, nil
}

func (a *app) GetAllApps() ([]domain.Application, error) {
	res, err := a.appRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var apps []domain.Application
	for _, appResp := range res {
		apps = append(apps, mappers.AppResponseToDomain(appResp))
	}

	return apps, nil
}

func (a *app) DeleteApp(app domain.Application) error {
	if err := a.appRepo.Delete(app.ID); err != nil {
		return err
	}

	return nil
}

func (a *app) GetAppByID(id string) (domain.Application, error) {
	res, err := a.appRepo.GetByID(id)
	if err != nil {
		return domain.Application{}, err
	}

	app := mappers.AppResponseToDomain(res)
	return app, nil
}

func (a *app) ReadAppEnvTypes() ([]domain.EnvironmentType, error) {
	res, err := a.envTypeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var envTypes []domain.EnvironmentType
	for _, envTypeResp := range res {
		envTypes = append(envTypes, mappers.EnvTypeResponseToDomain(envTypeResp))
	}

	return envTypes, nil
}
