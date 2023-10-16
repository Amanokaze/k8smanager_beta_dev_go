package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func ProcessTableVersion(tablever int) {
	// nowVer := 0
	// if tablever == nowVer {
	// 	//현재 테이블이 최종버전일 경우
	// 	//fmt.Println("tableVer : " + strconv.Itoa(tablever))
	// } else {
	// 	//현재 테이블이 최종버전이 아닐 경우... 여기서는 버전별 처리를 따로 해야 할듯.
	// 	//fmt.Println("tableVer : " + strconv.Itoa(tablever))
	// }
}

func getCreateResourceTableStmt(tablename string) string {
	if tablename == TB_KUBE_MANAGER_INFO {
		return CREATE_TABLE_MANAGER_INFO
	} else if tablename == TB_KUBE_CLUSTER_INFO {
		return CREATE_TABLE_CLUSTER_INFO
	} else if tablename == TB_KUBE_RESOURCE_INFO {
		return CREATE_TABLE_RESOURCE_INFO
	} else if tablename == TB_KUBE_NS_INFO {
		return CREATE_TABLE_NS_INFO
	} else if tablename == TB_KUBE_NODE_INFO {
		return CREATE_TABLE_NODE_INFO
	} else if tablename == TB_KUBE_POD_INFO {
		return CREATE_TABLE_POD_INFO
	} else if tablename == TB_KUBE_CONTAINER_INFO {
		return CREATE_TABLE_CONTAINER_INFO
	} else if tablename == TB_KUBE_SVC_INFO {
		return CREATE_TABLE_SVC_INFO
	} else if tablename == TB_KUBE_PVC_INFO {
		return CREATE_TABLE_PVC_INFO
	} else if tablename == TB_KUBE_PV_INFO {
		return CREATE_TABLE_PV_INFO
	} else if tablename == TB_KUBE_EVENT_INFO {
		return CREATE_TABLE_EVENT_INFO
	} else if tablename == TB_KUBE_LOG_INFO {
		return CREATE_TABLE_LOG_INFO
	} else if tablename == TB_KUBE_DEPLOY_INFO {
		return CREATE_TABLE_DEPLOY_INFO
	} else if tablename == TB_KUBE_STS_INFO {
		return CREATE_TABLE_STS_INFO
	} else if tablename == TB_KUBE_DS_INFO {
		return CREATE_TABLE_DS_INFO
	} else if tablename == TB_KUBE_RS_INFO {
		return CREATE_TABLE_RS_INFO
	} else if tablename == TB_KUBE_ING_INFO {
		return CREATE_TABLE_ING_INFO
	} else if tablename == TB_KUBE_INGHOST_INFO {
		return CREATE_TABLE_INGHOST_INFO
	} else if tablename == TB_KUBE_SC_INFO {
		return CREATE_TABLE_SC_INFO
	} else if tablename == TB_KUBE_FS_DEVICE_INFO {
		return CREATE_TABLE_FS_DEVICE_INFO
	} else if tablename == TB_KUBE_NET_INTERFACE_INFO {
		return CREATE_TABLE_NET_INTERFACE_INFO
	} else if tablename == TB_KUBE_METRIC_ID_INFO {
		return CREATE_TABLE_METRIC_ID_INFO
	} else if tablename == TB_KUBE_TABLE_INFO {
		return CREATE_TABLE_TABLE_INFO
	}

	return ""
}

func createResourceTable(conn *pgxpool.Conn, ontunetime int64, tablename string, cluster_idx_flag bool, enabled_idx_flag bool) bool {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)

		return false
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
			return false
		}

		_, err = tx.Exec(context.Background(), getCreateResourceTableStmt(tablename)+getTableSpaceStmt(RESOURCE_TABLE, 0))
		if !errorCheck(err) {
			return false
		}

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, 0)
		if !errorCheck(err) {
			return false
		}

		if cluster_idx_flag {
			createIndex(tx, tablename, "clusterid", "clusterid")
		}

		if enabled_idx_flag {
			createIndex(tx, tablename, "enabled", "enabled")
		}

		err = tx.Commit(context.Background())

		return errorCheck(errors.Wrap(err, "Commit error"))
	}
}

func createResourceView(conn *pgxpool.Conn, ontunetime int64, viewname string, viewstmt string, creation_flag bool) {
	if creation_flag {
		//Create View
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
			return
		}

		_, err = tx.Exec(context.Background(), viewstmt)
		if !errorCheck(err) {
			return
		}

		err = tx.Commit(context.Background())
		if !errorCheck(errors.Wrap(err, "Commit error")) {
			return
		}
	}
}

