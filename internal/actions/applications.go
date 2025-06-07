package actions

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/urfave/cli/v2"
)

func CreateApplication() cli.ActionFunc {
	return func(c *cli.Context) error {
		scanner := bufio.NewScanner(os.Stdin)

		// Step 1: Get application name
		fmt.Print("📛 Enter application name: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading application name: %w", err)
			}
			return fmt.Errorf("failed to read application name")
		}
		name := strings.TrimSpace(scanner.Text())
		if name == "" {
			return fmt.Errorf("application name cannot be empty")
		}

		// Step 2: Get application description
		fmt.Print("📝 Enter application description: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading application description: %w", err)
			}
			return fmt.Errorf("failed to read application description")
		}
		description := strings.TrimSpace(scanner.Text())
		if description == "" {
			return fmt.Errorf("application description cannot be empty")
		}

		// Step 3: Get optional metadata
		metadata := make(map[string]any)
		fmt.Print("🏷️  Enter metadata (key=value,key2=value2 format, or press Enter to skip): ")
		if scanner.Scan() {
			metadataStr := strings.TrimSpace(scanner.Text())
			if metadataStr != "" {
				// Simple key=value parsing, separated by commas
				// Format: "key1=value1,key2=value2"
				pairs := strings.Split(metadataStr, ",")
				for _, pair := range pairs {
					kv := strings.SplitN(pair, "=", 2)
					if len(kv) == 2 {
						metadata[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
					}
				}
			}
		} else if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading metadata: %w", err)
		}

		// Step 4: Create domain object
		app := domain.NewApplication(name, description, metadata)

		// Step 5: Initialize application service
		as := services.NewAppService()

		// Step 6: Create application
		if err := as.CreateApp(*app); err != nil {
			return fmt.Errorf("failed to create application: %w", err)
		}

		// Step 7: Success message
		c.App.Writer.Write([]byte("✅ Application created successfully!\n"))
		c.App.Writer.Write([]byte(fmt.Sprintf("📛 Name: %s\n", name)))
		c.App.Writer.Write([]byte(fmt.Sprintf("📝 Description: %s\n", description)))
		if len(metadata) > 0 {
			c.App.Writer.Write([]byte("🏷️  Metadata:\n"))
			for key, value := range metadata {
				c.App.Writer.Write([]byte(fmt.Sprintf("   %s: %v\n", key, value)))
			}
		}

		return nil
	}
}

func ListApplications() cli.ActionFunc {
	return func(c *cli.Context) error {
		// Step1: Initialize application service
		as := services.NewAppService()

		// Step2: Get all applications
		apps, err := as.GetAllApps()
		if err != nil {
			return err
		}

		// Step3: Print applications
		c.App.Writer.Write([]byte("🚀 Available Applications:\n"))
		for _, app := range apps {
			c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
			c.App.Writer.Write([]byte(fmt.Sprintf("📛 Name: %s\n", app.Name)))
			c.App.Writer.Write([]byte(fmt.Sprintf("🆔 ID: %s\n", app.ID)))
			c.App.Writer.Write([]byte(fmt.Sprintf("📝 Description: %s\n", app.Description)))
			c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		}

		return nil
	}
}

func DeleteApplication() cli.ActionFunc {
	return func(c *cli.Context) error {
		scanner := bufio.NewScanner(os.Stdin)

		// Step 1: Initialize application service
		as := services.NewAppService()

		// Step 2: Get all applications
		apps, err := as.GetAllApps()
		if err != nil {
			return fmt.Errorf("failed to fetch applications: %w", err)
		}

		if len(apps) == 0 {
			c.App.Writer.Write([]byte("📭 No applications found.\n"))
			return nil
		}

		// Step 3: Display available applications
		c.App.Writer.Write([]byte("🚀 Available Applications:\n"))
		for i, app := range apps {
			c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
			c.App.Writer.Write([]byte(fmt.Sprintf("%d. 📛 Name: %s\n", i+1, app.Name)))
			c.App.Writer.Write([]byte(fmt.Sprintf("   🆔 ID: %s\n", app.ID)))
			c.App.Writer.Write([]byte(fmt.Sprintf("   📝 Description: %s\n", app.Description)))
		}
		c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))

		// Step 4: Get user selection
		fmt.Print("🎯 Enter application name or ID to delete: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading input: %w", err)
			}
			return fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			return fmt.Errorf("input cannot be empty")
		}

		// Step 5: Find the application to delete
		var selectedApp *domain.Application
		for _, app := range apps {
			if app.Name == input || app.ID == input {
				selectedApp = &app
				break
			}
		}

		if selectedApp == nil {
			return fmt.Errorf("application with name or ID '%s' not found", input)
		}

		// Step 6: Confirm deletion
		fmt.Printf("\n⚠️  Are you sure you want to delete application '%s'? This action cannot be undone.\n", selectedApp.Name)
		fmt.Print("Type 'yes' to confirm: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading confirmation: %w", err)
			}
			return fmt.Errorf("failed to read confirmation")
		}

		confirmation := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if confirmation != "yes" && confirmation != "y" {
			c.App.Writer.Write([]byte("❎ Deletion cancelled.\n"))
			return nil
		}

		// Step 7: Delete the application
		if err := as.DeleteApp(*selectedApp); err != nil {
			return fmt.Errorf("failed to delete application: %w", err)
		}

		// Step 8: Success message
		c.App.Writer.Write([]byte("✅ Application deleted successfully!\n"))
		c.App.Writer.Write([]byte(fmt.Sprintf("📛 Deleted: %s\n", selectedApp.Name)))

		return nil
	}
}
