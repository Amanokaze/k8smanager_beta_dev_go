package database

import (
	"sync"

	"github.com/lib/pq"
)

type DBInfo struct {
	Host   string
	Port   string
	User   string
	Pwd    string
	DBname string
}

type Managerinfo struct {
	ArrManagername []string
	ArrDescription []string
	Arr_IP         []string
}

type Clusterinfo struct {
	ArrManagerid   []int
	ArrClustername []string
	ArrDescription []string
	Arr_IP         []string
}

type Resourceinfo struct {
	ArrClusterid    []int
	ArrResourcename []string
	ArrApiclass     []string
	ArrVersion      []string
	ArrEndpoint     []string
	ArrEnabled      []int
	ArrCreateTime   []int64
	ArrUpdateTime   []int64
}

type Scinfo struct {
	ArrScid              []int
	ArrClusterid         []int
	ArrScname            []string
	ArrUid               []string
	ArrStarttime         []int64
	ArrLabels            []string
	ArrProvisionor       []string
	ArrReclaimolicy      []string
	ArrVolumebindingmode []string
	ArrAllowvolumeexp    []int
	ArrEnabled           []int
	ArrCreateTime        []int64
	ArrUpdateTime        []int64
}

func (s *Scinfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrScname))
	data = append(data, pq.StringArray(s.ArrUid))
	data = append(data, pq.Int64Array(s.ArrStarttime))
	data = append(data, pq.StringArray(s.ArrLabels))
	data = append(data, pq.StringArray(s.ArrProvisionor))
	data = append(data, pq.StringArray(s.ArrReclaimolicy))
	data = append(data, pq.StringArray(s.ArrVolumebindingmode))
	data = append(data, pq.Array(s.ArrAllowvolumeexp))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type Namespaceinfo struct {
	ArrNsid       []int
	ArrNsuid      []string
	ArrClusterid  []int
	ArrNsname     []string
	ArrStarttime  []int64
	ArrLabels     []string
	ArrStatus     []string
	ArrEnabled    []int
	ArrCreateTime []int64
	ArrUpdateTime []int64
}

func (s *Namespaceinfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrNsname))
	data = append(data, pq.Int64Array(s.ArrStarttime))
	data = append(data, pq.StringArray(s.ArrLabels))
	data = append(data, pq.StringArray(s.ArrStatus))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type Nodeinfo struct {
	ArrNodeid              []int
	ArrManagerid           []int
	ArrClusterid           []int
	ArrNodeUid             []string
	ArrNodename            []string
	ArrNodenameext         []string
	ArrNodetype            []string
	ArrEnabled             []int
	ArrStarttime           []int64
	ArrLabels              []string
	ArrKernelversion       []string
	Arr_OSimage            []string
	Arr_OSname             []string
	ArrContainerruntimever []string
	ArrKubeletver          []string
	ArrKubeproxyver        []string
	ArrCpuarch             []string
	ArrCpucount            []int
	ArrEphemeralstorage    []int64
	ArrMemorysize          []int64
	ArrPods                []int64
	Arr_IP                 []string
	ArrStatus              []int
	ArrCreateTime          []int64
	ArrUpdateTime          []int64
}

func (s *Nodeinfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrManagerid))
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrNodeUid))
	data = append(data, pq.StringArray(s.ArrNodename))
	data = append(data, pq.StringArray(s.ArrNodenameext))
	data = append(data, pq.StringArray(s.ArrNodetype))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrStarttime))
	data = append(data, pq.StringArray(s.ArrLabels))
	data = append(data, pq.StringArray(s.ArrKernelversion))
	data = append(data, pq.StringArray(s.Arr_OSimage))
	data = append(data, pq.StringArray(s.Arr_OSname))
	data = append(data, pq.StringArray(s.ArrContainerruntimever))
	data = append(data, pq.StringArray(s.ArrKubeletver))
	data = append(data, pq.StringArray(s.ArrKubeproxyver))
	data = append(data, pq.StringArray(s.ArrCpuarch))
	data = append(data, pq.Array(s.ArrCpucount))
	data = append(data, pq.Array(s.ArrEphemeralstorage))
	data = append(data, pq.Array(s.ArrMemorysize))
	data = append(data, pq.Array(s.ArrPods))
	data = append(data, pq.StringArray(s.Arr_IP))
	data = append(data, pq.Array(s.ArrStatus))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type Podinfo struct {
	ArrPodid          []int
	ArrClusterid      []int
	ArrUid            []string
	ArrNodeUid        []string
	ArrNsUid          []string
	ArrAnnotationuid  []string
	ArrPodname        []string
	ArrStarttime      []int64
	ArrLabels         []string
	ArrSelector       []string
	ArrRestartpolicy  []string
	ArrServiceaccount []string
	ArrStatus         []string
	ArrHostip         []string
	ArrPodip          []string
	ArrRestartcount   []int64
	ArrRestarttime    []int64
	ArrCondition      []string
	ArrStaticpod      []string
	ArrRefkind        []string
	ArrRefuid         []string
	ArrPvcuid         []string
	ArrEnabled        []int
	ArrCreateTime     []int64
	ArrUpdateTime     []int64
}

func (s *Podinfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrUid))
	data = append(data, pq.StringArray(s.ArrNodeUid))
	data = append(data, pq.StringArray(s.ArrNsUid))
	data = append(data, pq.StringArray(s.ArrAnnotationuid))
	data = append(data, pq.StringArray(s.ArrPodname))
	data = append(data, pq.Int64Array(s.ArrStarttime))
	data = append(data, pq.StringArray(s.ArrLabels))
	data = append(data, pq.StringArray(s.ArrSelector))
	data = append(data, pq.StringArray(s.ArrRestartpolicy))
	data = append(data, pq.StringArray(s.ArrServiceaccount))
	data = append(data, pq.StringArray(s.ArrStatus))
	data = append(data, pq.StringArray(s.ArrHostip))
	data = append(data, pq.StringArray(s.ArrPodip))
	data = append(data, pq.Array(s.ArrRestartcount))
	data = append(data, pq.Array(s.ArrRestarttime))
	data = append(data, pq.StringArray(s.ArrCondition))
	data = append(data, pq.StringArray(s.ArrStaticpod))
	data = append(data, pq.StringArray(s.ArrRefkind))
	data = append(data, pq.StringArray(s.ArrRefuid))
	data = append(data, pq.StringArray(s.ArrPvcuid))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type Containerinfo struct {
	ArrContainerid      []int
	ArrClusterid        []int
	ArrPodUid           []string
	ArrContainername    []string
	ArrImage            []string
	ArrPorts            []string
	ArrEnv              []string
	ArrLimitCpu         []int64
	ArrLimitMemory      []int64
	ArrLimitStorage     []int64
	ArrLimitEphemeral   []int64
	ArrRequestCpu       []int64
	ArrRequestMemory    []int64
	ArrRequestStorage   []int64
	ArrRequestEphemeral []int64
	ArrVolumemounts     []string
	ArrState            []string
	ArrEnabled          []int
	ArrCreateTime       []int64
	ArrUpdateTime       []int64
}

func (s *Containerinfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrPodUid))
	data = append(data, pq.StringArray(s.ArrContainername))
	data = append(data, pq.StringArray(s.ArrImage))
	data = append(data, pq.StringArray(s.ArrPorts))
	data = append(data, pq.StringArray(s.ArrEnv))
	data = append(data, pq.Array(s.ArrLimitCpu))
	data = append(data, pq.Array(s.ArrLimitMemory))
	data = append(data, pq.Array(s.ArrLimitStorage))
	data = append(data, pq.Array(s.ArrLimitEphemeral))
	data = append(data, pq.Array(s.ArrRequestCpu))
	data = append(data, pq.Array(s.ArrRequestMemory))
	data = append(data, pq.Array(s.ArrRequestStorage))
	data = append(data, pq.Array(s.ArrRequestEphemeral))
	data = append(data, pq.StringArray(s.ArrVolumemounts))
	data = append(data, pq.StringArray(s.ArrState))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type Serviceinfo struct {
	ArrSvcid       []int
	ArrClusterid   []int
	ArrNsuid       []string
	ArrSvcname     []string
	ArrUid         []string
	ArrStarttime   []int64
	ArrLabels      []string
	ArrSelector    []string
	ArrServicetype []string
	ArrClusterip   []string
	ArrPorts       []string
	ArrEnabled     []int
	ArrCreateTime  []int64
	ArrUpdateTime  []int64
}

