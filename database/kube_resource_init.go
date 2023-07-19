package database

import (
	"onTuneKubeManager/common"
	"onTuneKubeManager/kubeapi"
	"sync"
	"time"
)

func InitNamespaceInfo(host string, clusterid int) {
	mapNamespaceInfo := make(map[string]kubeapi.MappingNamespace)
	rowsCnt := select_row_count_enabled(TB_KUBE_NS_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapNamespaceInfo) == 0 {
			var uid string
			var name string
			var starttime int64
			var status string
			var labels string
			var enabled int
			rows := select_row_enabled("nsuid, nsname, starttime, labels, status, enabled", TB_KUBE_NS_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&uid, &name, &starttime, &labels, &status, &enabled)
				errorCheck(err)
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingNamespace
					resource_data_temp.UID = uid
					resource_data_temp.Name = name
					resource_data_temp.Host = host
					resource_data_temp.StartTime = time.Unix(starttime, 0)
					resource_data_temp.Labels = labels
					resource_data_temp.Status = status
					mapNamespaceInfo[resource_data_temp.UID] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		apiresource.namespace = mapNamespaceInfo
		mapApiResource.Store(host, apiresource)
	}
}

func InitNodeInfo(host string, clusterid int) {
	mapNodeInfo := make(map[string]kubeapi.MappingNode)
	rowsCnt := select_row_count_enabled(TB_KUBE_NODE_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapNodeInfo) == 0 {
			var name string
			var uid string
			var nodetype string
			var enabled int
			var starttime int64
			var labels string
			var kernelversion string
			var osimage string
			var osname string
			var containerruntimever string
			var kubeletver string
			var kubeproxyver string
			var cpuarch string
			var cpucount int
			var ephemeralstorage int64
			var memorysize int64
			var pods int64
			var ip string
			var status int

			rows := select_row_enabled("nodename, nodeuid, nodetype, enabled, starttime, labels, kernelversion, osimage, osname, containerruntimever, kubeletver, kubeproxyver, cpuarch, cpucount, ephemeralstorage, memorysize, pods, ip, status", TB_KUBE_NODE_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&name, &uid, &nodetype, &enabled, &starttime, &labels, &kernelversion, &osimage, &osname, &containerruntimever, &kubeletver,
					&kubeproxyver, &cpuarch, &cpucount, &ephemeralstorage, &memorysize, &pods, &ip, &status)
				errorCheck(err)
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingNode
					resource_data_temp.UID = uid
					resource_data_temp.Name = name
					resource_data_temp.NodeType = nodetype
					resource_data_temp.Host = host
					resource_data_temp.StartTime = time.Unix(starttime, 0)
					resource_data_temp.Labels = labels
					resource_data_temp.KernelVersion = kernelversion
					resource_data_temp.OSImage = osimage
					resource_data_temp.OSName = osname
					resource_data_temp.ContainerRuntimeVersion = containerruntimever
					resource_data_temp.KubeletVersion = kubeletver
					resource_data_temp.KubeProxyVersion = kubeproxyver
					resource_data_temp.CPUArch = cpuarch
					resource_data_temp.CPUCount = int64(cpucount)
					resource_data_temp.EphemeralStorage = ephemeralstorage
					resource_data_temp.MemorySize = memorysize
					resource_data_temp.Pods = pods
					resource_data_temp.IP = ip
					resource_data_temp.Status = status
					mapNodeInfo[resource_data_temp.UID] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		apiresource.node = mapNodeInfo
		mapApiResource.Store(host, apiresource)
	}
}

