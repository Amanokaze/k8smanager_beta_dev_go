package database

const INSERT_MANAGER_INFO_SQL = `
INSERT INTO kubemanagerinfo (
		managername,
		description,
		ip,
		createdtime,
		updatetime
	) VALUES (
		$1, $2, $3, $4, $5
	) RETURNING managerid
	`
const SELECT_MANAGER_INFO = `
	SELECT managerid from ` + TB_KUBE_MANAGER_INFO + ` WHERE managername = $1
	`
const INSERT_CLUSTER_INFO_SQL = `
	INSERT INTO kubeclusterinfo (
			managerid,
			clustername,
			context,
			ip,
			createdtime,
			updatetime
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING clusterid, enabled
		`

const SELECT_CLUSTER_INFO = `
	SELECT clusterid, enabled, status from ` + TB_KUBE_CLUSTER_INFO + ` WHERE ip = $1
`

const SELECT_UNUSED_CLUSTER_INFO = `
	SELECT clusterid from ` + TB_KUBE_CLUSTER_INFO + ` WHERE enabled = 1 and ip not in (%s)
`

const UPDATE_CLUSTER_RESET = `
	UPDATE ` + TB_KUBE_CLUSTER_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
`

const UPDATE_CLUSTER_ENABLED = `
	UPDATE ` + TB_KUBE_CLUSTER_INFO + ` SET enabled = 1, updatetime = $2 WHERE clusterid = $1
`

const UPDATE_CLUSTER_NODE_RESET = `
	UPDATE ` + TB_KUBE_NODE_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
`

const UPDATE_CLUSTER_NS_RESET = `
	UPDATE ` + TB_KUBE_NS_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
`

const UPDATE_CLUSTER_POD_RESET = `
	UPDATE ` + TB_KUBE_POD_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
`

const UPDATE_CLUSTER_CONTAINER_RESET = `
	UPDATE ` + TB_KUBE_CONTAINER_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
`

const UPDATE_CLUSTER_SVC_RESET = `
	UPDATE ` + TB_KUBE_SVC_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
`

const UPDATE_CLUSTER_PV_RESET = `
	UPDATE ` + TB_KUBE_PV_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
`

const UPDATE_CLUSTER_PVC_RESET = `
	UPDATE ` + TB_KUBE_PVC_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
`

const UPDATE_CLUSTER_SC_RESET = `
	UPDATE ` + TB_KUBE_SC_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
`

const UPDATE_CLUSTER_EVENT_RESET = `
	UPDATE ` + TB_KUBE_EVENT_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
	`

const UPDATE_CLUSTER_DEPLOY_RESET = `
	UPDATE ` + TB_KUBE_DEPLOY_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
	`

const UPDATE_CLUSTER_RS_RESET = `
	UPDATE ` + TB_KUBE_RS_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
	`

const UPDATE_CLUSTER_DS_RESET = `
	UPDATE ` + TB_KUBE_DS_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
	`

const UPDATE_CLUSTER_STS_RESET = `
	UPDATE ` + TB_KUBE_STS_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
	`

const UPDATE_CLUSTER_ING_RESET = `
	UPDATE ` + TB_KUBE_ING_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
	`

const UPDATE_CLUSTER_INGHOST_RESET = `
	UPDATE ` + TB_KUBE_INGHOST_INFO + ` SET enabled = 0, updatetime = $2 WHERE clusterid = $1
	`

const SELECT_ROW = `
	SELECT * AS cnt from
`

const UPDATE_TABLEINFO = `
	UPDATE ` + TB_KUBE_TABLE_INFO + ` set updatetime = $1 where tablename = $2
`

const SELECT_ONTUNE_TIME = `
	SELECT _time, _bias from ontuneinfo
`

const SELECT_CURRENT_TIME = `
	SELECT floor(extract(epoch from now()))::bigint currenttime, _bias from ontuneinfo
`

const INSERT_RESOURCE_INFO = `
	INSERT INTO ` + TB_KUBE_RESOURCE_INFO + ` (clusterid, resourcename, apiclass, version, endpoint, enabled, createdtime, updatetime)
	(select * from unnest($1::int[], $2::text[], $3::text[], $4::text[], $5::text[], $6::int[], $7::bigint[], $8::bigint[]))
`
const INSERT_UNNEST_NAMESPACE_INFO = `
	INSERT INTO ` + TB_KUBE_NS_INFO + ` (nsuid, clusterid, nsname, starttime, labels, status, enabled, createdtime, updatetime)
	(select * from unnest($1::text[], $2::int[], $3::text[], $4::bigint[], $5::text[], $6::text[], $7::int[], $8::bigint[], $9::bigint[]))
`

const UPDATE_NAMESPACE_INFO = `
	UPDATE ` + TB_KUBE_NS_INFO + ` SET clusterid=$1, nsname=$2, starttime=$3, labels=$4, status=$5, enabled=1, updatetime=$6 WHERE nsuid = $7
