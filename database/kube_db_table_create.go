package database

import (
	"context"
	"fmt"
	"onTuneKubeManager/common"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
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

func checkKubemanagerinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_MANAGER_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_MANAGER_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_MANAGER_INFO+` (
				managerid     serial NOT NULL PRIMARY KEY,
				managername   text NOT NULL,
				description   text NULL,
				ip    		  text NOT NULL,
				createdtime   bigint NOT NULL,
				updatetime    bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_MANAGER_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubeclusterinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_CLUSTER_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_CLUSTER_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_CLUSTER_INFO+` (
				clusterid     serial NOT NULL PRIMARY KEY,
				managerid     integer NOT NULL,
				clustername   text NOT NULL,
				context		  text NULL,
				ip    		  text NOT NULL,
				enabled		  integer NOT NULL DEFAULT 1,
				status		  integer NOT NULL DEFAULT 1,
				createdtime   bigint NOT NULL,
				updatetime    bigint NOT NULL				
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_CLUSTER_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKuberesourceinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_RESOURCE_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_RESOURCE_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_RESOURCE_INFO+` (
				resourceid    	serial NOT NULL PRIMARY KEY,
				clusterid     	integer NOT NULL,
				resourcename  	text NOT NULL,
				apiclass   		text NOT NULL,
				version    		text NOT NULL,
				endpoint    	text NOT NULL,
				enabled    		integer NOT NULL DEFAULT 1,
				createdtime   	bigint NOT NULL,
				updatetime    	bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_RESOURCE_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubescinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_SC_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_SC_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_SC_INFO+` (
				scid    		  	serial NOT NULL PRIMARY KEY,
				clusterid     		integer NOT NULL,
				scname        		text NOT NULL,
				uid   		    	text NOT NULL,
				starttime    		bigint NOT NULL,
				labels				text NULL,
				provisioner   		text NULL,
				reclaimpolicy 		text NULL,
				volumebindingmode 	text NULL,
				allowvolumeexp 		integer NULL,
				enabled    			integer NOT NULL DEFAULT 0,
				createdtime   		bigint NOT NULL,
				updatetime   		bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_SC_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_SC_INFO, "clusterid")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubensinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_NS_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_NS_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_NS_INFO+` (
				nsid    		  	serial NOT NULL PRIMARY KEY,
				nsuid        		text NOT NULL,
				clusterid     		integer NOT NULL,
				nsname        		text NOT NULL,
				starttime			bigint NOT NULL,
				labels				text NULL,
				status   		    text NULL,
				enabled    			integer NOT NULL DEFAULT 0,
				createdtime   		bigint NOT NULL,
				updatetime   		bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_NS_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_NS_INFO, "clusterid, nsname")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubenodeinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_NODE_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_NODE_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_NODE_INFO+` (
				nodeid    		  			serial NOT NULL PRIMARY KEY,
				managerid     				integer NOT NULL,
				clusterid     				integer NOT NULL,
				nodeuid     				text NOT NULL,
				nodename        			text NULL,
				nodenameext        			text NULL,
				nodetype        			text NULL,
				enabled    					integer NOT NULL DEFAULT 1,
				starttime    				bigint NOT NULL,
				labels						text NULL,
				kernelversion   			text NULL,
				osimage 					text NULL,
				osname 						text NULL,
				containerruntimever 		text NULL,
				kubeletver 					text NULL,
				kubeproxyver 				text NULL,
				cpuarch 					text NULL,
				cpucount 					integer NULL,
				ephemeralstorage   			bigint NOT NULL default 0,
				memorysize   				bigint NOT NULL default 0,
				pods   						bigint NOT NULL default 0,
				ip 							text NULL,
				status						integer NOT NULL default 0,
				createdtime   				bigint NOT NULL,
				updatetime   				bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_NODE_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_NODE_INFO, "clusterid")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubepodinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_POD_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_POD_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_POD_INFO+` (
				podid    		  		serial NOT NULL PRIMARY KEY,
				clusterid     	integer NOT NULL,
				uid     					text NOT NULL,
				nodeuid     			text NOT NULL,
				nsuid     				text NOT NULL,
				annotationuid   	text NOT NULL,
				podname        		text NULL,
				starttime    			bigint NOT NULL,
				labels				text NULL,
				selector			text NULL,
				restartpolicy   	text NULL,
				serviceaccount 		text NULL,
				status 						text NULL,
				hostip 						text NULL,
				podip 						text NULL,
				restartcount 			bigint NOT NULL default 0,
				restarttime 			bigint NOT NULL default 0,
				podcondition 			text NULL,
				staticpod   			text NULL,
				refkind   				text NULL,
				refuid   					text NULL,
				pvcuid					text NULL,
				enabled 					integer NULL DEFAULT 0,
				createdtime   		bigint NOT NULL,
				updatetime   			bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_POD_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_POD_INFO, "nodeuid, nsuid, uid, refkind, refuid")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubecontainerinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_CONTAINER_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_CONTAINER_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_CONTAINER_INFO+` (
				containerid    		serial NOT NULL PRIMARY KEY,
				clusterid     	integer NOT NULL,
				poduid     			text NOT NULL,
				containername     	text NOT NULL,
				image   		    text NOT NULL,
				ports   		    text NULL,
				env   		    	text NULL,
				limitcpu 			bigint null,
				limitmemory			bigint null,
				limitstorage		bigint null,
				limitephemeral		bigint null,
				reqcpu 				bigint null,
				reqmemory			bigint null,
				reqstorage			bigint null,
				reqephemeral		bigint null,
				volumemounts   		text NULL,
				state               text NULL,
				enabled    			integer NOT NULL DEFAULT 0,
				createdtime   		bigint NOT NULL,
				updatetime   		bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_CONTAINER_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_CONTAINER_INFO, "poduid")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubefsdeviceinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_FS_DEVICE_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_FS_DEVICE_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_FS_DEVICE_INFO+` (
				deviceid      serial NOT NULL PRIMARY KEY,
				devicename    text NOT NULL,
				createdtime   bigint NOT NULL,
				updatetime   	bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_FS_DEVICE_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubenetinterfaceinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_NET_INTERFACE_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_NET_INTERFACE_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_NET_INTERFACE_INFO+` (
				interfaceid      serial NOT NULL PRIMARY KEY,
				interfacename    text NOT NULL,
				createdtime   bigint NOT NULL,
				updatetime   	bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_NET_INTERFACE_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubemetricidinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_METRIC_ID_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_METRIC_ID_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_METRIC_ID_INFO+` (
				metricid      serial NOT NULL PRIMARY KEY,
				metricname    text NOT NULL,
				image		  text NOT NULL,
				createdtime   bigint NOT NULL,
				updatetime   	bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_METRIC_ID_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubeinginfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_ING_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_ING_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_ING_INFO+` (
				ingid     		serial NOT NULL PRIMARY KEY,
				clusterid     	integer NOT NULL,
				nsuid     		text NOT NULL,
				ingname   		text NOT NULL,
				uid   				text NOT NULL,
				starttime   	bigint NOT NULL,
				labels				text NULL,
				classname   	text NULL,
				enabled    		integer NOT NULL DEFAULT 0,
				createdtime   bigint NOT NULL,
				updatetime    bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_ING_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_ING_INFO, "nsuid")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubeinghostinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_INGHOST_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_INGHOST_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_INGHOST_INFO+` (
				inghostid     serial NOT NULL PRIMARY KEY,
				clusterid     	integer NOT NULL,
				inguid     		text NOT NULL,
				backendtype   text NOT NULL,
				backendname   text NULL,
				hostname   		text NOT NULL DEFAULT '*',
				pathtype   		text NULL,
				path   				text NULL,
				serviceport   integer NULL,
				rscapigroup   text NULL,
				rsckind   		text NULL,
				enabled    		integer NOT NULL DEFAULT 0,
				createdtime   bigint NOT NULL,
				updatetime    bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		//Create View
		stmt_str := `create or replace view ` + TB_KUBE_ING_INFO_V + ` as
		select *,
		coalesce((select sum(svcv.pods) 
		   from ` + TB_KUBE_INGHOST_INFO + ` ih,
		   ` + TB_KUBE_SVC_INFO_V + ` svcv
		  where ih.inguid=ing.uid 
			and ih.backendtype='service'
			and svcv.svcname = ih.backendname),0) pods,
		coalesce((select sum(svcv.available_pods) 
		   from ` + TB_KUBE_INGHOST_INFO + ` ih,
		   ` + TB_KUBE_SVC_INFO_V + ` svcv
		  where ih.inguid=ing.uid 
			and ih.backendtype='service'
			and svcv.svcname = ih.backendname),0) available_pods
		from ` + TB_KUBE_ING_INFO + ` ing`
		_, err = tx.Exec(context.Background(), stmt_str)
		errorCheck(err)

		// Create View
		stmt_str = `create or replace view ` + TB_KUBE_INGPOD_INFO_V + ` as
		SELECT ing.ingid,ing.clusterid,ing.nsuid,ing.ingname,ing.uid inguid,ing.enabled,
				ih.backendtype, ih.backendname, ih.hostname, ih.pathtype, ih.path, ih.serviceport,
				svcpod.svcid, svcpod.svcname, svcpod.svcuid, svcpod.selector, svcpod.podid, svcpod.poduid, svcpod.podname
		FROM ` + TB_KUBE_ING_INFO + ` ing, ` + TB_KUBE_INGHOST_INFO + ` ih, ` + TB_KUBE_SVCPOD_INFO_V + ` svcpod
		where ing.enabled =1
			and svcpod.enabled=1
			and ing.nsuid =svcpod.nsuid
			and ing.clusterid =svcpod.clusterid
			and ih.inguid = ing.uid 
			and ih.backendtype ='service'
			and svcpod.svcname = ih.backendname`
		_, err = tx.Exec(context.Background(), stmt_str)
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_INGHOST_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubesvcinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_SVC_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_SVC_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_SVC_INFO+` (
				svcid     		serial NOT NULL PRIMARY KEY,
				clusterid     	integer NOT NULL,
				nsuid     		text NOT NULL,
				svcname   		text NOT NULL,
				uid   				text NOT NULL,
				starttime   	bigint NOT NULL,
				labels			text NULL,
				selector		text NULL,
				servicetype   text NULL,
				clusterip   	text NULL,
				ports   			text NULL,
				enabled    		integer NOT NULL DEFAULT 0,
				createdtime   bigint NOT NULL,
				updatetime    bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		//Create View
		stmt_str := `create or replace view ` + TB_KUBE_SVC_INFO_V + ` as
		select *,
		coalesce((select count(*) from (select nsuid, string_to_array(labels,',') as label_arr from ` + TB_KUBE_POD_INFO + ` where enabled=1) pod
			where pod.nsuid=svc.nsuid
			  and pod.label_arr@>string_to_array(svc.selector,',')
			  and array_length(string_to_array(svc.selector,','),1)>0),0) pods,
		coalesce((select count(*) from (select nsuid, string_to_array(labels,',') as label_arr from ` + TB_KUBE_POD_INFO + ` where enabled=1 and status='Running') pod
			where pod.nsuid=svc.nsuid
			  and pod.label_arr@>string_to_array(svc.selector,',')
			  and array_length(string_to_array(svc.selector,','),1)>0),0) available_pods
		from ` + TB_KUBE_SVC_INFO + ` svc
		where enabled=1`
		_, err = tx.Exec(context.Background(), stmt_str)
		errorCheck(err)

		//Create View
		stmt_str = `create or replace view ` + TB_KUBE_SVCPOD_INFO_V + ` as
		SELECT svc.svcid,svc.clusterid,svc.nsuid,
				svc.svcname,svc.uid svcuid, svc.selector,svc.enabled, pod.podid, pod.uid poduid, pod.podname, pod.labels 
			FROM ` + TB_KUBE_SVC_INFO + ` svc, ` + TB_KUBE_POD_INFO + ` pod
			WHERE svc.enabled = 1
				and pod.enabled = 1
				and svc.nsuid = pod.nsuid 
				and svc.clusterid =pod.clusterid 
				and string_to_array(pod.labels,','::text) @> string_to_array(svc.selector,','::text)  
				and array_length(string_to_array(svc.selector, ','::text), 1) > 0;`
		_, err = tx.Exec(context.Background(), stmt_str)
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_SVC_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_SVC_INFO, "uid, svcname, nsuid")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubedeployinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_DEPLOY_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_DEPLOY_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_DEPLOY_INFO+` (
				deployid     		serial NOT NULL PRIMARY KEY,
				clusterid     	integer NOT NULL,
				nsuid     			text NOT NULL,
				deployname   		text NOT NULL,
				uid   					text NOT NULL,
				starttime   		bigint NOT NULL,
				labels			text NULL,
				selector		text NULL,
				serviceaccount  text NULL,
				replicas   			bigint NULL DEFAULT 0,
				updatedrs   		bigint NULL DEFAULT 0,
				readyrs   			bigint NULL DEFAULT 0,
				availablers   	bigint NULL DEFAULT 0,
				observedgen   	bigint NULL DEFAULT 0,
				enabled    			integer NOT NULL DEFAULT 0,
				createdtime   	bigint NOT NULL,
				updatetime    	bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_DEPLOY_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_DEPLOY_INFO, "nsuid")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubestsinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_STS_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_STS_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_STS_INFO+` (
				stsid     			serial NOT NULL PRIMARY KEY,
				clusterid     	integer NOT NULL,
				nsuid     			text NOT NULL,
				stsname   			text NOT NULL,
				uid   					text NOT NULL,
				starttime   		bigint NOT NULL,
				labels			text NULL,
				selector		text NULL,
				serviceaccount  text NULL,
				replicas   			bigint NULL DEFAULT 0,
				updatedrs   		bigint NULL DEFAULT 0,
				readyrs   			bigint NULL DEFAULT 0,
				availablers   	bigint NULL DEFAULT 0,
				enabled    			integer NOT NULL DEFAULT 0,
				createdtime   	bigint NOT NULL,
				updatetime    	bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_STS_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_STS_INFO, "nsuid")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubedsinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_DS_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_DS_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_DS_INFO+` (
				dsid     				serial NOT NULL PRIMARY KEY,
				clusterid     	integer NOT NULL,
				nsuid     			text NOT NULL,
				dsname   				text NOT NULL,
				uid   					text NOT NULL,
				starttime   		bigint NOT NULL,
				labels			text NULL,
				selector		text NULL,
				serviceaccount  text NULL,
				current   			bigint NULL DEFAULT 0,
				desired   			bigint NULL DEFAULT 0,
				ready   				bigint NULL DEFAULT 0,
				updated   			bigint NULL DEFAULT 0,
				available   		bigint NULL DEFAULT 0,
				enabled    			integer NOT NULL DEFAULT 0,
				createdtime   	bigint NOT NULL,
				updatetime    	bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_DS_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_DS_INFO, "nsuid")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubersinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_RS_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_RS_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_RS_INFO+` (
				rsid     				serial NOT NULL PRIMARY KEY,
				clusterid     	integer NOT NULL,
				nsuid     			text NOT NULL,
				rsname   				text NOT NULL,
				uid   					text NOT NULL,
				starttime   		bigint NOT NULL,
				labels			text NULL,
				selector		text NULL,
				replicas  			bigint NULL DEFAULT 0,
				fullylabeledrs  bigint NULL DEFAULT 0,
				readyrs   			bigint NULL DEFAULT 0,
				availablers   	bigint NULL DEFAULT 0,
				observedgen   	bigint NULL DEFAULT 0,
				refkind   			text NULL,
				refuid   				text NULL,
				enabled    			integer NOT NULL DEFAULT 0,
				createdtime   	bigint NOT NULL,
				updatetime    	bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_RS_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_RS_INFO, "nsuid, refkind, refuid")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubepvcinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_PVC_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_PVC_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_PVC_INFO+` (
				pvcid     		serial NOT NULL PRIMARY KEY,
				clusterid     	integer NOT NULL,
				nsuid     		text NOT NULL,
				pvcname   		text NOT NULL,
				uid   		text NOT NULL,
				starttime   	bigint NOT NULL,
				labels			text NULL,
				selector		text NULL,
				accessmodes  	text NULL,
				reqstorage  	bigint NOT NULL default 0,
				status   		text NULL,
				scuid   			text NULL,
				enabled    		integer NOT NULL DEFAULT 0,
				createdtime   	bigint NOT NULL,
				updatetime    	bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_PVC_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_PVC_INFO, "nsuid, scuid")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubepvinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_PV_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_PV_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_PV_INFO+` (
				pvid     			serial NOT NULL PRIMARY KEY,
				clusterid     	integer NOT NULL,
				pvname   			text NOT NULL,
				pvuid   			text NOT NULL,
				pvcuid   			text NOT NULL,
				starttime   		bigint NOT NULL,
				labels			text NULL,
				accessmodes   		text NULL,
				capacity   			bigint NOT NULL default 0,
				reclaimpolicy   	text NULL,
				status   			text NULL,
				enabled    			integer NOT NULL DEFAULT 0,
				createdtime   		bigint NOT NULL,
				updatetime    		bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_PV_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_PV_INFO, "pvcuid")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubeeventinfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_EVENT_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_EVENT_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_EVENT_INFO+` (
				eventid     			serial NOT NULL PRIMARY KEY,
				clusterid     	integer NOT NULL,
				nsuid     				text NOT NULL,
				eventname   			text NOT NULL,
				uid   						text NOT NULL,
				firsttime   			bigint NOT NULL,
				lasttime   				bigint NOT NULL DEFAULT 0,
				labels					text NULL,
				eventtype  				text NOT NULL,
				eventcount   			bigint NULL DEFAULT 0,
				objkind  					text NOT NULL,
				objuid   			text NOT NULL,
				srccomponent   		text NOT NULL,
				srchost   				text NOT NULL,
				reason   					text NOT NULL,
				message   				text NULL,
				enabled    				integer NOT NULL DEFAULT 0,
				createdtime   		bigint NOT NULL,
				updatetime    		bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_EVENT_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_EVENT_INFO, "nsuid, objkind, objuid")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubeloginfo(conn *pgxpool.Conn, ontunetime int64) {
	if existKubeTableinfo(conn, TB_KUBE_LOG_INFO) {
		tablever := getKubeTableinfoVer(conn, TB_KUBE_LOG_INFO)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(),
			`CREATE TABLE IF NOT EXISTS `+TB_KUBE_LOG_INFO+` (
				logid   		serial NOT NULL PRIMARY KEY,
				logtype     	text NOT NULL DEFAULT 'pod',
				nsuid     		text NULL,
				poduid			text NULL,
				starttime   	bigint NOT NULL,
				message   		text NULL,
				createdtime   	bigint NOT NULL,
				updatetime    	bigint NOT NULL
			) `+getTableSpaceStmt(RESOURCE_TABLE, 0))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", TB_KUBE_LOG_INFO, 0, ontunetime, ontunetime, 0)
		errorCheck(err)

		createIndex(tx, TB_KUBE_LOG_INFO, "starttime, poduid, nsuid, logtype")

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

// realtime perf
func getTableName(tablename string) string {
	now := time.Now()
	createDay := now.Format(DATE_FORMAT)

	return tablename + "_" + createDay
}

func getMetricColumns(col_type string) string {
	switch col_type {
	case "cpu_raw":
		return `
		cpuusagesecondstotal	double precision NOT NULL default 0,
		cpusystemsecondstotal 	double precision NOT NULL default 0,
		cpuusersecondstotal 	double precision NOT NULL default 0,	
		`
	case "cpu":
		return `
		cpuusage 				double precision NOT NULL default 0,
		cpusystem 				double precision NOT NULL default 0,
		cpuuser 				double precision NOT NULL default 0,
		cpuusagecores 			double precision NOT NULL default 0,
		`
	case "cputotal":
		return `
		cputotalcores 			double precision NOT NULL default 0,
		`
	case "cpureqlimit":
		return `
		cpurequestcores			double precision NOT NULL default 0,
		cpulimitcores			double precision NOT NULL default 0,
		`
	case "memory_raw":
		return `
		memoryusagebytes 		double precision NOT NULL default 0,
		memoryworkingsetbytes 	double precision NOT NULL default 0,
		memorycache 			double precision NOT NULL default 0,
		memoryswap 				double precision NOT NULL default 0,
		memoryrss 				double precision NOT NULL default 0,
		`
	case "memory":
		return `
		memoryusage 			double precision NOT NULL default 0,
		memoryusagebytes		double precision NOT NULL default 0,			
		`
	case "memorysize":
		return `
		memorysizebytes 		double precision NOT NULL default 0,
		memoryswap 				double precision NOT NULL default 0,
		`
	case "memoryreqlimit":
		return `
		memoryrequestbytes 		double precision NOT NULL default 0,
		memorylimitbytes		double precision NOT NULL default 0,
		memoryswap 				double precision NOT NULL default 0,
		`
	case "net_raw":
		return `
		networkreceivebytestotal 	double precision NOT NULL default 0,
		networkreceiveerrorstotal 	double precision NOT NULL default 0,
		networktransmitbytestotal 	double precision NOT NULL default 0,
		networktransmiterrorstotal 	double precision NOT NULL default 0,
		`
	case "net":
		return `
		netiorate 				double precision NOT NULL default 0,
		netioerrors 			double precision NOT NULL default 0,
		netreceiverate 			double precision NOT NULL default 0,
		netreceiveerrors 		double precision NOT NULL default 0,
		nettransmitrate 		double precision NOT NULL default 0,
		nettransmiterrors 		double precision NOT NULL default 0,
		`
	case "net_noprefix":
		return `
		iorate 					double precision NOT NULL default 0,
		ioerrors 				double precision NOT NULL default 0,
		receiverate 			double precision NOT NULL default 0,
		receiveerrors		 	double precision NOT NULL default 0,
		transmitrate 			double precision NOT NULL default 0,
		transmiterrors		 	double precision NOT NULL default 0,
		`
	case "fs_pre_raw":
		return `
		fsinodesfree 				double precision NULL default 0,
		fsinodestotal 				double precision NULL default 0,
		fslimitbytes 				double precision NULL default 0,
		`
	case "fsrw_raw":
		return `
		fsreadsbytestotal 		double precision NOT NULL default 0,
		fswritesbytestotal 		double precision NOT NULL default 0,
		`
	case "fs_post_raw":
		return `
		fsusagebytes 				double precision NULL default 0,
		`
	case "fs":
		return `
		fsiorate 				double precision NOT NULL default 0,
		fsreadrate 				double precision NOT NULL default 0,
		fswriterate 			double precision NOT NULL default 0,
		`
	case "fs_noprefix":
		return `
		iorate 					double precision NOT NULL default 0,
		readrate 				double precision NOT NULL default 0,
		writerate	 			double precision NOT NULL default 0,
		`
	case "process":
		return `
		processes 				double precision NOT NULL default 0,
		`
	}

	return ""
}

func createKubenodePerfRaw(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			nodeuid					text NOT NULL default '',
			ontunetime   			bigint NOT NULL,
			metricid				int NOT NULL default 0,
			%s
			%s
			%s
			%s
			timestampms   			bigint NOT NULL,
			CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, nodeuid, metricid)
		) %s`, tablename, getMetricColumns("cpu_raw"), getMetricColumns("memory_raw"), getMetricColumns("fsrw_raw"), getMetricColumns("process"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func createKubenodePerf(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			nodeuid					text NOT NULL default '',
			ontunetime   			bigint NOT NULL,
			metricid				int NOT NULL default 0,
			%s
			%s
			%s
			%s
			%s
			%s
			%s
			CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, nodeuid, metricid)
		) %s`, tablename, getMetricColumns("cpu"), getMetricColumns("cputotal"), getMetricColumns("memory"), getMetricColumns("memorysize"), getMetricColumns("net"), getMetricColumns("fs"), getMetricColumns("process"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		createIndex(tx, tablename, "ontunetime, nodeuid")
		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubenodePerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubenodePerfRaw(conn, TB_KUBE_LAST_NODE_PERF_RAW, ontunetime)
	createKubenodePerfRaw(conn, getTableName(TB_KUBE_NODE_PERF_RAW), ontunetime)
	createKubenodePerf(conn, TB_KUBE_LAST_NODE_PERF, ontunetime)
	createKubenodePerf(conn, getTableName(TB_KUBE_NODE_PERF), ontunetime)
}

func createKubepodPerfRaw(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			poduid     				text NOT NULL default '',
			ontunetime   			bigint NOT NULL,
			nsuid					text NOT NULL default '',
			nodeuid					text NOT NULL default '',
			metricid				int NOT NULL default 0,
			%s
			%s
			%s
			%s
			timestampms   			bigint NOT NULL,
			CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, poduid, nsuid, nodeuid, metricid)
		) %s`, tablename, getMetricColumns("cpu_raw"), getMetricColumns("memory_raw"), getMetricColumns("fsrw_raw"), getMetricColumns("process"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func createKubepodPerf(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			poduid     				text NOT NULL default '',
			ontunetime   			bigint NOT NULL,
			nsuid					text NOT NULL default '',
			nodeuid					text NOT NULL default '',
			metricid				int NOT NULL default 0,
			%s
			%s
			%s
			%s
			%s
			%s
			%s
			CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, poduid, nsuid, nodeuid, metricid)
		) %s`, tablename, getMetricColumns("cpu"), getMetricColumns("cpureqlimit"), getMetricColumns("memory"), getMetricColumns("memoryreqlimit"), getMetricColumns("net"), getMetricColumns("fs"), getMetricColumns("process"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		createIndex(tx, tablename, "ontunetime, poduid, nsuid")

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func createKubeClusterPerf(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			clusterid     			int NOT NULL,
			ontunetime   			bigint NOT NULL,
			podcount			    int NOT NULL,
			%s
			%s
			%s
			%s
			%s
			%s
			%s
			CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, clusterid)
		) %s`, tablename, getMetricColumns("cpu"), getMetricColumns("cputotal"), getMetricColumns("memory"), getMetricColumns("memorysize"), getMetricColumns("net"), getMetricColumns("fs"), getMetricColumns("process"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func createKubeNsPerf(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			nsuid					text NOT NULL default '',
			ontunetime   			bigint NOT NULL,
			podcount			    int NOT NULL,
			%s
			%s
			%s
			%s
			%s
			%s
			%s
			CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, nsuid)
		) %s`, tablename, getMetricColumns("cpu"), getMetricColumns("cputotal"), getMetricColumns("memory"), getMetricColumns("memorysize"), getMetricColumns("net"), getMetricColumns("fs"), getMetricColumns("process"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func createKubeWorkloadPerf(conn *pgxpool.Conn, tablename string, uidname string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			%s 						text NOT NULL default '',
			ontunetime   			bigint NOT NULL,
			podcount			    int NOT NULL,
			%s
			%s
			%s
			%s
			%s
			%s
			%s
			CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, %s)
		) %s`, tablename, uidname, getMetricColumns("cpu"), getMetricColumns("cputotal"), getMetricColumns("memory"), getMetricColumns("memorysize"), getMetricColumns("net"), getMetricColumns("fs"), getMetricColumns("process"), tablename, uidname, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func createKubeSvcPerf(conn *pgxpool.Conn, tablename string, ontunetime int64) {
}

func createKubeIngPerf(conn *pgxpool.Conn, tablename string, ontunetime int64) {
}

func checkKubepodPerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubepodPerfRaw(conn, TB_KUBE_LAST_POD_PERF_RAW, ontunetime)
	createKubepodPerfRaw(conn, getTableName(TB_KUBE_POD_PERF_RAW), ontunetime)
	createKubepodPerf(conn, TB_KUBE_LAST_POD_PERF, ontunetime)
	createKubepodPerf(conn, getTableName(TB_KUBE_POD_PERF), ontunetime)

	createKubeClusterPerf(conn, getTableName(TB_KUBE_CLUSTER_PERF), ontunetime)
	createKubeNsPerf(conn, getTableName(TB_KUBE_NAMESPACE_PERF), ontunetime)
	createKubeWorkloadPerf(conn, getTableName(TB_KUBE_REPLICASET_PERF), "rsuid", ontunetime)
	createKubeWorkloadPerf(conn, getTableName(TB_KUBE_DAEMONSET_PERF), "dsuid", ontunetime)
	createKubeWorkloadPerf(conn, getTableName(TB_KUBE_STATEFULSET_PERF), "stsuid", ontunetime)
	createKubeWorkloadPerf(conn, getTableName(TB_KUBE_DEPLOYMENT_PERF), "deployuid", ontunetime)
	createKubeSvcPerf(conn, getTableName(TB_KUBE_SERVICE_PERF), ontunetime)
	createKubeIngPerf(conn, getTableName(TB_KUBE_INGRESS_PERF), ontunetime)

	createKubeClusterPerf(conn, TB_KUBE_LAST_CLUSTER_PERF, ontunetime)
	createKubeNsPerf(conn, TB_KUBE_LAST_NAMESPACE_PERF, ontunetime)
	createKubeWorkloadPerf(conn, TB_KUBE_LAST_REPLICASET_PERF, "rsuid", ontunetime)
	createKubeWorkloadPerf(conn, TB_KUBE_LAST_DAEMONSET_PERF, "dsuid", ontunetime)
	createKubeWorkloadPerf(conn, TB_KUBE_LAST_STATEFULSET_PERF, "stsuid", ontunetime)
	createKubeWorkloadPerf(conn, TB_KUBE_LAST_DEPLOYMENT_PERF, "deployuid", ontunetime)
	createKubeSvcPerf(conn, TB_KUBE_LAST_SERVICE_PERF, ontunetime)
	createKubeIngPerf(conn, TB_KUBE_LAST_INGRESS_PERF, ontunetime)
}

func createKubecontainerPerfRaw(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				containername     		text NOT NULL default '',
				ontunetime   			bigint NOT NULL,
				poduid					text NOT NULL default '',
				nsuid					text NOT NULL default '',
				nodeuid					text NOT NULL default '',
				metricid				int NOT NULL default 0,
				%s
				%s
				%s
				timestampms   			bigint NOT NULL,
				CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, poduid, nsuid, nodeuid, containername, metricid)
			) %s`, tablename, getMetricColumns("cpu_raw"), getMetricColumns("memory_raw"), getMetricColumns("process"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func createKubecontainerPerf(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				containername     		text NOT NULL default '',
				ontunetime   			bigint NOT NULL,
				poduid					text NOT NULL default '',
				nsuid					text NOT NULL default '',
				nodeuid					text NOT NULL default '',
				metricid				int NOT NULL default 0,
				%s
				%s
				%s
				%s
				%s
				%s
				CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, poduid, nsuid, nodeuid, containername, metricid)
			) %s`, tablename, getMetricColumns("cpu"), getMetricColumns("cpureqlimit"), getMetricColumns("memory"), getMetricColumns("memoryreqlimit"), getMetricColumns("fs"), getMetricColumns("process"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		createIndex(tx, tablename, "ontunetime, poduid, containername")

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubecontainerPerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubecontainerPerfRaw(conn, TB_KUBE_LAST_CONTAINER_PERF_RAW, ontunetime)
	createKubecontainerPerfRaw(conn, getTableName(TB_KUBE_CONTAINER_PERF_RAW), ontunetime)
	createKubecontainerPerf(conn, TB_KUBE_LAST_CONTAINER_PERF, ontunetime)
	createKubecontainerPerf(conn, getTableName(TB_KUBE_CONTAINER_PERF), ontunetime)
}

func createKubeNodeNetPerfRaw(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				nodeuid						text NOT NULL default '',
				ontunetime		   			bigint NOT NULL,
				metricid					int NOT NULL default 0,
				interfaceid               	int NOT NULL default 0,
				%s
				timestampms		   			bigint NOT NULL,
				CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, nodeuid, metricid, interfaceid)
			) %s`, tablename, getMetricColumns("net_raw"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func createKubeNodeNetPerf(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			nodeuid					text NOT NULL default '',
			ontunetime   			bigint NOT NULL,
			metricid				int NOT NULL default 0,
			interfaceid             int NOT NULL default 0,
			%s
		CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, nodeuid, metricid, interfaceid)
		) %s`, tablename, getMetricColumns("net_noprefix"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		createIndex(tx, tablename, "ontunetime, nodeuid, interfaceid")

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubeNodeNetPerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubeNodeNetPerfRaw(conn, getTableName(TB_KUBE_NODE_NET_PERF_RAW), ontunetime)
	createKubeNodeNetPerf(conn, TB_KUBE_LAST_NODE_NET_PERF, ontunetime)
	createKubeNodeNetPerf(conn, getTableName(TB_KUBE_NODE_NET_PERF), ontunetime)
}

