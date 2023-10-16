package database

import (
	"context"
	"fmt"
	"math"
	"onTuneKubeManager/common"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func setNodeRealtimeperf(current_data *sync.Map, prev_data *sync.Map, net_map *sync.Map, fs_map *sync.Map, nodeuid string) *NodePerfStat {
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

			current_net_rbt := getSummaryData(net_map, summary_key, "current", "rbt")
			current_net_rbe := getSummaryData(net_map, summary_key, "current", "rbe")
			current_net_tbt := getSummaryData(net_map, summary_key, "current", "tbt")
			current_net_tbe := getSummaryData(net_map, summary_key, "current", "tbe")
			previous_net_rbt := getSummaryData(net_map, summary_key, "previous", "rbt")
			previous_net_rbe := getSummaryData(net_map, summary_key, "previous", "rbe")
			previous_net_tbt := getSummaryData(net_map, summary_key, "previous", "tbt")
			previous_net_tbe := getSummaryData(net_map, summary_key, "previous", "tbe")
			rcv_rate := getStatData("rate", current_net_rbt, previous_net_rbt, 0)
			rcv_errors := getStatData("subtract", current_net_rbe, previous_net_rbe, 0)
			trans_rate := getStatData("rate", current_net_tbt, previous_net_tbt, 0)
			trans_errors := getStatData("subtract", current_net_tbe, previous_net_tbe, 0)

			insert_stat_data.ArrNetworkreceiverate = append(insert_stat_data.ArrNetworkreceiverate, rcv_rate)
			insert_stat_data.ArrNetworkreceiveerrors = append(insert_stat_data.ArrNetworkreceiveerrors, rcv_errors)
			insert_stat_data.ArrNetworktransmitrate = append(insert_stat_data.ArrNetworktransmitrate, trans_rate)
			insert_stat_data.ArrNetworktransmiterrors = append(insert_stat_data.ArrNetworktransmiterrors, trans_errors)
			insert_stat_data.ArrNetworkiorate = append(insert_stat_data.ArrNetworkiorate, rcv_rate+trans_rate)
			insert_stat_data.ArrNetworkioerrors = append(insert_stat_data.ArrNetworkioerrors, rcv_errors+trans_errors)

			current_fs_rbt := getSummaryData(fs_map, summary_key, "current", "rbt")
			current_fs_wbt := getSummaryData(fs_map, summary_key, "current", "wbt")
			previous_fs_rbt := getSummaryData(fs_map, summary_key, "previous", "rbt")
			previous_fs_wbt := getSummaryData(fs_map, summary_key, "previous", "wbt")
			fs_read_rate := getStatData("rate", current_fs_rbt, previous_fs_rbt, 0)
			fs_write_rate := getStatData("rate", current_fs_wbt, previous_fs_wbt, 0)

			insert_stat_data.ArrFsreadrate = append(insert_stat_data.ArrFsreadrate, fs_read_rate)
			insert_stat_data.ArrFswriterate = append(insert_stat_data.ArrFswriterate, fs_write_rate)
			insert_stat_data.ArrFsiorate = append(insert_stat_data.ArrFsiorate, fs_read_rate+fs_write_rate)

			insert_stat_data.ArrCpuusage = append(insert_stat_data.ArrCpuusage, getStatData("usage", current_node_data.ArrCpuusagesecondstotal[i], prev_node_data.ArrCpuusagesecondstotal[pidx], float64(cpucount)))
			insert_stat_data.ArrCpusystem = append(insert_stat_data.ArrCpusystem, getStatData("usage", current_node_data.ArrCpusystemsecondstotal[i], prev_node_data.ArrCpusystemsecondstotal[pidx], float64(cpucount)))
			insert_stat_data.ArrCpuuser = append(insert_stat_data.ArrCpuuser, getStatData("usage", current_node_data.ArrCpuusersecondstotal[i], prev_node_data.ArrCpuusersecondstotal[pidx], float64(cpucount)))
			insert_stat_data.ArrCpuusagecores = append(insert_stat_data.ArrCpuusagecores, getStatData("core", current_node_data.ArrCpuusagesecondstotal[i], prev_node_data.ArrCpuusagesecondstotal[pidx], 0))
			insert_stat_data.ArrCputotalcores = append(insert_stat_data.ArrCputotalcores, float64(cpucount))
			insert_stat_data.ArrMemoryusage = append(insert_stat_data.ArrMemoryusage, getStatData("current_usage", current_node_data.ArrMemoryworkingsetbytes[i], 0, float64(memorysize)))
			insert_stat_data.ArrMemoryusagebytes = append(insert_stat_data.ArrMemoryusagebytes, math.Round(current_node_data.ArrMemoryworkingsetbytes[i]))
			insert_stat_data.ArrMemorysizebytes = append(insert_stat_data.ArrMemorysizebytes, math.Round(float64(memorysize)))
			insert_stat_data.ArrMemoryswap = append(insert_stat_data.ArrMemoryswap, getStatData("current_usage", current_node_data.ArrMemoryswap[i], 0, float64(memorysize)))

			break
		}
	}

	lsdm, _ := last_stat_data_map.LoadOrStore(METRIC_VAR_NODE, &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store(METRIC_VAR_NODE, lsdm)

	return insert_stat_data
}

