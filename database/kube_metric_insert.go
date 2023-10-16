package database

import (
	"fmt"
	"onTuneKubeManager/common"
	"sync"
	"time"
)

// Metric Insert
func MetricSender() {
	data_queue_map := &sync.Map{}

	go MetricInsertMap(data_queue_map)

	for {
		metric_data := <-Channelmetric_insert

		var nodeuid string
		var nodename string = metric_data[METRIC_VAR_NODE].(string)
		var hostname string = metric_data[METRIC_VAR_HOST].(string)
		var ontunetime int64 = metric_data["ontunetime"].(int64)

		common.LogManager.Debug(fmt.Sprintf("Manager DB - Metric Insert Data Store Queue Map start: %s %s %d", hostname, nodename, ontunetime))
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
		metric_data_map.Store("hostname", hostname)
		metric_data_map.Store("ontunetime", ontunetime)
		metric_data_map.Store("statflag", false)
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
		QueueMapProcessing(nodeuid, &metric_arr)
		data_queue_map.Store(nodeuid, metric_arr)

		common.LogManager.Debug(fmt.Sprintf("Manager DB - Metric Insert Data Store Queue Map end: %s %s %d", hostname, nodename, ontunetime))
		time.Sleep(time.Millisecond * 10)
	}
}

func QueueMapProcessing(nodeuid string, metric_arr *[]*sync.Map) {
	if cot, cok := (*metric_arr)[len(*metric_arr)-1].Load("ontunetime"); cok {
		for i := len(*metric_arr) - 1; i >= 0; i-- {
			if pot, pok := (*metric_arr)[i].Load("ontunetime"); pok {
				if pot.(int64) <= cot.(int64)-int64(common.RateInterval) {
					(*metric_arr)[i].Store("statflag", true)
					*metric_arr = (*metric_arr)[i:]

					break
				}
			}
		}
	}
}

