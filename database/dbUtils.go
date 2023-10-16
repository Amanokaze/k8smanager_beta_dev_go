package database

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"onTuneKubeManager/common"
	"onTuneKubeManager/kubeapi"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func getUID(map_type string, parent_key interface{}, key interface{}) string {
	key_str := fmt.Sprintf("%s/%s", parent_key, key)
	map_key := strings.ToLower(map_type)

	wl_map, _ := common.ResourceMap.Load(map_key)
	if returnVal, ok := wl_map.(*sync.Map).Load(key_str); ok {
		return returnVal.(string)
	} else {
		return ""
	}
}

func getID(map_type string, data interface{}) int {
	data_str := data.(string)
	if data_str == "" {
		return 0
	}

	switch map_type {
	case "fsdevice":
		if _, ok := mapFsDeviceInfo[data_str]; !ok {
			updateFsDeviceId(data_str)
		}

		return mapFsDeviceInfo[data_str]
	case "netinterface":
		if _, ok := mapNetInterfaceInfo[data_str]; !ok {
			updateNetInterfaceId(data_str)
		}

		return mapNetInterfaceInfo[data_str]
	}

	return 0
}

func getMetricID(id interface{}, image interface{}) int {
	id_str := id.(string)
	image_str := image.(string)
	if id_str == "" {
		return 0
	}

	if _, ok := mapMetricIdInfo[id_str]; !ok {
		updateMetricId(id_str, image_str)
	}

	return mapMetricIdInfo[id_str]
}

func getStarttime(src_time int64, biastime int64) int64 {
	var starttime int64 = src_time - biastime
	if starttime > 0 {
		return starttime
	} else {
		return 0
	}
}

func getUidHost(uidhostinfo map[string]string) (string, string) {
	var uid string
	var hostip string

	for u, h := range uidhostinfo {
		if uid == "" {
			uid = "('" + u + "'"
			hostip = h
		} else {
			uid = uid + ",'" + u + "'"
		}
	}
	uid = uid + ")"

	return uid, hostip
}

func getStatTimes(current_data *sync.Map) (int64, bool) {
	var ontunetime int64
	if o, ok := current_data.Load("ontunetime"); ok {
		// prev_data의 ontunetime을 빼지 않고 RateInterval만큼의 수치를 대상으로 함
		ontunetime = o.(int64)
		return ontunetime, true
	} else {
		return 0, false
	}
}

func getMetricKeyColumns(perf_type string, col_type string) string {
	var col_stmt string
	switch col_type {
	case METRIC_VAR_NODE:
		col_stmt = `
		nodeuid					text NOT NULL default '',
		ontunetime   			bigint NOT NULL,
		`
	case METRIC_VAR_POD:
		col_stmt = `
		poduid					text NOT NULL default '',
		ontunetime   			bigint NOT NULL,
		nsuid					text NOT NULL default '',
		nodeuid					text NOT NULL default '',
		`
	case METRIC_VAR_CONTAINER:
		col_stmt = `
		containername     		text NOT NULL default '',
		ontunetime   			bigint NOT NULL,
		poduid					text NOT NULL default '',
		nsuid					text NOT NULL default '',
		nodeuid					text NOT NULL default '',
		`
	case METRIC_VAR_CLUSTER:
		col_stmt = "clusterid				int NOT NULL default 0,"
	case METRIC_VAR_NAMESPACE:
		col_stmt = "nsuid					text NOT NULL default '',"
	case METRIC_VAR_DEPLOYMENT:
		col_stmt = "deployuid				text NOT NULL default '',"
	case METRIC_VAR_STATEFULSET:
		col_stmt = "stsuid					text NOT NULL default '',"
	case METRIC_VAR_DAEMONSET:
		col_stmt = "dsuid					text NOT NULL default '',"
	case METRIC_VAR_REPLICASET:
		col_stmt = "rsuid					text NOT NULL default '',"
	}

	var postfix_stmt string
	switch perf_type {
	case "basic":
		postfix_stmt = "metricid				int NOT NULL default 0,"
	case METRIC_VAR_NET:
		postfix_stmt = `
		metricid					int NOT NULL default 0,
		interfaceid               	int NOT NULL default 0,
		`
	case "fs":
		postfix_stmt = `
		metricid					int NOT NULL default 0,
		deviceid               		int NOT NULL default 0,
		`
	case "stat":
		postfix_stmt = `
		ontunetime   			bigint NOT NULL,
		podcount			    int NOT NULL,
		`
	}

	return col_stmt + postfix_stmt
}

