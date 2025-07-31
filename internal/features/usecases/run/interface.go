package run

import "context"

type InjectEnvUseCase interface {
	Execute(context.Context) (map[string]string, error)
}

type InjectSecretsUseCase interface {
	Execute(context.Context) (map[string]string, error)
}

type RedactUseCase interface {
	Execute(context.Context, []string, []string) int
}
