package config

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"onTuneKubeManager/common"
	"onTuneKubeManager/logger"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

const DEFAULT_ONTUNE = "ontune"
const SECTION_CONNECTION = "connection"
const SECTION_DATABASE = "database"
const SECTION_TABLE = "table"
const ITEM_USER = "user"
const ITEM_PASSWORD = "password"

var EncodeKey = []byte("kuObeNmaTnaUgeNrE0981029384765@!")
var Nonce = []byte("enutno!^@%#$")

type tomlConfigManager struct {
	Name                string `toml:"name"`
	DisconnectFailCount int    `toml:"disconnect_fail_count"`
}

type tomlConfigDatabase struct {
	Docker         bool   `toml:"docker"`
	DockerPath     string `toml:"docker_path"`
	User           string `toml:"user"`
	MaxConnections int    `toml:"max_connections"`
}

type tomlConfigConnection struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Database string `toml:"database"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Code     string `toml:"code"`
	Stamp    string `toml:"stamp"`
	Sslmode  string `toml:"sslmode"`
}

type tomlConfigTableTerm struct {
	Duration       int    `toml:"duration"`
	TablespaceName string `toml:"tablespace_name"`
	Tablespace     bool   `toml:"tablespace"`
	TablespacePath string `toml:"tablespace_path"`
}

type tomlConfigTableResource struct {
	TablespaceName string `toml:"tablespace_name"`
	Tablespace     bool   `toml:"tablespace"`
	TablespacePath string `toml:"tablespace_path"`
}

type tomlConfigTable struct {
	Autovacuum bool                    `toml:"autovacuum"`
	Initvacuum bool                    `toml:"initvacuum"`
	ShortTerm  tomlConfigTableTerm     `toml:"short_term"`
	LongTerm   tomlConfigTableTerm     `toml:"long_term"`
	Resource   tomlConfigTableResource `toml:"resource"`
}

type tomlConfigKubernetesCluster struct {
	ConfigFile string `toml:"config_file"`
}

type tomlConfigKubernetes struct {
	ClusterCount int                           `toml:"cluster_count"`
	Clusters     []tomlConfigKubernetesCluster `toml:"clusters"`
}

type tomlConfigInterval struct {
	Eventlog        int `toml:"eventlog"`
	Rate            int `toml:"rate"`
	Realtime        int `toml:"realtime"`
	Avg             int `toml:"avg"`
	Resource        int `toml:"resource"`
	ResourceChange  int `toml:"resource_change"`
	RetryConnection int `toml:"retry_connection"`
}

type tomlConfigLog struct {
	Debug bool `toml:"debug"`
}

type tomlConfig struct {
	Manager    tomlConfigManager    `toml:"manager"`
	Database   tomlConfigDatabase   `toml:"database"`
	Connection tomlConfigConnection `toml:"connection"`
	Table      tomlConfigTable      `toml:"table"`
	Kubernetes tomlConfigKubernetes `toml:"kubernetes"`
	Interval   tomlConfigInterval   `toml:"interval"`
	Log        tomlConfigLog        `toml:"log"`
}

type tomlConfigHandler struct {
	filepath string
	tomlConfig
}

// WriteDefValuesInNotExistsVlues is a function that writes default values to the config file.
func (t *tomlConfigHandler) WriteDefValuesInNotExistsVlues() error {
	var buf bytes.Buffer = bytes.Buffer{}
	enc := toml.NewEncoder(&buf)
	err := enc.Encode(t.tomlConfig)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	file, err := os.OpenFile(t.filepath, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, &buf)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	t.writeBeginDefKey(nil)
	t.writeEndDefKey()

	return nil
}

// Empty Method
func (t *tomlConfigHandler) writeBeginDefKey(key *string) {
	fmt.Printf("")
}

// Empty Method
func (t *tomlConfigHandler) writeEndDefKey() {
	fmt.Printf("")
}

// ChangeKey is a function that changes the key in the config file.
func (t *tomlConfigHandler) ChangeKey(key string, changeKey string, val string) error {
	keys := strings.Split(key, ".")
	changeKeys := strings.Split(changeKey, ".")
	if len(keys) == 2 {
		switch keys[0] {
		case SECTION_CONNECTION:
			switch keys[1] {
			case ITEM_USER:
				if changeKeys[0] == keys[0] && keys[1] == "user" {
					t.tomlConfig.Connection.User = ""
					t.tomlConfig.Connection.Stamp = val
				}
			case ITEM_PASSWORD:
				if changeKeys[0] == keys[0] && keys[1] == "password" {
					t.tomlConfig.Connection.Password = ""
					t.tomlConfig.Connection.Code = val
				}
			}
		}
	}

	return nil
}

// EncodeKey is a function that encodes the key in the config file.
func (t *tomlConfigHandler) EncodeKey(val string) (string, error) {
	key_text := []byte(val)

	block, err := aes.NewCipher(EncodeKey)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, Nonce, key_text, nil)
	ciphertext = append(ciphertext, Nonce...)

	encoded := hex.EncodeToString(ciphertext)

	return encoded, nil
}

// DecodeKey is a function that decodes the key in the config file.
func (t *tomlConfigHandler) DecodeKey(val string) (string, error) {
	decoded, err := hex.DecodeString(val)
	if err != nil {
		return "", err
	}

	nonce := decoded[len(decoded)-12:]
	ciphertext := decoded[:len(decoded)-12]

	block, err := aes.NewCipher(EncodeKey)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// ChangeKeyAfterOnceLoadConf is a function that changes the key in the config file after loading the config file.
// path parameter is not used.
func (t *tomlConfigHandler) ChangeKeyAfterOnceLoadConf(key string, changeKey string, defVal string, path string) (string, error) {
	encoded, err := t.EncodeKey(defVal)
	if err != nil {
		return "", err
	}

	err = t.ChangeKey(key, changeKey, encoded)
	if err != nil {
		return "", err
	}

	return encoded, nil
}

// GetStrConf is a function that returns a string value from the config file.
func (t *tomlConfigHandler) GetStrConf(key string, defVal string) string {
	keys := strings.Split(key, ".")
	if len(keys) >= 2 {
		switch keys[0] {
		case "manager":
			switch keys[1] {
			case "name":
				return t.tomlConfig.Manager.Name
			}
		case SECTION_DATABASE:
			switch keys[1] {
			case ITEM_USER:
				return t.tomlConfig.Database.User
			case "docker_path":
				return t.tomlConfig.Database.DockerPath
			}
		case SECTION_CONNECTION:
			switch keys[1] {
			case "host":
				return t.tomlConfig.Connection.Host
			case SECTION_DATABASE:
				return t.tomlConfig.Connection.Database
			case ITEM_USER:
				val, err := t.DecodeKey(t.tomlConfig.Connection.Stamp)
				if err != nil {
					log.Println(err.Error())
					return ""
				}

				return val
			case ITEM_PASSWORD:
				val, err := t.DecodeKey(t.tomlConfig.Connection.Code)
				if err != nil {
					log.Println(err.Error())
					return ""
				}

				return val
			case "code":
				return t.tomlConfig.Connection.Code
			case "stamp":
				return t.tomlConfig.Connection.Stamp
			case "sslmode":
				return t.tomlConfig.Connection.Sslmode
			}
		case SECTION_TABLE:
			switch fmt.Sprintf("%s.%s", keys[1], keys[2]) {
			case "short_term.tablespace_name":
				return t.tomlConfig.Table.ShortTerm.TablespaceName
			case "short_term.tablespace_path":
				return t.tomlConfig.Table.ShortTerm.TablespacePath
			case "long_term.tablespace_name":
				return t.tomlConfig.Table.LongTerm.TablespaceName
			case "long_term.tablespace_path":
				return t.tomlConfig.Table.LongTerm.TablespacePath
			case "resource.tablespace_name":
				return t.tomlConfig.Table.Resource.TablespaceName
			case "resource.tablespace_path":
				return t.tomlConfig.Table.Resource.TablespacePath
			}
		case "kubernetes":
			switch keys[1] {
			case "clusters":
				if len(keys) == 4 && keys[3] == "config_file" {
					clusterNum, err := strconv.Atoi(keys[2])
					if err != nil {
						log.Println(err.Error())
						return ""
					}

					return t.tomlConfig.Kubernetes.Clusters[clusterNum].ConfigFile
				}
			}
		case "log":
			switch keys[1] {
			case "debug":
				return strconv.FormatBool(t.tomlConfig.Log.Debug)
			}
		}
	}

	return defVal
}

// GetUint16Conf is a function that returns a uint16 value from the config file.
func (t *tomlConfigHandler) GetUint16Conf(key string, defVal uint16) uint16 {
	keys := strings.Split(key, ".")
	if len(keys) >= 2 {
		switch keys[0] {
		case "manager":
			switch keys[1] {
			case "disconnect_fail_count":
				return uint16(t.tomlConfig.Manager.DisconnectFailCount)
			}
		case SECTION_DATABASE:
			switch keys[1] {
			case "max_connections":
				return uint16(t.tomlConfig.Database.MaxConnections)
			}
		case SECTION_CONNECTION:
			switch keys[1] {
			case "port":
				return uint16(t.tomlConfig.Connection.Port)
			}
		case SECTION_TABLE:
			switch fmt.Sprintf("%s.%s", keys[1], keys[2]) {
			case "short_term.duration":
				return uint16(t.tomlConfig.Table.ShortTerm.Duration)
			case "long_term.duration":
				return uint16(t.tomlConfig.Table.LongTerm.Duration)
			}
		case "kubernetes":
			switch keys[1] {
			case "cluster_count":
				cluster_size := len(t.tomlConfig.Kubernetes.Clusters)
				if cluster_size < t.tomlConfig.Kubernetes.ClusterCount {
					t.tomlConfig.Kubernetes.ClusterCount = cluster_size
				}

				return uint16(t.tomlConfig.Kubernetes.ClusterCount)
			}
		case "interval":
			switch keys[1] {
			case "eventlog":
				return uint16(t.tomlConfig.Interval.Eventlog)
			case "rate":
				return uint16(t.tomlConfig.Interval.Rate)
			case "realtime":
				return uint16(t.tomlConfig.Interval.Realtime)
			case "avg":
				return uint16(t.tomlConfig.Interval.Avg)
			case "resource":
				return uint16(t.tomlConfig.Interval.Resource)
			case "resource_change":
				return uint16(t.tomlConfig.Interval.ResourceChange)
			case "retry_connection":
				return uint16(t.tomlConfig.Interval.RetryConnection)
			}
		}
	}

	return uint16(0)
}

// Empty Method
func (t *tomlConfigHandler) GetUint32Conf(key string, defVal uint32) uint32 {
	return uint32(0)
}

// GetBoolConf is a function that returns a bool value from the config file.
func (t *tomlConfigHandler) GetBoolConf(key string, defVal bool) bool {
	keys := strings.Split(key, ".")
	if len(keys) >= 2 {
		switch keys[0] {
		case SECTION_DATABASE:
			switch keys[1] {
			case "docker":
				return t.tomlConfig.Database.Docker
			}
		case SECTION_TABLE:
			if len(keys) == 2 {
				switch keys[1] {
				case "autovacuum":
					return t.tomlConfig.Table.Autovacuum
				case "initvacuum":
					return t.tomlConfig.Table.Initvacuum
				}
			}

			if len(keys) > 2 {
				switch fmt.Sprintf("%s.%s", keys[1], keys[2]) {
				case "short_term.tablespace":
					return t.tomlConfig.Table.ShortTerm.Tablespace
				case "long_term.tablespace":
					return t.tomlConfig.Table.LongTerm.Tablespace
				case "resource.tablespace":
					return t.tomlConfig.Table.Resource.Tablespace
				}
			}
		case "log":
			switch keys[1] {
			case "debug":
				return t.tomlConfig.Log.Debug
			}
		}
	}

	return defVal
}

// SetDefaultValue is a function that sets the default value.
func (t *tomlConfigHandler) SetDefaultValue() {
	if t.tomlConfig.Manager.Name == "" {
		t.tomlConfig.Manager.Name = "KubeManager"
	}

	if t.tomlConfig.Manager.DisconnectFailCount == 0 {
		t.tomlConfig.Manager.DisconnectFailCount = 10
	}

	if t.tomlConfig.Database.User == "" {
		t.tomlConfig.Database.User = "postgres"
	}

	if t.tomlConfig.Database.MaxConnections == 0 {
		t.tomlConfig.Database.MaxConnections = 100
	}

	if t.tomlConfig.Database.DockerPath == "" {
		t.tomlConfig.Database.DockerPath = "/var/lib/postgresql/data"
	}

	if t.tomlConfig.Connection.Host == "" {
		t.tomlConfig.Connection.Host = "127.0.0.1"
	}

	if t.tomlConfig.Connection.Port == 0 {
		t.tomlConfig.Connection.Port = 5432
	}

	if t.tomlConfig.Connection.Database == "" {
		t.tomlConfig.Connection.Database = DEFAULT_ONTUNE
	}

	if t.tomlConfig.Connection.User == "" && t.tomlConfig.Connection.Stamp == "" {
		t.tomlConfig.Connection.User = DEFAULT_ONTUNE
	}

	if t.tomlConfig.Connection.User != "" {
		_, err := t.ChangeKeyAfterOnceLoadConf("connection.user", "connection.stamp", t.tomlConfig.Connection.User, "")
		if err != nil {
			log.Println(err.Error())
		}
	}

	if t.tomlConfig.Connection.Password == "" && t.tomlConfig.Connection.Code == "" {
		t.tomlConfig.Connection.Password = DEFAULT_ONTUNE
	}

	if t.tomlConfig.Connection.Password != "" {
		_, err := t.ChangeKeyAfterOnceLoadConf("connection.password", "connection.code", t.tomlConfig.Connection.Password, "")
		if err != nil {
			log.Println(err.Error())
		}
	}

	if t.tomlConfig.Connection.Sslmode == "" {
		t.tomlConfig.Connection.Sslmode = "disable"
	}

	if t.tomlConfig.Table.ShortTerm.Duration == 0 {
		t.tomlConfig.Table.ShortTerm.Duration = 10
	}

	if t.tomlConfig.Table.LongTerm.Duration == 0 {
		t.tomlConfig.Table.LongTerm.Duration = 600
	}

	// 차후 사용 가능성을 고려하여 주석처리
	// if t.tomlConfig.Table.ShortTerm.TablespaceName == "" {
	// 	t.tomlConfig.Table.ShortTerm.TablespaceName = "ontuneshortterm"
	// }

	// if t.tomlConfig.Table.ShortTerm.Path == "" {
	// 	t.tomlConfig.Table.ShortTerm.Path = getDefaultTablespacePath()
	// }

	// if t.tomlConfig.Table.Resource.TablespaceName == "" {
	// 	t.tomlConfig.Table.Resource.TablespaceName = "ontunedata"
	// }

	// if t.tomlConfig.Table.Resource.Path == "" {
	// 	t.tomlConfig.Table.Resource.Path = getDefaultTablespacePath()
	// }

	if t.tomlConfig.Kubernetes.ClusterCount == 0 {
		t.tomlConfig.Kubernetes.ClusterCount = 1
	}

	if t.tomlConfig.Kubernetes.Clusters == nil {
		t.tomlConfig.Kubernetes.Clusters = make([]tomlConfigKubernetesCluster, t.tomlConfig.Kubernetes.ClusterCount)

		for i := 0; i < t.tomlConfig.Kubernetes.ClusterCount; i++ {
			if t.tomlConfig.Kubernetes.Clusters[i].ConfigFile == "" {
				t.tomlConfig.Kubernetes.Clusters[i].ConfigFile = fmt.Sprintf(".kube.%d.conf", i+1)
			}
		}
	}

	if t.tomlConfig.Interval.Eventlog == 0 {
		t.tomlConfig.Interval.Eventlog = 10
	}

	if t.tomlConfig.Interval.Rate == 0 {
		t.tomlConfig.Interval.Rate = 60
	}

	if t.tomlConfig.Interval.Realtime == 0 {
		t.tomlConfig.Interval.Realtime = 10
	}

	if t.tomlConfig.Interval.Avg == 0 {
		t.tomlConfig.Interval.Avg = 600
	}

	if t.tomlConfig.Interval.Resource == 0 {
		t.tomlConfig.Interval.Resource = 600
	}

	if t.tomlConfig.Interval.ResourceChange == 0 {
		t.tomlConfig.Interval.ResourceChange = 60
	}

	if t.tomlConfig.Interval.RetryConnection == 0 {
		t.tomlConfig.Interval.RetryConnection = 600
	}
}

func (t *tomlConfigHandler) TomlRead(confPath string) (tomlConfig, error) {
	var tomlConfig tomlConfig
	_, err := toml.DecodeFile(confPath, &tomlConfig)
	if err != nil {
		log.Println("Error loading .toml file", err.Error())
		return tomlConfig, err
	}

	return tomlConfig, nil
}

// LoadConf is a function that loads the config file.
// confPath parameter is a path to the config file.
func (t *tomlConfigHandler) LoadConf(confPath string) error {
	t.filepath = confPath

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

	tomlConfig, err := t.TomlRead(confPath)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	t.tomlConfig = tomlConfig

	return nil
}

// createTomlHandler is a function that returns a pointer to a tomlConfigHandler instance.
func createTomlHandler() ConfigHandler {
	return &tomlConfigHandler{}
}
