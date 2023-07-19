package kubeapi

import (
	"fmt"
	"net"
	"net/url"
	"onTuneKubeManager/common"
	"onTuneKubeManager/config"
	"onTuneKubeManager/logger"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	certutil "k8s.io/client-go/util/cert"
)

type KubeClusters struct {
	Name string `yaml:"name"`
}

type KubeConfig struct {
	Clusters       []KubeClusters `yaml:"clusters"`
	CurrentContext string         `yaml:"current-context"`
}

// GetName returns the name of the cluster
func (k *KubeConfig) GetName() string {
	return k.Clusters[0].Name
}

// GetKubeClusterConfFile returns the name of the cluster
func (k *KubeConfig) GetContext() string {
	return k.CurrentContext
}

// GetKubeClusterConfFile returns the name of the cluster
// logMgr is used to write log
// config is used to get the name of the cluster
func GetClientsetConfig(logMgr *logger.OnTuneLogManager, config []config.KubeConfig) ([]*kubernetes.Clientset, []string, []*KubeConfig) {
	cs_arr := make([]*kubernetes.Clientset, 0)
	host_arr := make([]string, 0)
	cfg_arr := make([]*KubeConfig, 0)

	for i := 0; i < int(common.KubeClusterCount); i++ {
		fmt.Printf("%s\n", filepath.Join("config", config[i].GetKubeClusterConfFile()))
		var kubeconfig = filepath.Join("config", config[i].GetKubeClusterConfFile())

		// use the current context in kubeconfig
		restconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			logMgr.WriteLog(fmt.Sprintf("Kubernetes Clientset Configuration Error %v - %v", config[i].GetKubeClusterConfFile(), err.Error()))
			continue
		}
		host_url, _ := url.Parse(restconfig.Host)
		host, _, _ := net.SplitHostPort(host_url.Host)

		buf, err := os.ReadFile(kubeconfig)
		if err != nil {
			logMgr.WriteLog(fmt.Sprintf("Kubernetes Clientset Configuration File Read Error %v - %v", config[i].GetKubeClusterConfFile(), err.Error()))
			continue
		}

		k := &KubeConfig{}
		err = yaml.Unmarshal(buf, &k)
		if err != nil {
			logMgr.WriteLog(fmt.Sprintf("Kubernetes Clientset Configuration File Decoding Error %v - %v", config[i].GetKubeClusterConfFile(), err.Error()))
			continue
		}

		// create the clientset
		clientset, err := kubernetes.NewForConfig(restconfig)
		if err != nil {
			logMgr.WriteLog(fmt.Sprintf("Kubernetes Clientset Configuration Setting Error %v - %v", config[i].GetKubeClusterConfFile(), err.Error()))
			continue
		}

		cs_arr = append(cs_arr, clientset)
		host_arr = append(host_arr, host)
		cfg_arr = append(cfg_arr, k)
	}

	return cs_arr, host_arr, cfg_arr
}

// 해당 함수는 실행하지 않고 GetClientsetConfig() 함수만 실행하도록 하나,
// 차후 사용 가능성을 염두에 두고 아래와 같이 함수는 유지토록 함
func GetClientsetRest() *kubernetes.Clientset {
	token := "token값 직접 입력"
	host := "192.168.0.138"
	port := "6443"

	tlsClientConfig := rest.TLSClientConfig{}
	rootCAFile := filepath.Join("config", "ca.crt")

	if _, err := certutil.NewPool(rootCAFile); err != nil {
		fmt.Printf("Expected to load root CA config from %s, but got err: %v", rootCAFile, err)
	} else {
		tlsClientConfig.CAFile = rootCAFile
	}

	config := rest.Config{
		Host:            "https://" + net.JoinHostPort(host, port),
		BearerToken:     token,
		TLSClientConfig: tlsClientConfig,
	}

	// create the clientset
	var err error
	clientset, err := kubernetes.NewForConfig(&config)
	if err != nil {
		common.LogManager.WriteLog(fmt.Sprintf("Kubernetes Clientset Configuration Error - %v", err.Error()))
		panic(err.Error())
	}

	return clientset
}
