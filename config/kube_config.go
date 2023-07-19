package config

type KubeConfig struct {
	kube_cluster_conffile string
	kube_cluster_userid   string
	kube_cluster_password string
	kube_cluster_sshport  uint16
	kube_cluster_rate     uint16
	//kube_enabledresources uint32 //= 현재 구성되어 있는 아래 항목 bit flag로 처리
	//[Namespace', 'Node', 'Pod', 'Service', 'Ingress', 'Deployment', 'StatefulSet', 'DaemonSet', 'ReplicaSet', 'Event', 'PersistentVolumeClaim', 'StorageClass', 'PersistentVolume' ]
}

// GetKubeClusterConfFile is a function that returns the kube cluster conf file.
func (kubeConf *KubeConfig) GetKubeClusterConfFile() string {
	return kubeConf.kube_cluster_conffile
}

// GetKubeClusterUserID is a function that returns the kube cluster user id.
func (kubeConf *KubeConfig) GetKubeClusterUserID() string {
	return kubeConf.kube_cluster_userid
}

// GetKubeClusterPassword is a function that returns the kube cluster password.
func (kubeConfg *KubeConfig) GetKubeClusterPassword() string {
	return kubeConfg.kube_cluster_password
}

// GetKubeClusterSSHPort is a function that returns the kube cluster ssh port.
func (kubeConf *KubeConfig) GetKubeClusterSSHPort() uint16 {
	return kubeConf.kube_cluster_sshport
}

// GetKubeClusterRate is a function that returns the kube cluster rate.
func (kubeConf *KubeConfig) GetKubeClusterRate() uint16 {
	return kubeConf.kube_cluster_rate
}

// func (kubeConf *KubeConfig) GetKubeEnabledResources() uint32 {
// 	return kubeConf.kube_enabledresources
// }