func createKubePodNetPerfRaw(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				poduid     					text NOT NULL default '',
				ontunetime		   			bigint NOT NULL,
				nsuid						text NOT NULL default '',
				nodeuid						text NOT NULL default '',
				metricid					int NOT NULL default 0,
				interfaceid               	int NOT NULL default 0,
				%s
				timestampms		   			bigint NOT NULL,
				CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, poduid, nsuid, nodeuid, metricid, interfaceid)
			) %s`, tablename, getMetricColumns("net_raw"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func createKubePodNetPerf(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			poduid     				text NOT NULL default '',
			ontunetime   			bigint NOT NULL,
			nsuid					text NOT NULL default '',
			nodeuid					text NOT NULL default '',
			metricid				int NOT NULL default 0,
			interfaceid             int NOT NULL default 0,
			%s
			CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, poduid, nsuid, nodeuid, metricid, interfaceid)
		) %s`, tablename, getMetricColumns("net_noprefix"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		createIndex(tx, tablename, "ontunetime, poduid, nsuid, interfaceid")

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubePodNetPerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubePodNetPerfRaw(conn, getTableName(TB_KUBE_POD_NET_PERF_RAW), ontunetime)
	createKubePodNetPerf(conn, TB_KUBE_LAST_POD_NET_PERF, ontunetime)
	createKubePodNetPerf(conn, getTableName(TB_KUBE_POD_NET_PERF), ontunetime)
}

