package database

import (
	"context"
	"errors"
	"fmt"
	"math"
	"onTuneKubeManager/common"
	"onTuneKubeManager/kubeapi"
	"strconv"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

func getUID(map_type string, parent_key interface{}, key interface{}) string {
	key_str := fmt.Sprintf("%s/%s", parent_key, key)
	map_key := strings.ToLower(map_type)

	wl_map, _ := common.ResourceMap.Load(map_key)
	if returnVal, ok := wl_map.(*sync.Map).Load(key_str); ok {
		return returnVal.(string)
	} else {
		return ""
	}
}

func getID(map_type string, data interface{}) int {
	data_str := data.(string)
	if data_str == "" {
		return 0
	}

	switch map_type {
	case "fsdevice":
		if _, ok := mapFsDeviceInfo[data_str]; !ok {
			updateFsDeviceId(data_str)
		}

		return mapFsDeviceInfo[data_str]
	case "netinterface":
		if _, ok := mapNetInterfaceInfo[data_str]; !ok {
			updateNetInterfaceId(data_str)
		}

		return mapNetInterfaceInfo[data_str]
	}

	return 0
}

func getMetricID(id interface{}, image interface{}) int {
	id_str := id.(string)
	image_str := image.(string)
	if id_str == "" {
		return 0
	}

	if _, ok := mapMetricIdInfo[id_str]; !ok {
		updateMetricId(id_str, image_str)
	}

	return mapMetricIdInfo[id_str]
}

func getStatTimes(pod_data *sync.Map) (int64, bool) {
	var ontunetime int64
	if o, ok := pod_data.Load("ontunetime"); ok {
		// prev_data의 ontunetime을 빼지 않고 RateInterval만큼의 수치를 대상으로 함
		ontunetime = o.(int64)
		return ontunetime, true
	} else {
		return 0, false
	}
}

func update_tableinfo(tablename string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	_, err = tx.Exec(context.Background(), UPDATE_TABLEINFO, ontunetime, tablename)
	errorCheck(err)

	err = tx.Commit(context.Background())
	errorCheck(err)
}

func select_row_enabled(colName string, tablename string, clusterid int) pgx.Rows {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	common.LogManager.Debug(fmt.Sprintf("select %s from %s where enabled=1 and clusterid=%d", colName, tablename, clusterid))
	rows, err := conn.Query(context.Background(), fmt.Sprintf("select %s from %s where enabled=1 and clusterid=%d", colName, tablename, clusterid))
	errorCheck(err)

	return rows
}

func select_row_count_enabled(tablename string, clusterid int) int {
	var cnt int
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), fmt.Sprintf("select count(*) as cnt from %s where enabled=1 and clusterid=%d", tablename, clusterid))
	errorCheck(err)
	for rows.Next() {
		err := rows.Scan(&cnt)
		errorCheck(err)
	}

	return cnt
}

func select_row_count(tablename string) int {
	var cnt int
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), "select count(*) as cnt from "+tablename)
	errorCheck(err)
	for rows.Next() {
		err := rows.Scan(&cnt)
		errorCheck(err)
	}

	return cnt
}

func select_one_value_condition(colName string, tablename string, condition string, conditionVal string) string {
	var colValue string
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), fmt.Sprintf("select %s from %s where %s='%s'", colName, tablename, condition, conditionVal))
	errorCheck(err)
	for rows.Next() {
		err := rows.Scan(&colValue)
		errorCheck(err)
	}

	return colValue
}

func select_condition(colName string, tablename string, condition string, operator string, conditionVal string) pgx.Rows {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	errorCheck(err)

	defer conn.Release()

	common.LogManager.Debug(fmt.Sprintf("select %s from %s where %s%s'%s'", colName, tablename, condition, operator, conditionVal))
	rows, err := conn.Query(context.Background(), fmt.Sprintf("select %s from %s where %s%s'%s'", colName, tablename, condition, operator, conditionVal))
	errorCheck(err)

	return rows
}

func errorCheck(err error) bool {
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			common.LogManager.WriteLog(fmt.Sprintf("Database Error - %v %v", pqerr.Code, err.Error()))
			return false
		} else if errors.Is(err, pgx.ErrNoRows) {
			// No Message
			return false
		} else {
			common.LogManager.WriteLog(fmt.Sprintf("Database Error - %v", err.Error()))
			return false
		}
	}

	return true
}