`

const INSERT_NAMESPACE_INFO = `
	INSERT INTO ` + TB_KUBE_NS_INFO + ` (nsuid, clusterid, nsname, starttime, labels, status, enabled, createdtime, updatetime) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`
const INSERT_UNNEST_NODE_INFO = `
	INSERT INTO ` + TB_KUBE_NODE_INFO + ` (managerid, clusterid, nodeuid, nodename, nodenameext, nodetype
	, enabled, starttime, labels, kernelversion, osimage, osname, containerruntimever, kubeletver, kubeproxyver, cpuarch, cpucount, ephemeralstorage, memorysize, pods, ip, status, createdtime, updatetime)
	(select * from unnest($1::int[], $2::int[], $3::text[], $4::text[], $5::text[], $6::text[], $7::int[], $8::bigint[], $9::text[], $10::text[], $11::text[], $12::text[], 
		$13::text[], $14::text[], $15::text[], $16::text[], $17::int[], $18::bigint[], $19::bigint[], $20::bigint[], $21::text[], $22::int[], $23::bigint[], $24::bigint[]))
`

const UPDATE_NODE_INFO = `
	UPDATE ` + TB_KUBE_NODE_INFO + ` SET managerid = $1, clusterid = $2, enabled = 1, nodename = $3, nodetype = $4, starttime = $5, labels = $6
	, kernelversion = $7, osimage = $8, osname = $9, containerruntimever = $10, kubeletver = $11, kubeproxyver = $12, cpuarch = $13, cpucount = $14, ephemeralstorage = $15, memorysize = $16
	, pods = $17, ip = $18, status=$19, updatetime = $20 WHERE nodeuid = $21 and enabled=1
`

const INSERT_NODE_INFO = `
	INSERT INTO ` + TB_KUBE_NODE_INFO + ` (managerid, clusterid, nodeuid, nodename, nodenameext, nodetype, enabled, starttime, labels, kernelversion, osimage, osname, 
		containerruntimever, kubeletver, kubeproxyver, cpuarch, cpucount, ephemeralstorage, memorysize, pods, ip, status, createdtime, updatetime) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
`

const UPDATE_POD_INFO = `
	UPDATE ` + TB_KUBE_POD_INFO + ` SET clusterid = $22, nodeuid=$1, nsuid=$2, annotationuid=$3, podname=$4, starttime=$5, labels=$6, selector=$7, restartpolicy=$8, serviceaccount=$9, status=$10, 
	hostip=$11, podip=$12, restartcount=$13, restarttime=$14, podcondition=$15, staticpod=$16, refkind=$17, refuid=$18, pvcuid=$19, enabled=1, updatetime=$20 WHERE uid = $21 and enabled=1
`

const INSERT_POD_INFO = `
	INSERT INTO ` + TB_KUBE_POD_INFO + ` (uid, nodeuid, nsuid, annotationuid, podname, starttime, labels, selector, restartpolicy, serviceaccount, status, hostip, 
		podip, restartcount, restarttime, podcondition, staticpod, refkind, refuid, pvcuid, enabled, createdtime, updatetime, clusterid)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
`

const INSERT_UNNEST_POD_INFO = `
	INSERT INTO ` + TB_KUBE_POD_INFO + ` (clusterid, uid, nodeuid, nsuid, annotationuid, podname, starttime, labels, selector, restartpolicy, serviceaccount, status, hostip, 
		podip, restartcount, restarttime, podcondition, staticpod, refkind, refuid, pvcuid, enabled, createdtime, updatetime)
	(select * from  unnest($1::int[], $2::text[], $3::text[], $4::text[], $5::text[], $6::text[], $7::bigint[], $8::text[], $9::text[], $10::text[], $11::text[], $12::text[], $13::text[], 
		$14::text[], $15::bigint[], $16::bigint[], $17::text[], $18::text[], $19::text[], $20::text[], $21::text[], $22::int[], $23::bigint[], $24::bigint[]))
`

const UPDATE_CONTAINER_INFO = `
	UPDATE ` + TB_KUBE_CONTAINER_INFO + ` SET clusterid = $17, image=$2, ports=$3, env=$4, limitcpu=$5, limitmemory=$6, limitstorage=$7, limitephemeral=$8, 
		reqcpu=$9, reqmemory=$10, reqstorage=$11, reqephemeral=$12, volumemounts=$13, state=$14, enabled=1, updatetime=$15 WHERE containername = $16 and poduid=$1 and enabled=1
`

const INSERT_CONTAINER_INFO = `
	INSERT INTO ` + TB_KUBE_CONTAINER_INFO + ` (poduid, containername, image, ports, env, limitcpu, limitmemory, limitstorage, limitephemeral, reqcpu, reqmemory, reqstorage, 
		reqephemeral, volumemounts, state, enabled, createdtime, updatetime, clusterid) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
