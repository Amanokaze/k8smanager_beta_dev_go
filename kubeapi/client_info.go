package kubeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"onTuneKubeManager/common"
	"onTuneKubeManager/config"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Referenced by APIResourceList.APIResouce
type ResourceSource struct {
	Kind     string
	Group    string
	Version  string
	Host     string
	Endpoint string
	Enabled  bool
}

type Resources struct {
	core       []ResourceSource
	apps       []ResourceSource
	networking []ResourceSource
	storage    []ResourceSource
}

type GroupVersions struct {
	core       metav1.GroupVersionForDiscovery
	apps       metav1.GroupVersionForDiscovery
	networking metav1.GroupVersionForDiscovery
	storage    metav1.GroupVersionForDiscovery
}

func (g *GroupVersions) GetVersion(group string) metav1.GroupVersionForDiscovery {
	switch group {
	case "core":
		return g.core
	case "apps":
		return g.apps
	case "networking":
		return g.networking
	case "storage":
		return g.storage
	default:
		return metav1.GroupVersionForDiscovery{}
	}
}

type ClientInfo struct {
	Clientset               *kubernetes.Clientset
	Host                    string
	Versions                *GroupVersions
	KubeCfg                 *config.KubeConfig
	Enabled                 bool
	LastRequestTime         time.Time
	RetryConnectionInterval uint16
	ResourceChangeInterval  uint16
	ResourceProcessFlag     bool
	SleepFlag               bool
	RequestFlag             bool
	RequestMutex            *sync.Mutex
	ResourceGroup           *Resources
}

// Set Group Version Info
func (c *ClientInfo) SetGroupVersions() error {
	countFlag := true

	groupVersions, err := c.getGroupVersions()
	if err != nil {
		common.LogManager.WriteLog(fmt.Sprintf("Kubernetes Version Error - %v", err.Error()))
		c.ErrorDisabledCheck(err, &countFlag)

		return err
	}

	c.Versions = groupVersions

	return nil
}

// Get Group Versions
func (c *ClientInfo) getGroupVersions() (*GroupVersions, error) {
	coreVersion, err := c.getGroupVersion("")
	if err != nil {
		return nil, err
	}

	appsVersion, err := c.getGroupVersion("apps")
	if err != nil {
		return nil, err
	}

	networkingVersion, err := c.getGroupVersion("networking.k8s.io")
	if err != nil {
		return nil, err
	}

	storageVersion, err := c.getGroupVersion("storage.k8s.io")
	if err != nil {
		return nil, err
	}

	groupVersions := GroupVersions{
		core:       coreVersion,
		apps:       appsVersion,
		networking: networkingVersion,
		storage:    storageVersion,
	}

	return &groupVersions, nil
}

// Get Single Group Version
func (c *ClientInfo) getGroupVersion(name string) (metav1.GroupVersionForDiscovery, error) {
	glist, err := c.Clientset.ServerGroups()
	if err != nil {
		return metav1.GroupVersionForDiscovery{}, err
	}

	for _, g := range glist.Groups {
		if g.Name == name {
			return g.PreferredVersion, nil
		}
	}

	return metav1.GroupVersionForDiscovery{}, nil
}

// Set Resource Info
func (c *ClientInfo) SetResourceGroups() {
	countFlag := true

	resources, err := c.getResourceGroups()
	if err != nil {
		common.LogManager.WriteLog(fmt.Sprintf("Kubernetes Version Error - %v", err.Error()))
		c.ErrorDisabledCheck(err, &countFlag)
	}

	c.ResourceGroup = resources
}

func (c *ClientInfo) getResourceGroups() (*Resources, error) {
	coreResources, err := c.getResources("core")
	if err != nil {
		return nil, err
	}

	appsResources, err := c.getResources("apps")
	if err != nil {
		return nil, err
	}

	networkingResources, err := c.getResources("networking")
	if err != nil {
		return nil, err
	}

	storageResources, err := c.getResources("storage")
	if err != nil {
		return nil, err
	}

	resources := Resources{
		core:       coreResources,
		apps:       appsResources,
		networking: networkingResources,
		storage:    storageResources,
	}

	return &resources, nil
}

