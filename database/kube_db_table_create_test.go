package database_test

import (
	"context"
	"database/sql"
	"onTuneKubeManager/database"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func dbConnection() (*sql.DB, error) {
	connString := "host=127.0.0.1 user=ontune password=ontune dbname=kube_test_db port=5432 sslmode=disable TimeZone=Asia/Seoul"
	conn, err := sql.Open("postgres", connString)
	if err != nil || conn.Ping() != nil {
		return nil, err
	}

	return conn, nil
}

func dbConnectionPool() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig("postgres://ontune:ontune@127.0.0.1:5432/kube_test_db?sslmode=disable")
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(100)

	return pgxpool.NewWithConfig(context.Background(), config)
}

func ontuneTableExists(conn *pgxpool.Conn, tablename string) bool {
	// Check Table Exist using select table query or pg_tables
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return false
	}

	SQL := "select tablename from pg_tables where tablename = '" + tablename + "'"
	rows, err := tx.Query(context.Background(), SQL)
	if err != nil {
		return false
	}

	rows.Close()

	defer rows.Close()

	return true
}

func tableExists(conn *pgxpool.Conn, tablename string) bool {
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return false
	}
	SQL := "select tablename from " + database.TB_KUBE_TABLE_INFO + " where tablename = '" + tablename + "'"
	rows, err := tx.Query(context.Background(), SQL)
	if err != nil {
		return false
	}

	rows.Close()

	// Check Table Exist using select table query or pg_tables
	SQL = "select tablename from pg_tables where tablename = '" + tablename + "'"
	rows, err = tx.Query(context.Background(), SQL)
	if err != nil {
		return false
	}

	rows.Close()

	return true
}

func viewExists(conn *pgxpool.Conn, viewname string) bool {
	// Check Table Exist using select table query or pg_tables
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return false
	}

	SQL := "select viewname from pg_views where viewname = '" + viewname + "'"
	rows, err := tx.Query(context.Background(), SQL)
	if err != nil {
		return false
	}

	defer rows.Close()

	return true
}

func getTableName(tablename string) string {
	now := time.Now()
	createDay := now.Format(database.DATE_FORMAT)

	return tablename + "_" + createDay
}

func TestOntuneTableCreation(t *testing.T) {
	pool, err := dbConnectionPool()
	if err != nil {
		t.Errorf("DB Connection Pool Fail")
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		t.Errorf("DB Connection Fail")
	}

	defer conn.Release()

	if !ontuneTableExists(conn, database.TB_ONTUNE_INFO) {
		t.Errorf("Table %s does not exist", database.TB_ONTUNE_INFO)
	}
}

func TestInfoTableCreation(t *testing.T) {
	pool, err := dbConnectionPool()
	if err != nil {
		t.Errorf("DB Connection Pool Fail")
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		t.Errorf("DB Connection Fail")
	}

	defer conn.Release()

	if !tableExists(conn, database.TB_KUBE_MANAGER_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_MANAGER_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_CLUSTER_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_CLUSTER_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_RESOURCE_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_RESOURCE_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_SC_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_SC_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_NS_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_NS_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_NODE_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_NODE_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_POD_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_POD_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_CONTAINER_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_CONTAINER_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_ING_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_ING_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_INGHOST_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_INGHOST_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_SVC_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_SVC_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_DEPLOY_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_DEPLOY_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_STS_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_STS_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_DS_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_DS_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_RS_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_RS_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_PV_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_PV_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_PVC_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_PVC_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_EVENT_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_EVENT_INFO)
	}
}

func TestInfoViewCreation(t *testing.T) {
	pool, err := dbConnectionPool()
	if err != nil {
		t.Errorf("DB Connection Pool Fail")
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		t.Errorf("DB Connection Fail")
	}

	defer conn.Release()

	if !viewExists(conn, database.TB_KUBE_SVC_INFO_V) {
		t.Errorf("View %s does not exist", database.TB_KUBE_SVC_INFO_V)
	}

	if !viewExists(conn, database.TB_KUBE_ING_INFO_V) {
		t.Errorf("View %s does not exist", database.TB_KUBE_ING_INFO_V)
	}

	if !viewExists(conn, database.TB_KUBE_SVCPOD_INFO_V) {
		t.Errorf("View %s does not exist", database.TB_KUBE_SVCPOD_INFO_V)
	}

	if !viewExists(conn, database.TB_KUBE_INGPOD_INFO_V) {
		t.Errorf("View %s does not exist", database.TB_KUBE_INGPOD_INFO_V)
	}
}

