package database

import (
	"context"
	"fmt"
	"onTuneKubeManager/common"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func insertNodeRawRealtimeperf(raw_node_map map[string]NodePerf, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_NODE_PERF_RAW)
	var lasttablename string = TB_KUBE_LAST_NODE_PERF_RAW

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, insert_stat_data := range raw_node_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODE_PERF, tablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - RDI Node insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	// Insert Last Data
	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(TRUNCATE_DATA, lasttablename))
	if !errorCheck(del_err) {
		return
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - RDI LastNode Delete %d", processid, time.Now().UnixNano()))

	for _, insert_stat_data := range raw_node_map {
		_, err = tx.Exec(context.Background(), INSERT_UNNEST_LAST_NODE_PERF_RAW, insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - RDI LastNode Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func insertPodRawRealtimeperf(raw_pod_map map[string]PodPerf, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_POD_PERF_RAW)
	var lasttablename string = TB_KUBE_LAST_POD_PERF_RAW

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, insert_stat_data := range raw_pod_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_POD_PERF, tablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - RDI Pod insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	// Insert Last Data
	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(TRUNCATE_DATA, lasttablename))
	if !errorCheck(del_err) {
		return
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - RDI LastPod Delete %d", processid, time.Now().UnixNano()))

	for _, insert_stat_data := range raw_pod_map {
		_, err = tx.Exec(context.Background(), INSERT_UNNEST_LAST_POD_PERF_RAW, insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - RDI LastPod Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func insertContainerRawRealtimeperf(raw_container_map map[string]ContainerPerf, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_CONTAINER_PERF_RAW)
	var lasttablename string = TB_KUBE_LAST_CONTAINER_PERF_RAW

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, insert_stat_data := range raw_container_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CONTAINER_PERF, tablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - RDI Container insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	// Insert Last Data
	tx, err = conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, del_err := tx.Exec(context.Background(), fmt.Sprintf(TRUNCATE_DATA, lasttablename))
	if !errorCheck(del_err) {
		return
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - RDI LastContainer Delete %d", processid, time.Now().UnixNano()))

	for _, insert_stat_data := range raw_container_map {
		_, err = tx.Exec(context.Background(), INSERT_UNNEST_LAST_CONTAINER_PERF_RAW, insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - RDI LastContainer Insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func insertNodeNetRawRealtimeperf(raw_nodenet_map map[string]NodeNetPerf, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_NODE_NET_PERF_RAW)
	var lasttablename string = TB_KUBE_LAST_NODE_NET_PERF_RAW

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, insert_stat_data := range raw_nodenet_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODE_NET_PERF, tablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - RDI NodeNet insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func insertPodNetRawRealtimeperf(raw_podnet_map map[string]PodNetPerf, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_POD_NET_PERF_RAW)
	var lasttablename string = TB_KUBE_LAST_POD_NET_PERF_RAW

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, insert_stat_data := range raw_podnet_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_POD_NET_PERF, tablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - RDI PodNet insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func insertNodeFsRawRealtimeperf(raw_nodefs_map map[string]NodeFsPerf, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_NODE_FS_PERF_RAW)
	var lasttablename string = TB_KUBE_LAST_NODE_FS_PERF_RAW

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, insert_stat_data := range raw_nodefs_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_NODE_FS_PERF, tablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - RDI NodeFs insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func insertPodFsRawRealtimeperf(raw_podfs_map map[string]PodFsPerf, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_POD_FS_PERF_RAW)
	var lasttablename string = TB_KUBE_LAST_POD_FS_PERF_RAW

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, insert_stat_data := range raw_podfs_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_POD_FS_PERF, tablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - RDI PodFs insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func insertContainerFsRawRealtimeperf(raw_containerfs_map map[string]ContainerFsPerf, ontunetime int64, processid string) {
	var tablename string = getTableName(TB_KUBE_CONTAINER_FS_PERF_RAW)
	var lasttablename string = TB_KUBE_LAST_CONTAINER_FS_PERF_RAW

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	// Insert Data
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	for _, insert_stat_data := range raw_containerfs_map {
		_, err = tx.Exec(context.Background(), fmt.Sprintf(INSERT_UNNEST_CONTAINER_FS_PERF, tablename), insert_stat_data.GetArgs()...)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("[%s] Manager DB - RDI ContainerFs insert %d", processid, time.Now().UnixNano()))

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	conn.Release()

	updateTableinfo(tablename, ontunetime)
	updateTableinfo(lasttablename, ontunetime)
}

func setRawOntunetime(size int, ontunetime int64) []int64 {
	var arr []int64
	for i := 0; i < size; i++ {
		arr = append(arr, ontunetime)
	}

	return arr
}

func rawDataSetting(metric_data *sync.Map, raw_map *BasicRawMap, nodeuid string, ontunetime int64, processid string) {
	if node_data, ok := metric_data.Load("node_insert"); ok {
		node_insert_data := node_data.(NodePerf)
		node_insert_data.ArrOntunetime = setRawOntunetime(len(node_insert_data.ArrOntunetime), ontunetime)
		raw_map.Node[nodeuid] = node_insert_data
		metric_data.Store("node_insert", node_insert_data)
	}
	if pod_data, ok := metric_data.Load("pod_insert"); ok {
		pod_insert_data := pod_data.(PodPerf)
		pod_insert_data.ArrOntunetime = setRawOntunetime(len(pod_insert_data.ArrOntunetime), ontunetime)
		raw_map.Pod[nodeuid] = pod_insert_data
		metric_data.Store("pod_insert", pod_insert_data)
	}
	if container_data, ok := metric_data.Load("container_insert"); ok {
		container_insert_data := container_data.(ContainerPerf)
		container_insert_data.ArrOntunetime = setRawOntunetime(len(container_insert_data.ArrOntunetime), ontunetime)
		raw_map.Container[nodeuid] = container_insert_data
		metric_data.Store("container_insert", container_insert_data)
	}
	if nodenet_data, ok := metric_data.Load("nodnet_insert"); ok {
		nodenet_insert_data := nodenet_data.(NodeNetPerf)
		nodenet_insert_data.ArrOntunetime = setRawOntunetime(len(nodenet_insert_data.ArrOntunetime), ontunetime)
		raw_map.NodeNet[nodeuid] = nodenet_insert_data
		metric_data.Store("nodnet_insert", nodenet_insert_data)
	}
	if podnet_data, ok := metric_data.Load("podnet_insert"); ok {
		podnet_insert_data := podnet_data.(PodNetPerf)
		podnet_insert_data.ArrOntunetime = setRawOntunetime(len(podnet_insert_data.ArrOntunetime), ontunetime)
		raw_map.PodNet[nodeuid] = podnet_insert_data
		metric_data.Store("podnet_insert", podnet_insert_data)
	}
	if nodefs_data, ok := metric_data.Load("nodefs_insert"); ok {
		nodefs_insert_data := nodefs_data.(NodeFsPerf)
		nodefs_insert_data.ArrOntunetime = setRawOntunetime(len(nodefs_insert_data.ArrOntunetime), ontunetime)
		raw_map.NodeFs[nodeuid] = nodefs_insert_data
		metric_data.Store("nodefs_insert", nodefs_insert_data)
	}
	if podfs_data, ok := metric_data.Load("podfs_insert"); ok {
		podfs_insert_data := podfs_data.(PodFsPerf)
		podfs_insert_data.ArrOntunetime = setRawOntunetime(len(podfs_insert_data.ArrOntunetime), ontunetime)
		raw_map.PodFs[nodeuid] = podfs_insert_data
		metric_data.Store("podfs_insert", podfs_insert_data)
	}
	if containerfs_data, ok := metric_data.Load("containerfs_insert"); ok {
		containerfs_insert_data := containerfs_data.(ContainerFsPerf)
		containerfs_insert_data.ArrOntunetime = setRawOntunetime(len(containerfs_insert_data.ArrOntunetime), ontunetime)
		raw_map.ContainerFs[nodeuid] = containerfs_insert_data
		metric_data.Store("containerfs_insert", containerfs_insert_data)
	}
}

func insertRawRealtimeData(raw_map *BasicRawMap, ontunetime int64, processid string) {
	insertNodeRawRealtimeperf(raw_map.Node, ontunetime, processid)
	insertPodRawRealtimeperf(raw_map.Pod, ontunetime, processid)
	insertContainerRawRealtimeperf(raw_map.Container, ontunetime, processid)
	insertNodeNetRawRealtimeperf(raw_map.NodeNet, ontunetime, processid)
	insertPodNetRawRealtimeperf(raw_map.PodNet, ontunetime, processid)
	insertNodeFsRawRealtimeperf(raw_map.NodeFs, ontunetime, processid)
	insertPodFsRawRealtimeperf(raw_map.PodFs, ontunetime, processid)
	insertContainerFsRawRealtimeperf(raw_map.ContainerFs, ontunetime, processid)
}