func createKubeNodeFsPerfRaw(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				nodeuid						text NOT NULL default '',
				ontunetime  		 		bigint NOT NULL,
				metricid					int NOT NULL default 0,				
				deviceid					int NOT NULL default 0,
				%s
				%s
				%s
				timestampms		   			bigint NOT NULL,
				CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, nodeuid, metricid, deviceid)
			) %s`, tablename, getMetricColumns("fs_pre_raw"), getMetricColumns("fsrw_raw"), getMetricColumns("fs_post_raw"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func createKubeNodeFsPerf(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			nodeuid					text NOT NULL default '',
			ontunetime   			bigint NOT NULL,
			metricid				int NOT NULL default 0,
			deviceid             	int NOT NULL default 0,
			%s
		CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, nodeuid, metricid, deviceid)
		) %s`, tablename, getMetricColumns("fs_noprefix"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		createIndex(tx, tablename, "ontunetime, nodeuid, deviceid")

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubeNodeFsPerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubeNodeFsPerfRaw(conn, getTableName(TB_KUBE_NODE_FS_PERF_RAW), ontunetime)
	createKubeNodeFsPerf(conn, TB_KUBE_LAST_NODE_FS_PERF, ontunetime)
	createKubeNodeFsPerf(conn, getTableName(TB_KUBE_NODE_FS_PERF), ontunetime)
}

