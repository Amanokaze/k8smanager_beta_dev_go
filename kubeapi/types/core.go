package types

import (
	corev1 "k8s.io/api/core/v1"
)

type CoreV1Data struct {
	Namespace             *corev1.NamespaceList
	Node                  *corev1.NodeList
	Pod                   *corev1.PodList
	Service               *corev1.ServiceList
	PersistentVolumeClaim *corev1.PersistentVolumeClaimList
	PersistentVolume      *corev1.PersistentVolumeList
}

type CoreDataInterface interface {
	CoreV1Data
}

type CoreV1EventlogData struct {
	Event *corev1.EventList
}
