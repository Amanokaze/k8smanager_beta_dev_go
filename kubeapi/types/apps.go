package types

import (
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
)

type AppsV1Data struct {
	Deployment  *appsv1.DeploymentList
	StatefulSet *appsv1.StatefulSetList
	DaemonSet   *appsv1.DaemonSetList
	ReplicaSet  *appsv1.ReplicaSetList
}

type AppsV1beta1Data struct {
	Deployment  *appsv1beta1.DeploymentList
	StatefulSet *appsv1beta1.StatefulSetList
}

type AppsV1beta2Data struct {
	Deployment  *appsv1beta2.DeploymentList
	StatefulSet *appsv1beta2.StatefulSetList
	DaemonSet   *appsv1beta2.DaemonSetList
	ReplicaSet  *appsv1beta2.ReplicaSetList
}

type AppsData struct {
}

type AppsDataInterface interface {
	AppsV1Data | AppsV1beta1Data | AppsV1beta2Data
}
