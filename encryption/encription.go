package encryption

type encryptionHandler interface {
	Encrypt(secretKey []byte, plaintext []byte) ([]byte, error)
	Decrypt(secretKey []byte, ciphertext []byte) ([]byte, error)
}

func CreateencryptionHandler() encryptionHandler {
	return createAes256GcmHandler()
}