func getCreateMetricTableArgsMetric(argtype string) []any {
	args := make([]any, 0)

	switch argtype {
	case "basic_raw":
		args = append(args, getMetricColumns("cpu_raw"), getMetricColumns("memory_raw"), getMetricColumns("fsrw_raw"), getMetricColumns("process"))
	case "container_raw":
		args = append(args, getMetricColumns("cpu_raw"), getMetricColumns("memory_raw"), getMetricColumns("process"))
	case "net_raw":
		args = append(args, getMetricColumns("net_raw"))
	case "nodefs_raw":
		args = append(args, getMetricColumns("fs_pre_raw"), getMetricColumns("fsrw_raw"), getMetricColumns("fs_post_raw"))
	case "fs_raw":
		args = append(args, getMetricColumns("fsrw_raw"))
	case "basic":
		args = append(args, getMetricColumns("cpu"), getMetricColumns("cputotal"), getMetricColumns("memory"), getMetricColumns("memorysize"), getMetricColumns(METRIC_VAR_NET), getMetricColumns("fs"), trimLastComma(getMetricColumns("process")))
	case METRIC_VAR_POD:
		args = append(args, getMetricColumns("cpu"), getMetricColumns("cpureqlimit"), getMetricColumns("memory"), getMetricColumns("memoryreqlimit"), getMetricColumns(METRIC_VAR_NET), getMetricColumns("fs"), trimLastComma(getMetricColumns("process")))
	case METRIC_VAR_CONTAINER:
		args = append(args, getMetricColumns("cpu"), getMetricColumns("cpureqlimit"), getMetricColumns("memory"), getMetricColumns("memoryreqlimit"), getMetricColumns("fs"), trimLastComma(getMetricColumns("process")))
	case METRIC_VAR_NET:
		args = append(args, trimLastComma(getMetricColumns("net_noprefix")))
	case "fs":
		args = append(args, trimLastComma(getMetricColumns("fs_noprefix")))
	}

	return args
}

func getCreateMetricTableArgs(argtype string) []any {
	args := make([]any, 0)

	switch argtype {
	case "node_raw":
		args = append(args, getMetricKeyColumns("basic", METRIC_VAR_NODE))
		args = append(args, getCreateMetricTableArgsMetric("basic_raw")...)
	case "pod_raw":
		args = append(args, getMetricKeyColumns("basic", METRIC_VAR_POD))
		args = append(args, getCreateMetricTableArgsMetric("basic_raw")...)
	case "container_raw":
		args = append(args, getMetricKeyColumns("basic", METRIC_VAR_CONTAINER))
		args = append(args, getCreateMetricTableArgsMetric("container_raw")...)
	case "nodenet_raw":
		args = append(args, getMetricKeyColumns(METRIC_VAR_NET, METRIC_VAR_NODE))
		args = append(args, getCreateMetricTableArgsMetric("net_raw")...)
	case "podnet_raw":
		args = append(args, getMetricKeyColumns(METRIC_VAR_NET, METRIC_VAR_POD))
		args = append(args, getCreateMetricTableArgsMetric("net_raw")...)
	case "nodefs_raw":
		args = append(args, getMetricKeyColumns("fs", METRIC_VAR_NODE))
		args = append(args, getCreateMetricTableArgsMetric("nodefs_raw")...)
	case "podfs_raw":
		args = append(args, getMetricKeyColumns("fs", METRIC_VAR_POD))
		args = append(args, getCreateMetricTableArgsMetric("fs_raw")...)
	case "containerfs_raw":
		args = append(args, getMetricKeyColumns("fs", METRIC_VAR_CONTAINER))
		args = append(args, getCreateMetricTableArgsMetric("fs_raw")...)
	case METRIC_VAR_NODE:
		args = append(args, getMetricKeyColumns("basic", METRIC_VAR_NODE))
		args = append(args, getCreateMetricTableArgsMetric("basic")...)
	case METRIC_VAR_POD:
		args = append(args, getMetricKeyColumns("basic", METRIC_VAR_POD))
		args = append(args, getCreateMetricTableArgsMetric(METRIC_VAR_POD)...)
	case METRIC_VAR_CONTAINER:
		args = append(args, getMetricKeyColumns("basic", METRIC_VAR_CONTAINER))
		args = append(args, getCreateMetricTableArgsMetric(METRIC_VAR_CONTAINER)...)
	case METRIC_VAR_NODENET:
		args = append(args, getMetricKeyColumns(METRIC_VAR_NET, METRIC_VAR_NODE))
		args = append(args, getCreateMetricTableArgsMetric(METRIC_VAR_NET)...)
	case METRIC_VAR_PODNET:
		args = append(args, getMetricKeyColumns(METRIC_VAR_NET, METRIC_VAR_POD))
		args = append(args, getCreateMetricTableArgsMetric(METRIC_VAR_NET)...)
	case METRIC_VAR_NODEFS:
		args = append(args, getMetricKeyColumns("fs", METRIC_VAR_NODE))
		args = append(args, getCreateMetricTableArgsMetric("fs")...)
	case METRIC_VAR_PODFS:
		args = append(args, getMetricKeyColumns("fs", METRIC_VAR_POD))
		args = append(args, getCreateMetricTableArgsMetric("fs")...)
	case METRIC_VAR_CONTAINERFS:
		args = append(args, getMetricKeyColumns("fs", METRIC_VAR_CONTAINER))
		args = append(args, getCreateMetricTableArgsMetric("fs")...)
	case METRIC_VAR_CLUSTER:
		args = append(args, getMetricKeyColumns("stat", METRIC_VAR_CLUSTER))
		args = append(args, getCreateMetricTableArgsMetric("basic")...)
	case METRIC_VAR_NAMESPACE:
		args = append(args, getMetricKeyColumns("stat", METRIC_VAR_NAMESPACE))
		args = append(args, getCreateMetricTableArgsMetric("basic")...)
	case METRIC_VAR_DEPLOYMENT:
		args = append(args, getMetricKeyColumns("stat", METRIC_VAR_DEPLOYMENT))
		args = append(args, getCreateMetricTableArgsMetric("basic")...)
	case METRIC_VAR_STATEFULSET:
		args = append(args, getMetricKeyColumns("stat", METRIC_VAR_STATEFULSET))
		args = append(args, getCreateMetricTableArgsMetric("basic")...)
	case METRIC_VAR_DAEMONSET:
		args = append(args, getMetricKeyColumns("stat", METRIC_VAR_DAEMONSET))
		args = append(args, getCreateMetricTableArgsMetric("basic")...)
	case METRIC_VAR_REPLICASET:
		args = append(args, getMetricKeyColumns("stat", METRIC_VAR_REPLICASET))
		args = append(args, getCreateMetricTableArgsMetric("basic")...)
	}

	return args
}

