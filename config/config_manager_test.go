package config

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const FILE_DIR = "config"

func LoadTomlFile(filename string) (ConfigHandler, error) {
	var hConf = CreateConfHandlerToml()
	err := hConf.LoadConf(filename)
	if err != nil {
		return nil, err
	}

	return hConf, nil
}

func CreateTomlFile(filename string) error {
	hConf, err := LoadTomlFile(filename)
	if err != nil {
		return err
	}

	hConf.SetDefaultValue()
	err = hConf.WriteDefValuesInNotExistsVlues()
	if err != nil {
		return err
	}

	return nil
}

func CreateAndLoadTomlFile(filename string) (ConfigHandler, error) {
	err := CreateTomlFile(filename)
	if err != nil {
		return nil, err
	}

	hConf, err := LoadTomlFile(filename)
	if err != nil {
		return nil, err
	}

	return hConf, nil
}

func ChangeTomlFile(filename string) error {
	var lines []string
	var remove_lines map[string]string = map[string]string{
		"realtime": "",
		"eventlog": "",
	}
	var replace_lines map[string]string = map[string]string{
		"disconnect_fail_count": "disconnect_fail_count = 100",
	}

	readFile, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	filescanner := bufio.NewScanner(readFile)

	for filescanner.Scan() {
		var flag bool
		var txt string = filescanner.Text()
		for k := range remove_lines {
			if strings.Contains(txt, k) {
				flag = true
				break
			}
		}

		for k, v := range replace_lines {
			if strings.Contains(txt, k) {
				txt = v
				break
			}
		}

		if flag {
			continue
		}

		lines = append(lines, txt)
	}

	writeFile, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(writeFile)
	for _, line := range lines {
		_, err = w.Write([]byte(line + "\n"))
		if err != nil {
			return err
		}
	}
	w.Flush()

	err = readFile.Close()
	if err != nil {
		return err
	}

	err = writeFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func EncodeKeyFunc(val string) (string, error) {
	key_text := []byte(val)

	block, err := aes.NewCipher(EncodeKey)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	fmt.Printf("Nonce : %v\n", len(Nonce))
	ciphertext := aesgcm.Seal(nil, Nonce, key_text, nil)
	ciphertext = append(ciphertext, Nonce...)

	encoded := hex.EncodeToString(ciphertext)

	return encoded, nil
}

func DecodeKeyFunc(val string) (string, error) {
	ciphertext, err := hex.DecodeString(val)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(EncodeKey)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, Nonce, ciphertext[:len(ciphertext)-len(Nonce)], nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func TestLoadConfToml(t *testing.T) {
	var filename = "config_test_loadconf.toml"
	err := CreateTomlFile(filename)
	if err != nil {
		t.Errorf("CreateTomlFile(%s) is failed.", filename)
	}

	err = os.Remove(filename)
	if err != nil {
		t.Errorf("RemoveTestDir(%s) is failed.", filename)
	}
}

func TestLoadConfAndModifyToml(t *testing.T) {
	var filename = "config_test_loadconf_modify.toml"
	err := CreateTomlFile(filename)
	if err != nil {
		t.Errorf("CreateTomlFile(%s) is failed.", filename)
	}

	err = ChangeTomlFile(filename)
	if err != nil {
		t.Errorf("ChangeTomlFile(%s) is failed.", filename)
	}

	hConf, err := LoadTomlFile(filename)
	if err != nil {
		t.Errorf("LoadTomlFile(%s) is failed.", filename)
	}

	hConf.SetDefaultValue()
	err = hConf.WriteDefValuesInNotExistsVlues()
	if err != nil {
		t.Errorf("WriteDefValuesInNotExistsVlues() is failed.")
	}

	assert.Equal(t, hConf.GetUint16Conf("manager.disconnect_fail_count", 100), uint16(100))

	err = os.Remove(filename)
	if err != nil {
		t.Errorf("RemoveTestDir(%s) is failed.", filename)
	}
}

func TestLoadConfAndCheckEncryption(t *testing.T) {
	var filename = "config_test_loadconf_checkencryption.toml"
	hConf, err := CreateAndLoadTomlFile(filename)
	if err != nil {
		t.Errorf("CreateAndLoadTomlFile(%s) is failed.", filename)
	}

	connection_user := DEFAULT_ONTUNE
	encode, err := EncodeKeyFunc(connection_user)
	if err != nil {
		t.Errorf("EncodeKeyFunc(%s) is failed.", connection_user)
	}

	stamp := hConf.GetStrConf("connection.stamp", "")
	assert.Equal(t, stamp, encode)

	err = os.Remove(filename)
	if err != nil {
		t.Errorf("RemoveTestDir(%s) is failed.", filename)
	}
}

func TestLoadConfAndCheckDecryption(t *testing.T) {
	var filename = "config_test_loadconf_checkdecryption.toml"
	hConf, err := CreateAndLoadTomlFile(filename)
	if err != nil {
		t.Errorf("CreateAndLoadTomlFile(%s) is failed.", filename)
	}

	stamp := hConf.GetStrConf("connection.stamp", "")
	plaintext, err := DecodeKeyFunc(stamp)
	if err != nil {
		t.Errorf("DecodeKeyFunc(%s) is failed.", stamp)
	}

	connection_user := DEFAULT_ONTUNE
	assert.Equal(t, string(plaintext), connection_user)

	err = os.Remove(filename)
	if err != nil {
		t.Errorf("RemoveTestDir(%s) is failed.", filename)
	}
}
