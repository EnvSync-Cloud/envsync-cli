package mappers

import (
	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/repository/responses"
)

func EnvironmentVariableToDomain(env responses.EnvironmentVariables) *domain.EnvironmentVariable {
	return &domain.EnvironmentVariable{
		Key:   env.Key,
		Value: env.Value,
	}
}

func EnvironmentVariablesToDomain(envs []responses.EnvironmentVariables) []*domain.EnvironmentVariable {
	var domainEnvs []*domain.EnvironmentVariable
	for _, env := range envs {
		domainEnvs = append(domainEnvs, EnvironmentVariableToDomain(env))
	}
	return domainEnvs
}
