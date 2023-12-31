package database

import (
	"context"
	"fmt"
	"onTuneKubeManager/common"
	"onTuneKubeManager/kubeapi"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func UpdateClusterStatusinfo() {
	for {
		status_map := <-common.ChannelClusterStatus

		for k, v := range status_map {
			conn, err := common.DBConnectionPool.Acquire(context.Background())
			if conn == nil || err != nil {
				errorDisconnect(errors.Wrap(err, "Acquire connection error"))
				return
			}

			var flag string
			if v {
				flag = "1"
				common.ChannelRequestChangeHost <- k
			} else {
				flag = "0"
			}
			common.ClusterStatusMap.Store(k, v)

			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			ontunetime, _ := GetOntuneTime()
			if ontunetime == 0 {
				return
			}

			clusterid := common.ClusterID[k]
			_, err = tx.Exec(context.Background(), fmt.Sprintf(UPDATE_STATUS, TB_KUBE_CLUSTER_INFO, "status", flag, ontunetime, clusterid))
			if !errorCheck(err) {
				return
			}

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}

			conn.Release()

			if flag == "0" {
				updateResourceStatusByCluster(ontunetime, k)
				updateTableinfo(TB_KUBE_NODE_INFO, ontunetime)
				updateTableinfo(TB_KUBE_POD_INFO, ontunetime)
			}

			updateTableinfo(TB_KUBE_CLUSTER_INFO, ontunetime)
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func ResourceSender() {
	delResourceinfo()
	for {
		resource_data := <-ChannelResourceInsert
		for key, data := range resource_data {
			if key == "resourceinfo_insert" {
				rsc_data := data.(Resourceinfo)
				insertResourceinfo(rsc_data)
			} else if key == "namespace_update" {
				ns_data := data.(map[string]kubeapi.MappingNamespace)
				ontunetime, _ := GetOntuneTime()
				if ontunetime == 0 {
					return
				}

				var host string
				for _, v := range ns_data {
					host = v.Host
					break
				}

				if ar, ok := mapApiResource.Load(host); ok {
					apiresource := ar.(*ApiResource)
					mapNamespaceInfo := apiresource.namespace
					common.LogManager.Debug(fmt.Sprintf("mapnamespaceinfo : %v", mapNamespaceInfo))
					common.LogManager.Debug(fmt.Sprintf("ns_data : %v", ns_data))

					UpdateList := namespace_updateCheck(ns_data, mapNamespaceInfo)
					updateCnt := updateEnableNamespaceinfo(ontunetime, ns_data)
					if len(UpdateList) > 0 {
						updateNamespaceinfo(ns_data, UpdateList, ontunetime)
					}
					if updateCnt > 0 || len(UpdateList) > 0 {
						updateTableinfo(TB_KUBE_NS_INFO, ontunetime)
					}
				}
			} else if key == "namespace_insert" {
				ns_data := data.(Namespaceinfo)
				insertNamespaceinfo(ns_data)
			} else if key == "node_update" {
				node_data := data.(map[string]kubeapi.MappingNode)
				ontunetime, _ := GetOntuneTime()
				if ontunetime == 0 {
					return
				}

				var host string
				for _, v := range node_data {
					host = v.Host
					break
				}

				if ar, ok := mapApiResource.Load(host); ok {
					apiresource := ar.(*ApiResource)
					mapNodeInfo := apiresource.node

					UpdateList := node_updateCheck(node_data, mapNodeInfo)
					updateCnt := updateEnableNodeinfo(ontunetime, node_data)
					if len(UpdateList) > 0 {
						updateNodeinfo(node_data, UpdateList, ontunetime)
					}
					if updateCnt > 0 || len(UpdateList) > 0 {
						updateTableinfo(TB_KUBE_NODE_INFO, ontunetime)
					}
				}
			} else if key == "node_insert" {
				node_data := data.(Nodeinfo)
				insertNodeinfo(node_data)
			} else if key == "pods_update" {
				pod_data := data.(map[string]kubeapi.MappingPod)
				ontunetime, _ := GetOntuneTime()
				if ontunetime == 0 {
					return
				}

				var host string
				for _, v := range pod_data {
					host = v.Host
					break
				}

				if ar, ok := mapApiResource.Load(host); ok {
					apiresource := ar.(*ApiResource)
					mapPodInfo := apiresource.pod

					UpdateList := pod_updateCheck(pod_data, mapPodInfo)
					updateCnt := updateEnablePodinfo(ontunetime, pod_data)
					if len(UpdateList) > 0 {
						updatePodinfo(pod_data, UpdateList, ontunetime)
					}
					if updateCnt > 0 || len(UpdateList) > 0 {
						updateTableinfo(TB_KUBE_POD_INFO, ontunetime)
					}
				}
			} else if key == "pods_insert" {
				pod_data := data.(Podinfo)
				insertPodinfo(pod_data)
			} else if key == "container_update" {
				container_data := data.(map[string]kubeapi.MappingContainer)
				ontunetime, _ := GetOntuneTime()
				if ontunetime == 0 {
					return
				}

				var host string
				for _, v := range container_data {
					host = v.Host
					break
				}

				if ar, ok := mapApiResource.Load(host); ok {
					apiresource := ar.(*ApiResource)
					mapContainerInfo := apiresource.container

					UpdateList := container_updateCheck(container_data, mapContainerInfo)
					updateCnt := updateEnableContainerinfo(ontunetime, container_data)
					if len(UpdateList) > 0 {
						updateContainerinfo(container_data, UpdateList, ontunetime)
					}
					if updateCnt > 0 || len(UpdateList) > 0 {
						updateTableinfo(TB_KUBE_CONTAINER_INFO, ontunetime)
					}
				}
			} else if key == "container_insert" {
				container_data := data.(Containerinfo)
				insertContainerinfo(container_data)
			} else if key == "service_update" {
				service_data := data.(map[string]kubeapi.MappingService)
				ontunetime, _ := GetOntuneTime()
				if ontunetime == 0 {
					return
				}

				var host string
				for _, v := range service_data {
					host = v.Host
					break
				}

				if ar, ok := mapApiResource.Load(host); ok {
					apiresource := ar.(*ApiResource)
					mapServiceInfo := apiresource.service

					UpdateList := service_updateCheck(service_data, mapServiceInfo)
					updateCnt := updateEnableServiceinfo(ontunetime, service_data)
					if len(UpdateList) > 0 {
						updateServiceinfo(service_data, UpdateList, ontunetime)
					}
					if updateCnt > 0 || len(UpdateList) > 0 {
						updateTableinfo(TB_KUBE_SVC_INFO, ontunetime)
					}
				}
			} else if key == "service_insert" {
				service_data := data.(Serviceinfo)
				insertServiceinfo(service_data)
			} else if key == "pvc_update" {
				pvc_data := data.(map[string]kubeapi.MappingPvc)
				ontunetime, _ := GetOntuneTime()
				if ontunetime == 0 {
					return
				}

				var host string
				for _, v := range pvc_data {
					host = v.Host
					break
				}

				if ar, ok := mapApiResource.Load(host); ok {
					apiresource := ar.(*ApiResource)
					mapPvcInfo := apiresource.persistentvolumeclaim

					UpdateList := pvc_updateCheck(pvc_data, mapPvcInfo)
					updateCnt := updateEnablePvcinfo(ontunetime, pvc_data)
					if len(UpdateList) > 0 {
						updatePvcinfo(pvc_data, UpdateList, ontunetime)
					}
					if updateCnt > 0 || len(UpdateList) > 0 {
						updateTableinfo(TB_KUBE_PVC_INFO, ontunetime)
					}
				}
			} else if key == "pvc_insert" {
				pvc_data := data.(Pvcinfo)
				insertPvcinfo(pvc_data)
			} else if key == "pv_update" {
				pv_data := data.(map[string]kubeapi.MappingPv)
				ontunetime, _ := GetOntuneTime()
				if ontunetime == 0 {
					return
				}

				var host string
				for _, v := range pv_data {
					host = v.Host
					break
				}

				if ar, ok := mapApiResource.Load(host); ok {
					apiresource := ar.(*ApiResource)
					mapPvInfo := apiresource.persistentvolume

					UpdateList := pv_updateCheck(pv_data, mapPvInfo)
					updateCnt := updateEnablePvinfo(ontunetime, pv_data)
					if len(UpdateList) > 0 {
						updatePvinfo(pv_data, UpdateList, ontunetime)
					}
					if updateCnt > 0 || len(UpdateList) > 0 {
						updateTableinfo(TB_KUBE_PV_INFO, ontunetime)
					}
				}
			} else if key == "pv_insert" {
				pv_data := data.(Pvinfo)
				insertPvinfo(pv_data)
			} else if key == "deploy_update" {
				deploy_data := data.(map[string]kubeapi.MappingDeployment)
				ontunetime, _ := GetOntuneTime()
				if ontunetime == 0 {
					return
				}

				var host string
				for _, v := range deploy_data {
					host = v.Host
					break
				}

				if ar, ok := mapApiResource.Load(host); ok {
					apiresource := ar.(*ApiResource)
					mapDeployInfo := apiresource.deployment

					UpdateList := deploy_updateCheck(deploy_data, mapDeployInfo)
					updateCnt := updateEnableDeployinfo(ontunetime, deploy_data)
					if len(UpdateList) > 0 {
						updateDeployinfo(deploy_data, UpdateList, ontunetime)
					}
					if updateCnt > 0 || len(UpdateList) > 0 {
						updateTableinfo(TB_KUBE_DEPLOY_INFO, ontunetime)
					}
				}
			} else if key == "deploy_insert" {
				deploy_data := data.(Deployinfo)
				insertDeployinfo(deploy_data)
			} else if key == "stateful_update" {
				stateful_data := data.(map[string]kubeapi.MappingStatefulSet)
				ontunetime, _ := GetOntuneTime()
				if ontunetime == 0 {
					return
				}

				var host string
				for _, v := range stateful_data {
					host = v.Host
					break
				}

				if ar, ok := mapApiResource.Load(host); ok {
					apiresource := ar.(*ApiResource)
					mapStatefulSetInfo := apiresource.statefulset

					UpdateList := stateful_updateCheck(stateful_data, mapStatefulSetInfo)
					updateCnt := updateEnableStatefulinfo(ontunetime, stateful_data)
					if len(UpdateList) > 0 {
						updateStatefulinfo(stateful_data, UpdateList, ontunetime)
					}
					if updateCnt > 0 || len(UpdateList) > 0 {
						updateTableinfo(TB_KUBE_STS_INFO, ontunetime)
					}
				}
			} else if key == "stateful_insert" {
				stateful_data := data.(StateFulSetinfo)
				insertStatefulinfo(stateful_data)
			} else if key == "daemonset_update" {
				daemonset_data := data.(map[string]kubeapi.MappingDaemonSet)
				ontunetime, _ := GetOntuneTime()
				if ontunetime == 0 {
					return
				}

				var host string
				for _, v := range daemonset_data {
					host = v.Host
					break
				}

				if ar, ok := mapApiResource.Load(host); ok {
					apiresource := ar.(*ApiResource)
					mapDaemonSetInfo := apiresource.daemonset

					UpdateList := daemonset_updateCheck(daemonset_data, mapDaemonSetInfo)
					updateCnt := updateEnableDaemonsetinfo(ontunetime, daemonset_data)
					if len(UpdateList) > 0 {
						updateDaemonsetinfo(daemonset_data, UpdateList, ontunetime)
					}
					if updateCnt > 0 || len(UpdateList) > 0 {
						updateTableinfo(TB_KUBE_DS_INFO, ontunetime)
					}
				}
			} else if key == "daemonset_insert" {
				daemonset_data := data.(DaemonSetinfo)
				insertDaemonsetinfo(daemonset_data)
			} else if key == "replicaset_update" {
				replicaset_data := data.(map[string]kubeapi.MappingReplicaSet)
				ontunetime, _ := GetOntuneTime()
				if ontunetime == 0 {
					return
				}

				var host string
				for _, v := range replicaset_data {
					host = v.Host
					break
				}

				if ar, ok := mapApiResource.Load(host); ok {
					apiresource := ar.(*ApiResource)
					mapReplicaSetInfo := apiresource.replicaset

					UpdateList := replicaset_updateCheck(replicaset_data, mapReplicaSetInfo)
					updateCnt := updateEnableReplicasetinfo(ontunetime, replicaset_data)
					if len(UpdateList) > 0 {
						updateReplicasetinfo(replicaset_data, UpdateList, ontunetime)
					}
					if updateCnt > 0 || len(UpdateList) > 0 {
						updateTableinfo(TB_KUBE_RS_INFO, ontunetime)
					}
				}
			} else if key == "replicaset_insert" {
				replicaset_data := data.(ReplicaSetinfo)
				insertReplicasetinfo(replicaset_data)
			} else if key == "sc_update" {
				sc_data := data.(map[string]kubeapi.MappingStorageClass)
				ontunetime, _ := GetOntuneTime()
				if ontunetime == 0 {
					return
				}

				var host string
				for _, v := range sc_data {
					host = v.Host
					break
				}

				if ar, ok := mapApiResource.Load(host); ok {
					apiresource := ar.(*ApiResource)
					mapStorageClassInfo := apiresource.storageclass

					UpdateList := sc_updateCheck(sc_data, mapStorageClassInfo)
					updateCnt := updateEnableScinfo(ontunetime, sc_data)
					if len(UpdateList) > 0 {
						updateScinfo(sc_data, UpdateList, ontunetime)
					}
					if updateCnt > 0 || len(UpdateList) > 0 {
						updateTableinfo(TB_KUBE_SC_INFO, ontunetime)
					}
				}
			} else if key == "sc_insert" {
				sc_data := data.(Scinfo)
				insertScinfo(sc_data)
			} else if key == "ing_update" {
				ing_data := data.(map[string]kubeapi.MappingIngress)
				ontunetime, _ := GetOntuneTime()
				if ontunetime == 0 {
					return
				}

				var host string
				for _, v := range ing_data {
					host = v.Host
					break
				}

				if ar, ok := mapApiResource.Load(host); ok {
					apiresource := ar.(*ApiResource)
					mapIngressInfo := apiresource.ingress

					UpdateList := ing_updateCheck(ing_data, mapIngressInfo)
					updateCnt := updateEnableInginfo(ontunetime, ing_data)
					if len(UpdateList) > 0 {
						updateInginfo(ing_data, UpdateList, ontunetime)
					}
					if updateCnt > 0 || len(UpdateList) > 0 {
						updateTableinfo(TB_KUBE_ING_INFO, ontunetime)
					}
				}
			} else if key == "ing_insert" {
				ing_data := data.(Inginfo)
				insertInginfo(ing_data)
			} else if key == "inghost_update" {
				inghost_data := data.(map[string]kubeapi.MappingIngressHost)
				ontunetime, _ := GetOntuneTime()
				if ontunetime == 0 {
					return
				}

				var host string
				for _, v := range inghost_data {
					host = v.Host
					break
				}

				if ar, ok := mapApiResource.Load(host); ok {
					apiresource := ar.(*ApiResource)
					mapIngressHostInfo := apiresource.ingresshost

					UpdateList := inghost_updateCheck(inghost_data, mapIngressHostInfo)
					updateCnt := updateEnableInghostinfo(ontunetime, inghost_data)
					if len(UpdateList) > 0 {
						updateInghostinfo(inghost_data, UpdateList, ontunetime)
					}
					if updateCnt > 0 || len(UpdateList) > 0 {
						updateTableinfo(TB_KUBE_INGHOST_INFO, ontunetime)
					}
				}
			} else if key == "inghost_insert" {
				inghost_data := data.(IngHostinfo)
				insertInghostinfo(inghost_data)
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func delResourceinfo() {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), "delete from kuberesourceinfo")
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}
}

func insertResourceinfo(Arr_rcs Resourceinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(Arr_rcs.ArrClusterid); i++ {
		Arr_rcs.ArrCreateTime[i] = ontunetime
		Arr_rcs.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_RESOURCE_INFO, pq.Array(Arr_rcs.ArrClusterid), pq.StringArray(Arr_rcs.ArrResourcename), pq.StringArray(Arr_rcs.ArrApiclass), pq.StringArray(Arr_rcs.ArrVersion), pq.StringArray(Arr_rcs.ArrEndpoint), pq.Array(Arr_rcs.ArrEnabled), pq.Int64Array(Arr_rcs.ArrCreateTime), pq.Int64Array(Arr_rcs.ArrUpdateTime))
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(INSERT_RESOURCE_INFO, ontunetime)
}

func updateResourceStatusByCluster(ontunetime int64, hostip string) {
	clusterid := common.ClusterID[hostip]

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), fmt.Sprintf(UPDATE_STATUS, TB_KUBE_NODE_INFO, "status", "0", ontunetime, clusterid))
	if !errorCheck(err) {
		return
	}

	_, err = tx.Exec(context.Background(), fmt.Sprintf(UPDATE_STATUS, TB_KUBE_POD_INFO, "status", "Cluster Disconnected", ontunetime, clusterid))
	if !errorCheck(err) {
		return
	}

	_, err = tx.Exec(context.Background(), fmt.Sprintf(UPDATE_STATUS, TB_KUBE_CONTAINER_INFO, "state", "cluster disconnected", ontunetime, clusterid))
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapNodeInfo := apiresource.node
		for nk, nv := range mapNodeInfo {
			nv.Status = 0
			mapNodeInfo[nk] = nv
		}

		mapPodInfo := apiresource.pod
		for pk, pv := range mapPodInfo {
			pv.Status = "Cluster Disconnected"
			mapPodInfo[pk] = pv
		}

		mapContainerInfo := apiresource.container
		for ck, cv := range mapContainerInfo {
			cv.State = "cluster disconnected"
			mapContainerInfo[ck] = cv
		}

		apiresource.node = mapNodeInfo
		apiresource.pod = mapPodInfo
		apiresource.container = mapContainerInfo

		mapApiResource.Store(hostip, apiresource)
	}
}

func namespace_updateCheck(new_info map[string]kubeapi.MappingNamespace, old_info map[string]kubeapi.MappingNamespace) []string {
	var UpdateList []string
	if len(new_info) != len(old_info) { // 기존데이터와 새로운 데이터의 갯수가 다르다면 업데이트 필요
		for key, d := range new_info {
			old_data := old_info[key]
			if old_data.UID == "" {
				UpdateList = append(UpdateList, d.UID)
			}
		}

		return UpdateList
	}
	for i, d := range new_info {
		old_infodata := old_info[i]
		if old_infodata.UID == "" || !reflect.DeepEqual(d, old_infodata) {
			UpdateList = append(UpdateList, d.UID)
		}
	}

	return UpdateList
}

func updateEnableNamespaceinfo(ontunetime int64, update_info map[string]kubeapi.MappingNamespace) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	var uidhostinfo map[string]string = make(map[string]string)
	for uid, data := range update_info {
		uidhostinfo[uid] = data.Host
	}

	uid, hostip := getUidHost(uidhostinfo)
	result, err := conn.Query(context.Background(), "update "+TB_KUBE_NS_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where nsuid not in "+uid+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING nsuid")
	if !errorCheck(err) {
		return 0
	}

	var updateUid string
	var returnVal int64
	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapNamespaceInfo := apiresource.namespace
		for result.Next() {
			err := result.Scan(&updateUid)
			if !errorCheck(err) {
				return 0
			}
			delete(mapNamespaceInfo, updateUid)
			returnVal++
		}
		apiresource.namespace = mapNamespaceInfo
		mapApiResource.Store(hostip, apiresource)
	}

	result.Close()

	return returnVal
}