func InitPodInfo(host string, clusterid int) {
	mapPodInfo := make(map[string]kubeapi.MappingPod)
	rowsCnt := select_row_count_enabled(TB_KUBE_POD_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapPodInfo) == 0 {
			var podname string
			var uid string
			var nodeuid string
			var nsuid string
			var annotationuid string
			var starttime int64
			var labels string
			var selector string
			var restartpolicy string
			var serviceaccount string
			var status string
			var hostip string
			var podip string
			var restartcount int
			var restarttime int64
			var podcondition string
			var staticpod string
			var refkind string
			var refuid string
			var enabled int

			rows := select_row_enabled("podname, uid, nodeuid, nsuid, annotationuid, starttime, labels, selector, restartpolicy, serviceaccount, status, hostip, podip, restartcount, restarttime, podcondition, staticpod, refkind, refuid, enabled", TB_KUBE_POD_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&podname, &uid, &nodeuid, &nsuid, &annotationuid, &starttime, &labels, &selector, &restartpolicy, &serviceaccount, &status, &hostip, &podip, &restartcount, &restarttime, &podcondition, &staticpod, &refkind, &refuid, &enabled)
				errorCheck(err)
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingPod
					resource_data_temp.UID = uid
					resource_data_temp.Name = podname
					resource_data_temp.NodeName = nodeuid
					resource_data_temp.NamespaceName = nsuid
					resource_data_temp.AnnotationUID = annotationuid
					tm := time.Unix(starttime, 0)
					resource_data_temp.StartTime = tm
					resource_data_temp.Host = host
					resource_data_temp.Labels = labels
					resource_data_temp.Selector = selector
					resource_data_temp.RestartPolicy = restartpolicy
					resource_data_temp.ServiceAccount = serviceaccount
					resource_data_temp.Status = status
					resource_data_temp.HostIP = hostip
					resource_data_temp.PodIP = podip
					resource_data_temp.RestartCount = int32(restartcount)
					restart_tm := time.Unix(restarttime, 0)
					resource_data_temp.RestartTime = restart_tm
					resource_data_temp.Condition = podcondition
					resource_data_temp.StaticPod = staticpod
					resource_data_temp.ReferenceKind = refkind
					resource_data_temp.ReferenceUID = refuid
					mapPodInfo[resource_data_temp.UID] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		apiresource.pod = mapPodInfo
		mapApiResource.Store(host, apiresource)
	}
}

func InitContainerInfo(host string, clusterid int) {
	mapContainerInfo := make(map[string]kubeapi.MappingContainer)
	rowsCnt := select_row_count_enabled(TB_KUBE_CONTAINER_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapContainerInfo) == 0 {
			var poduid string
			var containername string
			var image string
			var ports string
			var env string
			var limitcpu int64
			var limitmemory int64
			var limitstorage int64
			var limitephemeral int64
			var reqcpu int64
			var reqmemory int64
			var reqstorage int64
			var reqephemeral int64
			var volumemount string
			var state string
			var enabled int

			rows := select_row_enabled("poduid, containername, image, ports, env, limitcpu, limitmemory, limitstorage, limitephemeral, reqcpu, reqmemory, reqstorage, reqephemeral, volumemounts, state, enabled", TB_KUBE_CONTAINER_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&poduid, &containername, &image, &ports, &env, &limitcpu, &limitmemory, &limitstorage, &limitephemeral, &reqcpu, &reqmemory, &reqstorage, &reqephemeral, &volumemount, &state, &enabled)
				errorCheck(err)
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingContainer
					resource_data_temp.UID = poduid
					resource_data_temp.Name = containername
					resource_data_temp.Host = host
					resource_data_temp.Image = image
					resource_data_temp.Ports = ports
					resource_data_temp.Env = env
					resource_data_temp.LimitCpu = limitcpu
					resource_data_temp.LimitMemory = limitmemory
					resource_data_temp.LimitStorage = limitstorage
					resource_data_temp.LimitEphemeral = limitephemeral
					resource_data_temp.RequestCpu = reqcpu
					resource_data_temp.RequestMemory = reqmemory
					resource_data_temp.RequestStorage = reqstorage
					resource_data_temp.RequestEphemeral = reqephemeral
					resource_data_temp.VolumeMounts = volumemount
					resource_data_temp.State = state

					container_key := poduid + ":" + containername
					mapContainerInfo[container_key] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		apiresource.container = mapContainerInfo
		mapApiResource.Store(host, apiresource)
	}
}

func InitServiceInfo(host string, clusterid int) {
	mapServiceInfo := make(map[string]kubeapi.MappingService)
	rowsCnt := select_row_count_enabled(TB_KUBE_SVC_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapServiceInfo) == 0 {
			var nsuid string
			var svcname string
			var uid string
			var starttime int64
			var labels string
			var selector string
			var servicetype string
			var clusterip string
			var ports string
			var enabled int
			rows := select_row_enabled("nsuid, svcname, uid, starttime, labels, selector, servicetype, clusterip, ports, enabled", TB_KUBE_SVC_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&nsuid, &svcname, &uid, &starttime, &labels, &selector, &servicetype, &clusterip, &ports, &enabled)
				errorCheck(err)
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingService

					if ns_map, ok := common.ResourceMap.Load("namespace"); ok {
						ns_map.(*sync.Map).Range(func(key, value any) bool {
							if value.(string) == nsuid {
								resource_data_temp.NamespaceName = key.(string)
								return false
							}

							return true
						})
					}

					resource_data_temp.Name = svcname
					resource_data_temp.UID = uid
					tm := time.Unix(starttime, 0)
					resource_data_temp.StartTime = tm
					resource_data_temp.Host = host
					resource_data_temp.Labels = labels
					resource_data_temp.Selector = selector
					resource_data_temp.ServiceType = servicetype
					resource_data_temp.ClusterIP = clusterip
					resource_data_temp.Ports = ports
					mapServiceInfo[resource_data_temp.UID] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		apiresource.service = mapServiceInfo
		mapApiResource.Store(host, apiresource)
	}
}

