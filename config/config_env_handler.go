package config

import (
	"bytes"
	"encoding/hex"
	"log"
	"onTuneKubeManager/common"
	"onTuneKubeManager/encryption"
	"onTuneKubeManager/logger"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/joho/godotenv"
)

var hEnc = encryption.CreateencryptionHandler()

type envConfigHandler struct {
	filepath string
	defBuf   bytes.Buffer
}

const SECRET_KEY = "HelloOntune"

// removeEncryptionKey is a function that removes encryption keys from the config file.
// configMap parameter is a map that contains the key and value of the config file.
// It will delete the keys, such as kube_cluster_userid_1, kube_cluster_password_1, db_user, db_password.
func removeEncryptionKey(configMap *map[string]string) {
	var clustercount int

	cc := (*configMap)["kube_cluster_count"]
	clustercount, err := strconv.Atoi(cc)
	if err != nil {
		return
	}

	for i := 0; i < clustercount; i++ {
		if _, ok := (*configMap)["kube_cluster_stamp_"+strconv.Itoa(i+1)]; ok {
			delete(*configMap, "kube_cluster_userid_"+strconv.Itoa(i+1))
			os.Unsetenv("kube_cluster_userid_" + strconv.Itoa(i+1))
		}

		if _, ok := (*configMap)["kube_cluster_code_"+strconv.Itoa(i+1)]; ok {
			delete(*configMap, "kube_cluster_password_"+strconv.Itoa(i+1))
			os.Unsetenv("kube_cluster_password_" + strconv.Itoa(i+1))
		}
	}

	if _, ok := (*configMap)["db_stamp"]; ok {
		delete(*configMap, "db_user")
		os.Unsetenv("db_user")
	}

	if _, ok := (*configMap)["db_code"]; ok {
		delete(*configMap, "db_password")
		os.Unsetenv("db_password")
	}
}

