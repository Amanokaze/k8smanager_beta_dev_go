package database

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"onTuneKubeManager/common"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

const (
	NODE_ROOT_METRIC_ID = "/"
)

type MapUIDInfo struct {
	currentNodeUid      map[string]struct{}
	currentPodUid       map[string]struct{}
	currentNamespaceUid map[string]struct{}
	prevNodeUid         map[string]struct{}
	prevPodUid          map[string]struct{}
	prevNamespaceUid    map[string]struct{}
}

var (
	mapFsDeviceInfo     map[string]int        = make(map[string]int)
	mapNetInterfaceInfo map[string]int        = make(map[string]int)
	mapMetricIdInfo     map[string]int        = make(map[string]int)
	mapMetricImageInfo  map[int]string        = make(map[int]string)
	mapUidInfo          map[string]MapUIDInfo = make(map[string]MapUIDInfo)
	// mapCurrentNodeUIDInfo      map[string]struct{} = make(map[string]struct{})
	// mapCurrentPodUIDInfo       map[string]struct{} = make(map[string]struct{})
	// mapCurrentNamespaceUIDInfo map[string]struct{} = make(map[string]struct{})
	// mapPrevNodeUIDInfo         map[string]struct{} = make(map[string]struct{})
	// mapPrevPodUIDInfo          map[string]struct{} = make(map[string]struct{})
	// mapPrevNamespaceUIDInfo    map[string]struct{} = make(map[string]struct{})
	fs_metric_names       []string  = []string{"rbt", "wbt"}
	net_metric_names      []string  = []string{"rbt", "rbe", "tbt", "tbe"}
	cluster_map           *sync.Map = &sync.Map{}
	last_stat_data_map    *sync.Map = &sync.Map{}
	last_summary_data_map *sync.Map = &sync.Map{}
)

var Channelmetric_insert chan map[string]interface{} = make(chan map[string]interface{})

func MetricInsert() {
	for {
		metric_msg := <-common.ChannelMetricData
		writeLog("mapping metric data receive complete")
		dataParsing(metric_msg)
		time.Sleep(time.Millisecond * 10)
	}
}

func InitMapMetricData() {
	// fsdevice info
	var rowsCnt int
	rowsCnt = select_row_count(TB_KUBE_FS_DEVICE_INFO)
	if rowsCnt > 0 {
		if len(mapFsDeviceInfo) == 0 {
			var id int
			var name string
			rows := select_condition("deviceid, devicename", TB_KUBE_FS_DEVICE_INFO, "devicename", "!=", "")
			for rows.Next() {
				err := rows.Scan(&id, &name)
				errorCheck(err)
				mapFsDeviceInfo[name] = id
			}
		}
	}

	// netinterface info
	rowsCnt = select_row_count(TB_KUBE_NET_INTERFACE_INFO)
	if rowsCnt > 0 {
		if len(mapNetInterfaceInfo) == 0 {
			var id int
			var name string
			rows := select_condition("interfaceid, interfacename", TB_KUBE_NET_INTERFACE_INFO, "interfacename", "!=", "")
			for rows.Next() {
				err := rows.Scan(&id, &name)
				errorCheck(err)
				mapNetInterfaceInfo[name] = id
			}
		}
	}

	// metricid info
	rowsCnt = select_row_count(TB_KUBE_METRIC_ID_INFO)
	if rowsCnt > 0 {
		if len(mapMetricIdInfo) == 0 {
			var id int
			var name string
			var image string
			rows := select_condition("metricid, metricname, image", TB_KUBE_METRIC_ID_INFO, "metricname", "!=", "")
			for rows.Next() {
				err := rows.Scan(&id, &name, &image)
				errorCheck(err)
				mapMetricIdInfo[name] = id
				mapMetricImageInfo[id] = image
			}
		}
	}
}

func updateFsDeviceId(name string) {
	ontunetime, _ := GetOntuneTime()
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	_, insert_err := tx.Exec(context.Background(), INSERT_FSDEVICE_INFO, name, ontunetime, ontunetime)
	errorCheck(insert_err)

	err = tx.Commit(context.Background())
	errorCheck(err)

	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	idvalue := select_one_value_condition("deviceid", TB_KUBE_FS_DEVICE_INFO, "devicename", name)
	idvalue_num, err := strconv.Atoi(idvalue)
	errorCheck(err)

	mapFsDeviceInfo[name] = idvalue_num

	err = tx.Commit(context.Background())
	errorCheck(err)
}

func updateNetInterfaceId(name string) {
	ontunetime, _ := GetOntuneTime()
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	_, insert_err := tx.Exec(context.Background(), INSERT_NETINTERFACE_INFO, name, ontunetime, ontunetime)
	errorCheck(insert_err)

	err = tx.Commit(context.Background())
	errorCheck(err)

	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	idvalue := select_one_value_condition("interfaceid", TB_KUBE_NET_INTERFACE_INFO, "interfacename", name)
	idvalue_num, err := strconv.Atoi(idvalue)
	errorCheck(err)

	mapNetInterfaceInfo[name] = idvalue_num

	err = tx.Commit(context.Background())
	errorCheck(err)
}

func updateMetricId(name string, image string) {
	ontunetime, _ := GetOntuneTime()
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	_, insert_err := tx.Exec(context.Background(), INSERT_METRICID_INFO, name, image, ontunetime, ontunetime)
	errorCheck(insert_err)

	err = tx.Commit(context.Background())
	errorCheck(err)

	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	idvalue := select_one_value_condition("metricid", TB_KUBE_METRIC_ID_INFO, "metricname", name)
	idvalue_num, err := strconv.Atoi(idvalue)
	errorCheck(err)

	mapMetricIdInfo[name] = idvalue_num
	mapMetricImageInfo[idvalue_num] = image

	err = tx.Commit(context.Background())
	errorCheck(err)
}

func dataParsing(metric_data []byte) {
	var data interface{}

	err := json.Unmarshal([]byte(metric_data), &data)
	errorCheck(err)

	map_datas := data.(map[string]interface{})
	insert_metric_data := make(map[string]interface{})
	insert_metric_data["node"] = map_datas["node"]
	insert_metric_data["host"] = map_datas["host"]

	var flag bool = false
	if _, ok := map_datas["nodes"]; ok {
		mapflag := setNodeMetric(map_datas["nodes"], &insert_metric_data)
		flag = flag || mapflag
	}

	if _, ok := map_datas["pods"]; ok {
		mapflag := setPodMetric(map_datas["pods"], &insert_metric_data)
		flag = flag || mapflag
	}

	if _, ok := map_datas["containers"]; ok {
		mapflag := setContainerMetric(map_datas["containers"], &insert_metric_data)
		flag = flag || mapflag
	}

	if _, ok := map_datas["nodenetworks"]; ok {
		mapflag := setNodeNetMetric(map_datas["nodenetworks"], &insert_metric_data)
		flag = flag || mapflag
	}

	if _, ok := map_datas["podnetworks"]; ok {
		mapflag := setPodNetMetric(map_datas["podnetworks"], &insert_metric_data)
		flag = flag || mapflag
	}

	if _, ok := map_datas["nodefilesystems"]; ok {
		mapflag := setNodeFsMetric(map_datas["nodefilesystems"], &insert_metric_data)
		flag = flag || mapflag
	}

	if _, ok := map_datas["podfilesystems"]; ok {
		mapflag := setPodFsMetric(map_datas["podfilesystems"], &insert_metric_data)
		flag = flag || mapflag
	}

	if _, ok := map_datas["containerfilesystems"]; ok {
		mapflag := setContainerFsMetric(map_datas["containerfilesystems"], &insert_metric_data)
		flag = flag || mapflag
	}

	if flag {
		common.ChannelRequestChangeHost <- insert_metric_data["host"].(string)
	}

	Channelmetric_insert <- insert_metric_data
}

func setNodeMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var nodeInsert NodePerf
	var flag bool = false

	for _, node_data := range map_data.([]interface{}) {
		node_data_map := node_data.(map[string]interface{})
		nodeuid := getUID("node", (*map_metric)["host"], node_data_map["node"])
		image := node_data.(map[string]interface{})["image"]

		// getUID()가 Null이라고 즉시 불러오지는 말고, 아래와 같이 별도 로직을 호출해서 불러오는게 좋을 것 같습니다.
		if nodeuid == "" {
			flag = true
			continue
		}

		for item, data := range node_data_map {
			switch item {
			case METRIC_VAR_NODE:
				nodeInsert.ArrNodeuid = append(nodeInsert.ArrNodeuid, nodeuid)
			case METRIC_VAR_METRICID:
				nodeInsert.ArrMetricid = append(nodeInsert.ArrMetricid, getMetricID(data, image))
			case METRIC_VAR_TIMESTAMPMS:
				f_data := data.(float64)
				nodeInsert.ArrTimestampMs = append(nodeInsert.ArrTimestampMs, int64(f_data))
			case METRIC_VAR_CPUUSAGESEONDSTOTAL:
				nodeInsert.ArrCpuusagesecondstotal = append(nodeInsert.ArrCpuusagesecondstotal, data.(float64))
			case METRIC_VAR_CPUSYSTEMSECONDSTOTAL:
				nodeInsert.ArrCpusystemsecondstotal = append(nodeInsert.ArrCpusystemsecondstotal, data.(float64))
			case METRIC_VAR_CPUUSERSECONDSTOTAL:
				nodeInsert.ArrCpuusersecondstotal = append(nodeInsert.ArrCpuusersecondstotal, data.(float64))
			case METRIC_VAR_MEMORYUSAGEBYTES:
				nodeInsert.ArrMemoryusagebytes = append(nodeInsert.ArrMemoryusagebytes, data.(float64))
			case METRIC_VAR_MEMORYWORKINGSETBYTES:
				nodeInsert.ArrMemoryworkingsetbytes = append(nodeInsert.ArrMemoryworkingsetbytes, data.(float64))
			case METRIC_VAR_MEMORYCACHE:
				nodeInsert.ArrMemorycache = append(nodeInsert.ArrMemorycache, data.(float64))
			case METRIC_VAR_MEMORYSWAP:
				nodeInsert.ArrMemoryswap = append(nodeInsert.ArrMemoryswap, data.(float64))
			case METRIC_VAR_MEMORYRSS:
				nodeInsert.ArrMemoryrss = append(nodeInsert.ArrMemoryrss, data.(float64))
			case METRIC_VAR_FSREADSBYTESTOTAL:
				nodeInsert.ArrFsreadsbytestotal = append(nodeInsert.ArrFsreadsbytestotal, data.(float64))
			case METRIC_VAR_FSWRITESBYTESTOTAL:
				nodeInsert.ArrFswritesbytestotal = append(nodeInsert.ArrFswritesbytestotal, data.(float64))
			case METRIC_VAR_PROCESSES:
				nodeInsert.ArrProcesses = append(nodeInsert.ArrProcesses, data.(float64))
			}
		}
		nodeInsert.ArrOntunetime = append(nodeInsert.ArrOntunetime, 0)
	}
	var m_data interface{} = nodeInsert
	(*map_metric)["node_insert"] = m_data

	return flag
}

func insertNodePerf(insert_data NodePerf, ontunetime int64, nodeuid string) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	perftablename := getTableName(TB_KUBE_NODE_PERF_RAW)

	for i := 0; i < len(insert_data.ArrOntunetime); i++ {
		insert_data.ArrOntunetime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODE_PERF, perftablename), pq.StringArray(insert_data.ArrNodeuid), pq.Int64Array(insert_data.ArrOntunetime),
		pq.Array(insert_data.ArrMetricid), pq.Float64Array(insert_data.ArrCpuusagesecondstotal), pq.Float64Array(insert_data.ArrCpusystemsecondstotal),
		pq.Float64Array(insert_data.ArrCpuusersecondstotal), pq.Float64Array(insert_data.ArrMemoryusagebytes), pq.Float64Array(insert_data.ArrMemoryworkingsetbytes),
		pq.Float64Array(insert_data.ArrMemorycache), pq.Float64Array(insert_data.ArrMemoryswap), pq.Float64Array(insert_data.ArrMemoryrss), pq.Float64Array(insert_data.ArrFsreadsbytestotal),
		pq.Float64Array(insert_data.ArrFswritesbytestotal), pq.Float64Array(insert_data.ArrProcesses), pq.Int64Array(insert_data.ArrTimestampMs))
	errorCheck(err)

	err = tx.Commit(context.Background())
	errorCheck(err)

	conn.Release()

	update_tableinfo(perftablename, ontunetime)
}

func insertLastNodePerf(insert_data NodePerf, ontunetime int64, nodeuid string) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA, TB_KUBE_LAST_NODE_PERF_RAW, nodeuid, ontunetime))
	errorCheck(del_err)

	for i := 0; i < len(insert_data.ArrOntunetime); i++ {
		insert_data.ArrOntunetime[i] = ontunetime
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_LAST_NODE_PERF_RAW, pq.StringArray(insert_data.ArrNodeuid), pq.Int64Array(insert_data.ArrOntunetime),
		pq.Array(insert_data.ArrMetricid), pq.Float64Array(insert_data.ArrCpuusagesecondstotal), pq.Float64Array(insert_data.ArrCpusystemsecondstotal),
		pq.Float64Array(insert_data.ArrCpuusersecondstotal), pq.Float64Array(insert_data.ArrMemoryusagebytes), pq.Float64Array(insert_data.ArrMemoryworkingsetbytes),
		pq.Float64Array(insert_data.ArrMemorycache), pq.Float64Array(insert_data.ArrMemoryswap), pq.Float64Array(insert_data.ArrMemoryrss), pq.Float64Array(insert_data.ArrFsreadsbytestotal),
		pq.Float64Array(insert_data.ArrFswritesbytestotal), pq.Float64Array(insert_data.ArrProcesses), pq.Int64Array(insert_data.ArrTimestampMs))
	errorCheck(err)

	err = tx.Commit(context.Background())
	errorCheck(err)

	conn.Release()

	update_tableinfo(TB_KUBE_LAST_NODE_PERF_RAW, ontunetime)
}

