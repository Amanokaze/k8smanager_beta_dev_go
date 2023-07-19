package common

import (
	"onTuneKubeManager/logger"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestInfo struct {
	Kind      string
	Host      string
	Name      string
	Namespace string
}

var EndChan chan string = make(chan string)

// var cancelChan chan os.Signal = make(chan os.Signal)
var DBConnectionPool *pgxpool.Pool
var ChannelClusterStatus chan map[string]bool = make(chan map[string]bool)
var ChannelResourceData chan map[string][]byte = make(chan map[string][]byte)
var ChannelEventlogData chan map[string][]byte = make(chan map[string][]byte)
var ChannelMetricData chan []byte = make(chan []byte)
var ChannelRequestChangeHost chan string = make(chan string)

var ManagerID int
var ClusterID map[string]int = make(map[string]int)
var LogManager *logger.OnTuneLogManager
var ActiveFlag = true
var TerminatedCount = 0
var TemporaryMultiply = 10000
var MutexNode = &sync.Mutex{}
var MutexPod = &sync.Mutex{}
var MutexNs = &sync.Mutex{}
var Once = &sync.Once{}

var ResourceMap *sync.Map = &sync.Map{}
var PodListMap *sync.Map = &sync.Map{}
var NodeStatusMap *sync.Map = &sync.Map{}
var ClusterStatusMap *sync.Map = &sync.Map{}

var ShorttermDuration int
var LongtermDuration int
var RealtimeInterval int
var ResourceInterval int
var RateInterval int
var EventlogInterval int
var InitVacuum bool
var AutoVacuum bool

var TableSpace bool
var TableSpaceName string
var TableSpacePath string
var DbOsUser string
var ShorttermTableSpaceName string
var ShorttermTableSpace bool
var ShorttermTableSpacePath string
var LongtermTableSpace bool
var LongtermTableSpaceName string
var LongtermTableSpacePath string
var DisconnectFailCount uint16
var KubeClusterCount uint16

// MapInit : init sync.Map
func MapInit() {
	ResourceMap.Store("namespace", &sync.Map{})
	ResourceMap.Store("node", &sync.Map{})
	ResourceMap.Store("pod", &sync.Map{})
	ResourceMap.Store("service", &sync.Map{})
	ResourceMap.Store("persistentvolumeclaim", &sync.Map{})
	ResourceMap.Store("persistentvolume", &sync.Map{})
	ResourceMap.Store("deployment", &sync.Map{})
	ResourceMap.Store("statefulset", &sync.Map{})
	ResourceMap.Store("daemonset", &sync.Map{})
	ResourceMap.Store("replicaset", &sync.Map{})
	ResourceMap.Store("ingress", &sync.Map{})
	ResourceMap.Store("storageclass", &sync.Map{})
	ResourceMap.Store("node_cpu", &sync.Map{})
	ResourceMap.Store("node_memory", &sync.Map{})
	ResourceMap.Store("node_cluster", &sync.Map{})
	ResourceMap.Store("pod_status", &sync.Map{})
	ResourceMap.Store("container_status", &sync.Map{})
	ResourceMap.Store("container_cpu", &sync.Map{})
	ResourceMap.Store("container_memory", &sync.Map{})
	ResourceMap.Store("replicaset_label", &sync.Map{})
	ResourceMap.Store("deployment_label", &sync.Map{})
	ResourceMap.Store("statefulset_label", &sync.Map{})
	ResourceMap.Store("daemonset_label", &sync.Map{})
	ResourceMap.Store("pod_label", &sync.Map{})
	ResourceMap.Store("service_label", &sync.Map{})
	ResourceMap.Store("ingress_label", &sync.Map{})
	ResourceMap.Store("persistentvolumeclaim_label", &sync.Map{})
	ResourceMap.Store("persistentvolume_label", &sync.Map{})
	ResourceMap.Store("storageclass_label", &sync.Map{})
	ResourceMap.Store("node_label", &sync.Map{})
	ResourceMap.Store("namespace_label", &sync.Map{})
	ResourceMap.Store("replicaset_selector", &sync.Map{})
	ResourceMap.Store("deployment_selector", &sync.Map{})
	ResourceMap.Store("statefulset_selector", &sync.Map{})
	ResourceMap.Store("daemonset_selector", &sync.Map{})
	ResourceMap.Store("service_selector", &sync.Map{})
	ResourceMap.Store("pod_reference", &sync.Map{})
	ResourceMap.Store("replicaset_reference", &sync.Map{})
}
