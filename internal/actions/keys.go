package actions

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

// GenerateKeyPairs generates RSA key pairs with optimized settings
func GenerateKeyPairs() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		outDir := cmd.String("outDir")

		// Initialize key generation service
		kgs := services.NewKeyGenerationService()

		// Generate RSA key pair
		privateKey, publicKey, err := kgs.GenerateRSAKeyPair()
		if err != nil {
			return err
		}

		// Encode to PEM format
		privateKeyPEM, publicKeyPEM, err := kgs.EncodeToPEM(privateKey, publicKey)
		if err != nil {
			return err
		}

		// Save to files
		privateKeyPath, publicKeyPath, err := kgs.SaveKeyPairToFiles(privateKeyPEM, publicKeyPEM, outDir)
		if err != nil {
			return err
		}

		cmd.Writer.Write([]byte("✅ RSA key pair generated successfully!\n"))
		cmd.Writer.Write([]byte(fmt.Sprintf("🔐 Private key: %s\n", privateKeyPath)))
		cmd.Writer.Write([]byte(fmt.Sprintf("🔑 Public key: %s\n", publicKeyPath)))

		return nil
	}
}