func insertNamespaceinfo(ArrResource Namespaceinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(ArrResource.ArrClusterid); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_NAMESPACE_INFO, ArrResource.GetArgs()...)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	common.LogManager.WriteLog(fmt.Sprintf("Manager DB - namespace insertion is completed: %s", strings.Join(ArrResource.ArrNsUid, ",")))

	updateTableinfo(TB_KUBE_NS_INFO, ontunetime)
}

func updateNamespaceinfo(update_info map[string]kubeapi.MappingNamespace, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	hostname := update_info[update_list[0]].Host

	if ar, ok := mapApiResource.Load(hostname); ok {
		apiresource := ar.(*ApiResource)
		mapNamespaceInfo := apiresource.namespace

		for _, key := range update_list {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			result, err := tx.Exec(context.Background(), UPDATE_NAMESPACE_INFO, common.ClusterID[update_info[key].Host], update_info[key].Name, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].Status, ontunetime, update_info[key].UID)
			if !errorCheck(err) {
				return
			}
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_NAMESPACE_INFO, update_info[key].UID, common.ClusterID[update_info[key].Host], update_info[key].Name, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].Status, 1, ontunetime, ontunetime)
				if !errorCheck(insert_err) {
					return
				}
			}
			var update_data kubeapi.MappingNamespace
			update_data.UID = update_info[key].UID
			update_data.Name = update_info[key].Name
			update_data.StartTime = update_info[key].StartTime
			update_data.Labels = update_info[key].Labels
			update_data.Status = update_info[key].Status
			mapNamespaceInfo[update_info[key].UID] = update_data

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		common.LogManager.WriteLog(fmt.Sprintf("Manager DB - %s namespace update is completed", hostname))

		apiresource.namespace = mapNamespaceInfo
		mapApiResource.Store(hostname, apiresource)
	}
}

