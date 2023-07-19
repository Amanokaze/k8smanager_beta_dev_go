package database

import (
	"encoding/json"
	"fmt"
	"onTuneKubeManager/common"
	"onTuneKubeManager/kubeapi"
	"sync"
	"time"
)

func ResourceReceive() {
	for {
		info_msg := <-common.ChannelResourceData
		_, biastime = GetOntuneTime()
		if val, ok := info_msg["resourceinfo"]; ok {
			writeLog("this data is resource info!!")
			var data map[string]interface{}
			setResourceinfo(val, data)
			writeLog("resource info End")
		} else if val, ok := info_msg["namespaceinfo"]; ok {
			writeLog("this data is namespace info!!")
			setNamespaceinfo(val)
			writeLog("namespace info End")
		} else if val, ok := info_msg["nodeinfo"]; ok {
			writeLog("this data is node info!!")
			setNodeinfo(val)
			writeLog("node info End")
		} else if val, ok := info_msg["podsinfo"]; ok {
			writeLog("this data is pods info!!")
			PodContainerinfo := setPodsinfo(val)
			setContainerinfo(PodContainerinfo)
			writeLog("pods info End")
		} else if val, ok := info_msg["serviceinfo"]; ok {
			writeLog("this data is service info!!")
			setServiceinfo(val)
			writeLog("service info End")
		} else if val, ok := info_msg["pvcinfo"]; ok {
			writeLog("this data is pvc info!!")
			setPvcinfo(val)
			writeLog("pvc info End")
		} else if val, ok := info_msg["pvinfo"]; ok {
			writeLog("this data is pv info!!")
			setPvinfo(val)
			writeLog("pv info End")
		} else if val, ok := info_msg["deployinfo"]; ok {
			writeLog("this data is deploy info!!")
			setDeployinfo(val)
			writeLog("deploy info End")
		} else if val, ok := info_msg["statefulset"]; ok {
			writeLog("this data is statefulset info!!")
			setStatefulSetinfo(val)
			writeLog("statefulset info End")
		} else if val, ok := info_msg["daemonset"]; ok {
			writeLog("this data is daemonset info!!")
			setDaemonSetinfo(val)
			writeLog("daemonset info End")
		} else if val, ok := info_msg["replicaset"]; ok {
			writeLog("this data is replicaset info!!")
			setReplicaSetinfo(val)
			writeLog("replicaset info End")
		} else if val, ok := info_msg["sc"]; ok {
			writeLog("this data is sc info!!")
			setStorageclassinfo(val)
			writeLog("sc info End")
		} else if val, ok := info_msg["ingress_info"]; ok {
			writeLog("this data is ingress info!!")
			ingressHostinfo := setIngressinfo(val)
			setIngressHostinfo(ingressHostinfo)
			writeLog("ingress info End")
		} else {
			writeLog("no data??")
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func setResourceinfo(info_msg []byte, data map[string]interface{}) {
	var Arr_rcs Resourceinfo

	err := json.Unmarshal([]byte(info_msg), &data)
	errorCheck(err)

	for k2 := range data {
		rcs_map := data[k2].(map[string]interface{})
		Kind := rcs_map["Kind"].(string)
		Group := rcs_map["Group"].(string)
		Version := rcs_map["Version"].(string)
		Endpoint := rcs_map["Endpoint"].(string)
		Enabled := rcs_map["Enabled"].(bool)
		HostIP := rcs_map["Host"].(string)

		Arr_rcs.ArrClusterid = append(Arr_rcs.ArrClusterid, common.ClusterID[HostIP])
		Arr_rcs.ArrResourcename = append(Arr_rcs.ArrResourcename, Kind)
		Arr_rcs.ArrApiclass = append(Arr_rcs.ArrApiclass, Group)
		Arr_rcs.ArrVersion = append(Arr_rcs.ArrVersion, Version)
		Arr_rcs.ArrEndpoint = append(Arr_rcs.ArrEndpoint, Endpoint)
		if Enabled {
			Arr_rcs.ArrEnabled = append(Arr_rcs.ArrEnabled, 1)
		} else {
			Arr_rcs.ArrEnabled = append(Arr_rcs.ArrEnabled, 0)
		}
		Arr_rcs.ArrCreateTime = append(Arr_rcs.ArrCreateTime, 0)
		Arr_rcs.ArrUpdateTime = append(Arr_rcs.ArrUpdateTime, 0)
	}

	var i_data interface{}
	map_info := make(map[string]interface{})
	i_data = Arr_rcs
	map_info["resourceinfo_insert"] = i_data
	ChannelResourceInsert <- map_info
}

func setNamespaceinfo(info_msg []byte) {
	var tempNamespaceInfo map[string]kubeapi.MappingNamespace = make(map[string]kubeapi.MappingNamespace)
	var map_data []kubeapi.MappingNamespace
	var ArrResource Namespaceinfo
	var host string

	err := json.Unmarshal([]byte(info_msg), &map_data)
	errorCheck(err)

	for _, resource_data := range map_data {
		host = resource_data.Host

		ArrResource.ArrNsuid = append(ArrResource.ArrNsuid, resource_data.UID)
		ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[resource_data.Host])
		ArrResource.ArrNsname = append(ArrResource.ArrNsname, resource_data.Name)
		ArrResource.ArrStarttime = append(ArrResource.ArrStarttime, getStarttime(resource_data.StartTime.Unix(), biastime))
		ArrResource.ArrLabels = append(ArrResource.ArrLabels, resource_data.Labels)
		ArrResource.ArrStatus = append(ArrResource.ArrStatus, resource_data.Status)
		ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1)
		ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
		ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)

		resource_data.StartTime = time.Unix(getStarttime(resource_data.StartTime.Unix(), biastime), 0)
		tempNamespaceInfo[resource_data.UID] = resource_data // 신규로 들어온 데이터의 맵

		ns_map, _ := common.ResourceMap.Load("namespace")
		ns_label_map, _ := common.ResourceMap.Load("namespace_label")

		ns_key := fmt.Sprintf("%s/%s", resource_data.Host, resource_data.Name)
		ns_map.(*sync.Map).Store(ns_key, resource_data.UID)
		ns_label_map.(*sync.Map).Store(ns_key, resource_data.Labels)
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		mapNamespaceInfo := apiresource.namespace

		if len(mapNamespaceInfo) > 0 {
			i_data = tempNamespaceInfo
			map_info["namespace_update"] = i_data
			common.LogManager.Debug(fmt.Sprintf("map_data : %v", map_data))
			ChannelResourceInsert <- map_info
		} else {
			i_data = ArrResource
			map_info["namespace_insert"] = i_data
			ChannelResourceInsert <- map_info
			if len(mapNamespaceInfo) == 0 {
				for _, resource_data := range map_data {
					mapNamespaceInfo[resource_data.UID] = resource_data
				}
			}
		}
	}
}