func (c *ClientInfo) getResources(name string) ([]ResourceSource, error) {
	var req rest.Request
	var group string
	var epPrefix string

	_, err := c.getGroupVersions()
	if err != nil {
		return nil, err
	}

	switch name {
	case "core":
		if c.Versions.GetVersion(name).Version == "v1" {
			req = *c.Clientset.CoreV1().RESTClient().Get()
			group = "CoreV1"
			epPrefix = ENDPOINT_PREFIX_API
		}
	case "apps":
		if c.Versions.GetVersion(name).Version == "v1" {
			req = *c.Clientset.AppsV1().RESTClient().Get()
			group = "AppsV1"
			epPrefix = ENDPOINT_PREFIX_APIS
		}
	case "networking":
		if c.Versions.GetVersion(name).Version == "v1" {
			req = *c.Clientset.NetworkingV1().RESTClient().Get()
			group = "NetworkingV1"
			epPrefix = ENDPOINT_PREFIX_APIS
		}
	case "storage":
		if c.Versions.GetVersion(name).Version == "v1" {
			req = *c.Clientset.StorageV1().RESTClient().Get()
			group = "StorageV1"
			epPrefix = ENDPOINT_PREFIX_APIS
		}
	}
	res, err := req.Do(context.Background()).Raw()
	if err != nil {
		return nil, err
	}

	var rsc_data *metav1.APIResourceList
	err = json.Unmarshal(res, &rsc_data)
	if err != nil {
		return nil, err
	}

	rsc_sources := make([]ResourceSource, 0)
	rsc_map := make(map[string]ResourceSource)
	for _, r := range rsc_data.APIResources {
		if _, ok := rsc_map[r.Kind]; ok {
			continue
		} else {
			host, _, _ := net.SplitHostPort(c.Clientset.DiscoveryClient.RESTClient().Get().URL().Host)
			rsc_map[r.Kind] = ResourceSource{
				Kind:     r.Kind,
				Group:    group,
				Version:  c.Versions.GetVersion(name).Version,
				Host:     host,
				Endpoint: epPrefix + c.Versions.GetVersion(name).GroupVersion,
				Enabled:  Contains(r.Kind),
			}
			rsc_sources = append(rsc_sources, rsc_map[r.Kind])
		}
	}
	rcs_json, _ := json.Marshal(rsc_map)
	rcs_resource_map := make(map[string][]byte)
	rcs_resource_map["resourceinfo"] = rcs_json
	common.ChannelResourceData <- rcs_resource_map

	return rsc_sources, nil
}

func (c *ClientInfo) GetInitData() {
	c.GetResourceData()
	c.GetEventlogData(true)
	c.GetMetricData()
	c.SetResourceGroups()
}

func (c *ClientInfo) GetResourceDataInterval() {
	var last_updatetime int64 = time.Now().Unix()
	for {
		c.SleepGetData()

		var interval int
		if c.RequestFlag {
			interval = common.RealtimeInterval
		} else {
			interval = common.ResourceInterval
		}

		if last_updatetime+int64(interval) <= time.Now().Unix() {
			last_updatetime = time.Now().Unix()
			c.GetResourceData()
		}
		time.Sleep(time.Second * 1)
	}
}

func (c *ClientInfo) GetEventlogDataInterval() {
	var last_updatetime int64 = time.Now().Unix()
	for {
		c.SleepGetData()

		if last_updatetime+int64(common.EventlogInterval) <= time.Now().Unix() {
			last_updatetime = time.Now().Unix()
			c.GetEventlogData(false)
		}
		time.Sleep(time.Second * 1)
	}
}

func (c *ClientInfo) GetMetricDataInterval() {
	var last_updatetime int64 = time.Now().Unix()
	for {
		c.SleepGetData()

		if last_updatetime+int64(common.RealtimeInterval) <= time.Now().Unix() {
			last_updatetime = time.Now().Unix()
			c.GetMetricData()
		}
		time.Sleep(time.Second * 1)
	}
}

func (c *ClientInfo) GetResourceData() {
	err := c.SetGroupVersions()
	if err != nil {
		return
	}

	c.ResourceProcessFlag = true
	var k8sApiResourceSource = KubernetesAPIResourceSource{}
	k8sApiResourceSource.GetData(c)

	mappingResourceSource := MappingResourceSource{
		ClientInfo: *c,
	}

	mappingResourceSource.MakeResourceData(&k8sApiResourceSource)
	c.ResourceProcessFlag = false
}

