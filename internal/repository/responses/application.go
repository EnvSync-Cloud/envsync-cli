package responses

import "time"

type AppResponse struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Metadata    map[string]any `json:"metadata"`
	OrgID       string         `json:"org_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// NewAppResponse creates a new AppResponse instance
func NewAppResponse(id, name, description, orgID string, metadata map[string]any, createdAt, updatedAt time.Time) *AppResponse {
	return &AppResponse{
		ID:          id,
		Name:        name,
		Description: description,
		Metadata:    metadata,
		OrgID:       orgID,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// EnvTypeResponse represents the response structure for environment types
type EnvTypeResponse struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