func insertNodeStat(current_data *sync.Map, prev_data *sync.Map, net_map *sync.Map, fs_map *sync.Map, nodeuid string) *NodePerfStat {
	_, result := getStatTimes(current_data)
	if !result {
		return &NodePerfStat{}
	}

	node_stat_metric_id := getMetricID(NODE_ROOT_METRIC_ID, "")

	var current_node_data NodePerf
	var prev_node_data NodePerf

	if c, ok := current_data.Load("node_insert"); ok {
		current_node_data = c.(NodePerf)
	} else {
		return &NodePerfStat{}
	}

	if p, ok := prev_data.Load("node_insert"); ok {
		prev_node_data = p.(NodePerf)
	} else {
		return &NodePerfStat{}
	}

	insert_stat_data := new(NodePerfStat)
	insert_stat_data.Init()

	for i := 0; i < len(current_node_data.ArrOntunetime); i++ {
		if current_node_data.ArrMetricid[i] == node_stat_metric_id {
			var pidx int = -1
			for j := 0; j < len(prev_node_data.ArrOntunetime); j++ {
				if prev_node_data.ArrMetricid[j] == node_stat_metric_id {
					pidx = j
					break
				}
			}
			if pidx == -1 {
				break
			}

			node_cpu_map, _ := common.ResourceMap.Load("node_cpu")
			node_memory_map, _ := common.ResourceMap.Load("node_memory")
			cc, _ := node_cpu_map.(*sync.Map).LoadOrStore(current_node_data.ArrNodeuid[i], 0)
			ms, _ := node_memory_map.(*sync.Map).LoadOrStore(current_node_data.ArrNodeuid[i], int64(0))
			cpucount := cc.(int)
			memorysize := ms.(int64)

			insert_stat_data.ArrNodeuid = append(insert_stat_data.ArrNodeuid, current_node_data.ArrNodeuid[i])
			insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_node_data.ArrOntunetime[i])
			insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, current_node_data.ArrMetricid[i])
			insert_stat_data.ArrProcesses = append(insert_stat_data.ArrProcesses, current_node_data.ArrProcesses[i])

			summary_key := SummaryIndex{
				NodeUID: current_node_data.ArrNodeuid[i],
			}

			current_net_rbt := GetSummaryData(net_map, summary_key, "current", "rbt")
			current_net_rbe := GetSummaryData(net_map, summary_key, "current", "rbe")
			current_net_tbt := GetSummaryData(net_map, summary_key, "current", "tbt")
			current_net_tbe := GetSummaryData(net_map, summary_key, "current", "tbe")
			previous_net_rbt := GetSummaryData(net_map, summary_key, "previous", "rbt")
			previous_net_rbe := GetSummaryData(net_map, summary_key, "previous", "rbe")
			previous_net_tbt := GetSummaryData(net_map, summary_key, "previous", "tbt")
			previous_net_tbe := GetSummaryData(net_map, summary_key, "previous", "tbe")
			rcv_rate := GetStatData("rate", current_net_rbt, previous_net_rbt, 0)
			rcv_errors := GetStatData("subtract", current_net_rbe, previous_net_rbe, 0)
			trans_rate := GetStatData("rate", current_net_tbt, previous_net_tbt, 0)
			trans_errors := GetStatData("subtract", current_net_tbe, previous_net_tbe, 0)

			insert_stat_data.ArrNetworkreceiverate = append(insert_stat_data.ArrNetworkreceiverate, rcv_rate)
			insert_stat_data.ArrNetworkreceiveerrors = append(insert_stat_data.ArrNetworkreceiveerrors, rcv_errors)
			insert_stat_data.ArrNetworktransmitrate = append(insert_stat_data.ArrNetworktransmitrate, trans_rate)
			insert_stat_data.ArrNetworktransmiterrors = append(insert_stat_data.ArrNetworktransmiterrors, trans_errors)
			insert_stat_data.ArrNetworkiorate = append(insert_stat_data.ArrNetworkiorate, rcv_rate+trans_rate)
			insert_stat_data.ArrNetworkioerrors = append(insert_stat_data.ArrNetworkioerrors, rcv_errors+trans_errors)

			current_fs_rbt := GetSummaryData(fs_map, summary_key, "current", "rbt")
			current_fs_wbt := GetSummaryData(fs_map, summary_key, "current", "wbt")
			previous_fs_rbt := GetSummaryData(fs_map, summary_key, "previous", "rbt")
			previous_fs_wbt := GetSummaryData(fs_map, summary_key, "previous", "wbt")
			fs_read_rate := GetStatData("rate", current_fs_rbt, previous_fs_rbt, 0)
			fs_write_rate := GetStatData("rate", current_fs_wbt, previous_fs_wbt, 0)

			insert_stat_data.ArrFsreadrate = append(insert_stat_data.ArrFsreadrate, fs_read_rate)
			insert_stat_data.ArrFswriterate = append(insert_stat_data.ArrFswriterate, fs_write_rate)
			insert_stat_data.ArrFsiorate = append(insert_stat_data.ArrFsiorate, fs_read_rate+fs_write_rate)

			insert_stat_data.ArrCpuusage = append(insert_stat_data.ArrCpuusage, GetStatData("usage", current_node_data.ArrCpuusagesecondstotal[i], prev_node_data.ArrCpuusagesecondstotal[pidx], float64(cpucount)))
			insert_stat_data.ArrCpusystem = append(insert_stat_data.ArrCpusystem, GetStatData("usage", current_node_data.ArrCpusystemsecondstotal[i], prev_node_data.ArrCpusystemsecondstotal[pidx], float64(cpucount)))
			insert_stat_data.ArrCpuuser = append(insert_stat_data.ArrCpuuser, GetStatData("usage", current_node_data.ArrCpuusersecondstotal[i], prev_node_data.ArrCpuusersecondstotal[pidx], float64(cpucount)))
			insert_stat_data.ArrCpuusagecores = append(insert_stat_data.ArrCpuusagecores, GetStatData("core", current_node_data.ArrCpuusagesecondstotal[i], prev_node_data.ArrCpuusagesecondstotal[pidx], 0))
			insert_stat_data.ArrCputotalcores = append(insert_stat_data.ArrCputotalcores, float64(cpucount))
			insert_stat_data.ArrMemoryusage = append(insert_stat_data.ArrMemoryusage, GetStatData("current_usage", current_node_data.ArrMemoryworkingsetbytes[i], 0, float64(memorysize)))
			insert_stat_data.ArrMemoryusagebytes = append(insert_stat_data.ArrMemoryusagebytes, math.Round(current_node_data.ArrMemoryworkingsetbytes[i]))
			insert_stat_data.ArrMemorysizebytes = append(insert_stat_data.ArrMemorysizebytes, math.Round(float64(memorysize)))
			insert_stat_data.ArrMemoryswap = append(insert_stat_data.ArrMemoryswap, GetStatData("current_usage", current_node_data.ArrMemoryswap[i], 0, float64(memorysize)))

			break
		}
	}

	lsdm, _ := last_stat_data_map.LoadOrStore("node", &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store("node", lsdm)

	return insert_stat_data
}

func insertNodeStatTable(nodeuid string, insert_stat_data NodePerfStat, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	tablename := getTableName(TB_KUBE_NODE_PERF)
	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODE_STAT, tablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	lasttablename := TB_KUBE_LAST_NODE_PERF
	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA, lasttablename, nodeuid, ontunetime))
	errorCheck(del_err)

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODE_STAT, lasttablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	err = tx.Commit(context.Background())
	errorCheck(err)

	conn.Release()

	update_tableinfo(tablename, ontunetime)
	update_tableinfo(lasttablename, ontunetime)
}

func setPodMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var arrInsert PodPerf
	var flag bool = false

	for _, pod_data := range map_data.([]interface{}) {
		pod_data_map := pod_data.(map[string]interface{})

		nsuid := getUID("namespace", (*map_metric)["host"], pod_data_map["namespace"])
		poduid := getUID("pod", pod_data_map["namespace"].(string), pod_data_map["pod"].(string))
		nodeuid := getUID("node", (*map_metric)["host"], pod_data_map["node"])
		image := pod_data_map["image"]

		if nodeuid == "" {
			flag = true
			continue
		}

		if nsuid == "" {
			flag = true
			continue
		}

		if poduid == "" {
			flag = true
			continue
		}

		for item, data := range pod_data_map {
			switch item {
			case METRIC_VAR_NAMESPACE:
				arrInsert.ArrNsuid = append(arrInsert.ArrNsuid, nsuid)
			case METRIC_VAR_POD:
				arrInsert.ArrPoduid = append(arrInsert.ArrPoduid, poduid)
			case METRIC_VAR_NODE:
				arrInsert.ArrNodeuid = append(arrInsert.ArrNodeuid, nodeuid)
			case METRIC_VAR_METRICID:
				arrInsert.ArrMetricid = append(arrInsert.ArrMetricid, getMetricID(data, image))
			case METRIC_VAR_TIMESTAMPMS:
				f_data := data.(float64)
				arrInsert.ArrTimestampMs = append(arrInsert.ArrTimestampMs, int64(f_data))
			case METRIC_VAR_CPUUSAGESEONDSTOTAL:
				arrInsert.ArrCpuusagesecondstotal = append(arrInsert.ArrCpuusagesecondstotal, data.(float64))
			case METRIC_VAR_CPUSYSTEMSECONDSTOTAL:
				arrInsert.ArrCpusystemsecondstotal = append(arrInsert.ArrCpusystemsecondstotal, data.(float64))
			case METRIC_VAR_CPUUSERSECONDSTOTAL:
				arrInsert.ArrCpuusersecondstotal = append(arrInsert.ArrCpuusersecondstotal, data.(float64))
			case METRIC_VAR_MEMORYUSAGEBYTES:
				arrInsert.ArrMemoryusagebytes = append(arrInsert.ArrMemoryusagebytes, data.(float64))
			case METRIC_VAR_MEMORYWORKINGSETBYTES:
				arrInsert.ArrMemoryworkingsetbytes = append(arrInsert.ArrMemoryworkingsetbytes, data.(float64))
			case METRIC_VAR_MEMORYCACHE:
				arrInsert.ArrMemorycache = append(arrInsert.ArrMemorycache, data.(float64))
			case METRIC_VAR_MEMORYSWAP:
				arrInsert.ArrMemoryswap = append(arrInsert.ArrMemoryswap, data.(float64))
			case METRIC_VAR_MEMORYRSS:
				arrInsert.ArrMemoryrss = append(arrInsert.ArrMemoryrss, data.(float64))
			case METRIC_VAR_PROCESSES:
				arrInsert.ArrProcesses = append(arrInsert.ArrProcesses, data.(float64))
			}
		}
		arrInsert.ArrOntunetime = append(arrInsert.ArrOntunetime, 0)
	}
	var m_data interface{} = arrInsert
	(*map_metric)["pod_insert"] = m_data

	return flag
}

func insertPodPerf(insert_data PodPerf, ontunetime int64, nodeuid string) {
	perftablename := getTableName(TB_KUBE_POD_PERF_RAW)

	for i := 0; i < len(insert_data.ArrOntunetime); i++ {
		insert_data.ArrOntunetime[i] = ontunetime
	}

	if len(insert_data.ArrCpuusagesecondstotal) > 0 {
		conn, err := common.DBConnectionPool.Acquire(context.Background())
		if err != nil {
			errorCheck(err)
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_POD_PERF, perftablename), insert_data.GetArgs()...)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)

		conn.Release()

		update_tableinfo(perftablename, ontunetime)
	}
}

func insertLastPodPerf(insert_data PodPerf, ontunetime int64, nodeuid string) {
	for i := 0; i < len(insert_data.ArrOntunetime); i++ {
		insert_data.ArrOntunetime[i] = ontunetime
	}

	if len(insert_data.ArrCpuusagesecondstotal) > 0 {
		conn, err := common.DBConnectionPool.Acquire(context.Background())
		if err != nil {
			errorCheck(err)
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, del_err := tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA, TB_KUBE_LAST_POD_PERF_RAW, nodeuid, ontunetime))
		errorCheck(del_err)

		_, err = tx.Exec(context.Background(), INSERT_UNNEST_LAST_POD_PERF_RAW, insert_data.GetArgs()...)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)

		conn.Release()

		update_tableinfo(TB_KUBE_LAST_POD_PERF_RAW, ontunetime)
	}
}

