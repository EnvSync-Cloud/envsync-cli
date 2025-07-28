package services

import (
	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/mappers"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository"
)

type EnvTypeService interface {
	CreateEnvType(envType *domain.EnvType) (domain.EnvType, error)
	GetEnvTypeByID(id string) (domain.EnvType, error)
	GetEnvTypeByAppID(appID string) ([]domain.EnvType, error)
}

type envTypeService struct {
	repo repository.EnvTypeRepository
}

func NewEnvTypeService() EnvTypeService {
	r := repository.NewEnvTypeRepository()

	return &envTypeService{
		repo: r,
	}
}

func (e *envTypeService) CreateEnvType(envType *domain.EnvType) (domain.EnvType, error) {
	req := mappers.EnvTypeDomainToRequest(envType)

	res, err := e.repo.Create(&req)
	if err != nil {
		return domain.EnvType{}, err
	}

	return mappers.EnvTypeResponseToDomain(res), nil
}

func (e *envTypeService) GetEnvTypeByID(id string) (domain.EnvType, error) {
	res, err := e.repo.GetByID(id)
	if err != nil {
		return domain.EnvType{}, err
	}

	return mappers.EnvTypeResponseToDomain(res), nil
}

func (e *envTypeService) GetEnvTypeByAppID(appID string) ([]domain.EnvType, error) {
	res, err := e.repo.GetByAppID(appID)
	if err != nil {
		return nil, err
	}

	var envTypes []domain.EnvType
	for _, envTypeResp := range res {
		envTypes = append(envTypes, mappers.EnvTypeResponseToDomain(envTypeResp))
	}

	return envTypes, nil
}