func node_updateCheck(new_info map[string]kubeapi.MappingNode, old_info map[string]kubeapi.MappingNode) []string {
	var UpdateList []string
	if len(new_info) != len(old_info) { // 기존데이터와 새로운 데이터의 갯수가 다르다면 업데이트 필요
		for key, d := range new_info {
			old_data := old_info[key]
			if old_data.UID == "" {
				UpdateList = append(UpdateList, d.UID)
			}
		}

		return UpdateList
	}
	for i, d := range new_info {
		old_infodata := old_info[i]
		if old_infodata.UID == "" || !reflect.DeepEqual(d, old_infodata) {
			UpdateList = append(UpdateList, d.UID)
		}
	}

	return UpdateList
}

func updateEnableNodeinfo(ontunetime int64, update_info map[string]kubeapi.MappingNode) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	var uidhostinfo map[string]string = make(map[string]string)
	for uid, data := range update_info {
		uidhostinfo[uid] = data.Host
	}

	uid, hostip := getUidHost(uidhostinfo)
	result, err := conn.Query(context.Background(), "update "+TB_KUBE_NODE_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where nodeuid not in "+uid+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING nodeuid")
	if !errorCheck(err) {
		return 0
	}

	var updateUid string
	var returnVal int64
	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapNodeInfo := apiresource.node
		for result.Next() {
			err := result.Scan(&updateUid)
			if !errorCheck(err) {
				return 0
			}
			delete(mapNodeInfo, updateUid)
			returnVal++
		}
		apiresource.node = mapNodeInfo
		mapApiResource.Store(hostip, apiresource)
	}

	result.Close()

	return returnVal
}

func updateNodeinfo(update_info map[string]kubeapi.MappingNode, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	hostname := update_info[update_list[0]].Host

	if ar, ok := mapApiResource.Load(hostname); ok {
		apiresource := ar.(*ApiResource)
		mapNodeInfo := apiresource.node

		for _, key := range update_list {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			result, err := tx.Exec(context.Background(), UPDATE_NODE_INFO, common.ManagerID, common.ClusterID[update_info[key].Host], update_info[key].Name, update_info[key].NodeType, getStarttime(update_info[key].StartTime.Unix(), biastime),
				update_info[key].Labels, update_info[key].KernelVersion, update_info[key].OSImage, update_info[key].OSName, update_info[key].ContainerRuntimeVersion, update_info[key].KubeletVersion,
				update_info[key].KubeProxyVersion, update_info[key].CPUArch, update_info[key].CPUCount, update_info[key].EphemeralStorage, update_info[key].MemorySize,
				update_info[key].Pods, update_info[key].IP, update_info[key].Status, ontunetime, update_info[key].UID)
			if !errorCheck(err) {
				return
			}
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_NODE_INFO, common.ManagerID, common.ClusterID[update_info[key].Host], update_info[key].UID, update_info[key].Name, update_info[key].Name,
					update_info[key].NodeType, 1, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels,
					update_info[key].KernelVersion, update_info[key].OSImage, update_info[key].OSName, update_info[key].ContainerRuntimeVersion, update_info[key].KubeletVersion,
					update_info[key].KubeProxyVersion, update_info[key].CPUArch, update_info[key].CPUCount, update_info[key].EphemeralStorage, update_info[key].MemorySize,
					update_info[key].Pods, update_info[key].IP, update_info[key].Status, ontunetime, ontunetime)
				if !errorCheck(insert_err) {
					return
				}
			}
			var update_data kubeapi.MappingNode
			update_data.UID = update_info[key].UID
			update_data.Name = update_info[key].Name
			update_data.NodeType = update_info[key].NodeType
			update_data.StartTime = update_info[key].StartTime
			update_data.Host = update_info[key].Host
			update_data.Labels = update_info[key].Labels
			update_data.KernelVersion = update_info[key].KernelVersion
			update_data.OSImage = update_info[key].OSImage
			update_data.OSName = update_info[key].OSName
			update_data.ContainerRuntimeVersion = update_info[key].ContainerRuntimeVersion
			update_data.KubeletVersion = update_info[key].KubeletVersion
			update_data.KubeProxyVersion = update_info[key].KubeProxyVersion
			update_data.CPUArch = update_info[key].CPUArch
			update_data.CPUCount = int64(update_info[key].CPUCount)
			update_data.EphemeralStorage = update_info[key].EphemeralStorage
			update_data.MemorySize = update_info[key].MemorySize
			update_data.Pods = update_info[key].Pods
			update_data.IP = update_info[key].IP
			update_data.Status = update_info[key].Status
			mapNodeInfo[update_info[key].UID] = update_data

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		common.LogManager.WriteLog(fmt.Sprintf("Manager DB - %s node update is completed", hostname))

		apiresource.node = mapNodeInfo
		mapApiResource.Store(hostname, apiresource)
	}
}

func insertNodeinfo(ArrResource Nodeinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(ArrResource.ArrManagerid); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_NODE_INFO, ArrResource.GetArgs()...)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	common.LogManager.WriteLog(fmt.Sprintf("Manager DB - node insertion is completed: %s", strings.Join(ArrResource.ArrNodeUid, ",")))

	updateTableinfo(TB_KUBE_NODE_INFO, ontunetime)
}

func pod_updateCheck(new_info map[string]kubeapi.MappingPod, old_info map[string]kubeapi.MappingPod) []string {
	var UpdateList []string
	if len(new_info) != len(old_info) { // 기존데이터와 새로운 데이터의 갯수가 다르다면 업데이트 필요
		for key, d := range new_info {
			old_data := old_info[key]
			if old_data.UID == "" {
				UpdateList = append(UpdateList, d.UID)
			}
		}

		return UpdateList
	}
	for i, d := range new_info {
		old_infodata := old_info[i]
		if old_infodata.Name == "" || !reflect.DeepEqual(d, old_infodata) {
			UpdateList = append(UpdateList, d.UID)
		}
	}

	return UpdateList
}

func updateEnablePodinfo(ontunetime int64, update_info map[string]kubeapi.MappingPod) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	var uidhostinfo map[string]string = make(map[string]string)
	for uid, data := range update_info {
		uidhostinfo[uid] = data.Host
	}

	uid, hostip := getUidHost(uidhostinfo)
	result, err := conn.Query(context.Background(), "update "+TB_KUBE_POD_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where enabled = 1 and uid not in "+uid+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING uid")
	if !errorCheck(err) {
		return 0
	}

	var updateUid string
	var returnVal int64
	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapPodInfo := apiresource.pod
		for result.Next() {
			err := result.Scan(&updateUid)
			if !errorCheck(err) {
				return 0
			}
			delete(mapPodInfo, updateUid)
			returnVal++
		}
		apiresource.pod = mapPodInfo
		mapApiResource.Store(hostip, apiresource)
	}

	result.Close()

	return returnVal
}