func insertPodStat(current_data *sync.Map, prev_data *sync.Map, net_map *sync.Map, fs_map *sync.Map, nodeuid string) (*PodPerfStat, *sync.Map) {
	_, result := getStatTimes(current_data)
	if !result {
		return &PodPerfStat{}, &sync.Map{}
	}

	var current_pod_summary_data PodSummaryPerf
	var prev_pod_summary_data PodSummaryPerf
	var summary_data *sync.Map = &sync.Map{}

	if c, ok := current_data.Load("pod_insert"); ok {
		current_pod_data := c.(PodPerf)

		var current_container_data ContainerPerf
		if cc, ok2 := current_data.Load("container_insert"); ok2 {
			current_container_data = cc.(ContainerPerf)
		}
		current_pod_summary_data = PodSummaryPreprocessing(&current_pod_data, &current_container_data)
	} else {
		return &PodPerfStat{}, &sync.Map{}
	}

	if p, ok := prev_data.Load("pod_insert"); ok {
		prev_pod_data := p.(PodPerf)

		var prev_container_data ContainerPerf
		if pc, ok2 := prev_data.Load("container_insert"); ok2 {
			prev_container_data = pc.(ContainerPerf)
		}
		prev_pod_summary_data = PodSummaryPreprocessing(&prev_pod_data, &prev_container_data)
	} else {
		return &PodPerfStat{}, &sync.Map{}
	}

	insert_stat_data := new(PodPerfStat)
	insert_stat_data.Init()

	for i := 0; i < len(current_pod_summary_data.ArrOntunetime); i++ {
		var pidx int = -1
		for j := 0; j < len(prev_pod_summary_data.ArrOntunetime); j++ {
			if prev_pod_summary_data.ArrPoduid[j] == current_pod_summary_data.ArrPoduid[i] {
				pidx = j
				break
			}
		}
		if pidx == -1 {
			continue
		}
		pod_status_map, _ := common.ResourceMap.Load("pod_status")
		if pval, ok := pod_status_map.(*sync.Map).Load(current_pod_summary_data.ArrPoduid[i]); !ok || pval.(string) != POD_STATUS_RUNNING {
			continue
		}

		summary_key := SummaryIndex{
			NsUID:  current_pod_summary_data.ArrNsuid[i],
			PodUID: current_pod_summary_data.ArrPoduid[i],
		}

		insert_stat_data.ArrPoduid = append(insert_stat_data.ArrPoduid, current_pod_summary_data.ArrPoduid[i])
		insert_stat_data.ArrNodeuid = append(insert_stat_data.ArrNodeuid, current_pod_summary_data.ArrNodeuid[i])
		insert_stat_data.ArrNsuid = append(insert_stat_data.ArrNsuid, current_pod_summary_data.ArrNsuid[i])
		insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_pod_summary_data.ArrOntunetime[i])
		insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, 0)
		insert_stat_data.ArrProcesses = append(insert_stat_data.ArrProcesses, current_pod_summary_data.ArrProcesses[i])
		SetPodSummary(summary_data, summary_key, "current", "process", current_pod_summary_data.ArrProcesses[i])
		SetPodSummary(summary_data, summary_key, "previous", "process", prev_pod_summary_data.ArrProcesses[pidx])

		insert_stat_data.ArrCpuusage = append(insert_stat_data.ArrCpuusage, GetStatData("usage", current_pod_summary_data.ArrCpuusagesecondstotal[i], prev_pod_summary_data.ArrCpuusagesecondstotal[pidx], current_pod_summary_data.ArrCpucount[i]))
		insert_stat_data.ArrCpusystem = append(insert_stat_data.ArrCpusystem, GetStatData("usage", current_pod_summary_data.ArrCpusystemsecondstotal[i], prev_pod_summary_data.ArrCpusystemsecondstotal[pidx], current_pod_summary_data.ArrCpucount[i]))
		insert_stat_data.ArrCpuuser = append(insert_stat_data.ArrCpuuser, GetStatData("usage", current_pod_summary_data.ArrCpuusersecondstotal[i], prev_pod_summary_data.ArrCpuusersecondstotal[pidx], current_pod_summary_data.ArrCpucount[i]))
		insert_stat_data.ArrCpuusagecores = append(insert_stat_data.ArrCpuusagecores, GetStatData("core", current_pod_summary_data.ArrCpuusagesecondstotal[i], prev_pod_summary_data.ArrCpuusagesecondstotal[pidx], 0))
		insert_stat_data.ArrCpurequestcores = append(insert_stat_data.ArrCpurequestcores, current_pod_summary_data.ArrCpurequest[i])
		insert_stat_data.ArrCpulimitcores = append(insert_stat_data.ArrCpulimitcores, current_pod_summary_data.ArrCpulimit[i])
		insert_stat_data.ArrMemoryusage = append(insert_stat_data.ArrMemoryusage, GetStatData("current_usage", current_pod_summary_data.ArrMemoryworkingsetbytes[i], 0, current_pod_summary_data.ArrMemorysize[i]))
		insert_stat_data.ArrMemoryusagebytes = append(insert_stat_data.ArrMemoryusagebytes, math.Round(current_pod_summary_data.ArrMemoryworkingsetbytes[i]))
		insert_stat_data.ArrMemoryrequestbytes = append(insert_stat_data.ArrMemoryrequestbytes, math.Round(current_pod_summary_data.ArrMemoryrequest[i]))
		insert_stat_data.ArrMemorylimitbytes = append(insert_stat_data.ArrMemorylimitbytes, math.Round(current_pod_summary_data.ArrMemorylimit[i]))
		insert_stat_data.ArrMemoryswap = append(insert_stat_data.ArrMemoryswap, GetStatData("current_usage", current_pod_summary_data.ArrMemoryswap[i], 0, current_pod_summary_data.ArrMemorysize[i]))

		SetPodSummary(summary_data, summary_key, "current", "cpuusage", current_pod_summary_data.ArrCpuusagesecondstotal[i])
		SetPodSummary(summary_data, summary_key, "current", "cpusystem", current_pod_summary_data.ArrCpusystemsecondstotal[i])
		SetPodSummary(summary_data, summary_key, "current", "cpuuser", current_pod_summary_data.ArrCpuusersecondstotal[i])
		SetPodSummary(summary_data, summary_key, "current", "cpucount", current_pod_summary_data.ArrCpucount[i])
		SetPodSummary(summary_data, summary_key, "current", "cpurequest", current_pod_summary_data.ArrCpurequest[i])
		SetPodSummary(summary_data, summary_key, "current", "cpulimit", current_pod_summary_data.ArrCpulimit[i])
		SetPodSummary(summary_data, summary_key, "current", "memoryws", current_pod_summary_data.ArrMemoryworkingsetbytes[i])
		SetPodSummary(summary_data, summary_key, "current", "memoryswap", current_pod_summary_data.ArrMemoryswap[i])
		SetPodSummary(summary_data, summary_key, "current", "memorysize", current_pod_summary_data.ArrMemorysize[i])
		SetPodSummary(summary_data, summary_key, "current", "memoryrequest", current_pod_summary_data.ArrMemoryrequest[i])
		SetPodSummary(summary_data, summary_key, "current", "memorylimit", current_pod_summary_data.ArrMemorylimit[i])

		SetPodSummary(summary_data, summary_key, "previous", "cpuusage", prev_pod_summary_data.ArrCpuusagesecondstotal[pidx])
		SetPodSummary(summary_data, summary_key, "previous", "cpusystem", prev_pod_summary_data.ArrCpusystemsecondstotal[pidx])
		SetPodSummary(summary_data, summary_key, "previous", "cpuuser", prev_pod_summary_data.ArrCpuusersecondstotal[pidx])
		SetPodSummary(summary_data, summary_key, "previous", "memoryws", prev_pod_summary_data.ArrMemoryworkingsetbytes[pidx])
		SetPodSummary(summary_data, summary_key, "previous", "memoryswap", prev_pod_summary_data.ArrMemoryswap[pidx])

		current_net_rbt := GetSummaryData(net_map, summary_key, "current", "rbt")
		current_net_rbe := GetSummaryData(net_map, summary_key, "current", "rbe")
		current_net_tbt := GetSummaryData(net_map, summary_key, "current", "tbt")
		current_net_tbe := GetSummaryData(net_map, summary_key, "current", "tbe")
		previous_net_rbt := GetSummaryData(net_map, summary_key, "previous", "rbt")
		previous_net_rbe := GetSummaryData(net_map, summary_key, "previous", "rbe")
		previous_net_tbt := GetSummaryData(net_map, summary_key, "previous", "tbt")
		previous_net_tbe := GetSummaryData(net_map, summary_key, "previous", "tbe")
		rcv_rate := GetStatData("rate", current_net_rbt, previous_net_rbt, 0)
		rcv_errors := GetStatData("subtract", current_net_rbe, previous_net_rbe, 0)
		trans_rate := GetStatData("rate", current_net_tbt, previous_net_tbt, 0)
		trans_errors := GetStatData("subtract", current_net_tbe, previous_net_tbe, 0)
		SetPodSummary(summary_data, summary_key, "current", "netrbt", current_net_rbt)
		SetPodSummary(summary_data, summary_key, "current", "netrbe", current_net_rbe)
		SetPodSummary(summary_data, summary_key, "current", "nettbt", current_net_tbt)
		SetPodSummary(summary_data, summary_key, "current", "nettbe", current_net_tbe)
		SetPodSummary(summary_data, summary_key, "previous", "netrbt", previous_net_rbt)
		SetPodSummary(summary_data, summary_key, "previous", "netrbe", previous_net_rbe)
		SetPodSummary(summary_data, summary_key, "previous", "nettbt", previous_net_tbt)
		SetPodSummary(summary_data, summary_key, "previous", "nettbe", previous_net_tbe)

		insert_stat_data.ArrNetworkreceiverate = append(insert_stat_data.ArrNetworkreceiverate, rcv_rate)
		insert_stat_data.ArrNetworkreceiveerrors = append(insert_stat_data.ArrNetworkreceiveerrors, rcv_errors)
		insert_stat_data.ArrNetworktransmitrate = append(insert_stat_data.ArrNetworktransmitrate, trans_rate)
		insert_stat_data.ArrNetworktransmiterrors = append(insert_stat_data.ArrNetworktransmiterrors, trans_errors)
		insert_stat_data.ArrNetworkiorate = append(insert_stat_data.ArrNetworkiorate, rcv_rate+trans_rate)
		insert_stat_data.ArrNetworkioerrors = append(insert_stat_data.ArrNetworkioerrors, rcv_errors+trans_errors)

		current_fs_rbt := GetSummaryData(fs_map, summary_key, "current", "rbt")
		current_fs_wbt := GetSummaryData(fs_map, summary_key, "current", "wbt")
		previous_fs_rbt := GetSummaryData(fs_map, summary_key, "previous", "rbt")
		previous_fs_wbt := GetSummaryData(fs_map, summary_key, "previous", "wbt")
		fs_read_rate := GetStatData("rate", current_fs_rbt, previous_fs_rbt, 0)
		fs_write_rate := GetStatData("rate", current_fs_wbt, previous_fs_wbt, 0)
		SetPodSummary(summary_data, summary_key, "current", "fsrbt", current_fs_rbt)
		SetPodSummary(summary_data, summary_key, "current", "fswbt", current_fs_wbt)
		SetPodSummary(summary_data, summary_key, "previous", "fsrbt", previous_fs_rbt)
		SetPodSummary(summary_data, summary_key, "previous", "fswbt", previous_fs_wbt)

		insert_stat_data.ArrFsreadrate = append(insert_stat_data.ArrFsreadrate, fs_read_rate)
		insert_stat_data.ArrFswriterate = append(insert_stat_data.ArrFswriterate, fs_write_rate)
		insert_stat_data.ArrFsiorate = append(insert_stat_data.ArrFsiorate, fs_read_rate+fs_write_rate)
	}

	lsdm, _ := last_stat_data_map.LoadOrStore("pod", &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store("pod", lsdm)
	last_summary_data_map.Store(nodeuid, summary_data)

	return insert_stat_data, summary_data
}

func insertPodStatTable(nodeuid string, insert_stat_data PodPerfStat, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	tablename := getTableName(TB_KUBE_POD_PERF)
	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_POD_STAT, tablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	lasttablename := TB_KUBE_LAST_POD_PERF
	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA, lasttablename, nodeuid, ontunetime))
	errorCheck(del_err)

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_POD_STAT, lasttablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	err = tx.Commit(context.Background())
	errorCheck(err)

	conn.Release()
	update_tableinfo(tablename, ontunetime)
	update_tableinfo(lasttablename, ontunetime)
}

func setContainerMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var arrInsert ContainerPerf
	var flag bool = false

	for _, container_data := range map_data.([]interface{}) {
		container_data_map := container_data.(map[string]interface{})

		nsuid := getUID("namespace", (*map_metric)["host"], container_data_map["namespace"])
		poduid := getUID("pod", container_data_map["namespace"].(string), container_data_map["pod"].(string))
		nodeuid := getUID("node", (*map_metric)["host"], container_data_map["node"])
		image := container_data_map["image"]

		if nodeuid == "" {
			flag = true
			continue
		}

		if nsuid == "" {
			flag = true
			continue
		}

		if poduid == "" {
			flag = true
			continue
		}

		for item, data := range container_data_map {
			switch item {
			case METRIC_VAR_CONTAINER:
				arrInsert.ArrContainername = append(arrInsert.ArrContainername, data.(string))
			case METRIC_VAR_NAMESPACE:
				arrInsert.ArrNsuid = append(arrInsert.ArrNsuid, nsuid)
			case METRIC_VAR_POD:
				arrInsert.ArrPoduid = append(arrInsert.ArrPoduid, poduid)
			case METRIC_VAR_NODE:
				arrInsert.ArrNodeuid = append(arrInsert.ArrNodeuid, nodeuid)
			case METRIC_VAR_METRICID:
				arrInsert.ArrMetricid = append(arrInsert.ArrMetricid, getMetricID(data, image))
			case METRIC_VAR_TIMESTAMPMS:
				f_data := data.(float64)
				arrInsert.ArrTimestampMs = append(arrInsert.ArrTimestampMs, int64(f_data))
			case METRIC_VAR_CPUUSAGESEONDSTOTAL:
				arrInsert.ArrCpuusagesecondstotal = append(arrInsert.ArrCpuusagesecondstotal, data.(float64))
			case METRIC_VAR_CPUSYSTEMSECONDSTOTAL:
				arrInsert.ArrCpusystemsecondstotal = append(arrInsert.ArrCpusystemsecondstotal, data.(float64))
			case METRIC_VAR_CPUUSERSECONDSTOTAL:
				arrInsert.ArrCpuusersecondstotal = append(arrInsert.ArrCpuusersecondstotal, data.(float64))
			case METRIC_VAR_MEMORYUSAGEBYTES:
				arrInsert.ArrMemoryusagebytes = append(arrInsert.ArrMemoryusagebytes, data.(float64))
			case METRIC_VAR_MEMORYWORKINGSETBYTES:
				arrInsert.ArrMemoryworkingsetbytes = append(arrInsert.ArrMemoryworkingsetbytes, data.(float64))
			case METRIC_VAR_MEMORYCACHE:
				arrInsert.ArrMemorycache = append(arrInsert.ArrMemorycache, data.(float64))
			case METRIC_VAR_MEMORYSWAP:
				arrInsert.ArrMemoryswap = append(arrInsert.ArrMemoryswap, data.(float64))
			case METRIC_VAR_MEMORYRSS:
				arrInsert.ArrMemoryrss = append(arrInsert.ArrMemoryrss, data.(float64))
			case METRIC_VAR_PROCESSES:
				arrInsert.ArrProcesses = append(arrInsert.ArrProcesses, data.(float64))
			}
		}
		arrInsert.ArrOntunetime = append(arrInsert.ArrOntunetime, 0)
	}
	var m_data interface{} = arrInsert
	(*map_metric)["container_insert"] = m_data

	return flag
}

func insertContainerPerf(insert_data ContainerPerf, ontunetime int64) {
	perftablename := getTableName(TB_KUBE_CONTAINER_PERF_RAW)

	for i := 0; i < len(insert_data.ArrOntunetime); i++ {
		insert_data.ArrOntunetime[i] = ontunetime
	}

	if len(insert_data.ArrCpuusagesecondstotal) > 0 {
		conn, err := common.DBConnectionPool.Acquire(context.Background())
		if err != nil {
			errorCheck(err)
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CONTAINER_PERF, perftablename), insert_data.GetArgs()...)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)

		conn.Release()

		update_tableinfo(perftablename, ontunetime)
	}
}

func insertLastContainerPerf(insert_data ContainerPerf, ontunetime int64, nodeuid string) {
	for i := 0; i < len(insert_data.ArrOntunetime); i++ {
		insert_data.ArrOntunetime[i] = ontunetime
	}

	if len(insert_data.ArrCpuusagesecondstotal) > 0 {
		conn, err := common.DBConnectionPool.Acquire(context.Background())
		if err != nil {
			errorCheck(err)
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, del_err := tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA, TB_KUBE_LAST_CONTAINER_PERF_RAW, nodeuid, ontunetime))
		errorCheck(del_err)

		_, err = tx.Exec(context.Background(), INSERT_UNNEST_LAST_CONTAINER_PERF_RAW, insert_data.GetArgs()...)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)

		conn.Release()

		update_tableinfo(TB_KUBE_LAST_CONTAINER_PERF_RAW, ontunetime)
	}
}