func createMetricTable(conn *pgxpool.Conn, ontunetime int64, tablename string, tabletype int, stmt string, args []any) bool {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)

		return false
	} else {
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
			return false
		}

		insert_args := make([]any, 0)
		insert_args = append(insert_args, tablename)
		insert_args = append(insert_args, args...)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(stmt, insert_args...)+getTableSpaceStmt(tabletype, ontunetime))
		if !errorCheck(err) {
			return false
		}

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, getTableDuration(tabletype))
		if !errorCheck(err) {
			return false
		}

		if tabletype == LONGTERM_TABLE {
			setAutoVacuum(tx, tablename)
		}

		err = tx.Commit(context.Background())

		return errorCheck(errors.Wrap(err, "Commit error"))
	}
}

func createMetricIndex(conn *pgxpool.Conn, tablename string, midfix string, columns string, creation_flag bool) {
	if creation_flag {
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
			return
		}

		createIndex(tx, tablename, midfix, columns)

		err = tx.Commit(context.Background())
		if !errorCheck(errors.Wrap(err, "Commit error")) {
			return
		}
	}
}

func checkKubenodePerf(conn *pgxpool.Conn, ontunetime int64) {
	lastnode_raw_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_NODE_PERF_RAW, SHORTTERM_TABLE, CREATE_TABLE_BASIC_RAW, getCreateMetricTableArgs("node_raw"))
	node_raw_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_NODE_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_BASIC_RAW, getCreateMetricTableArgs("node_raw"))
	createMetricIndex(conn, TB_KUBE_LAST_NODE_PERF_RAW, "nodeuid", "nodeuid", lastnode_raw_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_NODE_PERF_RAW), "nodeuid", "nodeuid, ontunetime", node_raw_flag)

	lastnode_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_NODE_PERF, SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_NODE))
	node_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_NODE_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_NODE))
	createMetricIndex(conn, TB_KUBE_LAST_NODE_PERF, "nodeuid", "nodeuid", lastnode_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_NODE_PERF), "nodeuid", "nodeuid, ontunetime", node_flag)
}