func updatePodinfo(update_info map[string]kubeapi.MappingPod, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	hostname := update_info[update_list[0]].Host

	if ar, ok := mapApiResource.Load(hostname); ok {
		apiresource := ar.(*ApiResource)
		mapPodInfo := apiresource.pod

		for _, key := range update_list {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			result, err := tx.Exec(context.Background(), UPDATE_POD_INFO, getUID(METRIC_VAR_NODE, update_info[key].Host, update_info[key].NodeName), getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName),
				update_info[key].AnnotationUID, update_info[key].Name, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].Selector,
				update_info[key].RestartPolicy, update_info[key].ServiceAccount, update_info[key].Status, update_info[key].HostIP, update_info[key].PodIP,
				int64(update_info[key].RestartCount), getStarttime(update_info[key].RestartTime.Unix(), biastime), update_info[key].Condition, update_info[key].StaticPod,
				update_info[key].ReferenceKind, update_info[key].ReferenceUID, getUID("persistentvolumeclaim", update_info[key].NamespaceName, update_info[key].PersistentVolumeClaim),
				ontunetime, update_info[key].UID, common.ClusterID[update_info[key].Host])
			if !errorCheck(err) {
				return
			}
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_POD_INFO, update_info[key].UID, getUID(METRIC_VAR_NODE, update_info[key].Host, update_info[key].NodeName), getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName),
					update_info[key].AnnotationUID, update_info[key].Name, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].Selector,
					update_info[key].RestartPolicy, update_info[key].ServiceAccount,
					update_info[key].Status, update_info[key].HostIP, update_info[key].PodIP, int64(update_info[key].RestartCount), getStarttime(update_info[key].RestartTime.Unix(), biastime),
					update_info[key].Condition, update_info[key].StaticPod, update_info[key].ReferenceKind, update_info[key].ReferenceUID, getUID("persistentvolumeclaim", update_info[key].NamespaceName, update_info[key].PersistentVolumeClaim),
					1, ontunetime, ontunetime, common.ClusterID[update_info[key].Host])
				if !errorCheck(insert_err) {
					return
				}
			}
			var update_data kubeapi.MappingPod
			update_data.UID = update_info[key].UID
			update_data.Name = update_info[key].Name
			update_data.NodeName = update_info[key].NodeName
			update_data.NamespaceName = update_info[key].NamespaceName
			update_data.AnnotationUID = update_info[key].AnnotationUID
			update_data.StartTime = update_info[key].StartTime
			update_data.Host = update_info[key].Host
			update_data.Labels = update_info[key].Labels
			update_data.Selector = update_info[key].Selector
			update_data.RestartPolicy = update_info[key].RestartPolicy
			update_data.ServiceAccount = update_info[key].ServiceAccount
			update_data.Status = update_info[key].Status
			update_data.HostIP = update_info[key].HostIP
			update_data.PodIP = update_info[key].PodIP
			update_data.RestartCount = int32(update_info[key].RestartCount)
			update_data.RestartTime = update_info[key].RestartTime
			update_data.Condition = update_info[key].Condition
			update_data.StaticPod = update_info[key].StaticPod
			update_data.ReferenceKind = update_info[key].ReferenceKind
			update_data.ReferenceUID = update_info[key].ReferenceUID
			update_data.PersistentVolumeClaim = update_info[key].PersistentVolumeClaim
			update_data.Host = update_info[key].Host
			mapPodInfo[update_info[key].UID] = update_data

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		common.LogManager.WriteLog(fmt.Sprintf("Manager DB - %s pod update is completed", hostname))

		apiresource.pod = mapPodInfo
		mapApiResource.Store(hostname, apiresource)
	}
}

func insertPodinfo(ArrResource Podinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(ArrResource.ArrUid); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_POD_INFO, ArrResource.GetArgs()...)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	common.LogManager.WriteLog(fmt.Sprintf("Manager DB - pod insertion is completed: %s", strings.Join(ArrResource.ArrUid, ",")))

	updateTableinfo(TB_KUBE_POD_INFO, ontunetime)
}

func container_updateCheck(new_info map[string]kubeapi.MappingContainer, old_info map[string]kubeapi.MappingContainer) []string {
	var UpdateList []string
	if len(new_info) != len(old_info) { // 기존데이터와 새로운 데이터의 갯수가 다르다면 업데이트 필요
		for key, d := range new_info {
			old_data := old_info[key]
			container_key := d.UID + ":" + d.Name
			if old_data.Name == "" || old_data.UID == "" {
				UpdateList = append(UpdateList, container_key)
			}
		}

		return UpdateList
	}
	for i, d := range new_info {
		old_infodata := old_info[i]
		container_key := old_infodata.UID + ":" + old_infodata.Name
		if old_infodata.Name == "" || !reflect.DeepEqual(d, old_infodata) {
			UpdateList = append(UpdateList, container_key)
		}
	}

	return UpdateList
}

func updateEnableContainerinfo(ontunetime int64, update_info map[string]kubeapi.MappingContainer) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	var uidhostinfo map[string]string = make(map[string]string)
	for uid, data := range update_info {
		uidhostinfo[uid] = data.Host
	}

	container_key, hostip := getUidHost(uidhostinfo)
	result, err := conn.Query(context.Background(), "update "+TB_KUBE_CONTAINER_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where poduid||':'||containername not in "+container_key+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING containername, poduid")
	if !errorCheck(err) {
		return 0
	}

	var updateName string
	var poduid string
	var returnVal int64
	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapContainerInfo := apiresource.container
		for result.Next() {
			err := result.Scan(&updateName, &poduid)
			if !errorCheck(err) {
				return 0
			}

			key := poduid + ":" + updateName
			delete(mapContainerInfo, key)
			returnVal++
		}
		apiresource.container = mapContainerInfo
		mapApiResource.Store(hostip, apiresource)
	}

	result.Close()

	return returnVal
}

func updateContainerinfo(update_info map[string]kubeapi.MappingContainer, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	hostname := update_info[update_list[0]].Host

	if ar, ok := mapApiResource.Load(hostname); ok {
		apiresource := ar.(*ApiResource)
		mapContainerInfo := apiresource.container

		for _, key := range update_list {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			result, err := tx.Exec(context.Background(), UPDATE_CONTAINER_INFO, update_info[key].UID, update_info[key].Image, update_info[key].Ports, update_info[key].Env, update_info[key].LimitCpu,
				update_info[key].LimitMemory, update_info[key].LimitStorage, update_info[key].LimitEphemeral, update_info[key].RequestCpu, update_info[key].RequestMemory, update_info[key].RequestStorage,
				update_info[key].RequestEphemeral, update_info[key].VolumeMounts, update_info[key].State, ontunetime, update_info[key].Name, common.ClusterID[update_info[key].Host])
			if !errorCheck(err) {
				return
			}
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_CONTAINER_INFO, update_info[key].UID, update_info[key].Name, update_info[key].Image, update_info[key].Ports, update_info[key].Env, update_info[key].LimitCpu,
					update_info[key].LimitMemory, update_info[key].LimitStorage, update_info[key].LimitEphemeral, update_info[key].RequestCpu, update_info[key].RequestMemory, update_info[key].RequestStorage,
					update_info[key].RequestEphemeral, update_info[key].VolumeMounts, update_info[key].State, 1, ontunetime, ontunetime, common.ClusterID[update_info[key].Host])
				if !errorCheck(insert_err) {
					return
				}
			}
			var update_data kubeapi.MappingContainer
			update_data.UID = update_info[key].UID
			update_data.Name = update_info[key].Name
			update_data.Host = update_info[key].Host
			update_data.Image = update_info[key].Image
			update_data.Ports = update_info[key].Ports
			update_data.Env = update_info[key].Env
			update_data.LimitCpu = update_info[key].LimitCpu
			update_data.LimitMemory = update_info[key].LimitMemory
			update_data.LimitStorage = update_info[key].LimitStorage
			update_data.LimitEphemeral = update_info[key].LimitEphemeral
			update_data.RequestCpu = update_info[key].RequestCpu
			update_data.RequestMemory = update_info[key].RequestMemory
			update_data.RequestStorage = update_info[key].RequestStorage
			update_data.RequestEphemeral = update_info[key].RequestEphemeral
			update_data.VolumeMounts = update_info[key].VolumeMounts
			update_data.State = update_info[key].State
			update_data.Host = update_info[key].Host

			container_key := update_info[key].UID + ":" + update_info[key].Name
			mapContainerInfo[container_key] = update_data

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		common.LogManager.WriteLog(fmt.Sprintf("Manager DB - %s container update is completed", hostname))

		apiresource.container = mapContainerInfo
		mapApiResource.Store(hostname, apiresource)
	}
}

func insertContainerinfo(ArrResource Containerinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(ArrResource.ArrContainername); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_CONTAINER_INFO, ArrResource.GetArgs()...)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	common.LogManager.WriteLog(fmt.Sprintf("Manager DB - container insertion is completed: %s", strings.Join(ArrResource.ArrContainername, ",")))

	updateTableinfo(TB_KUBE_CONTAINER_INFO, ontunetime)
}

func service_updateCheck(new_info map[string]kubeapi.MappingService, old_info map[string]kubeapi.MappingService) []string {
	var UpdateList []string
	if len(new_info) != len(old_info) { // 기존데이터와 새로운 데이터의 갯수가 다르다면 업데이트 필요
		for key, d := range new_info {
			old_data := old_info[key]
			if old_data.UID == "" {
				UpdateList = append(UpdateList, d.UID)
			}
		}

		return UpdateList
	}
	for i, d := range new_info {
		old_infodata := old_info[i]
		if old_infodata.UID == "" || !reflect.DeepEqual(d, old_infodata) {
			UpdateList = append(UpdateList, d.UID)
		}
	}

	return UpdateList
}

func updateEnableServiceinfo(ontunetime int64, update_info map[string]kubeapi.MappingService) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	var uidhostinfo map[string]string = make(map[string]string)
	for uid, data := range update_info {
		uidhostinfo[uid] = data.Host
	}

	uid, hostip := getUidHost(uidhostinfo)
	result, err := conn.Query(context.Background(), "update "+TB_KUBE_SVC_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where uid not in "+uid+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING uid")
	if !errorCheck(err) {
		return 0
	}

	var updateUid string
	var returnVal int64
	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapServiceInfo := apiresource.service
		for result.Next() {
			err := result.Scan(&updateUid)
			if !errorCheck(err) {
				return 0
			}
			delete(mapServiceInfo, updateUid)
			returnVal++
		}
		apiresource.service = mapServiceInfo
		mapApiResource.Store(hostip, apiresource)
	}

	result.Close()

	return returnVal
}