func getMetricColumns(col_type string) string {
	switch col_type {
	case "cpu_raw":
		return `
		cpuusagesecondstotal	double precision NOT NULL default 0,
		cpusystemsecondstotal 	double precision NOT NULL default 0,
		cpuusersecondstotal 	double precision NOT NULL default 0,`
	case "cpu":
		return `
		cpuusage 				double precision NOT NULL default 0,
		cpusystem 				double precision NOT NULL default 0,
		cpuuser 				double precision NOT NULL default 0,
		cpuusagecores 			double precision NOT NULL default 0,`
	case "cputotal":
		return `
		cputotalcores 			double precision NOT NULL default 0,`
	case "cpureqlimit":
		return `
		cpurequestcores			double precision NOT NULL default 0,
		cpulimitcores			double precision NOT NULL default 0,`
	case "memory_raw":
		return `
		memoryusagebytes 		double precision NOT NULL default 0,
		memoryworkingsetbytes 	double precision NOT NULL default 0,
		memorycache 			double precision NOT NULL default 0,
		memoryswap 				double precision NOT NULL default 0,
		memoryrss 				double precision NOT NULL default 0,`
	case "memory":
		return `
		memoryusage 			double precision NOT NULL default 0,
		memoryusagebytes		double precision NOT NULL default 0,`
	case "memorysize":
		return `
		memorysizebytes 		double precision NOT NULL default 0,
		memoryswap 				double precision NOT NULL default 0,`
	case "memoryreqlimit":
		return `
		memoryrequestbytes 		double precision NOT NULL default 0,
		memorylimitbytes		double precision NOT NULL default 0,
		memoryswap 				double precision NOT NULL default 0,`
	case "net_raw":
		return `
		networkreceivebytestotal 	double precision NOT NULL default 0,
		networkreceiveerrorstotal 	double precision NOT NULL default 0,
		networktransmitbytestotal 	double precision NOT NULL default 0,
		networktransmiterrorstotal 	double precision NOT NULL default 0,`
	case METRIC_VAR_NET:
		return `
		netiorate 				double precision NOT NULL default 0,
		netioerrors 			double precision NOT NULL default 0,
		netreceiverate 			double precision NOT NULL default 0,
		netreceiveerrors 		double precision NOT NULL default 0,
		nettransmitrate 		double precision NOT NULL default 0,
		nettransmiterrors 		double precision NOT NULL default 0,`
	case "net_noprefix":
		return `
		iorate 					double precision NOT NULL default 0,
		ioerrors 				double precision NOT NULL default 0,
		receiverate 			double precision NOT NULL default 0,
		receiveerrors		 	double precision NOT NULL default 0,
		transmitrate 			double precision NOT NULL default 0,
		transmiterrors		 	double precision NOT NULL default 0,`
	case "fs_pre_raw":
		return `
		fsinodesfree 				double precision NULL default 0,
		fsinodestotal 				double precision NULL default 0,
		fslimitbytes 				double precision NULL default 0,`
	case "fsrw_raw":
		return `
		fsreadsbytestotal 		double precision NOT NULL default 0,
		fswritesbytestotal 		double precision NOT NULL default 0,`
	case "fs_post_raw":
		return `
		fsusagebytes 				double precision NULL default 0,`
	case "fs":
		return `
		fsiorate 				double precision NOT NULL default 0,
		fsreadrate 				double precision NOT NULL default 0,
		fswriterate 			double precision NOT NULL default 0,`
	case "fs_noprefix":
		return `
		iorate 					double precision NOT NULL default 0,
		readrate 				double precision NOT NULL default 0,
		writerate	 			double precision NOT NULL default 0,`
	case "process":
		return `
		processes 				double precision NOT NULL default 0,`
	}

	return ""
}

