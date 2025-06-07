package mappers

import (
	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/responses"
)

func EnvTypeResponseToDomain(envTypeResp responses.EnvTypeResponse) domain.EnvType {
	return domain.EnvType{
		ID:   envTypeResp.ID,
		Name: envTypeResp.Name,
	}
}