func errorCheckQueryRow(err error) {
	if err != nil {
		common.LogManager.WriteLog(fmt.Sprintf("Database Error - %v", err.Error()))
	}
}

func errorRecover() {
	if r := recover(); r != nil {
		err := fmt.Errorf("Database Error - %v", r)
		common.LogManager.WriteLog(err.Error())
	}
}

func RowCountCheck(err error) bool {
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return false
	} else {
		return true
	}
}

func writeLog(message string) {
	logmsg := "\t" + message
	common.LogManager.WriteLog(logmsg)
}

func PodSummaryPreprocessing(pod_data *PodPerf, container_data *ContainerPerf) PodSummaryPerf {
	var summary_data PodSummaryPerf = PodSummaryPerf{}
	var podindex map[string]int = make(map[string]int)
	var idx int

	container_cpu_request_map, _ := common.ResourceMap.Load("container_cpu_request")
	container_cpu_limit_map, _ := common.ResourceMap.Load("container_cpu_limit")
	container_memory_request_map, _ := common.ResourceMap.Load("container_memory_request")
	container_memory_limit_map, _ := common.ResourceMap.Load("container_memory_limit")

	for i := 0; i < len(container_data.ArrOntunetime); i++ {
		var container_key string = container_data.ArrPoduid[i] + ":" + container_data.ArrContainername[i]
		cpurequest, _ := container_cpu_request_map.(*sync.Map).LoadOrStore(container_key, int64(0))
		cpulimit, _ := container_cpu_limit_map.(*sync.Map).LoadOrStore(container_key, int64(0))
		memoryrequest, _ := container_memory_request_map.(*sync.Map).LoadOrStore(container_key, int64(0))
		memorylimit, _ := container_memory_limit_map.(*sync.Map).LoadOrStore(container_key, int64(0))

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

		var key string = container_data.ArrPoduid[i] + ":" + container_data.ArrNsuid[i]
		if val, ok := podindex[key]; !ok {
			podindex[key] = idx
			summary_data.ArrPoduid = append(summary_data.ArrPoduid, container_data.ArrPoduid[i])
			summary_data.ArrOntunetime = append(summary_data.ArrOntunetime, container_data.ArrOntunetime[i])
			summary_data.ArrTimestampMs = append(summary_data.ArrTimestampMs, container_data.ArrTimestampMs[i])
			summary_data.ArrNsuid = append(summary_data.ArrNsuid, container_data.ArrNsuid[i])
			summary_data.ArrNodeuid = append(summary_data.ArrNodeuid, container_data.ArrNodeuid[i])
			summary_data.ArrMetricid = append(summary_data.ArrMetricid, 0)
			summary_data.ArrCpuusagesecondstotal = append(summary_data.ArrCpuusagesecondstotal, container_data.ArrCpuusagesecondstotal[i])
			summary_data.ArrCpusystemsecondstotal = append(summary_data.ArrCpusystemsecondstotal, container_data.ArrCpusystemsecondstotal[i])
			summary_data.ArrCpuusersecondstotal = append(summary_data.ArrCpuusersecondstotal, container_data.ArrCpuusersecondstotal[i])
			summary_data.ArrMemoryusagebytes = append(summary_data.ArrMemoryusagebytes, container_data.ArrMemoryusagebytes[i])
			summary_data.ArrMemoryworkingsetbytes = append(summary_data.ArrMemoryworkingsetbytes, container_data.ArrMemoryworkingsetbytes[i])
			summary_data.ArrMemorycache = append(summary_data.ArrMemorycache, container_data.ArrMemorycache[i])
			summary_data.ArrMemoryswap = append(summary_data.ArrMemoryswap, container_data.ArrMemoryswap[i])
			summary_data.ArrMemoryrss = append(summary_data.ArrMemoryrss, container_data.ArrMemoryrss[i])
			summary_data.ArrProcesses = append(summary_data.ArrProcesses, container_data.ArrProcesses[i])
			summary_data.ArrCpucount = append(summary_data.ArrCpucount, cpucount)
			summary_data.ArrMemorysize = append(summary_data.ArrMemorysize, memorysize)
			summary_data.ArrCpurequest = append(summary_data.ArrCpurequest, float64(cpurequest.(int64)))
			summary_data.ArrCpulimit = append(summary_data.ArrCpulimit, float64(cpulimit.(int64)))
			summary_data.ArrMemoryrequest = append(summary_data.ArrMemoryrequest, float64(memoryrequest.(int64)))
			summary_data.ArrMemorylimit = append(summary_data.ArrMemorylimit, float64(memorylimit.(int64)))
			idx++
		} else {
			summary_data.ArrCpuusagesecondstotal[val] += container_data.ArrCpuusagesecondstotal[i]
			summary_data.ArrCpusystemsecondstotal[val] += container_data.ArrCpusystemsecondstotal[i]
			summary_data.ArrCpuusersecondstotal[val] += container_data.ArrCpuusersecondstotal[i]
			summary_data.ArrMemoryusagebytes[val] += container_data.ArrMemoryusagebytes[i]
			summary_data.ArrMemoryworkingsetbytes[val] += container_data.ArrMemoryworkingsetbytes[i]
			summary_data.ArrMemorycache[val] += container_data.ArrMemorycache[i]
			summary_data.ArrMemoryswap[val] += container_data.ArrMemoryswap[i]
			summary_data.ArrMemoryrss[val] += container_data.ArrMemoryrss[i]
			summary_data.ArrProcesses[val] += container_data.ArrProcesses[i]
			summary_data.ArrCpucount[val] += float64(cpucount)
			summary_data.ArrMemorysize[val] += float64(memorysize)
			summary_data.ArrCpurequest[val] += float64(cpurequest.(int64))
			summary_data.ArrCpulimit[val] += float64(cpulimit.(int64))
			summary_data.ArrMemoryrequest[val] += float64(memoryrequest.(int64))
			summary_data.ArrMemorylimit[val] += float64(memorylimit.(int64))
		}
	}

	return summary_data
}