func checkKubepodPerf(conn *pgxpool.Conn, ontunetime int64) {
	lastpod_raw_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_POD_PERF_RAW, SHORTTERM_TABLE, CREATE_TABLE_BASIC_RAW, getCreateMetricTableArgs("pod_raw"))
	pod_raw_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_POD_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_BASIC_RAW, getCreateMetricTableArgs("pod_raw"))
	createMetricIndex(conn, TB_KUBE_LAST_POD_PERF_RAW, "nodeuid", "nodeuid", lastpod_raw_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_POD_PERF_RAW), "poduid", "poduid, ontunetime", pod_raw_flag)

	lastpod_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_POD_PERF, SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_POD))
	pod_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_POD_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_POD))
	createMetricIndex(conn, TB_KUBE_LAST_POD_PERF, "nodeuid", "nodeuid", lastpod_flag)
	createMetricIndex(conn, TB_KUBE_LAST_POD_PERF, "poduid", "poduid", lastpod_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_POD_PERF), "poduid", "poduid, ontunetime", pod_flag)

	lastcluster_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_CLUSTER_PERF, SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_CLUSTER))
	lastns_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_NAMESPACE_PERF, SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_NAMESPACE))
	lastdeploy_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_DEPLOYMENT_PERF, SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_DEPLOYMENT))
	laststs_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_STATEFULSET_PERF, SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_STATEFULSET))
	lastds_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_DAEMONSET_PERF, SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_DAEMONSET))
	lastrs_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_REPLICASET_PERF, SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_REPLICASET))
	createMetricIndex(conn, TB_KUBE_LAST_CLUSTER_PERF, "clusterid", "clusterid", lastcluster_flag)
	createMetricIndex(conn, TB_KUBE_LAST_NAMESPACE_PERF, "nsuid", "nsuid", lastns_flag)
	createMetricIndex(conn, TB_KUBE_LAST_DEPLOYMENT_PERF, "deployuid", "deployuid", lastdeploy_flag)
	createMetricIndex(conn, TB_KUBE_LAST_STATEFULSET_PERF, "stsuid", "stsuid", laststs_flag)
	createMetricIndex(conn, TB_KUBE_LAST_DAEMONSET_PERF, "dsuid", "dsuid", lastds_flag)
	createMetricIndex(conn, TB_KUBE_LAST_REPLICASET_PERF, "rsuid", "rsuid", lastrs_flag)

	cluster_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_CLUSTER_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_CLUSTER))
	ns_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_NAMESPACE_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_NAMESPACE))
	deploy_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_DEPLOYMENT_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_DEPLOYMENT))
	sts_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_STATEFULSET_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_STATEFULSET))
	ds_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_DAEMONSET_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_DAEMONSET))
	rs_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_REPLICASET_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_REPLICASET))
	createMetricIndex(conn, getTableName(TB_KUBE_CLUSTER_PERF), "clusterid", "clusterid, ontunetime", cluster_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_NAMESPACE_PERF), "nsuid", "nsuid, ontunetime", ns_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_DEPLOYMENT_PERF), "deployuid", "deployuid, ontunetime", deploy_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_STATEFULSET_PERF), "stsuid", "stsuid, ontunetime", sts_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_DAEMONSET_PERF), "dsuid", "dsuid, ontunetime", ds_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_REPLICASET_PERF), "rsuid", "rsuid, ontunetime", rs_flag)
}

func checkKubecontainerPerf(conn *pgxpool.Conn, ontunetime int64) {
	lastcontainer_raw_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_CONTAINER_PERF_RAW, SHORTTERM_TABLE, CREATE_TABLE_FOUR_PARAMS_RAW, getCreateMetricTableArgs("container_raw"))
	container_raw_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_CONTAINER_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_FOUR_PARAMS_RAW, getCreateMetricTableArgs("container_raw"))
	createMetricIndex(conn, TB_KUBE_LAST_CONTAINER_PERF_RAW, "nodeuid", "nodeuid", lastcontainer_raw_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_CONTAINER_PERF_RAW), "containername", "containername, poduid, ontunetime", container_raw_flag)

	lastcontainer_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_CONTAINER_PERF, SHORTTERM_TABLE, CREATE_TABLE_CONTAINER_PERF, getCreateMetricTableArgs(METRIC_VAR_CONTAINER))
	container_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_CONTAINER_PERF), SHORTTERM_TABLE, CREATE_TABLE_CONTAINER_PERF, getCreateMetricTableArgs(METRIC_VAR_CONTAINER))
	createMetricIndex(conn, TB_KUBE_LAST_CONTAINER_PERF, "nodeuid", "nodeuid", lastcontainer_flag)
	createMetricIndex(conn, TB_KUBE_LAST_CONTAINER_PERF, "containername", "containername, poduid", lastcontainer_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_CONTAINER_PERF), "containername", "containername, poduid, ontunetime", container_flag)
}