`

const INSERT_UNNEST_CONTAINER_INFO = `
	INSERT INTO ` + TB_KUBE_CONTAINER_INFO + ` (clusterid, poduid, containername, image, ports, env, limitcpu, limitmemory, limitstorage, limitephemeral, reqcpu, reqmemory, reqstorage, 
		reqephemeral, volumemounts, state, enabled, createdtime, updatetime)
	(select * from  unnest($1::int[], $2::text[], $3::text[], $4::text[], $5::text[], $6::text[], $7::bigint[], $8::bigint[], $9::bigint[], $10::bigint[], $11::bigint[], $12::bigint[], 
		$13::bigint[], $14::bigint[], $15::text[], $16::text[], $17::int[], $18::bigint[], $19::bigint[]))
`

const UPDATE_SVC_INFO = `
	UPDATE ` + TB_KUBE_SVC_INFO + ` SET clusterid=$11, nsuid=$1, svcname=$2, starttime=$3, labels=$4, selector=$5, servicetype=$6, clusterip=$7, ports=$8, enabled=1, updatetime = $9 WHERE uid = $10 and enabled=1
`

const INSERT_SVC_INFO = `
	INSERT INTO ` + TB_KUBE_SVC_INFO + ` (nsuid, svcname, uid, starttime, labels, selector, servicetype, clusterip, ports, enabled, createdtime, updatetime, clusterid)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
`

const INSERT_UNNEST_SVC_INFO = `
	INSERT INTO ` + TB_KUBE_SVC_INFO + ` (clusterid, nsuid, svcname, uid, starttime, labels, selector, servicetype, clusterip, ports, enabled, createdtime, updatetime)
	(select * from  unnest($1::int[], $2::text[], $3::text[], $4::text[], $5::bigint[], $6::text[], $7::text[], $8::text[], $9::text[], $10::text[], $11::int[], $12::bigint[], $13::bigint[]))
`

const UPDATE_DEPLOY_INFO = `
	UPDATE ` + TB_KUBE_DEPLOY_INFO + ` SET clusterid = $14, nsuid=$1, deployname=$2, starttime=$3, labels=$4, selector=$5, serviceaccount=$6, replicas=$7, updatedrs=$8, readyrs=$9, availablers=$10, observedgen=$11, enabled=1, updatetime = $12 WHERE uid = $13 and enabled=1
`

const INSERT_DEPLOY_INFO = `
INSERT INTO ` + TB_KUBE_DEPLOY_INFO + ` (nsuid, deployname, uid, starttime, labels, selector, serviceaccount, replicas, updatedrs, readyrs, availablers, observedgen, enabled, createdtime, updatetime, clusterid)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
`

const INSERT_UNNEST_DEPLOY_INFO = `
	INSERT INTO ` + TB_KUBE_DEPLOY_INFO + ` (clusterid, nsuid, deployname, uid, starttime, labels, selector, serviceaccount, replicas, updatedrs, readyrs, availablers, observedgen, enabled, createdtime, updatetime)
	(select * from  unnest($1::int[], $2::text[], $3::text[], $4::text[], $5::bigint[], $6::text[], $7::text[], $8::text[], $9::bigint[], $10::bigint[], $11::bigint[], $12::bigint[], $13::bigint[], $14::int[], $15::bigint[], $16::bigint[]))
`

const UPDATE_STATEFUL_INFO = `
	UPDATE ` + TB_KUBE_STS_INFO + ` SET clusterid = $13, nsuid=$1, stsname=$2, starttime=$3, labels=$4, selector=$5, serviceaccount=$6, replicas=$7, updatedrs=$8, readyrs=$9, availablers=$10, enabled=1, updatetime = $11 WHERE uid = $12 and enabled=1
`

const INSERT_STATEFUL_INFO = `
INSERT INTO ` + TB_KUBE_STS_INFO + ` (nsuid, stsname, uid, starttime, labels, selector, serviceaccount, replicas, updatedrs, readyrs, availablers, enabled, createdtime, updatetime, clusterid)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
`
const INSERT_UNNEST_STATEFUL_INFO = `
	INSERT INTO ` + TB_KUBE_STS_INFO + ` (clusterid, nsuid, stsname, uid, starttime, labels, selector, serviceaccount, replicas, updatedrs, readyrs, availablers, enabled, createdtime, updatetime)
	(select * from  unnest($1::int[], $2::text[], $3::text[], $4::text[], $5::bigint[], $6::text[], $7::text[], $8::text[], $9::bigint[], $10::bigint[], $11::bigint[], $12::bigint[], $13::int[], $14::bigint[], $15::bigint[]))
`

const UPDATE_DAEMONSET_INFO = `
	UPDATE ` + TB_KUBE_DS_INFO + ` SET clusterid = $14, nsuid=$1, dsname=$2, starttime=$3, labels=$4, selector=$5, serviceaccount=$6, current=$7, desired=$8, ready=$9, updated=$10, available=$11, enabled=1, updatetime = $12 WHERE uid = $13 and enabled=1
