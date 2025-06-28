package app

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/app"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/cli/formatters"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/tui/factory"
	"github.com/EnvSync-Cloud/envsync-cli/internal/shared/utils"
)

type Handler struct {
	createUseCase app.CreateAppUseCase
	deleteUseCase app.DeleteAppUseCase
	listUseCase   app.ListAppsUseCase
	getUseCase    app.GetAppUseCase
	formatter     *formatters.AppFormatter
	tuiFactory    *factory.AppFactory
}

func NewHandler(
	createUseCase app.CreateAppUseCase,
	deleteUseCase app.DeleteAppUseCase,
	listUseCase app.ListAppsUseCase,
	getUseCase app.GetAppUseCase,
	formatter *formatters.AppFormatter,
	tuiFactory *factory.AppFactory,
) *Handler {
	return &Handler{
		createUseCase: createUseCase,
		deleteUseCase: deleteUseCase,
		listUseCase:   listUseCase,
		getUseCase:    getUseCase,
		formatter:     formatter,
		tuiFactory:    tuiFactory,
	}
}

func (h *Handler) Create(ctx context.Context, cmd *cli.Command) error {
	// Check if interactive mode should be used
	useInteractive := utils.IsInteractiveMode(cmd.Bool("json"))

	var application domain.Application
	if cmd.IsSet("name") {
		application.Name = cmd.String("name")
	}
	if cmd.IsSet("description") {
		application.Description = cmd.String("description")
	}
	if cmd.IsSet("metadata") {
		metadata := cmd.String("metadata")
		if metadata != "" {
			// Parse metadata JSON string into a map
			// Assuming metadata is in format key-value pairs like "key1=value1,key2=value2"
			metadataMap := make(map[string]any)
			// pairs := utils.ParseKeyValuePairs(metadata)
			// for _, pair := range pairs {
			// 	parts := utils.SplitKeyValue(pair)
			// 	if len(parts) == 2 {
			// 		metadataMap[parts[0]] = parts[1]
			// 	} else {
			// 		return h.formatter.FormatError(cmd.Writer, "Invalid metadata format. Use key=value pairs.")
			// 	}
			// }
			application.Metadata = metadataMap
		}
	}

	if useInteractive {
		return h.createInteractive(ctx)
	}
	return h.createCLI(ctx, cmd)
}

func (h *Handler) Delete(ctx context.Context, cmd *cli.Command) error {
	appID := cmd.String("id")
	appName := cmd.String("name")

	// If no ID or name provided and we're in a terminal, use interactive mode
	if appID == "" && appName == "" && utils.IsTerminal() {
		return h.deleteInteractive(ctx)
	}

	return h.deleteCLI(ctx, cmd, appID, appName)
}

func (h *Handler) List(ctx context.Context, cmd *cli.Command) error {
	useInteractive := utils.IsInteractiveMode(cmd.Bool("json"))

	if useInteractive {
		return h.listInteractive(ctx)
	}
	return h.listCLI(ctx, cmd)
}

func (h *Handler) Select(ctx context.Context, cmd *cli.Command) error {
	// Select is always interactive
	return h.selectInteractive(ctx, cmd)
}

// Interactive implementations

func (h *Handler) createInteractive(ctx context.Context) error {
	return h.tuiFactory.CreateAppInteractive(ctx)
}

func (h *Handler) deleteInteractive(ctx context.Context) error {
	// return nil
	return h.tuiFactory.DeleteAppInteractive(ctx)
}

func (h *Handler) listInteractive(ctx context.Context) error {
	// return nil
	return h.tuiFactory.ListAppsInteractive(ctx)
}

func (h *Handler) selectInteractive(ctx context.Context, cmd *cli.Command) error {
	// return nil
	return h.tuiFactory.SelectAppInteractive(ctx)
}

// CLI implementations

func (h *Handler) createCLI(ctx context.Context, cmd *cli.Command) error {
	// For CLI mode, we could prompt for input or require flags
	// For now, let's return an error asking for interactive mode
	return h.formatter.FormatError(cmd.Writer, "CLI mode for app creation not implemented. Use interactive mode or add required flags.")
}

func (h *Handler) deleteCLI(ctx context.Context, cmd *cli.Command, appID, appName string) error {
	// Build delete request
	req := app.DeleteAppRequest{
		ID:   appID,
		Name: appName,
	}

	// Execute use case
	if err := h.deleteUseCase.Execute(ctx, req); err != nil {
		return h.formatUseCaseError(cmd, err)
	}

	// Format success message
	identifier := appID
	if identifier == "" {
		identifier = appName
	}
	return h.formatter.FormatSuccess(cmd.Writer, "Application '"+identifier+"' deleted successfully!")
}

func (h *Handler) listCLI(ctx context.Context, cmd *cli.Command) error {
	// Build list request
	req := app.ListAppsRequest{
		Limit:  cmd.Int("limit"),
		Offset: cmd.Int("offset"),
	}

	// Execute use case
	apps, err := h.listUseCase.Execute(ctx, req)
	if err != nil {
		return h.formatUseCaseError(cmd, err)
	}

	// Format output based on requested format
	if cmd.Bool("json") {
		return h.formatter.FormatJSON(cmd.Writer, apps)
	}

	format := cmd.String("format")
	switch format {
	case "compact":
		return h.formatter.FormatCompact(cmd.Writer, apps)
	case "list":
		return h.formatter.FormatList(cmd.Writer, apps)
	default:
		return h.formatter.FormatTable(cmd.Writer, apps)
	}
}

// Helper methods

func (h *Handler) formatUseCaseError(cmd *cli.Command, err error) error {
	// Handle different types of use case errors
	switch e := err.(type) {
	case *app.AppError:
		switch e.Code {
		case app.AppErrorCodeNotFound:
			return h.formatter.FormatError(cmd.Writer, "Application not found: "+e.Message)
		case app.AppErrorCodeAlreadyExists:
			return h.formatter.FormatError(cmd.Writer, "Application already exists: "+e.Message)
		case app.AppErrorCodeValidation:
			return h.formatter.FormatError(cmd.Writer, "Validation error: "+e.Message)
		case app.AppErrorCodeAccessDenied:
			return h.formatter.FormatError(cmd.Writer, "Access denied: "+e.Message)
		case app.AppErrorCodeInUse:
			return h.formatter.FormatWarning(cmd.Writer, "Cannot complete operation: "+e.Message)
		default:
			return h.formatter.FormatError(cmd.Writer, "Service error: "+e.Message)
		}
	default:
		return h.formatter.FormatError(cmd.Writer, "Unexpected error: "+err.Error())
	}
}