func insertContainerStat(current_data *sync.Map, prev_data *sync.Map, fs_map *sync.Map, nodeuid string) *ContainerPerfStat {
	_, result := getStatTimes(current_data)
	if !result {
		return &ContainerPerfStat{}
	}

	var current_container_data ContainerPerf
	var prev_container_data ContainerPerf

	if c, ok := current_data.Load("container_insert"); ok {
		current_container_data = c.(ContainerPerf)
	} else {
		return &ContainerPerfStat{}
	}

	if p, ok := prev_data.Load("container_insert"); ok {
		prev_container_data = p.(ContainerPerf)
	} else {
		return &ContainerPerfStat{}
	}

	insert_stat_data := new(ContainerPerfStat)
	insert_stat_data.Init()

	for i := 0; i < len(current_container_data.ArrOntunetime); i++ {
		var pidx int = -1
		for j := 0; j < len(prev_container_data.ArrOntunetime); j++ {
			if prev_container_data.ArrMetricid[j] == current_container_data.ArrMetricid[i] {
				pidx = j
				break
			}
		}
		if pidx == -1 {
			break
		}
		key := current_container_data.ArrPoduid[i] + ":" + current_container_data.ArrContainername[i]
		container_status_map, _ := common.ResourceMap.Load("container_status")
		if cval, ok := container_status_map.(*sync.Map).Load(key); !ok || cval.(string) != CONTAINER_STATUS_RUNNING {
			continue
		}

		container_cpu_request_map, _ := common.ResourceMap.Load("container_cpu_request")
		container_cpu_limit_map, _ := common.ResourceMap.Load("container_cpu_limit")
		container_memory_request_map, _ := common.ResourceMap.Load("container_memory_request")
		container_memory_limit_map, _ := common.ResourceMap.Load("container_memory_limit")
		cpurequest, _ := container_cpu_request_map.(*sync.Map).LoadOrStore(key, int64(0))
		memoryrequest, _ := container_memory_request_map.(*sync.Map).LoadOrStore(key, int64(0))
		cpulimit, _ := container_cpu_limit_map.(*sync.Map).LoadOrStore(key, int64(0))
		memorylimit, _ := container_memory_limit_map.(*sync.Map).LoadOrStore(key, int64(0))

		var cpucount float64
		if cpulimit.(int64) > 0 {
			cpucount = float64(cpulimit.(int64))
		} else if cpurequest.(int64) > 0 {
			cpucount = float64(cpurequest.(int64))
		}

		var memorysize float64
		if memorylimit.(int64) > 0 {
			memorysize = float64(memorylimit.(int64))
		} else if memoryrequest.(int64) > 0 {
			memorysize = float64(memoryrequest.(int64))
		}

		insert_stat_data.ArrContainername = append(insert_stat_data.ArrContainername, current_container_data.ArrContainername[i])
		insert_stat_data.ArrPoduid = append(insert_stat_data.ArrPoduid, current_container_data.ArrPoduid[i])
		insert_stat_data.ArrNodeuid = append(insert_stat_data.ArrNodeuid, current_container_data.ArrNodeuid[i])
		insert_stat_data.ArrNsuid = append(insert_stat_data.ArrNsuid, current_container_data.ArrNsuid[i])
		insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_container_data.ArrOntunetime[i])
		insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, current_container_data.ArrMetricid[i])
		insert_stat_data.ArrProcesses = append(insert_stat_data.ArrProcesses, current_container_data.ArrProcesses[i])

		insert_stat_data.ArrCpuusage = append(insert_stat_data.ArrCpuusage, GetStatData("usage", current_container_data.ArrCpuusagesecondstotal[i], prev_container_data.ArrCpuusagesecondstotal[pidx], cpucount))
		insert_stat_data.ArrCpusystem = append(insert_stat_data.ArrCpusystem, GetStatData("usage", current_container_data.ArrCpusystemsecondstotal[i], prev_container_data.ArrCpusystemsecondstotal[pidx], cpucount))
		insert_stat_data.ArrCpuuser = append(insert_stat_data.ArrCpuuser, GetStatData("usage", current_container_data.ArrCpuusersecondstotal[i], prev_container_data.ArrCpuusersecondstotal[pidx], cpucount))
		insert_stat_data.ArrCpuusagecores = append(insert_stat_data.ArrCpuusagecores, GetStatData("core", current_container_data.ArrCpuusagesecondstotal[i], prev_container_data.ArrCpuusagesecondstotal[pidx], 0))
		insert_stat_data.ArrCpurequestcores = append(insert_stat_data.ArrCpurequestcores, float64(cpurequest.(int64)))
		insert_stat_data.ArrCpulimitcores = append(insert_stat_data.ArrCpulimitcores, float64(cpulimit.(int64)))
		insert_stat_data.ArrMemoryusage = append(insert_stat_data.ArrMemoryusage, GetStatData("current_usage", current_container_data.ArrMemoryworkingsetbytes[i], 0, float64(memorysize)))
		insert_stat_data.ArrMemoryusagebytes = append(insert_stat_data.ArrMemoryusagebytes, math.Round(current_container_data.ArrMemoryworkingsetbytes[i]))
		insert_stat_data.ArrMemoryrequestbytes = append(insert_stat_data.ArrMemoryrequestbytes, math.Round(float64(memoryrequest.(int64))))
		insert_stat_data.ArrMemorylimitbytes = append(insert_stat_data.ArrMemorylimitbytes, math.Round(float64(memorylimit.(int64))))
		insert_stat_data.ArrMemoryswap = append(insert_stat_data.ArrMemoryswap, GetStatData("current_usage", current_container_data.ArrMemoryswap[i], 0, memorysize))

		summary_key := SummaryIndex{
			ContainerName: current_container_data.ArrContainername[i],
			PodUID:        current_container_data.ArrPoduid[i],
		}

		current_fs_rbt := GetSummaryData(fs_map, summary_key, "current", "rbt")
		current_fs_wbt := GetSummaryData(fs_map, summary_key, "current", "wbt")
		previous_fs_rbt := GetSummaryData(fs_map, summary_key, "previous", "rbt")
		previous_fs_wbt := GetSummaryData(fs_map, summary_key, "previous", "wbt")
		fs_read_rate := GetStatData("rate", current_fs_rbt, previous_fs_rbt, 0)
		fs_write_rate := GetStatData("rate", current_fs_wbt, previous_fs_wbt, 0)

		insert_stat_data.ArrFsreadrate = append(insert_stat_data.ArrFsreadrate, fs_read_rate)
		insert_stat_data.ArrFswriterate = append(insert_stat_data.ArrFswriterate, fs_write_rate)
		insert_stat_data.ArrFsiorate = append(insert_stat_data.ArrFsiorate, fs_read_rate+fs_write_rate)
	}

	lsdm, _ := last_stat_data_map.LoadOrStore("container", &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store("container", lsdm)

	return insert_stat_data
}

func insertContainerStatTable(nodeuid string, insert_stat_data ContainerPerfStat, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	tablename := getTableName(TB_KUBE_CONTAINER_PERF)
	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CONTAINER_STAT, tablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	lasttablename := TB_KUBE_LAST_CONTAINER_PERF
	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA, lasttablename, nodeuid, ontunetime))
	errorCheck(del_err)

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CONTAINER_STAT, lasttablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	err = tx.Commit(context.Background())
	errorCheck(err)

	conn.Release()

	update_tableinfo(tablename, ontunetime)
	update_tableinfo(lasttablename, ontunetime)
}

func setNodeNetMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var arrInsert NodeNetPerf
	var flag bool = false

	for _, nodenet_data := range map_data.([]interface{}) {
		nodenet_data_map := nodenet_data.(map[string]interface{})
		nodeuid := getUID("node", (*map_metric)["host"], nodenet_data_map["node"])
		image := nodenet_data_map["image"]

		if nodeuid == "" {
			flag = true
			continue
		}

		for item, data := range nodenet_data_map {
			switch item {
			case METRIC_VAR_INTERFACE:
				arrInsert.ArrInterfaceid = append(arrInsert.ArrInterfaceid, getID("netinterface", data))
			case METRIC_VAR_NODE:
				arrInsert.ArrNodeuid = append(arrInsert.ArrNodeuid, nodeuid)
			case METRIC_VAR_METRICID:
				arrInsert.ArrMetricid = append(arrInsert.ArrMetricid, getMetricID(data, image))
			case METRIC_VAR_TIMESTAMPMS:
				f_data := data.(float64)
				arrInsert.ArrTimestampMs = append(arrInsert.ArrTimestampMs, int64(f_data))
			case METRIC_VAR_NETWORKRECEIVEBYTESTOTAL:
				arrInsert.ArrNetworkreceivebytestotal = append(arrInsert.ArrNetworkreceivebytestotal, data.(float64))
			case METRIC_VAR_NETWORKRECEIVEERRORSTOTAL:
				arrInsert.ArrNetworkreceiveerrorstotal = append(arrInsert.ArrNetworkreceiveerrorstotal, data.(float64))
			case METRIC_VAR_NETWORKTRANSMITBYTESTOTAL:
				arrInsert.ArrNetworktransmitbytestotal = append(arrInsert.ArrNetworktransmitbytestotal, data.(float64))
			case METRIC_VAR_NETWORKTRANSMITERRORSTOTAL:
				arrInsert.ArrNetworktransmiterrorstotal = append(arrInsert.ArrNetworktransmiterrorstotal, data.(float64))
			}
		}
		arrInsert.ArrOntunetime = append(arrInsert.ArrOntunetime, 0)
	}
	var m_data interface{} = arrInsert
	(*map_metric)["nodnet_insert"] = m_data

	return flag
}

func insertNodeNetPerf(insert_data NodeNetPerf, ontunetime int64) {
	perftablename := getTableName(TB_KUBE_NODE_NET_PERF_RAW)

	for i := 0; i < len(insert_data.ArrOntunetime); i++ {
		insert_data.ArrOntunetime[i] = ontunetime
	}

	if len(insert_data.ArrNetworkreceivebytestotal) > 0 {
		conn, err := common.DBConnectionPool.Acquire(context.Background())
		if err != nil {
			errorCheck(err)
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODE_NET_PERF, perftablename), insert_data.GetArgs()...)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)

		conn.Release()

		update_tableinfo(perftablename, ontunetime)
	}
}

func insertNodeNetStat(current_data *sync.Map, prev_data *sync.Map, nodeuid string) (*NodeNetPerfStat, *sync.Map) {
	_, result := getStatTimes(current_data)
	if !result {
		return &NodeNetPerfStat{}, &sync.Map{}
	}

	node_stat_metric_id := getMetricID(NODE_ROOT_METRIC_ID, "")

	var current_nodenet_data NodeNetPerf
	var prev_nodenet_data NodeNetPerf
	var summary_data *sync.Map = &sync.Map{}

	if c, ok := current_data.Load("nodnet_insert"); ok {
		current_nodenet_data = c.(NodeNetPerf)
	} else {
		return &NodeNetPerfStat{}, &sync.Map{}
	}

	if p, ok := prev_data.Load("nodnet_insert"); ok {
		prev_nodenet_data = p.(NodeNetPerf)
	} else {
		return &NodeNetPerfStat{}, &sync.Map{}
	}

	insert_stat_data := new(NodeNetPerfStat)
	insert_stat_data.Init()

	for i := 0; i < len(current_nodenet_data.ArrOntunetime); i++ {
		if current_nodenet_data.ArrMetricid[i] == node_stat_metric_id {
			var pidx int = -1
			for j := 0; j < len(prev_nodenet_data.ArrOntunetime); j++ {
				if prev_nodenet_data.ArrMetricid[j] == node_stat_metric_id && current_nodenet_data.ArrInterfaceid[i] == prev_nodenet_data.ArrInterfaceid[j] {
					pidx = j
					break
				}
			}
			if pidx == -1 {
				continue
			}

			SetSummaryData(summary_data, "net", &current_nodenet_data, &prev_nodenet_data, i, pidx)

			insert_stat_data.ArrNodeuid = append(insert_stat_data.ArrNodeuid, current_nodenet_data.ArrNodeuid[i])
			insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_nodenet_data.ArrOntunetime[i])
			insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, current_nodenet_data.ArrMetricid[i])
			insert_stat_data.ArrInterfaceid = append(insert_stat_data.ArrInterfaceid, current_nodenet_data.ArrInterfaceid[i])

			rcv_rate := GetStatData("rate", current_nodenet_data.ArrNetworkreceivebytestotal[i], prev_nodenet_data.ArrNetworkreceivebytestotal[pidx], 0)
			rcv_errors := GetStatData("subtract", current_nodenet_data.ArrNetworkreceiveerrorstotal[i], prev_nodenet_data.ArrNetworkreceiveerrorstotal[pidx], 0)
			trans_rate := GetStatData("rate", current_nodenet_data.ArrNetworktransmitbytestotal[i], prev_nodenet_data.ArrNetworktransmitbytestotal[pidx], 0)
			trans_errors := GetStatData("subtract", current_nodenet_data.ArrNetworktransmiterrorstotal[i], prev_nodenet_data.ArrNetworktransmiterrorstotal[pidx], 0)

			insert_stat_data.ArrNetworkreceiverate = append(insert_stat_data.ArrNetworkreceiverate, rcv_rate)
			insert_stat_data.ArrNetworkreceiveerrors = append(insert_stat_data.ArrNetworkreceiveerrors, rcv_errors)
			insert_stat_data.ArrNetworktransmitrate = append(insert_stat_data.ArrNetworktransmitrate, trans_rate)
			insert_stat_data.ArrNetworktransmiterrors = append(insert_stat_data.ArrNetworktransmiterrors, trans_errors)

			insert_stat_data.ArrNetworkiorate = append(insert_stat_data.ArrNetworkiorate, rcv_rate+trans_rate)
			insert_stat_data.ArrNetworkioerrors = append(insert_stat_data.ArrNetworkioerrors, rcv_errors+trans_errors)
		}
	}

	lsdm, _ := last_stat_data_map.LoadOrStore("nodenet", &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store("nodenet", lsdm)

	return insert_stat_data, summary_data
}

func insertNodeNetStatTable(nodeuid string, insert_stat_data NodeNetPerfStat, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	tablename := getTableName(TB_KUBE_NODE_NET_PERF)
	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODENET_STAT, tablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	lasttablename := TB_KUBE_LAST_NODE_NET_PERF
	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA, lasttablename, nodeuid, ontunetime))
	errorCheck(del_err)

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODENET_STAT, lasttablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	err = tx.Commit(context.Background())
	errorCheck(err)

	conn.Release()

	update_tableinfo(tablename, ontunetime)
	update_tableinfo(lasttablename, ontunetime)
}

func setPodNetMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var arrInsert PodNetPerf
	var flag bool = false

	for _, podnet_data := range map_data.([]interface{}) {
		podnet_data_map := podnet_data.(map[string]interface{})

		nsuid := getUID("namespace", (*map_metric)["host"], podnet_data_map["namespace"])
		poduid := getUID("pod", podnet_data_map["namespace"].(string), podnet_data_map["pod"].(string))
		nodeuid := getUID("node", (*map_metric)["host"], podnet_data_map["node"])
		image := podnet_data_map["image"]

		if nodeuid == "" {
			flag = true
			continue
		}

		if nsuid == "" {
			flag = true
			continue
		}

		if poduid == "" {
			flag = true
			continue
		}

		for item, data := range podnet_data_map {
			switch item {
			case METRIC_VAR_NAMESPACE:
				arrInsert.ArrNsuid = append(arrInsert.ArrNsuid, nsuid)
			case METRIC_VAR_POD:
				arrInsert.ArrPoduid = append(arrInsert.ArrPoduid, poduid)
			case METRIC_VAR_NODE:
				arrInsert.ArrNodeuid = append(arrInsert.ArrNodeuid, nodeuid)
			case METRIC_VAR_METRICID:
				arrInsert.ArrMetricid = append(arrInsert.ArrMetricid, getMetricID(data, image))
			case METRIC_VAR_INTERFACE:
				arrInsert.ArrInterfaceid = append(arrInsert.ArrInterfaceid, getID("netinterface", data))
			case METRIC_VAR_TIMESTAMPMS:
				f_data := data.(float64)
				arrInsert.ArrTimestampMs = append(arrInsert.ArrTimestampMs, int64(f_data))
			case METRIC_VAR_NETWORKRECEIVEBYTESTOTAL:
				arrInsert.ArrNetworkreceivebytestotal = append(arrInsert.ArrNetworkreceivebytestotal, data.(float64))
			case METRIC_VAR_NETWORKRECEIVEERRORSTOTAL:
				arrInsert.ArrNetworkreceiveerrorstotal = append(arrInsert.ArrNetworkreceiveerrorstotal, data.(float64))
			case METRIC_VAR_NETWORKTRANSMITBYTESTOTAL:
				arrInsert.ArrNetworktransmitbytestotal = append(arrInsert.ArrNetworktransmitbytestotal, data.(float64))
			case METRIC_VAR_NETWORKTRANSMITERRORSTOTAL:
				arrInsert.ArrNetworktransmiterrorstotal = append(arrInsert.ArrNetworktransmiterrorstotal, data.(float64))
			}
		}
		arrInsert.ArrOntunetime = append(arrInsert.ArrOntunetime, 0)
	}
	var m_data interface{} = arrInsert
	(*map_metric)["podnet_insert"] = m_data

	return flag
}

func insertPodNetPerf(insert_data PodNetPerf, ontunetime int64) {
	perftablename := getTableName(TB_KUBE_POD_NET_PERF_RAW)

	for i := 0; i < len(insert_data.ArrOntunetime); i++ {
		insert_data.ArrOntunetime[i] = ontunetime
	}

	if len(insert_data.ArrNetworkreceivebytestotal) > 0 {
		conn, err := common.DBConnectionPool.Acquire(context.Background())
		if err != nil {
			errorCheck(err)
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_POD_NET_PERF, perftablename), insert_data.GetArgs()...)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)

		conn.Release()

		update_tableinfo(perftablename, ontunetime)
	}
}

func insertPodNetStat(current_data *sync.Map, prev_data *sync.Map, nodeuid string) (*PodNetPerfStat, *sync.Map) {
	_, result := getStatTimes(current_data)
	if !result {
		return &PodNetPerfStat{}, &sync.Map{}
	}

	var current_podnet_summary_data PodNetPerf
	var prev_podnet_summary_data PodNetPerf
	var summary_data *sync.Map = &sync.Map{}

	if c, ok := current_data.Load("podnet_insert"); ok {
		current_podnet_data := c.(PodNetPerf)
		current_podnet_summary_data = PodNetSummaryPreprocessing(&current_podnet_data)
	} else {
		return &PodNetPerfStat{}, &sync.Map{}
	}

	if p, ok := prev_data.Load("podnet_insert"); ok {
		prev_podnet_data := p.(PodNetPerf)
		prev_podnet_summary_data = PodNetSummaryPreprocessing(&prev_podnet_data)
	} else {
		return &PodNetPerfStat{}, &sync.Map{}
	}

	insert_stat_data := new(PodNetPerfStat)
	insert_stat_data.Init()

	for i := 0; i < len(current_podnet_summary_data.ArrOntunetime); i++ {
		var pidx int = -1
		for j := 0; j < len(prev_podnet_summary_data.ArrOntunetime); j++ {
			if prev_podnet_summary_data.ArrPoduid[j] == current_podnet_summary_data.ArrPoduid[i] && prev_podnet_summary_data.ArrInterfaceid[j] == current_podnet_summary_data.ArrInterfaceid[i] {
				pidx = j
				break
			}
		}
		if pidx == -1 {
			continue
		}
		pod_status_map, _ := common.ResourceMap.Load("pod_status")
		if pval, ok := pod_status_map.(*sync.Map).Load(current_podnet_summary_data.ArrPoduid[i]); !ok || pval.(string) != POD_STATUS_RUNNING {
			continue
		}

		SetSummaryData(summary_data, "net", &current_podnet_summary_data, &prev_podnet_summary_data, i, pidx)

		insert_stat_data.ArrPoduid = append(insert_stat_data.ArrPoduid, current_podnet_summary_data.ArrPoduid[i])
		insert_stat_data.ArrNodeuid = append(insert_stat_data.ArrNodeuid, current_podnet_summary_data.ArrNodeuid[i])
		insert_stat_data.ArrNsuid = append(insert_stat_data.ArrNsuid, current_podnet_summary_data.ArrNsuid[i])
		insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_podnet_summary_data.ArrOntunetime[i])
		insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, 0)
		insert_stat_data.ArrInterfaceid = append(insert_stat_data.ArrInterfaceid, current_podnet_summary_data.ArrInterfaceid[i])

		rcv_rate := GetStatData("rate", current_podnet_summary_data.ArrNetworkreceivebytestotal[i], prev_podnet_summary_data.ArrNetworkreceivebytestotal[pidx], 0)
		rcv_errors := GetStatData("subtract", current_podnet_summary_data.ArrNetworkreceiveerrorstotal[i], prev_podnet_summary_data.ArrNetworkreceiveerrorstotal[pidx], 0)
		trans_rate := GetStatData("rate", current_podnet_summary_data.ArrNetworktransmitbytestotal[i], prev_podnet_summary_data.ArrNetworktransmitbytestotal[pidx], 0)
		trans_errors := GetStatData("subtract", current_podnet_summary_data.ArrNetworktransmiterrorstotal[i], prev_podnet_summary_data.ArrNetworktransmiterrorstotal[pidx], 0)

		insert_stat_data.ArrNetworkreceiverate = append(insert_stat_data.ArrNetworkreceiverate, rcv_rate)
		insert_stat_data.ArrNetworkreceiveerrors = append(insert_stat_data.ArrNetworkreceiveerrors, rcv_errors)
		insert_stat_data.ArrNetworktransmitrate = append(insert_stat_data.ArrNetworktransmitrate, trans_rate)
		insert_stat_data.ArrNetworktransmiterrors = append(insert_stat_data.ArrNetworktransmiterrors, trans_errors)

		insert_stat_data.ArrNetworkiorate = append(insert_stat_data.ArrNetworkiorate, rcv_rate+trans_rate)
		insert_stat_data.ArrNetworkioerrors = append(insert_stat_data.ArrNetworkioerrors, rcv_errors+trans_errors)
	}

	lsdm, _ := last_stat_data_map.LoadOrStore("podnet", &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store("podnet", lsdm)

	return insert_stat_data, summary_data
}

func insertPodNetStatTable(nodeuid string, insert_stat_data PodNetPerfStat, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	tablename := getTableName(TB_KUBE_POD_NET_PERF)
	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_PODNET_STAT, tablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	lasttablename := TB_KUBE_LAST_POD_NET_PERF
	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA, lasttablename, nodeuid, ontunetime))
	errorCheck(del_err)

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_PODNET_STAT, lasttablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	err = tx.Commit(context.Background())
	errorCheck(err)

	conn.Release()

	update_tableinfo(tablename, ontunetime)
	update_tableinfo(lasttablename, ontunetime)
}

func setNodeFsMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var arrInsert NodeFsPerf
	var flag bool = false

	for _, nodefs_data := range map_data.([]interface{}) {
		nodefs_data_map := nodefs_data.(map[string]interface{})
		nodeuid := getUID("node", (*map_metric)["host"], nodefs_data_map["node"])
		image := nodefs_data_map["image"]
		if nodeuid == "" {
			flag = true
			continue
		}

		for item, data := range nodefs_data_map {
			switch item {
			case METRIC_VAR_NODE:
				arrInsert.ArrNodeuid = append(arrInsert.ArrNodeuid, nodeuid)
			case METRIC_VAR_METRICID:
				arrInsert.ArrMetricid = append(arrInsert.ArrMetricid, getMetricID(data, image))
			case METRIC_VAR_DEVICE:
				arrInsert.ArrDeviceid = append(arrInsert.ArrDeviceid, getID("fsdevice", data))
			case METRIC_VAR_TIMESTAMPMS:
				f_data := data.(float64)
				arrInsert.ArrTimestampMs = append(arrInsert.ArrTimestampMs, int64(f_data))
			case METRIC_VAR_FSINODESFREE:
				arrInsert.ArrFsinodesfree = append(arrInsert.ArrFsinodesfree, data.(float64))
			case METRIC_VAR_FSINODESTOTAL:
				arrInsert.ArrFsinodestotal = append(arrInsert.ArrFsinodestotal, data.(float64))
			case METRIC_VAR_FSLIMITBYTES:
				arrInsert.ArrFslimitbytes = append(arrInsert.ArrFslimitbytes, data.(float64))
			case METRIC_VAR_FSREADSBYTESTOTAL:
				arrInsert.ArrFsreadsbytestotal = append(arrInsert.ArrFsreadsbytestotal, data.(float64))
			case METRIC_VAR_FSWRITESBYTESTOTAL:
				arrInsert.ArrFswritesbytestotal = append(arrInsert.ArrFswritesbytestotal, data.(float64))
			case METRIC_VAR_FSUSAGEBYTES:
				arrInsert.ArrFsusagebytes = append(arrInsert.ArrFsusagebytes, data.(float64))
			}
		}
		arrInsert.ArrOntunetime = append(arrInsert.ArrOntunetime, 0)
	}
	var m_data interface{} = arrInsert
	(*map_metric)["nodefs_insert"] = m_data

	return flag
}

func insertNodeFsPerf(insert_data NodeFsPerf, ontunetime int64) {
	perftablename := getTableName(TB_KUBE_NODE_FS_PERF_RAW)

	for i := 0; i < len(insert_data.ArrOntunetime); i++ {
		insert_data.ArrOntunetime[i] = ontunetime
	}

	if len(insert_data.ArrFsreadsbytestotal) > 0 {
		conn, err := common.DBConnectionPool.Acquire(context.Background())
		if err != nil {
			errorCheck(err)
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODE_FS_PERF, perftablename), insert_data.GetArgs()...)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)

		conn.Release()

		update_tableinfo(perftablename, ontunetime)
	}
}

func insertNodeFsStat(current_data *sync.Map, prev_data *sync.Map, nodeuid string) (*NodeFsPerfStat, *sync.Map) {
	_, result := getStatTimes(current_data)
	if !result {
		return &NodeFsPerfStat{}, &sync.Map{}
	}

	node_stat_metric_id := getMetricID(NODE_ROOT_METRIC_ID, "")

	var current_nodefs_data NodeFsPerf
	var prev_nodefs_data NodeFsPerf
	var summary_data *sync.Map = &sync.Map{}

	if c, ok := current_data.Load("nodefs_insert"); ok {
		current_nodefs_data = c.(NodeFsPerf)
	} else {
		return &NodeFsPerfStat{}, &sync.Map{}
	}

	if p, ok := prev_data.Load("nodefs_insert"); ok {
		prev_nodefs_data = p.(NodeFsPerf)
	} else {
		return &NodeFsPerfStat{}, &sync.Map{}
	}

	insert_stat_data := new(NodeFsPerfStat)
	insert_stat_data.Init()

	for i := 0; i < len(current_nodefs_data.ArrOntunetime); i++ {
		if current_nodefs_data.ArrMetricid[i] == node_stat_metric_id {
			var pidx int = -1
			for j := 0; j < len(prev_nodefs_data.ArrOntunetime); j++ {
				if prev_nodefs_data.ArrMetricid[j] == node_stat_metric_id && prev_nodefs_data.ArrDeviceid[j] == current_nodefs_data.ArrDeviceid[i] {
					pidx = j
					break
				}
			}
			if pidx == -1 {
				continue
			}

			SetSummaryData(summary_data, "fs", &current_nodefs_data, &prev_nodefs_data, i, pidx)

			insert_stat_data.ArrNodeuid = append(insert_stat_data.ArrNodeuid, current_nodefs_data.ArrNodeuid[i])
			insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_nodefs_data.ArrOntunetime[i])
			insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, current_nodefs_data.ArrMetricid[i])
			insert_stat_data.ArrDeviceid = append(insert_stat_data.ArrDeviceid, current_nodefs_data.ArrDeviceid[i])

			read_rate := GetStatData("rate", current_nodefs_data.ArrFsreadsbytestotal[i], prev_nodefs_data.ArrFsreadsbytestotal[pidx], 0)
			write_rate := GetStatData("rate", current_nodefs_data.ArrFswritesbytestotal[i], prev_nodefs_data.ArrFswritesbytestotal[pidx], 0)

			insert_stat_data.ArrFsreadrate = append(insert_stat_data.ArrFsreadrate, read_rate)
			insert_stat_data.ArrFswriterate = append(insert_stat_data.ArrFswriterate, write_rate)
			insert_stat_data.ArrFsiorate = append(insert_stat_data.ArrFsiorate, read_rate+write_rate)
		}
	}

	lsdm, _ := last_stat_data_map.LoadOrStore("nodefs", &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store("nodefs", lsdm)

	return insert_stat_data, summary_data
}

func insertNodeFsStatTable(nodeuid string, insert_stat_data NodeFsPerfStat, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	tablename := getTableName(TB_KUBE_NODE_FS_PERF)
	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODEFS_STAT, tablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	lasttablename := TB_KUBE_LAST_NODE_FS_PERF
	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA, lasttablename, nodeuid, ontunetime))
	errorCheck(del_err)

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODEFS_STAT, lasttablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	err = tx.Commit(context.Background())
	errorCheck(err)

	conn.Release()

	update_tableinfo(tablename, ontunetime)
	update_tableinfo(lasttablename, ontunetime)
}

func setPodFsMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var arrInsert PodFsPerf
	var flag bool = false

	for _, podfs_data := range map_data.([]interface{}) {
		podfs_data_map := podfs_data.(map[string]interface{})

		nsuid := getUID("namespace", (*map_metric)["host"], podfs_data_map["namespace"])
		poduid := getUID("pod", podfs_data_map["namespace"].(string), podfs_data_map["pod"].(string))
		nodeuid := getUID("node", (*map_metric)["host"], podfs_data_map["node"])
		image := podfs_data.(map[string]interface{})["image"]

		if nodeuid == "" {
			flag = true
			continue
		}

		if nsuid == "" {
			flag = true
			continue
		}

		if poduid == "" {
			flag = true
			continue
		}

		for item, data := range podfs_data.(map[string]interface{}) {
			switch item {
			case METRIC_VAR_POD:
				arrInsert.ArrPoduid = append(arrInsert.ArrPoduid, poduid)
			case METRIC_VAR_NAMESPACE:
				arrInsert.ArrNsuid = append(arrInsert.ArrNsuid, nsuid)
			case METRIC_VAR_NODE:
				arrInsert.ArrNodeuid = append(arrInsert.ArrNodeuid, nodeuid)
			case METRIC_VAR_METRICID:
				arrInsert.ArrMetricid = append(arrInsert.ArrMetricid, getMetricID(data, image))
			case METRIC_VAR_DEVICE:
				arrInsert.ArrDeviceid = append(arrInsert.ArrDeviceid, getID("fsdevice", data))
			case METRIC_VAR_TIMESTAMPMS:
				f_data := data.(float64)
				arrInsert.ArrTimestampMs = append(arrInsert.ArrTimestampMs, int64(f_data))
			case METRIC_VAR_FSREADSBYTESTOTAL:
				arrInsert.ArrFsreadsbytestotal = append(arrInsert.ArrFsreadsbytestotal, data.(float64))
			case METRIC_VAR_FSWRITESBYTESTOTAL:
				arrInsert.ArrFswritesbytestotal = append(arrInsert.ArrFswritesbytestotal, data.(float64))
			}
		}
		arrInsert.ArrOntunetime = append(arrInsert.ArrOntunetime, 0)
	}
	var m_data interface{} = arrInsert
	(*map_metric)["podfs_insert"] = m_data

	return flag
}

