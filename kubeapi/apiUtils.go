package kubeapi

import (
	"bytes"
	"fmt"
	"net"
	"onTuneKubeManager/common"
	"onTuneKubeManager/config"
	"sort"
	"strings"
	"sync"
	"time"

	dto "github.com/prometheus/client_model/go"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	EnabledResources = []string{"Namespace", "Node", "Pod", "Service", "Ingress", "Deployment", "StatefulSet", "DaemonSet", "ReplicaSet", "Event", "PersistentVolumeClaim", "StorageClass", "PersistentVolume"}
	ClientInfoArr    []*ClientInfo
)

func ErrorRecover() {
	if r := recover(); r != nil {
		err := fmt.Errorf("Kubernetes API Error - %v", r)
		common.LogManager.WriteLog(err.Error())
	}
}

func GetRefData(refs []metav1.OwnerReference) map[string]string {
	if refs != nil {
		return map[string]string{
			"kind": refs[0].Kind,
			"UID":  string(refs[0].UID),
			"name": refs[0].Name,
		}
	} else {
		return map[string]string{
			"kind": "",
			"UID":  "",
			"name": "",
		}
	}
}

func GetSelector(selector *metav1.LabelSelector) string {
	var selector_str string

	if selector != nil {
		selector_str = GetLabelSelector(selector.MatchLabels)
	}

	return selector_str
}

func GetLabelSelector(dtmap map[string]string) string {
	sortKeys := make([]string, 0, len(dtmap))
	for k := range dtmap {
		sortKeys = append(sortKeys, k)
	}
	sort.Strings(sortKeys)

	b := new(bytes.Buffer)
	for _, k := range sortKeys {
		fmt.Fprintf(b, "\"%s\":\"%s\",", k, dtmap[k])
	}

	return b.String()
}

func GetAccessModes(acs_source []corev1.PersistentVolumeAccessMode) string {
	accessmodes := make([]string, 0)
	for _, ac := range acs_source {
		switch ac {
		case corev1.ReadWriteOnce:
			accessmodes = append(accessmodes, "RWO")
		case corev1.ReadOnlyMany:
			accessmodes = append(accessmodes, "ROM")
		case corev1.ReadWriteOncePod:
			accessmodes = append(accessmodes, "RWOP")
		case corev1.ReadWriteMany:
			accessmodes = append(accessmodes, "RWM")
		}
	}

	return strings.Join(accessmodes, ",")
}

func Contains(name string) bool {
	for _, e := range EnabledResources {
		if e == name {
			return true
		}
	}

	return false
}

type LabelFilter struct {
	namespace string
	pod       string
	container string
	id        string
	itfc      string
	device    string
	image     string
	flag      bool
}