func PodNetSummaryPreprocessing(podnet_data *PodNetPerf) PodNetPerf {
	var summary_data PodNetPerf = PodNetPerf{}
	var podnetindex map[string]int = make(map[string]int)
	var idx int
	for i := 0; i < len(podnet_data.ArrOntunetime); i++ {
		var key string = podnet_data.ArrPoduid[i] + ":" + strconv.Itoa(podnet_data.ArrInterfaceid[i]) + ":" + podnet_data.ArrNsuid[i]
		if val, ok := podnetindex[key]; !ok {
			podnetindex[key] = idx
			summary_data.ArrPoduid = append(summary_data.ArrPoduid, podnet_data.ArrPoduid[i])
			summary_data.ArrOntunetime = append(summary_data.ArrOntunetime, podnet_data.ArrOntunetime[i])
			summary_data.ArrTimestampMs = append(summary_data.ArrTimestampMs, podnet_data.ArrTimestampMs[i])
			summary_data.ArrNsuid = append(summary_data.ArrNsuid, podnet_data.ArrNsuid[i])
			summary_data.ArrNodeuid = append(summary_data.ArrNodeuid, podnet_data.ArrNodeuid[i])
			summary_data.ArrMetricid = append(summary_data.ArrMetricid, 0)
			summary_data.ArrInterfaceid = append(summary_data.ArrInterfaceid, podnet_data.ArrInterfaceid[i])
			summary_data.ArrNetworkreceivebytestotal = append(summary_data.ArrNetworkreceivebytestotal, podnet_data.ArrNetworkreceivebytestotal[i])
			summary_data.ArrNetworkreceiveerrorstotal = append(summary_data.ArrNetworkreceiveerrorstotal, podnet_data.ArrNetworkreceiveerrorstotal[i])
			summary_data.ArrNetworktransmitbytestotal = append(summary_data.ArrNetworktransmitbytestotal, podnet_data.ArrNetworktransmitbytestotal[i])
			summary_data.ArrNetworktransmiterrorstotal = append(summary_data.ArrNetworktransmiterrorstotal, podnet_data.ArrNetworktransmiterrorstotal[i])
			idx++
		} else {
			summary_data.ArrNetworkreceivebytestotal[val] += podnet_data.ArrNetworkreceivebytestotal[i]
			summary_data.ArrNetworkreceiveerrorstotal[val] += podnet_data.ArrNetworkreceiveerrorstotal[i]
			summary_data.ArrNetworktransmitbytestotal[val] += podnet_data.ArrNetworktransmitbytestotal[i]
			summary_data.ArrNetworktransmiterrorstotal[val] += podnet_data.ArrNetworktransmiterrorstotal[i]
		}
	}

	return summary_data
}