func setNodeinfo(info_msg []byte) {
	var tempNodeInfo map[string]kubeapi.MappingNode = make(map[string]kubeapi.MappingNode)
	var map_data []kubeapi.MappingNode
	var ArrResource Nodeinfo
	var host string

	err := json.Unmarshal([]byte(info_msg), &map_data)
	errorCheck(err)

	for _, resource_data := range map_data {
		host = resource_data.Host

		ArrResource.ArrManagerid = append(ArrResource.ArrManagerid, common.ManagerID)
		ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[resource_data.Host])
		ArrResource.ArrNodeUid = append(ArrResource.ArrNodeUid, resource_data.UID)
		ArrResource.ArrNodename = append(ArrResource.ArrNodename, resource_data.Name)
		ArrResource.ArrNodenameext = append(ArrResource.ArrNodenameext, resource_data.Name)
		ArrResource.ArrNodetype = append(ArrResource.ArrNodetype, resource_data.NodeType)
		ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1)
		ArrResource.ArrStarttime = append(ArrResource.ArrStarttime, getStarttime(resource_data.StartTime.Unix(), biastime))
		ArrResource.ArrLabels = append(ArrResource.ArrLabels, resource_data.Labels)
		ArrResource.ArrKernelversion = append(ArrResource.ArrKernelversion, resource_data.KernelVersion)
		ArrResource.Arr_OSimage = append(ArrResource.Arr_OSimage, resource_data.OSImage)
		ArrResource.Arr_OSname = append(ArrResource.Arr_OSname, resource_data.OSName)
		ArrResource.ArrContainerruntimever = append(ArrResource.ArrContainerruntimever, resource_data.ContainerRuntimeVersion)
		ArrResource.ArrKubeletver = append(ArrResource.ArrKubeletver, resource_data.KubeletVersion)
		ArrResource.ArrKubeproxyver = append(ArrResource.ArrKubeproxyver, resource_data.KubeProxyVersion)
		ArrResource.ArrCpuarch = append(ArrResource.ArrCpuarch, resource_data.CPUArch)
		ArrResource.ArrCpucount = append(ArrResource.ArrCpucount, int(resource_data.CPUCount))
		ArrResource.ArrEphemeralstorage = append(ArrResource.ArrEphemeralstorage, resource_data.EphemeralStorage)
		ArrResource.ArrMemorysize = append(ArrResource.ArrMemorysize, resource_data.MemorySize)
		ArrResource.ArrPods = append(ArrResource.ArrPods, resource_data.Pods)
		ArrResource.Arr_IP = append(ArrResource.Arr_IP, resource_data.IP)
		ArrResource.ArrStatus = append(ArrResource.ArrStatus, resource_data.Status)
		ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
		ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)

		resource_data.StartTime = time.Unix(getStarttime(resource_data.StartTime.Unix(), biastime), 0)
		tempNodeInfo[resource_data.UID] = resource_data // 신규로 들어온 데이터의 맵

		node_map, _ := common.ResourceMap.Load("node")
		node_cpu_map, _ := common.ResourceMap.Load("node_cpu")
		node_memory_map, _ := common.ResourceMap.Load("node_memory")
		node_cluster_map, _ := common.ResourceMap.Load("node_cluster")
		node_label_map, _ := common.ResourceMap.Load("node_label")

		node_key := fmt.Sprintf("%s/%s", resource_data.Host, resource_data.Name)
		node_map.(*sync.Map).Store(node_key, resource_data.UID)
		node_cpu_map.(*sync.Map).Store(resource_data.UID, int(resource_data.CPUCount))
		node_memory_map.(*sync.Map).Store(resource_data.UID, resource_data.MemorySize)
		node_cluster_map.(*sync.Map).Store(resource_data.UID, common.ClusterID[resource_data.Host])
		node_label_map.(*sync.Map).Store(resource_data.UID, resource_data.Labels)
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		mapNodeInfo := apiresource.node
		if len(mapNodeInfo) > 0 {
			i_data = tempNodeInfo
			map_info["node_update"] = i_data
			ChannelResourceInsert <- map_info
		} else {
			i_data = ArrResource
			map_info["node_insert"] = i_data
			ChannelResourceInsert <- map_info
			if len(mapNodeInfo) == 0 {
				for _, resource_data := range map_data {
					mapNodeInfo[resource_data.UID] = resource_data
				}
			}
		}
	}
}

