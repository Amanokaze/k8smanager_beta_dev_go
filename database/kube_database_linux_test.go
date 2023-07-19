package database_test

import (
	"context"
	"onTuneKubeManager/logger"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func dbConnectionPoolLinux() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig("postgres://ontune:ontune@127.0.0.1:25432/ontune?sslmode=disable")
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(100)

	return pgxpool.NewWithConfig(context.Background(), config)
}

func TestCreateTablespaceLinux(t *testing.T) {
	pool, err := dbConnectionPoolLinux()
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

	path := "/var/lib/postgresql/data/"
	docker_path := "/workspace/docker/k8s-188/pgdata/"

	err = logger.CheckOrCreateFilePathDir(docker_path+tablespacename+"/", "polkitd")
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

	err = RemoveDir(docker_path + tablespacename + "/")
	if err != nil {
		t.Errorf("Remove Tablespace Path Fail")
	}
}