func (s *Serviceinfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrSvcname))
	data = append(data, pq.StringArray(s.ArrUid))
	data = append(data, pq.Int64Array(s.ArrStarttime))
	data = append(data, pq.StringArray(s.ArrLabels))
	data = append(data, pq.StringArray(s.ArrSelector))
	data = append(data, pq.StringArray(s.ArrServicetype))
	data = append(data, pq.StringArray(s.ArrClusterip))
	data = append(data, pq.StringArray(s.ArrPorts))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type Deployinfo struct {
	ArrDeployid       []int
	ArrClusterid      []int
	ArrNsuid          []string
	ArrDeployname     []string
	ArrUid            []string
	ArrStarttime      []int64
	ArrLabels         []string
	ArrSelector       []string
	ArrServiceaccount []string
	ArrReplicas       []int64
	ArrUpdatedrs      []int64
	ArrReadyrs        []int64
	ArrAvailablers    []int64
	ArrObservedgen    []int64
	ArrEnabled        []int
	ArrCreateTime     []int64
	ArrUpdateTime     []int64
}

func (s *Deployinfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrDeployname))
	data = append(data, pq.StringArray(s.ArrUid))
	data = append(data, pq.Int64Array(s.ArrStarttime))
	data = append(data, pq.StringArray(s.ArrLabels))
	data = append(data, pq.StringArray(s.ArrSelector))
	data = append(data, pq.StringArray(s.ArrServiceaccount))
	data = append(data, pq.Array(s.ArrReplicas))
	data = append(data, pq.Array(s.ArrUpdatedrs))
	data = append(data, pq.Array(s.ArrReadyrs))
	data = append(data, pq.Array(s.ArrAvailablers))
	data = append(data, pq.Array(s.ArrObservedgen))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type StateFulSetinfo struct {
	ArrDeployid       []int
	ArrClusterid      []int
	ArrNsuid          []string
	ArrStsname        []string
	ArrUid            []string
	ArrStarttime      []int64
	ArrLabels         []string
	ArrSelector       []string
	ArrServiceaccount []string
	ArrReplicas       []int64
	ArrUpdatedrs      []int64
	ArrReadyrs        []int64
	ArrAvailablers    []int64
	ArrEnabled        []int
	ArrCreateTime     []int64
	ArrUpdateTime     []int64
}

func (s *StateFulSetinfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrStsname))
	data = append(data, pq.StringArray(s.ArrUid))
	data = append(data, pq.Int64Array(s.ArrStarttime))
	data = append(data, pq.StringArray(s.ArrLabels))
	data = append(data, pq.StringArray(s.ArrSelector))
	data = append(data, pq.StringArray(s.ArrServiceaccount))
	data = append(data, pq.Array(s.ArrReplicas))
	data = append(data, pq.Array(s.ArrUpdatedrs))
	data = append(data, pq.Array(s.ArrReadyrs))
	data = append(data, pq.Array(s.ArrAvailablers))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type DaemonSetinfo struct {
	ArrDsid           []int
	ArrClusterid      []int
	ArrNsuid          []string
	ArrDsname         []string
	ArrUid            []string
	ArrStarttime      []int64
	ArrLabels         []string
	ArrSelector       []string
	ArrServiceaccount []string
	ArrCurrent        []int64
	ArrDesired        []int64
	ArrReady          []int64
	ArrUpdated        []int64
	ArrAvailable      []int64
	ArrEnabled        []int
	ArrCreateTime     []int64
	ArrUpdateTime     []int64
}

func (s *DaemonSetinfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrDsname))
	data = append(data, pq.StringArray(s.ArrUid))
	data = append(data, pq.Int64Array(s.ArrStarttime))
	data = append(data, pq.StringArray(s.ArrLabels))
	data = append(data, pq.StringArray(s.ArrSelector))
	data = append(data, pq.StringArray(s.ArrServiceaccount))
	data = append(data, pq.Array(s.ArrCurrent))
	data = append(data, pq.Array(s.ArrDesired))
	data = append(data, pq.Array(s.ArrReady))
	data = append(data, pq.Array(s.ArrUpdated))
	data = append(data, pq.Array(s.ArrAvailable))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type ReplicaSetinfo struct {
	ArrRsid          []int
	ArrClusterid     []int
	ArrNsuid         []string
	ArrRsname        []string
	ArrUid           []string
	ArrStarttime     []int64
	ArrLabels        []string
	ArrSelector      []string
	ArrReplicas      []int64
	ArrFullylabeldrs []int64
	ArrReadyrs       []int64
	ArrAvailablers   []int64
	ArrObservedgen   []int64
	ArrRefkind       []string
	ArrRefuid        []string
	ArrEnabled       []int
	ArrCreateTime    []int64
	ArrUpdateTime    []int64
}

func (s *ReplicaSetinfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrRsname))
	data = append(data, pq.StringArray(s.ArrUid))
	data = append(data, pq.Int64Array(s.ArrStarttime))
	data = append(data, pq.StringArray(s.ArrLabels))
	data = append(data, pq.StringArray(s.ArrSelector))
	data = append(data, pq.Array(s.ArrReplicas))
	data = append(data, pq.Array(s.ArrFullylabeldrs))
	data = append(data, pq.Array(s.ArrReadyrs))
	data = append(data, pq.Array(s.ArrAvailablers))
	data = append(data, pq.Array(s.ArrObservedgen))
	data = append(data, pq.StringArray(s.ArrRefkind))
	data = append(data, pq.StringArray(s.ArrRefuid))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type Pvcinfo struct {
	ArrPvcid       []int
	ArrClusterid   []int
	ArrNsuid       []string
	ArrPvcname     []string
	ArrPvcUid      []string
	ArrStarttime   []int64
	ArrLabels      []string
	ArrSelector    []string
	ArrAccessmodes []string
	ArrReqstorage  []int64
	ArrStatus      []string
	ArrScuid       []string
	ArrEnabled     []int
	ArrCreateTime  []int64
	ArrUpdateTime  []int64
}

func (s *Pvcinfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrPvcname))
	data = append(data, pq.StringArray(s.ArrPvcUid))
	data = append(data, pq.Int64Array(s.ArrStarttime))
	data = append(data, pq.StringArray(s.ArrLabels))
	data = append(data, pq.StringArray(s.ArrSelector))
	data = append(data, pq.StringArray(s.ArrAccessmodes))
	data = append(data, pq.Array(s.ArrReqstorage))
	data = append(data, pq.StringArray(s.ArrStatus))
	data = append(data, pq.StringArray(s.ArrScuid))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type Pvinfo struct {
	ArrPvid          []int
	ArrClusterid     []int
	ArrPvname        []string
	ArrPvUid         []string
	ArrPvcUid        []string
	ArrStarttime     []int64
	ArrLabels        []string
	ArrAccessmodes   []string
	ArrCapacity      []int64
	ArrReclaimpolicy []string
	ArrStatus        []string
	ArrEnabled       []int
	ArrCreateTime    []int64
	ArrUpdateTime    []int64
}

func (s *Pvinfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrPvname))
	data = append(data, pq.StringArray(s.ArrPvUid))
	data = append(data, pq.StringArray(s.ArrPvcUid))
	data = append(data, pq.Int64Array(s.ArrStarttime))
	data = append(data, pq.StringArray(s.ArrLabels))
	data = append(data, pq.StringArray(s.ArrAccessmodes))
	data = append(data, pq.Array(s.ArrCapacity))
	data = append(data, pq.StringArray(s.ArrReclaimpolicy))
	data = append(data, pq.StringArray(s.ArrStatus))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type Inginfo struct {
	ArrIngid      []int
	ArrClusterid  []int
	ArrNsUid      []string
	ArrName       []string
	ArrUid        []string
	ArrStarttime  []int64
	ArrLabels     []string
	ArrClassname  []string
	ArrEnabled    []int
	ArrCreateTime []int64
	ArrUpdateTime []int64
}

func (s *Inginfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrNsUid))
	data = append(data, pq.StringArray(s.ArrName))
	data = append(data, pq.StringArray(s.ArrUid))
	data = append(data, pq.Int64Array(s.ArrStarttime))
	data = append(data, pq.StringArray(s.ArrLabels))
	data = append(data, pq.StringArray(s.ArrClassname))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type IngHostinfo struct {
	ArrInghostid   []int
	ArrClusterid   []int
	ArrInguid      []string
	ArrBackendtype []string
	ArrBackendname []string
	ArrHostname    []string
	ArrPathtype    []string
	ArrPath        []string
	ArrSvcport     []int
	ArrRscApiGroup []string
	ArrRsckind     []string
	ArrEnabled     []int
	ArrCreateTime  []int64
	ArrUpdateTime  []int64
}

