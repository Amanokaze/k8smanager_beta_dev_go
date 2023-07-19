package encryption

import (
	"onTuneKubeManager/logger"
	"strings"
	"testing"
)

var (
	passphraseForSecretKey = "champKimWantsEveryoneToGG20230106"
	plaintext              = "ontune1"
	encScr                 string
	confPath               = "config.txt"
	hEnc                   = CreateencryptionHandler()
)

func TestAES256GSM2(t *testing.T) {
	ciphertext, err := hEnc.Encrypt([]byte(passphraseForSecretKey), []byte(plaintext))
	//ciphertext, err := AES256GSMEncrypt([]byte(passphraseForSecretKey), []byte(plaintext))
	if err != nil {
		t.Error(err)
	}

	plaintextBytes, err := hEnc.Decrypt([]byte(passphraseForSecretKey), ciphertext)
	//plaintextBytes, err := AES256GSMDecrypt([]byte(passphraseForSecretKey), ciphertext)
	if err != nil {
		t.Error(err)
	}

	if plaintext != string(plaintextBytes) {
		t.Errorf("plaintext %s is differ from decrypted ciphertext %s", plaintext, string(plaintextBytes))
	}
}

func TestAES256GSM2_Enc(t *testing.T) {
	//hEnc := CreateencryptionHandler()

	ciphertext, err := hEnc.Encrypt([]byte(passphraseForSecretKey), []byte(plaintext))
	if err != nil {
		t.Error(err)
	}

	encScr = string(ciphertext)

	//logger.CreateFile(confPath)
	err = logger.WriteFile(confPath, encScr)
	if err != nil {
		t.Error(err)
	}
	//logger.WriteBytes(confPath, ciphertext)

	println("[enc] ", plaintext, " ----> ", encScr)
}

func TestAES256GSM2_Dec(t *testing.T) {
	//hEnc := CreateencryptionHandler()
	//logger.WriteFile(confPath, encScr)
	encScr, err := logger.ReadFile(confPath)
	if err != nil {
		t.Errorf("readfile(%s) error %s", confPath, err.Error())
		return
	}

	encScr = strings.TrimRight(encScr, "\x00")
	endBytes := []byte(encScr)

	plaintextBytes, err := hEnc.Decrypt([]byte(passphraseForSecretKey), endBytes)
	//plaintextBytes, err := AES256GSMDecrypt([]byte(passphraseForSecretKey), endBytes)
	if err != nil {
		t.Error(err)
	}
	println("[dec] ", encScr, " -----> ", string(plaintextBytes))

	if plaintext != string(plaintextBytes) {
		t.Errorf("plaintext %s is differ from decrypted ciphertext %s", plaintext, string(plaintextBytes))
	}
}