func insertNodeRealtimeperf(stat_node_map map[string]NodePerfStat, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_NODE_PERF)
	var lasttablename string = TB_KUBE_LAST_NODE_PERF

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, insert_stat_data := range stat_node_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODE_STAT, tablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI Node Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	// Insert Last Data
	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(TRUNCATE_DATA, lasttablename))
	if !errorCheck(del_err) {
		return
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastNode Delete %d", processid, time.Now().UnixNano()))

	for _, insert_stat_data := range stat_node_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODE_STAT, lasttablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastNode Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func setPodRealtimeperf(current_data *sync.Map, prev_data *sync.Map, net_map *sync.Map, fs_map *sync.Map, nodeuid string) (*PodPerfStat, *sync.Map) {
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
		current_pod_summary_data = podSummaryPreprocessing(&current_pod_data, &current_container_data)
	} else {
		return &PodPerfStat{}, &sync.Map{}
	}

	if p, ok := prev_data.Load("pod_insert"); ok {
		prev_pod_data := p.(PodPerf)

		var prev_container_data ContainerPerf
		if pc, ok2 := prev_data.Load("container_insert"); ok2 {
			prev_container_data = pc.(ContainerPerf)
		}
		prev_pod_summary_data = podSummaryPreprocessing(&prev_pod_data, &prev_container_data)
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
			NsUID:  current_pod_summary_data.ArrNsUid[i],
			PodUID: current_pod_summary_data.ArrPoduid[i],
		}

		insert_stat_data.ArrPoduid = append(insert_stat_data.ArrPoduid, current_pod_summary_data.ArrPoduid[i])
		insert_stat_data.ArrNodeuid = append(insert_stat_data.ArrNodeuid, current_pod_summary_data.ArrNodeuid[i])
		insert_stat_data.ArrNsUid = append(insert_stat_data.ArrNsUid, current_pod_summary_data.ArrNsUid[i])
		insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_pod_summary_data.ArrOntunetime[i])
		insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, 0)
		insert_stat_data.ArrProcesses = append(insert_stat_data.ArrProcesses, current_pod_summary_data.ArrProcesses[i])
		setPodSummary(summary_data, summary_key, "current", "process", current_pod_summary_data.ArrProcesses[i])
		setPodSummary(summary_data, summary_key, "previous", "process", prev_pod_summary_data.ArrProcesses[pidx])

		insert_stat_data.ArrCpuusage = append(insert_stat_data.ArrCpuusage, getStatData("usage", current_pod_summary_data.ArrCpuusagesecondstotal[i], prev_pod_summary_data.ArrCpuusagesecondstotal[pidx], current_pod_summary_data.ArrCpucount[i]))
		insert_stat_data.ArrCpusystem = append(insert_stat_data.ArrCpusystem, getStatData("usage", current_pod_summary_data.ArrCpusystemsecondstotal[i], prev_pod_summary_data.ArrCpusystemsecondstotal[pidx], current_pod_summary_data.ArrCpucount[i]))
		insert_stat_data.ArrCpuuser = append(insert_stat_data.ArrCpuuser, getStatData("usage", current_pod_summary_data.ArrCpuusersecondstotal[i], prev_pod_summary_data.ArrCpuusersecondstotal[pidx], current_pod_summary_data.ArrCpucount[i]))
		insert_stat_data.ArrCpuusagecores = append(insert_stat_data.ArrCpuusagecores, getStatData("core", current_pod_summary_data.ArrCpuusagesecondstotal[i], prev_pod_summary_data.ArrCpuusagesecondstotal[pidx], 0))
		insert_stat_data.ArrCpurequestcores = append(insert_stat_data.ArrCpurequestcores, current_pod_summary_data.ArrCpurequest[i])
		insert_stat_data.ArrCpulimitcores = append(insert_stat_data.ArrCpulimitcores, current_pod_summary_data.ArrCpulimit[i])
		insert_stat_data.ArrMemoryusage = append(insert_stat_data.ArrMemoryusage, getStatData("current_usage", current_pod_summary_data.ArrMemoryworkingsetbytes[i], 0, current_pod_summary_data.ArrMemorysize[i]))
		insert_stat_data.ArrMemoryusagebytes = append(insert_stat_data.ArrMemoryusagebytes, math.Round(current_pod_summary_data.ArrMemoryworkingsetbytes[i]))
		insert_stat_data.ArrMemoryrequestbytes = append(insert_stat_data.ArrMemoryrequestbytes, math.Round(current_pod_summary_data.ArrMemoryrequest[i]))
		insert_stat_data.ArrMemorylimitbytes = append(insert_stat_data.ArrMemorylimitbytes, math.Round(current_pod_summary_data.ArrMemorylimit[i]))
		insert_stat_data.ArrMemoryswap = append(insert_stat_data.ArrMemoryswap, getStatData("current_usage", current_pod_summary_data.ArrMemoryswap[i], 0, current_pod_summary_data.ArrMemorysize[i]))

		setPodSummary(summary_data, summary_key, "current", "cpuusage", current_pod_summary_data.ArrCpuusagesecondstotal[i])
		setPodSummary(summary_data, summary_key, "current", "cpusystem", current_pod_summary_data.ArrCpusystemsecondstotal[i])
		setPodSummary(summary_data, summary_key, "current", "cpuuser", current_pod_summary_data.ArrCpuusersecondstotal[i])
		setPodSummary(summary_data, summary_key, "current", "cpucount", current_pod_summary_data.ArrCpucount[i])
		setPodSummary(summary_data, summary_key, "current", "cpurequest", current_pod_summary_data.ArrCpurequest[i])
		setPodSummary(summary_data, summary_key, "current", "cpulimit", current_pod_summary_data.ArrCpulimit[i])
		setPodSummary(summary_data, summary_key, "current", "memoryws", current_pod_summary_data.ArrMemoryworkingsetbytes[i])
		setPodSummary(summary_data, summary_key, "current", "memoryswap", current_pod_summary_data.ArrMemoryswap[i])
		setPodSummary(summary_data, summary_key, "current", "memorysize", current_pod_summary_data.ArrMemorysize[i])
		setPodSummary(summary_data, summary_key, "current", "memoryrequest", current_pod_summary_data.ArrMemoryrequest[i])
		setPodSummary(summary_data, summary_key, "current", "memorylimit", current_pod_summary_data.ArrMemorylimit[i])

		setPodSummary(summary_data, summary_key, "previous", "cpuusage", prev_pod_summary_data.ArrCpuusagesecondstotal[pidx])
		setPodSummary(summary_data, summary_key, "previous", "cpusystem", prev_pod_summary_data.ArrCpusystemsecondstotal[pidx])
		setPodSummary(summary_data, summary_key, "previous", "cpuuser", prev_pod_summary_data.ArrCpuusersecondstotal[pidx])
		setPodSummary(summary_data, summary_key, "previous", "memoryws", prev_pod_summary_data.ArrMemoryworkingsetbytes[pidx])
		setPodSummary(summary_data, summary_key, "previous", "memoryswap", prev_pod_summary_data.ArrMemoryswap[pidx])

		current_net_rbt := getSummaryData(net_map, summary_key, "current", "rbt")
		current_net_rbe := getSummaryData(net_map, summary_key, "current", "rbe")
		current_net_tbt := getSummaryData(net_map, summary_key, "current", "tbt")
		current_net_tbe := getSummaryData(net_map, summary_key, "current", "tbe")
		previous_net_rbt := getSummaryData(net_map, summary_key, "previous", "rbt")
		previous_net_rbe := getSummaryData(net_map, summary_key, "previous", "rbe")
		previous_net_tbt := getSummaryData(net_map, summary_key, "previous", "tbt")
		previous_net_tbe := getSummaryData(net_map, summary_key, "previous", "tbe")
		rcv_rate := getStatData("rate", current_net_rbt, previous_net_rbt, 0)
		rcv_errors := getStatData("subtract", current_net_rbe, previous_net_rbe, 0)
		trans_rate := getStatData("rate", current_net_tbt, previous_net_tbt, 0)
		trans_errors := getStatData("subtract", current_net_tbe, previous_net_tbe, 0)
		setPodSummary(summary_data, summary_key, "current", "netrbt", current_net_rbt)
		setPodSummary(summary_data, summary_key, "current", "netrbe", current_net_rbe)
		setPodSummary(summary_data, summary_key, "current", "nettbt", current_net_tbt)
		setPodSummary(summary_data, summary_key, "current", "nettbe", current_net_tbe)
		setPodSummary(summary_data, summary_key, "previous", "netrbt", previous_net_rbt)
		setPodSummary(summary_data, summary_key, "previous", "netrbe", previous_net_rbe)
		setPodSummary(summary_data, summary_key, "previous", "nettbt", previous_net_tbt)
		setPodSummary(summary_data, summary_key, "previous", "nettbe", previous_net_tbe)

		insert_stat_data.ArrNetworkreceiverate = append(insert_stat_data.ArrNetworkreceiverate, rcv_rate)
		insert_stat_data.ArrNetworkreceiveerrors = append(insert_stat_data.ArrNetworkreceiveerrors, rcv_errors)
		insert_stat_data.ArrNetworktransmitrate = append(insert_stat_data.ArrNetworktransmitrate, trans_rate)
		insert_stat_data.ArrNetworktransmiterrors = append(insert_stat_data.ArrNetworktransmiterrors, trans_errors)
		insert_stat_data.ArrNetworkiorate = append(insert_stat_data.ArrNetworkiorate, rcv_rate+trans_rate)
		insert_stat_data.ArrNetworkioerrors = append(insert_stat_data.ArrNetworkioerrors, rcv_errors+trans_errors)

		current_fs_rbt := getSummaryData(fs_map, summary_key, "current", "rbt")
		current_fs_wbt := getSummaryData(fs_map, summary_key, "current", "wbt")
		previous_fs_rbt := getSummaryData(fs_map, summary_key, "previous", "rbt")
		previous_fs_wbt := getSummaryData(fs_map, summary_key, "previous", "wbt")
		fs_read_rate := getStatData("rate", current_fs_rbt, previous_fs_rbt, 0)
		fs_write_rate := getStatData("rate", current_fs_wbt, previous_fs_wbt, 0)
		setPodSummary(summary_data, summary_key, "current", "fsrbt", current_fs_rbt)
		setPodSummary(summary_data, summary_key, "current", "fswbt", current_fs_wbt)
		setPodSummary(summary_data, summary_key, "previous", "fsrbt", previous_fs_rbt)
		setPodSummary(summary_data, summary_key, "previous", "fswbt", previous_fs_wbt)

		insert_stat_data.ArrFsreadrate = append(insert_stat_data.ArrFsreadrate, fs_read_rate)
		insert_stat_data.ArrFswriterate = append(insert_stat_data.ArrFswriterate, fs_write_rate)
		insert_stat_data.ArrFsiorate = append(insert_stat_data.ArrFsiorate, fs_read_rate+fs_write_rate)
	}

	lsdm, _ := last_stat_data_map.LoadOrStore(METRIC_VAR_POD, &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store(METRIC_VAR_POD, lsdm)
	last_summary_data_map.Store(nodeuid, summary_data)

	return insert_stat_data, summary_data
}