func setPodsinfo(info_msg []byte) map[string][]kubeapi.MappingContainer {
	var PodContainerinfo map[string][]kubeapi.MappingContainer = make(map[string][]kubeapi.MappingContainer)
	var tempPodInfo map[string]kubeapi.MappingPod = make(map[string]kubeapi.MappingPod)
	var map_data []kubeapi.MappingPod
	var host string

	err := json.Unmarshal([]byte(info_msg), &map_data)
	errorCheck(err)

	var ArrResource Podinfo
	for _, resource_data := range map_data {
		host = resource_data.Host

		ArrResource.ArrUid = append(ArrResource.ArrUid, resource_data.UID)
		ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[resource_data.Host])
		ArrResource.ArrNodeUid = append(ArrResource.ArrNodeUid, getUID("node", resource_data.Host, resource_data.NodeName))       // name 갖고 UID를 가져와야 함
		ArrResource.ArrNsUid = append(ArrResource.ArrNsUid, getUID("namespace", resource_data.Host, resource_data.NamespaceName)) // name 갖고 UID를 가져와야 함
		ArrResource.ArrAnnotationuid = append(ArrResource.ArrAnnotationuid, resource_data.AnnotationUID)
		ArrResource.ArrPodname = append(ArrResource.ArrPodname, resource_data.Name)
		ArrResource.ArrStarttime = append(ArrResource.ArrStarttime, getStarttime(resource_data.StartTime.Unix(), biastime))
		ArrResource.ArrLabels = append(ArrResource.ArrLabels, resource_data.Labels)
		ArrResource.ArrSelector = append(ArrResource.ArrSelector, resource_data.Selector)
		ArrResource.ArrRestartpolicy = append(ArrResource.ArrRestartpolicy, resource_data.RestartPolicy)
		ArrResource.ArrServiceaccount = append(ArrResource.ArrServiceaccount, resource_data.ServiceAccount)
		ArrResource.ArrStatus = append(ArrResource.ArrStatus, resource_data.Status)
		ArrResource.ArrHostip = append(ArrResource.ArrHostip, resource_data.HostIP)
		ArrResource.ArrPodip = append(ArrResource.ArrPodip, resource_data.PodIP)
		ArrResource.ArrRestartcount = append(ArrResource.ArrRestartcount, int64(resource_data.RestartCount))
		ArrResource.ArrRestarttime = append(ArrResource.ArrRestarttime, getStarttime(resource_data.RestartTime.Unix(), biastime))

		ArrResource.ArrCondition = append(ArrResource.ArrCondition, resource_data.Condition)
		ArrResource.ArrStaticpod = append(ArrResource.ArrStaticpod, resource_data.StaticPod)
		ArrResource.ArrRefkind = append(ArrResource.ArrRefkind, resource_data.ReferenceKind)
		ArrResource.ArrRefuid = append(ArrResource.ArrRefuid, resource_data.ReferenceUID)
		ArrResource.ArrPvcuid = append(ArrResource.ArrPvcuid, getUID("persistentvolumeclaim", resource_data.NamespaceName, resource_data.PersistentVolumeClaim))

		ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1) // Label도 따로 테이블에 저장해야 하는건가?

		ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
		ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)

		var arr_container []kubeapi.MappingContainer
		arr_container = append(arr_container, resource_data.Containers...)
		PodContainerinfo[resource_data.UID] = arr_container //Pod별 Container 배열을 map에 저장

		resource_data.StartTime = time.Unix(getStarttime(resource_data.StartTime.Unix(), biastime), 0)
		tempPodInfo[resource_data.UID] = resource_data // 신규로 들어온 데이터의 맵

		pod_map, _ := common.ResourceMap.Load("pod")
		pod_status_map, _ := common.ResourceMap.Load("pod_status")
		pod_label_map, _ := common.ResourceMap.Load("pod_label")
		pod_reference_map, _ := common.ResourceMap.Load("pod_reference")

		pod_key := fmt.Sprintf("%s/%s", resource_data.NamespaceName, resource_data.Name)
		pod_map.(*sync.Map).Store(pod_key, resource_data.UID)
		pod_status_map.(*sync.Map).Store(resource_data.UID, resource_data.Status)
		pod_label_map.(*sync.Map).Store(resource_data.UID, resource_data.Labels)
		pod_reference_map.(*sync.Map).Store(resource_data.UID, map[string]string{
			resource_data.ReferenceKind: resource_data.ReferenceUID,
		})
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		mapPodInfo := apiresource.pod
		if len(mapPodInfo) > 0 {
			i_data = tempPodInfo
			map_info["pods_update"] = i_data
			ChannelResourceInsert <- map_info
		} else {
			i_data = ArrResource
			map_info["pods_insert"] = i_data
			ChannelResourceInsert <- map_info
			if len(mapPodInfo) == 0 {
				for _, resource_data := range map_data {
					mapPodInfo[resource_data.UID] = resource_data
				}
			}
		}
	}

	return PodContainerinfo
}

func setContainerinfo(info_msg map[string][]kubeapi.MappingContainer) {
	var tempContainerInfo map[string]kubeapi.MappingContainer = make(map[string]kubeapi.MappingContainer)
	var ArrResource Containerinfo
	var host string

	for _, resource_data := range info_msg {
		for _, con_data := range resource_data {
			host = con_data.Host
			container_key := con_data.UID + ":" + con_data.Name

			ArrResource.ArrPodUid = append(ArrResource.ArrPodUid, con_data.UID)
			ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[con_data.Host])
			ArrResource.ArrContainername = append(ArrResource.ArrContainername, con_data.Name)
			ArrResource.ArrImage = append(ArrResource.ArrImage, con_data.Image)
			ArrResource.ArrPorts = append(ArrResource.ArrPorts, con_data.Ports)
			ArrResource.ArrEnv = append(ArrResource.ArrEnv, con_data.Env)
			ArrResource.ArrLimitCpu = append(ArrResource.ArrLimitCpu, con_data.LimitCpu)
			ArrResource.ArrLimitMemory = append(ArrResource.ArrLimitMemory, con_data.LimitMemory)
			ArrResource.ArrLimitStorage = append(ArrResource.ArrLimitStorage, con_data.LimitStorage)
			ArrResource.ArrLimitEphemeral = append(ArrResource.ArrLimitEphemeral, con_data.LimitEphemeral)
			ArrResource.ArrRequestCpu = append(ArrResource.ArrRequestCpu, con_data.RequestCpu)
			ArrResource.ArrRequestMemory = append(ArrResource.ArrRequestMemory, con_data.RequestMemory)
			ArrResource.ArrRequestStorage = append(ArrResource.ArrRequestStorage, con_data.RequestStorage)
			ArrResource.ArrRequestEphemeral = append(ArrResource.ArrRequestEphemeral, con_data.RequestEphemeral)
			ArrResource.ArrVolumemounts = append(ArrResource.ArrVolumemounts, con_data.VolumeMounts)
			ArrResource.ArrState = append(ArrResource.ArrState, con_data.State)
			ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1)
			ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
			ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)

			tempContainerInfo[container_key] = con_data // 신규로 들어온 데이터의 맵

			container_status_map, _ := common.ResourceMap.LoadOrStore("container_status", &sync.Map{})
			container_cpu_request_map, _ := common.ResourceMap.LoadOrStore("container_cpu_request", &sync.Map{})
			container_cpu_limit_map, _ := common.ResourceMap.LoadOrStore("container_cpu_limit", &sync.Map{})
			container_memory_request_map, _ := common.ResourceMap.LoadOrStore("container_memory_request", &sync.Map{})
			container_memory_limit_map, _ := common.ResourceMap.LoadOrStore("container_memory_limit", &sync.Map{})

			container_status_map.(*sync.Map).Store(container_key, con_data.State)
			container_cpu_request_map.(*sync.Map).Store(container_key, con_data.RequestCpu)
			container_cpu_limit_map.(*sync.Map).Store(container_key, con_data.LimitCpu)
			container_memory_request_map.(*sync.Map).Store(container_key, con_data.RequestMemory)
			container_memory_limit_map.(*sync.Map).Store(container_key, con_data.LimitMemory)
		}
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		mapContainerInfo := apiresource.container
		if len(mapContainerInfo) > 0 {
			i_data = tempContainerInfo
			map_info["container_update"] = i_data
			ChannelResourceInsert <- map_info
		} else {
			i_data = ArrResource
			map_info["container_insert"] = i_data
			ChannelResourceInsert <- map_info

			if len(mapContainerInfo) == 0 {
				for _, resource_data := range info_msg {
					for _, con_data := range resource_data {
						container_key := con_data.UID + ":" + con_data.Name
						mapContainerInfo[container_key] = con_data
					}
				}
			}
		}
	}
}