func insertPodFsPerf(insert_data PodFsPerf, ontunetime int64) {
	perftablename := getTableName(TB_KUBE_POD_FS_PERF_RAW)
	for i := 0; i < len(insert_data.ArrOntunetime); i++ {
		insert_data.ArrOntunetime[i] = ontunetime
	}

	if len(insert_data.ArrFsreadsbytestotal) > 0 {
		conn, err := common.DBConnectionPool.Acquire(context.Background())
		if err != nil {
			errorCheck(err)
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_POD_FS_PERF, perftablename), insert_data.GetArgs()...)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)

		conn.Release()

		update_tableinfo(perftablename, ontunetime)
	}
}

func insertPodFsStat(current_data *sync.Map, prev_data *sync.Map, nodeuid string) (*PodFsPerfStat, *sync.Map) {
	_, result := getStatTimes(current_data)
	if !result {
		return &PodFsPerfStat{}, &sync.Map{}
	}

	var current_podfs_summary_data PodFsPerf
	var prev_podfs_summary_data PodFsPerf
	var summary_data *sync.Map = &sync.Map{}

	if c, ok := current_data.Load("podfs_insert"); ok {
		current_podfs_data := c.(PodFsPerf)

		var current_containerfs_data ContainerFsPerf
		if cc, ok2 := current_data.Load("containerfs_insert"); ok2 {
			current_containerfs_data = cc.(ContainerFsPerf)
		}
		current_podfs_summary_data = PodFsSummaryPreprocessing(&current_podfs_data, &current_containerfs_data)
	} else {
		return &PodFsPerfStat{}, &sync.Map{}
	}

	if p, ok := prev_data.Load("podfs_insert"); ok {
		prev_podfs_data := p.(PodFsPerf)

		var prev_containerfs_data ContainerFsPerf
		if pc, ok2 := prev_data.Load("containerfs_insert"); ok2 {
			prev_containerfs_data = pc.(ContainerFsPerf)
		}
		prev_podfs_summary_data = PodFsSummaryPreprocessing(&prev_podfs_data, &prev_containerfs_data)
	} else {
		return &PodFsPerfStat{}, &sync.Map{}
	}

	insert_stat_data := new(PodFsPerfStat)
	insert_stat_data.Init()

	for i := 0; i < len(current_podfs_summary_data.ArrOntunetime); i++ {
		var pidx int = -1
		for j := 0; j < len(prev_podfs_summary_data.ArrOntunetime); j++ {
			if prev_podfs_summary_data.ArrPoduid[j] == current_podfs_summary_data.ArrPoduid[i] && prev_podfs_summary_data.ArrDeviceid[j] == current_podfs_summary_data.ArrDeviceid[i] {
				pidx = j
				break
			}
		}
		if pidx == -1 {
			continue
		}
		pod_status_map, _ := common.ResourceMap.Load("pod_status")
		if pval, ok := pod_status_map.(*sync.Map).Load(current_podfs_summary_data.ArrPoduid[i]); !ok || pval.(string) != POD_STATUS_RUNNING {
			continue
		}

		SetSummaryData(summary_data, "fs", &current_podfs_summary_data, &prev_podfs_summary_data, i, pidx)

		insert_stat_data.ArrPoduid = append(insert_stat_data.ArrPoduid, current_podfs_summary_data.ArrPoduid[i])
		insert_stat_data.ArrNodeuid = append(insert_stat_data.ArrNodeuid, current_podfs_summary_data.ArrNodeuid[i])
		insert_stat_data.ArrNsuid = append(insert_stat_data.ArrNsuid, current_podfs_summary_data.ArrNsuid[i])
		insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_podfs_summary_data.ArrOntunetime[i])
		insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, 0)
		insert_stat_data.ArrDeviceid = append(insert_stat_data.ArrDeviceid, current_podfs_summary_data.ArrDeviceid[i])

		read_rate := GetStatData("rate", current_podfs_summary_data.ArrFsreadsbytestotal[i], prev_podfs_summary_data.ArrFsreadsbytestotal[pidx], 0)
		write_rate := GetStatData("rate", current_podfs_summary_data.ArrFswritesbytestotal[i], prev_podfs_summary_data.ArrFswritesbytestotal[pidx], 0)

		insert_stat_data.ArrFsreadrate = append(insert_stat_data.ArrFsreadrate, read_rate)
		insert_stat_data.ArrFswriterate = append(insert_stat_data.ArrFswriterate, write_rate)
		insert_stat_data.ArrFsiorate = append(insert_stat_data.ArrFsiorate, read_rate+write_rate)
	}

	lsdm, _ := last_stat_data_map.LoadOrStore("podfs", &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store("podfs", lsdm)

	return insert_stat_data, summary_data
}

func insertPodFsStatTable(nodeuid string, insert_stat_data PodFsPerfStat, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	tablename := getTableName(TB_KUBE_POD_FS_PERF)
	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_PODFS_STAT, tablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	lasttablename := TB_KUBE_LAST_POD_FS_PERF
	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA, lasttablename, nodeuid, ontunetime))
	errorCheck(del_err)

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_PODFS_STAT, lasttablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	err = tx.Commit(context.Background())
	errorCheck(err)

	conn.Release()

	update_tableinfo(tablename, ontunetime)
	update_tableinfo(lasttablename, ontunetime)
}

func setContainerFsMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var arrInsert ContainerFsPerf
	var flag bool = false

	for _, containerfs_data := range map_data.([]interface{}) {
		containerfs_data_map := containerfs_data.(map[string]interface{})

		nsuid := getUID("namespace", (*map_metric)["host"], containerfs_data_map["namespace"])
		poduid := getUID("pod", containerfs_data_map["namespace"].(string), containerfs_data_map["pod"].(string))
		nodeuid := getUID("node", (*map_metric)["host"], containerfs_data_map["node"])
		image := containerfs_data_map["image"]

		if nodeuid == "" {
			flag = true
			continue
		}

		if nsuid == "" {
			flag = true
			continue
		}

		if poduid == "" {
			flag = true
			continue
		}

		for item, data := range containerfs_data_map {
			switch item {
			case METRIC_VAR_CONTAINER:
				arrInsert.ArrContainername = append(arrInsert.ArrContainername, data.(string))
			case METRIC_VAR_NAMESPACE:
				arrInsert.ArrNsuid = append(arrInsert.ArrNsuid, nsuid)
			case METRIC_VAR_POD:
				arrInsert.ArrPoduid = append(arrInsert.ArrPoduid, poduid)
			case METRIC_VAR_NODE:
				arrInsert.ArrNodeuid = append(arrInsert.ArrNodeuid, nodeuid)
			case METRIC_VAR_METRICID:
				arrInsert.ArrMetricid = append(arrInsert.ArrMetricid, getMetricID(data, image))
			case METRIC_VAR_DEVICE:
				arrInsert.ArrDeviceid = append(arrInsert.ArrDeviceid, getID("fsdevice", data))
			case METRIC_VAR_TIMESTAMPMS:
				f_data := data.(float64)
				arrInsert.ArrTimestampMs = append(arrInsert.ArrTimestampMs, int64(f_data))
			case METRIC_VAR_FSREADSBYTESTOTAL:
				arrInsert.ArrFsreadsbytestotal = append(arrInsert.ArrFsreadsbytestotal, data.(float64))
			case METRIC_VAR_FSWRITESBYTESTOTAL:
				arrInsert.ArrFswritesbytestotal = append(arrInsert.ArrFswritesbytestotal, data.(float64))
			}
		}
		arrInsert.ArrOntunetime = append(arrInsert.ArrOntunetime, 0)
	}
	var m_data interface{} = arrInsert
	(*map_metric)["containerfs_insert"] = m_data

	return flag
}

func insertContainerFsPerf(insert_data ContainerFsPerf, ontunetime int64) {
	perftablename := getTableName(TB_KUBE_CONTAINER_FS_PERF_RAW)
	for i := 0; i < len(insert_data.ArrOntunetime); i++ {
		insert_data.ArrOntunetime[i] = ontunetime
	}

	if len(insert_data.ArrFsreadsbytestotal) > 0 {
		conn, err := common.DBConnectionPool.Acquire(context.Background())
		if err != nil {
			errorCheck(err)
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CONTAINER_FS_PERF, perftablename), insert_data.GetArgs()...)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)

		conn.Release()

		update_tableinfo(perftablename, ontunetime)
	}
}

func insertContainerFsStat(current_data *sync.Map, prev_data *sync.Map, nodeuid string) (*ContainerFsPerfStat, *sync.Map) {
	_, result := getStatTimes(current_data)
	if !result {
		return &ContainerFsPerfStat{}, &sync.Map{}
	}

	var current_containerfs_data ContainerFsPerf
	var prev_containerfs_data ContainerFsPerf
	var summary_data *sync.Map = &sync.Map{}

	if c, ok := current_data.Load("containerfs_insert"); ok {
		current_containerfs_data = c.(ContainerFsPerf)
	} else {
		return &ContainerFsPerfStat{}, &sync.Map{}
	}

	if p, ok := prev_data.Load("containerfs_insert"); ok {
		prev_containerfs_data = p.(ContainerFsPerf)
	} else {
		return &ContainerFsPerfStat{}, &sync.Map{}
	}

	insert_stat_data := new(ContainerFsPerfStat)
	insert_stat_data.Init()

	for i := 0; i < len(current_containerfs_data.ArrOntunetime); i++ {
		var pidx int = -1
		for j := 0; j < len(prev_containerfs_data.ArrOntunetime); j++ {
			if prev_containerfs_data.ArrMetricid[j] == current_containerfs_data.ArrMetricid[i] && prev_containerfs_data.ArrDeviceid[j] == current_containerfs_data.ArrDeviceid[i] {
				pidx = j
				break
			}
		}
		if pidx == -1 {
			break
		}
		key := current_containerfs_data.ArrPoduid[i] + ":" + current_containerfs_data.ArrContainername[i]
		container_status_map, _ := common.ResourceMap.Load("container_status")
		if cval, ok := container_status_map.(*sync.Map).Load(key); !ok || cval.(string) != CONTAINER_STATUS_RUNNING {
			continue
		}

		SetSummaryData(summary_data, "fs", &current_containerfs_data, &prev_containerfs_data, i, pidx)

		insert_stat_data.ArrContainername = append(insert_stat_data.ArrContainername, current_containerfs_data.ArrContainername[i])
		insert_stat_data.ArrPoduid = append(insert_stat_data.ArrPoduid, current_containerfs_data.ArrPoduid[i])
		insert_stat_data.ArrNodeuid = append(insert_stat_data.ArrNodeuid, current_containerfs_data.ArrNodeuid[i])
		insert_stat_data.ArrNsuid = append(insert_stat_data.ArrNsuid, current_containerfs_data.ArrNsuid[i])
		insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_containerfs_data.ArrOntunetime[i])
		insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, current_containerfs_data.ArrMetricid[i])
		insert_stat_data.ArrDeviceid = append(insert_stat_data.ArrDeviceid, current_containerfs_data.ArrDeviceid[i])

		read_rate := GetStatData("rate", current_containerfs_data.ArrFsreadsbytestotal[i], prev_containerfs_data.ArrFsreadsbytestotal[pidx], 0)
		write_rate := GetStatData("rate", current_containerfs_data.ArrFswritesbytestotal[i], prev_containerfs_data.ArrFswritesbytestotal[pidx], 0)

		insert_stat_data.ArrFsreadrate = append(insert_stat_data.ArrFsreadrate, read_rate)
		insert_stat_data.ArrFswriterate = append(insert_stat_data.ArrFswriterate, write_rate)
		insert_stat_data.ArrFsiorate = append(insert_stat_data.ArrFsiorate, read_rate+write_rate)
	}

	lsdm, _ := last_stat_data_map.LoadOrStore("containerfs", &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store("containerfs", lsdm)

	return insert_stat_data, summary_data
}

func insertContainerFsStatTable(nodeuid string, insert_stat_data ContainerFsPerfStat, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	tablename := getTableName(TB_KUBE_CONTAINER_FS_PERF)
	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CONTAINERFS_STAT, tablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	lasttablename := TB_KUBE_LAST_CONTAINER_FS_PERF
	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA, lasttablename, nodeuid, ontunetime))
	errorCheck(del_err)

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CONTAINERFS_STAT, lasttablename), insert_stat_data.GetArgs()...)
	errorCheck(err)

	err = tx.Commit(context.Background())
	errorCheck(err)

	conn.Release()

	update_tableinfo(tablename, ontunetime)
	update_tableinfo(lasttablename, ontunetime)
}