func PodFsSummaryPreprocessing(podfs_data *PodFsPerf, containerfs_data *ContainerFsPerf) PodFsPerf {
	var summary_data PodFsPerf = PodFsPerf{}
	var podfsindex map[string]int = make(map[string]int)
	var idx int

	for i := 0; i < len(containerfs_data.ArrOntunetime); i++ {
		var key string = containerfs_data.ArrPoduid[i] + ":" + strconv.Itoa(containerfs_data.ArrDeviceid[i]) + ":" + containerfs_data.ArrNsuid[i]
		if val, ok := podfsindex[key]; !ok {
			podfsindex[key] = idx
			summary_data.ArrPoduid = append(summary_data.ArrPoduid, containerfs_data.ArrPoduid[i])
			summary_data.ArrOntunetime = append(summary_data.ArrOntunetime, containerfs_data.ArrOntunetime[i])
			summary_data.ArrTimestampMs = append(summary_data.ArrTimestampMs, containerfs_data.ArrTimestampMs[i])
			summary_data.ArrNsuid = append(summary_data.ArrNsuid, containerfs_data.ArrNsuid[i])
			summary_data.ArrNodeuid = append(summary_data.ArrNodeuid, containerfs_data.ArrNodeuid[i])
			summary_data.ArrMetricid = append(summary_data.ArrMetricid, 0)
			summary_data.ArrDeviceid = append(summary_data.ArrDeviceid, containerfs_data.ArrDeviceid[i])
			summary_data.ArrFsreadsbytestotal = append(summary_data.ArrFsreadsbytestotal, containerfs_data.ArrFsreadsbytestotal[i])
			summary_data.ArrFswritesbytestotal = append(summary_data.ArrFswritesbytestotal, containerfs_data.ArrFswritesbytestotal[i])
			idx++
		} else {
			summary_data.ArrFsreadsbytestotal[val] += containerfs_data.ArrFsreadsbytestotal[i]
			summary_data.ArrFswritesbytestotal[val] += containerfs_data.ArrFswritesbytestotal[i]
		}
	}

	return summary_data
}

func GetStatData(stat_type string, current float64, previous float64, divider float64) float64 {
	const PERCENT = 100
	switch stat_type {
	case "rate":
		return math.Round((current - previous) / float64(common.RateInterval))
	case "subtract":
		return math.Round(current - previous)
	case "usage":
		if divider == 0 {
			return 0
		}

		return math.Round((current - previous) / float64(common.RateInterval) / divider * PERCENT * float64(common.TemporaryMultiply))
	case "current_usage":
		if divider == 0 {
			return 0
		}

		return math.Round(current / divider * PERCENT * float64(common.TemporaryMultiply))
	case "core":
		return math.Round((current - previous) / float64(common.RateInterval) * float64(common.TemporaryMultiply))
	}

	return 0
}

func SetSummaryIndex(src SummaryPerfInterface, idx int) SummaryIndex {
	summary_key := SummaryIndex{}
	if val, ok := src.GetSummaryMap().Load("nodeuid"); ok {
		summary_key.NodeUID = val.([]string)[idx]
	}
	if val, ok := src.GetSummaryMap().Load("nsuid"); ok {
		summary_key.NsUID = val.([]string)[idx]
	}
	if val, ok := src.GetSummaryMap().Load("poduid"); ok {
		summary_key.PodUID = val.([]string)[idx]
	}
	if val, ok := src.GetSummaryMap().Load("containername"); ok {
		summary_key.ContainerName = val.([]string)[idx]
	}

	return summary_key
}

func SetSummaryData(summary_data *sync.Map, datatype string, current SummaryPerfInterface, previous SummaryPerfInterface, cidx int, pidx int) {
	summary_key := SetSummaryIndex(current, cidx)
	var metric_names []string
	if datatype == "fs" {
		metric_names = append(metric_names, fs_metric_names...)
	} else if datatype == "net" {
		metric_names = append(metric_names, net_metric_names...)
	} else {
		return
	}

	summary_key.DataType = SUMMARY_DATATYPE_CURR
	for _, mn := range metric_names {
		summary_key.MetricName = mn
		if val, ok := current.GetSummaryMap().Load(mn); ok {
			sdata, _ := summary_data.LoadOrStore(summary_key, float64(0))
			summary_data.Store(summary_key, sdata.(float64)+val.([]float64)[cidx])
		}
	}

	summary_key.DataType = SUMMARY_DATATYPE_PREV
	for _, mn := range metric_names {
		summary_key.MetricName = mn
		if val, ok := previous.GetSummaryMap().Load(mn); ok {
			sdata, _ := summary_data.LoadOrStore(summary_key, float64(0))
			summary_data.Store(summary_key, sdata.(float64)+val.([]float64)[pidx])
		}
	}
}

