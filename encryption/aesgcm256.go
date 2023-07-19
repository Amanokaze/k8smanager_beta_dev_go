package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

type Aes256GcmHandler struct {
}

func getSecretKeytoHash(secretKey []byte) ([]byte, error) {
	// gen 32 byte secret key
	hash := sha256.New()
	_, err := hash.Write(secretKey)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

func checkSecretKey(secretKey []byte) error {
	if len(secretKey) != 32 {
		return fmt.Errorf("secret key is not for AES-256: total %d bits", 8*len(secretKey))
	} else {
		return nil
	}
}

// func AES256GSMEncrypt(secretKey []byte, plaintext []byte) ([]byte, error) {
func (a *Aes256GcmHandler) Encrypt(secretKey []byte, plaintext []byte) ([]byte, error) {
	secretKey, err := getSecretKeytoHash(secretKey)
	if err != nil {
		return nil, err
	}

	err = checkSecretKey(secretKey)
	if err != nil {
		return nil, err
	}

	// prepare AES-256-GSM cipher
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// make random nonce
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// encrypt plaintext
	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)
	//fmt.Println("--\nencryption start")
	//fmt.Printf("nonce to use: %x\n", nonce)
	//fmt.Printf("ciphertext: %x\n", ciphertext)

	return ciphertext, nil // nonce is included in ciphertext. no need to return
}

// func AES256GSMDecrypt(secretKey []byte, ciphertext []byte) ([]byte, error) {
func (a *Aes256GcmHandler) Decrypt(secretKey []byte, ciphertext []byte) ([]byte, error) {
	secretKey, err := getSecretKeytoHash(secretKey)
	if err != nil {
		return nil, err
	}

	err = checkSecretKey(secretKey)
	if err != nil {
		return nil, err
	}

	// prepare AES-256-GSM cipher
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	nonce, pureCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// decrypt ciphertext
	plaintext, err := aesgcm.Open(nil, nonce, pureCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func createAes256GcmHandler() encryptionHandler {
	return &Aes256GcmHandler{}
}