func (s *IngHostinfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrInguid))
	data = append(data, pq.StringArray(s.ArrBackendtype))
	data = append(data, pq.StringArray(s.ArrBackendname))
	data = append(data, pq.StringArray(s.ArrHostname))
	data = append(data, pq.StringArray(s.ArrPathtype))
	data = append(data, pq.StringArray(s.ArrPath))
	data = append(data, pq.Array(s.ArrSvcport))
	data = append(data, pq.StringArray(s.ArrRscApiGroup))
	data = append(data, pq.StringArray(s.ArrRsckind))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type Eventinfo struct {
	ArrEventid      []int
	ArrClusterid    []int
	ArrNsUid        []string
	ArrEventname    []string
	ArrEventUid     []string
	ArrFirsttime    []int64
	ArrLasttime     []int64
	ArrLabels       []string
	ArrEventtype    []string
	ArrEventcount   []int
	ArrObjkind      []string
	ArrObjUid       []string
	ArrSrccomponent []string
	ArrSrchost      []string
	ArrReason       []string
	ArrMessage      []string
	ArrEnabled      []int
	ArrCreateTime   []int64
	ArrUpdateTime   []int64
}

func (s *Eventinfo) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.StringArray(s.ArrNsUid))
	data = append(data, pq.StringArray(s.ArrEventname))
	data = append(data, pq.StringArray(s.ArrEventUid))
	data = append(data, pq.Int64Array(s.ArrFirsttime))
	data = append(data, pq.Int64Array(s.ArrLasttime))
	data = append(data, pq.StringArray(s.ArrLabels))
	data = append(data, pq.StringArray(s.ArrEventtype))
	data = append(data, pq.Array(s.ArrEventcount))
	data = append(data, pq.StringArray(s.ArrObjkind))
	data = append(data, pq.StringArray(s.ArrObjUid))
	data = append(data, pq.StringArray(s.ArrSrccomponent))
	data = append(data, pq.StringArray(s.ArrSrchost))
	data = append(data, pq.StringArray(s.ArrReason))
	data = append(data, pq.StringArray(s.ArrMessage))
	data = append(data, pq.Array(s.ArrEnabled))
	data = append(data, pq.Int64Array(s.ArrCreateTime))
	data = append(data, pq.Int64Array(s.ArrUpdateTime))

	return data
}

type Loginfo struct {
	ArrLogid      []int
	ArrLogType    []string
	ArrNsUid      []string
	ArrPodUid     []string
	ArrStarttime  []int64
	ArrMessage    []string
	ArrCreateTime []int64
	ArrUpdateTime []int64
}

type FsDeviceinfo struct {
	ArrDeviceid   []int
	ArrDevicename []string
	ArrCreateTime []int64
	ArrUpdateTime []int64
}

type NetInterfaceinfo struct {
	ArrInterfaceid   []int
	ArrInterfacename []string
	ArrCreateTime    []int64
	ArrUpdateTime    []int64
}

type Metricidinfo struct {
	ArrMetricid   []int
	ArrMetricname []string
	ArrCreateTime []int64
	ArrUpdateTime []int64
}

type NodePerf struct {
	ArrNodeuid               []string
	ArrOntunetime            []int64
	ArrMetricid              []int
	ArrCpuusagesecondstotal  []float64
	ArrCpusystemsecondstotal []float64
	ArrCpuusersecondstotal   []float64
	ArrMemoryusagebytes      []float64
	ArrMemoryworkingsetbytes []float64
	ArrMemorycache           []float64
	ArrMemoryswap            []float64
	ArrMemoryrss             []float64
	ArrFsreadsbytestotal     []float64
	ArrFswritesbytestotal    []float64
	ArrProcesses             []float64
	ArrTimestampMs           []int64
}

func (s *NodePerf) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Array(s.ArrCpuusagesecondstotal))
	data = append(data, pq.Array(s.ArrCpusystemsecondstotal))
	data = append(data, pq.Array(s.ArrCpuusersecondstotal))
	data = append(data, pq.Array(s.ArrMemoryusagebytes))
	data = append(data, pq.Array(s.ArrMemoryworkingsetbytes))
	data = append(data, pq.Array(s.ArrMemorycache))
	data = append(data, pq.Array(s.ArrMemoryswap))
	data = append(data, pq.Array(s.ArrMemoryrss))
	data = append(data, pq.Array(s.ArrFsreadsbytestotal))
	data = append(data, pq.Array(s.ArrFswritesbytestotal))
	data = append(data, pq.Array(s.ArrProcesses))
	data = append(data, pq.Int64Array(s.ArrTimestampMs))

	return data
}

type PodPerf struct {
	ArrPoduid                []string
	ArrOntunetime            []int64
	ArrNsuid                 []string
	ArrNodeuid               []string
	ArrMetricid              []int
	ArrCpuusagesecondstotal  []float64
	ArrCpusystemsecondstotal []float64
	ArrCpuusersecondstotal   []float64
	ArrMemoryusagebytes      []float64
	ArrMemoryworkingsetbytes []float64
	ArrMemorycache           []float64
	ArrMemoryswap            []float64
	ArrMemoryrss             []float64
	ArrProcesses             []float64
	ArrTimestampMs           []int64
}

func (s *PodPerf) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrPoduid))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Array(s.ArrCpuusagesecondstotal))
	data = append(data, pq.Array(s.ArrCpusystemsecondstotal))
	data = append(data, pq.Array(s.ArrCpuusersecondstotal))
	data = append(data, pq.Array(s.ArrMemoryusagebytes))
	data = append(data, pq.Array(s.ArrMemoryworkingsetbytes))
	data = append(data, pq.Array(s.ArrMemorycache))
	data = append(data, pq.Array(s.ArrMemoryswap))
	data = append(data, pq.Array(s.ArrMemoryrss))
	data = append(data, pq.Array(s.ArrProcesses))
	data = append(data, pq.Int64Array(s.ArrTimestampMs))

	return data
}

type ContainerPerf struct {
	ArrContainername         []string
	ArrOntunetime            []int64
	ArrPoduid                []string
	ArrNsuid                 []string
	ArrNodeuid               []string
	ArrMetricid              []int
	ArrCpuusagesecondstotal  []float64
	ArrCpusystemsecondstotal []float64
	ArrCpuusersecondstotal   []float64
	ArrMemoryusagebytes      []float64
	ArrMemoryworkingsetbytes []float64
	ArrMemorycache           []float64
	ArrMemoryswap            []float64
	ArrMemoryrss             []float64
	ArrProcesses             []float64
	ArrTimestampMs           []int64
}

func (s *ContainerPerf) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrContainername))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.StringArray(s.ArrPoduid))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Array(s.ArrCpuusagesecondstotal))
	data = append(data, pq.Array(s.ArrCpusystemsecondstotal))
	data = append(data, pq.Array(s.ArrCpuusersecondstotal))
	data = append(data, pq.Array(s.ArrMemoryusagebytes))
	data = append(data, pq.Array(s.ArrMemoryworkingsetbytes))
	data = append(data, pq.Array(s.ArrMemorycache))
	data = append(data, pq.Array(s.ArrMemoryswap))
	data = append(data, pq.Array(s.ArrMemoryrss))
	data = append(data, pq.Array(s.ArrProcesses))
	data = append(data, pq.Int64Array(s.ArrTimestampMs))

	return data
}

type NodeNetPerf struct {
	ArrNodeuid                    []string
	ArrOntunetime                 []int64
	ArrMetricid                   []int
	ArrInterfaceid                []int
	ArrNetworkreceivebytestotal   []float64
	ArrNetworkreceiveerrorstotal  []float64
	ArrNetworktransmitbytestotal  []float64
	ArrNetworktransmiterrorstotal []float64
	ArrTimestampMs                []int64
}

func (s *NodeNetPerf) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Array(s.ArrInterfaceid))
	data = append(data, pq.Array(s.ArrNetworkreceivebytestotal))
	data = append(data, pq.Array(s.ArrNetworkreceiveerrorstotal))
	data = append(data, pq.Array(s.ArrNetworktransmitbytestotal))
	data = append(data, pq.Array(s.ArrNetworktransmiterrorstotal))
	data = append(data, pq.Int64Array(s.ArrTimestampMs))

	return data
}