func insertClusterStat(cluster_map *sync.Map, ontunetime int64) {
	insert_stat_data := new(ClusterPerfStat)
	insert_stat_data.Init()

	cluster_map.Range(func(key, nodevalue any) bool {
		clusterid := key.(int)
		insert_stat_data.ArrClusterid = append(insert_stat_data.ArrClusterid, clusterid)
		insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, ontunetime)

		var cal_map *sync.Map = &sync.Map{}
		cal_map.Store("calculation", &sync.Map{})
		cal_map.Store("count", 0)

		nodevalue.(*sync.Map).Range(func(nuid, pmap any) bool {
			nodeuid := nuid.(string)
			pod_map := pmap.(*sync.Map)

			cal_map.Store("current", &sync.Map{})
			cal_map.Store("previous", &sync.Map{})

			pod_map.Range(func(sumkey, pod_value any) bool {
				summary_index := sumkey.(SummaryIndex)

				if p, ok := cal_map.Load("count"); ok && summary_index.DataType == "current" {
					cal_map.Store("count", p.(int)+1)
				}

				if cm, ok := cal_map.Load(summary_index.DataType); ok {
					stat_map := cm.(*sync.Map)
					value_map := pod_value.(*sync.Map)

					value_map.Range(func(mname, mvalue any) bool {
						metricname := mname.(string)
						sm, _ := stat_map.LoadOrStore(metricname, float64(0))
						stat_value := sm.(float64) + mvalue.(float64)
						stat_map.Store(metricname, stat_value)

						return true
					})
				}

				return true
			})

			CalculateNodeSummary(cal_map, nodeuid)

			return true
		})

		cnt, _ := cal_map.Load("count")
		insert_stat_data.ArrPodcount = append(insert_stat_data.ArrPodcount, cnt.(int))

		calc, _ := cal_map.Load("calculation")
		calculation_map := calc.(*sync.Map)

		cpuusage, _ := calculation_map.Load("cpuusage")
		insert_stat_data.ArrCpuusage = append(insert_stat_data.ArrCpuusage, cpuusage.(float64))

		cpusystem, _ := calculation_map.Load("cpusystem")
		insert_stat_data.ArrCpusystem = append(insert_stat_data.ArrCpusystem, cpusystem.(float64))

		cpuuser, _ := calculation_map.Load("cpuuser")
		insert_stat_data.ArrCpuuser = append(insert_stat_data.ArrCpuuser, cpuuser.(float64))

		cpuusagecores, _ := calculation_map.Load("cpuusagecores")
		insert_stat_data.ArrCpuusagecores = append(insert_stat_data.ArrCpuusagecores, cpuusagecores.(float64))

		cputotalcores, _ := calculation_map.Load("cputotalcores")
		insert_stat_data.ArrCputotalcores = append(insert_stat_data.ArrCputotalcores, cputotalcores.(float64))

		memoryusage, _ := calculation_map.Load("memoryusage")
		insert_stat_data.ArrMemoryusage = append(insert_stat_data.ArrMemoryusage, memoryusage.(float64))

		memoryusagebytes, _ := calculation_map.Load("memoryusagebytes")
		insert_stat_data.ArrMemoryusagebytes = append(insert_stat_data.ArrMemoryusagebytes, memoryusagebytes.(float64))

		memorysizebytes, _ := calculation_map.Load("memorysizebytes")
		insert_stat_data.ArrMemorysizebytes = append(insert_stat_data.ArrMemorysizebytes, memorysizebytes.(float64))

		memoryswap, _ := calculation_map.Load("memoryswap")
		insert_stat_data.ArrMemoryswap = append(insert_stat_data.ArrMemoryswap, memoryswap.(float64))

		netrcvrate, _ := calculation_map.Load("netrcvrate")
		insert_stat_data.ArrNetworkreceiverate = append(insert_stat_data.ArrNetworkreceiverate, netrcvrate.(float64))

		netrcverrors, _ := calculation_map.Load("netrcverrors")
		insert_stat_data.ArrNetworkreceiveerrors = append(insert_stat_data.ArrNetworkreceiveerrors, netrcverrors.(float64))

		nettransrate, _ := calculation_map.Load("nettransrate")
		insert_stat_data.ArrNetworktransmitrate = append(insert_stat_data.ArrNetworktransmitrate, nettransrate.(float64))

		nettranserrors, _ := calculation_map.Load("nettranserrors")
		insert_stat_data.ArrNetworktransmiterrors = append(insert_stat_data.ArrNetworktransmiterrors, nettranserrors.(float64))

		netiorate, _ := calculation_map.Load("netiorate")
		insert_stat_data.ArrNetworkiorate = append(insert_stat_data.ArrNetworkiorate, netiorate.(float64))

		netioerrors, _ := calculation_map.Load("netioerrors")
		insert_stat_data.ArrNetworkioerrors = append(insert_stat_data.ArrNetworkioerrors, netioerrors.(float64))

		fsreadrate, _ := calculation_map.Load("fsreadrate")
		insert_stat_data.ArrFsreadrate = append(insert_stat_data.ArrFsreadrate, fsreadrate.(float64))

		fswriterate, _ := calculation_map.Load("fswriterate")
		insert_stat_data.ArrFswriterate = append(insert_stat_data.ArrFswriterate, fswriterate.(float64))

		fsiorate, _ := calculation_map.Load("fsiorate")
		insert_stat_data.ArrFsiorate = append(insert_stat_data.ArrFsiorate, fsiorate.(float64))

		processcount, _ := calculation_map.Load("processcount")
		insert_stat_data.ArrProcesses = append(insert_stat_data.ArrProcesses, processcount.(float64))

		return true
	})

	if len(insert_stat_data.ArrClusterid) > 0 {
		conn, err := common.DBConnectionPool.Acquire(context.Background())
		if err != nil {
			errorCheck(err)
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		tablename := getTableName(TB_KUBE_CLUSTER_PERF)
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CLUSTER_STAT, tablename), insert_stat_data.GetArgs()...)
		if ok := errorCheck(err); !ok {
			return
		}

		lasttablename := TB_KUBE_LAST_CLUSTER_PERF

		cluster_ids_map := make(map[string]struct{})
		cluster_ids := make([]string, 0)
		for _, clusterid := range insert_stat_data.ArrClusterid {
			cluster_id_str := strconv.Itoa(clusterid)
			if _, ok := cluster_ids_map[cluster_id_str]; !ok {
				cluster_ids_map[cluster_id_str] = struct{}{}
				cluster_ids = append(cluster_ids, cluster_id_str)
			}
		}
		join_clusterids := strings.Join(cluster_ids, ",")

		_, del_err := tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA_CONDITION, lasttablename, "clusterid", join_clusterids, ontunetime))
		errorCheck(del_err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CLUSTER_STAT, lasttablename), insert_stat_data.GetArgs()...)
		if ok := errorCheck(err); !ok {
			return
		}

		err = tx.Commit(context.Background())
		errorCheck(err)

		conn.Release()

		update_tableinfo(tablename, ontunetime)
		update_tableinfo(lasttablename, ontunetime)
	}
}

func insertNamespaceStat(cluster_map *sync.Map, ontunetime int64) {
	insert_stat_data := new(NamespacePerfStat)
	insert_stat_data.Init()

	cluster_map.Range(func(key, nodevalue any) bool {
		var cal_ns_map *sync.Map = &sync.Map{}

		nodevalue.(*sync.Map).Range(func(nuid, pmap any) bool {
			nodeuid := nuid.(string)
			pod_map := pmap.(*sync.Map)

			pod_map.Range(func(sumkey, pod_value any) bool {
				summary_index := sumkey.(SummaryIndex)

				if summary_index.NsUID == "" {
					return true
				}

				var ns_map *sync.Map
				if ns, ok := cal_ns_map.LoadOrStore(summary_index.NsUID, &sync.Map{}); !ok {
					ns_map = ns.(*sync.Map)
					ns_map.Store("calculation", &sync.Map{})
					ns_map.Store("current", &sync.Map{})
					ns_map.Store("previous", &sync.Map{})
					ns_map.Store("count", 0)
				} else {
					ns_map = ns.(*sync.Map)
					ns_map.LoadOrStore("current", &sync.Map{})
					ns_map.LoadOrStore("previous", &sync.Map{})
				}

				if p, ok := ns_map.Load("count"); ok && summary_index.DataType == "current" {
					ns_map.Store("count", p.(int)+1)
				}

				if cm, ok := ns_map.Load(summary_index.DataType); ok {
					stat_map := cm.(*sync.Map)
					value_map := pod_value.(*sync.Map)

					value_map.Range(func(mname, mvalue any) bool {
						metricname := mname.(string)
						sm, _ := stat_map.LoadOrStore(metricname, float64(0))
						stat_value := sm.(float64) + mvalue.(float64)
						stat_map.Store(metricname, stat_value)

						return true
					})
				}

				return true
			})

			cal_ns_map.Range(func(nsuid, nsvalue any) bool {
				ns_map := nsvalue.(*sync.Map)

				if _, ok := ns_map.Load("current"); ok {
					CalculateNodeSummary(ns_map, nodeuid)
				}

				return true
			})

			return true
		})

		cal_ns_map.Range(func(nsuid, nsvalue any) bool {
			ns_map := nsvalue.(*sync.Map)

			insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, ontunetime)
			insert_stat_data.ArrNsuid = append(insert_stat_data.ArrNsuid, nsuid.(string))

			cnt, _ := ns_map.Load("count")
			insert_stat_data.ArrPodcount = append(insert_stat_data.ArrPodcount, cnt.(int))

			calc, _ := ns_map.Load("calculation")
			calculation_map := calc.(*sync.Map)

			cpuusage, _ := calculation_map.Load("cpuusage")
			insert_stat_data.ArrCpuusage = append(insert_stat_data.ArrCpuusage, cpuusage.(float64))

			cpusystem, _ := calculation_map.Load("cpusystem")
			insert_stat_data.ArrCpusystem = append(insert_stat_data.ArrCpusystem, cpusystem.(float64))

			cpuuser, _ := calculation_map.Load("cpuuser")
			insert_stat_data.ArrCpuuser = append(insert_stat_data.ArrCpuuser, cpuuser.(float64))

			cpuusagecores, _ := calculation_map.Load("cpuusagecores")
			insert_stat_data.ArrCpuusagecores = append(insert_stat_data.ArrCpuusagecores, cpuusagecores.(float64))

			cputotalcores, _ := calculation_map.Load("cputotalcores")
			insert_stat_data.ArrCputotalcores = append(insert_stat_data.ArrCputotalcores, cputotalcores.(float64))

			memoryusage, _ := calculation_map.Load("memoryusage")
			insert_stat_data.ArrMemoryusage = append(insert_stat_data.ArrMemoryusage, memoryusage.(float64))

			memoryusagebytes, _ := calculation_map.Load("memoryusagebytes")
			insert_stat_data.ArrMemoryusagebytes = append(insert_stat_data.ArrMemoryusagebytes, memoryusagebytes.(float64))

			memorysizebytes, _ := calculation_map.Load("memorysizebytes")
			insert_stat_data.ArrMemorysizebytes = append(insert_stat_data.ArrMemorysizebytes, memorysizebytes.(float64))

			memoryswap, _ := calculation_map.Load("memoryswap")
			insert_stat_data.ArrMemoryswap = append(insert_stat_data.ArrMemoryswap, memoryswap.(float64))

			netrcvrate, _ := calculation_map.Load("netrcvrate")
			insert_stat_data.ArrNetworkreceiverate = append(insert_stat_data.ArrNetworkreceiverate, netrcvrate.(float64))

			netrcverrors, _ := calculation_map.Load("netrcverrors")
			insert_stat_data.ArrNetworkreceiveerrors = append(insert_stat_data.ArrNetworkreceiveerrors, netrcverrors.(float64))

			nettransrate, _ := calculation_map.Load("nettransrate")
			insert_stat_data.ArrNetworktransmitrate = append(insert_stat_data.ArrNetworktransmitrate, nettransrate.(float64))

			nettranserrors, _ := calculation_map.Load("nettranserrors")
			insert_stat_data.ArrNetworktransmiterrors = append(insert_stat_data.ArrNetworktransmiterrors, nettranserrors.(float64))

			netiorate, _ := calculation_map.Load("netiorate")
			insert_stat_data.ArrNetworkiorate = append(insert_stat_data.ArrNetworkiorate, netiorate.(float64))

			netioerrors, _ := calculation_map.Load("netioerrors")
			insert_stat_data.ArrNetworkioerrors = append(insert_stat_data.ArrNetworkioerrors, netioerrors.(float64))

			fsreadrate, _ := calculation_map.Load("fsreadrate")
			insert_stat_data.ArrFsreadrate = append(insert_stat_data.ArrFsreadrate, fsreadrate.(float64))

			fswriterate, _ := calculation_map.Load("fswriterate")
			insert_stat_data.ArrFswriterate = append(insert_stat_data.ArrFswriterate, fswriterate.(float64))

			fsiorate, _ := calculation_map.Load("fsiorate")
			insert_stat_data.ArrFsiorate = append(insert_stat_data.ArrFsiorate, fsiorate.(float64))

			processcount, _ := calculation_map.Load("processcount")
			insert_stat_data.ArrProcesses = append(insert_stat_data.ArrProcesses, processcount.(float64))

			cal_ns_map.Delete(nsuid)

			return true
		})

		return true
	})

	if len(insert_stat_data.ArrNsuid) > 0 {
		conn, err := common.DBConnectionPool.Acquire(context.Background())
		if err != nil {
			errorCheck(err)
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		tablename := getTableName(TB_KUBE_NAMESPACE_PERF)
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NS_STAT, tablename), insert_stat_data.GetArgs()...)
		if ok := errorCheck(err); !ok {
			return
		}

		lasttablename := TB_KUBE_LAST_NAMESPACE_PERF

		ns_uids_map := make(map[string]struct{})
		ns_uids := make([]string, 0)
		for _, nsuid := range insert_stat_data.ArrNsuid {
			if _, ok := ns_uids_map[nsuid]; !ok {
				ns_uids_map[nsuid] = struct{}{}
				uid_mark := fmt.Sprintf("'%s'", nsuid)
				ns_uids = append(ns_uids, uid_mark)
			}
		}
		join_nsuids := strings.Join(ns_uids, ",")

		_, del_err := tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA_CONDITION, lasttablename, "nsuid", join_nsuids, ontunetime))
		errorCheck(del_err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NS_STAT, lasttablename), insert_stat_data.GetArgs()...)
		if ok := errorCheck(err); !ok {
			return
		}

		err = tx.Commit(context.Background())
		errorCheck(err)

		conn.Release()

		update_tableinfo(tablename, ontunetime)
		update_tableinfo(lasttablename, ontunetime)
	}
}

