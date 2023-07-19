package database

import (
	"context"
	"fmt"
	"onTuneKubeManager/common"
	"onTuneKubeManager/config"
	"onTuneKubeManager/logger"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"
)

// MakeDBConn is make database connection
// db: database info
func MakeDBConn(db *config.DBConfig) {
	config, err := pgxpool.ParseConfig(db.GetDsn())
	if err != nil {
		common.LogManager.WriteLog(fmt.Sprintf("DB Pool Configuration Error - %v", err.Error()))
		panic(err)
	}

	config.MaxConns = int32(db.GetMaxConnection())
	common.DBConnectionPool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		common.LogManager.WriteLog(fmt.Sprintf("DB Pool Error - %v", err.Error()))
		panic(err)
	}

	resourceTableCheck(db)
	ShorttermTableCheck(db)

	//InitVacuum
	if common.InitVacuum {
		setInitVacuum()
	}
}

// setInitVacuum is set init vacuum
func setInitVacuum() {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	SQL := "select tablename from " + TB_KUBE_TABLE_INFO + " where durationmin > 0 and tablename like '%kubeavg%' "

	rows, err := conn.Query(context.Background(), SQL)
	errorCheck(err)

	tablename_arr := make([]string, 0)
	for rows.Next() {
		var tablename string
		err := rows.Scan(&tablename)
		errorCheck(err)

		tablename_arr = append(tablename_arr, tablename)
	}

	rows.Close()

	for _, tablename := range tablename_arr {
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		SQL = "alter table " + tablename + " set (autovacuum_enabled=true)"
		_, err = tx.Exec(context.Background(), SQL)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

// setAutoVacuum is set auto vacuum
func setAutoVacuum(tx pgx.Tx, tablename string) {
	var SQL string
	if common.AutoVacuum {
		SQL = "alter table " + tablename + " set (autovacuum_enabled=true)"
	} else {
		SQL = "alter table " + tablename + " set (autovacuum_enabled=false)"
	}

	_, err := tx.Exec(context.Background(), SQL)
	errorCheck(err)
}

// getTableSpaceStmt is get table space statement
func getTableSpaceStmt(table_type int, unixtime int64) string {
	if table_type == RESOURCE_TABLE && common.TableSpace {
		return "tablespace " + getTableSpace(RESOURCE_TABLE, unixtime)
	} else if table_type == SHORTERM_TABLE && common.ShorttermTableSpace {
		return "tablespace " + getTableSpace(SHORTERM_TABLE, unixtime)
	} else if table_type == LONGTERM_TABLE && common.LongtermTableSpace {
		return "tablespace " + getTableSpace(LONGTERM_TABLE, unixtime)
	}

	return ""
}

// getTableSpace is make table space
func getTableSpace(tableType int, unixtime int64) string {
	var returnVal string
	if tableType == RESOURCE_TABLE {
		if common.TableSpaceName != "" {
			returnVal = common.TableSpaceName
		}
	} else if tableType == SHORTERM_TABLE {
		if common.ShorttermTableSpaceName != "" {
			returnVal = common.ShorttermTableSpaceName
		}
	} else if tableType == LONGTERM_TABLE {
		if common.LongtermTableSpaceName != "" {
			returnVal = common.LongtermTableSpaceName
		}
	}

	return returnVal
}

// ShorttermTableCheck checks metric tables
// db: database info
func ShorttermTableCheck(db *config.DBConfig) {
	createTablespace(SHORTERM_TABLE, db.GetDBDocker(), db.GetDBDockerpath())

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	ontunetime, _ := GetOntuneTime()
	//realtime node perf table
	checkKubenodePerf(conn, ontunetime)
	//realtime pod perf table
	checkKubepodPerf(conn, ontunetime)
	//realtime container perf table
	checkKubecontainerPerf(conn, ontunetime)
	//realtime node network perf table
	checkKubeNodeNetPerf(conn, ontunetime)
	//realtime pod network perf table
	checkKubePodNetPerf(conn, ontunetime)
	//realtime node filesystem perf table
	checkKubeNodeFsPerf(conn, ontunetime)
	//realtime container filesystem perf table
	checkKubePodFsPerf(conn, ontunetime)
	//realtime container filesystem perf table
	checkKubeContainerFsPerf(conn, ontunetime)
}

// DailyPerfTable checks metric tables daily and creates new tables
func DailyPerfTable() {
	var bTableCreate bool
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	for {
		hour := time.Now().Hour()
		if hour >= 23 && !bTableCreate {
			bTableCreate = true
			ontunetime, _ := GetOntuneTime()
			//realtime node perf table
			dailyKubenodePerf(conn, ontunetime)
			//realtime pod perf table
			dailyKubepodPerf(conn, ontunetime)
			//realtime container perf table
			dailyKubecontainerPerf(conn, ontunetime)
			//realtime node network perf table
			dailyKubeNodeNetPerf(conn, ontunetime)
			//realtime pod network perf table
			dailyKubePodNetPerf(conn, ontunetime)
			//realtime node filesystem perf table
			dailyKubeNodeFsPerf(conn, ontunetime)
			//realtime pod filesystem perf table
			dailyKubePodFsPerf(conn, ontunetime)
			//realtime container filesystem perf table
			dailyKubeContainerFsPerf(conn, ontunetime)
		} else if hour != 23 {
			bTableCreate = false
		}
		time.Sleep(time.Second * 10)
	}
}

// DropTable drops metric tables
func DropTable() {
	var SQL string
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	for {
		SQL = "select tablename from " + TB_KUBE_TABLE_INFO + " where tablename like '%kuberealtime%' and createdtime < (select _time from " + TB_ONTUNE_INFO +
			") - " + strconv.Itoa(UNIXTIME_TO_DAY) + "*" + strconv.Itoa(common.ShorttermDuration)
		rows, err := conn.Query(context.Background(), SQL)
		errorCheck(err)

		var tablenameArr []string = make([]string, 0)
		for rows.Next() {
			var tablename string
			err := rows.Scan(&tablename)
			errorCheck(err)

			tablenameArr = append(tablenameArr, tablename)
		}

		rows.Close()

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		for _, t := range tablenameArr {
			SQL = "drop table " + t
			_, err = tx.Exec(context.Background(), SQL)
			errorCheck(err)

			SQL = "delete from " + TB_KUBE_TABLE_INFO + " where tablename = '" + t + "'"
			_, err = tx.Exec(context.Background(), SQL)
			errorCheck(err)
		}

		err = tx.Commit(context.Background())
		errorCheck(err)

		time.Sleep(time.Minute * 1)
	}
}

// createTablespace creates table space
func createTablespace(table_type int, docker_flag bool, docker_path string) {
	var err error
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	var path string
	var tablespace_path string

	if table_type == RESOURCE_TABLE {
		path = common.TableSpacePath

		if !common.TableSpace {
			return
		}
	} else if table_type == SHORTERM_TABLE {
		path = common.ShorttermTableSpacePath

		if !common.ShorttermTableSpace {
			return
		}
	} else if table_type == LONGTERM_TABLE {
		path = common.LongtermTableSpacePath

		if !common.LongtermTableSpace {
			return
		}
	}

	if docker_flag {
		tablespace_path = docker_path
	} else {
		tablespace_path = path
	}

	var tablespacecount int
	ontunetime, _ := GetOntuneTime()
	tablespacename := getTableSpace(table_type, ontunetime)
	SQL := "select count(*) from pg_tablespace where spcname='" + tablespacename + "'"
	rows, err := conn.Query(context.Background(), SQL)
	errorCheck(err)
	for rows.Next() {
		err := rows.Scan(&tablespacecount)
		errorCheck(err)
	}
	common.LogManager.Debug(fmt.Sprintf("%v %v", SQL, rows))

	rows.Close()

	if tablespacecount == 0 {
		err := logger.CheckOrCreateFilePathDir(path+tablespacename+"/", common.DbOsUser)
		errorCheck(err)
		common.LogManager.WriteLog(fmt.Sprintf("Create Directory %v", path+tablespacename))

		SQL = "create tablespace " + tablespacename + " location '" + tablespace_path + tablespacename + "'"
		_, err = conn.Exec(context.Background(), SQL)
		errorCheck(err)
		common.LogManager.WriteLog(fmt.Sprintf("Create Tablespace %v", tablespace_path+tablespacename))
	}
}

// resourceTableCheck checks resource tables
// db: database info
func resourceTableCheck(db *config.DBConfig) {
	createTablespace(RESOURCE_TABLE, db.GetDBDocker(), db.GetDBDockerpath())

	var err error
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorCheck(err)
	}

	defer conn.Release()

	//kubetableinfo Create
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	errorCheck(err)

	_, err = tx.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS  `+TB_KUBE_TABLE_INFO+` (
			tablename     varchar(64) NOT NULL PRIMARY KEY,
			version       integer NOT NULL,
			createdtime   bigint NOT NULL,
			updatetime    bigint NOT NULL,
			durationmin   integer NULL
		) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
	errorCheck(err)

	createIndex(tx, TB_KUBE_TABLE_INFO, "createdtime, updatetime")

	err = tx.Commit(context.Background())
	errorCheck(err)

	ontunetime, _ := GetOntuneTime()

	//kubemanager info
	checkKubemanagerinfo(conn, ontunetime)
	//cluster info
	checkKubeclusterinfo(conn, ontunetime)
	//resource info
	checkKuberesourceinfo(conn, ontunetime)
	//sc info
	checkKubescinfo(conn, ontunetime)
	//ns info
	checkKubensinfo(conn, ontunetime)
	//node info
	checkKubenodeinfo(conn, ontunetime)
	//pod info
	checkKubepodinfo(conn, ontunetime)
	//container info
	checkKubecontainerinfo(conn, ontunetime)
	//fsdevice info
	checkKubefsdeviceinfo(conn, ontunetime)
	//netinterface info
	checkKubenetinterfaceinfo(conn, ontunetime)
	//metricid info
	checkKubemetricidinfo(conn, ontunetime)
	//service info
	checkKubesvcinfo(conn, ontunetime)
	//ing info
	checkKubeinginfo(conn, ontunetime)
	//inghost info
	checkKubeinghostinfo(conn, ontunetime)
	//deployment info
	checkKubedeployinfo(conn, ontunetime)
	//statefulSet info
	checkKubestsinfo(conn, ontunetime)
	//daemonSet info
	checkKubedsinfo(conn, ontunetime)
	//replicaSet info
	checkKubersinfo(conn, ontunetime)
	//pvc info
	checkKubepvcinfo(conn, ontunetime)
	//pv info
	checkKubepvinfo(conn, ontunetime)
	//event info
	checkKubeeventinfo(conn, ontunetime)
	//log info
	checkKubeloginfo(conn, ontunetime)
}

// GetDBConnString is get database connection string
func (db *DBInfo) GetDBConnString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		db.Host, db.User, db.Pwd, db.DBname, db.Port)
}