type PodNetPerf struct {
	ArrPoduid                     []string
	ArrOntunetime                 []int64
	ArrNsuid                      []string
	ArrNodeuid                    []string
	ArrMetricid                   []int
	ArrInterfaceid                []int
	ArrNetworkreceivebytestotal   []float64
	ArrNetworkreceiveerrorstotal  []float64
	ArrNetworktransmitbytestotal  []float64
	ArrNetworktransmiterrorstotal []float64
	ArrTimestampMs                []int64
}

func (s *PodNetPerf) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrPoduid))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Array(s.ArrInterfaceid))
	data = append(data, pq.Array(s.ArrNetworkreceivebytestotal))
	data = append(data, pq.Array(s.ArrNetworkreceiveerrorstotal))
	data = append(data, pq.Array(s.ArrNetworktransmitbytestotal))
	data = append(data, pq.Array(s.ArrNetworktransmiterrorstotal))
	data = append(data, pq.Int64Array(s.ArrTimestampMs))

	return data
}

type NodeFsPerf struct {
	ArrNodeuid            []string
	ArrOntunetime         []int64
	ArrMetricid           []int
	ArrDeviceid           []int
	ArrFsinodesfree       []float64
	ArrFsinodestotal      []float64
	ArrFslimitbytes       []float64
	ArrFsreadsbytestotal  []float64
	ArrFswritesbytestotal []float64
	ArrFsusagebytes       []float64
	ArrTimestampMs        []int64
}

func (s *NodeFsPerf) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Array(s.ArrDeviceid))
	data = append(data, pq.Array(s.ArrFsinodesfree))
	data = append(data, pq.Array(s.ArrFsinodestotal))
	data = append(data, pq.Array(s.ArrFslimitbytes))
	data = append(data, pq.Array(s.ArrFsreadsbytestotal))
	data = append(data, pq.Array(s.ArrFswritesbytestotal))
	data = append(data, pq.Array(s.ArrFsusagebytes))
	data = append(data, pq.Int64Array(s.ArrTimestampMs))

	return data
}

type PodFsPerf struct {
	ArrPoduid             []string
	ArrOntunetime         []int64
	ArrNsuid              []string
	ArrNodeuid            []string
	ArrMetricid           []int
	ArrDeviceid           []int
	ArrFsreadsbytestotal  []float64
	ArrFswritesbytestotal []float64
	ArrTimestampMs        []int64
}

func (s *PodFsPerf) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrPoduid))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Array(s.ArrDeviceid))
	data = append(data, pq.Array(s.ArrFsreadsbytestotal))
	data = append(data, pq.Array(s.ArrFswritesbytestotal))
	data = append(data, pq.Int64Array(s.ArrTimestampMs))

	return data
}

type ContainerFsPerf struct {
	ArrContainername      []string
	ArrOntunetime         []int64
	ArrPoduid             []string
	ArrNsuid              []string
	ArrNodeuid            []string
	ArrMetricid           []int
	ArrDeviceid           []int
	ArrFsreadsbytestotal  []float64
	ArrFswritesbytestotal []float64
	ArrTimestampMs        []int64
}

func (s *ContainerFsPerf) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrContainername))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.StringArray(s.ArrPoduid))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Array(s.ArrDeviceid))
	data = append(data, pq.Array(s.ArrFsreadsbytestotal))
	data = append(data, pq.Array(s.ArrFswritesbytestotal))
	data = append(data, pq.Int64Array(s.ArrTimestampMs))

	return data
}

type PodSummaryPerf struct {
	ArrPoduid                []string
	ArrOntunetime            []int64
	ArrTimestampMs           []int64
	ArrNsuid                 []string
	ArrNodeuid               []string
	ArrMetricid              []int
	ArrCpuusagesecondstotal  []float64
	ArrCpusystemsecondstotal []float64
	ArrCpuusersecondstotal   []float64
	ArrMemoryusagebytes      []float64
	ArrMemoryworkingsetbytes []float64
	ArrMemorycache           []float64
	ArrMemoryswap            []float64
	ArrMemoryrss             []float64
	ArrProcesses             []float64
	ArrCpucount              []float64
	ArrMemorysize            []float64
	ArrCpurequest            []float64
	ArrMemoryrequest         []float64
	ArrCpulimit              []float64
	ArrMemorylimit           []float64
}

type SummaryPerfInterface interface {
	GetSummaryMap() *sync.Map
}

func (s *NodeNetPerf) GetSummaryMap() *sync.Map {
	stat_map := &sync.Map{}
	stat_map.Store("nodeuid", s.ArrNodeuid)
	stat_map.Store("rbt", s.ArrNetworkreceivebytestotal)
	stat_map.Store("rbe", s.ArrNetworkreceiveerrorstotal)
	stat_map.Store("tbt", s.ArrNetworktransmitbytestotal)
	stat_map.Store("tbe", s.ArrNetworktransmiterrorstotal)

	return stat_map
}

func (s *PodNetPerf) GetSummaryMap() *sync.Map {
	stat_map := &sync.Map{}
	stat_map.Store("poduid", s.ArrPoduid)
	stat_map.Store("nsuid", s.ArrNsuid)
	stat_map.Store("rbt", s.ArrNetworkreceivebytestotal)
	stat_map.Store("rbe", s.ArrNetworkreceiveerrorstotal)
	stat_map.Store("tbt", s.ArrNetworktransmitbytestotal)
	stat_map.Store("tbe", s.ArrNetworktransmiterrorstotal)

	return stat_map
}

func (s *NodeFsPerf) GetSummaryMap() *sync.Map {
	stat_map := &sync.Map{}
	stat_map.Store("nodeuid", s.ArrNodeuid)
	stat_map.Store("rbt", s.ArrFsreadsbytestotal)
	stat_map.Store("wbt", s.ArrFswritesbytestotal)

	return stat_map
}

func (s *PodFsPerf) GetSummaryMap() *sync.Map {
	stat_map := &sync.Map{}
	stat_map.Store("poduid", s.ArrPoduid)
	stat_map.Store("nsuid", s.ArrNsuid)
	stat_map.Store("rbt", s.ArrFsreadsbytestotal)
	stat_map.Store("wbt", s.ArrFswritesbytestotal)

	return stat_map
}

func (s *ContainerFsPerf) GetSummaryMap() *sync.Map {
	stat_map := &sync.Map{}
	stat_map.Store("containername", s.ArrContainername)
	stat_map.Store("poduid", s.ArrPoduid)
	stat_map.Store("rbt", s.ArrFsreadsbytestotal)
	stat_map.Store("wbt", s.ArrFswritesbytestotal)

	return stat_map
}

type NodePerfStat struct {
	ArrNodeuid               []string
	ArrOntunetime            []int64
	ArrMetricid              []int
	ArrCpuusage              []float64
	ArrCpusystem             []float64
	ArrCpuuser               []float64
	ArrCpuusagecores         []float64
	ArrCputotalcores         []float64
	ArrMemoryusage           []float64
	ArrMemoryusagebytes      []float64
	ArrMemorysizebytes       []float64
	ArrMemoryswap            []float64
	ArrNetworkiorate         []float64
	ArrNetworkioerrors       []float64
	ArrNetworkreceiverate    []float64
	ArrNetworkreceiveerrors  []float64
	ArrNetworktransmitrate   []float64
	ArrNetworktransmiterrors []float64
	ArrFsiorate              []float64
	ArrFsreadrate            []float64
	ArrFswriterate           []float64
	ArrProcesses             []float64
}

func (s *NodePerfStat) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Float64Array(s.ArrCpuusage))
	data = append(data, pq.Float64Array(s.ArrCpusystem))
	data = append(data, pq.Float64Array(s.ArrCpuuser))
	data = append(data, pq.Float64Array(s.ArrCpuusagecores))
	data = append(data, pq.Float64Array(s.ArrCputotalcores))
	data = append(data, pq.Float64Array(s.ArrMemoryusage))
	data = append(data, pq.Float64Array(s.ArrMemoryusagebytes))
	data = append(data, pq.Float64Array(s.ArrMemorysizebytes))
	data = append(data, pq.Float64Array(s.ArrMemoryswap))
	data = append(data, pq.Float64Array(s.ArrNetworkiorate))
	data = append(data, pq.Float64Array(s.ArrNetworkioerrors))
	data = append(data, pq.Float64Array(s.ArrNetworkreceiverate))
	data = append(data, pq.Float64Array(s.ArrNetworkreceiveerrors))
	data = append(data, pq.Float64Array(s.ArrNetworktransmitrate))
	data = append(data, pq.Float64Array(s.ArrNetworktransmiterrors))
	data = append(data, pq.Float64Array(s.ArrFsiorate))
	data = append(data, pq.Float64Array(s.ArrFsreadrate))
	data = append(data, pq.Float64Array(s.ArrFswriterate))
	data = append(data, pq.Float64Array(s.ArrProcesses))

	return data
}

