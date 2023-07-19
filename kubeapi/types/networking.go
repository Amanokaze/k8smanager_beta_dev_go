package types

import (
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
)

type NetworkingV1Data struct {
	Ingress *networkingv1.IngressList
}

type NetworkingV1alpha1Data struct {
}

type NetworkingV1beta1Data struct {
	Ingress *networkingv1beta1.IngressList
}

type NetworkingData struct {
}

type NetworkingDataInterface interface {
	NetworkingV1Data
}