func updateServiceinfo(update_info map[string]kubeapi.MappingService, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	hostname := update_info[update_list[0]].Host

	if ar, ok := mapApiResource.Load(hostname); ok {
		apiresource := ar.(*ApiResource)
		mapServiceInfo := apiresource.service

		for _, key := range update_list {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			result, err := tx.Exec(context.Background(), UPDATE_SVC_INFO, getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, getStarttime(update_info[key].StartTime.Unix(), biastime),
				update_info[key].Labels, update_info[key].Selector, update_info[key].ServiceType, update_info[key].ClusterIP, update_info[key].Ports, ontunetime,
				update_info[key].UID, common.ClusterID[update_info[key].Host])
			if !errorCheck(err) {
				return
			}
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_SVC_INFO, getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, update_info[key].UID,
					getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].Selector, update_info[key].ServiceType, update_info[key].ClusterIP,
					update_info[key].Ports, 1, ontunetime, ontunetime, common.ClusterID[update_info[key].Host])
				if !errorCheck(insert_err) {
					return
				}
			}
			var update_data kubeapi.MappingService
			update_data.UID = update_info[key].UID
			update_data.Name = update_info[key].Name
			update_data.NamespaceName = update_info[key].NamespaceName
			update_data.StartTime = update_info[key].StartTime
			update_data.Host = update_info[key].Host
			update_data.Labels = update_info[key].Labels
			update_data.Selector = update_info[key].Selector
			update_data.ServiceType = update_info[key].ServiceType
			update_data.ClusterIP = update_info[key].ClusterIP
			update_data.Ports = update_info[key].Ports
			update_data.Host = update_info[key].Host
			mapServiceInfo[update_info[key].UID] = update_data

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		common.LogManager.WriteLog(fmt.Sprintf("Manager DB - %s service update is completed", hostname))

		apiresource.service = mapServiceInfo
		mapApiResource.Store(hostname, apiresource)
	}
}

func insertServiceinfo(ArrResource Serviceinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(ArrResource.ArrSvcname); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_SVC_INFO, ArrResource.GetArgs()...)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	common.LogManager.WriteLog(fmt.Sprintf("Manager DB - service insertion is completed: %s", strings.Join(ArrResource.ArrUid, ",")))

	updateTableinfo(TB_KUBE_SVC_INFO, ontunetime)
}

func deploy_updateCheck(new_info map[string]kubeapi.MappingDeployment, old_info map[string]kubeapi.MappingDeployment) []string {
	var UpdateList []string
	if len(new_info) != len(old_info) { // 기존데이터와 새로운 데이터의 갯수가 다르다면 업데이트 필요
		for key, d := range new_info {
			old_data := old_info[key]
			if old_data.UID == "" {
				UpdateList = append(UpdateList, d.UID)
			}
		}

		return UpdateList
	}
	for i, d := range new_info {
		old_infodata := old_info[i]
		if old_infodata.UID == "" || !reflect.DeepEqual(d, old_infodata) {
			UpdateList = append(UpdateList, d.UID)
		}
	}

	return UpdateList
}

func updateEnableDeployinfo(ontunetime int64, update_info map[string]kubeapi.MappingDeployment) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	var uidhostinfo map[string]string = make(map[string]string)
	for uid, data := range update_info {
		uidhostinfo[uid] = data.Host
	}

	uid, hostip := getUidHost(uidhostinfo)
	result, err := conn.Query(context.Background(), "update "+TB_KUBE_DEPLOY_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where uid not in "+uid+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING uid")
	if !errorCheck(err) {
		return 0
	}

	var updateUid string
	var returnVal int64
	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapDeployInfo := apiresource.deployment
		for result.Next() {
			err := result.Scan(&updateUid)
			if !errorCheck(err) {
				return 0
			}
			delete(mapDeployInfo, updateUid)
			returnVal++
		}
		apiresource.deployment = mapDeployInfo
		mapApiResource.Store(hostip, apiresource)
	}

	result.Close()

	return returnVal
}

func updateDeployinfo(update_info map[string]kubeapi.MappingDeployment, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	hostname := update_info[update_list[0]].Host

	if ar, ok := mapApiResource.Load(hostname); ok {
		apiresource := ar.(*ApiResource)
		mapDeployInfo := apiresource.deployment

		for _, key := range update_list {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			result, err := tx.Exec(context.Background(), UPDATE_DEPLOY_INFO, getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].Selector, update_info[key].ServiceAccount,
				update_info[key].Replicas, update_info[key].UpdatedReplicas, update_info[key].ReadyReplicas, update_info[key].AvailableReplicas, update_info[key].ObservedGeneneration,
				ontunetime, update_info[key].UID, common.ClusterID[update_info[key].Host])
			if !errorCheck(err) {
				return
			}
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_DEPLOY_INFO, getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, update_info[key].UID, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].Selector,
					update_info[key].ServiceAccount, update_info[key].Replicas, update_info[key].UpdatedReplicas, update_info[key].ReadyReplicas, update_info[key].AvailableReplicas,
					update_info[key].ObservedGeneneration, 1, ontunetime, ontunetime, common.ClusterID[update_info[key].Host])
				if !errorCheck(insert_err) {
					return
				}
			}
			var update_data kubeapi.MappingDeployment
			update_data.UID = update_info[key].UID
			update_data.Name = update_info[key].Name
			update_data.NamespaceName = update_info[key].NamespaceName
			update_data.StartTime = update_info[key].StartTime
			update_data.Host = update_info[key].Host
			update_data.Labels = update_info[key].Labels
			update_data.Selector = update_info[key].Selector
			update_data.ServiceAccount = update_info[key].ServiceAccount
			update_data.Replicas = update_info[key].Replicas
			update_data.UpdatedReplicas = update_info[key].UpdatedReplicas
			update_data.ReadyReplicas = update_info[key].ReadyReplicas
			update_data.AvailableReplicas = update_info[key].AvailableReplicas
			update_data.ObservedGeneneration = update_info[key].ObservedGeneneration
			update_data.Host = update_info[key].Host
			mapDeployInfo[update_info[key].UID] = update_data

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		common.LogManager.WriteLog(fmt.Sprintf("Manager DB - %s deployment update is completed", hostname))

		apiresource.deployment = mapDeployInfo
		mapApiResource.Store(hostname, apiresource)
	}
}

func insertDeployinfo(ArrResource Deployinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(ArrResource.ArrDeployname); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_DEPLOY_INFO, ArrResource.GetArgs()...)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	common.LogManager.WriteLog(fmt.Sprintf("Manager DB - deployment insertion is completed: %s", strings.Join(ArrResource.ArrUid, ",")))

	updateTableinfo(TB_KUBE_DEPLOY_INFO, ontunetime)
}

func stateful_updateCheck(new_info map[string]kubeapi.MappingStatefulSet, old_info map[string]kubeapi.MappingStatefulSet) []string {
	var UpdateList []string
	if len(new_info) != len(old_info) { // 기존데이터와 새로운 데이터의 갯수가 다르다면 업데이트 필요
		for key, d := range new_info {
			old_data := old_info[key]
			if old_data.UID == "" {
				UpdateList = append(UpdateList, d.UID)
			}
		}

		return UpdateList
	}
	for i, d := range new_info {
		old_infodata := old_info[i]
		if old_infodata.UID == "" || !reflect.DeepEqual(d, old_infodata) {
			UpdateList = append(UpdateList, d.UID)
		}
	}

	return UpdateList
}

func updateEnableStatefulinfo(ontunetime int64, update_info map[string]kubeapi.MappingStatefulSet) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	var uidhostinfo map[string]string = make(map[string]string)
	for uid, data := range update_info {
		uidhostinfo[uid] = data.Host
	}

	uid, hostip := getUidHost(uidhostinfo)
	result, err := conn.Query(context.Background(), "update "+TB_KUBE_STS_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where uid not in "+uid+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING uid")
	if !errorCheck(err) {
		return 0
	}

	var updateUid string
	var returnVal int64
	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapStatefulSetInfo := apiresource.statefulset
		for result.Next() {
			err := result.Scan(&updateUid)
			if !errorCheck(err) {
				return 0
			}
			delete(mapStatefulSetInfo, updateUid)
			returnVal++
		}
		apiresource.statefulset = mapStatefulSetInfo
		mapApiResource.Store(hostip, apiresource)
	}

	result.Close()

	return returnVal
}

func updateStatefulinfo(update_info map[string]kubeapi.MappingStatefulSet, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	hostname := update_info[update_list[0]].Host

	if ar, ok := mapApiResource.Load(hostname); ok {
		apiresource := ar.(*ApiResource)
		mapStatefulSetInfo := apiresource.statefulset

		for _, key := range update_list {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			result, err := tx.Exec(context.Background(), UPDATE_STATEFUL_INFO, getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].Selector, update_info[key].ServiceAccount,
				update_info[key].Replicas, update_info[key].UpdatedReplicas, update_info[key].ReadyReplicas, update_info[key].AvailableReplicas, ontunetime, update_info[key].UID,
				common.ClusterID[update_info[key].Host])
			if !errorCheck(err) {
				return
			}
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_STATEFUL_INFO, getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, update_info[key].UID, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].Selector,
					update_info[key].ServiceAccount, update_info[key].Replicas, update_info[key].UpdatedReplicas, update_info[key].ReadyReplicas,
					update_info[key].AvailableReplicas, 1, ontunetime, ontunetime, common.ClusterID[update_info[key].Host])
				if !errorCheck(insert_err) {
					return
				}
			}
			var update_data kubeapi.MappingStatefulSet
			update_data.UID = update_info[key].UID
			update_data.Name = update_info[key].Name
			update_data.NamespaceName = update_info[key].NamespaceName
			update_data.StartTime = update_info[key].StartTime
			update_data.Host = update_info[key].Host
			update_data.Labels = update_info[key].Labels
			update_data.Selector = update_info[key].Selector
			update_data.ServiceAccount = update_info[key].ServiceAccount
			update_data.Replicas = update_info[key].Replicas
			update_data.UpdatedReplicas = update_info[key].UpdatedReplicas
			update_data.ReadyReplicas = update_info[key].ReadyReplicas
			update_data.AvailableReplicas = update_info[key].AvailableReplicas
			update_data.Host = update_info[key].Host
			mapStatefulSetInfo[update_info[key].UID] = update_data

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		common.LogManager.WriteLog(fmt.Sprintf("Manager DB - %s statefulset update is completed", hostname))

		apiresource.statefulset = mapStatefulSetInfo
		mapApiResource.Store(hostname, apiresource)
	}
}

