package kubeapi

import (
	"context"
	"fmt"
	"onTuneKubeManager/common"
	"onTuneKubeManager/kubeapi/types"
	"sync"

	"golang.org/x/crypto/ssh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	EventLogCount = 1
)

// Create SSH client struct
type SSHClient struct {
	Host string
	Port int
	User string
	Pass string
	*ssh.Client
}

// Create SSH client
func NewSSHClient(host string, port int, user string, pass string) *SSHClient {
	return &SSHClient{Host: host, Port: port, User: user, Pass: pass}
}

// Connect to SSH server
func (c *SSHClient) Connect() error {
	conf := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port), conf)
	if err != nil {
		return err
	}
	c.Client = conn

	return nil
}

// Close SSH connection
func (c *SSHClient) Close() error {
	if c.Client != nil {
		return c.Client.Close()
	}

	return nil
}

// Execute command on SSH server
func (c *SSHClient) Execute(cmd string) (string, error) {
	// Execute command on SSH server
	session, err := c.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.Output(cmd)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

type KubernetesAPIEventlogSource struct {
	clientInfo *ClientInfo
	SSHClient
	Event interface{}
	Logs  interface{}
}

func (k *KubernetesAPIEventlogSource) GetData(clientInfo *ClientInfo, init_flag bool) {
	k.clientInfo = clientInfo
	k.SSHClient = *NewSSHClient(clientInfo.GetHost(), int(clientInfo.GetKubeCfg().GetKubeClusterSSHPort()), clientInfo.GetKubeCfg().GetKubeClusterUserID(), clientInfo.GetKubeCfg().GetKubeClusterPassword())

	var wg sync.WaitGroup
	wg.Add(EventLogCount)
	go k.GetCoreData(&k.Event, &wg)
	//go k.GetLogData(&k.Logs, init_flag, &wg) 보류
	wg.Wait()
}

func (k *KubernetesAPIEventlogSource) GetCoreData(coreventdata *interface{}, wg *sync.WaitGroup) {
	if k.clientInfo.GetVersions().GetVersion("core").Version == "v1" {
		common.LogManager.WriteLog(fmt.Sprintf("Kubernetes Get Event Data %s start", k.clientInfo.GetVersions().GetVersion("core").Version))
		corev1 := k.clientInfo.GetClientset().CoreV1()
		countFlag := true

		events, err := corev1.Events("").List(context.TODO(), metav1.ListOptions{})
		k.clientInfo.ErrorDisabledCheck(err, &countFlag)
		if err != nil {
			return
		}

		corev1eventdata := &types.CoreV1EventlogData{
			Event: events,
		}

		*coreventdata = corev1eventdata
		//common.LogManager.WriteLog(fmt.Sprintf("%v\n", corev1eventdata.Event.Items))
		common.LogManager.WriteLog(fmt.Sprintf("Kubernetes Get Event Data %s end", k.clientInfo.GetVersions().GetVersion("core").Version))
	}

	defer wg.Done()
}

func (k *KubernetesAPIEventlogSource) GetLogData(logdata *interface{}, init_flag bool, wg *sync.WaitGroup) {
	err := k.SSHClient.Connect()
	if err != nil {
		common.LogManager.WriteLog(fmt.Sprintf("SSH Connect Error %s", err))
	}

	// Get log data
	common.LogManager.WriteLog("Kubernetes Get Log Data start")
	if pod_name_array, ok := common.PodListMap.Load(k.clientInfo.GetHost()); ok {
		log_map := &sync.Map{}
		pod_values := pod_name_array.([]string)
		for _, pod := range pod_values {
			var get_statement string
			namespace_name, pod_name := SplitPodName(pod)

			if init_flag {
				get_statement = fmt.Sprintf("kubectl logs %s -n %s --tail=-1 --timestamps=true --all-containers=true", pod_name, namespace_name)
			} else {
				since_interval := common.EventlogInterval + 1
				get_statement = fmt.Sprintf("kubectl logs %s -n %s --since=%ds --timestamps=true --all-containers=true", pod_name, namespace_name, since_interval)
			}

			log, err := k.SSHClient.Execute(get_statement)
			if err != nil {
				common.LogManager.WriteLog(fmt.Sprintf("SSH Execute Error %s - %s", err, get_statement))
			}

			if log != "" {
				log_map.Store(pod, log)
			}
		}
		*logdata = log_map
	}

	err = k.SSHClient.Close()
	if err != nil {
		common.LogManager.WriteLog(fmt.Sprintf("SSH Close Error %s", err))
	}

	defer wg.Done()
}