func setServiceinfo(info_msg []byte) {
	var tempServiceInfo map[string]kubeapi.MappingService = make(map[string]kubeapi.MappingService)
	var map_data []kubeapi.MappingService
	var ArrResource Serviceinfo
	var host string

	err := json.Unmarshal([]byte(info_msg), &map_data)
	errorCheck(err)

	for _, resource_data := range map_data {
		host = resource_data.Host

		ArrResource.ArrNsuid = append(ArrResource.ArrNsuid, getUID("namespace", resource_data.Host, resource_data.NamespaceName))
		ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[resource_data.Host])
		ArrResource.ArrSvcname = append(ArrResource.ArrSvcname, resource_data.Name)
		ArrResource.ArrUid = append(ArrResource.ArrUid, resource_data.UID)
		ArrResource.ArrStarttime = append(ArrResource.ArrStarttime, getStarttime(resource_data.StartTime.Unix(), biastime))
		ArrResource.ArrLabels = append(ArrResource.ArrLabels, resource_data.Labels)
		ArrResource.ArrSelector = append(ArrResource.ArrSelector, resource_data.Selector)
		ArrResource.ArrServicetype = append(ArrResource.ArrServicetype, resource_data.ServiceType)
		ArrResource.ArrClusterip = append(ArrResource.ArrClusterip, resource_data.ClusterIP)
		ArrResource.ArrPorts = append(ArrResource.ArrPorts, resource_data.Ports)
		ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1)
		ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
		ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)

		resource_data.StartTime = time.Unix(getStarttime(resource_data.StartTime.Unix(), biastime), 0)
		tempServiceInfo[resource_data.UID] = resource_data // 신규로 들어온 데이터의 맵

		svc_map, _ := common.ResourceMap.Load("service")
		svc_label_map, _ := common.ResourceMap.Load("service_label")
		svc_selector_map, _ := common.ResourceMap.Load("service_selector")

		svc_key := fmt.Sprintf("%s:%s", resource_data.NamespaceName, resource_data.Name)
		svc_map.(*sync.Map).Store(svc_key, resource_data.UID)
		svc_label_map.(*sync.Map).Store(svc_key, resource_data.Labels)
		svc_selector_map.(*sync.Map).Store(svc_key, resource_data.Selector)
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		mapServiceInfo := apiresource.service
		if len(mapServiceInfo) > 0 {
			i_data = tempServiceInfo
			map_info["service_update"] = i_data
			ChannelResourceInsert <- map_info
		} else {
			i_data = ArrResource
			map_info["service_insert"] = i_data
			ChannelResourceInsert <- map_info

			if len(mapServiceInfo) == 0 {
				for _, resource_data := range map_data {
					mapServiceInfo[resource_data.UID] = resource_data
				}
			}
		}
	}
}

func setPvcinfo(info_msg []byte) {
	var tempPvcInfo map[string]kubeapi.MappingPvc = make(map[string]kubeapi.MappingPvc)
	var map_data []kubeapi.MappingPvc
	var ArrResource Pvcinfo
	var host string

	err := json.Unmarshal([]byte(info_msg), &map_data)
	errorCheck(err)

	for _, resource_data := range map_data {
		host = resource_data.Host

		ArrResource.ArrNsuid = append(ArrResource.ArrNsuid, getUID("namespace", resource_data.Host, resource_data.NamespaceName))
		ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[resource_data.Host])
		ArrResource.ArrPvcname = append(ArrResource.ArrPvcname, resource_data.Name)
		ArrResource.ArrPvcUid = append(ArrResource.ArrPvcUid, resource_data.UID)
		ArrResource.ArrStarttime = append(ArrResource.ArrStarttime, getStarttime(resource_data.StartTime.Unix(), biastime))
		ArrResource.ArrLabels = append(ArrResource.ArrLabels, resource_data.Labels)
		ArrResource.ArrSelector = append(ArrResource.ArrSelector, resource_data.Selector)
		ArrResource.ArrAccessmodes = append(ArrResource.ArrAccessmodes, resource_data.AccessModes)
		ArrResource.ArrReqstorage = append(ArrResource.ArrReqstorage, resource_data.RequestStorage)
		ArrResource.ArrStatus = append(ArrResource.ArrStatus, resource_data.Status)
		ArrResource.ArrScuid = append(ArrResource.ArrScuid, getUID("storageclass", resource_data.Host, resource_data.StorageClassName))
		ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1)
		ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
		ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)

		resource_data.StartTime = time.Unix(getStarttime(resource_data.StartTime.Unix(), biastime), 0)
		tempPvcInfo[resource_data.UID] = resource_data // 신규로 들어온 데이터의 맵

		pvc_map, _ := common.ResourceMap.Load("persistentvolumeclaim")
		pvc_label_map, _ := common.ResourceMap.Load("persistentvolumeclaim_label")

		pvc_key := fmt.Sprintf("%s:%s", resource_data.NamespaceName, resource_data.Name)
		pvc_map.(*sync.Map).Store(pvc_key, resource_data.UID)
		pvc_label_map.(*sync.Map).Store(pvc_key, resource_data.Labels)
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		mapPvcInfo := apiresource.persistentvolumeclaim

		if len(mapPvcInfo) > 0 {
			i_data = tempPvcInfo
			map_info["pvc_update"] = i_data
			ChannelResourceInsert <- map_info
		} else {
			i_data = ArrResource
			map_info["pvc_insert"] = i_data
			ChannelResourceInsert <- map_info

			if len(mapPvcInfo) == 0 {
				for _, resource_data := range map_data {
					mapPvcInfo[resource_data.UID] = resource_data
				}
			}
		}
	}
}