func getTableDuration(tabletype int) int {
	if tabletype == SHORTTERM_TABLE {
		return common.ShorttermDuration
	} else if tabletype == LONGTERM_TABLE {
		return common.LongtermDuration
	} else {
		return 0
	}
}

func getTableName(tablename string) string {
	now := time.Now()
	createDay := now.Format(DATE_FORMAT)

	return tablename + "_" + createDay
}

func getPreviousTableName(tablename string) string {
	now := time.Now().Add(-1 * time.Hour)
	createDay := now.Format(DATE_FORMAT)

	return tablename + "_" + createDay
}

func getRealtimeTableName(tablename string, nobiastime int64) string {
	if nobiastime%DAY_SECONDS != 0 {
		return getTableName(tablename)
	} else {
		return getPreviousTableName(tablename)
	}
}

func updateFsDeviceId(name string) {
	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

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

	_, err = tx.Exec(context.Background(), INSERT_FSDEVICE_INFO, name, ontunetime, ontunetime)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	idvalue := selectOneValueCondition("deviceid", TB_KUBE_FS_DEVICE_INFO, "devicename", name)
	if idvalue == "" {
		return
	}

	idvalue_num, err := strconv.Atoi(idvalue)
	if !errorCheck(err) {
		return
	}

	mapFsDeviceInfo[name] = idvalue_num

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}
}

func updateNetInterfaceId(name string) {
	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

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

	_, insert_err := tx.Exec(context.Background(), INSERT_NETINTERFACE_INFO, name, ontunetime, ontunetime)
	if !errorCheck(insert_err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	idvalue := selectOneValueCondition("interfaceid", TB_KUBE_NET_INTERFACE_INFO, "interfacename", name)
	if idvalue == "" {
		return
	}

	idvalue_num, err := strconv.Atoi(idvalue)
	if !errorCheck(err) {
		return
	}

	mapNetInterfaceInfo[name] = idvalue_num

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}
}

func updateMetricId(name string, image string) {
	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

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

	_, insert_err := tx.Exec(context.Background(), INSERT_METRICID_INFO, name, image, ontunetime, ontunetime)
	if !errorCheck(insert_err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	idvalue := selectOneValueCondition("metricid", TB_KUBE_METRIC_ID_INFO, "metricname", name)
	if idvalue == "" {
		return
	}

	idvalue_num, err := strconv.Atoi(idvalue)
	if !errorCheck(err) {
		return
	}

	mapMetricIdInfo[name] = idvalue_num
	mapMetricImageInfo[idvalue_num] = image

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}
}

func updateTableinfo(tablename string, ontunetime int64) {
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

	_, err = tx.Exec(context.Background(), UPDATE_TABLEINFO, ontunetime, tablename)
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}
}

func selectRowEnabled(colName string, tablename string, clusterid int) pgx.Rows {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return nil
	}

	defer conn.Release()

	common.LogManager.Debug(fmt.Sprintf("select %s from %s where enabled=1 and clusterid=%d", colName, tablename, clusterid))
	rows, err := conn.Query(context.Background(), fmt.Sprintf("select %s from %s where enabled=1 and clusterid=%d", colName, tablename, clusterid))
	if !errorCheck(err) {
		return nil
	}

	return rows
}

func selectRowCountEnabled(tablename string, clusterid int) int {
	var cnt int
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), fmt.Sprintf("select count(*) as cnt from %s where enabled=1 and clusterid=%d", tablename, clusterid))
	if !errorCheck(err) {
		return 0
	}

	for rows.Next() {
		err := rows.Scan(&cnt)
		if !errorCheck(err) {
			return 0
		}
	}

	return cnt
}