func insertStatefulinfo(ArrResource StateFulSetinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(ArrResource.ArrStsname); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_STATEFUL_INFO, ArrResource.GetArgs()...)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	common.LogManager.WriteLog(fmt.Sprintf("Manager DB - statefulset insertion is completed: %s", strings.Join(ArrResource.ArrUid, ",")))

	updateTableinfo(TB_KUBE_STS_INFO, ontunetime)
}

func daemonset_updateCheck(new_info map[string]kubeapi.MappingDaemonSet, old_info map[string]kubeapi.MappingDaemonSet) []string {
	var UpdateList []string
	if len(new_info) != len(old_info) { // 기존데이터와 새로운 데이터의 갯수가 다르다면 업데이트 필요
		for key, d := range new_info {
			old_data := old_info[key]
			if old_data.UID == "" {
				UpdateList = append(UpdateList, d.UID)
			}
		}

		return UpdateList
	}
	for i, d := range new_info {
		old_infodata := old_info[i]
		if old_infodata.UID == "" || !reflect.DeepEqual(d, old_infodata) {
			UpdateList = append(UpdateList, d.UID)
		}
	}

	return UpdateList
}

func updateEnableDaemonsetinfo(ontunetime int64, update_info map[string]kubeapi.MappingDaemonSet) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	var uidhostinfo map[string]string = make(map[string]string)
	for uid, data := range update_info {
		uidhostinfo[uid] = data.Host
	}

	uid, hostip := getUidHost(uidhostinfo)
	result, err := conn.Query(context.Background(), "update "+TB_KUBE_DS_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where uid not in "+uid+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING uid")
	if !errorCheck(err) {
		return 0
	}

	var updateUid string
	var returnVal int64
	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapDaemonSetInfo := apiresource.daemonset
		for result.Next() {
			err := result.Scan(&updateUid)
			if !errorCheck(err) {
				return 0
			}
			delete(mapDaemonSetInfo, updateUid)
			returnVal++
		}
		apiresource.daemonset = mapDaemonSetInfo
		mapApiResource.Store(hostip, apiresource)
	}

	result.Close()

	return returnVal
}

func updateDaemonsetinfo(update_info map[string]kubeapi.MappingDaemonSet, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	hostname := update_info[update_list[0]].Host

	if ar, ok := mapApiResource.Load(hostname); ok {
		apiresource := ar.(*ApiResource)
		mapDaemonSetInfo := apiresource.daemonset

		for _, key := range update_list {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			result, err := tx.Exec(context.Background(), UPDATE_DAEMONSET_INFO, getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].Selector, update_info[key].ServiceAccount,
				update_info[key].CurrentNumberScheduled, update_info[key].DesiredNumberScheduled, update_info[key].NumberReady, update_info[key].UpdatedNumberScheduled,
				update_info[key].NumberAvailable, ontunetime, update_info[key].UID, common.ClusterID[update_info[key].Host])
			if !errorCheck(err) {
				return
			}
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_DAEMONSET_INFO, getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, update_info[key].UID, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].Selector,
					update_info[key].ServiceAccount, update_info[key].CurrentNumberScheduled, update_info[key].DesiredNumberScheduled, update_info[key].NumberReady, update_info[key].UpdatedNumberScheduled,
					update_info[key].NumberAvailable, 1, ontunetime, ontunetime, common.ClusterID[update_info[key].Host])
				if !errorCheck(insert_err) {
					return
				}
			}
			var update_data kubeapi.MappingDaemonSet
			update_data.UID = update_info[key].UID
			update_data.Name = update_info[key].Name
			update_data.NamespaceName = update_info[key].NamespaceName
			update_data.StartTime = update_info[key].StartTime
			update_data.Host = update_info[key].Host
			update_data.Labels = update_info[key].Labels
			update_data.Selector = update_info[key].Selector
			update_data.ServiceAccount = update_info[key].ServiceAccount
			update_data.CurrentNumberScheduled = update_info[key].CurrentNumberScheduled
			update_data.DesiredNumberScheduled = update_info[key].DesiredNumberScheduled
			update_data.NumberReady = update_info[key].NumberReady
			update_data.UpdatedNumberScheduled = update_info[key].UpdatedNumberScheduled
			update_data.NumberAvailable = update_info[key].NumberAvailable
			update_data.Host = update_info[key].Host
			mapDaemonSetInfo[update_info[key].UID] = update_data

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		common.LogManager.WriteLog(fmt.Sprintf("Manager DB - %s daemonset update is completed", hostname))

		apiresource.daemonset = mapDaemonSetInfo
		mapApiResource.Store(hostname, apiresource)
	}
}

func insertDaemonsetinfo(ArrResource DaemonSetinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(ArrResource.ArrDsname); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_DAEMONSET_INFO, ArrResource.GetArgs()...)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	common.LogManager.WriteLog(fmt.Sprintf("Manager DB - daemonset insertion is completed: %s", strings.Join(ArrResource.ArrUid, ",")))

	updateTableinfo(TB_KUBE_DS_INFO, ontunetime)
}

func replicaset_updateCheck(new_info map[string]kubeapi.MappingReplicaSet, old_info map[string]kubeapi.MappingReplicaSet) []string {
	var UpdateList []string
	if len(new_info) != len(old_info) { // 기존데이터와 새로운 데이터의 갯수가 다르다면 업데이트 필요
		for key, d := range new_info {
			old_data := old_info[key]
			if old_data.UID == "" {
				UpdateList = append(UpdateList, d.UID)
			}
		}

		return UpdateList
	}
	for i, d := range new_info {
		old_infodata := old_info[i]
		if old_infodata.UID == "" || !reflect.DeepEqual(d, old_infodata) {
			UpdateList = append(UpdateList, d.UID)
		}
	}

	return UpdateList
}

func updateEnableReplicasetinfo(ontunetime int64, update_info map[string]kubeapi.MappingReplicaSet) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	var uidhostinfo map[string]string = make(map[string]string)
	for uid, data := range update_info {
		uidhostinfo[uid] = data.Host
	}

	uid, hostip := getUidHost(uidhostinfo)
	result, err := conn.Query(context.Background(), "update "+TB_KUBE_RS_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where uid not in "+uid+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING uid")
	if !errorCheck(err) {
		return 0
	}

	var updateUid string
	var returnVal int64
	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapReplicaSetInfo := apiresource.replicaset
		for result.Next() {
			err := result.Scan(&updateUid)
			if !errorCheck(err) {
				return 0
			}
			delete(mapReplicaSetInfo, updateUid)
			returnVal++
		}
		apiresource.replicaset = mapReplicaSetInfo
		mapApiResource.Store(hostip, apiresource)
	}

	result.Close()

	return returnVal
}

func updateReplicasetinfo(update_info map[string]kubeapi.MappingReplicaSet, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	hostname := update_info[update_list[0]].Host

	if ar, ok := mapApiResource.Load(hostname); ok {
		apiresource := ar.(*ApiResource)
		mapReplicaSetInfo := apiresource.replicaset

		for _, key := range update_list {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			result, err := tx.Exec(context.Background(), UPDATE_REPLICASET_INFO, getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].Selector, update_info[key].Replicas,
				update_info[key].FullyLabeledReplicas, update_info[key].ReadyReplicas, update_info[key].AvailableReplicas, update_info[key].ObservedGeneneration,
				update_info[key].ReferenceKind, update_info[key].ReferenceUID, ontunetime, update_info[key].UID, common.ClusterID[update_info[key].Host])
			if !errorCheck(err) {
				return
			}
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_REPLICASET_INFO, getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, update_info[key].UID, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].Selector,
					update_info[key].Replicas, update_info[key].FullyLabeledReplicas, update_info[key].ReadyReplicas, update_info[key].AvailableReplicas, update_info[key].ObservedGeneneration,
					update_info[key].ReferenceKind, update_info[key].ReferenceUID, 1, ontunetime, ontunetime, common.ClusterID[update_info[key].Host])
				if !errorCheck(insert_err) {
					return
				}
			}
			var update_data kubeapi.MappingReplicaSet
			update_data.UID = update_info[key].UID
			update_data.Name = update_info[key].Name
			update_data.NamespaceName = update_info[key].NamespaceName
			update_data.StartTime = update_info[key].StartTime
			update_data.Host = update_info[key].Host
			update_data.Labels = update_info[key].Labels
			update_data.Selector = update_info[key].Selector
			update_data.Replicas = update_info[key].Replicas
			update_data.FullyLabeledReplicas = update_info[key].FullyLabeledReplicas
			update_data.ReadyReplicas = update_info[key].ReadyReplicas
			update_data.AvailableReplicas = update_info[key].AvailableReplicas
			update_data.ObservedGeneneration = update_info[key].ObservedGeneneration
			update_data.ReferenceKind = update_info[key].ReferenceKind
			update_data.ReferenceUID = update_info[key].ReferenceUID
			update_data.Host = update_info[key].Host
			mapReplicaSetInfo[update_info[key].UID] = update_data

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		common.LogManager.WriteLog(fmt.Sprintf("Manager DB - %s replicaset update is completed", hostname))

		apiresource.replicaset = mapReplicaSetInfo
		mapApiResource.Store(hostname, apiresource)
	}
}

func insertReplicasetinfo(ArrResource ReplicaSetinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(ArrResource.ArrRsname); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_REPLICASET_INFO, ArrResource.GetArgs()...)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	common.LogManager.WriteLog(fmt.Sprintf("Manager DB - replicaset insertion is completed: %s", strings.Join(ArrResource.ArrUid, ",")))

	updateTableinfo(TB_KUBE_RS_INFO, ontunetime)
}

func pvc_updateCheck(new_info map[string]kubeapi.MappingPvc, old_info map[string]kubeapi.MappingPvc) []string {
	var UpdateList []string
	if len(new_info) != len(old_info) { // 기존데이터와 새로운 데이터의 갯수가 다르다면 업데이트 필요
		for key, d := range new_info {
			old_data := old_info[key]
			if old_data.UID == "" {
				UpdateList = append(UpdateList, d.UID)
			}
		}

		return UpdateList
	}
	for i, d := range new_info {
		old_infodata := old_info[i]
		if old_infodata.UID == "" || !reflect.DeepEqual(d, old_infodata) {
			UpdateList = append(UpdateList, d.UID)
		}
	}

	return UpdateList
}