func (s *NodePerfStat) Init() {
	s.ArrNodeuid = make([]string, 0)
	s.ArrOntunetime = make([]int64, 0)
	s.ArrMetricid = make([]int, 0)
	s.ArrCpuusage = make([]float64, 0)
	s.ArrCpusystem = make([]float64, 0)
	s.ArrCpuuser = make([]float64, 0)
	s.ArrCpuusagecores = make([]float64, 0)
	s.ArrCputotalcores = make([]float64, 0)
	s.ArrMemoryusage = make([]float64, 0)
	s.ArrMemoryusagebytes = make([]float64, 0)
	s.ArrMemorysizebytes = make([]float64, 0)
	s.ArrMemoryswap = make([]float64, 0)
	s.ArrNetworkiorate = make([]float64, 0)
	s.ArrNetworkioerrors = make([]float64, 0)
	s.ArrNetworkreceiverate = make([]float64, 0)
	s.ArrNetworkreceiveerrors = make([]float64, 0)
	s.ArrNetworktransmitrate = make([]float64, 0)
	s.ArrNetworktransmiterrors = make([]float64, 0)
	s.ArrFsiorate = make([]float64, 0)
	s.ArrFsreadrate = make([]float64, 0)
	s.ArrFswriterate = make([]float64, 0)
	s.ArrProcesses = make([]float64, 0)
}

type PodPerfStat struct {
	ArrPoduid                []string
	ArrOntunetime            []int64
	ArrNsuid                 []string
	ArrNodeuid               []string
	ArrMetricid              []int
	ArrCpuusage              []float64
	ArrCpusystem             []float64
	ArrCpuuser               []float64
	ArrCpuusagecores         []float64
	ArrCpurequestcores       []float64
	ArrCpulimitcores         []float64
	ArrMemoryusage           []float64
	ArrMemoryusagebytes      []float64
	ArrMemoryrequestbytes    []float64
	ArrMemorylimitbytes      []float64
	ArrMemoryswap            []float64
	ArrNetworkiorate         []float64
	ArrNetworkioerrors       []float64
	ArrNetworkreceiverate    []float64
	ArrNetworkreceiveerrors  []float64
	ArrNetworktransmitrate   []float64
	ArrNetworktransmiterrors []float64
	ArrFsiorate              []float64
	ArrFsreadrate            []float64
	ArrFswriterate           []float64
	ArrProcesses             []float64
}

func (s *PodPerfStat) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrPoduid))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Float64Array(s.ArrCpuusage))
	data = append(data, pq.Float64Array(s.ArrCpusystem))
	data = append(data, pq.Float64Array(s.ArrCpuuser))
	data = append(data, pq.Float64Array(s.ArrCpuusagecores))
	data = append(data, pq.Float64Array(s.ArrCpurequestcores))
	data = append(data, pq.Float64Array(s.ArrCpulimitcores))
	data = append(data, pq.Float64Array(s.ArrMemoryusage))
	data = append(data, pq.Float64Array(s.ArrMemoryusagebytes))
	data = append(data, pq.Float64Array(s.ArrMemoryrequestbytes))
	data = append(data, pq.Float64Array(s.ArrMemorylimitbytes))
	data = append(data, pq.Float64Array(s.ArrMemoryswap))
	data = append(data, pq.Float64Array(s.ArrNetworkiorate))
	data = append(data, pq.Float64Array(s.ArrNetworkioerrors))
	data = append(data, pq.Float64Array(s.ArrNetworkreceiverate))
	data = append(data, pq.Float64Array(s.ArrNetworkreceiveerrors))
	data = append(data, pq.Float64Array(s.ArrNetworktransmitrate))
	data = append(data, pq.Float64Array(s.ArrNetworktransmiterrors))
	data = append(data, pq.Float64Array(s.ArrFsiorate))
	data = append(data, pq.Float64Array(s.ArrFsreadrate))
	data = append(data, pq.Float64Array(s.ArrFswriterate))
	data = append(data, pq.Float64Array(s.ArrProcesses))

	return data
}

func (s *PodPerfStat) Init() {
	s.ArrPoduid = make([]string, 0)
	s.ArrOntunetime = make([]int64, 0)
	s.ArrNsuid = make([]string, 0)
	s.ArrNodeuid = make([]string, 0)
	s.ArrMetricid = make([]int, 0)
	s.ArrCpuusage = make([]float64, 0)
	s.ArrCpusystem = make([]float64, 0)
	s.ArrCpuuser = make([]float64, 0)
	s.ArrCpuusagecores = make([]float64, 0)
	s.ArrCpurequestcores = make([]float64, 0)
	s.ArrCpulimitcores = make([]float64, 0)
	s.ArrMemoryusage = make([]float64, 0)
	s.ArrMemoryusagebytes = make([]float64, 0)
	s.ArrMemoryrequestbytes = make([]float64, 0)
	s.ArrMemorylimitbytes = make([]float64, 0)
	s.ArrMemoryswap = make([]float64, 0)
	s.ArrNetworkiorate = make([]float64, 0)
	s.ArrNetworkioerrors = make([]float64, 0)
	s.ArrNetworkreceiverate = make([]float64, 0)
	s.ArrNetworkreceiveerrors = make([]float64, 0)
	s.ArrNetworktransmitrate = make([]float64, 0)
	s.ArrNetworktransmiterrors = make([]float64, 0)
	s.ArrFsiorate = make([]float64, 0)
	s.ArrFsreadrate = make([]float64, 0)
	s.ArrFswriterate = make([]float64, 0)
	s.ArrProcesses = make([]float64, 0)
}

type ContainerPerfStat struct {
	ArrContainername      []string
	ArrOntunetime         []int64
	ArrPoduid             []string
	ArrNsuid              []string
	ArrNodeuid            []string
	ArrMetricid           []int
	ArrCpuusage           []float64
	ArrCpusystem          []float64
	ArrCpuuser            []float64
	ArrCpuusagecores      []float64
	ArrCpurequestcores    []float64
	ArrCpulimitcores      []float64
	ArrMemoryusage        []float64
	ArrMemoryusagebytes   []float64
	ArrMemoryrequestbytes []float64
	ArrMemorylimitbytes   []float64
	ArrMemoryswap         []float64
	ArrFsiorate           []float64
	ArrFsreadrate         []float64
	ArrFswriterate        []float64
	ArrProcesses          []float64
}

func (s *ContainerPerfStat) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrContainername))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.StringArray(s.ArrPoduid))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Float64Array(s.ArrCpuusage))
	data = append(data, pq.Float64Array(s.ArrCpusystem))
	data = append(data, pq.Float64Array(s.ArrCpuuser))
	data = append(data, pq.Float64Array(s.ArrCpuusagecores))
	data = append(data, pq.Float64Array(s.ArrCpurequestcores))
	data = append(data, pq.Float64Array(s.ArrCpulimitcores))
	data = append(data, pq.Float64Array(s.ArrMemoryusage))
	data = append(data, pq.Float64Array(s.ArrMemoryusagebytes))
	data = append(data, pq.Float64Array(s.ArrMemoryrequestbytes))
	data = append(data, pq.Float64Array(s.ArrMemorylimitbytes))
	data = append(data, pq.Float64Array(s.ArrMemoryswap))
	data = append(data, pq.Float64Array(s.ArrFsiorate))
	data = append(data, pq.Float64Array(s.ArrFsreadrate))
	data = append(data, pq.Float64Array(s.ArrFswriterate))
	data = append(data, pq.Float64Array(s.ArrProcesses))

	return data
}

