package actions

import (
	"encoding/json"
	"fmt"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

func LoginAction() cli.ActionFunc {
	return func(c *cli.Context) error {
		authService := services.NewAuthService()

		// Step 1: Initiate the login process
		credentials, err := authService.InitiateLogin()
		if err != nil {
			return fmt.Errorf("failed to initiate login: %w", err)
		}

		// Step 2: Display login information to user
		if err := displayLoginInstructions(credentials, c); err != nil {
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
func displayLoginInstructions(credentials interface{}, c *cli.Context) error {
	// Print as JSON if requested
	if c.Bool("json") {
		return printAsJSON(credentials)
	}

	// Type assert to our domain model
	creds, ok := credentials.(*domain.LoginCredentials)
	if !ok {
		return fmt.Errorf("invalid credentials type")
	}

	fmt.Println("🔐 Authentication Required")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("1. Open this URL in your browser: %s\n", creds.GetVerificationUri())
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
