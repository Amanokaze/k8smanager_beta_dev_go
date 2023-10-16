package database

import (
	"context"
	"fmt"
	"onTuneKubeManager/common"
	"onTuneKubeManager/config"
	"onTuneKubeManager/logger"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

// CreateConnectionPool is create connection pool
func CreateConnectionPool(db *config.DBConfig) {
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
}

// MakeDBConn is make database connection
// db: database info
func MakeDBConn(db *config.DBConfig) {
	CreateConnectionPool(db)

	resourceTableCheck(db)
	ShorttermTableCheck(db)
	LongtermTableCheck(db)

	//InitVacuum
	if common.InitVacuum {
		setInitVacuum()
	}
}

// setInitVacuum is set init vacuum
func setInitVacuum() {
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	SQL := "select tablename from " + TB_KUBE_TABLE_INFO + " where durationmin > 0 and tablename like '%kubeavg%' "

	rows, err := conn.Query(context.Background(), SQL)
	if !errorCheck(err) {
		return
	}

	tablename_arr := make([]string, 0)
	for rows.Next() {
		var tablename string
		err := rows.Scan(&tablename)
		if !errorCheck(err) {
			return
		}

		tablename_arr = append(tablename_arr, tablename)
	}

	rows.Close()

	for _, tablename := range tablename_arr {
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
			return
		}

		SQL = "alter table " + tablename + " set (autovacuum_enabled=true)"
		_, err = tx.Exec(context.Background(), SQL)
		if !errorCheck(err) {
			return
		}

		err = tx.Commit(context.Background())
		if !errorCheck(errors.Wrap(err, "Commit error")) {
			return
		}
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
	if !errorCheck(err) {
		return
	}
}

// getTableSpaceStmt is get table space statement
func getTableSpaceStmt(table_type int, unixtime int64) string {
	if table_type == RESOURCE_TABLE && common.TableSpace {
		return "tablespace " + getTableSpace(RESOURCE_TABLE, unixtime)
	} else if table_type == SHORTTERM_TABLE && common.ShorttermTableSpace {
		return "tablespace " + getTableSpace(SHORTTERM_TABLE, unixtime)
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
	} else if tableType == SHORTTERM_TABLE {
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
	createTablespace(SHORTTERM_TABLE, db.GetDBDocker(), db.GetDBDockerpath())

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

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

// LongtermTableCheck checks metric tables
// db: database info
func LongtermTableCheck(db *config.DBConfig) {
	createTablespace(LONGTERM_TABLE, db.GetDBDocker(), db.GetDBDockerpath())

	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	//avg perf table
	checkKubeAvgPerf(conn, ontunetime)
}

// DailyPerfTable checks metric tables daily and creates new tables
func DailyPerfTable() {
	var bTableCreate bool
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	for {
		hour := time.Now().Hour()
		if hour >= 23 && !bTableCreate {
			bTableCreate = true
			ontunetime, _ := GetOntuneTime()
			if ontunetime == 0 {
				return
			}

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

// DropTableSelect query to drop metric tables
func DropTableSelect(conn *pgxpool.Conn, tabletype int) []string {
	var duration int
	var tableprefix string

	if tabletype == SHORTTERM_TABLE {
		duration = common.ShorttermDuration
		tableprefix = "kuberealtime"
	} else if tabletype == LONGTERM_TABLE {
		duration = common.LongtermDuration
		tableprefix = "kubeavg"
	} else {
		return nil
	}

	SQL := fmt.Sprintf(SELECT_DROP_TABLE, TB_KUBE_TABLE_INFO, tableprefix, TB_ONTUNE_INFO, UNIXTIME_TO_DAY, duration)
	rows, err := conn.Query(context.Background(), SQL)
	if !errorCheck(err) {
		return nil
	}

	var tablenameArr []string = make([]string, 0)
	for rows.Next() {
		var tablename string
		err := rows.Scan(&tablename)
		if !errorCheck(err) {
			return nil
		}

		tablenameArr = append(tablenameArr, tablename)
	}

	rows.Close()

	return tablenameArr
}

// DropTable drops metric tables
func DropTable() {
	for {
		conn, err := common.DBConnectionPool.Acquire(context.Background())
		if conn == nil || err != nil {
			errorDisconnect(errors.Wrap(err, "Acquire connection error"))
			return
		}

		var tablenameArr []string = make([]string, 0)
		tablenameArr = append(tablenameArr, DropTableSelect(conn, SHORTTERM_TABLE)...)
		tablenameArr = append(tablenameArr, DropTableSelect(conn, LONGTERM_TABLE)...)

		if len(tablenameArr) > 0 {
			tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
			if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
				return
			}

			for _, t := range tablenameArr {
				SQL := "drop table " + t
				_, err = tx.Exec(context.Background(), SQL)
				if !errorCheck(err) {
					return
				}

				SQL = "delete from " + TB_KUBE_TABLE_INFO + " where tablename = '" + t + "'"
				_, err = tx.Exec(context.Background(), SQL)
				if !errorCheck(err) {
					return
				}
			}

			err = tx.Commit(context.Background())
			if !errorCheck(errors.Wrap(err, "Commit error")) {
				return
			}
		}

		conn.Release()

		time.Sleep(time.Minute * 1)
	}
}

// createTablespace creates table space
func createTablespace(table_type int, docker_flag bool, docker_path string) {
	var err error
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	var path string
	var tablespace_path string

	if table_type == RESOURCE_TABLE {
		path = common.TableSpacePath

		if !common.TableSpace {
			return
		}
	} else if table_type == SHORTTERM_TABLE {
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
	if ontunetime == 0 {
		return
	}

	tablespacename := getTableSpace(table_type, ontunetime)
	SQL := "select count(*) from pg_tablespace where spcname='" + tablespacename + "'"
	rows, err := conn.Query(context.Background(), SQL)
	if !errorCheck(err) {
		return
	}

	for rows.Next() {
		err := rows.Scan(&tablespacecount)
		if !errorCheck(err) {
			return
		}
	}
	common.LogManager.Debug(fmt.Sprintf("%v %v", SQL, rows))

	rows.Close()

	if tablespacecount == 0 {
		err := logger.CheckOrCreateFilePathDir(tablespace_path+tablespacename+"/", common.DbOsUser)
		if !errorCheck(err) {
			return
		}
		common.LogManager.WriteLog(fmt.Sprintf("Create Directory %v", tablespace_path+tablespacename))

		SQL = "create tablespace " + tablespacename + " location '" + path + tablespacename + "'"
		_, err = conn.Exec(context.Background(), SQL)
		if !errorCheck(err) {
			return
		}
		common.LogManager.WriteLog(fmt.Sprintf("Create Tablespace %v", path+tablespacename))
	}
}

// resourceTableCheck checks resource tables
// db: database info
func resourceTableCheck(db *config.DBConfig) {
	createTablespace(RESOURCE_TABLE, db.GetDBDocker(), db.GetDBDockerpath())

	var err error
	conn, err := common.DBConnectionPool.Acquire(context.Background())
	if conn == nil || err != nil {
		errorDisconnect(errors.Wrap(err, "Acquire connection error"))
		return
	}

	defer conn.Release()

	//kubetableinfo Create
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if !errorCheck(errors.Wrap(err, "Begin transaction error")) {
		return
	}

	_, err = tx.Exec(context.Background(), getCreateResourceTableStmt(TB_KUBE_TABLE_INFO)+getTableSpaceStmt(RESOURCE_TABLE, 0))
	if !errorCheck(err) {
		return
	}

	createIndex(tx, TB_KUBE_TABLE_INFO, "updatetime", "tablename, updatetime")

	err = tx.Commit(context.Background())
	if !errorCheck(errors.Wrap(err, "Commit error")) {
		return
	}

	ontunetime, _ := GetOntuneTime()
	if ontunetime == 0 {
		return
	}

	//kubemanager info
	createResourceTable(conn, ontunetime, TB_KUBE_MANAGER_INFO, false, false)
	//cluster info
	createResourceTable(conn, ontunetime, TB_KUBE_CLUSTER_INFO, false, false)
	//resource info
	createResourceTable(conn, ontunetime, TB_KUBE_RESOURCE_INFO, false, false)

	//sc info
	createResourceTable(conn, ontunetime, TB_KUBE_SC_INFO, true, true)
	//ns info
	createResourceTable(conn, ontunetime, TB_KUBE_NS_INFO, true, true)
	//node info
	createResourceTable(conn, ontunetime, TB_KUBE_NODE_INFO, true, true)
	//pod info
	createResourceTable(conn, ontunetime, TB_KUBE_POD_INFO, true, true)
	//container info
	createResourceTable(conn, ontunetime, TB_KUBE_CONTAINER_INFO, true, true)

	//service info
	svc_creation_flag := createResourceTable(conn, ontunetime, TB_KUBE_SVC_INFO, true, true)
	createResourceView(conn, ontunetime, TB_KUBE_SVC_INFO, CREATE_VIEW_SVC_INFO, svc_creation_flag)
	createResourceView(conn, ontunetime, TB_KUBE_SVC_INFO, CREATE_VIEW_SVCPOD_INFO, svc_creation_flag)

	//pvc info
	createResourceTable(conn, ontunetime, TB_KUBE_PVC_INFO, true, true)
	//pv info
	createResourceTable(conn, ontunetime, TB_KUBE_PV_INFO, true, true)
	//event info
	createResourceTable(conn, ontunetime, TB_KUBE_EVENT_INFO, true, true)
	//log info
	createResourceTable(conn, ontunetime, TB_KUBE_LOG_INFO, false, false)

	//deployment info
	createResourceTable(conn, ontunetime, TB_KUBE_DEPLOY_INFO, true, true)
	//statefulSet info
	createResourceTable(conn, ontunetime, TB_KUBE_STS_INFO, true, true)
	//daemonSet info
	createResourceTable(conn, ontunetime, TB_KUBE_DS_INFO, true, true)
	//replicaSet info
	createResourceTable(conn, ontunetime, TB_KUBE_RS_INFO, true, true)

	//ing info
	createResourceTable(conn, ontunetime, TB_KUBE_ING_INFO, true, true)

	//inghost info
	inghost_creation_flag := createResourceTable(conn, ontunetime, TB_KUBE_INGHOST_INFO, true, true)
	createResourceView(conn, ontunetime, TB_KUBE_INGHOST_INFO, CREATE_VIEW_ING_INFO, inghost_creation_flag)
	createResourceView(conn, ontunetime, TB_KUBE_INGHOST_INFO, CREATE_VIEW_INGPOD_INFO, inghost_creation_flag)

	//fsdevice info
	createResourceTable(conn, ontunetime, TB_KUBE_FS_DEVICE_INFO, false, false)
	//netinterface info
	createResourceTable(conn, ontunetime, TB_KUBE_NET_INTERFACE_INFO, false, false)
	//metricid info
	createResourceTable(conn, ontunetime, TB_KUBE_METRIC_ID_INFO, false, false)
}

// Manager DB Check checks DB connection
func ManagerDBCheck() {
	for {
		<-common.ChannelManagerDBConnection

		for {
			err := common.DBConnectionPool.Ping(context.Background())
			if err != nil {
				common.LogManager.WriteLog(fmt.Sprintf("Manager DB is disconnected now - %v", err.Error()))
			} else {
				common.LogManager.WriteLog("Manager DB is re-connected now!")
				common.ManagerDBConnected = true

				break
			}

			time.Sleep(time.Second * time.Duration(1))
		}

		time.Sleep(time.Second * time.Duration(1))
	}
}

// GetDBConnString is get database connection string
func (db *DBInfo) GetDBConnString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		db.Host, db.User, db.Pwd, db.DBname, db.Port)
}
