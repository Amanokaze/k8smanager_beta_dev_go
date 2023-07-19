package kubeapi

import (
	"bytes"
	"context"
	"fmt"
	"onTuneKubeManager/common"
	"strings"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

type KubernetesAPIMetricSource struct {
	NodeMetric map[string]map[string]*dto.MetricFamily
	clientInfo *ClientInfo
	mutex      *sync.Mutex
}

func (k *KubernetesAPIMetricSource) GetData(clientInfo *ClientInfo) {
	k.clientInfo = clientInfo
	k.mutex = &sync.Mutex{}

	if clientInfo.GetVersions().GetVersion("core").Version == "v1" {
		countFlag := true
		k.NodeMetric = make(map[string]map[string]*dto.MetricFamily)
		nodes, err := clientInfo.GetClientset().CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		k.clientInfo.ErrorDisabledCheck(err, &countFlag)
		if err != nil {
			return
		}

		var wg sync.WaitGroup
		wg.Add(len(nodes.Items))

		for _, n := range nodes.Items {
			go k.GetMetricNodeData(n.Name, &wg)
		}
		wg.Wait()
	}
}

func (k *KubernetesAPIMetricSource) GetMetricNodeData(nodename string, wg *sync.WaitGroup) {
	writeLog(fmt.Sprintf("Kubernetes Get Node %s %s Metric Data start", k.clientInfo.GetVersions().GetVersion("core").Version, nodename))
	countFlag := true

	node, err := k.clientInfo.GetClientset().CoreV1().Nodes().Get(context.Background(), nodename, metav1.GetOptions{})
	if err != nil {
		k.clientInfo.ErrorCheck(err, &countFlag)
		return
	}

	var status int = 0
	for _, condition := range node.Status.Conditions {
		if condition.Type == "Ready" {
			if condition.Status == "True" {
				status = 1
			}

			break
		}
	}

	if ns, ok := common.NodeStatusMap.LoadOrStore(nodename, status); ok {
		if ns.(int) != status {
			common.NodeStatusMap.Store(nodename, status)
			common.ChannelRequestChangeHost <- k.clientInfo.GetHost()
		}
	}

	if status == 1 {
		req := k.clientInfo.GetClientset().CoreV1().RESTClient().Get().Resource("nodes").Name(nodename).SubResource("proxy", "metrics", "cadvisor")
		res, err := req.Do(context.Background()).Raw()

		k.clientInfo.ErrorCheck(err, &countFlag)

		cadvisor_res_str := bytes.NewBuffer(res).String()
		cadvisor_reader := strings.NewReader(cadvisor_res_str)

		var parser expfmt.TextParser
		mf, err := parser.TextToMetricFamilies(cadvisor_reader)

		k.clientInfo.ErrorCheck(err, &countFlag)

		k.mutex.Lock()
		k.NodeMetric[nodename] = mf
		k.mutex.Unlock()

		writeLog(fmt.Sprintf("Kubernetes Get Node %s %s Metric Data end", k.clientInfo.GetVersions().GetVersion("core").Version, nodename))
	} else {
		writeLog(fmt.Sprintf("Kubernetes Get Node %s Metric Data is failed because the node's status is not ready", nodename))
	}
	defer wg.Done()
}