func insertPodRealtimeperf(stat_pod_map map[string]PodPerfStat, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_POD_PERF)
	var lasttablename string = TB_KUBE_LAST_POD_PERF

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, insert_stat_data := range stat_pod_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_POD_STAT, tablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI Pod Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	// Insert Last Data
	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(TRUNCATE_DATA, lasttablename))
	if !errorCheck(del_err) {
		return
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastPod Delete %d", processid, time.Now().UnixNano()))

	for _, insert_stat_data := range stat_pod_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_POD_STAT, lasttablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastPod Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func setContainerRealtimeperf(current_data *sync.Map, prev_data *sync.Map, fs_map *sync.Map, nodeuid string) *ContainerPerfStat {
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
		insert_stat_data.ArrNsUid = append(insert_stat_data.ArrNsUid, current_container_data.ArrNsUid[i])
		insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_container_data.ArrOntunetime[i])
		insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, current_container_data.ArrMetricid[i])
		insert_stat_data.ArrProcesses = append(insert_stat_data.ArrProcesses, current_container_data.ArrProcesses[i])

		insert_stat_data.ArrCpuusage = append(insert_stat_data.ArrCpuusage, getStatData("usage", current_container_data.ArrCpuusagesecondstotal[i], prev_container_data.ArrCpuusagesecondstotal[pidx], cpucount))
		insert_stat_data.ArrCpusystem = append(insert_stat_data.ArrCpusystem, getStatData("usage", current_container_data.ArrCpusystemsecondstotal[i], prev_container_data.ArrCpusystemsecondstotal[pidx], cpucount))
		insert_stat_data.ArrCpuuser = append(insert_stat_data.ArrCpuuser, getStatData("usage", current_container_data.ArrCpuusersecondstotal[i], prev_container_data.ArrCpuusersecondstotal[pidx], cpucount))
		insert_stat_data.ArrCpuusagecores = append(insert_stat_data.ArrCpuusagecores, getStatData("core", current_container_data.ArrCpuusagesecondstotal[i], prev_container_data.ArrCpuusagesecondstotal[pidx], 0))
		insert_stat_data.ArrCpurequestcores = append(insert_stat_data.ArrCpurequestcores, float64(cpurequest.(int64)))
		insert_stat_data.ArrCpulimitcores = append(insert_stat_data.ArrCpulimitcores, float64(cpulimit.(int64)))
		insert_stat_data.ArrMemoryusage = append(insert_stat_data.ArrMemoryusage, getStatData("current_usage", current_container_data.ArrMemoryworkingsetbytes[i], 0, float64(memorysize)))
		insert_stat_data.ArrMemoryusagebytes = append(insert_stat_data.ArrMemoryusagebytes, math.Round(current_container_data.ArrMemoryworkingsetbytes[i]))
		insert_stat_data.ArrMemoryrequestbytes = append(insert_stat_data.ArrMemoryrequestbytes, math.Round(float64(memoryrequest.(int64))))
		insert_stat_data.ArrMemorylimitbytes = append(insert_stat_data.ArrMemorylimitbytes, math.Round(float64(memorylimit.(int64))))
		insert_stat_data.ArrMemoryswap = append(insert_stat_data.ArrMemoryswap, getStatData("current_usage", current_container_data.ArrMemoryswap[i], 0, memorysize))

		summary_key := SummaryIndex{
			ContainerName: current_container_data.ArrContainername[i],
			PodUID:        current_container_data.ArrPoduid[i],
		}

		current_fs_rbt := getSummaryData(fs_map, summary_key, "current", "rbt")
		current_fs_wbt := getSummaryData(fs_map, summary_key, "current", "wbt")
		previous_fs_rbt := getSummaryData(fs_map, summary_key, "previous", "rbt")
		previous_fs_wbt := getSummaryData(fs_map, summary_key, "previous", "wbt")
		fs_read_rate := getStatData("rate", current_fs_rbt, previous_fs_rbt, 0)
		fs_write_rate := getStatData("rate", current_fs_wbt, previous_fs_wbt, 0)

		insert_stat_data.ArrFsreadrate = append(insert_stat_data.ArrFsreadrate, fs_read_rate)
		insert_stat_data.ArrFswriterate = append(insert_stat_data.ArrFswriterate, fs_write_rate)
		insert_stat_data.ArrFsiorate = append(insert_stat_data.ArrFsiorate, fs_read_rate+fs_write_rate)
	}

	lsdm, _ := last_stat_data_map.LoadOrStore(METRIC_VAR_CONTAINER, &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store(METRIC_VAR_CONTAINER, lsdm)

	return insert_stat_data
}