`

const INSERT_DAEMONSET_INFO = `
INSERT INTO ` + TB_KUBE_DS_INFO + ` (nsuid, dsname, uid, starttime, labels, selector, serviceaccount, current, desired, ready, updated, available, enabled, createdtime, updatetime, clusterid)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
`

const INSERT_UNNEST_DAEMONSET_INFO = `
	INSERT INTO ` + TB_KUBE_DS_INFO + ` (clusterid, nsuid, dsname, uid, starttime, labels, selector, serviceaccount, current, desired, ready, updated, available, enabled, createdtime, updatetime)
	(select * from  unnest($1::int[], $2::text[], $3::text[], $4::text[], $5::bigint[], $6::text[], $7::text[], $8::text[], $9::bigint[], $10::bigint[], $11::bigint[], $12::bigint[], $13::bigint[], $14::int[], $15::bigint[], $16::bigint[]))
`

const UPDATE_REPLICASET_INFO = `
	UPDATE ` + TB_KUBE_RS_INFO + ` SET 
	clusterid = $15, nsuid=$1, rsname=$2, starttime=$3, labels=$4, selector=$5, replicas=$6, fullylabeledrs=$7, readyrs=$8, availablers=$9, observedgen=$10, refkind=$11, refuid=$12, enabled=1, updatetime = $13
	 WHERE uid = $14 and enabled=1
`

const INSERT_REPLICASET_INFO = `
INSERT INTO ` + TB_KUBE_RS_INFO + ` (nsuid, rsname, uid, starttime, labels, selector, replicas, fullylabeledrs, readyrs, availablers, observedgen, refkind, refuid, enabled, createdtime, updatetime, clusterid)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
`

const INSERT_UNNEST_REPLICASET_INFO = `
	INSERT INTO ` + TB_KUBE_RS_INFO + ` (clusterid, nsuid, rsname, uid, starttime, labels, selector, replicas, fullylabeledrs, readyrs, availablers, observedgen, refkind, refuid, enabled, createdtime, updatetime)
	(select * from  unnest($1::int[], $2::text[], $3::text[], $4::text[], $5::bigint[], $6::text[], $7::text[], $8::bigint[], $9::bigint[], $10::bigint[], $11::bigint[], $12::bigint[], $13::text[], $14::text[], $15::int[], $16::bigint[], $17::bigint[]))
`

const UPDATE_PVC_INFO = `
	UPDATE ` + TB_KUBE_PVC_INFO + ` SET 
	clusterid = $12, nsuid=$1, pvcname=$2, starttime=$3, labels=$4, selector=$5, accessmodes=$6, reqstorage=$7, status=$8, scuid=$9, enabled=1, updatetime = $10
	 WHERE uid = $11 and enabled=1
`

const INSERT_PVC_INFO = `
INSERT INTO ` + TB_KUBE_PVC_INFO + ` (nsuid, pvcname, uid, starttime, labels, selector, accessmodes, reqstorage, status, scuid, enabled, createdtime, updatetime, clusterid)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
`

const INSERT_UNNEST_PVC_INFO = `
	INSERT INTO ` + TB_KUBE_PVC_INFO + ` (clusterid, nsuid, pvcname, uid, starttime, labels, selector, accessmodes, reqstorage, status, scuid, enabled, createdtime, updatetime)
	(select * from  unnest($1::int[], $2::text[], $3::text[], $4::text[], $5::bigint[], $6::text[], $7::text[], $8::text[], $9::bigint[], $10::text[], $11::text[], $12::int[], $13::bigint[], $14::bigint[]))
`

const UPDATE_PV_INFO = `
	UPDATE ` + TB_KUBE_PV_INFO + ` SET 
	clusterid = $11, pvname=$1, pvcuid=$2, starttime=$3, labels=$4, accessmodes=$5, capacity=$6, reclaimpolicy=$7, status=$8, enabled=1, updatetime = $9
	 WHERE pvuid = $10 and enabled=1
`
const INSERT_PV_INFO = `
INSERT INTO ` + TB_KUBE_PV_INFO + ` (pvname, pvuid, pvcuid, starttime, labels, accessmodes, capacity, reclaimpolicy, status, enabled, createdtime, updatetime, clusterid)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
`

const INSERT_UNNEST_PV_INFO = `
	INSERT INTO ` + TB_KUBE_PV_INFO + ` (clusterid, pvname, pvuid, pvcuid, starttime, labels, accessmodes, capacity, reclaimpolicy, status, enabled, createdtime, updatetime)
	(select * from  unnest($1::int[], $2::text[], $3::text[], $4::text[], $5::bigint[], $6::text[], $7::text[], $8::bigint[], $9::text[], $10::text[], $11::int[], $12::bigint[], $13::bigint[]))
