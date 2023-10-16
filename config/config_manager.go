package config

import (
	"fmt"
	"sync"
)

type ManagerConfig struct {
	managername             string
	debuglog                bool
	realtimeinterval        uint16
	avginterval             uint16
	resourceinterval        uint16
	resourcechangeinterval  uint16
	rateinterval            uint16
	eventloginterval        uint16
	disconnectfailcount     uint16
	retryconnectioninterval uint16
	clustercount            uint16
	kubeConf                []KubeConfig
	dbConf                  DBConfig
	tableConf               TableConfig
	//datalog          bool
	//primaryDBConf DBConfig
}

var (
	instance *ManagerConfig
	confOnce sync.Once
)

// GetOrCreateManagerConfig is a function that returns a pointer to a ManagerConfig instance.
func GetOrCreateManagerConfig() *ManagerConfig {
	confOnce.Do(func() {
		instance = new(ManagerConfig)
	})

	return instance
}

// LoadServerConfigurationToml is a function that loads the configuration file.
func (mConf *ManagerConfig) LoadServerConfigurationToml(confPath string) error {
	hConf := CreateConfHandlerToml()

	err := hConf.LoadConf(confPath)
	if err != nil {
		return err
	}

	hConf.SetDefaultValue()

	mConf.managername = hConf.GetStrConf("manager.name", "KubeManager")
	mConf.disconnectfailcount = hConf.GetUint16Conf("manager.disconnect_fail_count", 10)

	mConf.dbConf.databaseOwner = hConf.GetStrConf("database.user", "postgres")
	mConf.dbConf.databaseDocker = hConf.GetBoolConf("database.docker", false)
	mConf.dbConf.databaseDockerpath = hConf.GetStrConf("database.docker_path", "/var/lib/postgresql/data")
	mConf.dbConf.maxConnection = hConf.GetUint16Conf("database.max_connections", 100)

	mConf.dbConf.connectionHost = hConf.GetStrConf("connection.host", "127.0.0.1")
	mConf.dbConf.connectionPort = hConf.GetUint16Conf("connection.port", 5432)
	mConf.dbConf.connectionName = hConf.GetStrConf("connection.database", "ontune")
	mConf.dbConf.connectionUser = hConf.GetStrConf("connection.user", "ontune")
	mConf.dbConf.connectionPassword = hConf.GetStrConf("connection.password", "xxxxxx")
	mConf.dbConf.connectionSslmode = hConf.GetStrConf("connection.sslmode", "disable")

	mConf.tableConf.autovacuum = hConf.GetBoolConf("table.autovacuum", true)
	mConf.tableConf.initvacuum = hConf.GetBoolConf("table.initvacuum", false)
	mConf.tableConf.tablespace = hConf.GetBoolConf("table.resource.tablespace", false)
	mConf.tableConf.tablespacename = hConf.GetStrConf("table.resource.tablespace_name", "")
	mConf.tableConf.tablespacepath = hConf.GetStrConf("table.resource.tablespace_path", "")
	mConf.tableConf.shorttermduration = hConf.GetUint16Conf("table.short_term.duration", 10)
	mConf.tableConf.shorttermtablespace = hConf.GetBoolConf("table.short_term.tablespace", false)
	mConf.tableConf.shorttermtablespacename = hConf.GetStrConf("table.short_term.tablespace_name", "")
	mConf.tableConf.shorttermtablespacepath = hConf.GetStrConf("table.short_term.tablespace_path", "")
	mConf.tableConf.longtermduration = hConf.GetUint16Conf("table.long_term.duration", 600)
	mConf.tableConf.longtermtablespace = hConf.GetBoolConf("table.long_term.tablespace", false)
	mConf.tableConf.longtermtablespacename = hConf.GetStrConf("table.long_term.tablespace_name", "")
	mConf.tableConf.longtermtablespacepath = hConf.GetStrConf("table.long_term.tablespace_path", "")

	mConf.clustercount = hConf.GetUint16Conf("kubernetes.cluster_count", 1)
	mConf.kubeConf = make([]KubeConfig, mConf.clustercount)
	for i := 0; i < int(mConf.clustercount); i++ {
		mConf.kubeConf[i].kube_cluster_conffile = hConf.GetStrConf(fmt.Sprintf("kubernetes.clusters.%d.config_file", i), fmt.Sprintf(".kube.%d.conf", i+1))
	}

	mConf.realtimeinterval = hConf.GetUint16Conf("interval.realtime", 10)
	mConf.avginterval = hConf.GetUint16Conf("interval.avg", 600)
	mConf.resourceinterval = hConf.GetUint16Conf("interval.resource", 600)
	mConf.resourcechangeinterval = hConf.GetUint16Conf("interval.resource_change", 60)
	mConf.eventloginterval = hConf.GetUint16Conf("interval.eventlog", 10)
	mConf.retryconnectioninterval = hConf.GetUint16Conf("interval.retry_connection", 600)
	mConf.rateinterval = hConf.GetUint16Conf("interval.rate", 60)

	mConf.debuglog = hConf.GetBoolConf("log.debug", false)

	err = hConf.WriteDefValuesInNotExistsVlues()
	if err != nil {
		return err
	}

	return nil
}

