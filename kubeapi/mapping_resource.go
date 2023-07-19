package kubeapi

import (
	"encoding/json"
	"fmt"
	"net"
	"onTuneKubeManager/common"
	"onTuneKubeManager/kubeapi/types"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

type MappingInterface interface {
	GetHost() string
}

type MappingNamespace struct {
	Name      string    `json:"name"`
	UID       string    `json:"uid"`
	StartTime time.Time `json:"starttime"`
	Host      string    `json:"host"`
	Labels    string    `json:"labels"`
	Status    string    `json:"status"`
}

func (m *MappingNamespace) GetHost() string {
	return m.Host
}

type MappingNode struct {
	Name                    string    `json:"name"`
	UID                     string    `json:"uid"`
	StartTime               time.Time `json:"starttime"`
	Host                    string    `json:"host"`
	Labels                  string    `json:"labels"`
	NodeType                string    `json:"nodetype"`
	KernelVersion           string    `json:"kernelversion"`
	OSImage                 string    `json:"osimage"`
	OSName                  string    `json:"osname"`
	ContainerRuntimeVersion string    `json:"containerruntimever"`
	KubeletVersion          string    `json:"kubeletver"`
	KubeProxyVersion        string    `json:"kubeproxyver"`
	CPUArch                 string    `json:"cpuarch"`
	CPUCount                int64     `json:"cpucount"`
	EphemeralStorage        int64     `json:"ephemeralstorage"`
	MemorySize              int64     `json:"memorysize"`
	Pods                    int64     `json:"pods"`
	IP                      string    `json:"ip"`
	NamespaceName           string    `json:"nsname"`
	Status                  int       `json:"status"`
}

func (m *MappingNode) GetHost() string {
	return m.Host
}

type MappingContainer struct {
	Name             string `json:"name"`
	UID              string `json:"uid"`
	Host             string `json:"host"`
	Image            string `json:"image"`
	Ports            string `json:"ports"`
	Env              string `json:"env"`
	LimitCpu         int64  `json:"limitcpu"`
	LimitMemory      int64  `json:"limitmemory"`
	LimitStorage     int64  `json:"limitstorage"`
	LimitEphemeral   int64  `json:"limitephemeral"`
	RequestCpu       int64  `json:"reqcpu"`
	RequestMemory    int64  `json:"reqmemory"`
	RequestStorage   int64  `json:"reqstorage"`
	RequestEphemeral int64  `json:"reqephemeral"`
	VolumeMounts     string `json:"volumemounts"`
	State            string `json:"state"`
}

func (m *MappingContainer) GetHost() string {
	return m.Host
}

type MappingPod struct {
	Name                  string             `json:"name"`
	UID                   string             `json:"uid"`
	StartTime             time.Time          `json:"starttime"`
	Host                  string             `json:"host"`
	Labels                string             `json:"labels"`
	Selector              string             `json:"selector"`
	AnnotationUID         string             `json:"annotationuid"`
	RestartPolicy         string             `json:"restartpolicy"`
	ServiceAccount        string             `json:"serviceaccount"`
	Status                string             `json:"status"`
	HostIP                string             `json:"hostip"`
	PodIP                 string             `json:"podip"`
	RestartCount          int32              `json:"restartcount"`
	RestartTime           time.Time          `json:"restarttime"`
	Condition             string             `json:"condition"`
	StaticPod             string             `json:"staticpod"`
	NodeName              string             `json:"node"`
	NamespaceName         string             `json:"nsname"`
	ReferenceKind         string             `json:"refkind"`
	ReferenceUID          string             `json:"refuid"`
	PersistentVolumeClaim string             `json:"pvc"`
	Containers            []MappingContainer `json:"containers"`
}

func (m *MappingPod) GetHost() string {
	return m.Host
}

type MappingService struct {
	Name          string    `json:"name"`
	UID           string    `json:"uid"`
	StartTime     time.Time `json:"starttime"`
	Host          string    `json:"host"`
	Labels        string    `json:"labels"`
	Selector      string    `json:"selector"`
	ServiceType   string    `json:"servicetype"`
	ClusterIP     string    `json:"clusterip"`
	Ports         string    `json:"ports"`
	NamespaceName string    `json:"nsname"`
}

func (m *MappingService) GetHost() string {
	return m.Host
}

type MappingPvc struct {
	Name             string    `json:"name"`
	UID              string    `json:"uid"`
	StartTime        time.Time `json:"starttime"`
	Host             string    `json:"host"`
	Labels           string    `json:"labels"`
	Selector         string    `json:"selector"`
	AccessModes      string    `json:"accessmodes"`
	RequestStorage   int64     `json:"reqstorage"`
	Status           string    `json:"status"`
	StorageClassName string    `json:"scname"`
	NamespaceName    string    `json:"nsname"`
}

func (m *MappingPvc) GetHost() string {
	return m.Host
}

type MappingPv struct {
	Name          string    `json:"name"`
	UID           string    `json:"uid"`
	StartTime     time.Time `json:"starttime"`
	Host          string    `json:"host"`
	Labels        string    `json:"labels"`
	AccessModes   string    `json:"accessmodes"`
	Capacity      int64     `json:"capacity"`
	ReclaimPolicy string    `json:"reclaimpolicy"`
	Status        string    `json:"status"`
	PvcUID        string    `json:"pvcuid"`
}

func (m *MappingPv) GetHost() string {
	return m.Host
}

type MappingDeployment struct {
	Name                 string    `json:"name"`
	UID                  string    `json:"uid"`
	StartTime            time.Time `json:"starttime"`
	Host                 string    `json:"host"`
	Labels               string    `json:"labels"`
	Selector             string    `json:"selector"`
	ServiceAccount       string    `json:"serviceAccount"`
	Replicas             int32     `json:"replicas"`
	UpdatedReplicas      int32     `json:"updatedrs"`
	ReadyReplicas        int32     `json:"readyrs"`
	AvailableReplicas    int32     `json:"availablers"`
	ObservedGeneneration int64     `json:"observedgen"`
	NamespaceName        string    `json:"nsname"`
}

func (m *MappingDeployment) GetHost() string {
	return m.Host
}

type MappingStatefulSet struct {
	Name              string    `json:"name"`
	UID               string    `json:"uid"`
	StartTime         time.Time `json:"starttime"`
	Host              string    `json:"host"`
	Labels            string    `json:"labels"`
	Selector          string    `json:"selector"`
	ServiceAccount    string    `json:"serviceAccount"`
	Replicas          int32     `json:"replicas"`
	UpdatedReplicas   int32     `json:"updatedrs"`
	ReadyReplicas     int32     `json:"readyrs"`
	AvailableReplicas int32     `json:"availablers"`
	NamespaceName     string    `json:"nsname"`
}

func (m *MappingStatefulSet) GetHost() string {
	return m.Host
}

type MappingDaemonSet struct {
	Name                   string    `json:"name"`
	UID                    string    `json:"uid"`
	StartTime              time.Time `json:"starttime"`
	Host                   string    `json:"host"`
	Labels                 string    `json:"labels"`
	Selector               string    `json:"selector"`
	ServiceAccount         string    `json:"serviceAccount"`
	CurrentNumberScheduled int32     `json:"current"`
	DesiredNumberScheduled int32     `json:"desired"`
	NumberReady            int32     `json:"ready"`
	UpdatedNumberScheduled int32     `json:"updated"`
	NumberAvailable        int32     `json:"available"`
	NamespaceName          string    `json:"nsname"`
}

func (m *MappingDaemonSet) GetHost() string {
	return m.Host
}

type MappingReplicaSet struct {
	Name                 string    `json:"name"`
	UID                  string    `json:"uid"`
	StartTime            time.Time `json:"starttime"`
	Host                 string    `json:"host"`
	Labels               string    `json:"labels"`
	Selector             string    `json:"selector"`
	Replicas             int32     `json:"replicas"`
	FullyLabeledReplicas int32     `json:"fullylabeledrs"`
	ReadyReplicas        int32     `json:"readyrs"`
	AvailableReplicas    int32     `json:"availablers"`
	ObservedGeneneration int64     `json:"observedgen"`
	ReferenceKind        string    `json:"refkind"`
	ReferenceUID         string    `json:"refuid"`
	NamespaceName        string    `json:"nsname"`
}

func (m *MappingReplicaSet) GetHost() string {
	return m.Host
}

type MappingIngressHost struct {
	BackendType      string `json:"backendtype"`
	BackendName      string `json:"backendname"`
	UID              string `json:"uid"`
	Hostname         string `json:"hostname"`
	Host             string `json:"host"`
	PathType         string `json:"pathtype"`
	Path             string `json:"path"`
	ServicePort      int32  `json:"serviceport"`
	ResourceAPIGroup string `json:"resourceapigroup"`
	ResourceKind     string `json:"resourceKind"`
}

func (m *MappingIngressHost) GetHost() string {
	return m.Host
}

type MappingIngress struct {
	Name             string               `json:"name"`
	UID              string               `json:"uid"`
	StartTime        time.Time            `json:"starttime"`
	Host             string               `json:"host"`
	Labels           string               `json:"labels"`
	IngressClassName string               `json:"classname"`
	NamespaceName    string               `json:"nsname"`
	IngressHosts     []MappingIngressHost `json:"hostdata"`
}

func (m *MappingIngress) GetHost() string {
	return m.Host
}

type MappingStorageClass struct {
	Name                 string    `json:"name"`
	UID                  string    `json:"uid"`
	StartTime            time.Time `json:"starttime"`
	Host                 string    `json:"host"`
	Labels               string    `json:"labels"`
	Provisioner          string    `json:"provisioner"`
	ReclaimPolicy        string    `json:"reclaimpolicy"`
	VolumeBindingMode    string    `json:"volumebindingmode"`
	AllowVolumeExpansion bool      `json:"allowvolumeexp"`
}

func (m *MappingStorageClass) GetHost() string {
	return m.Host
}

type ContainerStatus struct {
	State        string
	RestartCount int32
	Ready        bool
	Started      bool
}

type GlobalInfo struct {
	masternodename string
}

type MappingResourceSource struct {
	ClientInfo     `json:"client"`
	GlobalInfo     `json:"globalinfo"`
	Namespaces     []MappingNamespace    `json:"namespaces"`
	Nodes          []MappingNode         `json:"nodes"`
	Pods           []MappingPod          `json:"pods"`
	Services       []MappingService      `json:"services"`
	Pvcs           []MappingPvc          `json:"persistentvolumeclaims"`
	Pvs            []MappingPv           `json:"persistentvolumes"`
	Deployments    []MappingDeployment   `json:"deployments"`
	StatefulSets   []MappingStatefulSet  `json:"statefulsets"`
	DaemonSets     []MappingDaemonSet    `json:"daemonsets"`
	ReplicaSets    []MappingReplicaSet   `json:"replicasets"`
	Ingresses      []MappingIngress      `json:"ingresses"`
	StorageClasses []MappingStorageClass `json:"storageclasses"`
}

func (m *MappingResourceSource) MakeResourceData(kresource *KubernetesAPIResourceSource) {
	if kresource.Core != nil && m.ClientInfo.GetVersions().GetVersion("core").Version == "v1" {
		m.MakeNamespaceV1Data(kresource.Core.(*types.CoreV1Data))
		m.MakeNodeV1Data(kresource.Core.(*types.CoreV1Data))
		m.MakePersistentVolumeClaimV1Data(kresource.Core.(*types.CoreV1Data))
		m.MakePersistentVolumeV1Data(kresource.Core.(*types.CoreV1Data))
		m.MakePodV1Data(kresource.Core.(*types.CoreV1Data))
		m.MakeServiceV1Data(kresource.Core.(*types.CoreV1Data))
	}
	if kresource.Storage != nil && m.ClientInfo.GetVersions().GetVersion("storage").Version == "v1" {
		m.MakeStorageClassV1Data(kresource.Storage.(*types.StorageV1Data))
	}
	if kresource.Apps != nil && m.ClientInfo.GetVersions().GetVersion("apps").Version == "v1" {
		m.MakeDeploymentV1Data(kresource.Apps.(*types.AppsV1Data))
		m.MakeStatefulSetV1Data(kresource.Apps.(*types.AppsV1Data))
		m.MakeDaemonSetV1Data(kresource.Apps.(*types.AppsV1Data))
		m.MakeReplicaSetV1Data(kresource.Apps.(*types.AppsV1Data))
	}
	if kresource.Networking != nil && m.ClientInfo.GetVersions().GetVersion("networking").Version == "v1" {
		m.MakeIngressV1Data(kresource.Networking.(*types.NetworkingV1Data))
	}
}

func (m *MappingResourceSource) MakeNamespaceV1Data(kcore *types.CoreV1Data) {
	defer ErrorRecover()

	if kcore.Namespace != nil && kcore.Namespace.Items != nil {
		m.Namespaces = make([]MappingNamespace, 0)
		for _, n := range kcore.Namespace.Items {
			host, _, _ := net.SplitHostPort(m.ClientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)
			m.Namespaces = append(m.Namespaces, MappingNamespace{
				Name:      n.GetObjectMeta().GetName(),
				UID:       string(n.GetObjectMeta().GetUID()),
				StartTime: n.GetObjectMeta().GetCreationTimestamp().Time,
				Host:      host,
				Labels:    GetLabelSelector(n.GetObjectMeta().GetLabels()),
				Status:    string(n.Status.Phase),
			})
		}

		if len(m.Namespaces) > 0 {
			ns_json, _ := json.Marshal(m.Namespaces)
			ns_resource_map := make(map[string][]byte)
			ns_resource_map["namespaceinfo"] = ns_json
			common.ChannelResourceData <- ns_resource_map
		}
	}
}

func (m *MappingResourceSource) MakeNodeV1Data(kcore *types.CoreV1Data) {
	defer ErrorRecover()

	m.Nodes = make([]MappingNode, 0)

	if kcore.Node != nil && kcore.Node.Items != nil {
		for _, n := range kcore.Node.Items {
			var nodetype string
			if _, ok := n.GetObjectMeta().GetLabels()["node-role.kubernetes.io/control-plane"]; ok {
				nodetype = "Control Plane"
				m.GlobalInfo.masternodename = n.GetObjectMeta().GetName()
			} else {
				nodetype = "Worker"
			}
			var status int = 0
			for _, condition := range n.Status.Conditions {
				if condition.Type == "Ready" {
					if condition.Status == corev1.ConditionTrue {
						status = 1
					}

					break
				}
			}
			host, _, _ := net.SplitHostPort(m.ClientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)
			m.Nodes = append(m.Nodes, MappingNode{
				Name:                    n.GetObjectMeta().GetName(),
				UID:                     string(n.GetObjectMeta().GetUID()),
				StartTime:               n.GetObjectMeta().GetCreationTimestamp().Time,
				Host:                    host,
				Labels:                  GetLabelSelector(n.GetObjectMeta().GetLabels()),
				NodeType:                nodetype,
				KernelVersion:           n.Status.NodeInfo.KernelVersion,
				OSImage:                 n.Status.NodeInfo.OSImage,
				OSName:                  n.Status.NodeInfo.OperatingSystem,
				ContainerRuntimeVersion: n.Status.NodeInfo.ContainerRuntimeVersion,
				KubeletVersion:          n.Status.NodeInfo.KubeletVersion,
				KubeProxyVersion:        n.Status.NodeInfo.KubeProxyVersion,
				CPUArch:                 n.Status.NodeInfo.Architecture,
				CPUCount:                n.Status.Capacity.Cpu().Value(),
				EphemeralStorage:        n.Status.Capacity.StorageEphemeral().Value(),
				MemorySize:              n.Status.Capacity.Memory().Value(),
				Pods:                    n.Status.Capacity.Pods().Value(),
				IP:                      n.Status.Addresses[0].Address,
				Status:                  status,
			})
		}

		if len(m.Nodes) > 0 {
			node_json, _ := json.Marshal(m.Nodes)
			node_resource_map := make(map[string][]byte)
			node_resource_map["nodeinfo"] = node_json
			common.ChannelResourceData <- node_resource_map
		}
	}
}

func (m *MappingResourceSource) MakePersistentVolumeClaimV1Data(kcore *types.CoreV1Data) {
	defer ErrorRecover()

	if kcore.PersistentVolumeClaim != nil && kcore.PersistentVolumeClaim.Items != nil {
		m.Pvcs = make([]MappingPvc, 0)
		for _, n := range kcore.PersistentVolumeClaim.Items {
			host, _, _ := net.SplitHostPort(m.ClientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)
			var scname string
			if n.Spec.StorageClassName != nil {
				scname = *n.Spec.StorageClassName
			}
			m.Pvcs = append(m.Pvcs, MappingPvc{
				Name:             n.GetObjectMeta().GetName(),
				UID:              string(n.GetObjectMeta().GetUID()),
				StartTime:        n.GetObjectMeta().GetCreationTimestamp().Time,
				Host:             host,
				Labels:           GetLabelSelector(n.GetObjectMeta().GetLabels()),
				Selector:         GetSelector(n.Spec.Selector),
				RequestStorage:   n.Spec.Resources.Requests.Storage().Value(),
				AccessModes:      GetAccessModes(n.Spec.AccessModes),
				Status:           string(n.Status.Phase),
				StorageClassName: scname,
				NamespaceName:    n.GetObjectMeta().GetNamespace(),
			})
		}

		if len(m.Pvcs) > 0 {
			pvc_json, _ := json.Marshal(m.Pvcs)
			pvc_resource_map := make(map[string][]byte)
			pvc_resource_map["pvcinfo"] = pvc_json
			common.ChannelResourceData <- pvc_resource_map
		}
	}
}

func (m *MappingResourceSource) MakePersistentVolumeV1Data(kcore *types.CoreV1Data) {
	defer ErrorRecover()

	if kcore.PersistentVolume != nil && kcore.PersistentVolume.Items != nil {
		m.Pvs = make([]MappingPv, 0)
		for _, n := range kcore.PersistentVolume.Items {
			host, _, _ := net.SplitHostPort(m.ClientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)
			m.Pvs = append(m.Pvs, MappingPv{
				Name:          n.GetObjectMeta().GetName(),
				UID:           string(n.GetObjectMeta().GetUID()),
				StartTime:     n.GetObjectMeta().GetCreationTimestamp().Time,
				Host:          host,
				Labels:        GetLabelSelector(n.GetObjectMeta().GetLabels()),
				AccessModes:   GetAccessModes(n.Spec.AccessModes),
				Capacity:      n.Spec.Capacity.Storage().Value(),
				ReclaimPolicy: string(n.Spec.PersistentVolumeReclaimPolicy),
				Status:        string(n.Status.Phase),
				PvcUID:        string(n.Spec.ClaimRef.UID),
			})
		}

		if len(m.Pvs) > 0 {
			pvs_json, _ := json.Marshal(m.Pvs)
			pvs_resource_map := make(map[string][]byte)
			pvs_resource_map["pvinfo"] = pvs_json
			common.ChannelResourceData <- pvs_resource_map
		}
	}
}

func (m *MappingResourceSource) MakePodV1Data(kcore *types.CoreV1Data) {
	defer ErrorRecover()
	pod_name_array := make([]string, 0)

	if kcore.Pod != nil && kcore.Pod.Items != nil {
		m.Pods = make([]MappingPod, 0)
		for _, n := range kcore.Pod.Items {
			var nodename = n.Spec.NodeName
			var podname = n.GetObjectMeta().GetName()
			var poduid = string(n.GetObjectMeta().GetUID())
			host, _, _ := net.SplitHostPort(m.ClientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)
			key := fmt.Sprintf("%s/%s", n.GetObjectMeta().GetNamespace(), n.GetObjectMeta().GetName())
			pod_name_array = append(pod_name_array, key)

			var annotation string
			if _, ok := n.GetObjectMeta().GetAnnotations()["kubernetes.io/config.hash"]; ok {
				annotation = n.GetObjectMeta().GetAnnotations()["kubernetes.io/config.hash"]
			} else if _, ok := n.GetObjectMeta().GetAnnotations()["kubernetes.io/config.mirror"]; ok {
				annotation = n.GetObjectMeta().GetAnnotations()["kubernetes.io/config.mirror"]
			}

			var restartcount int32
			var restarttime time.Time
			for _, c := range n.Status.ContainerStatuses {
				restartcount += c.RestartCount
				if c.State.Running != nil && (restarttime.IsZero() || c.State.Running.StartedAt.Time.Unix() > restarttime.Unix()) {
					restarttime = c.State.Running.StartedAt.Time
				}
			}

			podconditions := make([]string, 0)
			for _, c := range n.Status.Conditions {
				if c.Message != "" || c.Status == corev1.ConditionFalse {
					podcondition := fmt.Sprintf("%s(%s/%s):%s", c.Type, c.Status, c.LastTransitionTime.Format(time.RFC3339), c.Message)
					podconditions = append(podconditions, podcondition)
				}
			}

			var staticpodname string
			if nodename == m.GlobalInfo.masternodename && len(podname) > len(nodename) && podname[len(podname)-len(nodename):] == nodename {
				staticpodname = podname[:len(podname)-len(nodename)-1]
			}

			container_status := make(map[string]ContainerStatus)
			for _, cs := range n.Status.ContainerStatuses {
				var container_state string
				if cs.State.Running != nil {
					container_state = "running"
				} else if cs.State.Waiting != nil {
					container_state = "waiting"
				} else if cs.State.Terminated != nil {
					container_state = "terminated"
				}

				container_status[cs.Name] = ContainerStatus{
					State:        container_state,
					RestartCount: cs.RestartCount,
					Ready:        cs.Ready,
					Started:      *cs.Started,
				}
			}

			containers := make([]MappingContainer, 0)
			for _, c := range n.Spec.Containers {
				cports := make([]string, 0)
				for _, cp := range c.Ports {
					var cpn string
					if cp.Name != "" {
						cpn = cp.Name
					}

					var cphp string
					if cp.HostPort > 0 {
						cphp = fmt.Sprintf("-%d", cp.HostPort)
					}

					var cpflag string
					if cpn != "" || cphp != "" {
						cpflag = "/"
					}

					var cpp string
					if cp.Protocol == "" {
						cpp = "TCP"
					} else {
						cpp = string(cp.Protocol)
					}

					var cphi string
					if cp.HostIP != "" {
						cphi = fmt.Sprintf(":%s", cp.HostIP)
					}

					cports = append(cports, fmt.Sprintf("%s%s%s%s%s", cpn, cphp, cpflag, cpp, cphi))
				}

				cenvs := make([]string, 0)
				for _, ce := range c.Env {
					cenvs = append(cenvs, fmt.Sprintf("%s:%s", ce.Name, ce.Value))
				}

				cvolumemounts := make([]string, 0)
				for _, cvm := range c.VolumeMounts {
					var cvmro string
					if cvm.ReadOnly {
						cvmro = "(RO)"
					}

					// SubPathExpr and SubPath are mutually exclusive.
					var cvmsp string
					if cvm.SubPath != "" {
						cvmsp = fmt.Sprintf(":%s", cvm.SubPath)
					} else if cvm.SubPathExpr != "" {
						cvmsp = fmt.Sprintf(":%s", cvm.SubPathExpr)
					}
					cvolumemounts = append(cvolumemounts, fmt.Sprintf("%s%s(%s%s)", cvm.Name, cvmro, cvm.MountPath, cvmsp))
				}
				var cstate string
				if _, ok := container_status[c.Name]; ok {
					cstate = container_status[c.Name].State
				}

				containers = append(containers, MappingContainer{
					Name:             c.Name,
					UID:              poduid,
					Host:             host,
					Image:            c.Image,
					Ports:            strings.Join(cports, ","),
					Env:              strings.Join(cenvs, ","),
					LimitCpu:         c.Resources.Limits.Cpu().Value(),
					LimitMemory:      c.Resources.Limits.Memory().Value(),
					LimitStorage:     c.Resources.Limits.Storage().Value(),
					LimitEphemeral:   c.Resources.Limits.StorageEphemeral().Value(),
					RequestCpu:       c.Resources.Requests.Cpu().Value(),
					RequestMemory:    c.Resources.Requests.Memory().Value(),
					RequestStorage:   c.Resources.Requests.Storage().Value(),
					RequestEphemeral: c.Resources.Requests.StorageEphemeral().Value(),
					VolumeMounts:     strings.Join(cvolumemounts, ","),
					State:            cstate,
				})
			}

			var pvc string
			for _, v := range n.Spec.Volumes {
				if v.PersistentVolumeClaim != nil {
					pvc = v.PersistentVolumeClaim.ClaimName
				}
			}

			refs := GetRefData(n.GetObjectMeta().GetOwnerReferences())
			m.Pods = append(m.Pods, MappingPod{
				Name:                  podname,
				UID:                   poduid,
				StartTime:             n.GetObjectMeta().GetCreationTimestamp().Time,
				Host:                  host,
				Labels:                GetLabelSelector(n.GetObjectMeta().GetLabels()),
				Selector:              GetLabelSelector(n.Spec.NodeSelector),
				AnnotationUID:         annotation,
				RestartPolicy:         string(n.Spec.RestartPolicy),
				ServiceAccount:        n.Spec.ServiceAccountName,
				Status:                string(n.Status.Phase),
				HostIP:                n.Status.HostIP,
				PodIP:                 n.Status.PodIP,
				RestartCount:          restartcount,
				RestartTime:           restarttime,
				Condition:             strings.Join(podconditions, ","),
				StaticPod:             staticpodname,
				NodeName:              nodename,
				NamespaceName:         n.GetObjectMeta().GetNamespace(),
				ReferenceKind:         refs["kind"],
				ReferenceUID:          refs["UID"],
				PersistentVolumeClaim: pvc,
				Containers:            containers,
			})
		}
		common.PodListMap.Store(m.ClientInfo.GetHost(), pod_name_array)

		if len(m.Pods) > 0 {
			pod_json, _ := json.Marshal(m.Pods)
			pod_resource_map := make(map[string][]byte)
			pod_resource_map["podsinfo"] = pod_json
			common.ChannelResourceData <- pod_resource_map
		}
	}
}

func (m *MappingResourceSource) MakeServiceV1Data(kcore *types.CoreV1Data) {
	defer ErrorRecover()

	if kcore.Service != nil && kcore.Service.Items != nil {
		m.Services = make([]MappingService, 0)
		for _, n := range kcore.Service.Items {
			ports := make([]string, 0)
			for _, p := range n.Spec.Ports {
				var pt string
				if p.TargetPort.StrVal == "" || p.TargetPort.IntVal == p.Port {
					pt = fmt.Sprintf("%d", p.Port)
				} else {
					pt = fmt.Sprintf("%d:%s", p.Port, p.TargetPort.StrVal)
				}

				var npt string
				if p.NodePort != 0 {
					npt = fmt.Sprintf(":%d", p.NodePort)
				}
				port := fmt.Sprintf("%s(%s%s/%s)", p.Name, pt, npt, p.Protocol)
				ports = append(ports, port)
			}

			host, _, _ := net.SplitHostPort(m.ClientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)
			m.Services = append(m.Services, MappingService{
				Name:          n.GetObjectMeta().GetName(),
				UID:           string(n.GetObjectMeta().GetUID()),
				StartTime:     n.GetObjectMeta().GetCreationTimestamp().Time,
				Host:          host,
				Labels:        GetLabelSelector(n.GetObjectMeta().GetLabels()),
				Selector:      GetLabelSelector(n.Spec.Selector),
				ServiceType:   string(n.Spec.Type),
				ClusterIP:     n.Spec.ClusterIP,
				Ports:         strings.Join(ports, ","),
				NamespaceName: n.GetObjectMeta().GetNamespace(),
			})
		}

		if len(m.Services) > 0 {
			svc_json, _ := json.Marshal(m.Services)
			svc_resource_map := make(map[string][]byte)
			svc_resource_map["serviceinfo"] = svc_json
			common.ChannelResourceData <- svc_resource_map
		}
	}
}

func (m *MappingResourceSource) MakeDeploymentV1Data(kapps *types.AppsV1Data) {
	defer ErrorRecover()

	if kapps.Deployment != nil && kapps.Deployment.Items != nil {
		m.Deployments = make([]MappingDeployment, 0)
		for _, n := range kapps.Deployment.Items {
			host, _, _ := net.SplitHostPort(m.ClientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)
			m.Deployments = append(m.Deployments, MappingDeployment{
				Name:                 n.GetObjectMeta().GetName(),
				UID:                  string(n.GetObjectMeta().GetUID()),
				StartTime:            n.GetObjectMeta().GetCreationTimestamp().Time,
				Host:                 host,
				Labels:               GetLabelSelector(n.GetObjectMeta().GetLabels()),
				Selector:             GetSelector(n.Spec.Selector),
				ServiceAccount:       n.Spec.Template.Spec.ServiceAccountName,
				Replicas:             *n.Spec.Replicas,
				UpdatedReplicas:      n.Status.UpdatedReplicas,
				ReadyReplicas:        n.Status.ReadyReplicas,
				AvailableReplicas:    n.Status.AvailableReplicas,
				ObservedGeneneration: n.Status.ObservedGeneration,
				NamespaceName:        n.GetObjectMeta().GetNamespace(),
			})
		}

		if len(m.Deployments) > 0 {
			deploy_json, _ := json.Marshal(m.Deployments)
			deploy_resource_map := make(map[string][]byte)
			deploy_resource_map["deployinfo"] = deploy_json
			common.ChannelResourceData <- deploy_resource_map
		}
	}
}

func (m *MappingResourceSource) MakeStatefulSetV1Data(kapps *types.AppsV1Data) {
	defer ErrorRecover()

	if kapps.StatefulSet != nil && kapps.StatefulSet.Items != nil {
		m.StatefulSets = make([]MappingStatefulSet, 0)
		for _, n := range kapps.StatefulSet.Items {
			host, _, _ := net.SplitHostPort(m.ClientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)
			m.StatefulSets = append(m.StatefulSets, MappingStatefulSet{
				Name:              n.GetObjectMeta().GetName(),
				UID:               string(n.GetObjectMeta().GetUID()),
				StartTime:         n.GetObjectMeta().GetCreationTimestamp().Time,
				Host:              host,
				Labels:            GetLabelSelector(n.GetObjectMeta().GetLabels()),
				Selector:          GetSelector(n.Spec.Selector),
				ServiceAccount:    n.Spec.Template.Spec.ServiceAccountName,
				Replicas:          *n.Spec.Replicas,
				UpdatedReplicas:   n.Status.UpdatedReplicas,
				ReadyReplicas:     n.Status.ReadyReplicas,
				AvailableReplicas: n.Status.AvailableReplicas,
				NamespaceName:     n.GetObjectMeta().GetNamespace(),
			})
		}

		if len(m.StatefulSets) > 0 {
			stateful_json, _ := json.Marshal(m.StatefulSets)
			stateful_resource_map := make(map[string][]byte)
			stateful_resource_map["statefulset"] = stateful_json
			common.ChannelResourceData <- stateful_resource_map
		}
	}
}

func (m *MappingResourceSource) MakeDaemonSetV1Data(kapps *types.AppsV1Data) {
	defer ErrorRecover()

	if kapps.DaemonSet != nil && kapps.DaemonSet.Items != nil {
		m.DaemonSets = make([]MappingDaemonSet, 0)
		for _, n := range kapps.DaemonSet.Items {
			host, _, _ := net.SplitHostPort(m.ClientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)
			m.DaemonSets = append(m.DaemonSets, MappingDaemonSet{
				Name:                   n.GetObjectMeta().GetName(),
				UID:                    string(n.GetObjectMeta().GetUID()),
				StartTime:              n.GetObjectMeta().GetCreationTimestamp().Time,
				Host:                   host,
				Labels:                 GetLabelSelector(n.GetObjectMeta().GetLabels()),
				Selector:               GetSelector(n.Spec.Selector),
				ServiceAccount:         n.Spec.Template.Spec.ServiceAccountName,
				CurrentNumberScheduled: n.Status.CurrentNumberScheduled,
				DesiredNumberScheduled: n.Status.DesiredNumberScheduled,
				NumberReady:            n.Status.NumberReady,
				UpdatedNumberScheduled: n.Status.UpdatedNumberScheduled,
				NumberAvailable:        n.Status.NumberAvailable,
				NamespaceName:          n.GetObjectMeta().GetNamespace(),
			})
		}

		if len(m.DaemonSets) > 0 {
			daemonset_json, _ := json.Marshal(m.DaemonSets)
			daemonset_resource_map := make(map[string][]byte)
			daemonset_resource_map["daemonset"] = daemonset_json
			common.ChannelResourceData <- daemonset_resource_map
		}
	}
}

func (m *MappingResourceSource) MakeReplicaSetV1Data(kapps *types.AppsV1Data) {
	defer ErrorRecover()

	if kapps.ReplicaSet != nil && kapps.ReplicaSet.Items != nil {
		m.ReplicaSets = make([]MappingReplicaSet, 0)
		for _, n := range kapps.ReplicaSet.Items {
			refs := GetRefData(n.GetObjectMeta().GetOwnerReferences())
			host, _, _ := net.SplitHostPort(m.ClientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)
			m.ReplicaSets = append(m.ReplicaSets, MappingReplicaSet{
				Name:                 n.GetObjectMeta().GetName(),
				UID:                  string(n.GetObjectMeta().GetUID()),
				StartTime:            n.GetObjectMeta().GetCreationTimestamp().Time,
				Host:                 host,
				Labels:               GetLabelSelector(n.GetObjectMeta().GetLabels()),
				Selector:             GetSelector(n.Spec.Selector),
				Replicas:             *n.Spec.Replicas,
				FullyLabeledReplicas: n.Status.FullyLabeledReplicas,
				ReadyReplicas:        n.Status.ReadyReplicas,
				AvailableReplicas:    n.Status.AvailableReplicas,
				ReferenceKind:        refs["kind"],
				ReferenceUID:         refs["UID"],
				NamespaceName:        n.GetObjectMeta().GetNamespace(),
			})
		}

		if len(m.ReplicaSets) > 0 {
			replicaset_json, _ := json.Marshal(m.ReplicaSets)
			replicaset_resource_map := make(map[string][]byte)
			replicaset_resource_map["replicaset"] = replicaset_json
			common.ChannelResourceData <- replicaset_resource_map
		}
	}
}

func (m *MappingResourceSource) MakeIngressV1Data(knetworking *types.NetworkingV1Data) {
	defer ErrorRecover()

	if knetworking.Ingress != nil && knetworking.Ingress.Items != nil {
		m.Ingresses = make([]MappingIngress, 0)
		for _, n := range knetworking.Ingress.Items {
			host, _, _ := net.SplitHostPort(m.ClientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)

			ingresshosts := make([]MappingIngressHost, 0)
			if spec := n.Spec.DefaultBackend; spec != nil {
				ingress := MappingIngressHost{
					Hostname: "*",
					Host:     host,
					PathType: "",
					Path:     "",
					UID:      string(n.GetObjectMeta().GetUID()),
				}

				if svc := spec.Service; svc != nil {
					ingress.BackendType = "service"
					ingress.BackendName = svc.Name
					ingress.ServicePort = svc.Port.Number
				} else if rsc := spec.Resource; rsc != nil {
					ingress.BackendType = "resource"
					ingress.BackendName = rsc.Name
					ingress.ResourceAPIGroup = *rsc.APIGroup
					ingress.ResourceKind = rsc.Kind
				}

				ingresshosts = append(ingresshosts, ingress)
			} else {
				for _, r := range n.Spec.Rules {
					hostname := r.Host

					for _, p := range r.HTTP.Paths {
						ingress := MappingIngressHost{
							Hostname: hostname,
							Host:     host,
							PathType: string(*p.PathType),
							Path:     p.Path,
							UID:      string(n.GetObjectMeta().GetUID()),
						}

						if svc := p.Backend.Service; svc != nil {
							ingress.BackendType = "service"
							ingress.BackendName = svc.Name
							ingress.ServicePort = svc.Port.Number
						} else if rsc := p.Backend.Resource; rsc != nil {
							ingress.BackendType = "resource"
							ingress.BackendName = rsc.Name
							ingress.ResourceAPIGroup = *rsc.APIGroup
							ingress.ResourceKind = rsc.Kind
						}

						ingresshosts = append(ingresshosts, ingress)
					}
				}
			}

			m.Ingresses = append(m.Ingresses, MappingIngress{
				Name:             n.GetObjectMeta().GetName(),
				UID:              string(n.GetObjectMeta().GetUID()),
				StartTime:        n.GetObjectMeta().GetCreationTimestamp().Time,
				Host:             host,
				Labels:           GetLabelSelector(n.GetObjectMeta().GetLabels()),
				IngressClassName: *n.Spec.IngressClassName,
				NamespaceName:    n.GetObjectMeta().GetNamespace(),
				IngressHosts:     ingresshosts,
			})
		}

		if len(m.Ingresses) > 0 {
			ingress_json, _ := json.Marshal(m.Ingresses)
			ingress_resource_map := make(map[string][]byte)
			ingress_resource_map["ingress_info"] = ingress_json
			common.ChannelResourceData <- ingress_resource_map
		}
	}
}

func (m *MappingResourceSource) MakeStorageClassV1Data(kstorage *types.StorageV1Data) {
	defer ErrorRecover()

	if kstorage.StorageClass != nil && kstorage.StorageClass.Items != nil {
		m.StorageClasses = make([]MappingStorageClass, 0)
		for _, n := range kstorage.StorageClass.Items {
			var ave bool
			if n.AllowVolumeExpansion != nil {
				ave = *n.AllowVolumeExpansion
			}

			host, _, _ := net.SplitHostPort(m.ClientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)
			m.StorageClasses = append(m.StorageClasses, MappingStorageClass{
				Name:                 n.GetObjectMeta().GetName(),
				UID:                  string(n.GetObjectMeta().GetUID()),
				StartTime:            n.GetObjectMeta().GetCreationTimestamp().Time,
				Host:                 host,
				Labels:               GetLabelSelector(n.GetObjectMeta().GetLabels()),
				Provisioner:          n.Provisioner,
				ReclaimPolicy:        string(*n.ReclaimPolicy),
				VolumeBindingMode:    string(*n.VolumeBindingMode),
				AllowVolumeExpansion: ave,
			})
		}

		if len(m.StorageClasses) > 0 {
			sc_json, _ := json.Marshal(m.StorageClasses)
			sc_resource_map := make(map[string][]byte)
			sc_resource_map["sc"] = sc_json
			common.ChannelResourceData <- sc_resource_map
		}
	}
}