func InitPersistentVolumeClaimInfo(host string, clusterid int) {
	mapPvcInfo := make(map[string]kubeapi.MappingPvc)
	rowsCnt := select_row_count_enabled(TB_KUBE_PVC_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapPvcInfo) == 0 {
			var nsuid string
			var pvcname string
			var uid string
			var starttime int64
			var labels string
			var selector string
			var accessmodes string
			var status string
			var enabled int
			rows := select_row_enabled("nsuid, pvcname, uid, starttime, labels, selector, accessmodes, status, enabled", TB_KUBE_PVC_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&nsuid, &pvcname, &uid, &starttime, &labels, &selector, &accessmodes, &status, &enabled)
				errorCheck(err)
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingPvc
					if ns_map, ok := common.ResourceMap.Load("namespace"); ok {
						ns_map.(*sync.Map).Range(func(key, value any) bool {
							if value.(string) == nsuid {
								resource_data_temp.NamespaceName = key.(string)
								return false
							}

							return true
						})
					}

					resource_data_temp.Name = pvcname
					resource_data_temp.UID = uid
					tm := time.Unix(starttime, 0)
					resource_data_temp.StartTime = tm
					resource_data_temp.Host = host
					resource_data_temp.Labels = labels
					resource_data_temp.Selector = selector
					resource_data_temp.AccessModes = accessmodes
					resource_data_temp.Status = status
					mapPvcInfo[resource_data_temp.UID] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		apiresource.persistentvolumeclaim = mapPvcInfo
		mapApiResource.Store(host, apiresource)
	}
}

func InitPersistentVolumeInfo(host string, clusterid int) {
	mapPvInfo := make(map[string]kubeapi.MappingPv)
	rowsCnt := select_row_count_enabled(TB_KUBE_PV_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapPvInfo) == 0 {
			var pvname string
			var pvuid string
			var pvcuid string
			var starttime int64
			var labels string
			var accessmodes string
			var reclaimpolicy string
			var status string
			var enabled int
			rows := select_row_enabled("pvname, pvuid, pvcuid, starttime, labels, accessmodes, reclaimpolicy, status, enabled", TB_KUBE_PV_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&pvname, &pvuid, &pvcuid, &starttime, &labels, &accessmodes, &reclaimpolicy, &status, &enabled)
				errorCheck(err)
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingPv
					resource_data_temp.Name = pvname
					resource_data_temp.UID = pvuid
					resource_data_temp.PvcUID = pvcuid
					tm := time.Unix(starttime, 0)
					resource_data_temp.StartTime = tm
					resource_data_temp.Host = host
					resource_data_temp.Labels = labels
					resource_data_temp.AccessModes = accessmodes
					resource_data_temp.ReclaimPolicy = reclaimpolicy
					resource_data_temp.Status = status
					mapPvInfo[resource_data_temp.UID] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		apiresource.persistentvolume = mapPvInfo
		mapApiResource.Store(host, apiresource)
	}
}