func updateEnablePvcinfo(ontunetime int64, update_info map[string]kubeapi.MappingPvc) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	var uidhostinfo map[string]string = make(map[string]string)
	for uid, data := range update_info {
		uidhostinfo[uid] = data.Host
	}

	uid, hostip := getUidHost(uidhostinfo)
	result, err := conn.Query(context.Background(), "update "+TB_KUBE_PVC_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where uid not in "+uid+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING uid")
	if !errorCheck(err) {
		return 0
	}

	var updateUid string
	var returnVal int64
	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapPvcInfo := apiresource.persistentvolumeclaim
		for result.Next() {
			err := result.Scan(&updateUid)
			if !errorCheck(err) {
				return 0
			}
			delete(mapPvcInfo, updateUid)
			returnVal++
		}
		apiresource.persistentvolumeclaim = mapPvcInfo
		mapApiResource.Store(hostip, apiresource)
	}

	result.Close()

	return returnVal
}

func updatePvcinfo(update_info map[string]kubeapi.MappingPvc, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	hostname := update_info[update_list[0]].Host

	if ar, ok := mapApiResource.Load(hostname); ok {
		apiresource := ar.(*ApiResource)
		mapPvcInfo := apiresource.persistentvolumeclaim

		for _, key := range update_list {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			result, err := tx.Exec(context.Background(), UPDATE_PVC_INFO, getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].Selector, update_info[key].AccessModes,
				update_info[key].RequestStorage, update_info[key].Status, getUID("storageclass", update_info[key].Host, update_info[key].StorageClassName), ontunetime, update_info[key].UID, common.ClusterID[update_info[key].Host])
			if !errorCheck(err) {
				return
			}
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_PVC_INFO, getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, update_info[key].UID, getStarttime(update_info[key].StartTime.Unix(), biastime),
					update_info[key].Labels, update_info[key].Selector, update_info[key].AccessModes,
					update_info[key].RequestStorage, update_info[key].Status, getUID("storageclass", update_info[key].Host, update_info[key].StorageClassName), 1, ontunetime, ontunetime, common.ClusterID[update_info[key].Host])
				if !errorCheck(insert_err) {
					return
				}
			}
			var update_data kubeapi.MappingPvc
			update_data.UID = update_info[key].UID
			update_data.Name = update_info[key].Name
			update_data.NamespaceName = update_info[key].NamespaceName
			update_data.StartTime = update_info[key].StartTime
			update_data.Host = update_info[key].Host
			update_data.Labels = update_info[key].Labels
			update_data.Selector = update_info[key].Selector
			update_data.AccessModes = update_info[key].AccessModes
			update_data.RequestStorage = update_info[key].RequestStorage
			update_data.Status = update_info[key].Status
			update_data.StorageClassName = update_info[key].StorageClassName
			update_data.Host = update_info[key].Host
			mapPvcInfo[update_info[key].UID] = update_data

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		common.LogManager.WriteLog(fmt.Sprintf("Manager DB - %s persistentvolumeclaim update is completed", hostname))

		apiresource.persistentvolumeclaim = mapPvcInfo
		mapApiResource.Store(hostname, apiresource)
	}
}

func insertPvcinfo(ArrResource Pvcinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(ArrResource.ArrPvcUid); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_PVC_INFO, ArrResource.GetArgs()...)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	common.LogManager.WriteLog(fmt.Sprintf("Manager DB - persistentvolumeclaim insertion is completed: %s", strings.Join(ArrResource.ArrPvcUid, ",")))

	updateTableinfo(TB_KUBE_PVC_INFO, ontunetime)
}

func pv_updateCheck(new_info map[string]kubeapi.MappingPv, old_info map[string]kubeapi.MappingPv) []string {
	var UpdateList []string
	if len(new_info) != len(old_info) { // 기존데이터와 새로운 데이터의 갯수가 다르다면 업데이트 필요
		for key, d := range new_info {
			old_data := old_info[key]
			if old_data.UID == "" {
				UpdateList = append(UpdateList, d.UID)
			}
		}

		return UpdateList
	}
	for i, d := range new_info {
		old_infodata := old_info[i]
		if old_infodata.UID == "" || !reflect.DeepEqual(d, old_infodata) {
			UpdateList = append(UpdateList, d.UID)
		}
	}

	return UpdateList
}

func updateEnablePvinfo(ontunetime int64, update_info map[string]kubeapi.MappingPv) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	var uidhostinfo map[string]string = make(map[string]string)
	for uid, data := range update_info {
		uidhostinfo[uid] = data.Host
	}

	uid, hostip := getUidHost(uidhostinfo)
	result, err := conn.Query(context.Background(), "update "+TB_KUBE_PV_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where pvuid not in "+uid+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING pvuid")
	if !errorCheck(err) {
		return 0
	}

	var updateUid string
	var returnVal int64
	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapPvInfo := apiresource.persistentvolume
		for result.Next() {
			err := result.Scan(&updateUid)
			if !errorCheck(err) {
				return 0
			}
			delete(mapPvInfo, updateUid)
			returnVal++
		}
		apiresource.persistentvolume = mapPvInfo
		mapApiResource.Store(hostip, apiresource)
	}

	result.Close()

	return returnVal
}

func updatePvinfo(update_info map[string]kubeapi.MappingPv, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	hostname := update_info[update_list[0]].Host

	if ar, ok := mapApiResource.Load(hostname); ok {
		apiresource := ar.(*ApiResource)
		mapPvInfo := apiresource.persistentvolume

		for _, key := range update_list {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			result, err := tx.Exec(context.Background(), UPDATE_PV_INFO, update_info[key].Name, update_info[key].PvcUID, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels, update_info[key].AccessModes,
				update_info[key].Capacity, update_info[key].ReclaimPolicy, update_info[key].Status, ontunetime, update_info[key].UID, common.ClusterID[update_info[key].Host])
			if !errorCheck(err) {
				return
			}
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_PV_INFO, update_info[key].Name, update_info[key].UID, update_info[key].PvcUID, getStarttime(update_info[key].StartTime.Unix(), biastime),
					update_info[key].Labels, update_info[key].AccessModes,
					update_info[key].Capacity, update_info[key].ReclaimPolicy, update_info[key].Status, 1, ontunetime, ontunetime, common.ClusterID[update_info[key].Host])
				if !errorCheck(insert_err) {
					return
				}
			}
			var update_data kubeapi.MappingPv
			update_data.UID = update_info[key].UID
			update_data.Name = update_info[key].Name
			update_data.PvcUID = update_info[key].PvcUID
			update_data.StartTime = update_info[key].StartTime
			update_data.Host = update_info[key].Host
			update_data.Labels = update_info[key].Labels
			update_data.AccessModes = update_info[key].AccessModes
			update_data.Capacity = update_info[key].Capacity
			update_data.ReclaimPolicy = update_info[key].ReclaimPolicy
			update_data.Status = update_info[key].Status
			update_data.Host = update_info[key].Host
			mapPvInfo[update_info[key].UID] = update_data

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		common.LogManager.WriteLog(fmt.Sprintf("Manager DB - %s persistentvolume update is completed", hostname))

		apiresource.persistentvolume = mapPvInfo
		mapApiResource.Store(hostname, apiresource)
	}
}

func insertPvinfo(ArrResource Pvinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(ArrResource.ArrPvUid); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_PV_INFO, ArrResource.GetArgs()...)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	common.LogManager.WriteLog(fmt.Sprintf("Manager DB - persistentvolume insertion is completed: %s", strings.Join(ArrResource.ArrPvUid, ",")))

	updateTableinfo(TB_KUBE_PV_INFO, ontunetime)
}

func sc_updateCheck(new_info map[string]kubeapi.MappingStorageClass, old_info map[string]kubeapi.MappingStorageClass) []string {
	var UpdateList []string
	if len(new_info) != len(old_info) { // 기존데이터와 새로운 데이터의 갯수가 다르다면 업데이트 필요
		for key, d := range new_info {
			old_data := old_info[key]
			if old_data.UID == "" {
				UpdateList = append(UpdateList, d.UID)
			}
		}

		return UpdateList
	}
	for i, d := range new_info {
		old_infodata := old_info[i]
		if old_infodata.UID == "" || !reflect.DeepEqual(d, old_infodata) {
			UpdateList = append(UpdateList, d.UID)
		}
	}

	return UpdateList
}

func updateEnableScinfo(ontunetime int64, update_info map[string]kubeapi.MappingStorageClass) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	var uidhostinfo map[string]string = make(map[string]string)
	for uid, data := range update_info {
		uidhostinfo[uid] = data.Host
	}

	uid, hostip := getUidHost(uidhostinfo)
	result, err := conn.Query(context.Background(), "update "+TB_KUBE_SC_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where uid not in "+uid+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING uid")
	if !errorCheck(err) {
		return 0
	}

	var updateUid string
	var returnVal int64
	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapStorageClassInfo := apiresource.storageclass
		for result.Next() {
			err := result.Scan(&updateUid)
			if !errorCheck(err) {
				return 0
			}
			delete(mapStorageClassInfo, updateUid)
			returnVal++
		}
		apiresource.storageclass = mapStorageClassInfo
		mapApiResource.Store(hostip, apiresource)
	}

	result.Close()

	return returnVal
}

