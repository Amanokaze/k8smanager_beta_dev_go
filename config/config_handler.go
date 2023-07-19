package config

type ConfigHandler interface {
	GetStrConf(key string, defVal string) string
	GetUint16Conf(key string, defVal uint16) uint16
	GetUint32Conf(key string, defVal uint32) uint32
	GetBoolConf(key string, defVal bool) bool
	SetDefaultValue()
	WriteDefValuesInNotExistsVlues() error
	LoadConf(confPath string) error
	ChangeKeyAfterOnceLoadConf(key string, changeKey string, defVal string, path string) (string, error)
	ChangeKey(key string, changeKey string, val string) error
}

// CreateConfHandler is a function that returns a pointer to a ConfigHandler instance.
func CreateConfHandlerEnv() ConfigHandler {
	return createEnvHandler()
}

func CreateConfHandlerToml() ConfigHandler {
	return createTomlHandler()
}