func InitDeploymentInfo(host string, clusterid int) {
	mapDeployInfo := make(map[string]kubeapi.MappingDeployment)
	rowsCnt := select_row_count_enabled(TB_KUBE_DEPLOY_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapDeployInfo) == 0 {
			var nsuid string
			var deployname string
			var uid string
			var starttime int64
			var labels string
			var selector string
			var serviceaccount string
			var replicas int64
			var updatedrs int64
			var readyrs int64
			var availablers int64
			var observedgen int64
			var enabled int
			rows := select_row_enabled("nsuid, deployname, uid, starttime, labels, selector, serviceaccount, replicas, updatedrs, readyrs, availablers, observedgen, enabled", TB_KUBE_DEPLOY_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&nsuid, &deployname, &uid, &starttime, &labels, &selector, &serviceaccount, &replicas, &updatedrs, &readyrs, &availablers, &observedgen, &enabled)
				errorCheck(err)
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingDeployment
					if ns_map, ok := common.ResourceMap.Load("namespace"); ok {
						ns_map.(*sync.Map).Range(func(key, value any) bool {
							if value.(string) == nsuid {
								resource_data_temp.NamespaceName = key.(string)
								return false
							}

							return true
						})
					}
					resource_data_temp.Name = deployname
					resource_data_temp.UID = uid
					tm := time.Unix(starttime, 0)
					resource_data_temp.StartTime = tm
					resource_data_temp.Host = host
					resource_data_temp.Labels = labels
					resource_data_temp.Selector = selector
					resource_data_temp.ServiceAccount = serviceaccount
					resource_data_temp.Replicas = int32(replicas)
					resource_data_temp.UpdatedReplicas = int32(updatedrs)
					resource_data_temp.ReadyReplicas = int32(readyrs)
					resource_data_temp.AvailableReplicas = int32(availablers)
					resource_data_temp.ObservedGeneneration = observedgen
					mapDeployInfo[resource_data_temp.UID] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		apiresource.deployment = mapDeployInfo
		mapApiResource.Store(host, apiresource)
	}
}

func InitStatefulSetInfo(host string, clusterid int) {
	mapStatefulInfo := make(map[string]kubeapi.MappingStatefulSet)
	rowsCnt := select_row_count_enabled(TB_KUBE_STS_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapStatefulInfo) == 0 {
			var nsuid string
			var stsname string
			var uid string
			var starttime int64
			var labels string
			var selector string
			var serviceaccount string
			var replicas int64
			var updatedrs int64
			var readyrs int64
			var availablers int64
			var enabled int
			rows := select_row_enabled("nsuid, stsname, uid, starttime, labels, selector, serviceaccount, replicas, updatedrs, readyrs, availablers, enabled", TB_KUBE_STS_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&nsuid, &stsname, &uid, &starttime, &labels, &selector, &serviceaccount, &replicas, &updatedrs, &readyrs, &availablers, &enabled)
				errorCheck(err)
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingStatefulSet
					if ns_map, ok := common.ResourceMap.Load("namespace"); ok {
						ns_map.(*sync.Map).Range(func(key, value any) bool {
							if value.(string) == nsuid {
								resource_data_temp.NamespaceName = key.(string)
								return false
							}

							return true
						})
					}
					resource_data_temp.Name = stsname
					resource_data_temp.UID = uid
					tm := time.Unix(starttime, 0)
					resource_data_temp.StartTime = tm
					resource_data_temp.Host = host
					resource_data_temp.Labels = labels
					resource_data_temp.Selector = selector
					resource_data_temp.ServiceAccount = serviceaccount
					resource_data_temp.Replicas = int32(replicas)
					resource_data_temp.UpdatedReplicas = int32(updatedrs)
					resource_data_temp.ReadyReplicas = int32(readyrs)
					resource_data_temp.AvailableReplicas = int32(availablers)
					mapStatefulInfo[resource_data_temp.UID] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		apiresource.statefulset = mapStatefulInfo
		mapApiResource.Store(host, apiresource)
	}
}

func InitDaemonSetInfo(host string, clusterid int) {
	mapDaemonSetInfo := make(map[string]kubeapi.MappingDaemonSet)
	rowsCnt := select_row_count_enabled(TB_KUBE_DS_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapDaemonSetInfo) == 0 {
			var nsuid string
			var stsname string
			var uid string
			var starttime int64
			var labels string
			var selector string
			var serviceaccount string
			var current int64
			var desired int64
			var ready int64
			var updated int64
			var available int64
			var enabled int
			rows := select_row_enabled("nsuid, dsname, uid, starttime, labels, selector, serviceaccount, current, desired, ready, updated, available, enabled", TB_KUBE_DS_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&nsuid, &stsname, &uid, &starttime, &labels, &selector, &serviceaccount, &current, &desired, &ready, &updated, &available, &enabled)
				errorCheck(err)
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingDaemonSet
					if ns_map, ok := common.ResourceMap.Load("namespace"); ok {
						ns_map.(*sync.Map).Range(func(key, value any) bool {
							if value.(string) == nsuid {
								resource_data_temp.NamespaceName = key.(string)
								return false
							}

							return true
						})
					}

					resource_data_temp.Name = stsname
					resource_data_temp.UID = uid
					tm := time.Unix(starttime, 0)
					resource_data_temp.StartTime = tm
					resource_data_temp.Host = host
					resource_data_temp.Labels = labels
					resource_data_temp.Selector = selector
					resource_data_temp.ServiceAccount = serviceaccount
					resource_data_temp.CurrentNumberScheduled = int32(current)
					resource_data_temp.DesiredNumberScheduled = int32(desired)
					resource_data_temp.NumberReady = int32(ready)
					resource_data_temp.UpdatedNumberScheduled = int32(updated)
					resource_data_temp.NumberAvailable = int32(available)
					mapDaemonSetInfo[resource_data_temp.UID] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		apiresource.daemonset = mapDaemonSetInfo
		mapApiResource.Store(host, apiresource)
	}
}