`

const UPDATE_SC_INFO = `
	UPDATE ` + TB_KUBE_SC_INFO + ` SET 
	clusterid=$1, scname=$2, starttime=$3, labels=$4, provisioner=$5, reclaimpolicy=$6, volumebindingmode=$7, allowvolumeexp=$8, enabled=1, updatetime = $9
	 WHERE uid = $10 and enabled=1
`

const INSERT_SC_INFO = `
INSERT INTO ` + TB_KUBE_SC_INFO + ` (clusterid, scname, uid, starttime, labels, provisioner, reclaimpolicy, volumebindingmode, allowvolumeexp, enabled, createdtime, updatetime)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
`

const INSERT_FSDEVICE_INFO = `
INSERT INTO ` + TB_KUBE_FS_DEVICE_INFO + ` (devicename, createdtime, updatetime)
VALUES ($1, $2, $3)
`

const INSERT_NETINTERFACE_INFO = `
INSERT INTO ` + TB_KUBE_NET_INTERFACE_INFO + ` (interfacename, createdtime, updatetime)
VALUES ($1, $2, $3)
`

const INSERT_METRICID_INFO = `
INSERT INTO ` + TB_KUBE_METRIC_ID_INFO + ` (metricname, image, createdtime, updatetime)
VALUES ($1, $2, $3, $4)
`

const INSERT_UNNEST_SC_INFO = `
	INSERT INTO ` + TB_KUBE_SC_INFO + ` (clusterid, scname, uid, starttime, labels, provisioner, reclaimpolicy, volumebindingmode, allowvolumeexp, enabled, createdtime, updatetime)
	(select * from  unnest($1::int[], $2::text[], $3::text[], $4::bigint[], $5::text[], $6::text[], $7::text[], $8::text[], $9::int[], $10::int[], $11::bigint[], $12::bigint[]))
`

const UPDATE_ING_INFO = `
	UPDATE ` + TB_KUBE_ING_INFO + ` SET 
	clusterid=$8, nsuid=$1, ingname=$2, starttime=$3, labels=$4, classname=$5, enabled=1, updatetime = $6
	 WHERE uid = $7 and enabled=1
`
const INSERT_ING_INFO = `
INSERT INTO ` + TB_KUBE_ING_INFO + ` (nsuid, ingname, uid, starttime, labels, classname, enabled, createdtime, updatetime, clusterid)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

const INSERT_UNNEST_ING_INFO = `
	INSERT INTO ` + TB_KUBE_ING_INFO + ` (clusterid, nsuid, ingname, uid, starttime, labels, classname, enabled, createdtime, updatetime)
	(select * from  unnest($1::int[], $2::text[], $3::text[], $4::text[], $5::bigint[], $6::text[], $7::text[], $8::int[], $9::bigint[], $10::bigint[]))
`

const UPDATE_INGHOST_INFO = `
	UPDATE ` + TB_KUBE_INGHOST_INFO + ` SET 
	clusterid=$11, inguid=$1, backendtype=$2, backendname=$3, pathtype=$4, path=$5, serviceport=$6, rscapigroup=$7, rsckind=$8, enabled=1, updatetime = $9
	 WHERE hostname = $10 and enabled=1
`

const INSERT_INGHOST_INFO = `
INSERT INTO ` + TB_KUBE_INGHOST_INFO + ` (inguid, backendtype, backendname, hostname, pathtype, path, serviceport, rscapigroup, rsckind, enabled, createdtime, updatetime, clusterid)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
`

const INSERT_UNNEST_INGHOST_INFO = `
	INSERT INTO ` + TB_KUBE_INGHOST_INFO + ` (clusterid, inguid, backendtype, backendname, hostname, pathtype, path, serviceport, rscapigroup, rsckind, enabled, createdtime, updatetime)
	(select * from unnest($1::int[], $2::text[], $3::text[], $4::text[], $5::text[], $6::text[], $7::text[], $8::int[], $9::text[], $10::text[], $11::int[], $12::bigint[], $13::bigint[]))
`

const UPDATE_EVENT_INFO = `
	UPDATE ` + TB_KUBE_EVENT_INFO + ` SET clusterid=$1, nsuid=$2, eventname=$3, firsttime=$4, lasttime=$5, labels=$6, eventtype=$7,
	eventcount=$8, objkind=$9, objuid=$10, srccomponent=$11, srchost=$12, reason=$13, message=$14, enabled=1, updatetime=$15 WHERE uid = $16 and enabled=1
`

const INSERT_EVENT_INFO = `
	INSERT INTO ` + TB_KUBE_EVENT_INFO + ` (nsuid, eventname, uid, firsttime, lasttime, labels, eventtype, eventcount, 
		objkind, objuid, srccomponent, srchost, reason, message, enabled, createdtime, updatetime, clusterid)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
`

