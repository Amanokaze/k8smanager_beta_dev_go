package database

import (
	"context"
	"fmt"
	"math"
	"onTuneKubeManager/common"
	"strconv"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func loadAndInitMap(cm *sync.Map, pm *sync.Map, key string) (float64, float64) {
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

func insertClusterRealtimeperf(cluster_map *sync.Map, ontunetime int64) {
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

			calculateNodeSummary(cal_map, nodeuid)

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
		if conn == nil || err != nil {
			errorDisconnect(errors.Wrap(err, "Acquire connection error"))
			return
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
			return
		}

		tablename := getTableName(TB_KUBE_CLUSTER_PERF)
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CLUSTER_STAT, tablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
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

		_, err = tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA_CONDITION, lasttablename, "clusterid", join_clusterids, ontunetime))
		if !errorCheck(err) {
			return
		}

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CLUSTER_STAT, lasttablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}

		err = tx.Commit(context.Background())
		if !errorCheck(errors.Wrap(err, "Commit error")) {
			return
		}

		conn.Release()

		updateTableinfo(tablename, ontunetime)
		updateTableinfo(lasttablename, ontunetime)
	}
}

func insertNamespaceRealtimeperf(cluster_map *sync.Map, ontunetime int64) {
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
					calculateNodeSummary(ns_map, nodeuid)
				}

				return true
			})

			return true
		})

		cal_ns_map.Range(func(nsuid, nsvalue any) bool {
			ns_map := nsvalue.(*sync.Map)

			insert_stat_data.ArrOntunetime = append(insert_stat_data.ArrOntunetime, ontunetime)
			insert_stat_data.ArrNsUid = append(insert_stat_data.ArrNsUid, nsuid.(string))

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

	if len(insert_stat_data.ArrNsUid) > 0 {
		conn, err := common.DBConnectionPool.Acquire(context.Background())
		if conn == nil || err != nil {
			errorDisconnect(errors.Wrap(err, "Acquire connection error"))
			return
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
			return
		}

		tablename := getTableName(TB_KUBE_NAMESPACE_PERF)
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NS_STAT, tablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}

		lasttablename := TB_KUBE_LAST_NAMESPACE_PERF

		ns_uids_map := make(map[string]struct{})
		ns_uids := make([]string, 0)
		for _, nsuid := range insert_stat_data.ArrNsUid {
			if _, ok := ns_uids_map[nsuid]; !ok {
				ns_uids_map[nsuid] = struct{}{}
				uid_mark := fmt.Sprintf("'%s'", nsuid)
				ns_uids = append(ns_uids, uid_mark)
			}
		}
		join_nsuids := strings.Join(ns_uids, ",")

		_, err = tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA_CONDITION, lasttablename, "nsuid", join_nsuids, ontunetime))
		if !errorCheck(err) {
			return
		}

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NS_STAT, lasttablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}

		err = tx.Commit(context.Background())
		if !errorCheck(errors.Wrap(err, "Commit error")) {
			return
		}

		conn.Release()

		updateTableinfo(tablename, ontunetime)
		updateTableinfo(lasttablename, ontunetime)
	}
}

func insertWorkkloadRealtimeperf(cluster_map *sync.Map, kind string, ontunetime int64) {
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
		cal_wl_map := processWorkloadSummaryByNode(nodevalue.(*sync.Map), kind)

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
		if conn == nil || err != nil {
			errorDisconnect(errors.Wrap(err, "Acquire connection error"))
			return
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
			return
		}

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_WORKLOAD_STAT, tablename, uidname), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
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

		_, err = tx.Exec(context.Background(), fmt.Sprintf(DELETE_DATA_CONDITION, lasttablename, uidname, join_workload_uids, ontunetime))
		if !errorCheck(err) {
			return
		}

		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_WORKLOAD_STAT, lasttablename, uidname), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}

		err = tx.Commit(context.Background())
		if !errorCheck(errors.Wrap(err, "Commit error")) {
			return
		}

		conn.Release()

		updateTableinfo(tablename, ontunetime)
		updateTableinfo(lasttablename, ontunetime)
	}
}