func (s *ContainerPerfStat) Init() {
	s.ArrContainername = make([]string, 0)
	s.ArrPoduid = make([]string, 0)
	s.ArrNodeuid = make([]string, 0)
	s.ArrNsuid = make([]string, 0)
	s.ArrOntunetime = make([]int64, 0)
	s.ArrMetricid = make([]int, 0)
	s.ArrCpuusage = make([]float64, 0)
	s.ArrCpusystem = make([]float64, 0)
	s.ArrCpuuser = make([]float64, 0)
	s.ArrCpuusagecores = make([]float64, 0)
	s.ArrCpurequestcores = make([]float64, 0)
	s.ArrCpulimitcores = make([]float64, 0)
	s.ArrMemoryusage = make([]float64, 0)
	s.ArrMemoryusagebytes = make([]float64, 0)
	s.ArrMemoryrequestbytes = make([]float64, 0)
	s.ArrMemorylimitbytes = make([]float64, 0)
	s.ArrMemoryswap = make([]float64, 0)
	s.ArrFsiorate = make([]float64, 0)
	s.ArrFsreadrate = make([]float64, 0)
	s.ArrFswriterate = make([]float64, 0)
	s.ArrProcesses = make([]float64, 0)
}

type NodeNetPerfStat struct {
	ArrNodeuid               []string
	ArrOntunetime            []int64
	ArrMetricid              []int
	ArrInterfaceid           []int
	ArrNetworkiorate         []float64
	ArrNetworkioerrors       []float64
	ArrNetworkreceiverate    []float64
	ArrNetworkreceiveerrors  []float64
	ArrNetworktransmitrate   []float64
	ArrNetworktransmiterrors []float64
}

func (s *NodeNetPerfStat) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Array(s.ArrInterfaceid))
	data = append(data, pq.Float64Array(s.ArrNetworkiorate))
	data = append(data, pq.Float64Array(s.ArrNetworkioerrors))
	data = append(data, pq.Float64Array(s.ArrNetworkreceiverate))
	data = append(data, pq.Float64Array(s.ArrNetworkreceiveerrors))
	data = append(data, pq.Float64Array(s.ArrNetworktransmitrate))
	data = append(data, pq.Float64Array(s.ArrNetworktransmiterrors))

	return data
}

func (s *NodeNetPerfStat) Init() {
	s.ArrNodeuid = make([]string, 0)
	s.ArrOntunetime = make([]int64, 0)
	s.ArrMetricid = make([]int, 0)
	s.ArrInterfaceid = make([]int, 0)
	s.ArrNetworkiorate = make([]float64, 0)
	s.ArrNetworkioerrors = make([]float64, 0)
	s.ArrNetworkreceiverate = make([]float64, 0)
	s.ArrNetworkreceiveerrors = make([]float64, 0)
	s.ArrNetworktransmitrate = make([]float64, 0)
	s.ArrNetworktransmiterrors = make([]float64, 0)
}

type PodNetPerfStat struct {
	ArrPoduid                []string
	ArrOntunetime            []int64
	ArrNsuid                 []string
	ArrNodeuid               []string
	ArrMetricid              []int
	ArrInterfaceid           []int
	ArrNetworkiorate         []float64
	ArrNetworkioerrors       []float64
	ArrNetworkreceiverate    []float64
	ArrNetworkreceiveerrors  []float64
	ArrNetworktransmitrate   []float64
	ArrNetworktransmiterrors []float64
}

func (s *PodNetPerfStat) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrPoduid))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Array(s.ArrInterfaceid))
	data = append(data, pq.Float64Array(s.ArrNetworkiorate))
	data = append(data, pq.Float64Array(s.ArrNetworkioerrors))
	data = append(data, pq.Float64Array(s.ArrNetworkreceiverate))
	data = append(data, pq.Float64Array(s.ArrNetworkreceiveerrors))
	data = append(data, pq.Float64Array(s.ArrNetworktransmitrate))
	data = append(data, pq.Float64Array(s.ArrNetworktransmiterrors))

	return data
}

func (s *PodNetPerfStat) Init() {
	s.ArrPoduid = make([]string, 0)
	s.ArrOntunetime = make([]int64, 0)
	s.ArrNsuid = make([]string, 0)
	s.ArrNodeuid = make([]string, 0)
	s.ArrMetricid = make([]int, 0)
	s.ArrInterfaceid = make([]int, 0)
	s.ArrNetworkiorate = make([]float64, 0)
	s.ArrNetworkioerrors = make([]float64, 0)
	s.ArrNetworkreceiverate = make([]float64, 0)
	s.ArrNetworkreceiveerrors = make([]float64, 0)
	s.ArrNetworktransmitrate = make([]float64, 0)
	s.ArrNetworktransmiterrors = make([]float64, 0)
}

type NodeFsPerfStat struct {
	ArrNodeuid     []string
	ArrOntunetime  []int64
	ArrMetricid    []int
	ArrDeviceid    []int
	ArrFsiorate    []float64
	ArrFsreadrate  []float64
	ArrFswriterate []float64
}

func (s *NodeFsPerfStat) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Array(s.ArrDeviceid))
	data = append(data, pq.Float64Array(s.ArrFsiorate))
	data = append(data, pq.Float64Array(s.ArrFsreadrate))
	data = append(data, pq.Float64Array(s.ArrFswriterate))

	return data
}

func (s *NodeFsPerfStat) Init() {
	s.ArrNodeuid = make([]string, 0)
	s.ArrOntunetime = make([]int64, 0)
	s.ArrMetricid = make([]int, 0)
	s.ArrDeviceid = make([]int, 0)
	s.ArrFsiorate = make([]float64, 0)
	s.ArrFsreadrate = make([]float64, 0)
	s.ArrFswriterate = make([]float64, 0)
}

type PodFsPerfStat struct {
	ArrPoduid      []string
	ArrOntunetime  []int64
	ArrNsuid       []string
	ArrNodeuid     []string
	ArrMetricid    []int
	ArrDeviceid    []int
	ArrFsiorate    []float64
	ArrFsreadrate  []float64
	ArrFswriterate []float64
}

func (s *PodFsPerfStat) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrPoduid))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Array(s.ArrDeviceid))
	data = append(data, pq.Float64Array(s.ArrFsiorate))
	data = append(data, pq.Float64Array(s.ArrFsreadrate))
	data = append(data, pq.Float64Array(s.ArrFswriterate))

	return data
}

func (s *PodFsPerfStat) Init() {
	s.ArrPoduid = make([]string, 0)
	s.ArrOntunetime = make([]int64, 0)
	s.ArrNsuid = make([]string, 0)
	s.ArrNodeuid = make([]string, 0)
	s.ArrMetricid = make([]int, 0)
	s.ArrDeviceid = make([]int, 0)
	s.ArrFsiorate = make([]float64, 0)
	s.ArrFsreadrate = make([]float64, 0)
	s.ArrFswriterate = make([]float64, 0)
}

type ContainerFsPerfStat struct {
	ArrContainername []string
	ArrOntunetime    []int64
	ArrPoduid        []string
	ArrNsuid         []string
	ArrNodeuid       []string
	ArrMetricid      []int
	ArrDeviceid      []int
	ArrFsiorate      []float64
	ArrFsreadrate    []float64
	ArrFswriterate   []float64
}

func (s *ContainerFsPerfStat) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrContainername))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.StringArray(s.ArrPoduid))
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.StringArray(s.ArrNodeuid))
	data = append(data, pq.Array(s.ArrMetricid))
	data = append(data, pq.Array(s.ArrDeviceid))
	data = append(data, pq.Float64Array(s.ArrFsiorate))
	data = append(data, pq.Float64Array(s.ArrFsreadrate))
	data = append(data, pq.Float64Array(s.ArrFswriterate))

	return data
}

func (s *ContainerFsPerfStat) Init() {
	s.ArrContainername = make([]string, 0)
	s.ArrPoduid = make([]string, 0)
	s.ArrNodeuid = make([]string, 0)
	s.ArrNsuid = make([]string, 0)
	s.ArrOntunetime = make([]int64, 0)
	s.ArrMetricid = make([]int, 0)
	s.ArrDeviceid = make([]int, 0)
	s.ArrFsiorate = make([]float64, 0)
	s.ArrFsreadrate = make([]float64, 0)
	s.ArrFswriterate = make([]float64, 0)
}

