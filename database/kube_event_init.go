package database

import (
	"onTuneKubeManager/common"
	"onTuneKubeManager/kubeapi"
	"strings"
	"sync"
	"time"
)

func InitEventInfo(host string, clusterid int) {
	mapEventInfo := make(map[string]kubeapi.MappingEvent)
	rowsCnt := selectRowCountEnabled(TB_KUBE_EVENT_INFO, clusterid)
	if rowsCnt > 0 {
		if len(mapEventInfo) == 0 {
			var name string
			var uid string
			var nsuid string
			var firsttime int64
			var lasttime int64
			var labels string
			var eventtype string
			var eventcount int
			var objkind string
			var objuid string
			var srccomponent string
			var srchost string
			var reason string
			var message string
			var enabled int
			rows := selectRowEnabled("eventname, uid, nsuid, firsttime, lasttime, labels, eventtype, eventcount, objkind, objuid, srccomponent, srchost, reason, message, enabled", TB_KUBE_EVENT_INFO, clusterid)
			for rows.Next() {
				err := rows.Scan(&name, &uid, &nsuid, &firsttime, &lasttime, &labels, &eventtype, &eventcount, &objkind, &objuid, &srccomponent, &srchost, &reason, &message, &enabled)
				if !errorCheck(err) {
					return
				}
				if enabled == 1 {
					var resource_data_temp kubeapi.MappingEvent
					if ns_map, ok := common.ResourceMap.Load(METRIC_VAR_NAMESPACE); ok {
						ns_map.(*sync.Map).Range(func(key, value any) bool {
							if value.(string) == nsuid {
								resource_data_temp.NamespaceName = key.(string)
								return false
							}

							return true
						})
					}

					if kind_map, ok := common.ResourceMap.Load(strings.ToLower(objkind)); ok {
						kind_map.(*sync.Map).Range(func(key, value any) bool {
							if value.(string) == objuid {
								resource_data_temp.ObjectName = key.(string)
								return false
							}

							return true
						})
					}

					resource_data_temp.Name = name
					resource_data_temp.UID = uid
					resource_data_temp.Firsttime = time.Unix(firsttime, 0)
					resource_data_temp.Lasttime = time.Unix(lasttime, 0)
					resource_data_temp.Host = host
					resource_data_temp.Labels = labels
					resource_data_temp.Eventtype = eventtype
					resource_data_temp.Eventcount = eventcount
					resource_data_temp.ObjectKind = objkind
					resource_data_temp.SourceComponent = srccomponent
					resource_data_temp.SourceHost = srchost
					resource_data_temp.Reason = reason
					resource_data_temp.Message = message
					mapEventInfo[resource_data_temp.UID] = resource_data_temp //일단 DB에 있는 데이터를 가져와서 메모리에 저장...
				}
			}
		}
	}

	if ae, ok := mapApiEvent.Load(host); ok {
		apievent := ae.(*ApiEvent)
		apievent.event = mapEventInfo
		mapApiEvent.Store(host, apievent)
	}
}

// ############################ Init Data Load
func InitEventDataMap() {
	defer errorRecover()

	for k, v := range common.ClusterID {
		mapApiEvent.Store(k, &ApiEvent{
			event: make(map[string]kubeapi.MappingEvent),
		})

		InitEventInfo(k, v)
	}

	// rowsCnt = selectRowCountEnabled(TB_KUBE_LOG_INFO)
	// if rowsCnt > 0 {
	// 	var logtype string
	// 	var poduid string
	// 	var starttime int64
	// 	rows := selectRowEnabled("logtype, poduid, starttime", TB_KUBE_LOG_INFO)
	// 	for rows.Next() {
	// 		err := rows.Scan(&logtype, &poduid, &starttime)
	// 		if !errorCheck(err) { return }

	// 		if logtype == "pod" {
	// 			mapPodLogRecentTime[poduid] = time.Unix(starttime, 0)
	// 		}
	// 	}
	// }
}
