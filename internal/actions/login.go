package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/browser"
	"github.com/savioxavier/termlink"
	"github.com/urfave/cli/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/style"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

func LoginAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		authService := services.NewAuthService()

		// Step 1: Initiate the login process
		credentials, err := authService.InitiateLogin()
		if err != nil {
			return fmt.Errorf("failed to initiate login: %w", err)
		}

		// Step 2: Display login information to user
		if err := displayLoginInstructions(credentials, cmd); err != nil {
			return err
		}

		// Step 3: Open browser for user convenience
		if err := openBrowserForLogin(credentials.VerificationUri); err != nil {
			// Don't fail if browser can't be opened, just warn
			fmt.Printf("⚠️  Could not open browser automatically: %v\n", err)
		}

		// Step 4: Poll for authentication completion
		fmt.Println("⏳ Waiting for authentication...")
		token, err := authService.PollForToken(credentials)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		// Step 5: Save the token
		if err := authService.SaveToken(token); err != nil {
			return fmt.Errorf("failed to save authentication token: %w", err)
		}

		// Step 6: Display success message
		fmt.Println("✅ Login successful! You are now authenticated.")

		return nil
	}
}

// displayLoginInstructions shows the user what they need to do to authenticate
func displayLoginInstructions(credentials interface{}, cmd *cli.Command) error {
	// Print as JSON if requested
	if cmd.Bool("json") {
		return printAsJSON(credentials)
	}

	// Type assert to our domain model
	creds, ok := credentials.(*domain.LoginCredentials)
	if !ok {
		return fmt.Errorf("invalid credentials type")
	}

	fmt.Println("🔐 Authentication Required")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println(termlink.Link("1. Open this URL in your browser: ", style.LinkStyle.Render(creds.GetVerificationUri()), true))
	fmt.Printf("2. Enter this verification code: %s\n", creds.GetUserCode())
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	return nil
}

// openBrowserForLogin attempts to open the verification URL in the user's default browser
func openBrowserForLogin(verificationUri string) error {
	return browser.OpenURL(verificationUri)
}

// printAsJSON prints the credentials in JSON format for debugging or scripting purposes
func printAsJSON(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	fmt.Printf("Login Response JSON:\n%s\n", string(jsonData))
	return nil
}