func setPvinfo(info_msg []byte) {
	var tempPvInfo map[string]kubeapi.MappingPv = make(map[string]kubeapi.MappingPv)
	var map_data []kubeapi.MappingPv
	var ArrResource Pvinfo
	var host string

	err := json.Unmarshal([]byte(info_msg), &map_data)
	errorCheck(err)

	for _, resource_data := range map_data {
		host = resource_data.Host
		ArrResource.ArrPvname = append(ArrResource.ArrPvname, resource_data.Name)
		ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[resource_data.Host])
		ArrResource.ArrPvUid = append(ArrResource.ArrPvUid, resource_data.UID)
		ArrResource.ArrPvcUid = append(ArrResource.ArrPvcUid, resource_data.PvcUID)
		ArrResource.ArrStarttime = append(ArrResource.ArrStarttime, getStarttime(resource_data.StartTime.Unix(), biastime))
		ArrResource.ArrLabels = append(ArrResource.ArrLabels, resource_data.Labels)
		ArrResource.ArrAccessmodes = append(ArrResource.ArrAccessmodes, resource_data.AccessModes)
		ArrResource.ArrCapacity = append(ArrResource.ArrCapacity, resource_data.Capacity)
		ArrResource.ArrReclaimpolicy = append(ArrResource.ArrReclaimpolicy, resource_data.ReclaimPolicy)
		ArrResource.ArrStatus = append(ArrResource.ArrStatus, resource_data.Status)
		ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1)
		ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
		ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)

		resource_data.StartTime = time.Unix(getStarttime(resource_data.StartTime.Unix(), biastime), 0)
		tempPvInfo[resource_data.UID] = resource_data // 신규로 들어온 데이터의 맵

		pv_map, _ := common.ResourceMap.Load("persistentvolume")
		pv_label_map, _ := common.ResourceMap.Load("persistentvolume_label")

		pv_key := fmt.Sprintf("%s:%s", resource_data.Host, resource_data.Name)
		pv_map.(*sync.Map).Store(pv_key, resource_data.UID)
		pv_label_map.(*sync.Map).Store(pv_key, resource_data.Labels)
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		mapPvInfo := apiresource.persistentvolume
		if len(mapPvInfo) > 0 {
			i_data = tempPvInfo
			map_info["pv_update"] = i_data
			ChannelResourceInsert <- map_info
		} else {
			i_data = ArrResource
			map_info["pv_insert"] = i_data
			ChannelResourceInsert <- map_info

			if len(mapPvInfo) == 0 {
				for _, resource_data := range map_data {
					mapPvInfo[resource_data.UID] = resource_data
				}
			}
		}
	}
}

func setDeployinfo(info_msg []byte) {
	var tempDeployInfo map[string]kubeapi.MappingDeployment = make(map[string]kubeapi.MappingDeployment)
	var map_data []kubeapi.MappingDeployment
	var ArrResource Deployinfo
	var host string

	err := json.Unmarshal([]byte(info_msg), &map_data)
	errorCheck(err)

	for _, resource_data := range map_data {
		host = resource_data.Host

		ArrResource.ArrNsuid = append(ArrResource.ArrNsuid, getUID("namespace", resource_data.Host, resource_data.NamespaceName))
		ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[resource_data.Host])
		ArrResource.ArrDeployname = append(ArrResource.ArrDeployname, resource_data.Name)
		ArrResource.ArrUid = append(ArrResource.ArrUid, resource_data.UID)
		ArrResource.ArrStarttime = append(ArrResource.ArrStarttime, getStarttime(resource_data.StartTime.Unix(), biastime))
		ArrResource.ArrLabels = append(ArrResource.ArrLabels, resource_data.Labels)
		ArrResource.ArrSelector = append(ArrResource.ArrSelector, resource_data.Selector)
		ArrResource.ArrServiceaccount = append(ArrResource.ArrServiceaccount, resource_data.ServiceAccount)
		ArrResource.ArrReplicas = append(ArrResource.ArrReplicas, int64(resource_data.Replicas))
		ArrResource.ArrUpdatedrs = append(ArrResource.ArrUpdatedrs, int64(resource_data.UpdatedReplicas))
		ArrResource.ArrReadyrs = append(ArrResource.ArrReadyrs, int64(resource_data.ReadyReplicas))
		ArrResource.ArrAvailablers = append(ArrResource.ArrAvailablers, int64(resource_data.AvailableReplicas))
		ArrResource.ArrObservedgen = append(ArrResource.ArrObservedgen, resource_data.ObservedGeneneration)
		ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1)
		ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
		ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)

		resource_data.StartTime = time.Unix(getStarttime(resource_data.StartTime.Unix(), biastime), 0)
		tempDeployInfo[resource_data.UID] = resource_data // 신규로 들어온 데이터의 맵

		deploy_map, _ := common.ResourceMap.Load("deployment")
		deploy_label_map, _ := common.ResourceMap.Load("deployment_label")
		deploy_selector_map, _ := common.ResourceMap.Load("deployment_selector")

		deploy_key := fmt.Sprintf("%s:%s", resource_data.NamespaceName, resource_data.Name)
		deploy_map.(*sync.Map).Store(deploy_key, resource_data.UID)
		deploy_label_map.(*sync.Map).Store(deploy_key, resource_data.Labels)
		deploy_selector_map.(*sync.Map).Store(deploy_key, resource_data.Selector)
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		mapDeployInfo := apiresource.deployment

		if len(mapDeployInfo) > 0 {
			i_data = tempDeployInfo
			map_info["deploy_update"] = i_data
			ChannelResourceInsert <- map_info
		} else {
			i_data = ArrResource
			map_info["deploy_insert"] = i_data
			ChannelResourceInsert <- map_info

			if len(mapDeployInfo) == 0 {
				for _, resource_data := range map_data {
					mapDeployInfo[resource_data.UID] = resource_data
				}
			}
		}
	}
}

