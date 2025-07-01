package handlers

import (
	"context"

	inituc "github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/init"
	"github.com/urfave/cli/v3"
)

type InitHandler struct {
	initUseCase inituc.InitUseCase
}

func NewInitHandler(initUseCase inituc.InitUseCase) *InitHandler {
	return &InitHandler{
		initUseCase: initUseCase,
	}
}

func (h *InitHandler) Init(ctx context.Context, cmd *cli.Command) error {
	err := h.initUseCase.Execute(ctx, cmd.String("config"))
	if err != nil {
		return err
	}
	return nil
}
