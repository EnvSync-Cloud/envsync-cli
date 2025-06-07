package mappers

import (
	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/requests"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/responses"
)

func AppResponseToDomain(appResp responses.AppResponse) domain.Application {
	return domain.Application{
		ID:          appResp.ID,
		Name:        appResp.Name,
		Description: appResp.Description,
		Metadata:    appResp.Metadata,
		OrgID:       appResp.OrgID,
		CreatedAt:   appResp.CreatedAt,
		UpdatedAt:   appResp.UpdatedAt,
	}
}

func DomainToAppRequest(app *domain.Application) requests.ApplicationRequest {
	return requests.ApplicationRequest{
		Name:        app.Name,
		Description: app.Description,
		Metadata:    app.Metadata,
	}
}
