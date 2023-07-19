package kubeapi_test

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/prometheus/common/expfmt"
	cv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetMetricNodeData(t *testing.T) {
	clientset := authConfigfile(t)
	corev1 := clientset.CoreV1()
	nodes, err := corev1.Nodes().List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	for _, n := range nodes.Items {
		var status int = 0
		for _, condition := range n.Status.Conditions {
			if condition.Type == "Ready" {
				if condition.Status == cv1.ConditionTrue {
					status = 1
				}

				break
			}
		}

		if status == 1 {
			req := corev1.RESTClient().Get().Resource("nodes").Name(n.Name).SubResource("proxy", "metrics", "cadvisor")
			res, err := req.Do(context.Background()).Raw()
			errorCheck(err, t)

			cadvisor_res_str := bytes.NewBuffer(res).String()
			cadvisor_reader := strings.NewReader(cadvisor_res_str)

			var parser expfmt.TextParser
			mf, err := parser.TextToMetricFamilies(cadvisor_reader)
			errorCheck(err, t)

			fmt.Printf("MF Data Length: %d\n", len(mf))
		}
	}
}
