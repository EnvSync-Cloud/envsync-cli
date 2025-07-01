package run

import "context"

type RedactUseCase interface {
	Execute(context.Context, []string, []string) int
}

type InjectEnvUseCase interface {
	Execute(context.Context) (map[string]string, error)
}
