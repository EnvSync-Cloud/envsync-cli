package handlers

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/app"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/formatters"
)

type AppHandler struct {
	createUseCase app.CreateAppUseCase
	deleteUseCase app.DeleteAppUseCase
	listUseCase   app.ListAppsUseCase
	formatter     *formatters.AppFormatter
}

func NewAppHandler(
	createUseCase app.CreateAppUseCase,
	deleteUseCase app.DeleteAppUseCase,
	listUseCase app.ListAppsUseCase,
	formatter *formatters.AppFormatter,
) *AppHandler {
	return &AppHandler{
		createUseCase: createUseCase,
		deleteUseCase: deleteUseCase,
		listUseCase:   listUseCase,
		formatter:     formatter,
	}
}

func (h *AppHandler) Create(ctx context.Context, cmd *cli.Command) error {
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

	app, err := h.createUseCase.Execute(ctx, application)
	if err != nil {
		return h.formatUseCaseError(cmd, err)
	}

	if cmd.Bool("json") {
		// If JSON output is requested, format the application as JSON
		return h.formatter.BaseFormatter.FormatJSON(cmd.Writer, app)
	}

	// Display success message
	return h.formatter.FormatCreateSuccessMessage(cmd.Writer, *app)
}

func (h *AppHandler) Delete(ctx context.Context, cmd *cli.Command) error {
	_ = h.deleteUseCase.Execute(ctx)

	h.formatter.FormatSuccess(cmd.Writer, "Successfully deleted application!")

	return nil
}

func (h *AppHandler) List(ctx context.Context, cmd *cli.Command) error {
	_ = h.listUseCase.Execute(ctx)

	return nil
}

// Helper methods

func (h *AppHandler) formatUseCaseError(cmd *cli.Command, err error) error {
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