func checkKubeNodeNetPerf(conn *pgxpool.Conn, ontunetime int64) {
	nodenet_raw_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_NODE_NET_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_TWO_PARAMS_RAW, getCreateMetricTableArgs("nodenet_raw"))
	createMetricIndex(conn, getTableName(TB_KUBE_NODE_NET_PERF_RAW), "nodeuid", "nodeuid, ontunetime", nodenet_raw_flag)

	lastnodenet_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_NODE_NET_PERF, SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_NODENET))
	nodenet_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_NODE_NET_PERF), SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_NODENET))
	createMetricIndex(conn, TB_KUBE_LAST_NODE_NET_PERF, "nodeuid", "nodeuid", lastnodenet_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_NODE_NET_PERF), "nodeuid", "nodeuid, ontunetime", nodenet_flag)
}

func checkKubePodNetPerf(conn *pgxpool.Conn, ontunetime int64) {
	podnet_raw_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_POD_NET_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_TWO_PARAMS_RAW, getCreateMetricTableArgs("podnet_raw"))
	createMetricIndex(conn, getTableName(TB_KUBE_POD_NET_PERF_RAW), "poduid", "poduid, ontunetime", podnet_raw_flag)

	lastpodnet_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_POD_NET_PERF, SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_PODNET))
	podnet_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_POD_NET_PERF), SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_PODNET))
	createMetricIndex(conn, TB_KUBE_LAST_POD_NET_PERF, "nodeuid", "nodeuid", lastpodnet_flag)
	createMetricIndex(conn, TB_KUBE_LAST_POD_NET_PERF, "poduid", "poduid", lastpodnet_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_POD_NET_PERF), "poduid", "poduid, ontunetime", podnet_flag)
}

func checkKubeNodeFsPerf(conn *pgxpool.Conn, ontunetime int64) {
	nodefs_raw_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_NODE_FS_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_FOUR_PARAMS_RAW, getCreateMetricTableArgs("nodefs_raw"))
	createMetricIndex(conn, getTableName(TB_KUBE_NODE_FS_PERF_RAW), "nodeuid", "nodeuid, ontunetime", nodefs_raw_flag)

	lastnodefs_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_NODE_FS_PERF, SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_NODEFS))
	nodefs_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_NODE_FS_PERF), SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_NODEFS))
	createMetricIndex(conn, TB_KUBE_LAST_NODE_FS_PERF, "nodeuid", "nodeuid", lastnodefs_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_NODE_FS_PERF), "nodeuid", "nodeuid, ontunetime", nodefs_flag)
}

func checkKubePodFsPerf(conn *pgxpool.Conn, ontunetime int64) {
	podfs_raw_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_POD_FS_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_TWO_PARAMS_RAW, getCreateMetricTableArgs("podfs_raw"))
	createMetricIndex(conn, getTableName(TB_KUBE_POD_FS_PERF_RAW), "poduid", "poduid, ontunetime", podfs_raw_flag)

	lastpodfs_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_POD_FS_PERF, SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_PODFS))
	podfs_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_POD_FS_PERF), SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_PODFS))
	createMetricIndex(conn, TB_KUBE_LAST_POD_FS_PERF, "nodeuid", "nodeuid", lastpodfs_flag)
	createMetricIndex(conn, TB_KUBE_LAST_POD_FS_PERF, "poduid", "poduid", lastpodfs_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_POD_FS_PERF), "poduid", "poduid, ontunetime", podfs_flag)
}

func checkKubeContainerFsPerf(conn *pgxpool.Conn, ontunetime int64) {
	containerfs_raw_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_CONTAINER_FS_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_TWO_PARAMS_RAW, getCreateMetricTableArgs("containerfs_raw"))
	createMetricIndex(conn, getTableName(TB_KUBE_CONTAINER_FS_PERF_RAW), "containername", "containername, poduid, ontunetime", containerfs_raw_flag)

	lastcontainerfs_flag := createMetricTable(conn, ontunetime, TB_KUBE_LAST_CONTAINER_FS_PERF, SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_CONTAINERFS))
	containerfs_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_CONTAINER_FS_PERF), SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_CONTAINERFS))
	createMetricIndex(conn, TB_KUBE_LAST_CONTAINER_FS_PERF, "nodeuid", "nodeuid", lastcontainerfs_flag)
	createMetricIndex(conn, TB_KUBE_LAST_CONTAINER_FS_PERF, "containername", "containername, poduid", lastcontainerfs_flag)
	createMetricIndex(conn, getTableName(TB_KUBE_CONTAINER_FS_PERF), "containername", "containername, poduid, ontunetime", containerfs_flag)
}