func MetricInsertMap(data_queue_map *sync.Map) {
	var last_updatetime int64
	_, bias := GetOntuneTime()

	for {
		ontunetime := time.Now().Unix() - bias
		if last_updatetime+int64(common.RealtimeInterval) <= ontunetime {
			// Select Query는 최초 1회, 그리고 Interval 내 조건을 충족할 때에만 수행하는 것으로 함
			_, bias = GetOntuneTime()
			last_updatetime = ontunetime
			processid := getRandomProcessId()

			common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - Metric Queue Insertion start: %d", processid, ontunetime))

			basic_stat_map := &BasicStatMap{
				Node:        make(map[string]NodePerfStat),
				Pod:         make(map[string]PodPerfStat),
				Container:   make(map[string]ContainerPerfStat),
				NodeNet:     make(map[string]NodeNetPerfStat),
				PodNet:      make(map[string]PodNetPerfStat),
				NodeFs:      make(map[string]NodeFsPerfStat),
				PodFs:       make(map[string]PodFsPerfStat),
				ContainerFs: make(map[string]ContainerFsPerfStat),
			}

			basic_raw_map := &BasicRawMap{
				Node:        make(map[string]NodePerf),
				Pod:         make(map[string]PodPerf),
				Container:   make(map[string]ContainerPerf),
				NodeNet:     make(map[string]NodeNetPerf),
				PodNet:      make(map[string]PodNetPerf),
				NodeFs:      make(map[string]NodeFsPerf),
				PodFs:       make(map[string]PodFsPerf),
				ContainerFs: make(map[string]ContainerFsPerf),
			}

			// 초기화
			initMapUidInfo(data_queue_map)

			data_queue_map.Range(func(key, value any) bool {
				nodeuid := key.(string)
				metric_arr := value.([]*sync.Map)

				if len(metric_arr) == 0 {
					return true
				}

				metric_data := metric_arr[len(metric_arr)-1]

				var cluster_status bool = false
				var hostname string
				if h, ok := metric_data.Load("hostname"); ok {
					hostname = h.(string)

					if cs, ok := common.ClusterStatusMap.Load(hostname); ok {
						cluster_status = cs.(bool)
					}
				}

				if cluster_status {
					// Get Recent Queue Value
					common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - nodeuid %s Raw Data setting start %d", processid, nodeuid, time.Now().UnixNano()))
					rawDataSetting(metric_data, basic_raw_map, nodeuid, ontunetime, processid)
					common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - nodeuid %s Raw Data setting end %d", processid, nodeuid, time.Now().UnixNano()))

					// Get Previous Queue Value
					prev_data := metric_arr[0]
					var stat_flag bool
					if sf, ok := prev_data.Load("statflag"); ok {
						stat_flag = sf.(bool)
						if stat_flag {
							common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - nodeuid %s Basic Realtime Data setting start %d", processid, nodeuid, time.Now().UnixNano()))
							basicRealtimeDataSetting(metric_data, prev_data, nodeuid, ontunetime, basic_stat_map, processid)
							common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - nodeuid %s Basic Realtime Data setting end %d", processid, nodeuid, time.Now().UnixNano()))
						}
					}
				} else {
					clusterid := common.ClusterID[hostname]
					data_queue_map.Store(key, []*sync.Map{})
					cluster_map.Delete(clusterid)
				}

				return true
			})

			if len(basic_raw_map.Node) > 0 {
				common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - Raw Data insertion start", processid))
				insertRawRealtimeData(basic_raw_map, ontunetime, processid)
				common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - Raw Data insertion end", processid))
			}

			if len(basic_stat_map.Node) > 0 {
				common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - Perf Data insertion start", processid))
				insertBasicRealtimeData(basic_stat_map, ontunetime, processid)
				common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - Perf Data insertion end", processid))
			}

			// Summary Stat Insert
			if checkMapNil(cluster_map) {
				common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - Stat Data insertion start", processid))

				insertClusterRealtimeperf(cluster_map, ontunetime)
				common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - Cluster insertion %d", processid, time.Now().UnixNano()))

				insertNamespaceRealtimeperf(cluster_map, ontunetime)
				common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - Namespace insertion %d", processid, time.Now().UnixNano()))

				insertWorkkloadRealtimeperf(cluster_map, "ReplicaSet", ontunetime)
				common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - Replicaset insertion %d", processid, time.Now().UnixNano()))

				insertWorkkloadRealtimeperf(cluster_map, "DaemonSet", ontunetime)
				common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - Daemonset insertion %d", processid, time.Now().UnixNano()))

				insertWorkkloadRealtimeperf(cluster_map, "StatefulSet", ontunetime)
				common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - Statefulset insertion %d", processid, time.Now().UnixNano()))

				insertWorkkloadRealtimeperf(cluster_map, "Deployment", ontunetime)
				common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - Deployment insertion %d", processid, time.Now().UnixNano()))

				common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - Stat Data insertion end", processid))
			}

			checkMapUidInfo(data_queue_map)
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func AvgMetricDataInsert(ontunetime int64, bias int64) {
	nobiastime := ontunetime + bias
	starttime := ontunetime - int64(common.AvgInterval)
	endtime := ontunetime

	insertNodeAvgPerf(nobiastime, starttime, endtime)
	insertPodAvgPerf(nobiastime, starttime, endtime)
	insertContainerAvgPerf(nobiastime, starttime, endtime)
	insertNodeNetAvgPerf(nobiastime, starttime, endtime)
	insertPodNetAvgPerf(nobiastime, starttime, endtime)
	insertNodeFsAvgPerf(nobiastime, starttime, endtime)
	insertPodFsAvgPerf(nobiastime, starttime, endtime)
	insertContainerFsAvgPerf(nobiastime, starttime, endtime)

	insertClusterAvgPerf(nobiastime, starttime, endtime)
	insertNamespaceAvgPerf(nobiastime, starttime, endtime)
	insertWorkloadAvgPerf(nobiastime, starttime, endtime, METRIC_VAR_DEPLOYMENT)
	insertWorkloadAvgPerf(nobiastime, starttime, endtime, METRIC_VAR_DAEMONSET)
	insertWorkloadAvgPerf(nobiastime, starttime, endtime, METRIC_VAR_STATEFULSET)
	insertWorkloadAvgPerf(nobiastime, starttime, endtime, METRIC_VAR_REPLICASET)

	common.LogManager.WriteLog("Manager DB - Avg Stat Data insertion end")
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
			hostMapUidInfo.prevNodeUid = copyUIDMap(mapUidInfo[hostname].currentNodeUid)
			hostMapUidInfo.prevPodUid = copyUIDMap(mapUidInfo[hostname].currentPodUid)
			hostMapUidInfo.prevNamespaceUid = copyUIDMap(mapUidInfo[hostname].currentNamespaceUid)
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
		requestCheck(mapUidInfo[hostname].prevNodeUid, mapUidInfo[hostname].currentNodeUid, &requestflag)
		common.MutexNode.Unlock()

		// UIDMap Compare - Pod
		common.MutexPod.Lock()
		requestCheck(mapUidInfo[hostname].prevPodUid, mapUidInfo[hostname].currentPodUid, &requestflag)
		common.MutexPod.Unlock()

		// UIDMap Compare - Namespace
		common.MutexNs.Lock()
		requestCheck(mapUidInfo[hostname].prevNamespaceUid, mapUidInfo[hostname].currentNamespaceUid, &requestflag)
		common.MutexNs.Unlock()

		if requestflag {
			common.ChannelRequestChangeHost <- hostname
		}
	}
}

func requestCheck(srcmap map[string]struct{}, dstmap map[string]struct{}, flag *bool) {
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

func copyUIDMap(uid_map map[string]struct{}) map[string]struct{} {
	uid_copy := map[string]struct{}{}
	for k, v := range uid_map {
		uid_copy[k] = v
	}

	return uid_copy
}
