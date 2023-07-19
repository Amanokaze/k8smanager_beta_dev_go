package database

import (
	"encoding/json"
	"onTuneKubeManager/common"
	"onTuneKubeManager/kubeapi"
	"time"
)

func EventlogReceive() {
	for {
		eventlog_msg := <-common.ChannelEventlogData
		if val, ok := eventlog_msg["event"]; ok {
			writeLog("this data is event info!!")
			setEventinfo(val)
			writeLog("event info End")
			// } else if val, ok := eventlog_msg["log"]; ok {
			// 	writeLog("this data is log info!!")
			// 	setLoginfo(val)
			// 	writeLog("log info End")
		} else {
			writeLog("no data??")
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func setEventinfo(info_msg []byte) {
	var tempEventInfo map[string]kubeapi.MappingEvent = make(map[string]kubeapi.MappingEvent)
	var map_data []kubeapi.MappingEvent
	var ArrResource Eventinfo
	var host string

	err := json.Unmarshal([]byte(info_msg), &map_data)
	errorCheck(err)

	for _, resource_data := range map_data {
		host = resource_data.Host

		ArrResource.ArrEventUid = append(ArrResource.ArrEventUid, resource_data.UID)
		ArrResource.ArrClusterid = append(ArrResource.ArrClusterid, common.ClusterID[resource_data.Host])
		ArrResource.ArrNsUid = append(ArrResource.ArrNsUid, getUID("namespace", resource_data.Host, resource_data.NamespaceName))
		ArrResource.ArrEventname = append(ArrResource.ArrEventname, resource_data.Name)
		ArrResource.ArrFirsttime = append(ArrResource.ArrFirsttime, getStarttime(resource_data.Firsttime.Unix(), biastime))
		ArrResource.ArrLasttime = append(ArrResource.ArrLasttime, getStarttime(resource_data.Lasttime.Unix(), biastime))
		ArrResource.ArrLabels = append(ArrResource.ArrLabels, resource_data.Labels)
		ArrResource.ArrEventtype = append(ArrResource.ArrEventtype, resource_data.Eventtype)
		ArrResource.ArrEventcount = append(ArrResource.ArrEventcount, resource_data.Eventcount)
		ArrResource.ArrObjkind = append(ArrResource.ArrObjkind, resource_data.ObjectKind)
		ArrResource.ArrObjUid = append(ArrResource.ArrObjUid, GetResourceObjectUID(resource_data))
		ArrResource.ArrSrccomponent = append(ArrResource.ArrSrccomponent, resource_data.SourceComponent)
		ArrResource.ArrSrchost = append(ArrResource.ArrSrchost, resource_data.SourceHost)
		ArrResource.ArrReason = append(ArrResource.ArrReason, resource_data.Reason)
		ArrResource.ArrMessage = append(ArrResource.ArrMessage, resource_data.Message)
		ArrResource.ArrEnabled = append(ArrResource.ArrEnabled, 1)
		ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
		ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)

		resource_data.Firsttime = time.Unix(getStarttime(resource_data.Firsttime.Unix(), biastime), 0)
		resource_data.Lasttime = time.Unix(getStarttime(resource_data.Lasttime.Unix(), biastime), 0)
		tempEventInfo[resource_data.UID] = resource_data // 신규로 들어온 데이터의 맵
	}

	var i_data interface{}
	map_info := make(map[string]interface{})

	if ae, ok := mapApiEvent.Load(host); ok {
		apievent := ae.(*ApiEvent)
		mapEventInfo := apievent.event

		if len(mapEventInfo) > 0 {
			i_data = tempEventInfo
			map_info["event_update"] = i_data
			Channeleventlog_insert <- map_info
		} else {
			i_data = ArrResource
			map_info["event_insert"] = i_data
			Channeleventlog_insert <- map_info
			if len(mapEventInfo) == 0 {
				for _, resource_data := range map_data {
					mapEventInfo[resource_data.UID] = resource_data
				}
			}
		}
	}
}

// func setLoginfo(info_msg []byte) {
// 	var tempLogInfo map[string]kubeapi.MappingLog = make(map[string]kubeapi.MappingLog)
// 	var map_data []kubeapi.MappingLog
// 	var ArrResource Loginfo

// 	err := json.Unmarshal([]byte(info_msg), &map_data)
// 	errorCheck(err)

// 	for _, resource_data := range map_data {
// 		if resource_data.Logtype == "pod" {
// 			poduid := getUID("pod", resource_data.NamespaceName, resource_data.Name)
// 			if val, ok := mapPodLogRecentTime[poduid]; ok {
// 				if val.Before(resource_data.Starttime) {
// 					mapPodLogRecentTime[poduid] = resource_data.Starttime
// 				} else {
// 					continue
// 				}
// 			} else {
// 				mapPodLogRecentTime[poduid] = resource_data.Starttime
// 			}

// 			ArrResource.ArrLogType = append(ArrResource.ArrLogType, resource_data.Logtype)
// 			ArrResource.ArrNsUid = append(ArrResource.ArrNsUid, getUID("namespace", resource_data.Host, resource_data.NamespaceName))
// 			ArrResource.ArrPodUid = append(ArrResource.ArrPodUid, poduid)
// 			ArrResource.ArrStarttime = append(ArrResource.ArrStarttime, getStarttime(resource_data.Starttime.Unix(), biastime))
// 			ArrResource.ArrMessage = append(ArrResource.ArrMessage, resource_data.Message)
// 			ArrResource.ArrCreateTime = append(ArrResource.ArrCreateTime, 0)
// 			ArrResource.ArrUpdateTime = append(ArrResource.ArrUpdateTime, 0)
// 			tempLogInfo[poduid] = resource_data // 신규로 들어온 데이터의 맵
// 		} else {
// 			continue
// 		}
// 	}

// 	var i_data interface{} = ArrResource
// 	map_info := make(map[string]interface{})
// 	map_info["log_insert"] = i_data
// 	Channeleventlog_insert <- map_info
// }