func checkKubeAvgPerf(conn *pgxpool.Conn, ontunetime int64) {
	node_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_NODE_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_NODE))
	createMetricIndex(conn, getTableName(TB_KUBE_NODE_AVG_PERF), "nodeuid", "nodeuid, ontunetime", node_flag)

	pod_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_POD_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_POD))
	createMetricIndex(conn, getTableName(TB_KUBE_POD_AVG_PERF), "poduid", "poduid, ontunetime", pod_flag)

	cluster_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_CLUSTER_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_CLUSTER))
	createMetricIndex(conn, getTableName(TB_KUBE_CLUSTER_AVG_PERF), "clusterid", "clusterid, ontunetime", cluster_flag)

	ns_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_NAMESPACE_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_NAMESPACE))
	createMetricIndex(conn, getTableName(TB_KUBE_NAMESPACE_AVG_PERF), "nsuid", "nsuid, ontunetime", ns_flag)

	deployment_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_DEPLOYMENT_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_DEPLOYMENT))
	createMetricIndex(conn, getTableName(TB_KUBE_DEPLOYMENT_AVG_PERF), "deployuid", "deployuid, ontunetime", deployment_flag)

	statefulset_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_STATEFULSET_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_STATEFULSET))
	createMetricIndex(conn, getTableName(TB_KUBE_STATEFULSET_AVG_PERF), "stsuid", "stsuid, ontunetime", statefulset_flag)

	daemonset_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_DAEMONSET_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_DAEMONSET))
	createMetricIndex(conn, getTableName(TB_KUBE_DAEMONSET_AVG_PERF), "dsuid", "dsuid, ontunetime", daemonset_flag)

	replicaset_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_REPLICASET_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_REPLICASET))
	createMetricIndex(conn, getTableName(TB_KUBE_REPLICASET_AVG_PERF), "rsuid", "rsuid, ontunetime", replicaset_flag)

	container_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_CONTAINER_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_CONTAINER_PERF, getCreateMetricTableArgs(METRIC_VAR_CONTAINER))
	createMetricIndex(conn, getTableName(TB_KUBE_CONTAINER_AVG_PERF), "containername", "containername, poduid, ontunetime", container_flag)

	nodenet_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_NODE_NET_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_NODENET))
	createMetricIndex(conn, getTableName(TB_KUBE_NODE_NET_AVG_PERF), "nodeuid", "nodeuid, ontunetime", nodenet_flag)

	podnet_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_POD_NET_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_PODNET))
	createMetricIndex(conn, getTableName(TB_KUBE_POD_NET_AVG_PERF), "poduid", "poduid, ontunetime", podnet_flag)

	nodefs_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_NODE_FS_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_NODEFS))
	createMetricIndex(conn, getTableName(TB_KUBE_NODE_FS_AVG_PERF), "nodeuid", "nodeuid, ontunetime", nodefs_flag)

	podfs_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_POD_FS_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_PODFS))
	createMetricIndex(conn, getTableName(TB_KUBE_POD_FS_AVG_PERF), "poduid", "poduid, ontunetime", podfs_flag)

	containerfs_flag := createMetricTable(conn, ontunetime, getTableName(TB_KUBE_CONTAINER_FS_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_CONTAINERFS))
	createMetricIndex(conn, getTableName(TB_KUBE_CONTAINER_FS_AVG_PERF), "containername", "containername, poduid, ontunetime", containerfs_flag)
}

// Daily Table
func getDailyTableName(tablename string) string {
	now := time.Now()
	createDay := now.AddDate(0, 0, 1).Format(DATE_FORMAT)

	return tablename + "_" + createDay
}

func dailyKubenodePerf(conn *pgxpool.Conn, ontunetime int64) {
	node_raw_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_NODE_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_BASIC_RAW, getCreateMetricTableArgs("node_raw"))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_NODE_PERF_RAW), "nodeuid", "nodeuid, ontunetime", node_raw_realtime_flag)

	node_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_NODE_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_NODE))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_NODE_PERF), "nodeuid", "nodeuid, ontunetime", node_realtime_flag)

	node_avg_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_NODE_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_NODE))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_NODE_AVG_PERF), "nodeuid", "nodeuid, ontunetime", node_avg_flag)
}