func setStatefulSetinfo(info_msg []byte) {
	var tempStatefulSetInfo map[string]kubeapi.MappingStatefulSet = make(map[string]kubeapi.MappingStatefulSet)
	var map_data []kubeapi.MappingStatefulSet
	var ArrResource StateFulSetinfo
	var host string

	err := json.Unmarshal([]byte(info_msg), &map_data)
	errorCheck(err)

	for _, resource_data := range map_data {
		host = resource_data.Host

		ArrResource.ArrNsuid = append(ArrResource.ArrNsuid, getUID("namespace", resource_data.Host, resource_data.NamespaceName))
		ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[resource_data.Host])
		ArrResource.ArrStsname = append(ArrResource.ArrStsname, resource_data.Name)
		ArrResource.ArrUid = append(ArrResource.ArrUid, resource_data.UID)
		ArrResource.ArrStarttime = append(ArrResource.ArrStarttime, getStarttime(resource_data.StartTime.Unix(), biastime))
		ArrResource.ArrLabels = append(ArrResource.ArrLabels, resource_data.Labels)
		ArrResource.ArrSelector = append(ArrResource.ArrSelector, resource_data.Selector)
		ArrResource.ArrServiceaccount = append(ArrResource.ArrServiceaccount, resource_data.ServiceAccount)
		ArrResource.ArrReplicas = append(ArrResource.ArrReplicas, int64(resource_data.Replicas))
		ArrResource.ArrUpdatedrs = append(ArrResource.ArrUpdatedrs, int64(resource_data.UpdatedReplicas))
		ArrResource.ArrReadyrs = append(ArrResource.ArrReadyrs, int64(resource_data.ReadyReplicas))
		ArrResource.ArrAvailablers = append(ArrResource.ArrAvailablers, int64(resource_data.AvailableReplicas))
		ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1)
		ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
		ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)

		resource_data.StartTime = time.Unix(getStarttime(resource_data.StartTime.Unix(), biastime), 0)
		tempStatefulSetInfo[resource_data.UID] = resource_data // 신규로 들어온 데이터의 맵

		sts_map, _ := common.ResourceMap.Load("statefulset")
		sts_label_map, _ := common.ResourceMap.Load("statefulset_label")
		sts_selector_map, _ := common.ResourceMap.Load("statefulset_selector")

		sts_key := fmt.Sprintf("%s:%s", resource_data.NamespaceName, resource_data.Name)
		sts_map.(*sync.Map).Store(sts_key, resource_data.UID)
		sts_label_map.(*sync.Map).Store(sts_key, resource_data.Labels)
		sts_selector_map.(*sync.Map).Store(sts_key, resource_data.Selector)
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		mapStatefulInfo := apiresource.statefulset

		if len(mapStatefulInfo) > 0 {
			i_data = tempStatefulSetInfo
			map_info["stateful_update"] = i_data
			ChannelResourceInsert <- map_info
		} else {
			i_data = ArrResource
			map_info["stateful_insert"] = i_data
			ChannelResourceInsert <- map_info

			if len(mapStatefulInfo) == 0 {
				for _, resource_data := range map_data {
					mapStatefulInfo[resource_data.UID] = resource_data
				}
			}
		}
	}
}

func setDaemonSetinfo(info_msg []byte) {
	var tempDaemonSetInfo map[string]kubeapi.MappingDaemonSet = make(map[string]kubeapi.MappingDaemonSet)
	var map_data []kubeapi.MappingDaemonSet
	var ArrResource DaemonSetinfo
	var host string

	err := json.Unmarshal([]byte(info_msg), &map_data)
	errorCheck(err)

	for _, resource_data := range map_data {
		host = resource_data.Host

		ArrResource.ArrNsuid = append(ArrResource.ArrNsuid, getUID("namespace", resource_data.Host, resource_data.NamespaceName))
		ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[resource_data.Host])
		ArrResource.ArrDsname = append(ArrResource.ArrDsname, resource_data.Name)
		ArrResource.ArrUid = append(ArrResource.ArrUid, resource_data.UID)
		ArrResource.ArrStarttime = append(ArrResource.ArrStarttime, getStarttime(resource_data.StartTime.Unix(), biastime))
		ArrResource.ArrLabels = append(ArrResource.ArrLabels, resource_data.Labels)
		ArrResource.ArrSelector = append(ArrResource.ArrSelector, resource_data.Selector)
		ArrResource.ArrServiceaccount = append(ArrResource.ArrServiceaccount, resource_data.ServiceAccount)
		ArrResource.ArrCurrent = append(ArrResource.ArrCurrent, int64(resource_data.CurrentNumberScheduled))
		ArrResource.ArrDesired = append(ArrResource.ArrDesired, int64(resource_data.DesiredNumberScheduled))
		ArrResource.ArrReady = append(ArrResource.ArrReady, int64(resource_data.NumberReady))
		ArrResource.ArrUpdated = append(ArrResource.ArrUpdated, int64(resource_data.UpdatedNumberScheduled))
		ArrResource.ArrAvailable = append(ArrResource.ArrAvailable, int64(resource_data.NumberAvailable))
		ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1)
		ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
		ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)

		resource_data.StartTime = time.Unix(getStarttime(resource_data.StartTime.Unix(), biastime), 0)
		tempDaemonSetInfo[resource_data.UID] = resource_data // 신규로 들어온 데이터의 맵

		ds_map, _ := common.ResourceMap.Load("daemonset")
		ds_label_map, _ := common.ResourceMap.Load("daemonset_label")
		ds_selector_map, _ := common.ResourceMap.Load("daemonset_selector")

		ds_key := fmt.Sprintf("%s:%s", resource_data.NamespaceName, resource_data.Name)
		ds_map.(*sync.Map).Store(ds_key, resource_data.UID)
		ds_label_map.(*sync.Map).Store(ds_key, resource_data.Labels)
		ds_selector_map.(*sync.Map).Store(ds_key, resource_data.Selector)
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		mapDaemonSetInfo := apiresource.daemonset

		if len(mapDaemonSetInfo) > 0 {
			i_data = tempDaemonSetInfo
			map_info["daemonset_update"] = i_data
			ChannelResourceInsert <- map_info
		} else {
			i_data = ArrResource
			map_info["daemonset_insert"] = i_data
			ChannelResourceInsert <- map_info

			if len(mapDaemonSetInfo) == 0 {
				for _, resource_data := range map_data {
					mapDaemonSetInfo[resource_data.UID] = resource_data
				}
			}
		}
	}
}

