package handlers

import (
	"context"
	"fmt"

	"github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/sync"
	"github.com/urfave/cli/v3"
)

type SyncHandler struct {
	pullUseCase sync.PullUseCase
	pushUseCase sync.PushUseCase
}

func NewSyncHandler(
	pullUseCase sync.PullUseCase,
	pushUseCase sync.PushUseCase,
) *SyncHandler {
	return &SyncHandler{
		pullUseCase: pullUseCase,
		pushUseCase: pushUseCase,
	}
}

func (h *SyncHandler) Pull(ctx context.Context, cmd *cli.Command) error {
	config := cmd.String("config")

	diff, err := h.pullUseCase.Execute(ctx, config)
	if err != nil {
		return err
	}

	if len(diff.Warnings) > 0 {
		// Handle warnings, e.g., print or log them
		for _, warning := range diff.Warnings {
			fmt.Printf("Warning: %s\n", warning)
		}
		return nil
	}

	if len(diff.Added) > 0 || len(diff.Updated) > 0 || len(diff.Deleted) > 0 {
		// Handle the sync response, e.g., print or log the changes
		fmt.Printf("Sync completed with %d added, %d updated, and %d deleted variables.\n",
			len(diff.Added), len(diff.Updated), len(diff.Deleted))
	} else {
		fmt.Printf("No changes detected during sync.\n")
	}

	return nil
}

func (h *SyncHandler) Push(ctx context.Context, cmd *cli.Command) error {
	config := cmd.String("config")

	diff, err := h.pushUseCase.Execute(ctx, config)
	if err != nil {
		return err
	}

	if len(diff.Warnings) > 0 {
		// Handle warnings, e.g., print or log them
		for _, warning := range diff.Warnings {
			fmt.Printf("Warning: %s\n", warning)
		}
		return nil
	}

	if len(diff.Added) > 0 || len(diff.Updated) > 0 || len(diff.Deleted) > 0 {
		// Handle the sync response, e.g., print or log the changes
		fmt.Printf("Sync completed with %d added, %d updated, and %d deleted variables.\n",
			len(diff.Added), len(diff.Updated), len(diff.Deleted))
	} else {
		fmt.Printf("No changes detected during sync.\n")
	}

	return nil
}
