package database

import (
	"context"
	"fmt"
	"onTuneKubeManager/common"
	"strings"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

const SECONDS = 60

// GetOntuneTime is get ontune time
// return: ontunetime, bias
func GetOntuneTime() (int64, int64) {
	var ontunetime int64
	var bias int64
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	err = conn.QueryRow(context.Background(), SELECT_ONTUNE_TIME).Scan(&ontunetime, &bias)
	errorCheckQueryRow(err)

	return ontunetime, bias * SECONDS
}

// QueryManagerinfo is to check and insert manager info
// managername: manager name
// description: manager description
// ip: manager ip
//
// return: manager id
func QueryManagerinfo(managername string, description string, ip string) int {
	var returnVal int
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	ontunetime, _ := GetOntuneTime()
	err = conn.QueryRow(context.Background(), SELECT_MANAGER_INFO, managername).Scan(&returnVal)
	rowcheck := RowCountCheck(err)
	if !rowcheck {
		insert_err := conn.QueryRow(context.Background(), INSERT_MANAGER_INFO_SQL, managername, description, ip, ontunetime, ontunetime).Scan(&returnVal)
		errorCheckQueryRow(insert_err)
	}

	conn.Release()

	update_tableinfo(TB_KUBE_MANAGER_INFO, ontunetime)

	return returnVal
}

// QueryClusterinfo is to check and insert cluster info
// managerid: manager id
// clustername: cluster name
// ctx: cluster context
// ip: cluster ip
//
// return: cluster id
func QueryClusterinfo(managerid int, clustername string, ctx string, ip string) int {
	var clusterid int
	var enabled int
	var status int
	var flag bool

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	ontunetime, _ := GetOntuneTime()

	rows, err := conn.Query(context.Background(), SELECT_CLUSTER_INFO, ip)
	errorCheck(err)

	if rows.Next() {
		err := rows.Scan(&clusterid, &enabled, &status)
		errorCheck(err)

		rows.Close()

		if enabled == 0 {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			errorCheck(err)

			_, err = tx.Exec(context.Background(), UPDATE_CLUSTER_ENABLED, clusterid, ontunetime)
			errorCheck(err)

			err = tx.Commit(context.Background())
			errorCheck(err)
		}
	} else {
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		err = tx.QueryRow(context.Background(), INSERT_CLUSTER_INFO_SQL, common.ManagerID, clustername, ctx, ip, ontunetime, ontunetime).Scan(&clusterid, &enabled)
		errorCheckQueryRow(err)

		err = tx.Commit(context.Background())
		errorCheck(err)

		flag = true
	}

	if status == 1 {
		common.ClusterStatusMap.Store(ip, true)
	} else {
		common.ClusterStatusMap.Store(ip, false)
	}

	conn.Release()

	if flag {
		update_tableinfo(TB_KUBE_CLUSTER_INFO, ontunetime)
	}

	return clusterid
}

// QueryUnusedClusterReset is to reset unused cluster and related data
func QueryUnusedClusterReset() {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	ip_array := make([]string, 0)
	common.ClusterStatusMap.Range(func(key, value any) bool {
		ip_array_mark := fmt.Sprintf("'%s'", key.(string))
		ip_array = append(ip_array, ip_array_mark)

		return true
	})
	ip_array_str := strings.Join(ip_array, ",")

	common.LogManager.Debug("Select Unused Cluster Info Previous Execution")
	common.LogManager.Debug(fmt.Sprintf(SELECT_UNUSED_CLUSTER_INFO, ip_array_str))
	rows, err := conn.Query(context.Background(), fmt.Sprintf(SELECT_UNUSED_CLUSTER_INFO, ip_array_str))
	errorCheck(err)
	common.LogManager.Debug("Select Unused Cluster Info Execution is completed")

	clusterid_arr := make([]int, 0)
	for rows.Next() {
		var clusterid int
		err := rows.Scan(&clusterid)
		errorCheck(err)

		clusterid_arr = append(clusterid_arr, clusterid)
	}

	rows.Close()

	ontunetime, _ := GetOntuneTime()

	for _, clusterid := range clusterid_arr {
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		common.LogManager.Debug(fmt.Sprintf("Cluster %d Reset Previous Execution", clusterid))

		_, update_err := tx.Exec(context.Background(), UPDATE_CLUSTER_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		_, update_err = tx.Exec(context.Background(), UPDATE_CLUSTER_NODE_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		_, update_err = tx.Exec(context.Background(), UPDATE_CLUSTER_POD_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		_, update_err = tx.Exec(context.Background(), UPDATE_CLUSTER_CONTAINER_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		_, update_err = tx.Exec(context.Background(), UPDATE_CLUSTER_NS_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		_, update_err = tx.Exec(context.Background(), UPDATE_CLUSTER_DEPLOY_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		_, update_err = tx.Exec(context.Background(), UPDATE_CLUSTER_STS_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		_, update_err = tx.Exec(context.Background(), UPDATE_CLUSTER_DS_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		_, update_err = tx.Exec(context.Background(), UPDATE_CLUSTER_RS_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		_, update_err = tx.Exec(context.Background(), UPDATE_CLUSTER_SVC_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		_, update_err = tx.Exec(context.Background(), UPDATE_CLUSTER_ING_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		_, update_err = tx.Exec(context.Background(), UPDATE_CLUSTER_PVC_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		_, update_err = tx.Exec(context.Background(), UPDATE_CLUSTER_PV_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		_, update_err = tx.Exec(context.Background(), UPDATE_CLUSTER_SC_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		_, update_err = tx.Exec(context.Background(), UPDATE_CLUSTER_INGHOST_RESET, clusterid, ontunetime)
		errorCheck(update_err)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}