func setReplicaSetinfo(info_msg []byte) {
	var tempReplicaSetInfo map[string]kubeapi.MappingReplicaSet = make(map[string]kubeapi.MappingReplicaSet)
	var map_data []kubeapi.MappingReplicaSet
	var ArrResource ReplicaSetinfo
	var host string

	err := json.Unmarshal([]byte(info_msg), &map_data)
	errorCheck(err)

	for _, resource_data := range map_data {
		host = resource_data.Host

		ArrResource.ArrNsuid = append(ArrResource.ArrNsuid, getUID("namespace", resource_data.Host, resource_data.NamespaceName))
		ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[resource_data.Host])
		ArrResource.ArrRsname = append(ArrResource.ArrRsname, resource_data.Name)
		ArrResource.ArrUid = append(ArrResource.ArrUid, resource_data.UID)
		ArrResource.ArrStarttime = append(ArrResource.ArrStarttime, getStarttime(resource_data.StartTime.Unix(), biastime))
		ArrResource.ArrLabels = append(ArrResource.ArrLabels, resource_data.Labels)
		ArrResource.ArrSelector = append(ArrResource.ArrSelector, resource_data.Selector)
		ArrResource.ArrReplicas = append(ArrResource.ArrReplicas, int64(resource_data.Replicas))
		ArrResource.ArrFullylabeldrs = append(ArrResource.ArrFullylabeldrs, int64(resource_data.FullyLabeledReplicas))
		ArrResource.ArrReadyrs = append(ArrResource.ArrReadyrs, int64(resource_data.ReadyReplicas))
		ArrResource.ArrAvailablers = append(ArrResource.ArrAvailablers, int64(resource_data.AvailableReplicas))
		ArrResource.ArrObservedgen = append(ArrResource.ArrObservedgen, int64(resource_data.ObservedGeneneration))
		ArrResource.ArrRefkind = append(ArrResource.ArrRefkind, resource_data.ReferenceKind)
		ArrResource.ArrRefuid = append(ArrResource.ArrRefuid, resource_data.ReferenceUID)
		ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1)
		ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
		ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)

		resource_data.StartTime = time.Unix(getStarttime(resource_data.StartTime.Unix(), biastime), 0)
		tempReplicaSetInfo[resource_data.UID] = resource_data // 신규로 들어온 데이터의 맵

		rs_map, _ := common.ResourceMap.Load("replicaset")
		rs_label_map, _ := common.ResourceMap.Load("replicaset_label")
		rs_selector_map, _ := common.ResourceMap.Load("replicaset_selector")
		rs_reference_map, _ := common.ResourceMap.Load("replicaset_reference")

		rs_key := fmt.Sprintf("%s:%s", resource_data.NamespaceName, resource_data.Name)
		rs_map.(*sync.Map).Store(rs_key, resource_data.UID)
		rs_label_map.(*sync.Map).Store(rs_key, resource_data.Labels)
		rs_selector_map.(*sync.Map).Store(rs_key, resource_data.Selector)
		rs_reference_map.(*sync.Map).Store(resource_data.UID, map[string]string{
			resource_data.ReferenceKind: resource_data.ReferenceUID,
		})
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		mapReplicaSetInfo := apiresource.replicaset
		if len(mapReplicaSetInfo) > 0 {
			i_data = tempReplicaSetInfo
			map_info["replicaset_update"] = i_data
			ChannelResourceInsert <- map_info
		} else {
			i_data = ArrResource
			map_info["replicaset_insert"] = i_data
			ChannelResourceInsert <- map_info

			if len(mapReplicaSetInfo) == 0 {
				for _, resource_data := range map_data {
					mapReplicaSetInfo[resource_data.UID] = resource_data
				}
			}
		}
	}
}

func setIngressinfo(info_msg []byte) map[string][]kubeapi.MappingIngressHost {
	var ingress_host_info map[string][]kubeapi.MappingIngressHost = make(map[string][]kubeapi.MappingIngressHost)
	var tempIngressInfo map[string]kubeapi.MappingIngress = make(map[string]kubeapi.MappingIngress)
	var map_data []kubeapi.MappingIngress
	var ArrResource Inginfo
	var host string

	err := json.Unmarshal([]byte(info_msg), &map_data)
	errorCheck(err)

	for _, resource_data := range map_data {
		host = resource_data.Host

		ArrResource.ArrUid = append(ArrResource.ArrUid, resource_data.UID)
		ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[resource_data.Host])
		ArrResource.ArrNsUid = append(ArrResource.ArrNsUid, getUID("namespace", resource_data.Host, resource_data.NamespaceName))
		ArrResource.ArrName = append(ArrResource.ArrName, resource_data.Name)
		ArrResource.ArrStarttime = append(ArrResource.ArrStarttime, getStarttime(resource_data.StartTime.Unix(), biastime))
		ArrResource.ArrLabels = append(ArrResource.ArrLabels, resource_data.Labels)
		ArrResource.ArrClassname = append(ArrResource.ArrClassname, resource_data.IngressClassName)
		ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1)
		ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
		ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)
		var arr_ingHost []kubeapi.MappingIngressHost
		arr_ingHost = append(arr_ingHost, resource_data.IngressHosts...)
		ingress_host_info[resource_data.UID] = arr_ingHost //Ingress별 Host 배열을 map에 저장

		resource_data.StartTime = time.Unix(getStarttime(resource_data.StartTime.Unix(), biastime), 0)
		tempIngressInfo[resource_data.UID] = resource_data // 신규로 들어온 데이터의 맵

		ing_map, _ := common.ResourceMap.Load("ingress")
		ing_label_map, _ := common.ResourceMap.Load("ingress_label")

		ing_key := fmt.Sprintf("%s:%s", resource_data.NamespaceName, resource_data.Name)
		ing_map.(*sync.Map).Store(ing_key, resource_data.UID)
		ing_label_map.(*sync.Map).Store(ing_key, resource_data.Labels)
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		mapIngInfo := apiresource.ingress

		if len(mapIngInfo) > 0 {
			i_data = tempIngressInfo
			map_info["ing_update"] = i_data
			ChannelResourceInsert <- map_info
		} else {
			i_data = ArrResource
			map_info["ing_insert"] = i_data
			ChannelResourceInsert <- map_info
			if len(mapIngInfo) == 0 {
				for _, resource_data := range map_data {
					mapIngInfo[resource_data.UID] = resource_data
				}
			}
		}
	}

	return ingress_host_info
}

