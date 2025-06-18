package services

//type EnvTypeService interface {
//	GetAllEnvTypes() ([]domain.EnvType, error)
//}
//
//type envTypeService struct {
//	repo repository.EnvTypeRepository
//}
//
//func NewEnvTypeService() EnvTypeService {
//	r := repository.NewEnvTypeRepository()
//
//	return &envTypeService{
//		repo: r,
//	}
//}
//
//func (e *envTypeService) GetAllEnvTypes() ([]domain.EnvType, error) {
//	res, err := e.repo.GetAll()
//	if err != nil {
//		return nil, err
//	}
//
//	var envTypes []domain.EnvType
//	for _, envTypeResp := range res {
//		envTypes = append(envTypes, mappers.EnvTypeResponseToDomain(envTypeResp))
//	}
//
//	return envTypes, nil
//}
