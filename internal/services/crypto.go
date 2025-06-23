package services

// SmartEncrypt is a convenience function that uses the encryption service
func SmartEncrypt(data string, publicKeyPEM string) (string, error) {
	es := NewEncryptionService()
	return es.SmartEncrypt(data, publicKeyPEM)
}

// SmartDecrypt is a convenience function that uses the encryption service
func SmartDecrypt(encryptedData string, privateKeyPEM string) (string, error) {
	es := NewEncryptionService()
	return es.SmartDecrypt(encryptedData, privateKeyPEM)
}