func insertContainerRealtimeperf(stat_container_map map[string]ContainerPerfStat, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_CONTAINER_PERF)
	var lasttablename string = TB_KUBE_LAST_CONTAINER_PERF

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, insert_stat_data := range stat_container_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CONTAINER_STAT, tablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI Container Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	// Insert Last Data
	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(TRUNCATE_DATA, lasttablename))
	if !errorCheck(del_err) {
		return
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastContainer Delete %d", processid, time.Now().UnixNano()))

	for _, insert_stat_data := range stat_container_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CONTAINER_STAT, lasttablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastContainer Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func setNodeNetRealtimeperf(current_data *sync.Map, prev_data *sync.Map, nodeuid string) (*NodeNetPerfStat, *sync.Map) {
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

			setSummaryData(summary_data, "net", &current_nodenet_data, &prev_nodenet_data, i, pidx)

			insert_stat_data.ArrNodeuid = append(insert_stat_data.ArrNodeuid, current_nodenet_data.ArrNodeuid[i])
			insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_nodenet_data.ArrOntunetime[i])
			insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, current_nodenet_data.ArrMetricid[i])
			insert_stat_data.ArrInterfaceid = append(insert_stat_data.ArrInterfaceid, current_nodenet_data.ArrInterfaceid[i])

			rcv_rate := getStatData("rate", current_nodenet_data.ArrNetworkreceivebytestotal[i], prev_nodenet_data.ArrNetworkreceivebytestotal[pidx], 0)
			rcv_errors := getStatData("subtract", current_nodenet_data.ArrNetworkreceiveerrorstotal[i], prev_nodenet_data.ArrNetworkreceiveerrorstotal[pidx], 0)
			trans_rate := getStatData("rate", current_nodenet_data.ArrNetworktransmitbytestotal[i], prev_nodenet_data.ArrNetworktransmitbytestotal[pidx], 0)
			trans_errors := getStatData("subtract", current_nodenet_data.ArrNetworktransmiterrorstotal[i], prev_nodenet_data.ArrNetworktransmiterrorstotal[pidx], 0)

			insert_stat_data.ArrNetworkreceiverate = append(insert_stat_data.ArrNetworkreceiverate, rcv_rate)
			insert_stat_data.ArrNetworkreceiveerrors = append(insert_stat_data.ArrNetworkreceiveerrors, rcv_errors)
			insert_stat_data.ArrNetworktransmitrate = append(insert_stat_data.ArrNetworktransmitrate, trans_rate)
			insert_stat_data.ArrNetworktransmiterrors = append(insert_stat_data.ArrNetworktransmiterrors, trans_errors)

			insert_stat_data.ArrNetworkiorate = append(insert_stat_data.ArrNetworkiorate, rcv_rate+trans_rate)
			insert_stat_data.ArrNetworkioerrors = append(insert_stat_data.ArrNetworkioerrors, rcv_errors+trans_errors)
		}
	}

	lsdm, _ := last_stat_data_map.LoadOrStore(METRIC_VAR_NODENET, &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store(METRIC_VAR_NODENET, lsdm)

	return insert_stat_data, summary_data
}

