package database_test

import (
	"context"
	"fmt"
	"onTuneKubeManager/logger"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var viewList = []string{
	"kubeinginfo_v",
	"kubeingpodinfo_v",
	"kubesvcinfo_v",
	"kubesvcpodinfo_v",
}

func RemoveDir(filePath string) error {
	dir := strings.Replace(filePath, "\\", "/", -1)
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}

	return nil
}

func TestConnection(t *testing.T) {
	_, err := dbConnection()
	if err != nil {
		t.Errorf("DB Connection Fail")
	}
}

func TestCreateTablespace(t *testing.T) {
	pool, err := dbConnectionPool()
	if err != nil {
		t.Errorf("DB Connection Fail")
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		t.Errorf("DB Connection Fail")
	}

	defer conn.Release()

	//short term tablespace Check and create
	var tablespacecount int
	tablespacename := "test_tablespace"
	SQL := "select count(*) from pg_tablespace where spcname='" + tablespacename + "'"
	rows, err := conn.Query(context.Background(), SQL)
	if err != nil {
		t.Errorf("Query Fail")
	}

	for rows.Next() {
		err := rows.Scan(&tablespacecount)
		if err != nil {
			t.Errorf("Scan Fail")
		}
	}

	rows.Close()

	path := "D:\\"

	err = logger.CheckOrCreateFilePathDir(path+tablespacename+"/", "")
	if err != nil {
		t.Errorf("Create Tablespace Path Fail")
	}

	SQL = "create tablespace " + tablespacename + " location '" + path + tablespacename + "'"
	_, err = conn.Exec(context.Background(), SQL)
	if err != nil {
		t.Errorf("Create Tablespace Fail")
	}

	SQL = "select count(*) from pg_tablespace where spcname='" + tablespacename + "'"
	rows, err = conn.Query(context.Background(), SQL)
	if err != nil {
		t.Errorf("Query Fail")
	}

	for rows.Next() {
		err := rows.Scan(&tablespacecount)
		if err != nil {
			t.Errorf("Scan Fail")
		}
	}

	assert.Equal(t, 1, tablespacecount, "tablespace count is not 1")

	rows.Close()

	SQL = "drop tablespace " + tablespacename
	_, err = conn.Exec(context.Background(), SQL)
	if err != nil {
		t.Errorf("Drop Tablespace Fail")
	}

	err = RemoveDir(path + tablespacename + "/")
	if err != nil {
		t.Errorf("Remove Tablespace Path Fail")
	}
}

func TestTimestamp(t *testing.T) {
	pool, err := dbConnectionPool()
	if err != nil {
		t.Errorf("DB Connection Fail")
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		t.Errorf("DB Connection Fail")
	}

	defer conn.Release()

	SQL := "select round(extract(epoch from now()))"

	rows, err := conn.Query(context.Background(), SQL)
	if err != nil {
		t.Errorf("Query Fail")
	}

	var timestamp int
	for rows.Next() {
		err := rows.Scan(&timestamp)
		if err != nil {
			t.Errorf("Scan Fail")
		}
	}

	rows.Close()

	fmt.Println(timestamp, time.Now().Unix())
}

func TestCheckView(t *testing.T) {
	pool, err := dbConnectionPool()
	if err != nil {
		t.Errorf("DB Connection Fail")
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		t.Errorf("DB Connection Fail")
	}

	defer conn.Release()

	SQL := "select count(*) from pg_views where viewname like 'kube_%'"
	rows, err := conn.Query(context.Background(), SQL)
	if err != nil {
		t.Errorf("Query Fail")
	}

	var viewcount int
	for rows.Next() {
		err := rows.Scan(&viewcount)
		if err != nil {
			t.Errorf("Scan Fail")
		}
	}

	rows.Close()

	assert.Equal(t, len(viewList), viewcount, fmt.Sprintf("view count is not %d", len(viewList)))
}
