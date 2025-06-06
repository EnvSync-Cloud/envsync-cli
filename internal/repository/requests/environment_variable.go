package requests

type EnvVariableRequest struct {
	AppID     string `json:"app_id"`
	EnvTypeID string `json:"env_type_id"`
}
