package kubeapi

import (
	"encoding/json"
	"net"
	"onTuneKubeManager/common"

	dto "github.com/prometheus/client_model/go"
)

type MappingMetricContainer struct {
	ContainerName         string  `json:"container"`
	PodName               string  `json:"pod"`
	Namespace             string  `json:"namespace"`
	NodeName              string  `json:"node"`
	MetricID              string  `json:"metricid"`
	TimestampMs           int64   `json:"timestampms"`
	CpuUsageSecondsTotal  float64 `json:"cpuusagesecondstotal"`
	CpuSystemSecondsTotal float64 `json:"cpusystemsecondstotal"`
	CpuUserSecondsTotal   float64 `json:"cpuusersecondstotal"`
	MemoryUsageBytes      float64 `json:"memoryusagebytes"`
	MemoryWorkingSetBytes float64 `json:"memoryworkingsetbytes"`
	MemoryCache           float64 `json:"memorycache"`
	MemorySwap            float64 `json:"memoryswap"`
	MemoryRss             float64 `json:"memoryrss"`
	Processes             float64 `json:"processes"`
	MetricImage           string  `json:"image"`
}

func (mmc *MappingMetricContainer) setMetricValue(metric *dto.Metric, mname string) {
	switch mname {
	case METRIC_CONTAINER_CPU_USAGE_SECONDS_TOTAL:
		mmc.CpuUsageSecondsTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_CPU_SYSTEM_SECONDS_TOTAL:
		mmc.CpuSystemSecondsTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_CPU_USER_SECONDS_TOTAL:
		mmc.CpuUserSecondsTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_MEMORY_USAGE_BYTES:
		mmc.MemoryUsageBytes = metric.Gauge.GetValue()
	case METRIC_CONTAINER_MEMORY_WORKING_SET_BYTES:
		mmc.MemoryWorkingSetBytes = metric.Gauge.GetValue()
	case METRIC_CONTAINER_MEMORY_CACHE:
		mmc.MemoryCache = metric.Gauge.GetValue()
	case METRIC_CONTAINER_MEMORY_SWAP:
		mmc.MemorySwap = metric.Gauge.GetValue()
	case METRIC_CONTAINER_MEMORY_RSS:
		mmc.MemoryRss = metric.Gauge.GetValue()
	case METRIC_CONTAINER_PROCESSES:
		mmc.Processes = metric.Gauge.GetValue()
	}

	if mmc.TimestampMs == 0 {
		mmc.TimestampMs = metric.GetTimestampMs()
	}
}

type MappingMetricPod struct {
	PodName               string  `json:"pod"`
	Namespace             string  `json:"namespace"`
	NodeName              string  `json:"node"`
	MetricID              string  `json:"metricid"`
	TimestampMs           int64   `json:"timestampms"`
	CpuUsageSecondsTotal  float64 `json:"cpuusagesecondstotal"`
	CpuSystemSecondsTotal float64 `json:"cpusystemsecondstotal"`
	CpuUserSecondsTotal   float64 `json:"cpuusersecondstotal"`
	MemoryUsageBytes      float64 `json:"memoryusagebytes"`
	MemoryWorkingSetBytes float64 `json:"memoryworkingsetbytes"`
	MemoryCache           float64 `json:"memorycache"`
	MemorySwap            float64 `json:"memoryswap"`
	MemoryRss             float64 `json:"memoryrss"`
	Processes             float64 `json:"processes"`
	MetricImage           string  `json:"image"`
}

func (mmp *MappingMetricPod) setMetricValue(metric *dto.Metric, mname string) {
	switch mname {
	case METRIC_CONTAINER_CPU_USAGE_SECONDS_TOTAL:
		mmp.CpuUsageSecondsTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_CPU_SYSTEM_SECONDS_TOTAL:
		mmp.CpuSystemSecondsTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_CPU_USER_SECONDS_TOTAL:
		mmp.CpuUserSecondsTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_MEMORY_USAGE_BYTES:
		mmp.MemoryUsageBytes = metric.Gauge.GetValue()
	case METRIC_CONTAINER_MEMORY_WORKING_SET_BYTES:
		mmp.MemoryWorkingSetBytes = metric.Gauge.GetValue()
	case METRIC_CONTAINER_MEMORY_CACHE:
		mmp.MemoryCache = metric.Gauge.GetValue()
	case METRIC_CONTAINER_MEMORY_SWAP:
		mmp.MemorySwap = metric.Gauge.GetValue()
	case METRIC_CONTAINER_MEMORY_RSS:
		mmp.MemoryRss = metric.Gauge.GetValue()
	case METRIC_CONTAINER_PROCESSES:
		mmp.Processes = metric.Gauge.GetValue()
	}

	if mmp.TimestampMs == 0 {
		mmp.TimestampMs = metric.GetTimestampMs()
	}
}

