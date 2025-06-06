package requests

type EnvVariableRequest struct {
	AppID   string `json:"app_id"`
	EnvType string `json:"env_type_id"`
}