const INSERT_UNNEST_EVENT_INFO = `
	INSERT INTO ` + TB_KUBE_EVENT_INFO + ` (clusterid, nsuid, eventname, uid, firsttime, lasttime, labels, eventtype, eventcount, 
		objkind, objuid, srccomponent, srchost, reason, message, enabled, createdtime, updatetime)
	(select * from unnest($1::int[], $2::text[], $3::text[], $4::text[], $5::bigint[], $6::bigint[], $7::text[], $8::text[], $9::int[],
		$10::text[], $11::text[], $12::text[], $13::text[], $14::text[], $15::text[], $16::int[], $17::bigint[], $18::bigint[]))
`

const INSERT_UNNEST_LOG_INFO = `
	INSERT INTO ` + TB_KUBE_LOG_INFO + ` (logtype, nsuid, poduid, starttime, message, createdtime, updatetime)
	(select * from  unnest($1::text[], $2::text[], $3::text[], $4::bigint[], $5::text[], $6::bigint[], $7::bigint[]))

`

const INSERT_UNNEST_NODE_STAT = `
	INSERT INTO ` + "%s" + ` (nodeuid, ontunetime, metricid, cpuusage, cpusystem, cpuuser, cpuusagecores, cputotalcores, memoryusage, memoryusagebytes, memorysizebytes, memoryswap, netiorate, netioerrors, netreceiverate, netreceiveerrors, nettransmitrate, nettransmiterrors, fsiorate, fsreadrate, fswriterate, processes)
	(select * from unnest($1::text[], $2::bigint[], $3::int[], $4::double precision[], $5::double precision[], $6::double precision[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::double precision[], $12::double precision[], $13::double precision[], $14::double precision[], $15::double precision[], $16::double precision[], $17::double precision[], $18::double precision[], $19::double precision[], $20::double precision[], $21::double precision[], $22::double precision[]))
`

const INSERT_UNNEST_POD_STAT = `
	INSERT INTO ` + "%s" + ` (poduid, ontunetime, nsuid, nodeuid, metricid, cpuusage, cpusystem, cpuuser, cpuusagecores, cpurequestcores, cpulimitcores, memoryusage, memoryusagebytes, memoryrequestbytes, memorylimitbytes, memoryswap, netiorate, netioerrors, netreceiverate, netreceiveerrors, nettransmitrate, nettransmiterrors, fsiorate, fsreadrate, fswriterate, processes)
	(select * from  unnest($1::text[], $2::bigint[], $3::text[], $4::text[], $5::int[], $6::double precision[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::double precision[], $12::double precision[], $13::double precision[], $14::double precision[], $15::double precision[], $16::double precision[], $17::double precision[], $18::double precision[], $19::double precision[], $20::double precision[], $21::double precision[], $22::double precision[], $23::double precision[], $24::double precision[], $25::double precision[], $26::double precision[]))
`

const INSERT_UNNEST_CLUSTER_STAT = `
	INSERT INTO ` + "%s" + ` (clusterid, ontunetime, podcount, cpuusage, cpusystem, cpuuser, cpuusagecores, cputotalcores, memoryusage, memoryusagebytes, memorysizebytes, memoryswap, netiorate, netioerrors, netreceiverate, netreceiveerrors, nettransmitrate, nettransmiterrors, fsiorate, fsreadrate, fswriterate, processes)
	(select * from  unnest($1::int[], $2::bigint[], $3::int[], $4::double precision[], $5::double precision[], $6::double precision[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::double precision[], $12::double precision[], $13::double precision[], $14::double precision[], $15::double precision[], $16::double precision[], $17::double precision[], $18::double precision[], $19::double precision[], $20::double precision[], $21::double precision[], $22::double precision[]))
`

const INSERT_UNNEST_NS_STAT = `
	INSERT INTO ` + "%s" + ` (nsuid, ontunetime, podcount, cpuusage, cpusystem, cpuuser, cpuusagecores, cputotalcores, memoryusage, memoryusagebytes, memorysizebytes, memoryswap, netiorate, netioerrors, netreceiverate, netreceiveerrors, nettransmitrate, nettransmiterrors, fsiorate, fsreadrate, fswriterate, processes)
	(select * from  unnest($1::text[], $2::bigint[], $3::int[], $4::double precision[], $5::double precision[], $6::double precision[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::double precision[], $12::double precision[], $13::double precision[], $14::double precision[], $15::double precision[], $16::double precision[], $17::double precision[], $18::double precision[], $19::double precision[], $20::double precision[], $21::double precision[], $22::double precision[]))
`

const INSERT_UNNEST_WORKLOAD_STAT = `
	INSERT INTO %s (%s, ontunetime, podcount, cpuusage, cpusystem, cpuuser, cpuusagecores, cputotalcores, memoryusage, memoryusagebytes, memorysizebytes, memoryswap, netiorate, netioerrors, netreceiverate, netreceiveerrors, nettransmitrate, nettransmiterrors, fsiorate, fsreadrate, fswriterate, processes)
	(select * from  unnest($1::text[], $2::bigint[], $3::int[], $4::double precision[], $5::double precision[], $6::double precision[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::double precision[], $12::double precision[], $13::double precision[], $14::double precision[], $15::double precision[], $16::double precision[], $17::double precision[], $18::double precision[], $19::double precision[], $20::double precision[], $21::double precision[], $22::double precision[]))
`