func searchLabel(metric *dto.Metric, mtype string) LabelFilter {
	labelMap := make(map[string]string)
	for _, l := range metric.Label {
		labelMap[*l.Name] = *l.Value
	}

	if _, ok := labelMap["container"]; ok {
		switch mtype {
		case "node":
			if _, ok := labelMap["device"]; (!ok || labelMap["device"] == "") && labelMap["container"] == "" && labelMap["pod"] == "" {
				return LabelFilter{
					id:    labelMap["id"],
					flag:  true,
					image: labelMap["image"],
				}
			}
		case "pod":
			if labelMap["container"] == "" && labelMap["pod"] != "" {
				return LabelFilter{
					namespace: labelMap["namespace"],
					pod:       labelMap["pod"],
					id:        labelMap["id"],
					flag:      true,
					image:     labelMap["image"],
				}
			}
		case "container":
			if labelMap["container"] != "" {
				return LabelFilter{
					namespace: labelMap["namespace"],
					pod:       labelMap["pod"],
					container: labelMap["container"],
					id:        labelMap["id"],
					flag:      true,
					image:     labelMap["image"],
				}
			}
		case "nodenetwork":
			if _, ok := labelMap["interface"]; ok && labelMap["interface"] != "" && labelMap["container"] == "" && labelMap["pod"] == "" {
				return LabelFilter{
					id:    labelMap["id"],
					itfc:  labelMap["interface"],
					flag:  true,
					image: labelMap["image"],
				}
			}
		case "podnetwork":
			if _, ok := labelMap["interface"]; ok && labelMap["interface"] != "" && labelMap["container"] == "" && labelMap["pod"] != "" {
				return LabelFilter{
					namespace: labelMap["namespace"],
					pod:       labelMap["pod"],
					id:        labelMap["id"],
					itfc:      labelMap["interface"],
					flag:      true,
					image:     labelMap["image"],
				}
			}
		case "nodefilesystem":
			if _, ok := labelMap["device"]; ok && labelMap["device"] != "" && labelMap["container"] == "" && labelMap["pod"] == "" {
				return LabelFilter{
					id:     labelMap["id"],
					device: labelMap["device"],
					flag:   true,
					image:  labelMap["image"],
				}
			}
		case "podfilesystem":
			if _, ok := labelMap["device"]; ok && labelMap["device"] != "" && labelMap["container"] == "" && labelMap["pod"] != "" {
				return LabelFilter{
					namespace: labelMap["namespace"],
					pod:       labelMap["pod"],
					id:        labelMap["id"],
					device:    labelMap["device"],
					flag:      true,
					image:     labelMap["image"],
				}
			}
		case "containerfilesystem":
			if _, ok := labelMap["device"]; ok && labelMap["device"] != "" && labelMap["container"] != "" {
				return LabelFilter{
					namespace: labelMap["namespace"],
					pod:       labelMap["pod"],
					container: labelMap["container"],
					id:        labelMap["id"],
					device:    labelMap["device"],
					flag:      true,
					image:     labelMap["image"],
				}
			}
		}
	}

	return LabelFilter{flag: false}
}

func SplitPodName(podName string) (string, string) {
	s := strings.Split(podName, "/")
	if len(s) != 2 {
		return "", ""
	}

	return s[0], s[1]
}

func writeLog(message string) {
	common.LogManager.WriteLog(message)
}

func ReceiveChannelRequestChange() {
	for {
		host := <-common.ChannelRequestChangeHost

		for _, clientInfo := range ClientInfoArr {
			if clientInfo.GetHost() == host {
				clientInfo.ChangeResourceInterval()
				break
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func SetClientInfo(clustercount uint16, clientset []*kubernetes.Clientset, conf *config.ManagerConfig) {
	ClientInfoArr = make([]*ClientInfo, 0)
	go ReceiveChannelRequestChange()

	for i := 0; i < int(common.KubeClusterCount); i++ {
		host, _, _ := net.SplitHostPort(clientset[i].DiscoveryClient.RESTClient().Get().URL().Host)

		var clusterEnabled bool
		if ce, ok := common.ClusterStatusMap.Load(host); ok {
			clusterEnabled = ce.(bool)
		} else {
			common.LogManager.WriteLog(fmt.Sprintf("Kubernetes Cluster %s is not enabled", host))
			continue
		}

		clientInfo := &ClientInfo{
			Clientset:               clientset[i],
			Host:                    host,
			KubeCfg:                 &conf.GetKubeConfig()[i],
			Enabled:                 clusterEnabled,
			RetryConnectionInterval: conf.GetRetryConnectionInterval(),
			ResourceChangeInterval:  conf.GetResourceChangeInterval(),
			ResourceProcessFlag:     false,
			RequestFlag:             false,
			RequestMutex:            &sync.Mutex{},
		}

		clientInfo.GetInitData()

		go clientInfo.RecoverResourceInterval()
		go clientInfo.GetResourceDataInterval()
		go clientInfo.GetEventlogDataInterval()
		go clientInfo.GetMetricDataInterval()

		ClientInfoArr = append(ClientInfoArr, clientInfo)
	}
}
