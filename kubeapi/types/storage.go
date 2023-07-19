package types

import (
	storagev1 "k8s.io/api/storage/v1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
)

type StorageV1Data struct {
	StorageClass *storagev1.StorageClassList
}

type StorageV1alpha1Data struct {
}

type StorageV1beta1Data struct {
	StorageClass *storagev1beta1.StorageClassList
}

type StorageData struct {
}

type StorageDataInterface interface {
	StorageV1Data
}