func selectRowCount(tablename string) int {
	var cnt int
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return 0
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), "select count(*) as cnt from "+tablename)
	if !errorCheck(err) {
		return 0
	}

	for rows.Next() {
		err := rows.Scan(&cnt)
		if !errorCheck(err) {
			return 0
		}
	}

	return cnt
}

func selectOneValueCondition(colName string, tablename string, condition string, conditionVal string) string {
	var colValue string
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return ""
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), fmt.Sprintf("select %s from %s where %s='%s'", colName, tablename, condition, conditionVal))
	if !errorCheck(err) {
		return ""
	}

	for rows.Next() {
		err := rows.Scan(&colValue)
		if !errorCheck(err) {
			return ""
		}
	}

	return colValue
}

func selectCondition(colName string, tablename string, condition string, operator string, conditionVal string) pgx.Rows {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return nil
	}

	defer conn.Release()

	common.LogManager.Debug(fmt.Sprintf("select %s from %s where %s%s'%s'", colName, tablename, condition, operator, conditionVal))
	rows, err := conn.Query(context.Background(), fmt.Sprintf("select %s from %s where %s%s'%s'", colName, tablename, condition, operator, conditionVal))
	if !errorCheck(err) {
		return nil
	}

	return rows
}

func errorCheck(err error) bool {
	if err != nil {
		common.LogManager.WriteLog(fmt.Sprintf("Database Error - %v", err.Error()))
		return false
	}

	return true
}

func errorDisconnect(err error) {
	if common.ManagerDBConnected {
		common.LogManager.WriteLog(fmt.Sprintf("Database Error - %v", err.Error()))
		common.ManagerDBConnected = false
		common.ChannelManagerDBConnection <- false
	}
}

func errorRecover() {
	if r := recover(); r != nil {
		err := fmt.Errorf("Database Error - %v", r)
		common.LogManager.WriteLog(err.Error())
	}
}

func checkRowCount(err error) bool {
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return false
	} else {
		return true
	}
}

func checkMapNil(srcmap *sync.Map) bool {
	var flag bool = false
	srcmap.Range(func(k, v any) bool {
		flag = true
		return false
	})

	return flag
}

func trimLastComma(str string) string {
	if str == "" {
		return ""
	} else if str[len(str)-1:] == "," {
		return str[:len(str)-1]
	} else {
		return str
	}
}

func getStatData(stat_type string, current float64, previous float64, divider float64) float64 {
	const PERCENT = 100
	switch stat_type {
	case "rate":
		return math.Round((current - previous) / float64(common.RateInterval))
	case "subtract":
		return math.Round(current - previous)
	case "usage":
		if divider == 0 {
			return 0
		}

		return math.Round((current - previous) / float64(common.RateInterval) / divider * PERCENT * float64(common.TemporaryMultiply))
	case "current_usage":
		if divider == 0 {
			return 0
		}

		return math.Round(current / divider * PERCENT * float64(common.TemporaryMultiply))
	case "core":
		return math.Round((current - previous) / float64(common.RateInterval) * float64(common.TemporaryMultiply))
	}

	return 0
}

func getParentObject(kind string, host string, nsname string) string {
	parent_host_list := map[string]struct{}{
		METRIC_VAR_NAMESPACE: {},
		"node":               {},
		"persistentvolume":   {},
		"storageclass":       {},
	}

	if _, ok := parent_host_list[strings.ToLower(kind)]; ok {
		return host
	} else {
		return nsname
	}
}

func GetResourceObjectUID(resource_data kubeapi.MappingEvent) string {
	var enabledrsc bool = false
	for _, rsc := range kubeapi.EnabledResources {
		if resource_data.ObjectKind == rsc {
			enabledrsc = true
			break
		}
	}

	if enabledrsc {
		return getUID(resource_data.ObjectKind, getParentObject(resource_data.ObjectKind, resource_data.Host, resource_data.NamespaceName), resource_data.ObjectName)
	} else {
		return resource_data.ObjectName
	}
}

func getRandomProcessId() string {
	// charset use random string
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}
