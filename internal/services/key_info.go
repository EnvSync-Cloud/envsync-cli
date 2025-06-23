package services

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// KeyInfo represents information about a key
type KeyInfo struct {
	FilePath   string
	Type       string
	Algorithm  string
	KeySize    int
	CanEncrypt bool
	CanDecrypt bool
}

// KeyInfoService handles key information operations
type KeyInfoService struct{}

// NewKeyInfoService creates a new key info service
func NewKeyInfoService() *KeyInfoService {
	return &KeyInfoService{}
}

// GetKeyInfo extracts information from a key file
func (kis *KeyInfoService) GetKeyInfo(keyPath string) (*KeyInfo, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	keyInfo := &KeyInfo{
		FilePath: keyPath,
		Type:     block.Type,
	}

	if block.Type == "PRIVATE KEY" {
		privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}

		if rsaKey, ok := privateKey.(*rsa.PrivateKey); ok {
			keyInfo.Algorithm = "RSA"
			keyInfo.KeySize = rsaKey.Size() * 8
			keyInfo.CanEncrypt = true
			keyInfo.CanDecrypt = true
		}
	} else if block.Type == "PUBLIC KEY" {
		publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %w", err)
		}

		if rsaKey, ok := publicKey.(*rsa.PublicKey); ok {
			keyInfo.Algorithm = "RSA"
			keyInfo.KeySize = rsaKey.Size() * 8
			keyInfo.CanEncrypt = true
			keyInfo.CanDecrypt = false
		}
	}

	return keyInfo, nil
}