func InitReplicaSetInfo(host string, clusterid int) {
	mapReplicaSetInfo := make(map[string]kubeapi.MappingReplicaSet)
	rowsCnt := select_row_count_enabled(TB_KUBE_RS_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapReplicaSetInfo) == 0 {
			var nsuid string
			var rsname string
			var uid string
			var starttime int64
			var labels string
			var selector string
			var replicas int64
			var fullylabeledrs int64
			var readyrs int64
			var availablers int64
			var observedgen int64
			var refkind string
			var refuid string
			var enabled int
			rows := select_row_enabled("nsuid, rsname, uid, starttime, labels, selector, replicas, fullylabeledrs, readyrs, availablers, observedgen, refkind, refuid, enabled", TB_KUBE_RS_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&nsuid, &rsname, &uid, &starttime, &labels, &selector, &replicas, &fullylabeledrs, &readyrs, &availablers, &observedgen, &refkind, &refuid, &enabled)
				errorCheck(err)
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingReplicaSet
					if ns_map, ok := common.ResourceMap.Load("namespace"); ok {
						ns_map.(*sync.Map).Range(func(key, value any) bool {
							if value.(string) == nsuid {
								resource_data_temp.NamespaceName = key.(string)
								return false
							}

							return true
						})
					}

					resource_data_temp.Name = rsname
					resource_data_temp.UID = uid
					tm := time.Unix(starttime, 0)
					resource_data_temp.StartTime = tm
					resource_data_temp.Host = host
					resource_data_temp.Labels = labels
					resource_data_temp.Selector = selector
					resource_data_temp.Replicas = int32(replicas)
					resource_data_temp.FullyLabeledReplicas = int32(fullylabeledrs)
					resource_data_temp.ReadyReplicas = int32(readyrs)
					resource_data_temp.AvailableReplicas = int32(availablers)
					resource_data_temp.ObservedGeneneration = observedgen
					resource_data_temp.ReferenceKind = refkind
					resource_data_temp.ReferenceUID = refuid
					mapReplicaSetInfo[resource_data_temp.UID] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		apiresource.replicaset = mapReplicaSetInfo
		mapApiResource.Store(host, apiresource)
	}
}

func InitIngressInfo(host string, clusterid int) {
	mapIngInfo := make(map[string]kubeapi.MappingIngress)
	rowsCnt := select_row_count_enabled(TB_KUBE_ING_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapIngInfo) == 0 {
			var nsuid string
			var ingname string
			var uid string
			var starttime int64
			var labels string
			var classname string
			var enabled int
			rows := select_row_enabled("nsuid, ingname, uid, starttime, labels, classname, enabled", TB_KUBE_ING_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&nsuid, &ingname, &uid, &starttime, &labels, &classname, &enabled)
				errorCheck(err)
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingIngress
					resource_data_temp.Name = ingname
					resource_data_temp.UID = uid
					tm := time.Unix(starttime, 0)
					resource_data_temp.StartTime = tm
					resource_data_temp.Host = host
					resource_data_temp.Labels = labels
					resource_data_temp.IngressClassName = classname
					mapIngInfo[resource_data_temp.UID] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		apiresource.ingress = mapIngInfo
		mapApiResource.Store(host, apiresource)
	}
}

