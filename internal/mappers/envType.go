package mappers

import (
	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/responses"
)

func EnvTypeResponseToDomain(res responses.EnvTypeResponse) domain.EnvType {
	return domain.EnvType{
		ID:          res.ID,
		OrgID:       res.OrgID,
		AppID:       res.AppID,
		IsDefault:   res.IsDefault,
		IsProtected: res.IsProtected,
		Color:       res.Color,
		Name:        res.Name,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}
}
