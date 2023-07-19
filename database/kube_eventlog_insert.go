package database

import (
	"context"
	"onTuneKubeManager/common"
	"onTuneKubeManager/kubeapi"
	"reflect"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

// var mapPodLogRecentTime map[string]time.Time = make(map[string]time.Time)
var mapEventInfo map[string]kubeapi.MappingEvent = make(map[string]kubeapi.MappingEvent)
var Channeleventlog_insert chan map[string]interface{} = make(chan map[string]interface{})

func EventlogSender() {
	for {
		eventlog_data := <-Channeleventlog_insert
		for key, data := range eventlog_data {
			if key == "event_update" {
				event_data := data.(map[string]kubeapi.MappingEvent)
				ontunetime, _ := GetOntuneTime()
				UpdateList := event_updateCheck(event_data, mapEventInfo)
				updateCnt := updateEnableEventinfo(ontunetime, event_data)
				if len(UpdateList) > 0 {
					updateEventinfo(event_data, UpdateList, ontunetime)
				}
				if updateCnt > 0 || len(UpdateList) > 0 {
					update_tableinfo(TB_KUBE_EVENT_INFO, ontunetime)
				}
			} else if key == "event_insert" {
				event_data := data.(Eventinfo)
				insertEventinfo(event_data)
				// } else if key == "log_insert" {
				// 	log_data := data.(Loginfo)
				// 	insertLoginfo(log_data)
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func event_updateCheck(new_info map[string]kubeapi.MappingEvent, old_info map[string]kubeapi.MappingEvent) []string {
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

func updateEnableEventinfo(ontunetime int64, update_info map[string]kubeapi.MappingEvent) int64 {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	var uid string
	var hostip string
	var returnVal int64

	if len(update_info) > 0 {
		for k, data := range update_info {
			if uid == "" {
				uid = "('" + k + "'"
				hostip = data.Host
			} else {
				uid = uid + ",'" + k + "'"
			}
		}
		uid = uid + ")"

		result, err := conn.Query(context.Background(), "update "+TB_KUBE_EVENT_INFO+" set enabled = 0, updatetime = "+strconv.FormatInt(ontunetime, 10)+" where uid not in "+uid+" and clusterid = "+strconv.Itoa(common.ClusterID[hostip])+" RETURNING uid")
		errorCheck(err)

		var updateUid string
		if ae, ok := mapApiEvent.Load(hostip); ok {
			apievent := ae.(*ApiEvent)
			mapEventInfo = apievent.event
			for result.Next() {
				err := result.Scan(&updateUid)
				errorCheck(err)
				delete(mapEventInfo, updateUid)
				returnVal++
			}
			apievent.event = mapEventInfo
			mapApiEvent.Store(hostip, apievent)
		}

		result.Close()
	}

	return returnVal
}

func insertEventinfo(ArrResource Eventinfo) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	ontunetime, _ := GetOntuneTime()
	for i := 0; i < len(ArrResource.ArrEventUid); i++ {
		ArrResource.ArrCreateTime[i] = ontunetime
		ArrResource.ArrUpdateTime[i] = ontunetime
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	_, err = tx.Exec(context.Background(), INSERT_UNNEST_EVENT_INFO, ArrResource.GetArgs()...)
	errorCheck(err)

	err = tx.Commit(context.Background())
	errorCheck(err)

	conn.Release()

	update_tableinfo(TB_KUBE_EVENT_INFO, ontunetime)
}

func updateEventinfo(update_info map[string]kubeapi.MappingEvent, update_list []string, ontunetime int64) {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	if ae, ok := mapApiEvent.Load(update_info[update_list[0]].Host); ok {
		apievent := ae.(*ApiEvent)
		mapEventInfo = apievent.event

		for _, key := range update_list {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			errorCheck(err)

			result, err := tx.Exec(context.Background(), UPDATE_EVENT_INFO, common.ClusterID[update_info[key].Host], getUID("namespace", update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name,
				getStarttime(update_info[key].Firsttime.Unix(), biastime), getStarttime(update_info[key].Lasttime.Unix(), biastime), update_info[key].Labels, update_info[key].Eventtype, update_info[key].Eventcount,
				update_info[key].ObjectKind, GetResourceObjectUID(update_info[key]), update_info[key].SourceComponent, update_info[key].SourceHost,
				update_info[key].Reason, update_info[key].Message, ontunetime, update_info[key].UID)
			errorCheck(err)
			n := result.RowsAffected()
			if n == 0 {
				_, insert_err := tx.Exec(context.Background(), INSERT_EVENT_INFO, getUID("namespace", update_info[key].Host, update_info[key].NamespaceName), update_info[key].Name, update_info[key].UID,
					getStarttime(update_info[key].Firsttime.Unix(), biastime), getStarttime(update_info[key].Lasttime.Unix(), biastime), update_info[key].Labels, update_info[key].Eventtype, update_info[key].Eventcount,
					update_info[key].ObjectKind, GetResourceObjectUID(update_info[key]), update_info[key].SourceComponent, update_info[key].SourceHost,
					update_info[key].Reason, update_info[key].Message, 1, ontunetime, ontunetime, common.ClusterID[update_info[key].Host])
				errorCheck(insert_err)
			}
			var update_data kubeapi.MappingEvent
			update_data.NamespaceName = update_info[key].NamespaceName
			update_data.Host = update_info[key].Host
			update_data.UID = update_info[key].UID
			update_data.Name = update_info[key].Name
			update_data.Firsttime = update_info[key].Firsttime
			update_data.Lasttime = update_info[key].Lasttime
			update_data.Labels = update_info[key].Labels
			update_data.Eventtype = update_info[key].Eventtype
			update_data.Eventcount = update_info[key].Eventcount
			update_data.ObjectKind = update_info[key].ObjectKind
			update_data.ObjectName = update_info[key].ObjectName
			update_data.SourceComponent = update_info[key].SourceComponent
			update_data.SourceHost = update_info[key].SourceHost
			update_data.Reason = update_info[key].Reason
			update_data.Message = update_info[key].Message
			mapEventInfo[update_info[key].UID] = update_data

			err = tx.Commit(context.Background())
			errorCheck(err)
		}

		apievent.event = mapEventInfo
		mapApiEvent.Store(update_info[update_list[0]].Host, apievent)
	}
}

// func insertLoginfo(ArrResource Loginfo) {
// 	conn, err := common.DBConnectionPool.Acquire(context.Background())
// if err != nil {
// 	errorCheck(err)
// }

// 	ontunetime, _ := GetOntuneTime()
// 	for i := 0; i < len(ArrResource.ArrPodUid); i++ {
// 		ArrResource.ArrCreateTime[i] = ontunetime
// 		ArrResource.ArrUpdateTime[i] = ontunetime
// 	}

// 	_, err := tx.Exec(context.Background(), INSERT_UNNEST_LOG_INFO, pq.StringArray(ArrResource.ArrLogType), pq.StringArray(ArrResource.ArrNsUid), pq.StringArray(ArrResource.ArrPodUid),
// 		pq.Array(ArrResource.ArrStarttime), pq.StringArray(ArrResource.ArrMessage),
// 		pq.Int64Array(ArrResource.ArrCreateTime), pq.Int64Array(ArrResource.ArrUpdateTime))
// 	errorCheck(err)
// 	update_tableinfo(TB_KUBE_LOG_INFO, ontunetime)

// defer conn.Release()
// }