func updateScinfo(update_info map[string]kubeapi.MappingStorageClass, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	hostname := update_info[update_list[0]].Host

	if ar, ok := mapApiResource.Load(hostname); ok {
		apiresource := ar.(*ApiResource)
		mapStorageClassInfo := apiresource.storageclass

		var iVolexp int
		for _, key := range update_list {
			if update_info[key].AllowVolumeExpansion {
				iVolexp = 1
			} else {
				iVolexp = 0
			}

			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			result, err := tx.Exec(context.Background(), UPDATE_SC_INFO, common.ClusterID[update_info[key].Host], update_info[key].Name, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels,
				update_info[key].Provisioner, update_info[key].ReclaimPolicy, update_info[key].VolumeBindingMode, iVolexp, ontunetime, update_info[key].UID)
			if !errorCheck(err) {
				return
			}
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_SC_INFO, common.ClusterID[update_info[key].Host], update_info[key].Name, update_info[key].UID, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels,
					update_info[key].Provisioner, update_info[key].ReclaimPolicy, update_info[key].VolumeBindingMode, iVolexp, 1, ontunetime, ontunetime)
				if !errorCheck(insert_err) {
					return
				}
			}
			var update_data kubeapi.MappingStorageClass
			update_data.UID = update_info[key].UID
			update_data.Name = update_info[key].Name
			update_data.StartTime = update_info[key].StartTime
			update_data.Host = update_info[key].Host
			update_data.Labels = update_info[key].Labels
			update_data.Provisioner = update_info[key].Provisioner
			update_data.ReclaimPolicy = update_info[key].ReclaimPolicy
			update_data.VolumeBindingMode = update_info[key].VolumeBindingMode
			update_data.AllowVolumeExpansion = update_info[key].AllowVolumeExpansion
			mapStorageClassInfo[update_info[key].UID] = update_data

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		common.LogManager.WriteLog(fmt.Sprintf("Manager DB - %s storageclass update is completed", hostname))

		apiresource.storageclass = mapStorageClassInfo
		mapApiResource.Store(hostname, apiresource)
	}
}

func insertScinfo(ArrResource Scinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(ArrResource.ArrUid); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_SC_INFO, ArrResource.GetArgs()...)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	common.LogManager.WriteLog(fmt.Sprintf("Manager DB - storageclass insertion is completed: %s", strings.Join(ArrResource.ArrUid, ",")))

	updateTableinfo(TB_KUBE_SC_INFO, ontunetime)
}

func ing_updateCheck(new_info map[string]kubeapi.MappingIngress, old_info map[string]kubeapi.MappingIngress) []string {
	var UpdateList []string
	if len(new_info) != len(old_info) { // 기존데이터와 새로운 데이터의 갯수가 다르다면 업데이트 필요
		for key, d := range new_info {
			old_data := old_info[key]
			if old_data.UID == "" {
				UpdateList = append(UpdateList, d.UID)
			}
		}

		return UpdateList
	}
	for i, d := range new_info {
		old_infodata := old_info[i]
		if old_infodata.UID == "" || !reflect.DeepEqual(d, old_infodata) {
			UpdateList = append(UpdateList, d.UID)
		}
	}

	return UpdateList
}

func updateEnableInginfo(ontunetime int64, update_info map[string]kubeapi.MappingIngress) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	var uidhostinfo map[string]string = make(map[string]string)
	for uid, data := range update_info {
		uidhostinfo[uid] = data.Host
	}

	uid, hostip := getUidHost(uidhostinfo)
	result, err := conn.Query(context.Background(), "update "+TB_KUBE_ING_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where uid not in "+uid+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING uid")
	if !errorCheck(err) {
		return 0
	}

	var updateUid string
	var returnVal int64
	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapIngressInfo := apiresource.ingress
		for result.Next() {
			err := result.Scan(&updateUid)
			if !errorCheck(err) {
				return 0
			}
			delete(mapIngressInfo, updateUid)
			returnVal++
		}
		apiresource.ingress = mapIngressInfo
		mapApiResource.Store(hostip, apiresource)
	}

	result.Close()

	return returnVal
}

func updateInginfo(update_info map[string]kubeapi.MappingIngress, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	hostname := update_info[update_list[0]].Host

	if ar, ok := mapApiResource.Load(hostname); ok {
		apiresource := ar.(*ApiResource)
		mapIngressInfo := apiresource.ingress

		for _, key := range update_list {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			result, err := tx.Exec(context.Background(), UPDATE_ING_INFO, getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels,
				update_info[key].IngressClassName, ontunetime, update_info[key].UID, common.ClusterID[update_info[key].Host])
			if !errorCheck(err) {
				return
			}
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_ING_INFO, getUID(METRIC_VAR_NAMESPACE, update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, update_info[key].UID, getStarttime(update_info[key].StartTime.Unix(), biastime), update_info[key].Labels,
					update_info[key].IngressClassName, 1, ontunetime, ontunetime, common.ClusterID[update_info[key].Host])
				if !errorCheck(insert_err) {
					return
				}
			}
			var update_data kubeapi.MappingIngress
			update_data.UID = update_info[key].UID
			update_data.Name = update_info[key].Name
			update_data.StartTime = update_info[key].StartTime
			update_data.Host = update_info[key].Host
			update_data.Labels = update_info[key].Labels
			update_data.IngressClassName = update_info[key].IngressClassName
			mapIngressInfo[update_info[key].UID] = update_data

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		common.LogManager.WriteLog(fmt.Sprintf("Manager DB - %s ingress update is completed", hostname))

		apiresource.ingress = mapIngressInfo
		mapApiResource.Store(hostname, apiresource)
	}
}

func insertInginfo(ArrResource Inginfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(ArrResource.ArrUid); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_ING_INFO, ArrResource.GetArgs()...)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	common.LogManager.WriteLog(fmt.Sprintf("Manager DB - ingress insertion is completed: %s", strings.Join(ArrResource.ArrUid, ",")))

	updateTableinfo(TB_KUBE_ING_INFO, ontunetime)
}

func inghost_updateCheck(new_info map[string]kubeapi.MappingIngressHost, old_info map[string]kubeapi.MappingIngressHost) []string {
	var UpdateList []string
	if len(new_info) != len(old_info) { // 기존데이터와 새로운 데이터의 갯수가 다르다면 업데이트 필요
		for key, d := range new_info {
			old_data := old_info[key]
			if old_data.Hostname == "" {
				UpdateList = append(UpdateList, d.Hostname)
			}
		}

		return UpdateList
	}
	for i, d := range new_info {
		old_infodata := old_info[i]
		if old_infodata.Hostname == "" || !reflect.DeepEqual(d, old_infodata) {
			UpdateList = append(UpdateList, d.Hostname)
		}
	}

	return UpdateList
}

func updateEnableInghostinfo(ontunetime int64, update_info map[string]kubeapi.MappingIngressHost) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	var uidhostinfo map[string]string = make(map[string]string)
	for uid, data := range update_info {
		uidhostinfo[uid] = data.Host
	}

	name, hostip := getUidHost(uidhostinfo)
	result, err := conn.Query(context.Background(), "update "+TB_KUBE_INGHOST_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where hostname not in "+name+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING hostname")
	if !errorCheck(err) {
		return 0
	}

	var updateName string
	var returnVal int64
	if ar, ok := mapApiResource.Load(hostip); ok {
		apiresource := ar.(*ApiResource)
		mapIngressHostInfo := apiresource.ingresshost
		for result.Next() {
			err := result.Scan(&updateName)
			if !errorCheck(err) {
				return 0
			}
			delete(mapIngressHostInfo, updateName)
			returnVal++
		}
		apiresource.ingresshost = mapIngressHostInfo
		mapApiResource.Store(hostip, apiresource)
	}

	result.Close()

	return returnVal
}

func updateInghostinfo(update_info map[string]kubeapi.MappingIngressHost, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	hostname := update_info[update_list[0]].Host

	if ar, ok := mapApiResource.Load(hostname); ok {
		apiresource := ar.(*ApiResource)
		mapIngressHostInfo := apiresource.ingresshost

		for _, key := range update_list {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			result, err := tx.Exec(context.Background(), UPDATE_INGHOST_INFO, update_info[key].UID, update_info[key].BackendType, update_info[key].BackendName,
				update_info[key].PathType, update_info[key].Path, update_info[key].ServicePort, update_info[key].ResourceAPIGroup, update_info[key].ResourceKind,
				ontunetime, update_info[key].Hostname, common.ClusterID[update_info[key].Host])
			if !errorCheck(err) {
				return
			}
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_INGHOST_INFO, update_info[key].UID, update_info[key].BackendType, update_info[key].BackendName,
					update_info[key].Hostname, update_info[key].PathType, update_info[key].Path, update_info[key].ServicePort, update_info[key].ResourceAPIGroup,
					update_info[key].ResourceKind, 1, ontunetime, ontunetime, common.ClusterID[update_info[key].Host])
				if !errorCheck(insert_err) {
					return
				}
			}
			var update_data kubeapi.MappingIngressHost
			update_data.UID = update_info[key].UID
			update_data.Host = update_info[key].Host
			update_data.BackendType = update_info[key].BackendType
			update_data.BackendName = update_info[key].BackendName
			update_data.Hostname = update_info[key].Hostname
			update_data.PathType = update_info[key].PathType
			update_data.Path = update_info[key].Path
			update_data.ServicePort = update_info[key].ServicePort
			update_data.ResourceAPIGroup = update_info[key].ResourceAPIGroup
			update_data.ResourceKind = update_info[key].ResourceKind
			mapIngressHostInfo[update_info[key].Hostname] = update_data

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		common.LogManager.WriteLog(fmt.Sprintf("Manager DB - %s ingress host update is completed", hostname))

		apiresource.ingresshost = mapIngressHostInfo
		mapApiResource.Store(hostname, apiresource)
	}
}

func insertInghostinfo(ArrResource IngHostinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	for i := 0; i < len(ArrResource.ArrHostname); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_INGHOST_INFO, ArrResource.GetArgs()...)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	common.LogManager.WriteLog(fmt.Sprintf("Manager DB - ingress host insertion is completed: %s", strings.Join(ArrResource.ArrBackendname, ",")))

	updateTableinfo(TB_KUBE_INGHOST_INFO, ontunetime)
}
