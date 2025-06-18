package actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/config"
)

func SetConfigAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		args := cmd.Args()

		if args.Len() < 1 {
			return fmt.Errorf("❌ No argument provided. Usage: envsync config set key=value")
		}

		// Parse key=value pairs from all arguments
		// Read current config from file
		cfg, err := config.ReadConfigFile()
		if err != nil {
			// If config file doesn't exist, create a new one
			cfg = config.AppConfig{}
		}

		for i := 0; i < args.Len(); i++ {
			arg := args.Get(i)

			// Split on first '=' to handle values that might contain '='
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("❌ Invalid format: '%s'. Expected format: key=value", arg)
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if key == "" {
				return fmt.Errorf("❌ Empty key provided in: '%s'", arg)
			}

			// Set the configuration based on the key
			switch strings.ToLower(key) {
			case "access_token", "accesstoken":
				cfg.AccessToken = value
				fmt.Printf("🔑 Set access_token to: %s\n", value)
			case "backend_url", "backendurl":
				cfg.BackendURL = value
				fmt.Printf("🌐 Set backend_url to: %s\n", value)
			default:
				return fmt.Errorf("❌ Unknown configuration key: '%s'. Valid keys are: access_token, backend_url", key)
			}
		}

		// Write the updated configuration to file
		if err := (&cfg).WriteConfigFile(); err != nil {
			return fmt.Errorf("❌ Failed to write config file: %v", err)
		}

		fmt.Println("✅ Configuration updated successfully!")
		return nil
	}
}

func GetConfigAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		args := cmd.Args()

		// Read current configuration from file
		cfg, err := config.ReadConfigFile()
		if err != nil {
			return fmt.Errorf("❌ Failed to read config file: %v", err)
		}

		// If no arguments provided, show all config values
		if args.Len() == 0 {
			fmt.Println("📋 Current configuration:")
			fmt.Printf("🔑 access_token: %s\n", cfg.AccessToken)
			fmt.Printf("🌐 backend_url: %s\n", cfg.BackendURL)
			return nil
		}

		// Show specific config value(s)
		for i := 0; i < args.Len(); i++ {
			key := strings.ToLower(strings.TrimSpace(args.Get(i)))

			switch key {
			case "access_token", "accesstoken":
				fmt.Printf("🔑 access_token: %s\n", cfg.AccessToken)
			case "backend_url", "backendurl":
				fmt.Printf("🌐 backend_url: %s\n", cfg.BackendURL)
			default:
				return fmt.Errorf("❌ Unknown configuration key: '%s'. Valid keys are: access_token, backend_url", key)
			}
		}

		return nil
	}
}
