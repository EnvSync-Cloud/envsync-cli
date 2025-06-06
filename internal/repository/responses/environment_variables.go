package responses

type EnvironmentVariables struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	AppID     string `json:"app_id"`
	EnvTypeID string `json:"env_type_id"`
	OrgID     string `json:"org_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type EnvVariableList []struct {
	EnvironmentVariables
}

func (e *EnvVariableList) ToMap() map[string]string {
	result := make(map[string]string)
	for _, envVar := range *e {
		result[envVar.Key] = envVar.Value
	}
	return result
}