type ClusterPerfStat struct {
	ArrClusterid             []int
	ArrOntunetime            []int64
	ArrPodcount              []int
	ArrCpuusage              []float64
	ArrCpusystem             []float64
	ArrCpuuser               []float64
	ArrCpuusagecores         []float64
	ArrCputotalcores         []float64
	ArrMemoryusage           []float64
	ArrMemoryusagebytes      []float64
	ArrMemorysizebytes       []float64
	ArrMemoryswap            []float64
	ArrNetworkiorate         []float64
	ArrNetworkioerrors       []float64
	ArrNetworkreceiverate    []float64
	ArrNetworkreceiveerrors  []float64
	ArrNetworktransmitrate   []float64
	ArrNetworktransmiterrors []float64
	ArrFsiorate              []float64
	ArrFsreadrate            []float64
	ArrFswriterate           []float64
	ArrProcesses             []float64
}

func (s *ClusterPerfStat) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(s.ArrClusterid))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.Array(s.ArrPodcount))
	data = append(data, pq.Float64Array(s.ArrCpuusage))
	data = append(data, pq.Float64Array(s.ArrCpusystem))
	data = append(data, pq.Float64Array(s.ArrCpuuser))
	data = append(data, pq.Float64Array(s.ArrCpuusagecores))
	data = append(data, pq.Float64Array(s.ArrCputotalcores))
	data = append(data, pq.Float64Array(s.ArrMemoryusage))
	data = append(data, pq.Float64Array(s.ArrMemoryusagebytes))
	data = append(data, pq.Float64Array(s.ArrMemorysizebytes))
	data = append(data, pq.Float64Array(s.ArrMemoryswap))
	data = append(data, pq.Float64Array(s.ArrNetworkiorate))
	data = append(data, pq.Float64Array(s.ArrNetworkioerrors))
	data = append(data, pq.Float64Array(s.ArrNetworkreceiverate))
	data = append(data, pq.Float64Array(s.ArrNetworkreceiveerrors))
	data = append(data, pq.Float64Array(s.ArrNetworktransmitrate))
	data = append(data, pq.Float64Array(s.ArrNetworktransmiterrors))
	data = append(data, pq.Float64Array(s.ArrFsiorate))
	data = append(data, pq.Float64Array(s.ArrFsreadrate))
	data = append(data, pq.Float64Array(s.ArrFswriterate))
	data = append(data, pq.Float64Array(s.ArrProcesses))

	return data
}

func (s *ClusterPerfStat) Init() {
	s.ArrClusterid = make([]int, 0)
	s.ArrOntunetime = make([]int64, 0)
	s.ArrPodcount = make([]int, 0)
	s.ArrCpuusage = make([]float64, 0)
	s.ArrCpusystem = make([]float64, 0)
	s.ArrCpuuser = make([]float64, 0)
	s.ArrCpuusagecores = make([]float64, 0)
	s.ArrCputotalcores = make([]float64, 0)
	s.ArrMemoryusage = make([]float64, 0)
	s.ArrMemoryusagebytes = make([]float64, 0)
	s.ArrMemorysizebytes = make([]float64, 0)
	s.ArrMemoryswap = make([]float64, 0)
	s.ArrNetworkiorate = make([]float64, 0)
	s.ArrNetworkioerrors = make([]float64, 0)
	s.ArrNetworkreceiverate = make([]float64, 0)
	s.ArrNetworkreceiveerrors = make([]float64, 0)
	s.ArrNetworktransmitrate = make([]float64, 0)
	s.ArrNetworktransmiterrors = make([]float64, 0)
	s.ArrFsiorate = make([]float64, 0)
	s.ArrFsreadrate = make([]float64, 0)
	s.ArrFswriterate = make([]float64, 0)
	s.ArrProcesses = make([]float64, 0)
}

type NamespacePerfStat struct {
	ArrNsuid                 []string
	ArrOntunetime            []int64
	ArrPodcount              []int
	ArrCpuusage              []float64
	ArrCpusystem             []float64
	ArrCpuuser               []float64
	ArrCpuusagecores         []float64
	ArrCputotalcores         []float64
	ArrMemoryusage           []float64
	ArrMemoryusagebytes      []float64
	ArrMemorysizebytes       []float64
	ArrMemoryswap            []float64
	ArrNetworkiorate         []float64
	ArrNetworkioerrors       []float64
	ArrNetworkreceiverate    []float64
	ArrNetworkreceiveerrors  []float64
	ArrNetworktransmitrate   []float64
	ArrNetworktransmiterrors []float64
	ArrFsiorate              []float64
	ArrFsreadrate            []float64
	ArrFswriterate           []float64
	ArrProcesses             []float64
}

func (s *NamespacePerfStat) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.StringArray(s.ArrNsuid))
	data = append(data, pq.Int64Array(s.ArrOntunetime))
	data = append(data, pq.Array(s.ArrPodcount))
	data = append(data, pq.Float64Array(s.ArrCpuusage))
	data = append(data, pq.Float64Array(s.ArrCpusystem))
	data = append(data, pq.Float64Array(s.ArrCpuuser))
	data = append(data, pq.Float64Array(s.ArrCpuusagecores))
	data = append(data, pq.Float64Array(s.ArrCputotalcores))
	data = append(data, pq.Float64Array(s.ArrMemoryusage))
	data = append(data, pq.Float64Array(s.ArrMemoryusagebytes))
	data = append(data, pq.Float64Array(s.ArrMemorysizebytes))
	data = append(data, pq.Float64Array(s.ArrMemoryswap))
	data = append(data, pq.Float64Array(s.ArrNetworkiorate))
	data = append(data, pq.Float64Array(s.ArrNetworkioerrors))
	data = append(data, pq.Float64Array(s.ArrNetworkreceiverate))
	data = append(data, pq.Float64Array(s.ArrNetworkreceiveerrors))
	data = append(data, pq.Float64Array(s.ArrNetworktransmitrate))
	data = append(data, pq.Float64Array(s.ArrNetworktransmiterrors))
	data = append(data, pq.Float64Array(s.ArrFsiorate))
	data = append(data, pq.Float64Array(s.ArrFsreadrate))
	data = append(data, pq.Float64Array(s.ArrFswriterate))
	data = append(data, pq.Float64Array(s.ArrProcesses))

	return data
}

func (s *NamespacePerfStat) Init() {
	s.ArrNsuid = make([]string, 0)
	s.ArrOntunetime = make([]int64, 0)
	s.ArrPodcount = make([]int, 0)
	s.ArrCpuusage = make([]float64, 0)
	s.ArrCpusystem = make([]float64, 0)
	s.ArrCpuuser = make([]float64, 0)
	s.ArrCpuusagecores = make([]float64, 0)
	s.ArrCputotalcores = make([]float64, 0)
	s.ArrMemoryusage = make([]float64, 0)
	s.ArrMemoryusagebytes = make([]float64, 0)
	s.ArrMemorysizebytes = make([]float64, 0)
	s.ArrMemoryswap = make([]float64, 0)
	s.ArrNetworkiorate = make([]float64, 0)
	s.ArrNetworkioerrors = make([]float64, 0)
	s.ArrNetworkreceiverate = make([]float64, 0)
	s.ArrNetworkreceiveerrors = make([]float64, 0)
	s.ArrNetworktransmitrate = make([]float64, 0)
	s.ArrNetworktransmiterrors = make([]float64, 0)
	s.ArrFsiorate = make([]float64, 0)
	s.ArrFsreadrate = make([]float64, 0)
	s.ArrFswriterate = make([]float64, 0)
	s.ArrProcesses = make([]float64, 0)
}

type WorkloadPerfStat struct {
	ArrWorkloaduid           []string
	ArrOntunetime            []int64
	ArrPodcount              []int
	ArrCpuusage              []float64
	ArrCpusystem             []float64
	ArrCpuuser               []float64
	ArrCpuusagecores         []float64
	ArrCputotalcores         []float64
	ArrMemoryusage           []float64
	ArrMemoryusagebytes      []float64
	ArrMemorysizebytes       []float64
	ArrMemoryswap            []float64
	ArrNetworkiorate         []float64
	ArrNetworkioerrors       []float64
	ArrNetworkreceiverate    []float64
	ArrNetworkreceiveerrors  []float64
	ArrNetworktransmitrate   []float64
	ArrNetworktransmiterrors []float64
	ArrFsiorate              []float64
	ArrFsreadrate            []float64
	ArrFswriterate           []float64
	ArrProcesses             []float64
}