func InitIngressHostInfo(host string, clusterid int) {
	mapIngHostInfo := make(map[string]kubeapi.MappingIngressHost)
	rowsCnt := select_row_count_enabled(TB_KUBE_INGHOST_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapIngHostInfo) == 0 {
			var inguid string
			var backendtype string
			var backendname string
			var hostname string
			var pathtype string
			var path string
			var serviceport int32
			var rscapigroup string
			var rsckind string
			var enabled int
			rows := select_row_enabled("inguid, backendtype, backendname, hostname, pathtype, path, serviceport, rscapigroup, rsckind, enabled", TB_KUBE_INGHOST_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&inguid, &backendtype, &backendname, &hostname, &pathtype, &path, &serviceport, &rscapigroup, &rsckind, &enabled)
				errorCheck(err)
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingIngressHost
					resource_data_temp.UID = inguid
					resource_data_temp.BackendType = backendtype
					resource_data_temp.BackendName = backendname
					resource_data_temp.Hostname = hostname
					resource_data_temp.Host = host
					resource_data_temp.PathType = pathtype
					resource_data_temp.Path = path
					resource_data_temp.ServicePort = serviceport
					resource_data_temp.ResourceAPIGroup = rscapigroup
					resource_data_temp.ResourceKind = rsckind
					mapIngHostInfo[resource_data_temp.Hostname] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		apiresource.ingresshost = mapIngHostInfo
		mapApiResource.Store(host, apiresource)
	}
}

func InitStorageClassInfo(host string, clusterid int) {
	mapScInfo := make(map[string]kubeapi.MappingStorageClass)
	rowsCnt := select_row_count_enabled(TB_KUBE_SC_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapScInfo) == 0 {
			var scname string
			var uid string
			var starttime int64
			var labels string
			var provisioner string
			var reclaimpolicy string
			var volumebindingmode string
			var allowvolumeexp int
			var enabled int
			rows := select_row_enabled("scname, uid, starttime, labels, provisioner, reclaimpolicy, volumebindingmode, allowvolumeexp, enabled", TB_KUBE_SC_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&scname, &uid, &starttime, &labels, &provisioner, &reclaimpolicy, &volumebindingmode, &allowvolumeexp, &enabled)
				errorCheck(err)
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingStorageClass
					resource_data_temp.Name = scname
					resource_data_temp.UID = uid
					tm := time.Unix(starttime, 0)
					resource_data_temp.StartTime = tm
					resource_data_temp.Host = host
					resource_data_temp.Labels = labels
					resource_data_temp.Provisioner = provisioner
					resource_data_temp.ReclaimPolicy = reclaimpolicy
					resource_data_temp.VolumeBindingMode = volumebindingmode
					if allowvolumeexp == 1 {
						resource_data_temp.AllowVolumeExpansion = true
					} else {
						resource_data_temp.AllowVolumeExpansion = false
					}
					mapScInfo[resource_data_temp.UID] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ar, ok := mapApiResource.Load(host); ok {
		apiresource := ar.(*ApiResource)
		apiresource.storageclass = mapScInfo
		mapApiResource.Store(host, apiresource)
	}
}

// ############################ Init Data Load
func InitMapData() {
	defer errorRecover()

	for k, v := range common.ClusterID {
		mapApiResource.Store(k, &ApiResource{
			namespace:             make(map[string]kubeapi.MappingNamespace),
			node:                  make(map[string]kubeapi.MappingNode),
			pod:                   make(map[string]kubeapi.MappingPod),
			container:             make(map[string]kubeapi.MappingContainer),
			service:               make(map[string]kubeapi.MappingService),
			persistentvolumeclaim: make(map[string]kubeapi.MappingPvc),
			persistentvolume:      make(map[string]kubeapi.MappingPv),
			deployment:            make(map[string]kubeapi.MappingDeployment),
			statefulset:           make(map[string]kubeapi.MappingStatefulSet),
			daemonset:             make(map[string]kubeapi.MappingDaemonSet),
			replicaset:            make(map[string]kubeapi.MappingReplicaSet),
			ingress:               make(map[string]kubeapi.MappingIngress),
			ingresshost:           make(map[string]kubeapi.MappingIngressHost),
			storageclass:          make(map[string]kubeapi.MappingStorageClass),
		})

		InitNamespaceInfo(k, v)
		InitNodeInfo(k, v)
		InitPodInfo(k, v)
		InitContainerInfo(k, v)
		InitServiceInfo(k, v)
		InitPersistentVolumeClaimInfo(k, v)
		InitPersistentVolumeInfo(k, v)
		InitDeploymentInfo(k, v)
		InitStatefulSetInfo(k, v)
		InitDaemonSetInfo(k, v)
		InitReplicaSetInfo(k, v)
		InitIngressInfo(k, v)
		InitIngressHostInfo(k, v)
		InitStorageClassInfo(k, v)
	}
}
