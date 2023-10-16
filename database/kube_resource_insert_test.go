package database_test

import (
	"context"
	"testing"

	"onTuneKubeManager/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestResourceRowEnabled(t *testing.T) {
	// Connection Phase
	pool, err := dbConnectionPool()
	if err != nil {
		t.Errorf("DB Connection Pool Fail")
	}

	uid, err := row_enabled_check_result(pool, "nsuid", database.TB_KUBE_NS_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("namespace uid: %s", uid)

	uid, err = row_enabled_check_result(pool, "nodeuid", database.TB_KUBE_NODE_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("node uid: %s", uid)

	uid, err = row_enabled_check_result(pool, "uid", database.TB_KUBE_POD_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("pod uid: %s", uid)

	containername, err := row_enabled_check_result(pool, "containername", database.TB_KUBE_CONTAINER_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("containername: %s", containername)

	uid, err = row_enabled_check_result(pool, "uid", database.TB_KUBE_SC_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("sc uid: %s", uid)

	uid, err = row_enabled_check_result(pool, "uid", database.TB_KUBE_SVC_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("svc uid: %s", uid)

	uid, err = row_enabled_check_result(pool, "uid", database.TB_KUBE_ING_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("ing uid: %s", uid)

	uid, err = row_enabled_check_result(pool, "inguid", database.TB_KUBE_INGHOST_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("inghost's ingress uid: %s", uid)

	uid, err = row_enabled_check_result(pool, "uid", database.TB_KUBE_DEPLOY_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("deploy uid: %s", uid)

	uid, err = row_enabled_check_result(pool, "uid", database.TB_KUBE_STS_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("sts uid: %s", uid)

	uid, err = row_enabled_check_result(pool, "uid", database.TB_KUBE_DS_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("ds uid: %s", uid)

	uid, err = row_enabled_check_result(pool, "uid", database.TB_KUBE_RS_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("rs uid: %s", uid)

	uid, err = row_enabled_check_result(pool, "pvuid", database.TB_KUBE_PV_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("pv uid: %s", uid)

	uid, err = row_enabled_check_result(pool, "uid", database.TB_KUBE_PVC_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("pvc uid: %s", uid)

	uid, err = row_enabled_check_result(pool, "uid", database.TB_KUBE_EVENT_INFO)
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
	t.Logf("event uid: %s", uid)
}

func row_enabled_check_result(pool *pgxpool.Pool, colName string, tablename string) (string, error) {
	rows, err := selectRowEnabled(pool, colName, tablename)
	if err != nil {
		return "", err
	}

	col, err := rows_scan(rows)
	if err != nil {
		return "", err
	}

	return col, nil
}

func rows_scan(rows pgx.Rows) (string, error) {
	for rows.Next() {
		var col string
		err := rows.Scan(&col)
		if err != nil {
			return col, err
		}
	}

	return "", nil
}

func selectRowEnabled(pool *pgxpool.Pool, colName string, tablename string) (pgx.Rows, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), "select "+colName+" from "+tablename+" where enabled=1")
	if err != nil {
		return nil, err
	}

	return rows, nil
}
