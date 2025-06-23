package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
)

// EncryptionService handles encryption and decryption operations
type EncryptionService struct{}

// NewEncryptionService creates a new encryption service
func NewEncryptionService() *EncryptionService {
	return &EncryptionService{}
}

// hybridEncrypt encrypts data using AES-192-GCM + RSA hybrid encryption
func (es *EncryptionService) hybridEncrypt(data string, publicKeyPEM string) (string, error) {
	// Parse public key
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", errors.New("failed to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("not an RSA public key")
	}

	// Generate AES-192 key (24 bytes)
	aesKey := make([]byte, 24)
	if _, err := rand.Read(aesKey); err != nil {
		return "", fmt.Errorf("failed to generate AES key: %w", err)
	}

	// Generate 12-byte IV for GCM mode
	iv := make([]byte, 12)
	if _, err := rand.Read(iv); err != nil {
		return "", fmt.Errorf("failed to generate IV: %w", err)
	}

	// Create AES-192-GCM cipher
	block2, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block2)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Encrypt data
	ciphertext := gcm.Seal(nil, iv, []byte(data), nil)

	// Encrypt AES key with RSA using OAEP padding
	encryptedAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPublicKey, aesKey, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt AES key: %w", err)
	}

	// Pack everything: keyLength(2 bytes) + encryptedKey + iv(12) + ciphertext
	keyLength := make([]byte, 2)
	binary.BigEndian.PutUint16(keyLength, uint16(len(encryptedAESKey)))

	result := make([]byte, 0, 2+len(encryptedAESKey)+12+len(ciphertext))
	result = append(result, keyLength...)
	result = append(result, encryptedAESKey...)
	result = append(result, iv...)
	result = append(result, ciphertext...)

	return base64.StdEncoding.EncodeToString(result), nil
}

// hybridDecrypt decrypts data using AES-192-GCM + RSA hybrid decryption
func (es *EncryptionService) hybridDecrypt(encryptedData string, privateKeyPEM string) (string, error) {
	// Parse private key
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("not an RSA private key")
	}

	// Decode base64 data
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	if len(data) < 2 {
		return "", errors.New("invalid encrypted data format")
	}

	// Extract components
	keyLength := binary.BigEndian.Uint16(data[0:2])
	offset := 2

	if len(data) < int(2+keyLength+12) {
		return "", errors.New("invalid encrypted data format")
	}

	encryptedAESKey := data[offset : offset+int(keyLength)]
	offset += int(keyLength)

	iv := data[offset : offset+12]
	offset += 12

	ciphertext := data[offset:]

	// Decrypt AES key using RSA private key with OAEP padding
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivateKey, encryptedAESKey, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt AES key: %w", err)
	}

	// Decrypt data using AES-192-GCM
	block2, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block2)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data: %w", err)
	}

	return string(plaintext), nil
}

// rsaEncryptSmall encrypts small data directly with RSA
func (es *EncryptionService) rsaEncryptSmall(data string, publicKeyPEM string) (string, error) {
	dataBytes := []byte(data)

	// Check size limit for direct RSA encryption with OAEP padding
	if len(dataBytes) > 190 { // 3072/8 - 2*32 - 2 (OAEP overhead) - margin
		return "", errors.New("data too large for direct RSA encryption. Use hybrid encryption instead")
	}

	// Parse public key
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", errors.New("failed to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("not an RSA public key")
	}

	// Encrypt with RSA OAEP
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPublicKey, dataBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// rsaDecryptSmall decrypts small data directly with RSA
func (es *EncryptionService) rsaDecryptSmall(encryptedData string, privateKeyPEM string) (string, error) {
	// Parse private key
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("not an RSA private key")
	}

	// Decode base64
	encrypted, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Decrypt with RSA OAEP
	decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivateKey, encrypted, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data: %w", err)
	}

	return string(decrypted), nil
}

// SmartEncrypt chooses the best encryption method based on data size
func (es *EncryptionService) SmartEncrypt(data string, publicKeyPEM string) (string, error) {
	dataSize := len([]byte(data))

	// Use direct RSA for small data (faster, smaller output)
	if dataSize <= 150 { // Safe margin for OAEP
		encrypted, err := es.rsaEncryptSmall(data, publicKeyPEM)
		if err != nil {
			return "", err
		}
		return "RSA:" + encrypted, nil
	}

	// Use hybrid for larger data
	encrypted, err := es.hybridEncrypt(data, publicKeyPEM)
	if err != nil {
		return "", err
	}
	return "HYB:" + encrypted, nil
}

// SmartDecrypt decrypts data based on the encryption method prefix
func (es *EncryptionService) SmartDecrypt(encryptedData string, privateKeyPEM string) (string, error) {
	if len(encryptedData) < 4 {
		return "", errors.New("invalid encrypted data format")
	}

	method := encryptedData[:4]
	data := encryptedData[4:]

	switch method {
	case "RSA:":
		return es.rsaDecryptSmall(data, privateKeyPEM)
	case "HYB:":
		return es.hybridDecrypt(data, privateKeyPEM)
	default:
		return "", errors.New("unknown encryption method")
	}
}