// LoadServerConf is a function that loads the configuration file.
func (mConf *ManagerConfig) LoadServerConfEnv(confPath string) {
	hConf := CreateConfHandlerEnv()

	err := hConf.LoadConf(confPath)
	if err != nil {
		return
	}

	config_path := "config"

	mConf.dbConf.connectionUser, err = hConf.ChangeKeyAfterOnceLoadConf("db_user", "db_stamp", "ontune", config_path)
	if err != nil {
		fmt.Println(err)
		return
	}

	mConf.dbConf.connectionPassword, err = hConf.ChangeKeyAfterOnceLoadConf("db_password", "db_code", "xxxxxx", config_path)
	if err != nil {
		fmt.Println(err)
		return
	}

	mConf.managername = hConf.GetStrConf("managername", "KubeManager")
	mConf.debuglog = hConf.GetBoolConf("debuglog", false)
	mConf.realtimeinterval = hConf.GetUint16Conf("realtimeinterval", 10)
	mConf.resourceinterval = hConf.GetUint16Conf("resourceinterval", 600)
	mConf.resourcechangeinterval = hConf.GetUint16Conf("resourcechangeinterval", 60)
	mConf.rateinterval = hConf.GetUint16Conf("rateinterval", 60)
	mConf.eventloginterval = hConf.GetUint16Conf("eventloginterval", 10)
	mConf.disconnectfailcount = hConf.GetUint16Conf("disconnectfailcount", 10)
	mConf.retryconnectioninterval = hConf.GetUint16Conf("retryconnectioninterval", 600)

	mConf.clustercount = hConf.GetUint16Conf("kube_cluster_count", 1)
	mConf.kubeConf = make([]KubeConfig, mConf.clustercount)
	for i := 0; i < int(mConf.clustercount); i++ {
		mConf.kubeConf[i].kube_cluster_conffile = hConf.GetStrConf(fmt.Sprintf("kube_cluster_conffile_%d", i+1), fmt.Sprintf(".kube.%d.conf", i+1))
	}

	mConf.dbConf.connectionHost = hConf.GetStrConf("db_host", "127.0.0.1")
	mConf.dbConf.connectionPort = hConf.GetUint16Conf("db_port", 5432)
	mConf.dbConf.connectionName = hConf.GetStrConf("db_name", "ontune")
	mConf.dbConf.connectionSslmode = hConf.GetStrConf("db_sslmode", "disable")
	mConf.dbConf.maxConnection = hConf.GetUint16Conf("db_maxconnection", 100)
	mConf.dbConf.databaseOwner = hConf.GetStrConf("db_osuser", "postgres")
	mConf.dbConf.databaseDocker = hConf.GetBoolConf("db_docker", false)
	mConf.dbConf.databaseDockerpath = hConf.GetStrConf("db_dockerpath", "/var/lib/postgresql/data")

	mConf.tableConf.autovacuum = hConf.GetBoolConf("autovacuum", true)
	mConf.tableConf.tablespace = hConf.GetBoolConf("tablespace", false)
	mConf.tableConf.tablespacename = hConf.GetStrConf("tablespacename", "")
	mConf.tableConf.tablespacepath = hConf.GetStrConf("tablespacepath", getDefaultTablespacePath())
	mConf.tableConf.shorttermduration = hConf.GetUint16Conf("shorttermduration", 10)
	mConf.tableConf.shorttermtablespace = hConf.GetBoolConf("shorttermtablespace", false)
	mConf.tableConf.shorttermtablespacename = hConf.GetStrConf("shorttermtablespacename", "")
	mConf.tableConf.shorttermtablespacepath = hConf.GetStrConf("shorttermtablespacepath", getDefaultTablespacePath())

	err = hConf.WriteDefValuesInNotExistsVlues()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// GetTableSpaceName is a function that returns the name of the tablespace.
func (mgrConf *ManagerConfig) GetTableSpaceName() string {
	return mgrConf.tableConf.tablespacename
}

// GetTableSpace is a function that returns whether to use tablespace.
func (mgrConf *ManagerConfig) GetTableSpace() bool {
	return mgrConf.tableConf.tablespace
}

// GetTableSpacePath is a function that returns the path of the tablespace.
func (mgrConf *ManagerConfig) GetTableSpacePath() string {
	return mgrConf.tableConf.tablespacepath
}

// GetTableShorttermSpace is a function that returns whether to use short-term tablespace.
func (mgrConf *ManagerConfig) GetTableShorttermSpace() bool {
	return mgrConf.tableConf.shorttermtablespace
}

// GetTableLongtermSpace
func (mgrConf *ManagerConfig) GetTableLongtermSpace() bool {
	return mgrConf.tableConf.longtermtablespace
}

// GetTableShorttermSpaceCountName is a function that returns the name of the short-term tablespace.
func (mgrConf *ManagerConfig) GetTableShorttermSpaceName() string {
	return mgrConf.tableConf.GetShorttermTableSpaceName()
}

// GetTableLongtermSpaceCountName is a function that returns the name of the long-term tablespace.
func (mgrConf *ManagerConfig) GetTableLongtermSpaceName() string {
	return mgrConf.tableConf.GetLongtermTableSpaceName()
}

// GetTableShorttermSpacePath is a function that returns the path of the short-term tablespace.
func (mgrConf *ManagerConfig) GetTableShorttermSpacePath() string {
	return mgrConf.tableConf.GetShorttermTableSpacePath()
}

// GetTableLongtermSpacePath is a function that returns the path of the long-term tablespace.
func (mgrConf *ManagerConfig) GetTableLongtermSpacePath() string {
	return mgrConf.tableConf.GetLongtermTableSpacePath()
}

// GetTableShorttermSpacePath is a function that returns the path of the short-term tablespace.
func (mgrConf *ManagerConfig) GetInitVacuum() bool {
	return mgrConf.tableConf.initvacuum
}

// GetTableLongtermSpacePath is a function that returns the path of the long-term tablespace.
func (mgrConf *ManagerConfig) GetAutoVacuum() bool {
	return mgrConf.tableConf.autovacuum
}

// GetTableShorttermSpacePath is a function that returns the path of the short-term tablespace.
func (mgrConf *ManagerConfig) GetShorttermduration() int {
	return int(mgrConf.tableConf.shorttermduration)
}

// GetTableLongtermSpacePath is a function that returns the path of the long-term tablespace.
func (mgrConf *ManagerConfig) GetLongtermduration() int {
	return int(mgrConf.tableConf.longtermduration)
}

// GetManagerName is a function that returns the path of the long-term tablespace.
func (mgrConf *ManagerConfig) GetManagerName() string {
	return mgrConf.managername
}

// GetCluster is a function that returns the path of the long-term tablespace.
func (mgrConf *ManagerConfig) GetCluster() uint16 {
	return mgrConf.clustercount
}

// GetRealtimetimeInterval is a function that returns the real-time interval.
func (mgrConf *ManagerConfig) GetRealtimeInterval() uint16 {
	return mgrConf.realtimeinterval
}

// GetAvgInterval is a function that returns the average interval.
func (mgrConf *ManagerConfig) GetAvgInterval() uint16 {
	return mgrConf.avginterval
}

// GetEventlogInterval is a function that returns the event log interval.
func (mgrConf *ManagerConfig) GetEventlogInterval() uint16 {
	return mgrConf.eventloginterval
}

// GetResourceInterval is a function that returns the resource interval.
func (mgrConf *ManagerConfig) GetResourceInterval() uint16 {
	return mgrConf.resourceinterval
}

// GetResourceChangeInterval is a function that returns the resource change interval.
func (mgrConf *ManagerConfig) GetResourceChangeInterval() uint16 {
	return mgrConf.resourcechangeinterval
}

// GetRateInterval is a function that returns the rate interval.
func (mgrConf *ManagerConfig) GetRateInterval() uint16 {
	return mgrConf.rateinterval
}

// GetDisconnectFailCount is a function that returns the number of disconnect failures.
func (mgrConf *ManagerConfig) GetDisconnectFailCount() uint16 {
	return mgrConf.disconnectfailcount
}

// GetRetryConnectionInterval is a function that returns the retry connection interval.
func (mgrConf *ManagerConfig) GetRetryConnectionInterval() uint16 {
	return mgrConf.retryconnectioninterval
}

// GetDebugLog is a function that returns the Debug Log.
func (mgrConf *ManagerConfig) GetDebugLog() bool {
	return mgrConf.debuglog
}

// GetKubeConfig is a function that returns the Kube Configuration log.
func (mgrConf *ManagerConfig) GetKubeConfig() []KubeConfig {
	return mgrConf.kubeConf
}

// GetDBConfig is a function that returns the Database Configuration log.
func (mgrConf *ManagerConfig) GetDBConfig() DBConfig {
	return mgrConf.dbConf
}

// GetTableConfig is a function that returns the Table Configuration log.
func (mgrConf *ManagerConfig) GetTableConfig() TableConfig {
	return mgrConf.tableConf
}
