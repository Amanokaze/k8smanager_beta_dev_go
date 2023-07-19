//go:generate goversioninfo

package main

import (
	"fmt"
	"onTuneKubeManager/common"
	"onTuneKubeManager/config"
	"onTuneKubeManager/database"
	"onTuneKubeManager/event"
	"onTuneKubeManager/kubeapi"
	"onTuneKubeManager/logger"
	"time"
)

const (
	FILE_VERSION = "4.1.13.2"
)

func main() {
	conf := config.GetOrCreateManagerConfig()
	//conf.LoadServerConfEnv("config/config.env")
	err := conf.LoadServerConfigurationToml("config/config.toml")
	if err != nil {
		fmt.Println("LoadServerConfigurationToml failed")
		return
	}

	common.LogManager = logger.CreateOnTuneLogManager(conf.GetDebugLog())
	if common.LogManager == nil {
		fmt.Println("CreateOnTuneLogManager failed")
		return
	}
	common.LogManager.WriteLog("Start onTune Kube Manager " + FILE_VERSION)

	managername := conf.GetManagerName()
	dbinfo := conf.GetDBConfig()
	common.DbOsUser = dbinfo.GetDBOSUser()

	//table config
	common.InitVacuum = conf.GetInitVacuum()
	common.AutoVacuum = conf.GetAutoVacuum()
	common.TableSpace = conf.GetTableSpace()
	common.TableSpaceName = conf.GetTableSpaceName()
	common.TableSpacePath = conf.GetTableSpacePath()
	common.ShorttermDuration = conf.GetShorttermduration()
	common.ShorttermTableSpace = conf.GetTableShorttermSpace()
	common.ShorttermTableSpaceName = conf.GetTableShorttermSpaceName()
	common.ShorttermTableSpacePath = conf.GetTableShorttermSpacePath()
	common.DisconnectFailCount = conf.GetDisconnectFailCount()
	common.RealtimeInterval = int(conf.GetRealtimeInterval())
	common.ResourceInterval = int(conf.GetResourceInterval())
	common.RateInterval = int(conf.GetRateInterval())
	common.EventlogInterval = int(conf.GetEventlogInterval())
	common.KubeClusterCount = conf.GetCluster()
	common.MapInit()

	for ip := range common.ClusterID {
		fmt.Println("cluster IP == " + ip)
	}

	clientset, host, kubeconfig := kubeapi.GetClientsetConfig(common.LogManager, conf.GetKubeConfig())

	database.MakeDBConn(&dbinfo)
	common.LogManager.Debug("MakeDBConn Complete")

	go database.DropTable()

	common.LogManager.Debug("QueryManagerinfo Previous Exection")
	common.ManagerID = database.QueryManagerinfo(managername, "Description", dbinfo.GetHost())
	for i := 0; i < int(common.KubeClusterCount); i++ {
		common.ClusterID[host[i]] = database.QueryClusterinfo(common.ManagerID, kubeconfig[i].GetName(), kubeconfig[i].GetContext(), host[i])
	}
	common.LogManager.Debug("QueryManagerinfo is completed")

	database.QueryUnusedClusterReset()
	common.LogManager.Debug("QueryUnusedClusterReset is completed")

	database.InitMapData()
	database.InitMapEventlogData()
	database.InitMapMetricData()

	event.SetEvent()

	go database.ResourceSender()
	go database.EventlogSender()
	go database.MetricSender()

	go database.DailyPerfTable()
	go database.ResourceReceive()
	go database.EventlogReceive()
	go database.MetricInsert()
	go database.UpdateClusterStatusinfo()

	common.Once.Do(func() {
		kubeapi.SetClientInfo(common.KubeClusterCount, clientset, conf)
	})

	for common.ActiveFlag {
		time.Sleep(time.Second * time.Duration(1))
	}
}