func createKubePodFsPerfRaw(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				poduid						text NOT NULL default '',
				ontunetime  		 		bigint NOT NULL,
				nsuid						text NOT NULL default '',
				nodeuid						text NOT NULL default '',
				metricid					int NOT NULL default 0,				
				deviceid					int NOT NULL default 0,
				%s
				timestampms		   			bigint NOT NULL,
				CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, poduid, nsuid, nodeuid, metricid, deviceid)
			) %s`, tablename, getMetricColumns("fsrw_raw"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func createKubePodFsPerf(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			poduid     					text NOT NULL default '',
			ontunetime		   			bigint NOT NULL,
			nsuid						text NOT NULL default '',
			nodeuid						text NOT NULL default '',
			metricid					int NOT NULL default 0,
			deviceid             		int NOT NULL default 0,
			%s
		CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, poduid, nsuid, nodeuid, metricid, deviceid)
		) %s`, tablename, getMetricColumns("fs_noprefix"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		createIndex(tx, tablename, "ontunetime, poduid, nsuid, deviceid")

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubePodFsPerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubePodFsPerfRaw(conn, getTableName(TB_KUBE_POD_FS_PERF_RAW), ontunetime)
	createKubePodFsPerf(conn, TB_KUBE_LAST_POD_FS_PERF, ontunetime)
	createKubePodFsPerf(conn, getTableName(TB_KUBE_POD_FS_PERF), ontunetime)
}

