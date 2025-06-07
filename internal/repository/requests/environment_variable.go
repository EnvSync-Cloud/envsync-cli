package requests

type EnvVariableRequest struct {
	AppID     string `json:"app_id"`
	EnvTypeID string `json:"env_type_id"`
}

type BatchSyncEnvRequest struct {
	AppID     string        `json:"app_id"`
	EnvTypeID string        `json:"env_type_id"`
	Envs      []EnvVariable `json:"envs"`
}

type EnvVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