func insertNodeNetRealtimeperf(stat_nodenet_map map[string]NodeNetPerfStat, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_NODE_NET_PERF)
	var lasttablename string = TB_KUBE_LAST_NODE_NET_PERF

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, stat_data := range stat_nodenet_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODENET_STAT, tablename), stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI NodeNet Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	// Insert Last Data
	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(TRUNCATE_DATA, lasttablename))
	if !errorCheck(del_err) {
		return
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastNodeNet Delete %d", processid, time.Now().UnixNano()))

	for _, stat_data := range stat_nodenet_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODENET_STAT, lasttablename), stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastNodeNet Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func setPodNetRealtimeperf(current_data *sync.Map, prev_data *sync.Map, nodeuid string) (*PodNetPerfStat, *sync.Map) {
	_, result := getStatTimes(current_data)
	if !result {
		return &PodNetPerfStat{}, &sync.Map{}
	}

	var current_podnet_summary_data PodNetPerf
	var prev_podnet_summary_data PodNetPerf
	var summary_data *sync.Map = &sync.Map{}

	if c, ok := current_data.Load("podnet_insert"); ok {
		current_podnet_data := c.(PodNetPerf)
		current_podnet_summary_data = podNetSummaryPreprocessing(&current_podnet_data)
	} else {
		return &PodNetPerfStat{}, &sync.Map{}
	}

	if p, ok := prev_data.Load("podnet_insert"); ok {
		prev_podnet_data := p.(PodNetPerf)
		prev_podnet_summary_data = podNetSummaryPreprocessing(&prev_podnet_data)
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

		setSummaryData(summary_data, "net", &current_podnet_summary_data, &prev_podnet_summary_data, i, pidx)

		insert_stat_data.ArrPoduid = append(insert_stat_data.ArrPoduid, current_podnet_summary_data.ArrPoduid[i])
		insert_stat_data.ArrNodeuid = append(insert_stat_data.ArrNodeuid, current_podnet_summary_data.ArrNodeuid[i])
		insert_stat_data.ArrNsUid = append(insert_stat_data.ArrNsUid, current_podnet_summary_data.ArrNsUid[i])
		insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_podnet_summary_data.ArrOntunetime[i])
		insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, 0)
		insert_stat_data.ArrInterfaceid = append(insert_stat_data.ArrInterfaceid, current_podnet_summary_data.ArrInterfaceid[i])

		rcv_rate := getStatData("rate", current_podnet_summary_data.ArrNetworkreceivebytestotal[i], prev_podnet_summary_data.ArrNetworkreceivebytestotal[pidx], 0)
		rcv_errors := getStatData("subtract", current_podnet_summary_data.ArrNetworkreceiveerrorstotal[i], prev_podnet_summary_data.ArrNetworkreceiveerrorstotal[pidx], 0)
		trans_rate := getStatData("rate", current_podnet_summary_data.ArrNetworktransmitbytestotal[i], prev_podnet_summary_data.ArrNetworktransmitbytestotal[pidx], 0)
		trans_errors := getStatData("subtract", current_podnet_summary_data.ArrNetworktransmiterrorstotal[i], prev_podnet_summary_data.ArrNetworktransmiterrorstotal[pidx], 0)

		insert_stat_data.ArrNetworkreceiverate = append(insert_stat_data.ArrNetworkreceiverate, rcv_rate)
		insert_stat_data.ArrNetworkreceiveerrors = append(insert_stat_data.ArrNetworkreceiveerrors, rcv_errors)
		insert_stat_data.ArrNetworktransmitrate = append(insert_stat_data.ArrNetworktransmitrate, trans_rate)
		insert_stat_data.ArrNetworktransmiterrors = append(insert_stat_data.ArrNetworktransmiterrors, trans_errors)

		insert_stat_data.ArrNetworkiorate = append(insert_stat_data.ArrNetworkiorate, rcv_rate+trans_rate)
		insert_stat_data.ArrNetworkioerrors = append(insert_stat_data.ArrNetworkioerrors, rcv_errors+trans_errors)
	}

	lsdm, _ := last_stat_data_map.LoadOrStore(METRIC_VAR_PODNET, &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store(METRIC_VAR_PODNET, lsdm)

	return insert_stat_data, summary_data
}

