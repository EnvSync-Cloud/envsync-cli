package actions

import (
	"context"

	"github.com/urfave/cli/v3"
)

func IndexAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		// TODO: Implement index action to desiaply default help message
		cmd.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		cmd.Writer.Write([]byte("Welcome to EnvSync CLI!\n"))
		cmd.Writer.Write([]byte("Use 'envsync --help' to see available commands.\n"))
		cmd.Writer.Write([]byte("For more information, visit: https://envsync.cloud/docs\n"))
		cmd.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))

		return nil
	}
}