type ServicePerfStat struct {
	ArrSvcuid                []string
	ArrOntunetime            []int64
	ArrPodcount              []int
	ArrCpuusage              []float64
	ArrCpusystem             []float64
	ArrCpuuser               []float64
	ArrCpuusagecores         []float64
	ArrCputotalcores         []float64
	ArrMemoryusage           []float64
	ArrMemoryusagebytes      []float64
	ArrMemorysizebytes       []float64
	ArrMemoryswap            []float64
	ArrNetworkiorate         []float64
	ArrNetworkioerrors       []float64
	ArrNetworkreceiverate    []float64
	ArrNetworkreceiveerrors  []float64
	ArrNetworktransmitrate   []float64
	ArrNetworktransmiterrors []float64
	ArrFsiorate              []float64
	ArrFsreadrate            []float64
	ArrFswriterate           []float64
	ArrProcesses             []float64
}

type IngressPerfStat struct {
	ArrInguid                []string
	ArrOntunetime            []int64
	ArrPodcount              []int
	ArrCpuusage              []float64
	ArrCpusystem             []float64
	ArrCpuuser               []float64
	ArrCpuusagecores         []float64
	ArrCputotalcores         []float64
	ArrMemoryusage           []float64
	ArrMemoryusagebytes      []float64
	ArrMemorysizebytes       []float64
	ArrMemoryswap            []float64
	ArrNetworkiorate         []float64
	ArrNetworkioerrors       []float64
	ArrNetworkreceiverate    []float64
	ArrNetworkreceiveerrors  []float64
	ArrNetworktransmitrate   []float64
	ArrNetworktransmiterrors []float64
	ArrFsiorate              []float64
	ArrFsreadrate            []float64
	ArrFswriterate           []float64
	ArrProcesses             []float64
}

func (w *WorkloadPerfStat) AppendValue(wluid string, ontunetime int64, wlmap *sync.Map) {
	w.ArrOntunetime = append(w.ArrOntunetime, ontunetime)
	w.ArrWorkloaduid = append(w.ArrWorkloaduid, wluid)

	cnt, _ := wlmap.Load("count")
	w.ArrPodcount = append(w.ArrPodcount, cnt.(int))

	calc, _ := wlmap.Load("calculation")
	calculation_map := calc.(*sync.Map)

	cpuusage, _ := calculation_map.Load("cpuusage")
	w.ArrCpuusage = append(w.ArrCpuusage, cpuusage.(float64))

	cpusystem, _ := calculation_map.Load("cpusystem")
	w.ArrCpusystem = append(w.ArrCpusystem, cpusystem.(float64))

	cpuuser, _ := calculation_map.Load("cpuuser")
	w.ArrCpuuser = append(w.ArrCpuuser, cpuuser.(float64))

	cpuusagecores, _ := calculation_map.Load("cpuusagecores")
	w.ArrCpuusagecores = append(w.ArrCpuusagecores, cpuusagecores.(float64))

	cputotalcores, _ := calculation_map.Load("cputotalcores")
	w.ArrCputotalcores = append(w.ArrCputotalcores, cputotalcores.(float64))

	memoryusage, _ := calculation_map.Load("memoryusage")
	w.ArrMemoryusage = append(w.ArrMemoryusage, memoryusage.(float64))

	memoryusagebytes, _ := calculation_map.Load("memoryusagebytes")
	w.ArrMemoryusagebytes = append(w.ArrMemoryusagebytes, memoryusagebytes.(float64))

	memorysizebytes, _ := calculation_map.Load("memorysizebytes")
	w.ArrMemorysizebytes = append(w.ArrMemorysizebytes, memorysizebytes.(float64))

	memoryswap, _ := calculation_map.Load("memoryswap")
	w.ArrMemoryswap = append(w.ArrMemoryswap, memoryswap.(float64))

	netrcvrate, _ := calculation_map.Load("netrcvrate")
	w.ArrNetworkreceiverate = append(w.ArrNetworkreceiverate, netrcvrate.(float64))

	netrcverrors, _ := calculation_map.Load("netrcverrors")
	w.ArrNetworkreceiveerrors = append(w.ArrNetworkreceiveerrors, netrcverrors.(float64))

	nettransrate, _ := calculation_map.Load("nettransrate")
	w.ArrNetworktransmitrate = append(w.ArrNetworktransmitrate, nettransrate.(float64))

	nettranserrors, _ := calculation_map.Load("nettranserrors")
	w.ArrNetworktransmiterrors = append(w.ArrNetworktransmiterrors, nettranserrors.(float64))

	netiorate, _ := calculation_map.Load("netiorate")
	w.ArrNetworkiorate = append(w.ArrNetworkiorate, netiorate.(float64))

	netioerrors, _ := calculation_map.Load("netioerrors")
	w.ArrNetworkioerrors = append(w.ArrNetworkioerrors, netioerrors.(float64))

	fsreadrate, _ := calculation_map.Load("fsreadrate")
	w.ArrFsreadrate = append(w.ArrFsreadrate, fsreadrate.(float64))

	fswriterate, _ := calculation_map.Load("fswriterate")
	w.ArrFswriterate = append(w.ArrFswriterate, fswriterate.(float64))

	fsiorate, _ := calculation_map.Load("fsiorate")
	w.ArrFsiorate = append(w.ArrFsiorate, fsiorate.(float64))

	processcount, _ := calculation_map.Load("processcount")
	w.ArrProcesses = append(w.ArrProcesses, processcount.(float64))
}

func (w *WorkloadPerfStat) GetArgs() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, pq.Array(w.ArrWorkloaduid))
	data = append(data, pq.Int64Array(w.ArrOntunetime))
	data = append(data, pq.Array(w.ArrPodcount))
	data = append(data, pq.Float64Array(w.ArrCpuusage))
	data = append(data, pq.Float64Array(w.ArrCpusystem))
	data = append(data, pq.Float64Array(w.ArrCpuuser))
	data = append(data, pq.Float64Array(w.ArrCpuusagecores))
	data = append(data, pq.Float64Array(w.ArrCputotalcores))
	data = append(data, pq.Float64Array(w.ArrMemoryusage))
	data = append(data, pq.Float64Array(w.ArrMemoryusagebytes))
	data = append(data, pq.Float64Array(w.ArrMemorysizebytes))
	data = append(data, pq.Float64Array(w.ArrMemoryswap))
	data = append(data, pq.Float64Array(w.ArrNetworkiorate))
	data = append(data, pq.Float64Array(w.ArrNetworkioerrors))
	data = append(data, pq.Float64Array(w.ArrNetworkreceiverate))
	data = append(data, pq.Float64Array(w.ArrNetworkreceiveerrors))
	data = append(data, pq.Float64Array(w.ArrNetworktransmitrate))
	data = append(data, pq.Float64Array(w.ArrNetworktransmiterrors))
	data = append(data, pq.Float64Array(w.ArrFsiorate))
	data = append(data, pq.Float64Array(w.ArrFsreadrate))
	data = append(data, pq.Float64Array(w.ArrFswriterate))
	data = append(data, pq.Float64Array(w.ArrProcesses))

	return data
}

func (w *WorkloadPerfStat) Init() {
	w.ArrWorkloaduid = make([]string, 0)
	w.ArrOntunetime = make([]int64, 0)
	w.ArrPodcount = make([]int, 0)
	w.ArrCpuusage = make([]float64, 0)
	w.ArrCpusystem = make([]float64, 0)
	w.ArrCpuuser = make([]float64, 0)
	w.ArrCpuusagecores = make([]float64, 0)
	w.ArrCputotalcores = make([]float64, 0)
	w.ArrMemoryusage = make([]float64, 0)
	w.ArrMemoryusagebytes = make([]float64, 0)
	w.ArrMemorysizebytes = make([]float64, 0)
	w.ArrMemoryswap = make([]float64, 0)
	w.ArrNetworkiorate = make([]float64, 0)
	w.ArrNetworkioerrors = make([]float64, 0)
	w.ArrNetworkreceiverate = make([]float64, 0)
	w.ArrNetworkreceiveerrors = make([]float64, 0)
	w.ArrNetworktransmitrate = make([]float64, 0)
	w.ArrNetworktransmiterrors = make([]float64, 0)
	w.ArrFsiorate = make([]float64, 0)
	w.ArrFsreadrate = make([]float64, 0)
	w.ArrFswriterate = make([]float64, 0)
	w.ArrProcesses = make([]float64, 0)
}

type SummaryIndex struct {
	DataType      string
	NsUID         string
	NodeUID       string
	PodUID        string
	ContainerName string
	MetricId      int
	MetricName    string
}