func insertPodNetRealtimeperf(stat_podnet_map map[string]PodNetPerfStat, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_POD_NET_PERF)
	var lasttablename string = TB_KUBE_LAST_POD_NET_PERF

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, stat_data := range stat_podnet_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_PODNET_STAT, tablename), stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI PodNet Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	// Insert Last Data
	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(TRUNCATE_DATA, lasttablename))
	if !errorCheck(del_err) {
		return
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastPodNet Delete %d", processid, time.Now().UnixNano()))

	for _, stat_data := range stat_podnet_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_PODNET_STAT, lasttablename), stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastPodNet Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func setNodeFsRealtimeperf(current_data *sync.Map, prev_data *sync.Map, nodeuid string) (*NodeFsPerfStat, *sync.Map) {
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

			setSummaryData(summary_data, "fs", &current_nodefs_data, &prev_nodefs_data, i, pidx)

			insert_stat_data.ArrNodeuid = append(insert_stat_data.ArrNodeuid, current_nodefs_data.ArrNodeuid[i])
			insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_nodefs_data.ArrOntunetime[i])
			insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, current_nodefs_data.ArrMetricid[i])
			insert_stat_data.ArrDeviceid = append(insert_stat_data.ArrDeviceid, current_nodefs_data.ArrDeviceid[i])

			read_rate := getStatData("rate", current_nodefs_data.ArrFsreadsbytestotal[i], prev_nodefs_data.ArrFsreadsbytestotal[pidx], 0)
			write_rate := getStatData("rate", current_nodefs_data.ArrFswritesbytestotal[i], prev_nodefs_data.ArrFswritesbytestotal[pidx], 0)

			insert_stat_data.ArrFsreadrate = append(insert_stat_data.ArrFsreadrate, read_rate)
			insert_stat_data.ArrFswriterate = append(insert_stat_data.ArrFswriterate, write_rate)
			insert_stat_data.ArrFsiorate = append(insert_stat_data.ArrFsiorate, read_rate+write_rate)
		}
	}

	lsdm, _ := last_stat_data_map.LoadOrStore(METRIC_VAR_NODEFS, &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store(METRIC_VAR_NODEFS, lsdm)

	return insert_stat_data, summary_data
}

func insertNodeFsRealtimeperf(stat_nodefs_map map[string]NodeFsPerfStat, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_NODE_FS_PERF)
	var lasttablename string = TB_KUBE_LAST_NODE_FS_PERF

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, stat_data := range stat_nodefs_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODEFS_STAT, tablename), stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI NodeFs Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	// Insert Last Data
	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(TRUNCATE_DATA, lasttablename))
	if !errorCheck(del_err) {
		return
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastNodeFs Delete %d", processid, time.Now().UnixNano()))

	for _, stat_data := range stat_nodefs_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODEFS_STAT, lasttablename), stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastNodeFs Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func setPodFsRealtimeperf(current_data *sync.Map, prev_data *sync.Map, nodeuid string) (*PodFsPerfStat, *sync.Map) {
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
		current_podfs_summary_data = podFsSummaryPreprocessing(&current_podfs_data, &current_containerfs_data)
	} else {
		return &PodFsPerfStat{}, &sync.Map{}
	}

	if p, ok := prev_data.Load("podfs_insert"); ok {
		prev_podfs_data := p.(PodFsPerf)

		var prev_containerfs_data ContainerFsPerf
		if pc, ok2 := prev_data.Load("containerfs_insert"); ok2 {
			prev_containerfs_data = pc.(ContainerFsPerf)
		}
		prev_podfs_summary_data = podFsSummaryPreprocessing(&prev_podfs_data, &prev_containerfs_data)
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

		setSummaryData(summary_data, "fs", &current_podfs_summary_data, &prev_podfs_summary_data, i, pidx)

		insert_stat_data.ArrPoduid = append(insert_stat_data.ArrPoduid, current_podfs_summary_data.ArrPoduid[i])
		insert_stat_data.ArrNodeuid = append(insert_stat_data.ArrNodeuid, current_podfs_summary_data.ArrNodeuid[i])
		insert_stat_data.ArrNsUid = append(insert_stat_data.ArrNsUid, current_podfs_summary_data.ArrNsUid[i])
		insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_podfs_summary_data.ArrOntunetime[i])
		insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, 0)
		insert_stat_data.ArrDeviceid = append(insert_stat_data.ArrDeviceid, current_podfs_summary_data.ArrDeviceid[i])

		read_rate := getStatData("rate", current_podfs_summary_data.ArrFsreadsbytestotal[i], prev_podfs_summary_data.ArrFsreadsbytestotal[pidx], 0)
		write_rate := getStatData("rate", current_podfs_summary_data.ArrFswritesbytestotal[i], prev_podfs_summary_data.ArrFswritesbytestotal[pidx], 0)

		insert_stat_data.ArrFsreadrate = append(insert_stat_data.ArrFsreadrate, read_rate)
		insert_stat_data.ArrFswriterate = append(insert_stat_data.ArrFswriterate, write_rate)
		insert_stat_data.ArrFsiorate = append(insert_stat_data.ArrFsiorate, read_rate+write_rate)
	}

	lsdm, _ := last_stat_data_map.LoadOrStore(METRIC_VAR_PODFS, &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store(METRIC_VAR_PODFS, lsdm)

	return insert_stat_data, summary_data
}