func setIngressHostinfo(info_msg map[string][]kubeapi.MappingIngressHost) {
	var tempIngHostInfo map[string]kubeapi.MappingIngressHost = make(map[string]kubeapi.MappingIngressHost)
	var ArrResource IngHostinfo
	var host string

	for key, resource_data := range info_msg {
		for _, con_data := range resource_data {
			host = con_data.Host

			ArrResource.ArrInguid = append(ArrResource.ArrInguid, key)
			ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[con_data.Host])
			ArrResource.ArrBackendtype = append(ArrResource.ArrBackendtype, con_data.BackendType)
			ArrResource.ArrBackendname = append(ArrResource.ArrBackendname, con_data.BackendName)
			ArrResource.ArrHostname = append(ArrResource.ArrHostname, con_data.Hostname)
			ArrResource.ArrPathtype = append(ArrResource.ArrPathtype, con_data.PathType)
			ArrResource.ArrPath = append(ArrResource.ArrPath, con_data.Path)
			ArrResource.ArrRscApiGroup = append(ArrResource.ArrRscApiGroup, con_data.ResourceAPIGroup)
			ArrResource.ArrRsckind = append(ArrResource.ArrRsckind, con_data.ResourceKind)
			ArrResource.ArrSvcport = append(ArrResource.ArrSvcport, int(con_data.ServicePort))
			ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1)
			ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
			ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)

			tempIngHostInfo[con_data.Hostname] = con_data // 신규로 들어온 데이터의 맵
		}
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		mapIngHostInfo := apiresource.ingresshost

		if len(mapIngHostInfo) > 0 {
			i_data = tempIngHostInfo
			map_info["inghost_update"] = i_data
			ChannelResourceInsert <- map_info
		} else {
			i_data = ArrResource
			map_info["inghost_insert"] = i_data
			ChannelResourceInsert <- map_info

			if len(mapIngHostInfo) == 0 {
				for _, resource_data := range info_msg {
					for _, con_data := range resource_data {
						mapIngHostInfo[con_data.Hostname] = con_data
					}
				}
			}
		}
	}
}

func setStorageclassinfo(info_msg []byte) {
	var tempScInfo map[string]kubeapi.MappingStorageClass = make(map[string]kubeapi.MappingStorageClass)
	var map_data []kubeapi.MappingStorageClass
	var ArrResource Scinfo
	var host string

	err := json.Unmarshal([]byte(info_msg), &map_data)
	errorCheck(err)
	for _, resource_data := range map_data {
		host = resource_data.Host

		ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[resource_data.Host])
		ArrResource.ArrScname = append(ArrResource.ArrScname, resource_data.Name)
		ArrResource.ArrUid = append(ArrResource.ArrUid, resource_data.UID)
		ArrResource.ArrStarttime = append(ArrResource.ArrStarttime, getStarttime(resource_data.StartTime.Unix(), biastime))
		ArrResource.ArrLabels = append(ArrResource.ArrLabels, resource_data.Labels)
		ArrResource.ArrProvisionor = append(ArrResource.ArrProvisionor, resource_data.Provisioner)
		ArrResource.ArrReclaimolicy = append(ArrResource.ArrReclaimolicy, resource_data.ReclaimPolicy)
		ArrResource.ArrVolumebindingmode = append(ArrResource.ArrVolumebindingmode, resource_data.VolumeBindingMode)
		VolExt := resource_data.AllowVolumeExpansion
		if VolExt {
			ArrResource.ArrAllowvolumeexp = append(ArrResource.ArrAllowvolumeexp, 1)
		} else {
			ArrResource.ArrAllowvolumeexp = append(ArrResource.ArrAllowvolumeexp, 0)
		}
		ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1)
		ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
		ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)

		resource_data.StartTime = time.Unix(getStarttime(resource_data.StartTime.Unix(), biastime), 0)
		tempScInfo[resource_data.UID] = resource_data // 신규로 들어온 데이터의 맵

		sc_map, _ := common.ResourceMap.Load("storageclass")
		sc_label_map, _ := common.ResourceMap.Load("storageclass_label")

		sc_key := fmt.Sprintf("%s:%s", resource_data.Host, resource_data.Name)
		sc_map.(*sync.Map).Store(sc_key, resource_data.UID)
		sc_label_map.(*sync.Map).Store(sc_key, resource_data.Labels)
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		mapScInfo := apiresource.storageclass

		if len(mapScInfo) > 0 {
			i_data = tempScInfo
			map_info["sc_update"] = i_data
			ChannelResourceInsert <- map_info
		} else {
			i_data = ArrResource
			map_info["sc_insert"] = i_data
			ChannelResourceInsert <- map_info

			if len(mapScInfo) == 0 {
				for _, resource_data := range map_data {
					mapScInfo[resource_data.UID] = resource_data
				}
			}
		}
	}
}