func dailyKubepodPerf(conn *pgxpool.Conn, ontunetime int64) {
	pod_raw_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_POD_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_BASIC_RAW, getCreateMetricTableArgs("pod_raw"))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_POD_PERF_RAW), "poduid", "poduid, ontunetime", pod_raw_realtime_flag)

	pod_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_POD_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_POD))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_POD_PERF), "poduid", "poduid, ontunetime", pod_realtime_flag)

	pod_avg_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_POD_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_POD))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_POD_AVG_PERF), "poduid", "poduid, ontunetime", pod_avg_flag)

	cluster_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_CLUSTER_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_CLUSTER))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_CLUSTER_PERF), "clusterid", "clusterid, ontunetime", cluster_realtime_flag)

	ns_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_NAMESPACE_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_NAMESPACE))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_NAMESPACE_PERF), "nsuid", "nsuid, ontunetime", ns_realtime_flag)

	deploy_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_DEPLOYMENT_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_DEPLOYMENT))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_DEPLOYMENT_PERF), "deployuid", "deployuid, ontunetime", deploy_realtime_flag)

	sts_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_STATEFULSET_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_STATEFULSET))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_STATEFULSET_PERF), "stsuid", "stsuid, ontunetime", sts_realtime_flag)

	ds_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_DAEMONSET_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_DAEMONSET))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_DAEMONSET_PERF), "dsuid", "dsuid, ontunetime", ds_realtime_flag)

	rs_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_REPLICASET_PERF), SHORTTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_REPLICASET))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_REPLICASET_PERF), "rsuid", "rsuid, ontunetime", rs_realtime_flag)

	cluster_avg_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_CLUSTER_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_CLUSTER))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_CLUSTER_AVG_PERF), "clusterid", "clusterid, ontunetime", cluster_avg_flag)

	ns_avg_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_NAMESPACE_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_NAMESPACE))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_NAMESPACE_AVG_PERF), "nsuid", "nsuid, ontunetime", ns_avg_flag)

	deploy_avg_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_DEPLOYMENT_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_DEPLOYMENT))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_DEPLOYMENT_AVG_PERF), "deployuid", "deployuid, ontunetime", deploy_avg_flag)

	sts_avg_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_STATEFULSET_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_STATEFULSET))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_STATEFULSET_AVG_PERF), "stsuid", "stsuid, ontunetime", sts_avg_flag)

	ds_avg_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_DAEMONSET_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_DAEMONSET))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_DAEMONSET_AVG_PERF), "dsuid", "dsuid, ontunetime", ds_avg_flag)

	rs_avg_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_REPLICASET_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_BASIC_PERF, getCreateMetricTableArgs(METRIC_VAR_REPLICASET))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_REPLICASET_AVG_PERF), "rsuid", "rsuid, ontunetime", rs_avg_flag)
}

func dailyKubecontainerPerf(conn *pgxpool.Conn, ontunetime int64) {
	container_raw_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_CONTAINER_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_FOUR_PARAMS_RAW, getCreateMetricTableArgs("container_raw"))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_CONTAINER_PERF_RAW), "containername", "containername, poduid, ontunetime", container_raw_realtime_flag)

	container_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_CONTAINER_PERF), SHORTTERM_TABLE, CREATE_TABLE_CONTAINER_PERF, getCreateMetricTableArgs(METRIC_VAR_CONTAINER))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_CONTAINER_PERF), "containername", "containername, poduid, ontunetime", container_realtime_flag)

	container_avg_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_CONTAINER_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_CONTAINER_PERF, getCreateMetricTableArgs(METRIC_VAR_CONTAINER))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_CONTAINER_AVG_PERF), "containername", "containername, poduid, ontunetime", container_avg_flag)
}

func dailyKubeNodeNetPerf(conn *pgxpool.Conn, ontunetime int64) {
	nodenet_raw_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_NODE_NET_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_TWO_PARAMS_RAW, getCreateMetricTableArgs("nodenet_raw"))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_NODE_NET_PERF_RAW), "nodeuid", "nodeuid, ontunetime", nodenet_raw_realtime_flag)

	nodenet_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_NODE_NET_PERF), SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_NODENET))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_NODE_NET_PERF), "nodeuid", "nodeuid, ontunetime", nodenet_realtime_flag)

	nodenet_avg_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_NODE_NET_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_NODENET))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_NODE_NET_AVG_PERF), "nodeuid", "nodeuid, ontunetime", nodenet_avg_flag)
}

func dailyKubePodNetPerf(conn *pgxpool.Conn, ontunetime int64) {
	podnet_raw_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_POD_NET_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_TWO_PARAMS_RAW, getCreateMetricTableArgs("podnet_raw"))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_POD_NET_PERF_RAW), "poduid", "poduid, ontunetime", podnet_raw_realtime_flag)

	podnet_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_POD_NET_PERF), SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_PODNET))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_POD_NET_PERF), "poduid", "poduid, ontunetime", podnet_realtime_flag)

	podnet_avg_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_POD_NET_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_PODNET))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_POD_NET_AVG_PERF), "poduid", "poduid, ontunetime", podnet_avg_flag)
}