func insertPodFsRealtimeperf(stat_podfs_map map[string]PodFsPerfStat, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_POD_FS_PERF)
	var lasttablename string = TB_KUBE_LAST_POD_FS_PERF

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, stat_data := range stat_podfs_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_PODFS_STAT, tablename), stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI PodFs Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	// Insert Last Data
	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(TRUNCATE_DATA, lasttablename))
	if !errorCheck(del_err) {
		return
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastPodFs Delete %d", processid, time.Now().UnixNano()))

	for _, stat_data := range stat_podfs_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_PODFS_STAT, lasttablename), stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastPodFs Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func setContainerFsRealtimeperf(current_data *sync.Map, prev_data *sync.Map, nodeuid string) (*ContainerFsPerfStat, *sync.Map) {
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

		setSummaryData(summary_data, "fs", &current_containerfs_data, &prev_containerfs_data, i, pidx)

		insert_stat_data.ArrContainername = append(insert_stat_data.ArrContainername, current_containerfs_data.ArrContainername[i])
		insert_stat_data.ArrPoduid = append(insert_stat_data.ArrPoduid, current_containerfs_data.ArrPoduid[i])
		insert_stat_data.ArrNodeuid = append(insert_stat_data.ArrNodeuid, current_containerfs_data.ArrNodeuid[i])
		insert_stat_data.ArrNsUid = append(insert_stat_data.ArrNsUid, current_containerfs_data.ArrNsUid[i])
		insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, current_containerfs_data.ArrOntunetime[i])
		insert_stat_data.ArrMetricid = append(insert_stat_data.ArrMetricid, current_containerfs_data.ArrMetricid[i])
		insert_stat_data.ArrDeviceid = append(insert_stat_data.ArrDeviceid, current_containerfs_data.ArrDeviceid[i])

		read_rate := getStatData("rate", current_containerfs_data.ArrFsreadsbytestotal[i], prev_containerfs_data.ArrFsreadsbytestotal[pidx], 0)
		write_rate := getStatData("rate", current_containerfs_data.ArrFswritesbytestotal[i], prev_containerfs_data.ArrFswritesbytestotal[pidx], 0)

		insert_stat_data.ArrFsreadrate = append(insert_stat_data.ArrFsreadrate, read_rate)
		insert_stat_data.ArrFswriterate = append(insert_stat_data.ArrFswriterate, write_rate)
		insert_stat_data.ArrFsiorate = append(insert_stat_data.ArrFsiorate, read_rate+write_rate)
	}

	lsdm, _ := last_stat_data_map.LoadOrStore(METRIC_VAR_CONTAINERFS, &sync.Map{})
	lsdm.(*sync.Map).Store(nodeuid, *insert_stat_data)
	last_stat_data_map.Store(METRIC_VAR_CONTAINERFS, lsdm)

	return insert_stat_data, summary_data
}

func insertContainerFsRealtimeperf(stat_containerfs_map map[string]ContainerFsPerfStat, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_CONTAINER_FS_PERF)
	var lasttablename string = TB_KUBE_LAST_CONTAINER_FS_PERF

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, stat_data := range stat_containerfs_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CONTAINERFS_STAT, tablename), stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI ContainerFs Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	// Insert Last Data
	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(TRUNCATE_DATA, lasttablename))
	if !errorCheck(del_err) {
		return
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastContainerFs Delete %d", processid, time.Now().UnixNano()))

	for _, stat_data := range stat_containerfs_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CONTAINERFS_STAT, lasttablename), stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - BRDI LastContainerFs Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func basicRealtimeDataSetting(metric_data *sync.Map, prev_data *sync.Map, nodeuid string, ontunetime int64, stat_map *BasicStatMap, processid string) {
	if _, ok := prev_data.Load("node_insert"); ok {
		common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - nodeuid %s BRDS NodeNet %d", processid, nodeuid, time.Now().UnixNano()))
		nndata, net_map := setNodeNetRealtimeperf(metric_data, prev_data, nodeuid)

		common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - nodeuid %s BRDS NodeFs %d", processid, nodeuid, time.Now().UnixNano()))
		nfdata, fs_map := setNodeFsRealtimeperf(metric_data, prev_data, nodeuid)

		common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - nodeuid %s BRDS Node %d", processid, nodeuid, time.Now().UnixNano()))
		ndata := setNodeRealtimeperf(metric_data, prev_data, net_map, fs_map, nodeuid)

		stat_map.Node[nodeuid] = *ndata
		stat_map.NodeNet[nodeuid] = *nndata
		stat_map.NodeFs[nodeuid] = *nfdata
	}
	if _, ok := prev_data.Load("pod_insert"); ok {
		common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - nodeuid %s BRDS PodNet %d", processid, nodeuid, time.Now().UnixNano()))
		pndata, net_map := setPodNetRealtimeperf(metric_data, prev_data, nodeuid)

		common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - nodeuid %s BRDS PodFs %d", processid, nodeuid, time.Now().UnixNano()))
		pfdata, fs_map := setPodFsRealtimeperf(metric_data, prev_data, nodeuid)

		common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - nodeuid %s BRDS Pod %d", processid, nodeuid, time.Now().UnixNano()))
		pdata, pod_map := setPodRealtimeperf(metric_data, prev_data, net_map, fs_map, nodeuid)

		stat_map.Pod[nodeuid] = *pdata
		stat_map.PodNet[nodeuid] = *pndata
		stat_map.PodFs[nodeuid] = *pfdata

		storeClusterMap(cluster_map, pod_map, nodeuid)
	}
	if _, ok := prev_data.Load("container_insert"); ok {
		common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - nodeuid %s BRDS ContainerFs %d", processid, nodeuid, time.Now().UnixNano()))
		cfdata, fs_map := setContainerFsRealtimeperf(metric_data, prev_data, nodeuid)

		common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - nodeuid %s BRDS Container %d", processid, nodeuid, time.Now().UnixNano()))
		cdata := setContainerRealtimeperf(metric_data, prev_data, fs_map, nodeuid)

		stat_map.Container[nodeuid] = *cdata
		stat_map.ContainerFs[nodeuid] = *cfdata
	}
}