type MappingMetricNode struct {
	NodeName              string  `json:"node"`
	MetricID              string  `json:"metricid"`
	TimestampMs           int64   `json:"timestampms"`
	CpuUsageSecondsTotal  float64 `json:"cpuusagesecondstotal"`
	CpuSystemSecondsTotal float64 `json:"cpusystemsecondstotal"`
	CpuUserSecondsTotal   float64 `json:"cpuusersecondstotal"`
	MemoryUsageBytes      float64 `json:"memoryusagebytes"`
	MemoryWorkingSetBytes float64 `json:"memoryworkingsetbytes"`
	MemoryCache           float64 `json:"memorycache"`
	MemorySwap            float64 `json:"memoryswap"`
	MemoryRss             float64 `json:"memoryrss"`
	FsReadsBytesTotal     float64 `json:"fsreadsbytestotal"`
	FsWritesBytesTotal    float64 `json:"fswritesbytestotal"`
	Processes             float64 `json:"processes"`
	MetricImage           string  `json:"image"`
}

func (mmn *MappingMetricNode) setMetricValue(metric *dto.Metric, mname string) {
	switch mname {
	case METRIC_CONTAINER_CPU_USAGE_SECONDS_TOTAL:
		mmn.CpuUsageSecondsTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_CPU_SYSTEM_SECONDS_TOTAL:
		mmn.CpuSystemSecondsTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_CPU_USER_SECONDS_TOTAL:
		mmn.CpuUserSecondsTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_MEMORY_USAGE_BYTES:
		mmn.MemoryUsageBytes = metric.Gauge.GetValue()
	case METRIC_CONTAINER_MEMORY_WORKING_SET_BYTES:
		mmn.MemoryWorkingSetBytes = metric.Gauge.GetValue()
	case METRIC_CONTAINER_MEMORY_CACHE:
		mmn.MemoryCache = metric.Gauge.GetValue()
	case METRIC_CONTAINER_MEMORY_SWAP:
		mmn.MemorySwap = metric.Gauge.GetValue()
	case METRIC_CONTAINER_MEMORY_RSS:
		mmn.MemoryRss = metric.Gauge.GetValue()
	case METRIC_CONTAINER_FS_READS_BYTES_TOTAL:
		mmn.FsReadsBytesTotal = metric.Gauge.GetValue()
	case METRIC_CONTAINER_FS_WRITES_BYTES_TOTAL:
		mmn.FsWritesBytesTotal = metric.Gauge.GetValue()
	case METRIC_CONTAINER_PROCESSES:
		mmn.Processes = metric.Gauge.GetValue()
	}

	if mmn.TimestampMs == 0 {
		mmn.TimestampMs = metric.GetTimestampMs()
	}
}

type MappingMetricNodeNetwork struct {
	NodeName                   string  `json:"node"`
	MetricID                   string  `json:"metricid"`
	Interface                  string  `json:"interface"`
	TimestampMs                int64   `json:"timestampms"`
	NetworkReceiveBytesTotal   float64 `json:"networkreceivebytestotal"`
	NetworkReceiveErrorsTotal  float64 `json:"networkreceiveerrorstotal"`
	NetworkTransmitBytesTotal  float64 `json:"networktransmitbytestotal"`
	NetworkTransmitErrorsTotal float64 `json:"networktransmiterrorstotal"`
	MetricImage                string  `json:"image"`
}

func (mmn *MappingMetricNodeNetwork) setMetricValue(metric *dto.Metric, mname string) {
	switch mname {
	case METRIC_CONTAINER_NETWORK_RECEIVE_BYTES_TOTAL:
		mmn.NetworkReceiveBytesTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_NETWORK_RECEIVE_ERRORS_TOTAL:
		mmn.NetworkReceiveErrorsTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_NETWORK_TRANSMIT_BYTES_TOTAL:
		mmn.NetworkTransmitBytesTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_NETWORK_TRANSMIT_ERRORS_TOTAL:
		mmn.NetworkTransmitErrorsTotal = metric.Gauge.GetValue()
	}

	if mmn.TimestampMs == 0 {
		mmn.TimestampMs = metric.GetTimestampMs()
	}
}

