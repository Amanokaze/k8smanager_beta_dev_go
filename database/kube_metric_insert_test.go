package database_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func GetConnection() (*pgxpool.Conn, error) {
	// Connection Phase
	pool, err := dbConnectionPool()
	if err != nil {
		return nil, err
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func CreateTable(tablename string, conn *pgxpool.Conn) error {
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}

	// Create Phase
	create_stmt := `CREATE TABLE IF NOT EXISTS ` + tablename + ` (
		ontunetime bigint,
		podcount int,
		cpuusage double precision
	)`

	_, err = tx.Exec(context.Background(), create_stmt)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func DropTable(tablename string, conn *pgxpool.Conn) error {
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}

	// Drop Phase
	drop_stmt := `DROP TABLE IF EXISTS ` + tablename

	_, err = tx.Exec(context.Background(), drop_stmt)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func GetInsertData(tablename string) (string, []interface{}) {
	insert_stmt := `INSERT INTO ` + tablename + ` (ontunetime, podcount, cpuusage)
	(select * from unnest($1::bigint[], $2::int[], $3::double precision[]))`

	data := make([]interface{}, 0)
	data = append(data, pq.Int64Array([]int64{1, 2, 3, 4, 5}))
	data = append(data, pq.Array([]int{1, 2, 3, 4, 5}))
	data = append(data, pq.Float64Array([]float64{1, 2, 3, 4, 5}))

	return insert_stmt, data
}

func TestMetricInsertIdleintransaction(t *testing.T) {
	conn, err := GetConnection()
	if err != nil {
		t.Errorf("DB Connection Fail")
	}

	defer conn.Release()

	tablename := "kube_test_realtimeperf_iiit"
	err = CreateTable(tablename, conn)
	if err != nil {
		t.Errorf("DB Create Table Fail")
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		t.Errorf("DB BeginTx Fail")
	}

	// Insert Phase
	insert_stmt, data := GetInsertData(tablename)

	// One Time
	_, err = tx.Exec(context.Background(), insert_stmt, data...)
	if err != nil {
		t.Errorf("DB Insertion Exec Fail - %s", err.Error())
	}

	// Two Time
	_, err = tx.Exec(context.Background(), insert_stmt, data...)
	if err != nil {
		t.Errorf("DB Insertion Exec Fail - %s", err.Error())
	}

	err = tx.Commit(context.Background())
	if err != nil {
		t.Errorf("DB Insertion Commit Fail")
	}

	// Select Phase
	SQL := "select state from pg_stat_activity where query like '%" + insert_stmt + "%'"

	rows, err := conn.Query(context.Background(), SQL)
	if err != nil {
		t.Errorf("Query Fail")
	}

	for rows.Next() {
		var state string
		err := rows.Scan(&state)
		if err != nil {
			t.Errorf("Scan Fail - %s", err.Error())
		}

		if state == "idle in transaction (aborted)" {
			t.Errorf("DB State Fail")
		}
	}

	rows.Close()

	// Drop Phase
	err = DropTable(tablename, conn)
	if err != nil {
		t.Errorf("DB Drop Table Fail")
	}
}

func TestMetricInsertLastrealtimeperfDuplication(t *testing.T) {
	conn, err := GetConnection()
	if err != nil {
		t.Errorf("DB Connection Fail")
	}

	defer conn.Release()

	tablename := "kube_test_lastrealtimeperf_duplication"
	err = CreateTable(tablename, conn)
	if err != nil {
		t.Errorf("DB Create Table Fail")
	}

	for i := 0; i < 10; i++ {
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if err != nil {
			t.Errorf("DB BeginTx Fail")
		}

		// Delete Phase
		delete_stmt := `DELETE FROM %s`
		_, del_err := tx.Exec(context.Background(), fmt.Sprintf(delete_stmt, tablename))
		if del_err != nil {
			t.Errorf("DB Delete Fail")
		}

		// Insert Phase
		insert_stmt, data := GetInsertData(tablename)
		_, err = tx.Exec(context.Background(), insert_stmt, data...)
		if err != nil {
			t.Errorf("DB Insertion Exec Fail - %s", err.Error())
		}

		err = tx.Commit(context.Background())
		if err != nil {
			t.Errorf("DB Insertion Commit Fail")
		}

		// Select Phase
		SQL := "select count(*) from " + tablename

		rows, err := conn.Query(context.Background(), SQL)
		if err != nil {
			t.Errorf("Query Fail")
		}

		for rows.Next() {
			var count int
			err := rows.Scan(&count)
			if err != nil {
				t.Errorf("Scan Fail - %s", err.Error())
			}

			assert.Equal(t, len(data[0].(pq.Int64Array)), count)
		}

		time.Sleep(time.Duration(10) * time.Millisecond)
	}

	// Drop Phase
	err = DropTable(tablename, conn)
	if err != nil {
		t.Errorf("DB Drop Table Fail")
	}
}
