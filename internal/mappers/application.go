package mappers

import (
	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/requests"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/responses"
)

func AppResponseToDomain(res responses.AppResponse) domain.Application {
	return domain.Application{
		ID:          res.ID,
		Name:        res.Name,
		Description: res.Description,
		Metadata:    res.Metadata,
		OrgID:       res.OrgID,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}
}

func DomainToAppRequest(app *domain.Application) requests.ApplicationRequest {
	return requests.ApplicationRequest{
		Name:        app.Name,
		Description: app.Description,
		Metadata:    app.Metadata,
	}
}
