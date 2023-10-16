package database

import (
	"context"
	"fmt"
	"onTuneKubeManager/common"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func insertNodeAvgPerf(nobiastime int64, starttime int64, endtime int64) {
	var avgtablename string = getTableName(TB_KUBE_NODE_AVG_PERF)
	var realtimetablename string = getRealtimeTableName(TB_KUBE_NODE_PERF, nobiastime)

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_AVG_NODE_PERF, avgtablename, endtime, realtimetablename, starttime, endtime))
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(avgtablename, endtime)
}

func insertPodAvgPerf(nobiastime int64, starttime int64, endtime int64) {
	var avgtablename string = getTableName(TB_KUBE_POD_AVG_PERF)
	var realtimetablename string = getRealtimeTableName(TB_KUBE_POD_PERF, nobiastime)

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_AVG_POD_PERF, avgtablename, endtime, realtimetablename, starttime, endtime))
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(avgtablename, endtime)
}

func insertContainerAvgPerf(nobiastime int64, starttime int64, endtime int64) {
	var avgtablename string = getTableName(TB_KUBE_CONTAINER_AVG_PERF)
	var realtimetablename string = getRealtimeTableName(TB_KUBE_CONTAINER_PERF, nobiastime)

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_AVG_CONTAINER_PERF, avgtablename, endtime, realtimetablename, starttime, endtime))
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(avgtablename, endtime)
}

func insertNodeNetAvgPerf(nobiastime int64, starttime int64, endtime int64) {
	var avgtablename string = getTableName(TB_KUBE_NODE_NET_AVG_PERF)
	var realtimetablename string = getRealtimeTableName(TB_KUBE_NODE_NET_PERF, nobiastime)

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_AVG_NODE_NET_PERF, avgtablename, endtime, realtimetablename, starttime, endtime))
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(avgtablename, endtime)
}

func insertPodNetAvgPerf(nobiastime int64, starttime int64, endtime int64) {
	var avgtablename string = getTableName(TB_KUBE_POD_NET_AVG_PERF)
	var realtimetablename string = getRealtimeTableName(TB_KUBE_POD_NET_PERF, nobiastime)

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_AVG_POD_NET_PERF, avgtablename, endtime, realtimetablename, starttime, endtime))
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(avgtablename, endtime)
}

func insertNodeFsAvgPerf(nobiastime int64, starttime int64, endtime int64) {
	var avgtablename string = getTableName(TB_KUBE_NODE_FS_AVG_PERF)
	var realtimetablename string = getRealtimeTableName(TB_KUBE_NODE_FS_PERF, nobiastime)

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_AVG_NODE_FS_PERF, avgtablename, endtime, realtimetablename, starttime, endtime))
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(avgtablename, endtime)
}

func insertPodFsAvgPerf(nobiastime int64, starttime int64, endtime int64) {
	var avgtablename string = getTableName(TB_KUBE_POD_FS_AVG_PERF)
	var realtimetablename string = getRealtimeTableName(TB_KUBE_POD_FS_PERF, nobiastime)

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_AVG_POD_FS_PERF, avgtablename, endtime, realtimetablename, starttime, endtime))
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(avgtablename, endtime)
}

func insertContainerFsAvgPerf(nobiastime int64, starttime int64, endtime int64) {
	var avgtablename string = getTableName(TB_KUBE_CONTAINER_FS_AVG_PERF)
	var realtimetablename string = getRealtimeTableName(TB_KUBE_CONTAINER_FS_PERF, nobiastime)

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_AVG_CONTAINER_FS_PERF, avgtablename, endtime, realtimetablename, starttime, endtime))
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(avgtablename, endtime)
}

func insertClusterAvgPerf(nobiastime int64, starttime int64, endtime int64) {
	var avgtablename string = getTableName(TB_KUBE_CLUSTER_AVG_PERF)
	var realtimetablename string = getRealtimeTableName(TB_KUBE_CLUSTER_PERF, nobiastime)

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	tx, err := conn.Begin(context.Background())
	if !errorCheck(err) {
		return
	}

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_AVG_CLUSTER_PERF, avgtablename, endtime, realtimetablename, starttime, endtime))
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(avgtablename, endtime)
}

func insertNamespaceAvgPerf(nobiastime int64, starttime int64, endtime int64) {
	var avgtablename string = getTableName(TB_KUBE_NAMESPACE_AVG_PERF)
	var realtimetablename string = getRealtimeTableName(TB_KUBE_NAMESPACE_PERF, nobiastime)

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	tx, err := conn.Begin(context.Background())
	if !errorCheck(err) {
		return
	}

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_AVG_NS_PERF, avgtablename, endtime, realtimetablename, starttime, endtime))
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	updateTableinfo(avgtablename, endtime)

	conn.Release()
}

func insertWorkloadAvgPerf(nobiastime int64, starttime int64, endtime int64, kind string) {
	var avgtablename string
	var realtimetablename string
	var uidname string

	switch kind {
	case METRIC_VAR_DEPLOYMENT:
		avgtablename = getTableName(TB_KUBE_DEPLOYMENT_AVG_PERF)
		realtimetablename = getRealtimeTableName(TB_KUBE_DEPLOYMENT_PERF, nobiastime)
		uidname = "deployuid"
	case METRIC_VAR_STATEFULSET:
		avgtablename = getTableName(TB_KUBE_STATEFULSET_AVG_PERF)
		realtimetablename = getRealtimeTableName(TB_KUBE_STATEFULSET_PERF, nobiastime)
		uidname = "stsuid"
	case METRIC_VAR_DAEMONSET:
		avgtablename = getTableName(TB_KUBE_DAEMONSET_AVG_PERF)
		realtimetablename = getRealtimeTableName(TB_KUBE_DAEMONSET_PERF, nobiastime)
		uidname = "dsuid"
	case METRIC_VAR_REPLICASET:
		avgtablename = getTableName(TB_KUBE_REPLICASET_AVG_PERF)
		realtimetablename = getRealtimeTableName(TB_KUBE_REPLICASET_PERF, nobiastime)
		uidname = "rsuid"
	default:
		return
	}

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	tx, err := conn.Begin(context.Background())
	if !errorCheck(err) {
		return
	}

	_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_AVG_WORKLOAD_PERF, avgtablename, uidname, endtime, realtimetablename, starttime, endtime, uidname))
	if !errorCheck(err) {
		return
	}

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	updateTableinfo(avgtablename, endtime)

	conn.Release()
}