func insertBasicRealtimeData(stat_map *BasicStatMap, ontunetime int64, processid string) {
	insertNodeRealtimeperf(stat_map.Node, ontunetime, processid)
	insertPodRealtimeperf(stat_map.Pod, ontunetime, processid)
	insertContainerRealtimeperf(stat_map.Container, ontunetime, processid)
	insertNodeNetRealtimeperf(stat_map.NodeNet, ontunetime, processid)
	insertPodNetRealtimeperf(stat_map.PodNet, ontunetime, processid)
	insertNodeFsRealtimeperf(stat_map.NodeFs, ontunetime, processid)
	insertPodFsRealtimeperf(stat_map.PodFs, ontunetime, processid)
	insertContainerFsRealtimeperf(stat_map.ContainerFs, ontunetime, processid)
}

func podSummaryPreprocessing(pod_data *PodPerf, container_data *ContainerPerf) PodSummaryPerf {
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

		var key string = container_data.ArrPoduid[i] + ":" + container_data.ArrNsUid[i]
		if val, ok := podindex[key]; !ok {
			podindex[key] = idx
			summary_data.ArrPoduid = append(summary_data.ArrPoduid, container_data.ArrPoduid[i])
			summary_data.ArrOntunetime = append(summary_data.ArrOntunetime, container_data.ArrOntunetime[i])
			summary_data.ArrTimestampMs = append(summary_data.ArrTimestampMs, container_data.ArrTimestampMs[i])
			summary_data.ArrNsUid = append(summary_data.ArrNsUid, container_data.ArrNsUid[i])
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

func podNetSummaryPreprocessing(podnet_data *PodNetPerf) PodNetPerf {
	var summary_data PodNetPerf = PodNetPerf{}
	var podnetindex map[string]int = make(map[string]int)
	var idx int
	for i := 0; i < len(podnet_data.ArrOntunetime); i++ {
		var key string = podnet_data.ArrPoduid[i] + ":" + strconv.Itoa(podnet_data.ArrInterfaceid[i]) + ":" + podnet_data.ArrNsUid[i]
		if val, ok := podnetindex[key]; !ok {
			podnetindex[key] = idx
			summary_data.ArrPoduid = append(summary_data.ArrPoduid, podnet_data.ArrPoduid[i])
			summary_data.ArrOntunetime = append(summary_data.ArrOntunetime, podnet_data.ArrOntunetime[i])
			summary_data.ArrTimestampMs = append(summary_data.ArrTimestampMs, podnet_data.ArrTimestampMs[i])
			summary_data.ArrNsUid = append(summary_data.ArrNsUid, podnet_data.ArrNsUid[i])
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

func podFsSummaryPreprocessing(podfs_data *PodFsPerf, containerfs_data *ContainerFsPerf) PodFsPerf {
	var summary_data PodFsPerf = PodFsPerf{}
	var podfsindex map[string]int = make(map[string]int)
	var idx int

	for i := 0; i < len(containerfs_data.ArrOntunetime); i++ {
		var key string = containerfs_data.ArrPoduid[i] + ":" + strconv.Itoa(containerfs_data.ArrDeviceid[i]) + ":" + containerfs_data.ArrNsUid[i]
		if val, ok := podfsindex[key]; !ok {
			podfsindex[key] = idx
			summary_data.ArrPoduid = append(summary_data.ArrPoduid, containerfs_data.ArrPoduid[i])
			summary_data.ArrOntunetime = append(summary_data.ArrOntunetime, containerfs_data.ArrOntunetime[i])
			summary_data.ArrTimestampMs = append(summary_data.ArrTimestampMs, containerfs_data.ArrTimestampMs[i])
			summary_data.ArrNsUid = append(summary_data.ArrNsUid, containerfs_data.ArrNsUid[i])
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

func setPodSummary(summary_data *sync.Map, key SummaryIndex, datatype string, metricname string, value float64) {
	key.DataType = datatype

	v, _ := summary_data.LoadOrStore(key, &sync.Map{})
	value_map := v.(*sync.Map)
	value_map.Store(metricname, value)
	summary_data.Store(key, value_map)
}

func setSummaryIndex(src SummaryPerfInterface, idx int) SummaryIndex {
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

func setSummaryData(summary_data *sync.Map, datatype string, current SummaryPerfInterface, previous SummaryPerfInterface, cidx int, pidx int) {
	summary_key := setSummaryIndex(current, cidx)
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

func getSummaryData(summary_data *sync.Map, key SummaryIndex, datatype string, metricname string) float64 {
	key.DataType = datatype
	key.MetricName = metricname
	if val, ok := summary_data.Load(key); ok {
		return val.(float64)
	}

	return float64(0)
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