// removeDuplicateDefKey is a function that removes duplicate keys from the config file.
// path parameter is a path to the config file.
func removeDuplicateDefKey(path string) error {
	//TODO
	dat, err := os.ReadFile(path)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	lines := strings.Split(string(dat), "\n")

	var configMap = make(map[string]string)
	for _, line := range lines {
		if line != "" {
			index := strings.Index(line, "=")
			key := line[:index]
			configMap[key] = strings.Trim(line[index+1:], "\"")
		}
	}

	removeEncryptionKey(&configMap)

	var linearr []string = make([]string, 0)
	for key, val := range configMap {
		linearr = append(linearr, key+"="+val)
	}
	sort.Strings(linearr)

	var buf bytes.Buffer
	for _, val := range linearr {
		buf.WriteString(val)
		buf.WriteString("\n")
	}

	err = os.WriteFile(path, buf.Bytes(), 0644)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// WriteDefValuesInNotExistsVlues is a function that writes default values to the config file.
func (e *envConfigHandler) WriteDefValuesInNotExistsVlues() error {
	//defBuf
	err := logger.WriteBufFileLog(&e.defBuf, e.filepath)
	if err != nil {
		return err
	}

	err = removeDuplicateDefKey(e.filepath)
	if err != nil {
		return err
	}

	return nil
}

// writeBeginDefKey is a function that writes the key to the config file.
func (e *envConfigHandler) writeBeginDefKey(key *string) {
	e.defBuf.WriteString(*key)
	e.defBuf.WriteString("=")
	e.defBuf.WriteString("\"")
}

// writeEndDefKey is a function that writes the end of the key to the config file.
func (e *envConfigHandler) writeEndDefKey() {
	e.defBuf.WriteString("\"")
	e.defBuf.WriteString("\n")
}

// StringToAsciiBytes is a function that converts a string to a byte array.
// s parameter is a string to convert.
func StringToAsciiBytes(s string) []byte {
	t := make([]byte, utf8.RuneCountInString(s))
	i := 0
	for _, r := range s {
		t[i] = byte(r)
		i++
	}

	return t
}

// ChangeKey is a function that changes the key in the config file.
// key parameter is a key to change.
// changeKey parameter is a key to change.
// val parameter is a value to change.
func (e *envConfigHandler) ChangeKey(key string, changeKey string, val string) error {
	var envMap map[string]string
	envMap, err := godotenv.Read(e.filepath)
	if err != nil {
		return err
	}

	delete(envMap, key)
	envMap[changeKey] = val

	err = godotenv.Write(envMap, e.filepath)
	if err != nil {
		return err
	}

	return nil
}

// Empty Method
func (e *envConfigHandler) SetDefaultValue() {

}

// ChangeKeyAfterOnceLoadConf is a function that changes the key in the config file after loading the config file.
// key parameter is a key to change.
// changeKey parameter is a key to change.
// defVal parameter is a default value.
// path parameter is a path to the config file.
func (e *envConfigHandler) ChangeKeyAfterOnceLoadConf(key string, changeKey string, defVal string, path string) (string, error) {
	ret := os.Getenv(changeKey)
	if ret != "" {
		var hexDec, err = hex.DecodeString(ret)
		if err != nil {
			log.Println(err.Error())
			return defVal, err
		}

		ret = string(hexDec)
		index := strings.Index(ret, " ")
		filename := ret[:index]
		byteLen, err := strconv.Atoi(ret[index+1:])
		if err != nil {
			return defVal, err
		}

		encStr, err := logger.ReadFile(path + "/" + filename + ".dat")
		if err == nil {
			encBytes := []byte(encStr)
			encBytes = encBytes[:byteLen] // make([]byte, byteLen)
			retBytes, err := hEnc.Decrypt([]byte(SECRET_KEY), encBytes)

			if err == nil {
				return string(retBytes), nil
			} else {
				log.Println(err.Error())
				return defVal, err
			}
		} else {
			log.Println(err.Error())
			return defVal, err
		}
	} else {
		ret = os.Getenv(key)
		if ret == "" {
			//없는 경우는 key 하고 value 는 "" 으로 저장하자
			e.writeBeginDefKey(&key)
			e.defBuf.WriteString("")
			e.writeEndDefKey()

			return defVal, nil
		}

		//ret Encript 후 저장
		retBytes, err := hEnc.Encrypt([]byte(SECRET_KEY), []byte(ret))
		if err == nil {
			encdatfile := hex.EncodeToString([]byte(changeKey))
			bytesLen := strconv.Itoa(len(retBytes))

			err = logger.WriteFile(path+"/"+encdatfile+".dat", string(retBytes))
			if err == nil {
				var hexEnc = hex.EncodeToString([]byte(encdatfile + " " + bytesLen))
				err = e.ChangeKey(key, changeKey, hexEnc)
				if err != nil {
					log.Println(err.Error())
					return defVal, err
				}
			} else {
				common.LogManager.Error(err.Error())
				return defVal, err
			}
		} else {
			common.LogManager.Error(err.Error())
			return defVal, err
		}

		//TODO:  기존 key 삭제
		return ret, nil
	}
}

// GetStrConf is a function that returns a string value from the config file.
// key parameter is a key to change.
// defVal parameter is a default value.
func (e *envConfigHandler) GetStrConf(key string, defVal string) string {
	ret := os.Getenv(key)
	if ret == "" {
		e.writeBeginDefKey(&key)
		e.defBuf.WriteString(defVal)
		e.writeEndDefKey()

		return defVal
	}

	return ret
}

// GetUint16Conf is a function that returns a uint16 value from the config file.
// key parameter is a key to change.
// defVal parameter is a default value.
func (e *envConfigHandler) GetUint16Conf(key string, defVal uint16) uint16 {
	sVal := os.Getenv(key)
	if sVal == "" {
		e.writeBeginDefKey(&key)
		e.defBuf.WriteString(strconv.Itoa(int(defVal)))
		e.writeEndDefKey()
	}

	ret, err := strconv.ParseUint(sVal, 10, 16)
	if nil != err {
		return defVal
	}

	return uint16(ret)
}

// GetUint32Conf is a function that returns a uint32 value from the config file.
// key parameter is a key to change.
// defVal parameter is a default value.
func (e *envConfigHandler) GetUint32Conf(key string, defVal uint32) uint32 {
	sVal := os.Getenv(key)
	if sVal == "" {
		e.writeBeginDefKey(&key)
		e.defBuf.WriteString(strconv.Itoa(int(defVal)))
		e.writeEndDefKey()
	}

	ret, err := strconv.ParseUint(sVal, 10, 32)
	if nil != err {
		return defVal
	}

	return uint32(ret)
}

// GetBoolConf is a function that returns a bool value from the config file.
// key parameter is a key to change.
// defVal parameter is a default value.
func (e *envConfigHandler) GetBoolConf(key string, defVal bool) bool {
	sVal := os.Getenv(key)
	if sVal == "" {
		e.writeBeginDefKey(&key)
		e.defBuf.WriteString(strconv.FormatBool(defVal))
		e.writeEndDefKey()
	}

	ret, err := strconv.ParseBool(sVal)
	if err != nil {
		return defVal
	}

	return ret
}

// LoadConf is a function that loads the config file.
// confPath parameter is a path to the config file.
func (e *envConfigHandler) LoadConf(confPath string) error {
	e.filepath = confPath

	err := logger.CheckOrCreateFilePathDir(confPath, common.DbOsUser)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = logger.CreateFile(confPath)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = godotenv.Load(confPath)
	if err != nil {
		log.Println("Error loading .env file", err.Error())
		return err
	} else {
		return nil
	}
}

// createEnvHandler is a function that returns a pointer to a envConfigHandler instance.
func createEnvHandler() ConfigHandler {
	return &envConfigHandler{}
}