const INSERT_UNNEST_CONTAINER_STAT = `
	INSERT INTO ` + "%s" + ` (containername, ontunetime, poduid, nsuid, nodeuid, metricid, cpuusage, cpusystem, cpuuser, cpuusagecores, cpurequestcores, cpulimitcores, memoryusage, memoryusagebytes, memoryrequestbytes, memorylimitbytes, memoryswap, fsiorate, fsreadrate, fswriterate, processes)
	(select * from  unnest($1::text[], $2::bigint[], $3::text[], $4::text[], $5::text[], $6::int[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::double precision[], $12::double precision[], $13::double precision[], $14::double precision[], $15::double precision[], $16::double precision[], $17::double precision[], $18::double precision[], $19::double precision[], $20::double precision[], $21::double precision[]))
`

const INSERT_UNNEST_NODENET_STAT = `
	INSERT INTO ` + "%s" + ` (nodeuid, ontunetime, metricid, interfaceid, iorate, ioerrors, receiverate, receiveerrors, transmitrate, transmiterrors)
	(select * from  unnest($1::text[], $2::bigint[], $3::int[], $4::int[], $5::double precision[], $6::double precision[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[]))
`

const INSERT_UNNEST_PODNET_STAT = `
	INSERT INTO ` + "%s" + ` (poduid, ontunetime, nsuid, nodeuid, metricid, interfaceid, iorate, ioerrors, receiverate, receiveerrors, transmitrate, transmiterrors)
	(select * from  unnest($1::text[], $2::bigint[], $3::text[], $4::text[], $5::int[], $6::int[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::double precision[], $12::double precision[]))
`

const INSERT_UNNEST_NODEFS_STAT = `
	INSERT INTO ` + "%s" + ` (nodeuid, ontunetime, metricid, deviceid, iorate, readrate, writerate)
	(select * from  unnest($1::text[], $2::bigint[], $3::int[], $4::int[], $5::double precision[], $6::double precision[], $7::double precision[]))
`

const INSERT_UNNEST_PODFS_STAT = `
	INSERT INTO ` + "%s" + ` (poduid, ontunetime, nsuid, nodeuid, metricid, deviceid, iorate, readrate, writerate)
	(select * from  unnest($1::text[], $2::bigint[], $3::text[], $4::text[], $5::int[], $6::int[], $7::double precision[], $8::double precision[], $9::double precision[]))
`

const INSERT_UNNEST_CONTAINERFS_STAT = `
	INSERT INTO ` + "%s" + ` (containername, ontunetime, poduid, nsuid, nodeuid, metricid, deviceid, iorate, readrate, writerate)
	(select * from  unnest($1::text[], $2::bigint[], $3::text[], $4::text[], $5::text[], $6::int[], $7::int[], $8::double precision[], $9::double precision[], $10::double precision[]))
`

const INSERT_UNNEST_NODE_PERF = `
	INSERT INTO ` + "%s" + ` (nodeuid, ontunetime, metricid, cpuusagesecondstotal, cpusystemsecondstotal, cpuusersecondstotal, memoryusagebytes, memoryworkingsetbytes, memorycache, memoryswap, memoryrss, fsreadsbytestotal, fswritesbytestotal, processes, timestampms)
	(select * from  unnest($1::text[], $2::bigint[], $3::int[], $4::double precision[], $5::double precision[], $6::double precision[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::double precision[],  $12::double precision[], $13::double precision[], $14::double precision[], $15::bigint[]))
`

const INSERT_UNNEST_POD_PERF = `
	INSERT INTO ` + "%s" + ` (poduid, ontunetime, nsuid, nodeuid, metricid, cpuusagesecondstotal, cpusystemsecondstotal, cpuusersecondstotal, memoryusagebytes, memoryworkingsetbytes, memorycache, memoryswap, memoryrss, processes, timestampms)
	(select * from  unnest($1::text[], $2::bigint[], $3::text[], $4::text[], $5::int[], $6::double precision[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::double precision[], $12::double precision[],  $13::double precision[], $14::double precision[], $15::bigint[]))
`

const INSERT_UNNEST_CONTAINER_PERF = `
	INSERT INTO ` + "%s" + ` (containername, ontunetime, poduid, nsuid, nodeuid, metricid, cpuusagesecondstotal, cpusystemsecondstotal, cpuusersecondstotal, memoryusagebytes, memoryworkingsetbytes, memorycache, memoryswap, memoryrss, processes, timestampms)
	(select * from  unnest($1::text[], $2::bigint[], $3::text[], $4::text[], $5::text[], $6::int[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::double precision[], $12::double precision[], $13::double precision[],  $14::double precision[], $15::double precision[], $16::bigint[]))
`