func insertWorkloadStat(cluster_map *sync.Map, kind string, ontunetime int64) {
	var insert_stat_data WorkloadPerfStat
	insert_stat_data.Init()

	var tablename string
	var lasttablename string
	var uidname string

	switch kind {
	case "ReplicaSet":
		tablename = getTableName(TB_KUBE_REPLICASET_PERF)
		lasttablename = TB_KUBE_LAST_REPLICASET_PERF
		uidname = "rsuid"
	case "Deployment":
		tablename = getTableName(TB_KUBE_DEPLOYMENT_PERF)
		lasttablename = TB_KUBE_LAST_DEPLOYMENT_PERF
		uidname = "deployuid"

		if pod_reference_map, ok := common.ResourceMap.Load("pod_reference"); ok {
			pod_reference_map.(*sync.Map).Range(func(key, value any) bool {
				poduid := key.(string)
				refmap := value.(map[string]string)

				if rsuid, ok2 := refmap["ReplicaSet"]; ok2 {
					if rs_reference_map, ok3 := common.ResourceMap.Load("replicaset_reference"); ok3 {
						if deploymap, ok4 := rs_reference_map.(*sync.Map).Load(rsuid); ok4 {
							if deployuid, ok5 := deploymap.(map[string]string)["Deployment"]; ok5 {
								refmap["Deployment"] = deployuid
								pod_reference_map.(*sync.Map).Store(poduid, refmap)
							}
						}
					}
				}

				return true
			})
		}
	case "StatefulSet":
		tablename = getTableName(TB_KUBE_STATEFULSET_PERF)
		lasttablename = TB_KUBE_LAST_STATEFULSET_PERF
		uidname = "stsuid"
	case "DaemonSet":
		tablename = getTableName(TB_KUBE_DAEMONSET_PERF)
		lasttablename = TB_KUBE_LAST_DAEMONSET_PERF
		uidname = "dsuid"
	}

	cluster_map.Range(func(key, nodevalue any) bool {
		cal_wl_map := ProcessWorkloadSummaryByNode(nodevalue.(*sync.Map), kind)

		cal_wl_map.Range(func(wluid, wlvalue any) bool {
			wl_map := wlvalue.(*sync.Map)
			insert_stat_data.AppendValue(wluid.(string), ontunetime, wl_map)
			cal_wl_map.Delete(wluid)

			return true
		})

		return true
	})

	if len(insert_stat_data.ArrWorkloaduid) > 0 {
		conn, err := common.DBConnectionPool.Acquire(context.Background())
		if err != nil {
			errorCheck(err)
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_WORKLOAD_STAT, tablename, uidname), insert_stat_data.GetArgs()...)
		if ok := errorCheck(err); !ok {
			return
		}

		workload_uids_map := make(map[string]struct{})
		workload_uids := make([]string, 0)
		for _, wluid := range insert_stat_data.ArrWorkloaduid {
			if _, ok := workload_uids_map[wluid]; !ok {
				workload_uids_map[wluid] = struct{}{}
				uid_mark := fmt.Sprintf("'%s'", wluid)
				workload_uids = append(workload_uids, uid_mark)
			}
		}
		join_workload_uids := strings.Join(workload_uids, ",")

		_, del_err := tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA_CONDITION, lasttablename, uidname, join_workload_uids, ontunetime))
		errorCheck(del_err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_WORKLOAD_STAT, lasttablename, uidname), insert_stat_data.GetArgs()...)
		if ok := errorCheck(err); !ok {
			return
		}

		err = tx.Commit(context.Background())
		errorCheck(err)

		conn.Release()

		update_tableinfo(tablename, ontunetime)
		update_tableinfo(lasttablename, ontunetime)
	}
}

func insertServiceStat(cluster_map *sync.Map, ontunetime int64) {
}

func insertIngressStat(cluster_map *sync.Map, ontunetime int64) {
}

// Metric Insert
func MetricSender() {
	data_queue_map := &sync.Map{}

	go MetricInsertMap(data_queue_map)

	for {
		metric_data := <-Channelmetric_insert

		var nodeuid string
		var nodename string = metric_data["node"].(string)
		var hostname string = metric_data["host"].(string)

		common.MutexNode.Lock()
		if ar, ok := mapApiResource.Load(hostname); ok {
			apiresource := ar.(*ApiResource)
			mapNodeInfo := apiresource.node

			for k, v := range mapNodeInfo {
				if v.Name == nodename {
					nodeuid = k
					break
				}
			}
		}
		common.MutexNode.Unlock()

		m, _ := data_queue_map.LoadOrStore(nodeuid, make([]*sync.Map, 0))
		metric_data_map := &sync.Map{}
		ontunetime, _ := GetOntuneTime()
		metric_data_map.Store("hostname", hostname)
		metric_data_map.Store("ontunetime", ontunetime)
		metric_data_map.Store("node_insert", metric_data["node_insert"])
		metric_data_map.Store("pod_insert", metric_data["pod_insert"])
		metric_data_map.Store("container_insert", metric_data["container_insert"])
		metric_data_map.Store("nodnet_insert", metric_data["nodnet_insert"])
		metric_data_map.Store("podnet_insert", metric_data["podnet_insert"])
		metric_data_map.Store("nodefs_insert", metric_data["nodefs_insert"])
		metric_data_map.Store("podfs_insert", metric_data["podfs_insert"])
		metric_data_map.Store("containerfs_insert", metric_data["containerfs_insert"])

		metric_arr := make([]*sync.Map, 0)
		metric_arr = append(metric_arr, m.([]*sync.Map)...)
		metric_arr = append(metric_arr, metric_data_map)
		data_queue_map.Store(nodeuid, metric_arr)

		time.Sleep(time.Millisecond * 10)
	}
}

func MetricInsertMap(data_queue_map *sync.Map) {
	var last_updatetime int64
	_, bias := GetOntuneTime()

	for {
		ontunetime := time.Now().Unix() - int64(bias)
		if last_updatetime+int64(common.RealtimeInterval) <= ontunetime {
			// Select Query는 최초 1회, 그리고 Interval 내 조건을 충족할 때에만 수행하는 것으로 함
			_, bias = GetOntuneTime()
			last_updatetime = ontunetime

			// 초기화
			initMapUidInfo(data_queue_map)

			var insert_node_stat_data map[string]NodePerfStat = make(map[string]NodePerfStat)
			var insert_pod_stat_data map[string]PodPerfStat = make(map[string]PodPerfStat)
			var insert_container_stat_data map[string]ContainerPerfStat = make(map[string]ContainerPerfStat)
			var insert_nodnet_stat_data map[string]NodeNetPerfStat = make(map[string]NodeNetPerfStat)
			var insert_podnet_stat_data map[string]PodNetPerfStat = make(map[string]PodNetPerfStat)
			var insert_nodefs_stat_data map[string]NodeFsPerfStat = make(map[string]NodeFsPerfStat)
			var insert_podfs_stat_data map[string]PodFsPerfStat = make(map[string]PodFsPerfStat)
			var insert_containerfs_stat_data map[string]ContainerFsPerfStat = make(map[string]ContainerFsPerfStat)

			data_queue_map.Range(func(key, value any) bool {
				nodeuid := key.(string)
				metric_arr := value.([]*sync.Map)

				if len(metric_arr) == 0 {
					return true
				}

				// Get Recent Queue Value
				metric_data := metric_arr[len(metric_arr)-1]

				var hostname string
				if hn, ok := metric_data.Load("hostname"); ok {
					hostname = hn.(string)
				}

				if node_data, ok := metric_data.Load("node_insert"); ok {
					node_insert_data := node_data.(NodePerf)
					insertNodePerf(node_insert_data, ontunetime, nodeuid)
					insertLastNodePerf(node_insert_data, ontunetime, nodeuid)

					for _, v := range node_insert_data.ArrNodeuid {
						mapUidInfo[hostname].currentNodeUid[v] = struct{}{}
					}
				}
				if pod_data, ok := metric_data.Load("pod_insert"); ok {
					pod_insert_data := pod_data.(PodPerf)
					insertPodPerf(pod_insert_data, ontunetime, nodeuid)
					insertLastPodPerf(pod_insert_data, ontunetime, nodeuid)

					for _, v := range pod_insert_data.ArrPoduid {
						mapUidInfo[hostname].currentPodUid[v] = struct{}{}
					}

					for _, v := range pod_insert_data.ArrNsuid {
						mapUidInfo[hostname].currentNamespaceUid[v] = struct{}{}
					}
				}
				if container_data, ok := metric_data.Load("container_insert"); ok {
					insertContainerPerf(container_data.(ContainerPerf), ontunetime)
					insertLastContainerPerf(container_data.(ContainerPerf), ontunetime, nodeuid)
				}
				if nodenet_data, ok := metric_data.Load("nodnet_insert"); ok {
					insertNodeNetPerf(nodenet_data.(NodeNetPerf), ontunetime)
				}
				if podnet_data, ok := metric_data.Load("podnet_insert"); ok {
					insertPodNetPerf(podnet_data.(PodNetPerf), ontunetime)
				}
				if nodefs_data, ok := metric_data.Load("nodefs_insert"); ok {
					insertNodeFsPerf(nodefs_data.(NodeFsPerf), ontunetime)
				}
				if podfs_data, ok := metric_data.Load("podfs_insert"); ok {
					insertPodFsPerf(podfs_data.(PodFsPerf), ontunetime)
				}
				if containerfs_data, ok := metric_data.Load("containerfs_insert"); ok {
					insertContainerFsPerf(containerfs_data.(ContainerFsPerf), ontunetime)
				}

				var prev_data *sync.Map = &sync.Map{}
				var prev_ontunetime int64

				for i := len(metric_arr) - 1; i >= 0; i-- {
					if pot, ok := metric_arr[i].Load("ontunetime"); ok {
						prev_ontunetime = pot.(int64)
					}
					if prev_ontunetime <= ontunetime-int64(common.RateInterval) {
						// Dequeue
						prev_data = metric_arr[i]
						data_queue_map.Store(key, metric_arr[i+1:])

						break
					}
				}

				if MapNilCheck(prev_data) {
					if _, ok := prev_data.Load("node_insert"); ok {
						nndata, net_map := insertNodeNetStat(metric_data, prev_data, nodeuid)
						nfdata, fs_map := insertNodeFsStat(metric_data, prev_data, nodeuid)
						ndata := insertNodeStat(metric_data, prev_data, net_map, fs_map, nodeuid)

						insert_node_stat_data[nodeuid] = *ndata
						insert_nodnet_stat_data[nodeuid] = *nndata
						insert_nodefs_stat_data[nodeuid] = *nfdata
					}
					if _, ok := prev_data.Load("pod_insert"); ok {
						pndata, net_map := insertPodNetStat(metric_data, prev_data, nodeuid)
						pfdata, fs_map := insertPodFsStat(metric_data, prev_data, nodeuid)
						pdata, pod_map := insertPodStat(metric_data, prev_data, net_map, fs_map, nodeuid)

						insert_pod_stat_data[nodeuid] = *pdata
						insert_podnet_stat_data[nodeuid] = *pndata
						insert_podfs_stat_data[nodeuid] = *pfdata

						StoreClusterMap(cluster_map, pod_map, nodeuid)
					}
					if _, ok := prev_data.Load("container_insert"); ok {
						cfdata, fs_map := insertContainerFsStat(metric_data, prev_data, nodeuid)
						cdata := insertContainerStat(metric_data, prev_data, fs_map, nodeuid)

						insert_container_stat_data[nodeuid] = *cdata
						insert_containerfs_stat_data[nodeuid] = *cfdata
					}
				} else {
					if lsdm, ok := last_stat_data_map.Load("node"); ok {
						if lsdmn, ok2 := lsdm.(*sync.Map).Load(nodeuid); ok2 {
							lsdm_node := lsdmn.(NodePerfStat)
							refreshOntunetime(&lsdm_node.ArrOntunetime, ontunetime)

							insert_node_stat_data[nodeuid] = lsdm_node
						}
					}
					if lsdm, ok := last_stat_data_map.Load("nodenet"); ok {
						if lsdmn, ok2 := lsdm.(*sync.Map).Load(nodeuid); ok2 {
							lsdm_node := lsdmn.(NodeNetPerfStat)
							refreshOntunetime(&lsdm_node.ArrOntunetime, ontunetime)

							insert_nodnet_stat_data[nodeuid] = lsdm_node
						}
					}
					if lsdm, ok := last_stat_data_map.Load("nodefs"); ok {
						if lsdmn, ok2 := lsdm.(*sync.Map).Load(nodeuid); ok2 {
							lsdm_node := lsdmn.(NodeFsPerfStat)
							refreshOntunetime(&lsdm_node.ArrOntunetime, ontunetime)

							insert_nodefs_stat_data[nodeuid] = lsdm_node
						}
					}
					if lsdm, ok := last_stat_data_map.Load("pod"); ok {
						if lsdmn, ok2 := lsdm.(*sync.Map).Load(nodeuid); ok2 {
							lsdm_node := lsdmn.(PodPerfStat)
							refreshOntunetime(&lsdm_node.ArrOntunetime, ontunetime)

							pm, _ := last_summary_data_map.Load(nodeuid)
							pod_map := pm.(*sync.Map)

							insert_pod_stat_data[nodeuid] = lsdm_node

							StoreClusterMap(cluster_map, pod_map, nodeuid)
						}
					}
					if lsdm, ok := last_stat_data_map.Load("podnet"); ok {
						if lsdmn, ok2 := lsdm.(*sync.Map).Load(nodeuid); ok2 {
							lsdm_node := lsdmn.(PodNetPerfStat)
							refreshOntunetime(&lsdm_node.ArrOntunetime, ontunetime)

							insert_podnet_stat_data[nodeuid] = lsdm_node
						}
					}
					if lsdm, ok := last_stat_data_map.Load("podfs"); ok {
						if lsdmn, ok2 := lsdm.(*sync.Map).Load(nodeuid); ok2 {
							lsdm_node := lsdmn.(PodFsPerfStat)
							refreshOntunetime(&lsdm_node.ArrOntunetime, ontunetime)

							insert_podfs_stat_data[nodeuid] = lsdm_node
						}
					}
					if lsdm, ok := last_stat_data_map.Load("container"); ok {
						if lsdmn, ok2 := lsdm.(*sync.Map).Load(nodeuid); ok2 {
							lsdm_node := lsdmn.(ContainerPerfStat)
							refreshOntunetime(&lsdm_node.ArrOntunetime, ontunetime)

							insert_container_stat_data[nodeuid] = lsdm_node
						}
					}
					if lsdm, ok := last_stat_data_map.Load("containerfs"); ok {
						if lsdmn, ok2 := lsdm.(*sync.Map).Load(nodeuid); ok2 {
							lsdm_node := lsdmn.(ContainerFsPerfStat)
							refreshOntunetime(&lsdm_node.ArrOntunetime, ontunetime)

							insert_containerfs_stat_data[nodeuid] = lsdm_node
						}
					}
				}

				return true
			})

			for k, v := range insert_node_stat_data {
				insertNodeStatTable(k, v, ontunetime)
			}

			for k, v := range insert_nodnet_stat_data {
				insertNodeNetStatTable(k, v, ontunetime)
			}

			for k, v := range insert_nodefs_stat_data {
				insertNodeFsStatTable(k, v, ontunetime)
			}

			for k, v := range insert_pod_stat_data {
				insertPodStatTable(k, v, ontunetime)
			}

			for k, v := range insert_podnet_stat_data {
				insertPodNetStatTable(k, v, ontunetime)
			}

			for k, v := range insert_podfs_stat_data {
				insertPodFsStatTable(k, v, ontunetime)
			}

			for k, v := range insert_container_stat_data {
				insertContainerStatTable(k, v, ontunetime)
			}

			for k, v := range insert_containerfs_stat_data {
				insertContainerFsStatTable(k, v, ontunetime)
			}

			// Summary Stat Insert
			if MapNilCheck(cluster_map) {
				insertClusterStat(cluster_map, ontunetime)
				insertNamespaceStat(cluster_map, ontunetime)
				insertWorkloadStat(cluster_map, "ReplicaSet", ontunetime)
				insertWorkloadStat(cluster_map, "DaemonSet", ontunetime)
				insertWorkloadStat(cluster_map, "StatefulSet", ontunetime)
				insertWorkloadStat(cluster_map, "Deployment", ontunetime)
				insertServiceStat(cluster_map, ontunetime)
				insertIngressStat(cluster_map, ontunetime)
			}

			checkMapUidInfo(data_queue_map)
		}
		time.Sleep(time.Millisecond * 10)
	}
}
