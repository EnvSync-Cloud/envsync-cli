package auth

import (
	"context"

	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

type whoamiUseCase struct {
	authService services.AuthService
}

func NewWhoamiUseCase(authService services.AuthService) WhoamiUseCase {
	return &whoamiUseCase{
		authService: authService,
	}
}

func (uc *whoamiUseCase) Execute(ctx context.Context, req WhoamiRequest) (*WhoamiResponse, error) {
	// Initialize response with configuration info
	response := &WhoamiResponse{
		IsLoggedIn: false,
		TokenValid: false,
	}

	// Try to get user information
	userInfo, err := uc.authService.Whoami()
	if err != nil {
		// Handle different types of errors
		response.TokenValid = false
		response.IsLoggedIn = false

		// Don't return error here - just indicate the user is not logged in
		// The error might be due to network issues, expired token, etc.
		return response, nil
	}

	// If we successfully got user info, the user is logged in
	response.UserInfo = userInfo
	response.IsLoggedIn = true
	response.TokenValid = true

	return response, nil
}