const INSERT_UNNEST_NODE_NET_PERF = `
	INSERT INTO ` + "%s" + ` (nodeuid, ontunetime, metricid, interfaceid, networkreceivebytestotal, networkreceiveerrorstotal, networktransmitbytestotal, networktransmiterrorstotal, timestampms)
	(select * from  unnest($1::text[], $2::bigint[], $3::int[], $4::int[], $5::double precision[], $6::double precision[], $7::double precision[], $8::double precision[], $9::bigint[]))
`

const INSERT_UNNEST_POD_NET_PERF = `
	INSERT INTO ` + "%s" + ` (poduid, ontunetime, nsuid, nodeuid, metricid, interfaceid, networkreceivebytestotal, networkreceiveerrorstotal, networktransmitbytestotal, networktransmiterrorstotal, timestampms)
	(select * from  unnest($1::text[], $2::bigint[], $3::text[], $4::text[], $5::int[], $6::int[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::bigint[]))
`

const INSERT_UNNEST_NODE_FS_PERF = `
	INSERT INTO ` + "%s" + ` (nodeuid, ontunetime, metricid, deviceid, fsinodesfree, fsinodestotal, fslimitbytes, fsreadsbytestotal, fswritesbytestotal, fsusagebytes, timestampms)
	(select * from  unnest($1::text[], $2::bigint[], $3::int[], $4::int[], $5::double precision[], $6::double precision[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::bigint[]))
`

const INSERT_UNNEST_POD_FS_PERF = `
	INSERT INTO ` + "%s" + ` (poduid, ontunetime, nsuid, nodeuid, metricid, deviceid, fsreadsbytestotal, fswritesbytestotal, timestampms)
	(select * from  unnest($1::text[], $2::bigint[], $3::text[], $4::text[], $5::int[], $6::int[], $7::double precision[], $8::double precision[], $9::bigint[]))
`

const INSERT_UNNEST_CONTAINER_FS_PERF = `
	INSERT INTO ` + "%s" + ` (containername, ontunetime, poduid, nsuid, nodeuid, metricid, deviceid, fsreadsbytestotal, fswritesbytestotal, timestampms)
	(select * from  unnest($1::text[], $2::bigint[], $3::text[], $4::text[], $5::text[], $6::int[], $7::int[], $8::double precision[], $9::double precision[], $10::bigint[]))
`

const DELETE_DATA = `
	DELETE FROM %s WHERE nodeuid='%s' and ontunetime < %d
`

const DELETE_DATA_CONDITION = `
DELETE FROM %s WHERE %s in (%s) and ontunetime < %d
`

const INSERT_UNNEST_LAST_NODE_PERF_RAW = `
	INSERT INTO ` + TB_KUBE_LAST_NODE_PERF_RAW + ` (nodeuid, ontunetime, metricid, cpuusagesecondstotal, cpusystemsecondstotal, cpuusersecondstotal, memoryusagebytes, memoryworkingsetbytes, memorycache, memoryswap, memoryrss, fsreadsbytestotal, fswritesbytestotal, processes, timestampms)
	(select * from  unnest($1::text[], $2::bigint[], $3::int[], $4::double precision[], $5::double precision[], $6::double precision[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::double precision[], $12::double precision[], $13::double precision[], $14::double precision[], $15::bigint[]))
`

const INSERT_UNNEST_LAST_POD_PERF_RAW = `
	INSERT INTO ` + TB_KUBE_LAST_POD_PERF_RAW + ` (poduid, ontunetime, nsuid, nodeuid, metricid, cpuusagesecondstotal, cpusystemsecondstotal, cpuusersecondstotal, memoryusagebytes, memoryworkingsetbytes, memorycache, memoryswap, memoryrss, processes, timestampms)
	(select * from  unnest($1::text[], $2::bigint[], $3::text[], $4::text[], $5::int[], $6::double precision[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::double precision[], $12::double precision[], $13::double precision[], $14::double precision[], $15::bigint[]))
`

const INSERT_UNNEST_LAST_CONTAINER_PERF_RAW = `
	INSERT INTO ` + TB_KUBE_LAST_CONTAINER_PERF_RAW + ` (containername, ontunetime, poduid, nsuid, nodeuid, metricid, cpuusagesecondstotal, cpusystemsecondstotal, cpuusersecondstotal, memoryusagebytes, memoryworkingsetbytes, memorycache, memoryswap, memoryrss, processes, timestampms)
	(select * from  unnest($1::text[], $2::bigint[], $3::text[], $4::text[], $5::text[], $6::int[], $7::double precision[], $8::double precision[], $9::double precision[], $10::double precision[], $11::double precision[], $12::double precision[], $13::double precision[], $14::double precision[], $15::double precision[], $16::bigint[]))
`
