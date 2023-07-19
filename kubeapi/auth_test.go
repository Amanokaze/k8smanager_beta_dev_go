package kubeapi_test

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func errorCheck(err error, t *testing.T) {
	if err != nil {
		t.Error(err.Error())
	}
}

func authConfigfile(t *testing.T) *kubernetes.Clientset {
	kubeconfig := filepath.Join("..", "config", ".kube.2.conf")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	errorCheck(err, t)

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	errorCheck(err, t)

	return clientset
}

func getPodinfo(t *testing.T, clientset *kubernetes.Clientset) {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// Examples for error handling:
	// - Use helper functions like e.g. errors.IsNotFound()
	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	namespace := "kube-system"
	pod := "kube-apiserver-dev-masternode"

	var statuserror *k8serrors.StatusError
	_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
	if k8serrors.IsNotFound(err) {
		t.Errorf("Pod %s in namespace %s not found\n", pod, namespace)
	} else if errors.As(err, &statuserror) {
		t.Errorf("Error getting pod %s in namespace %s: %v\n", pod, namespace, statuserror.ErrStatus.Message)
	} else if err != nil {
		t.Errorf(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
	}
}

func TestAuthConfigfile(t *testing.T) {
	authConfigfile(t)
}

func TestBasicPodinfo(t *testing.T) {
	clientset := authConfigfile(t)
	getPodinfo(t, clientset)
}
