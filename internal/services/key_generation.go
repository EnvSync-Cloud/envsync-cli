package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

// KeyGenerationService handles RSA key pair generation
type KeyGenerationService struct{}

// NewKeyGenerationService creates a new key generation service
func NewKeyGenerationService() *KeyGenerationService {
	return &KeyGenerationService{}
}

// GenerateRSAKeyPair generates a 3072-bit RSA key pair
func (kgs *KeyGenerationService) GenerateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// Generate RSA key pair with 3072-bit modulus
	privateKey, err := rsa.GenerateKey(rand.Reader, 3072)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	return privateKey, &privateKey.PublicKey, nil
}

// EncodeToPEM encodes keys to PEM format
func (kgs *KeyGenerationService) EncodeToPEM(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) ([]byte, []byte, error) {
	// Encode private key to PKCS#8 PEM format
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Encode public key to SPKI PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return privateKeyPEM, publicKeyPEM, nil
}

// SaveKeyPairToFiles saves the key pair to files
func (kgs *KeyGenerationService) SaveKeyPairToFiles(privateKeyPEM, publicKeyPEM []byte, outDir string) (string, string, error) {
	// Ensure output directory exists, creating it if necessary
	outDir, err := filepath.Abs(outDir)
	if err != nil {
		return "", "", fmt.Errorf("failed to get absolute path for output directory: %w", err)
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Define file paths
	privateKeyPath := filepath.Join(outDir, "private_key.pem")
	publicKeyPath := filepath.Join(outDir, "public_key.pem")

	// Write private key
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		return "", "", fmt.Errorf("failed to write private key: %w", err)
	}

	// Write public key
	if err := os.WriteFile(publicKeyPath, publicKeyPEM, 0644); err != nil {
		return "", "", fmt.Errorf("failed to write public key: %w", err)
	}

	return privateKeyPath, publicKeyPath, nil
}