func GetSummaryData(summary_data *sync.Map, key SummaryIndex, datatype string, metricname string) float64 {
	key.DataType = datatype
	key.MetricName = metricname
	if val, ok := summary_data.Load(key); ok {
		return val.(float64)
	}

	return float64(0)
}

func SetPodSummary(summary_data *sync.Map, key SummaryIndex, datatype string, metricname string, value float64) {
	key.DataType = datatype

	v, _ := summary_data.LoadOrStore(key, &sync.Map{})
	value_map := v.(*sync.Map)
	value_map.Store(metricname, value)
	summary_data.Store(key, value_map)
}

func LoadAndInitMap(cm *sync.Map, pm *sync.Map, key string) (float64, float64) {
	var cval float64
	var pval float64

	if val, ok := cm.Load(key); ok {
		cval = val.(float64)
	} else {
		cval = 0
	}

	if val, ok := pm.Load(key); ok {
		pval = val.(float64)
	} else {
		pval = 0
	}

	return cval, pval
}

func CalculateNodeSummary(summary_map *sync.Map, nodeuid string) {
	node_cpu_map, _ := common.ResourceMap.Load("node_cpu")
	node_memory_map, _ := common.ResourceMap.Load("node_memory")
	cc, _ := node_cpu_map.(*sync.Map).LoadOrStore(nodeuid, 0)
	ms, _ := node_memory_map.(*sync.Map).LoadOrStore(nodeuid, int64(0))
	cpucount := cc.(int)
	memorysize := ms.(int64)
	cur, _ := summary_map.Load(SUMMARY_DATATYPE_CURR)
	prev, _ := summary_map.Load(SUMMARY_DATATYPE_PREV)
	calc, _ := summary_map.Load(SUMMARY_DATATYPE_CALC)

	current_map := cur.(*sync.Map)
	previous_map := prev.(*sync.Map)
	calculation_map := calc.(*sync.Map)

	cprocess, _ := LoadAndInitMap(current_map, previous_map, "process")
	ccpuusage, pcpuusage := LoadAndInitMap(current_map, previous_map, "cpuusage")
	ccpusystem, pcpusystem := LoadAndInitMap(current_map, previous_map, "cpusystem")
	ccpuuser, pcpuuser := LoadAndInitMap(current_map, previous_map, "cpuuser")
	cmemoryws, _ := LoadAndInitMap(current_map, previous_map, "memoryws")
	cmemoryswap, _ := LoadAndInitMap(current_map, previous_map, "memoryswap")

	cal_proc, _ := calculation_map.LoadOrStore("processcount", float64(0))
	calculation_map.Store("processcount", cal_proc.(float64)+cprocess)

	cal_cpuusage, _ := calculation_map.LoadOrStore("cpuusage", float64(0))
	calculation_map.Store("cpuusage", cal_cpuusage.(float64)+GetStatData("usage", ccpuusage, pcpuusage, float64(cpucount)))

	cal_cpusystem, _ := calculation_map.LoadOrStore("cpusystem", float64(0))
	calculation_map.Store("cpusystem", cal_cpusystem.(float64)+GetStatData("usage", ccpusystem, pcpusystem, float64(cpucount)))

	cal_cpuuser, _ := calculation_map.LoadOrStore("cpuuser", float64(0))
	calculation_map.Store("cpuuser", cal_cpuuser.(float64)+GetStatData("usage", ccpuuser, pcpuuser, float64(cpucount)))

	cal_cpuusagecores, _ := calculation_map.LoadOrStore("cpuusagecores", float64(0))
	calculation_map.Store("cpuusagecores", cal_cpuusagecores.(float64)+GetStatData("core", ccpuusage, pcpuusage, 0))

	cal_cputotalcores, _ := calculation_map.LoadOrStore("cputotalcores", float64(0))
	calculation_map.Store("cputotalcores", cal_cputotalcores.(float64)+math.Round(float64(cpucount)))

	cal_memoryusage, _ := calculation_map.LoadOrStore("memoryusage", float64(0))
	calculation_map.Store("memoryusage", cal_memoryusage.(float64)+GetStatData("current_usage", cmemoryws, 0, float64(memorysize)))

	cal_memoryusagebytes, _ := calculation_map.LoadOrStore("memoryusagebytes", float64(0))
	calculation_map.Store("memoryusagebytes", cal_memoryusagebytes.(float64)+math.Round(cmemoryws))

	cal_memorysizebytes, _ := calculation_map.LoadOrStore("memorysizebytes", float64(0))
	calculation_map.Store("memorysizebytes", cal_memorysizebytes.(float64)+math.Round(float64(memorysize)))

	cal_memoryswap, _ := calculation_map.LoadOrStore("memoryswap", float64(0))
	calculation_map.Store("memoryswap", cal_memoryswap.(float64)+GetStatData("current_usage", cmemoryswap, 0, float64(memorysize)))

	cnetrbt, pnetrbt := LoadAndInitMap(current_map, previous_map, "netrbt")
	cnetrbe, pnetrbe := LoadAndInitMap(current_map, previous_map, "netrbe")
	cnettbt, pnettbt := LoadAndInitMap(current_map, previous_map, "nettbt")
	cnettbe, pnettbe := LoadAndInitMap(current_map, previous_map, "nettbe")

	rcv_rate := GetStatData("rate", cnetrbt, pnetrbt, 0)
	rcv_errors := GetStatData("subtract", cnetrbe, pnetrbe, 0)
	trans_rate := GetStatData("rate", cnettbt, pnettbt, 0)
	trans_errors := GetStatData("rate", cnettbe, pnettbe, 0)

	cal_netrcvrate, _ := calculation_map.LoadOrStore("netrcvrate", float64(0))
	calculation_map.Store("netrcvrate", cal_netrcvrate.(float64)+rcv_rate)

	cal_netrcverrors, _ := calculation_map.LoadOrStore("netrcverrors", float64(0))
	calculation_map.Store("netrcverrors", cal_netrcverrors.(float64)+rcv_errors)

	cal_nettransrate, _ := calculation_map.LoadOrStore("nettransrate", float64(0))
	calculation_map.Store("nettransrate", cal_nettransrate.(float64)+trans_rate)

	cal_nettranserrors, _ := calculation_map.LoadOrStore("nettranserrors", float64(0))
	calculation_map.Store("nettranserrors", cal_nettranserrors.(float64)+trans_errors)

	cal_netiorate, _ := calculation_map.LoadOrStore("netiorate", float64(0))
	calculation_map.Store("netiorate", cal_netiorate.(float64)+rcv_rate+trans_rate)

	cal_netioerrors, _ := calculation_map.LoadOrStore("netioerrors", float64(0))
	calculation_map.Store("netioerrors", cal_netioerrors.(float64)+rcv_errors+trans_errors)

	cfsrbt, pfsrbt := LoadAndInitMap(current_map, previous_map, "fsrbt")
	cfswbt, pfswbt := LoadAndInitMap(current_map, previous_map, "fswbt")

	fs_read_rate := GetStatData("rate", cfsrbt, pfsrbt, 0)
	fs_write_rate := GetStatData("rate", cfswbt, pfswbt, 0)

	cal_fsreadrate, _ := calculation_map.LoadOrStore("fsreadrate", float64(0))
	calculation_map.Store("fsreadrate", cal_fsreadrate.(float64)+fs_read_rate)

	cal_fswriterate, _ := calculation_map.LoadOrStore("fswriterate", float64(0))
	calculation_map.Store("fswriterate", cal_fswriterate.(float64)+fs_write_rate)

	cal_fsiorate, _ := calculation_map.LoadOrStore("fsiorate", float64(0))
	calculation_map.Store("fsiorate", cal_fsiorate.(float64)+fs_read_rate+fs_write_rate)

	summary_map.Delete(SUMMARY_DATATYPE_CURR)
	summary_map.Delete(SUMMARY_DATATYPE_PREV)
}