func TestMetricTableCreation(t *testing.T) {
	pool, err := dbConnectionPool()
	if err != nil {
		t.Errorf("DB Connection Pool Fail")
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		t.Errorf("DB Connection Fail")
	}

	defer conn.Release()

	if !tableExists(conn, database.TB_KUBE_LAST_NODE_PERF) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_NODE_PERF)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_POD_PERF) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_POD_PERF)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_CONTAINER_PERF) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_CONTAINER_PERF)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_NODE_NET_PERF) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_NODE_NET_PERF)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_POD_NET_PERF) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_POD_NET_PERF)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_NODE_FS_PERF) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_NODE_FS_PERF)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_POD_FS_PERF) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_POD_FS_PERF)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_CONTAINER_FS_PERF) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_CONTAINER_FS_PERF)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_CLUSTER_PERF) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_CLUSTER_PERF)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_NAMESPACE_PERF) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_NAMESPACE_PERF)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_REPLICASET_PERF) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_REPLICASET_PERF)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_DAEMONSET_PERF) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_DAEMONSET_PERF)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_STATEFULSET_PERF) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_STATEFULSET_PERF)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_DEPLOYMENT_PERF) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_DEPLOYMENT_PERF)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_NODE_PERF_RAW) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_NODE_PERF_RAW)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_POD_PERF_RAW) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_POD_PERF_RAW)
	}

	if !tableExists(conn, database.TB_KUBE_LAST_CONTAINER_PERF_RAW) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_LAST_CONTAINER_PERF_RAW)
	}

	if !tableExists(conn, database.TB_KUBE_FS_DEVICE_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_FS_DEVICE_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_NET_INTERFACE_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_NET_INTERFACE_INFO)
	}

	if !tableExists(conn, database.TB_KUBE_METRIC_ID_INFO) {
		t.Errorf("Table %s does not exist", database.TB_KUBE_METRIC_ID_INFO)
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_NODE_PERF)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_NODE_PERF))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_POD_PERF)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_POD_PERF))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_CONTAINER_PERF)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_CONTAINER_PERF))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_NODE_NET_PERF)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_NODE_NET_PERF))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_POD_NET_PERF)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_POD_NET_PERF))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_NODE_FS_PERF)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_NODE_FS_PERF))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_POD_FS_PERF)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_POD_FS_PERF))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_CONTAINER_FS_PERF)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_CONTAINER_FS_PERF))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_CLUSTER_PERF)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_CLUSTER_PERF))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_NAMESPACE_PERF)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_NAMESPACE_PERF))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_REPLICASET_PERF)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_REPLICASET_PERF))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_DAEMONSET_PERF)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_DAEMONSET_PERF))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_STATEFULSET_PERF)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_STATEFULSET_PERF))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_DEPLOYMENT_PERF)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_DEPLOYMENT_PERF))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_NODE_PERF_RAW)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_NODE_PERF_RAW))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_POD_PERF_RAW)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_POD_PERF_RAW))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_CONTAINER_PERF_RAW)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_CONTAINER_PERF_RAW))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_NODE_NET_PERF_RAW)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_NODE_NET_PERF_RAW))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_POD_NET_PERF_RAW)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_POD_NET_PERF_RAW))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_NODE_FS_PERF_RAW)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_NODE_FS_PERF_RAW))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_POD_FS_PERF_RAW)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_POD_FS_PERF_RAW))
	}

	if !tableExists(conn, getTableName(database.TB_KUBE_CONTAINER_FS_PERF_RAW)) {
		t.Errorf("Table %s does not exist", getTableName(database.TB_KUBE_CONTAINER_FS_PERF_RAW))
	}
}
