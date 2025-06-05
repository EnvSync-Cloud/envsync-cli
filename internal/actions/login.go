package actions

import (
	"fmt"
	"time"

	"github.com/EnvSync-Cloud/envsync-cli/internal/config"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

func LoginAction() cli.ActionFunc {
	return func(c *cli.Context) error {
		s := services.NewAuthService()

		res, err := s.LoginDeviceCode()
		if err != nil {
			return err
		}

		fmt.Printf("Open this URL in your browser to login: %s\n", res.VerificationUri)
		fmt.Printf("\nVerify your code: %s\n", res.UserCode)

		if err = browser.OpenURL(res.VerificationUri); err != nil {
			return err
		}

		for {
			token, err := s.LoginToken(res.DeviceCode, res.ClientId, res.AuthDomain)

			if err == nil {
				fmt.Printf("✅ Login successful")

				cfg := config.New()
				cfg.AccessToken = token
				if err := cfg.WriteConfigFile(); err != nil {
					return err
				}

				break
			}

			time.Sleep(time.Duration(res.Interval) * time.Second)
		}

		return nil
	}
}
