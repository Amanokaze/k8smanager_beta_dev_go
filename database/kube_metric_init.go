package database

import (
	"fmt"
	"onTuneKubeManager/common"
	"sync"
	"time"
)

const (
	NODE_ROOT_METRIC_ID = "/"
	DAY_SECONDS         = 86400
	HOUR_SECONDS        = 3600
	DEFALUT_AVGINTERVAL = 600
)

type MapUIDInfo struct {
	currentNodeUid      map[string]struct{}
	currentPodUid       map[string]struct{}
	currentNamespaceUid map[string]struct{}
	prevNodeUid         map[string]struct{}
	prevPodUid          map[string]struct{}
	prevNamespaceUid    map[string]struct{}
}

type BasicStatMap struct {
	Node        map[string]NodePerfStat
	Pod         map[string]PodPerfStat
	Container   map[string]ContainerPerfStat
	NodeNet     map[string]NodeNetPerfStat
	PodNet      map[string]PodNetPerfStat
	NodeFs      map[string]NodeFsPerfStat
	PodFs       map[string]PodFsPerfStat
	ContainerFs map[string]ContainerFsPerfStat
}

type BasicRawMap struct {
	Node        map[string]NodePerf
	Pod         map[string]PodPerf
	Container   map[string]ContainerPerf
	NodeNet     map[string]NodeNetPerf
	PodNet      map[string]PodNetPerf
	NodeFs      map[string]NodeFsPerf
	PodFs       map[string]PodFsPerf
	ContainerFs map[string]ContainerFsPerf
}

var (
	mapFsDeviceInfo       map[string]int        = make(map[string]int)
	mapNetInterfaceInfo   map[string]int        = make(map[string]int)
	mapMetricIdInfo       map[string]int        = make(map[string]int)
	mapMetricImageInfo    map[int]string        = make(map[int]string)
	mapUidInfo            map[string]MapUIDInfo = make(map[string]MapUIDInfo)
	fs_metric_names       []string              = []string{"rbt", "wbt"}
	net_metric_names      []string              = []string{"rbt", "rbe", "tbt", "tbe"}
	cluster_map           *sync.Map             = &sync.Map{}
	last_stat_data_map    *sync.Map             = &sync.Map{}
	last_summary_data_map *sync.Map             = &sync.Map{}
)

var Channelmetric_insert chan map[string]interface{} = make(chan map[string]interface{})

func MetricDataReceive() {
	for {
		metric_msg := <-common.ChannelMetricData
		common.LogManager.Debug("mapping metric data receive complete")
		RealtimeMetricDataParsing(metric_msg)
		time.Sleep(time.Millisecond * 10)
	}
}

func InitMetricDataMap() {
	// fsdevice info
	var rowsCnt int
	rowsCnt = selectRowCount(TB_KUBE_FS_DEVICE_INFO)
	if rowsCnt > 0 {
		if len(mapFsDeviceInfo) == 0 {
			var id int
			var name string
			rows := selectCondition("deviceid, devicename", TB_KUBE_FS_DEVICE_INFO, "devicename", "!=", "")
			for rows.Next() {
				err := rows.Scan(&id, &name)
				if !errorCheck(err) {
					return
				}
				mapFsDeviceInfo[name] = id
			}
		}
	}

	// netinterface info
	rowsCnt = selectRowCount(TB_KUBE_NET_INTERFACE_INFO)
	if rowsCnt > 0 {
		if len(mapNetInterfaceInfo) == 0 {
			var id int
			var name string
			rows := selectCondition("interfaceid, interfacename", TB_KUBE_NET_INTERFACE_INFO, "interfacename", "!=", "")
			for rows.Next() {
				err := rows.Scan(&id, &name)
				if !errorCheck(err) {
					return
				}
				mapNetInterfaceInfo[name] = id
			}
		}
	}

	// metricid info
	rowsCnt = selectRowCount(TB_KUBE_METRIC_ID_INFO)
	if rowsCnt > 0 {
		if len(mapMetricIdInfo) == 0 {
			var id int
			var name string
			var image string
			rows := selectCondition("metricid, metricname, image", TB_KUBE_METRIC_ID_INFO, "metricname", "!=", "")
			for rows.Next() {
				err := rows.Scan(&id, &name, &image)
				if !errorCheck(err) {
					return
				}
				mapMetricIdInfo[name] = id
				mapMetricImageInfo[id] = image
			}
		}
	}
}

func GenerateLongTermData() {
	if HOUR_SECONDS%common.AvgInterval != 0 {
		common.LogManager.Error("AvgInterval is not a divisor of 3600. In this execution, AvgInterval is 600")
		common.AvgInterval = DEFALUT_AVGINTERVAL
	} else if common.AvgInterval < common.RealtimeInterval {
		common.LogManager.Error("AvgInterval is smaller than RealtimeInterval. In this execution, AvgInterval is 600")
		common.AvgInterval = DEFALUT_AVGINTERVAL
	}

	_, bias := GetOntuneTime()

	for {
		ontunetime := time.Now().Unix() - bias

		// 이전 시간과 비교하는 개념이 아니라, AvgInterval 단위로 정확하게 입력하도록 함
		if ontunetime%int64(common.AvgInterval) == 0 {
			// biasSelect Query는 최초 1회, 그리고 Interval 내 조건을 충족할 때에만 수행하는 것으로 함
			_, bias = GetOntuneTime()
			common.LogManager.WriteLog(fmt.Sprintf("GenerateLongTermData Start: %d", ontunetime))

			AvgMetricDataInsert(ontunetime, bias)
		}
		time.Sleep(time.Second * 1)
	}
}