func createKubeContainerFsPerfRaw(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				containername     			text NOT NULL default '',
				ontunetime  		 		bigint NOT NULL,
				poduid						text NOT NULL default '',
				nsuid						text NOT NULL default '',
				nodeuid						text NOT NULL default '',
				metricid					int NOT NULL default 0,				
				deviceid					int NOT NULL default 0,
				%s
				timestampms		   			bigint NOT NULL,
				CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, poduid, nsuid, nodeuid, containername, metricid, deviceid)
			) %s`, tablename, getMetricColumns("fsrw_raw"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		setAutoVacuum(tx, tablename)
	}
}

func createKubeContainerFsPerf(conn *pgxpool.Conn, tablename string, ontunetime int64) {
	if existKubeTableinfo(conn, tablename) {
		tablever := getKubeTableinfoVer(conn, tablename)
		ProcessTableVersion(tablever)
	} else {
		//Create Table
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				containername     		text NOT NULL default '',
				ontunetime   			bigint NOT NULL,
				poduid					text NOT NULL default '',
				nsuid					text NOT NULL default '',
				nodeuid					text NOT NULL default '',
				metricid				int NOT NULL default 0,
				deviceid             	int NOT NULL default 0,
				%s
				CONSTRAINT %s_pkey PRIMARY KEY (ontunetime, poduid, nsuid, nodeuid, containername, metricid, deviceid)
			) %s`, tablename, getMetricColumns("fs_noprefix"), tablename, getTableSpaceStmt(SHORTERM_TABLE, ontunetime)))
		errorCheck(err)

		_, err = tx.Exec(context.Background(), "INSERT INTO  "+TB_KUBE_TABLE_INFO+"  VALUES ($1, $2, $3, $4, $5)", tablename, 0, ontunetime, ontunetime, common.ShorttermDuration)
		errorCheck(err)

		createIndex(tx, tablename, "ontunetime, poduid, containername, deviceid")

		setAutoVacuum(tx, tablename)

		err = tx.Commit(context.Background())
		errorCheck(err)
	}
}

func checkKubeContainerFsPerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubeContainerFsPerfRaw(conn, getTableName(TB_KUBE_CONTAINER_FS_PERF_RAW), ontunetime)
	createKubeContainerFsPerf(conn, TB_KUBE_LAST_CONTAINER_FS_PERF, ontunetime)
	createKubeContainerFsPerf(conn, getTableName(TB_KUBE_CONTAINER_FS_PERF), ontunetime)
}

// Daily Table
func getDailyTableName(tablename string) string {
	now := time.Now()
	createDay := now.AddDate(0, 0, 1).Format(DATE_FORMAT)

	return tablename + "_" + createDay
}

func dailyKubenodePerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubenodePerfRaw(conn, getDailyTableName(TB_KUBE_NODE_PERF_RAW), ontunetime)
	createKubenodePerf(conn, getDailyTableName(TB_KUBE_NODE_PERF), ontunetime)
}

func dailyKubepodPerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubepodPerfRaw(conn, getDailyTableName(TB_KUBE_POD_PERF_RAW), ontunetime)
	createKubepodPerf(conn, getDailyTableName(TB_KUBE_POD_PERF), ontunetime)

	createKubeClusterPerf(conn, getDailyTableName(TB_KUBE_CLUSTER_PERF), ontunetime)
	createKubeNsPerf(conn, getDailyTableName(TB_KUBE_NAMESPACE_PERF), ontunetime)
	createKubeWorkloadPerf(conn, getDailyTableName(TB_KUBE_REPLICASET_PERF), "rsuid", ontunetime)
	createKubeWorkloadPerf(conn, getDailyTableName(TB_KUBE_DAEMONSET_PERF), "dsuid", ontunetime)
	createKubeWorkloadPerf(conn, getDailyTableName(TB_KUBE_STATEFULSET_PERF), "stsuid", ontunetime)
	createKubeWorkloadPerf(conn, getDailyTableName(TB_KUBE_DEPLOYMENT_PERF), "deployuid", ontunetime)
	createKubeSvcPerf(conn, getDailyTableName(TB_KUBE_SERVICE_PERF), ontunetime)
	createKubeIngPerf(conn, getDailyTableName(TB_KUBE_INGRESS_PERF), ontunetime)
}

func dailyKubecontainerPerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubecontainerPerfRaw(conn, getDailyTableName(TB_KUBE_CONTAINER_PERF_RAW), ontunetime)
	createKubecontainerPerf(conn, getDailyTableName(TB_KUBE_CONTAINER_PERF), ontunetime)
}

func dailyKubeNodeNetPerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubeNodeNetPerfRaw(conn, getDailyTableName(TB_KUBE_NODE_NET_PERF_RAW), ontunetime)
	createKubeNodeNetPerf(conn, getDailyTableName(TB_KUBE_NODE_NET_PERF), ontunetime)
}

func dailyKubePodNetPerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubePodNetPerfRaw(conn, getDailyTableName(TB_KUBE_POD_NET_PERF_RAW), ontunetime)
	createKubePodNetPerf(conn, getDailyTableName(TB_KUBE_POD_NET_PERF), ontunetime)
}

func dailyKubeNodeFsPerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubeNodeFsPerfRaw(conn, getDailyTableName(TB_KUBE_NODE_FS_PERF_RAW), ontunetime)
	createKubeNodeFsPerf(conn, getDailyTableName(TB_KUBE_NODE_FS_PERF), ontunetime)
}

func dailyKubePodFsPerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubePodFsPerfRaw(conn, getDailyTableName(TB_KUBE_POD_FS_PERF_RAW), ontunetime)
	createKubePodFsPerf(conn, getDailyTableName(TB_KUBE_POD_FS_PERF), ontunetime)
}