func (c *ClientInfo) GetEventlogData(init_flag bool) {
	err := c.SetGroupVersions()
	if err != nil {
		return
	}

	var k8sApiEventlogSource = KubernetesAPIEventlogSource{}
	k8sApiEventlogSource.GetData(c, init_flag)

	mappingEventlogSource := MappingEventlogSource{
		ClientInfo: *c,
	}

	mappingEventlogSource.MakeEventlogData(&k8sApiEventlogSource)
}

func (c *ClientInfo) GetMetricData() {
	err := c.SetGroupVersions()
	if err != nil {
		return
	}

	var k8sApiMetricSource = KubernetesAPIMetricSource{}
	k8sApiMetricSource.GetData(c)

	mappingMetricSource := MappingMetricSource{
		ClientInfo: *c,
	}

	mappingMetricSource.MakeMetricData(&k8sApiMetricSource)
}

func (c *ClientInfo) ErrorCheck(err error, flag *bool) {
	if err != nil {
		errMsg := err.Error()
		common.LogManager.WriteLog(fmt.Sprintf("Kubernetes API Error - %s", errMsg))
	} else if !c.Enabled {
		c.Enabled = true
		cluster_map := make(map[string]bool)
		cluster_map[c.Host] = true
		common.ChannelClusterStatus <- cluster_map
	}
}

func (c *ClientInfo) ErrorDisabledCheck(err error, flag *bool) {
	if err != nil {
		if !c.SleepFlag {
			errMsg := err.Error()
			common.LogManager.WriteLog(fmt.Sprintf("Kubernete API Error(%d) - %s", common.TerminatedCount, errMsg))
			//log.Panic(errMsg)
			if *flag {
				common.TerminatedCount += 1
				*flag = false
			}

			if c.Enabled {
				c.Enabled = false
				cluster_map := make(map[string]bool)
				cluster_map[c.Host] = false
				common.ChannelClusterStatus <- cluster_map
			}
		}
	} else if !c.Enabled {
		c.Enabled = true
		cluster_map := make(map[string]bool)
		cluster_map[c.Host] = true
		common.ChannelClusterStatus <- cluster_map
	}

	if common.TerminatedCount > int(common.DisconnectFailCount) {
		common.LogManager.WriteLog(fmt.Sprintf("Error Terminated Count is over %d, Retry after %d seconds.", common.TerminatedCount, c.RetryConnectionInterval))
		common.TerminatedCount = 0
		c.SleepFlag = true
	}
}

func (c *ClientInfo) RecoverResourceInterval() {
	for {
		if c.RequestFlag && time.Since(c.LastRequestTime) > time.Duration(int(c.ResourceChangeInterval)*int(time.Second)) {
			c.RequestMutex.Lock()
			c.RequestFlag = false
			common.LogManager.WriteLog(fmt.Sprintf("ChangeResourceInterval: -> resourceinterval %v %v", time.Since(c.LastRequestTime), time.Duration(int(c.ResourceChangeInterval)*int(time.Second))))
			c.RequestMutex.Unlock()
		}
		time.Sleep(1 * time.Second)
	}
}

func (c *ClientInfo) ChangeResourceInterval() {
	c.RequestMutex.Lock()
	c.LastRequestTime = time.Now()
	c.RequestFlag = true
	common.LogManager.WriteLog("ChangeResourceInterval: -> RealtimeInterval")
	c.RequestMutex.Unlock()
}

func (c *ClientInfo) SleepGetData() {
	if c.SleepFlag {
		time.Sleep(time.Second * time.Duration(c.RetryConnectionInterval))
		c.SleepFlag = false
		common.LogManager.WriteLog("Sleep time is over. Get API Data will restart.")
	}
}

func (c *ClientInfo) GetClientset() *kubernetes.Clientset {
	return c.Clientset
}

func (c *ClientInfo) GetHost() string {
	return c.Host
}

func (c *ClientInfo) GetKubeCfg() *config.KubeConfig {
	return c.KubeCfg
}

func (c *ClientInfo) GetResourceProcessFlag() bool {
	return c.ResourceProcessFlag
}

func (c *ClientInfo) GetEnabled() bool {
	return c.Enabled
}

func (c *ClientInfo) GetVersions() *GroupVersions {
	return c.Versions
}