func calculateNodeSummary(summary_map *sync.Map, nodeuid string) {
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

	cprocess, _ := loadAndInitMap(current_map, previous_map, "process")
	ccpuusage, pcpuusage := loadAndInitMap(current_map, previous_map, "cpuusage")
	ccpusystem, pcpusystem := loadAndInitMap(current_map, previous_map, "cpusystem")
	ccpuuser, pcpuuser := loadAndInitMap(current_map, previous_map, "cpuuser")
	cmemoryws, _ := loadAndInitMap(current_map, previous_map, "memoryws")
	cmemoryswap, _ := loadAndInitMap(current_map, previous_map, "memoryswap")

	cal_proc, _ := calculation_map.LoadOrStore("processcount", float64(0))
	calculation_map.Store("processcount", cal_proc.(float64)+cprocess)

	cal_cpuusage, _ := calculation_map.LoadOrStore("cpuusage", float64(0))
	calculation_map.Store("cpuusage", cal_cpuusage.(float64)+getStatData("usage", ccpuusage, pcpuusage, float64(cpucount)))

	cal_cpusystem, _ := calculation_map.LoadOrStore("cpusystem", float64(0))
	calculation_map.Store("cpusystem", cal_cpusystem.(float64)+getStatData("usage", ccpusystem, pcpusystem, float64(cpucount)))

	cal_cpuuser, _ := calculation_map.LoadOrStore("cpuuser", float64(0))
	calculation_map.Store("cpuuser", cal_cpuuser.(float64)+getStatData("usage", ccpuuser, pcpuuser, float64(cpucount)))

	cal_cpuusagecores, _ := calculation_map.LoadOrStore("cpuusagecores", float64(0))
	calculation_map.Store("cpuusagecores", cal_cpuusagecores.(float64)+getStatData("core", ccpuusage, pcpuusage, 0))

	cal_cputotalcores, _ := calculation_map.LoadOrStore("cputotalcores", float64(0))
	calculation_map.Store("cputotalcores", cal_cputotalcores.(float64)+math.Round(float64(cpucount)))

	cal_memoryusage, _ := calculation_map.LoadOrStore("memoryusage", float64(0))
	calculation_map.Store("memoryusage", cal_memoryusage.(float64)+getStatData("current_usage", cmemoryws, 0, float64(memorysize)))

	cal_memoryusagebytes, _ := calculation_map.LoadOrStore("memoryusagebytes", float64(0))
	calculation_map.Store("memoryusagebytes", cal_memoryusagebytes.(float64)+math.Round(cmemoryws))

	cal_memorysizebytes, _ := calculation_map.LoadOrStore("memorysizebytes", float64(0))
	calculation_map.Store("memorysizebytes", cal_memorysizebytes.(float64)+math.Round(float64(memorysize)))

	cal_memoryswap, _ := calculation_map.LoadOrStore("memoryswap", float64(0))
	calculation_map.Store("memoryswap", cal_memoryswap.(float64)+getStatData("current_usage", cmemoryswap, 0, float64(memorysize)))

	cnetrbt, pnetrbt := loadAndInitMap(current_map, previous_map, "netrbt")
	cnetrbe, pnetrbe := loadAndInitMap(current_map, previous_map, "netrbe")
	cnettbt, pnettbt := loadAndInitMap(current_map, previous_map, "nettbt")
	cnettbe, pnettbe := loadAndInitMap(current_map, previous_map, "nettbe")

	rcv_rate := getStatData("rate", cnetrbt, pnetrbt, 0)
	rcv_errors := getStatData("subtract", cnetrbe, pnetrbe, 0)
	trans_rate := getStatData("rate", cnettbt, pnettbt, 0)
	trans_errors := getStatData("rate", cnettbe, pnettbe, 0)

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

	cfsrbt, pfsrbt := loadAndInitMap(current_map, previous_map, "fsrbt")
	cfswbt, pfswbt := loadAndInitMap(current_map, previous_map, "fswbt")

	fs_read_rate := getStatData("rate", cfsrbt, pfsrbt, 0)
	fs_write_rate := getStatData("rate", cfswbt, pfswbt, 0)

	cal_fsreadrate, _ := calculation_map.LoadOrStore("fsreadrate", float64(0))
	calculation_map.Store("fsreadrate", cal_fsreadrate.(float64)+fs_read_rate)

	cal_fswriterate, _ := calculation_map.LoadOrStore("fswriterate", float64(0))
	calculation_map.Store("fswriterate", cal_fswriterate.(float64)+fs_write_rate)

	cal_fsiorate, _ := calculation_map.LoadOrStore("fsiorate", float64(0))
	calculation_map.Store("fsiorate", cal_fsiorate.(float64)+fs_read_rate+fs_write_rate)

	summary_map.Delete(SUMMARY_DATATYPE_CURR)
	summary_map.Delete(SUMMARY_DATATYPE_PREV)
}

func processWorkloadSummaryByNode(nodevalue *sync.Map, kind string) *sync.Map {
	var cal_map *sync.Map = &sync.Map{}
	nodevalue.Range(func(nuid, pmap any) bool {
		nodeuid := nuid.(string)
		pod_map := pmap.(*sync.Map)

		pod_map.Range(func(sumkey, pod_value any) bool {
			summary_index := sumkey.(SummaryIndex)

			wluid, flag := getPodReference(summary_index.PodUID, kind)
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
				calculateNodeSummary(wl_map, nodeuid)
			}

			return true
		})

		return true
	})

	return cal_map
}

func getPodReference(poduid string, kind string) (string, bool) {
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

func storeClusterMap(cluster_map *sync.Map, pod_map *sync.Map, nodeuid string) {
	node_cluster_map, _ := common.ResourceMap.Load("node_cluster")
	cid, _ := node_cluster_map.(*sync.Map).LoadOrStore(nodeuid, 0)
	clusterid := cid.(int)

	if checkMapNil(pod_map) {
		cm, _ := cluster_map.LoadOrStore(clusterid, &sync.Map{})
		cm_node_map := cm.(*sync.Map)
		cm_node_map.Store(nodeuid, pod_map)
		cluster_map.Store(clusterid, cm_node_map)
	}
}
