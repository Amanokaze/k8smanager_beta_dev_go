package kubeapi_test

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceGroupCount = 4
)

func TestGetCoreV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	corev1 := clientset.CoreV1()

	var err error
	_, err = corev1.Namespaces().List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	_, err = corev1.Nodes().List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	_, err = corev1.Pods("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	_, err = corev1.Services("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	_, err = corev1.PersistentVolumeClaims("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	_, err = corev1.PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)
}

func TestGetAppsV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	appsv1 := clientset.AppsV1()

	var err error
	_, err = appsv1.Deployments("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	_, err = appsv1.StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	_, err = appsv1.DaemonSets("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	_, err = appsv1.ReplicaSets("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)
}

func TestGetNetworkingV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	networkingv1 := clientset.NetworkingV1()

	_, err := networkingv1.Ingresses("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)
}

func TestGetStorageV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	storagev1 := clientset.StorageV1()

	_, err := storagev1.StorageClasses().List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)
}