func dailyKubeNodeFsPerf(conn *pgxpool.Conn, ontunetime int64) {
	nodefs_raw_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_NODE_FS_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_FOUR_PARAMS_RAW, getCreateMetricTableArgs("nodefs_raw"))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_NODE_FS_PERF_RAW), "nodeuid", "nodeuid, ontunetime", nodefs_raw_realtime_flag)

	nodefs_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_NODE_FS_PERF), SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_NODEFS))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_NODE_FS_PERF), "nodeuid", "nodeuid, ontunetime", nodefs_realtime_flag)

	nodefs_avg_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_NODE_FS_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_NODEFS))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_NODE_FS_AVG_PERF), "nodeuid", "nodeuid, ontunetime", nodefs_avg_flag)
}

func dailyKubePodFsPerf(conn *pgxpool.Conn, ontunetime int64) {
	podfs_raw_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_POD_FS_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_TWO_PARAMS_RAW, getCreateMetricTableArgs("podfs_raw"))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_POD_FS_PERF_RAW), "poduid", "poduid, ontunetime", podfs_raw_realtime_flag)

	podfs_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_POD_FS_PERF), SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_PODFS))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_POD_FS_PERF), "poduid", "poduid, ontunetime", podfs_realtime_flag)

	podfs_avg_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_POD_FS_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_PODFS))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_POD_FS_AVG_PERF), "poduid", "poduid, ontunetime", podfs_avg_flag)
}

func dailyKubeContainerFsPerf(conn *pgxpool.Conn, ontunetime int64) {
	containerfs_raw_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_CONTAINER_FS_PERF_RAW), SHORTTERM_TABLE, CREATE_TABLE_TWO_PARAMS_RAW, getCreateMetricTableArgs("containerfs_raw"))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_CONTAINER_FS_PERF_RAW), "containername", "containername, poduid, ontunetime", containerfs_raw_realtime_flag)

	containerfs_realtime_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_CONTAINER_FS_PERF), SHORTTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_CONTAINERFS))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_CONTAINER_FS_PERF), "containername", "containername, poduid, ontunetime", containerfs_realtime_flag)

	containerfs_avg_flag := createMetricTable(conn, ontunetime, getDailyTableName(TB_KUBE_CONTAINER_FS_AVG_PERF), LONGTERM_TABLE, CREATE_TABLE_DETAILS_PERF, getCreateMetricTableArgs(METRIC_VAR_CONTAINERFS))
	createMetricIndex(conn, getDailyTableName(TB_KUBE_CONTAINER_FS_AVG_PERF), "containername", "containername, poduid, ontunetime", containerfs_avg_flag)
}

func createIndex(tx pgx.Tx, tablename string, postfix string, column string) {
	SQL := "CREATE INDEX IF NOT EXISTS " + tablename + "_" + postfix + "_idx ON " + tablename + " using btree ( " + column + ")"
	_, err := tx.Exec(context.Background(), SQL)
	if !errorCheck(err) {
		return
	}
}

func existKubeTableinfo(conn *pgxpool.Conn, tablename string) bool {
	var kubetbcount int
	SQL := "select count(*) from " + TB_KUBE_TABLE_INFO + " where tablename = '" + tablename + "'"

	err := conn.QueryRow(context.Background(), SQL).Scan(&kubetbcount)
	if !errorCheck(err) {
		return false
	}

	if kubetbcount == 0 {
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
			return false
		}

		_, err = tx.Exec(context.Background(), fmt.Sprintf("DROP TABLE IF EXISTS %s", tablename))
		if !errorCheck(err) {
			return false
		}

		err = tx.Commit(context.Background())
		if !errorCheck(err) {
			return false
		}

		return false
	} else {
		var pgtbcount int
		err := conn.QueryRow(context.Background(), "select count(*) from pg_tables where tablename = $1", tablename).Scan(&pgtbcount)
		if !errorCheck(err) {
			return false
		}

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
			return false
		}

		if pgtbcount == 0 {
			_, err = tx.Exec(context.Background(), fmt.Sprintf("DELETE FROM %s where tablename = '%s'", TB_KUBE_TABLE_INFO, tablename))
			if !errorCheck(err) {
				return false
			}

			err = tx.Commit(context.Background())
			if !errorCheck(err) {
				return false
			}

			return false
		} else {
			return true
		}
	}
}

func getKubeTableinfoVer(conn *pgxpool.Conn, tablename string) int {
	var selVer int
	SQL := "select version from " + TB_KUBE_TABLE_INFO + " where tablename = '" + tablename + "'"

	err := conn.QueryRow(context.Background(), SQL).Scan(&selVer)
	if !errorCheck(err) {
		return 0
	}

	return selVer
}