func getParentObject(kind string, host string, nsname string) string {
	parent_host_list := map[string]struct{}{
		"namespace":        {},
		"node":             {},
		"persistentvolume": {},
		"storageclass":     {},
	}

	if _, ok := parent_host_list[strings.ToLower(kind)]; ok {
		return host
	} else {
		return nsname
	}
}

func GetPodReference(poduid string, kind string) (string, bool) {
	if pod_reference_map, ok := common.ResourceMap.Load("pod_reference"); !ok {
		return "", false
	} else {
		if refval, ok := pod_reference_map.(*sync.Map).Load(poduid); !ok {
			return "", false
		} else {
			refval_map := refval.(map[string]string)
			if _, ok := refval_map[kind]; !ok {
				return "", false
			} else {
				return refval_map[kind], true
			}
		}
	}
}

func ProcessWorkloadSummaryByNode(nodevalue *sync.Map, kind string) *sync.Map {
	var cal_map *sync.Map = &sync.Map{}
	nodevalue.Range(func(nuid, pmap any) bool {
		nodeuid := nuid.(string)
		pod_map := pmap.(*sync.Map)

		pod_map.Range(func(sumkey, pod_value any) bool {
			summary_index := sumkey.(SummaryIndex)

			wluid, flag := GetPodReference(summary_index.PodUID, kind)
			if !flag {
				return true
			}

			var wl_map *sync.Map
			if wl, ok := cal_map.LoadOrStore(wluid, &sync.Map{}); !ok {
				wl_map = wl.(*sync.Map)
				wl_map.Store(SUMMARY_DATATYPE_CALC, &sync.Map{})
				wl_map.Store(SUMMARY_DATATYPE_CURR, &sync.Map{})
				wl_map.Store(SUMMARY_DATATYPE_PREV, &sync.Map{})
				wl_map.Store(SUMMARY_DATATYPE_CNT, 0)
			} else {
				wl_map = wl.(*sync.Map)
				wl_map.LoadOrStore(SUMMARY_DATATYPE_CURR, &sync.Map{})
				wl_map.LoadOrStore(SUMMARY_DATATYPE_PREV, &sync.Map{})
			}

			if p, ok := wl_map.Load(SUMMARY_DATATYPE_CNT); ok && summary_index.DataType == SUMMARY_DATATYPE_CURR {
				wl_map.Store(SUMMARY_DATATYPE_CNT, p.(int)+1)
			}

			if cm, ok := wl_map.Load(summary_index.DataType); ok {
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

		cal_map.Range(func(wluid, wlvalue any) bool {
			wl_map := wlvalue.(*sync.Map)

			if _, ok := wl_map.Load(SUMMARY_DATATYPE_CURR); ok {
				CalculateNodeSummary(wl_map, nodeuid)
			}

			return true
		})

		return true
	})

	return cal_map
}

func GetResourceObjectUID(resource_data kubeapi.MappingEvent) string {
	var enabledrsc bool = false
	for _, rsc := range kubeapi.EnabledResources {
		if resource_data.ObjectKind == rsc {
			enabledrsc = true
			break
		}
	}

	if enabledrsc {
		return getUID(resource_data.ObjectKind, getParentObject(resource_data.ObjectKind, resource_data.Host, resource_data.NamespaceName), resource_data.ObjectName)
	} else {
		return resource_data.ObjectName
	}
}

func CopyUIDMap(uid_map map[string]struct{}) map[string]struct{} {
	uid_copy := map[string]struct{}{}
	for k, v := range uid_map {
		uid_copy[k] = v
	}

	return uid_copy
}

func RequestCheck(srcmap map[string]struct{}, dstmap map[string]struct{}, flag *bool) {
	if len(srcmap) != len(dstmap) && len(srcmap) > 0 && len(dstmap) > 0 {
		for k := range srcmap {
			if _, ok := dstmap[k]; !ok {
				*flag = true
				break
			}
		}

		for k := range dstmap {
			if _, ok := srcmap[k]; !ok {
				*flag = true
				break
			}
		}
	}
}

func MapNilCheck(srcmap *sync.Map) bool {
	var flag bool = false
	srcmap.Range(func(k, v any) bool {
		flag = true
		return false
	})

	return flag
}

func getStarttime(src_time int64, biastime int64) int64 {
	var starttime int64 = src_time - biastime
	if starttime > 0 {
		return starttime
	} else {
		return 0
	}
}

func refreshOntunetime(arr *[]int64, ontunetime int64) {
	for i := 0; i < len(*arr); i++ {
		(*arr)[i] = ontunetime
	}
}

func StoreClusterMap(cluster_map *sync.Map, pod_map *sync.Map, nodeuid string) {
	node_cluster_map, _ := common.ResourceMap.Load("node_cluster")
	cid, _ := node_cluster_map.(*sync.Map).LoadOrStore(nodeuid, 0)
	clusterid := cid.(int)

	if MapNilCheck(pod_map) {
		cm, _ := cluster_map.LoadOrStore(clusterid, &sync.Map{})
		cm_node_map := cm.(*sync.Map)
		cm_node_map.Store(nodeuid, pod_map)
		cluster_map.Store(clusterid, cm_node_map)
	}
}

func getUidHost(uidhostinfo map[string]string) (string, string) {
	var uid string
	var hostip string

	for u, h := range uidhostinfo {
		if uid == "" {
			uid = "('" + u + "'"
			hostip = h
		} else {
			uid = uid + ",'" + u + "'"
		}
	}
	uid = uid + ")"

	return uid, hostip
}

func getDataQueueMapHost(data_queue_map *sync.Map) []string {
	var hostnames []string = make([]string, 0)

	data_queue_map.Range(func(key, value any) bool {
		var hostname string
		metric_arr := value.([]*sync.Map)

		if len(metric_arr) == 0 {
			return true
		}

		// Get Recent Queue Value
		metric_data := metric_arr[len(metric_arr)-1]

		if hn, ok := metric_data.Load("hostname"); ok {
			hostname = hn.(string)
		}

		var duplicate_flag bool = false
		for i := 0; i < len(hostnames); i++ {
			if hostnames[i] == hostname {
				duplicate_flag = true
				break
			}
		}

		if !duplicate_flag {
			hostnames = append(hostnames, hostname)
		}

		return true
	})

	return hostnames
}

func initMapUidInfo(data_queue_map *sync.Map) {
	hostnames := getDataQueueMapHost(data_queue_map)

	for _, hostname := range hostnames {
		if _, ok := mapUidInfo[hostname]; !ok {
			mapUidInfo[hostname] = MapUIDInfo{
				currentNodeUid:      make(map[string]struct{}),
				currentPodUid:       make(map[string]struct{}),
				currentNamespaceUid: make(map[string]struct{}),
				prevNodeUid:         make(map[string]struct{}),
				prevPodUid:          make(map[string]struct{}),
				prevNamespaceUid:    make(map[string]struct{}),
			}
		}

		if hostMapUidInfo, ok := mapUidInfo[hostname]; ok {
			hostMapUidInfo.prevNodeUid = CopyUIDMap(mapUidInfo[hostname].currentNodeUid)
			hostMapUidInfo.prevPodUid = CopyUIDMap(mapUidInfo[hostname].currentPodUid)
			hostMapUidInfo.prevNamespaceUid = CopyUIDMap(mapUidInfo[hostname].currentNamespaceUid)
			hostMapUidInfo.currentNodeUid = make(map[string]struct{})
			hostMapUidInfo.currentPodUid = make(map[string]struct{})
			hostMapUidInfo.currentNamespaceUid = make(map[string]struct{})

			mapUidInfo[hostname] = hostMapUidInfo
		}
	}
}

func checkMapUidInfo(data_queue_map *sync.Map) {
	hostnames := getDataQueueMapHost(data_queue_map)

	for _, hostname := range hostnames {
		var requestflag bool = false

		// UIDMap Compare - Node
		common.MutexNode.Lock()
		RequestCheck(mapUidInfo[hostname].prevNodeUid, mapUidInfo[hostname].currentNodeUid, &requestflag)
		common.MutexNode.Unlock()

		// UIDMap Compare - Pod
		common.MutexPod.Lock()
		RequestCheck(mapUidInfo[hostname].prevPodUid, mapUidInfo[hostname].currentPodUid, &requestflag)
		common.MutexPod.Unlock()

		// UIDMap Compare - Namespace
		common.MutexNs.Lock()
		RequestCheck(mapUidInfo[hostname].prevNamespaceUid, mapUidInfo[hostname].currentNamespaceUid, &requestflag)
		common.MutexNs.Unlock()

		if requestflag {
			common.ChannelRequestChangeHost <- hostname
		}
	}
}