type MappingMetricPodNetwork struct {
	PodName                    string  `json:"pod"`
	Namespace                  string  `json:"namespace"`
	NodeName                   string  `json:"node"`
	MetricID                   string  `json:"metricid"`
	Interface                  string  `json:"interface"`
	TimestampMs                int64   `json:"timestampms"`
	NetworkReceiveBytesTotal   float64 `json:"networkreceivebytestotal"`
	NetworkReceiveErrorsTotal  float64 `json:"networkreceiveerrorstotal"`
	NetworkTransmitBytesTotal  float64 `json:"networktransmitbytestotal"`
	NetworkTransmitErrorsTotal float64 `json:"networktransmiterrorstotal"`
	MetricImage                string  `json:"image"`
}

func (mmn *MappingMetricPodNetwork) setMetricValue(metric *dto.Metric, mname string) {
	switch mname {
	case METRIC_CONTAINER_NETWORK_RECEIVE_BYTES_TOTAL:
		mmn.NetworkReceiveBytesTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_NETWORK_RECEIVE_ERRORS_TOTAL:
		mmn.NetworkReceiveErrorsTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_NETWORK_TRANSMIT_BYTES_TOTAL:
		mmn.NetworkTransmitBytesTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_NETWORK_TRANSMIT_ERRORS_TOTAL:
		mmn.NetworkTransmitErrorsTotal = metric.Gauge.GetValue()
	}

	if mmn.TimestampMs == 0 {
		mmn.TimestampMs = metric.GetTimestampMs()
	}
}

type MappingMetricNodeFilesystem struct {
	NodeName           string  `json:"node"`
	MetricID           string  `json:"metricid"`
	Device             string  `json:"device"`
	TimestampMs        int64   `json:"timestampms"`
	FsInodesFree       float64 `json:"fsinodesfree"`
	FsInodesTotal      float64 `json:"fsinodestotal"`
	FsLimitBytes       float64 `json:"fslimitbytes"`
	FsReadsBytesTotal  float64 `json:"fsreadsbytestotal"`
	FsWritesBytesTotal float64 `json:"fswritesbytestotal"`
	FsUsageBytes       float64 `json:"fsusagebytes"`
	MetricImage        string  `json:"image"`
}

func (mmf *MappingMetricNodeFilesystem) setMetricValue(metric *dto.Metric, mname string) {
	switch mname {
	case METRIC_CONTAINER_FS_INODES_FREE:
		mmf.FsInodesFree = metric.Gauge.GetValue()
	case METRIC_CONTAINER_FS_INODES_TOTAL:
		mmf.FsInodesTotal = metric.Gauge.GetValue()
	case METRIC_CONTAINER_FS_LIMIT_BYTES:
		mmf.FsLimitBytes = metric.Gauge.GetValue()
	case METRIC_CONTAINER_FS_READS_BYTES_TOTAL:
		mmf.FsReadsBytesTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_FS_WRITES_BYTES_TOTAL:
		mmf.FsWritesBytesTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_FS_USAGE_BYTES:
		mmf.FsUsageBytes = metric.Gauge.GetValue()
	}

	if mmf.TimestampMs == 0 {
		mmf.TimestampMs = metric.GetTimestampMs()
	}
}

type MappingMetricPodFilesystem struct {
	PodName            string  `json:"pod"`
	Namespace          string  `json:"namespace"`
	NodeName           string  `json:"node"`
	MetricID           string  `json:"metricid"`
	Device             string  `json:"device"`
	TimestampMs        int64   `json:"timestampms"`
	FsLimitBytes       float64 `json:"fslimitbytes"`
	FsReadsBytesTotal  float64 `json:"fsreadsbytestotal"`
	FsWritesBytesTotal float64 `json:"fswritesbytestotal"`
	FsUsageBytes       float64 `json:"fsusagebytes"`
	MetricImage        string  `json:"image"`
}

func (mmf *MappingMetricPodFilesystem) setMetricValue(metric *dto.Metric, mname string) {
	switch mname {
	case METRIC_CONTAINER_FS_READS_BYTES_TOTAL:
		mmf.FsReadsBytesTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_FS_WRITES_BYTES_TOTAL:
		mmf.FsWritesBytesTotal = metric.Counter.GetValue()
	}

	if mmf.TimestampMs == 0 {
		mmf.TimestampMs = metric.GetTimestampMs()
	}
}

