package kubeapi

import (
	"context"
	"fmt"
	"onTuneKubeManager/common"
	"onTuneKubeManager/kubeapi/types"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceGroupCount = 4
)

type KubernetesAPIResourceSource struct {
	clientInfo *ClientInfo
	Core       interface{}
	Apps       interface{}
	Networking interface{}
	Storage    interface{}
}

func (k *KubernetesAPIResourceSource) GetData(clientInfo *ClientInfo) {
	k.clientInfo = clientInfo

	var wg sync.WaitGroup
	wg.Add(ResourceGroupCount)

	go k.GetCoreData(&k.Core, &wg)
	go k.GetAppsData(&k.Apps, &wg)
	go k.GetNetworkingData(&k.Networking, &wg)
	go k.GetStorageData(&k.Storage, &wg)

	wg.Wait()
}

func (k *KubernetesAPIResourceSource) GetCoreData(coredata *interface{}, wg *sync.WaitGroup) {
	if k.clientInfo.GetVersions().GetVersion("core").Version == "v1" {
		common.LogManager.WriteLog(fmt.Sprintf("Kubernetes Get Core Data %s start", k.clientInfo.GetVersions().GetVersion("core").Version))
		corev1 := k.clientInfo.GetClientset().CoreV1()
		countFlag := true

		namespaces, err := corev1.Namespaces().List(context.TODO(), metav1.ListOptions{})
		k.clientInfo.ErrorDisabledCheck(err, &countFlag)
		if err != nil {
			return
		}

		nodes, err := corev1.Nodes().List(context.TODO(), metav1.ListOptions{})
		k.clientInfo.ErrorDisabledCheck(err, &countFlag)
		if err != nil {
			return
		}

		pods, err := corev1.Pods("").List(context.TODO(), metav1.ListOptions{})
		k.clientInfo.ErrorDisabledCheck(err, &countFlag)
		if err != nil {
			return
		}

		services, err := corev1.Services("").List(context.TODO(), metav1.ListOptions{})
		k.clientInfo.ErrorDisabledCheck(err, &countFlag)
		if err != nil {
			return
		}

		pvcs, err := corev1.PersistentVolumeClaims("").List(context.TODO(), metav1.ListOptions{})
		k.clientInfo.ErrorDisabledCheck(err, &countFlag)
		if err != nil {
			return
		}

		pvs, err := corev1.PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
		k.clientInfo.ErrorDisabledCheck(err, &countFlag)
		if err != nil {
			return
		}

		corev1data := &types.CoreV1Data{
			Namespace:             namespaces,
			Node:                  nodes,
			Pod:                   pods,
			Service:               services,
			PersistentVolumeClaim: pvcs,
			PersistentVolume:      pvs,
		}

		*coredata = corev1data
		common.LogManager.WriteLog(fmt.Sprintf("Kubernetes Get Core Data %s end", k.clientInfo.GetVersions().GetVersion("core").Version))
	}

	defer wg.Done()
}

func (k *KubernetesAPIResourceSource) GetAppsData(appsdata *interface{}, wg *sync.WaitGroup) {
	if k.clientInfo.GetVersions().GetVersion("apps").Version == "v1" {
		common.LogManager.WriteLog(fmt.Sprintf("Kubernetes Get Apps Data %s start", k.clientInfo.GetVersions().GetVersion("apps").Version))
		appsv1 := k.clientInfo.GetClientset().AppsV1()
		countFlag := true

		deploys, err := appsv1.Deployments("").List(context.TODO(), metav1.ListOptions{})
		k.clientInfo.ErrorDisabledCheck(err, &countFlag)
		if err != nil {
			return
		}

		stss, err := appsv1.StatefulSets("").List(context.TODO(), metav1.ListOptions{})
		k.clientInfo.ErrorDisabledCheck(err, &countFlag)
		if err != nil {
			return
		}

		dss, err := appsv1.DaemonSets("").List(context.TODO(), metav1.ListOptions{})
		k.clientInfo.ErrorDisabledCheck(err, &countFlag)
		if err != nil {
			return
		}

		rss, err := appsv1.ReplicaSets("").List(context.TODO(), metav1.ListOptions{})
		k.clientInfo.ErrorDisabledCheck(err, &countFlag)
		if err != nil {
			return
		}

		appsv1data := &types.AppsV1Data{
			Deployment:  deploys,
			StatefulSet: stss,
			DaemonSet:   dss,
			ReplicaSet:  rss,
		}

		*appsdata = appsv1data
		common.LogManager.WriteLog(fmt.Sprintf("Kubernetes Get Apps Data %s end", k.clientInfo.GetVersions().GetVersion("apps").Version))
	}

	defer wg.Done()
}

func (k *KubernetesAPIResourceSource) GetNetworkingData(networkingdata *interface{}, wg *sync.WaitGroup) {
	if k.clientInfo.GetVersions().GetVersion("networking").Version == "v1" {
		common.LogManager.WriteLog(fmt.Sprintf("Kubernetes Get Networking Data %s start", k.clientInfo.GetVersions().GetVersion("networking").Version))
		networkingv1 := k.clientInfo.GetClientset().NetworkingV1()
		countFlag := true

		ings, err := networkingv1.Ingresses("").List(context.TODO(), metav1.ListOptions{})
		k.clientInfo.ErrorDisabledCheck(err, &countFlag)
		if err != nil {
			return
		}

		networkingv1data := &types.NetworkingV1Data{
			Ingress: ings,
		}

		*networkingdata = networkingv1data
		common.LogManager.WriteLog(fmt.Sprintf("Kubernetes Get Networking Data %s end", k.clientInfo.GetVersions().GetVersion("networking").Version))
	}

	defer wg.Done()
}

func (k *KubernetesAPIResourceSource) GetStorageData(storagedata *interface{}, wg *sync.WaitGroup) {
	if k.clientInfo.GetVersions().GetVersion("storage").Version == "v1" {
		common.LogManager.WriteLog(fmt.Sprintf("Kubernetes Get Storage Data %s start", k.clientInfo.GetVersions().GetVersion("storage").Version))
		storagev1 := k.clientInfo.GetClientset().StorageV1()
		countFlag := true

		scs, err := storagev1.StorageClasses().List(context.TODO(), metav1.ListOptions{})
		k.clientInfo.ErrorDisabledCheck(err, &countFlag)
		if err != nil {
			return
		}

		storagev1data := &types.StorageV1Data{
			StorageClass: scs,
		}

		*storagedata = storagev1data
		common.LogManager.WriteLog(fmt.Sprintf("Kubernetes Get Storage Data %s end", k.clientInfo.GetVersions().GetVersion("storage").Version))
	}

	defer wg.Done()
}
