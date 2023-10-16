package database

import (
	"encoding/json"
	"fmt"
	"onTuneKubeManager/common"
)

func RealtimeMetricDataParsing(metric_data []byte) {
	var data interface{}

	err := json.Unmarshal([]byte(metric_data), &data)
	if !errorCheck(err) {
		return
	}

	map_datas := data.(map[string]interface{})
	insert_metric_data := make(map[string]interface{})
	insert_metric_data[METRIC_VAR_NODE] = map_datas[METRIC_VAR_NODE]
	insert_metric_data[METRIC_VAR_HOST] = map_datas[METRIC_VAR_HOST]

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	insert_metric_data["ontunetime"] = ontunetime

	common.LogManager.Debug(fmt.Sprintf("Manager DB - Metric Data Parsing start: %s %s %d", insert_metric_data[METRIC_VAR_HOST], insert_metric_data[METRIC_VAR_NODE], insert_metric_data["ontunetime"]))

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
		common.ChannelRequestChangeHost <- insert_metric_data[METRIC_VAR_HOST].(string)
	}

	common.LogManager.Debug(fmt.Sprintf("Manager DB - Metric Data Parsing end: %s %s %d", insert_metric_data[METRIC_VAR_HOST], insert_metric_data[METRIC_VAR_NODE], insert_metric_data["ontunetime"]))
	Channelmetric_insert <- insert_metric_data
}

func setNodeMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var nodeInsert NodePerf
	var flag bool = false

	for _, node_data := range map_data.([]interface{}) {
		node_data_map := node_data.(map[string]interface{})
		nodeuid := getUID(METRIC_VAR_NODE, (*map_metric)[METRIC_VAR_HOST], node_data_map[METRIC_VAR_NODE])
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

func setPodMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var arrInsert PodPerf
	var flag bool = false

	for _, pod_data := range map_data.([]interface{}) {
		pod_data_map := pod_data.(map[string]interface{})

		nsuid := getUID(METRIC_VAR_NAMESPACE, (*map_metric)[METRIC_VAR_HOST], pod_data_map[METRIC_VAR_NAMESPACE])
		poduid := getUID(METRIC_VAR_POD, pod_data_map[METRIC_VAR_NAMESPACE].(string), pod_data_map[METRIC_VAR_POD].(string))
		nodeuid := getUID(METRIC_VAR_NODE, (*map_metric)[METRIC_VAR_HOST], pod_data_map[METRIC_VAR_NODE])
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
				arrInsert.ArrNsUid = append(arrInsert.ArrNsUid, nsuid)
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

func setContainerMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var arrInsert ContainerPerf
	var flag bool = false

	for _, container_data := range map_data.([]interface{}) {
		container_data_map := container_data.(map[string]interface{})

		nsuid := getUID(METRIC_VAR_NAMESPACE, (*map_metric)[METRIC_VAR_HOST], container_data_map[METRIC_VAR_NAMESPACE])
		poduid := getUID(METRIC_VAR_POD, container_data_map[METRIC_VAR_NAMESPACE].(string), container_data_map[METRIC_VAR_POD].(string))
		nodeuid := getUID(METRIC_VAR_NODE, (*map_metric)[METRIC_VAR_HOST], container_data_map[METRIC_VAR_NODE])
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
				arrInsert.ArrNsUid = append(arrInsert.ArrNsUid, nsuid)
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

func setNodeNetMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var arrInsert NodeNetPerf
	var flag bool = false

	for _, nodenet_data := range map_data.([]interface{}) {
		nodenet_data_map := nodenet_data.(map[string]interface{})
		nodeuid := getUID(METRIC_VAR_NODE, (*map_metric)[METRIC_VAR_HOST], nodenet_data_map[METRIC_VAR_NODE])
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

func setPodNetMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var arrInsert PodNetPerf
	var flag bool = false

	for _, podnet_data := range map_data.([]interface{}) {
		podnet_data_map := podnet_data.(map[string]interface{})

		nsuid := getUID(METRIC_VAR_NAMESPACE, (*map_metric)[METRIC_VAR_HOST], podnet_data_map[METRIC_VAR_NAMESPACE])
		poduid := getUID(METRIC_VAR_POD, podnet_data_map[METRIC_VAR_NAMESPACE].(string), podnet_data_map[METRIC_VAR_POD].(string))
		nodeuid := getUID(METRIC_VAR_NODE, (*map_metric)[METRIC_VAR_HOST], podnet_data_map[METRIC_VAR_NODE])
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
				arrInsert.ArrNsUid = append(arrInsert.ArrNsUid, nsuid)
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

func setNodeFsMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var arrInsert NodeFsPerf
	var flag bool = false

	for _, nodefs_data := range map_data.([]interface{}) {
		nodefs_data_map := nodefs_data.(map[string]interface{})
		nodeuid := getUID(METRIC_VAR_NODE, (*map_metric)[METRIC_VAR_HOST], nodefs_data_map[METRIC_VAR_NODE])
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

func setPodFsMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var arrInsert PodFsPerf
	var flag bool = false

	for _, podfs_data := range map_data.([]interface{}) {
		podfs_data_map := podfs_data.(map[string]interface{})

		nsuid := getUID(METRIC_VAR_NAMESPACE, (*map_metric)[METRIC_VAR_HOST], podfs_data_map[METRIC_VAR_NAMESPACE])
		poduid := getUID(METRIC_VAR_POD, podfs_data_map[METRIC_VAR_NAMESPACE].(string), podfs_data_map[METRIC_VAR_POD].(string))
		nodeuid := getUID(METRIC_VAR_NODE, (*map_metric)[METRIC_VAR_HOST], podfs_data_map[METRIC_VAR_NODE])
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
				arrInsert.ArrNsUid = append(arrInsert.ArrNsUid, nsuid)
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

func setContainerFsMetric(map_data interface{}, map_metric *map[string]interface{}) bool {
	var arrInsert ContainerFsPerf
	var flag bool = false

	for _, containerfs_data := range map_data.([]interface{}) {
		containerfs_data_map := containerfs_data.(map[string]interface{})

		nsuid := getUID(METRIC_VAR_NAMESPACE, (*map_metric)[METRIC_VAR_HOST], containerfs_data_map[METRIC_VAR_NAMESPACE])
		poduid := getUID(METRIC_VAR_POD, containerfs_data_map[METRIC_VAR_NAMESPACE].(string), containerfs_data_map[METRIC_VAR_POD].(string))
		nodeuid := getUID(METRIC_VAR_NODE, (*map_metric)[METRIC_VAR_HOST], containerfs_data_map[METRIC_VAR_NODE])
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
				arrInsert.ArrNsUid = append(arrInsert.ArrNsUid, nsuid)
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