type MappingMetricContainerFilesystem struct {
	ContainerName      string  `json:"container"`
	PodName            string  `json:"pod"`
	Namespace          string  `json:"namespace"`
	NodeName           string  `json:"node"`
	MetricID           string  `json:"metricid"`
	Device             string  `json:"device"`
	TimestampMs        int64   `json:"timestampms"`
	FsLimitBytes       float64 `json:"fslimitbytes"`
	FsReadsBytesTotal  float64 `json:"fsreadsbytestotal"`
	FsWritesBytesTotal float64 `json:"fswritesbytestotal"`
	FsUsageBytes       float64 `json:"fsusagebytes"`
	MetricImage        string  `json:"image"`
}

func (mmf *MappingMetricContainerFilesystem) setMetricValue(metric *dto.Metric, mname string) {
	switch mname {
	case METRIC_CONTAINER_FS_READS_BYTES_TOTAL:
		mmf.FsReadsBytesTotal = metric.Counter.GetValue()
	case METRIC_CONTAINER_FS_WRITES_BYTES_TOTAL:
		mmf.FsWritesBytesTotal = metric.Counter.GetValue()
	}

	if mmf.TimestampMs == 0 {
		mmf.TimestampMs = metric.GetTimestampMs()
	}
}

type MappingMetricSource struct {
	ClientInfo              `json:"client"`
	Host                    string                             `json:"host"`
	Node                    string                             `json:"node"`
	NodeInfo                []MappingMetricNode                `json:"nodes"`
	PodInfo                 []MappingMetricPod                 `json:"pods"`
	ContainerInfo           []MappingMetricContainer           `json:"containers"`
	NodeNetworkInfo         []MappingMetricNodeNetwork         `json:"nodenetworks"`
	PodNetworkInfo          []MappingMetricPodNetwork          `json:"podnetworks"`
	NodeFilesystemInfo      []MappingMetricNodeFilesystem      `json:"nodefilesystems"`
	PodFilesystemInfo       []MappingMetricPodFilesystem       `json:"podfilesystems"`
	ContainerFilesystemInfo []MappingMetricContainerFilesystem `json:"containerfilesystems"`
}

func (m *MappingMetricSource) MakeMetricData(kmetric *KubernetesAPIMetricSource) {
	for k, v := range kmetric.NodeMetric {
		m.NodeInfo = make([]MappingMetricNode, 0)
		m.PodInfo = make([]MappingMetricPod, 0)
		m.ContainerInfo = make([]MappingMetricContainer, 0)
		m.NodeNetworkInfo = make([]MappingMetricNodeNetwork, 0)
		m.PodNetworkInfo = make([]MappingMetricPodNetwork, 0)
		m.NodeFilesystemInfo = make([]MappingMetricNodeFilesystem, 0)
		m.PodFilesystemInfo = make([]MappingMetricPodFilesystem, 0)
		m.ContainerFilesystemInfo = make([]MappingMetricContainerFilesystem, 0)
		m.Node = k

		host, _, _ := net.SplitHostPort(kmetric.clientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)
		m.Host = host

		m.MakeNodeMetricData(k, &v)
		metric_json, _ := json.Marshal(m) // Node 별 insert
		common.ChannelMetricData <- metric_json
	}
}

