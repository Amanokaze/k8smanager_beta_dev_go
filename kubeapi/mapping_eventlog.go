package kubeapi

import (
	"encoding/json"
	"net"
	"onTuneKubeManager/common"
	"onTuneKubeManager/kubeapi/types"
	"strings"
	"sync"
	"time"
)

const (
	LOG_DATE_FORMAT      = "2006-01-02T15:04:05.999999999-07:00"
	LOG_TIMESTAMP_LENGTH = len(LOG_DATE_FORMAT)
)

type MappingEvent struct {
	Name            string    `json:"name"`
	UID             string    `json:"uid"`
	Host            string    `json:"host"`
	Firsttime       time.Time `json:"firsttime"`
	Lasttime        time.Time `json:"lasttime"`
	Labels          string    `json:"labels"`
	Eventtype       string    `json:"eventtype"`
	Eventcount      int       `json:"eventcount"`
	ObjectKind      string    `json:"objectkind"`
	ObjectNamespace string    `json:"objectns"`
	ObjectName      string    `json:"objectname"`
	SourceComponent string    `json:"srccomponent"`
	SourceHost      string    `json:"srchost"`
	Reason          string    `json:"reason"`
	Message         string    `json:"message"`
	NamespaceName   string    `json:"nsname"`
}

type MappingLog struct {
	Logtype       string    `json:"logtype"`
	Host          string    `json:"host"`
	NamespaceName string    `json:"nsname"`
	Name          string    `json:"name"`
	Starttime     time.Time `json:"starttime"`
	Message       string    `json:"message"`
}

type MappingEventlogSource struct {
	ClientInfo `json:"client"`
	Events     []MappingEvent `json:"events"`
	Logs       []MappingLog   `json:"logs"`
}

func (m *MappingEventlogSource) MakeEventlogData(keventlog *KubernetesAPIEventlogSource) {
	if m.ClientInfo.GetVersions().GetVersion("core").Version == "v1" {
		m.MakeEventV1Data(keventlog.Event.(*types.CoreV1EventlogData))
	}

	if keventlog.Logs != nil {
		m.MakePodLogData(keventlog.Logs.(*sync.Map))
	}
}

func (m *MappingEventlogSource) MakeEventV1Data(kevent *types.CoreV1EventlogData) {
	defer ErrorRecover()

	m.Events = make([]MappingEvent, 0)
	for _, n := range kevent.Event.Items {
		host, _, _ := net.SplitHostPort(m.ClientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)

		m.Events = append(m.Events, MappingEvent{
			Name:            n.GetObjectMeta().GetName(),
			UID:             string(n.GetObjectMeta().GetUID()),
			Host:            host,
			Firsttime:       n.FirstTimestamp.Time,
			Lasttime:        n.LastTimestamp.Time,
			Labels:          GetLabelSelector(n.GetObjectMeta().GetLabels()),
			Eventtype:       n.Type,
			Eventcount:      int(n.Count),
			ObjectKind:      n.InvolvedObject.Kind,
			ObjectNamespace: n.InvolvedObject.Namespace,
			ObjectName:      n.InvolvedObject.Name,
			SourceComponent: n.Source.Component,
			SourceHost:      n.Source.Host,
			Reason:          n.Reason,
			Message:         n.Message,
			NamespaceName:   n.GetObjectMeta().GetNamespace(),
		})
	}

	event_json, _ := json.Marshal(m.Events)
	event_map := make(map[string][]byte)
	event_map["event"] = event_json
	common.ChannelEventlogData <- event_map
}

func (m *MappingEventlogSource) MakePodLogData(klog *sync.Map) {
	m.Logs = make([]MappingLog, 0)
	klog.Range(func(key, value interface{}) bool {
		log := value.(string)
		log_lines := strings.Split(log, "\r\n")
		for _, line := range log_lines {
			date, err := time.Parse(LOG_DATE_FORMAT, line[:LOG_TIMESTAMP_LENGTH])
			if err != nil {
				continue
			}

			host, _, _ := net.SplitHostPort(m.ClientInfo.GetClientset().DiscoveryClient.RESTClient().Get().URL().Host)

			namespace_name, pod_name := SplitPodName(key.(string))
			m.Logs = append(m.Logs, MappingLog{
				Logtype:       "pod",
				Host:          host,
				NamespaceName: namespace_name,
				Name:          pod_name,
				Starttime:     date,
				Message:       line[LOG_TIMESTAMP_LENGTH+1:],
			})
		}

		return true
	})

	if len(m.Logs) == 0 {
		return
	}

	pod_log_json, _ := json.Marshal(m.Logs)
	pod_log_map := make(map[string][]byte)
	pod_log_map["log"] = pod_log_json
	common.ChannelEventlogData <- pod_log_map
}
