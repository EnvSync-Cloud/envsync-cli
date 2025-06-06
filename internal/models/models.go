package models

type ProjectEnvConfig struct {
	AppID   string `toml:"app_id"`
	EnvType string `toml:"env_type"`
}