func (m *MappingMetricSource) MakeNodeMetricData(nodename string, metricmap *map[string]*dto.MetricFamily) {
	nodemap := make(map[LabelFilter]*MappingMetricNode)
	podmap := make(map[LabelFilter]*MappingMetricPod)
	containermap := make(map[LabelFilter]*MappingMetricContainer)
	nodenetworkmap := make(map[LabelFilter]*MappingMetricNodeNetwork)
	podnetworkmap := make(map[LabelFilter]*MappingMetricPodNetwork)
	nodefilesystemmap := make(map[LabelFilter]*MappingMetricNodeFilesystem)
	podfilesystemmap := make(map[LabelFilter]*MappingMetricPodFilesystem)
	containerfilesystemmap := make(map[LabelFilter]*MappingMetricContainerFilesystem)

	for k, mf := range *metricmap {
		for _, mfm := range mf.GetMetric() {
			if lf := searchLabel(mfm, "node"); lf.flag {
				if _, ok := nodemap[lf]; !ok {
					nodemap[lf] = &MappingMetricNode{
						NodeName:    nodename,
						MetricID:    lf.id,
						MetricImage: lf.image,
					}
				}
				nodemap[lf].setMetricValue(mfm, k)
			}
			if lf := searchLabel(mfm, "pod"); lf.flag {
				if _, ok := podmap[lf]; !ok {
					podmap[lf] = &MappingMetricPod{
						PodName:     lf.pod,
						Namespace:   lf.namespace,
						NodeName:    nodename,
						MetricID:    lf.id,
						MetricImage: lf.image,
					}
				}
				podmap[lf].setMetricValue(mfm, k)
			}
			if lf := searchLabel(mfm, "container"); lf.flag {
				if _, ok := containermap[lf]; !ok {
					containermap[lf] = &MappingMetricContainer{
						ContainerName: lf.container,
						PodName:       lf.pod,
						Namespace:     lf.namespace,
						NodeName:      nodename,
						MetricID:      lf.id,
						MetricImage:   lf.image,
					}
				}
				containermap[lf].setMetricValue(mfm, k)
			}
			if lf := searchLabel(mfm, "nodenetwork"); lf.flag {
				if _, ok := nodenetworkmap[lf]; !ok {
					nodenetworkmap[lf] = &MappingMetricNodeNetwork{
						NodeName:    nodename,
						MetricID:    lf.id,
						Interface:   lf.itfc,
						MetricImage: lf.image,
					}
				}
				nodenetworkmap[lf].setMetricValue(mfm, k)
			}
			if lf := searchLabel(mfm, "podnetwork"); lf.flag {
				if _, ok := podnetworkmap[lf]; !ok {
					podnetworkmap[lf] = &MappingMetricPodNetwork{
						PodName:     lf.pod,
						Namespace:   lf.namespace,
						NodeName:    nodename,
						MetricID:    lf.id,
						Interface:   lf.itfc,
						MetricImage: lf.image,
					}
				}
				podnetworkmap[lf].setMetricValue(mfm, k)
			}
			if lf := searchLabel(mfm, "nodefilesystem"); lf.flag {
				if _, ok := nodefilesystemmap[lf]; !ok {
					nodefilesystemmap[lf] = &MappingMetricNodeFilesystem{
						NodeName:    nodename,
						MetricID:    lf.id,
						Device:      lf.device,
						MetricImage: lf.image,
					}
				}
				nodefilesystemmap[lf].setMetricValue(mfm, k)
			}
			if lf := searchLabel(mfm, "podfilesystem"); lf.flag {
				if _, ok := podfilesystemmap[lf]; !ok {
					podfilesystemmap[lf] = &MappingMetricPodFilesystem{
						PodName:     lf.pod,
						Namespace:   lf.namespace,
						NodeName:    nodename,
						MetricID:    lf.id,
						Device:      lf.device,
						MetricImage: lf.image,
					}
				}
				podfilesystemmap[lf].setMetricValue(mfm, k)
			}
			if lf := searchLabel(mfm, "containerfilesystem"); lf.flag {
				if _, ok := containerfilesystemmap[lf]; !ok {
					containerfilesystemmap[lf] = &MappingMetricContainerFilesystem{
						ContainerName: lf.container,
						PodName:       lf.pod,
						Namespace:     lf.namespace,
						NodeName:      nodename,
						MetricID:      lf.id,
						Device:        lf.device,
						MetricImage:   lf.image,
					}
				}
				containerfilesystemmap[lf].setMetricValue(mfm, k)
			}
		}
	}

	// Invalid Swap Memory Size modify 0
	for _, nm := range nodemap {
		if nm.MemorySwap > nm.MemoryUsageBytes {
			nm.MemorySwap = 0
		}
		m.NodeInfo = append(m.NodeInfo, *nm)
	}

	for _, pm := range podmap {
		if pm.MemorySwap > pm.MemoryUsageBytes {
			pm.MemorySwap = 0
		}
		m.PodInfo = append(m.PodInfo, *pm)
	}

	for _, cm := range containermap {
		if cm.MemorySwap > cm.MemoryUsageBytes {
			cm.MemorySwap = 0
		}
		m.ContainerInfo = append(m.ContainerInfo, *cm)
	}

	for _, nntm := range nodenetworkmap {
		m.NodeNetworkInfo = append(m.NodeNetworkInfo, *nntm)
	}

	for _, pntm := range podnetworkmap {
		m.PodNetworkInfo = append(m.PodNetworkInfo, *pntm)
	}

	for _, nfm := range nodefilesystemmap {
		m.NodeFilesystemInfo = append(m.NodeFilesystemInfo, *nfm)
	}

	for _, pfm := range podfilesystemmap {
		m.PodFilesystemInfo = append(m.PodFilesystemInfo, *pfm)
	}

	for _, cfm := range containerfilesystemmap {
		m.ContainerFilesystemInfo = append(m.ContainerFilesystemInfo, *cfm)
	}
}