func dailyKubeContainerFsPerf(conn *pgxpool.Conn, ontunetime int64) {
	createKubeContainerFsPerfRaw(conn, getDailyTableName(TB_KUBE_CONTAINER_FS_PERF_RAW), ontunetime)
	createKubeContainerFsPerf(conn, getDailyTableName(TB_KUBE_CONTAINER_FS_PERF), ontunetime)
}

func createIndex(tx pgx.Tx, tablename string, column string) {
	SQL := "CREATE INDEX IF NOT EXISTS i" + tablename + " ON " + tablename + " using btree ( " + column + ")"
	_, err := tx.Exec(context.Background(), SQL)
	errorCheck(err)
}

func existKubeTableinfo(conn *pgxpool.Conn, tablename string) bool {
	var selName string
	SQL := "select tablename from " + TB_KUBE_TABLE_INFO + " where tablename = '" + tablename + "'"

	err := conn.QueryRow(context.Background(), SQL).Scan(&selName)
	errorCheck(err)
	if selName == "" {
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		_, err = tx.Exec(context.Background(), fmt.Sprintf("DROP TABLE IF EXISTS %s", tablename))
		errorCheck(err)

		err = tx.Commit(context.Background())
		errorCheck(err)

		return false
	} else {
		var selPgName string
		err := conn.QueryRow(context.Background(), "select tablename from pg_tables where tablename = $1", tablename).Scan(&selPgName)
		errorCheck(err)

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		errorCheck(err)

		if selPgName == "" {
			_, err = tx.Exec(context.Background(), fmt.Sprintf("DELETE FROM %s where tablename = '%s'", TB_KUBE_TABLE_INFO, tablename))
			errorCheck(err)

			err = tx.Commit(context.Background())
			errorCheck(err)

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
	errorCheck(err)

	return selVer
}
