package kubeapi_test

import (
	"context"
	"fmt"
	"onTuneKubeManager/kubeapi"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMakeNamespaceV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	corev1client := clientset.CoreV1()

	nss, err := corev1client.Namespaces().List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	namespaces := make([]kubeapi.MappingNamespace, 0)
	for _, n := range nss.Items {
		namespaces = append(namespaces, kubeapi.MappingNamespace{
			Name:   n.GetObjectMeta().GetName(),
			UID:    string(n.GetObjectMeta().GetUID()),
			Labels: kubeapi.GetLabelSelector(n.GetObjectMeta().GetLabels()),
			Status: string(n.Status.Phase),
		})
	}
	fmt.Printf("Namespace count is %d\n", len(namespaces))
}

func TestMakeNodeV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	corev1client := clientset.CoreV1()

	ns, err := corev1client.Nodes().List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	nodes := make([]kubeapi.MappingNode, 0)
	for _, n := range ns.Items {
		var nodetype string
		if _, ok := n.GetObjectMeta().GetLabels()["node-role.kubernetes.io/control-plane"]; ok {
			nodetype = "Control Plane"
		} else {
			nodetype = "Worker"
		}
		nodes = append(nodes, kubeapi.MappingNode{
			Name:                    n.GetObjectMeta().GetName(),
			UID:                     string(n.GetObjectMeta().GetUID()),
			StartTime:               n.GetObjectMeta().GetCreationTimestamp().Time,
			Labels:                  kubeapi.GetLabelSelector(n.GetObjectMeta().GetLabels()),
			NodeType:                nodetype,
			KernelVersion:           n.Status.NodeInfo.KernelVersion,
			OSImage:                 n.Status.NodeInfo.OSImage,
			OSName:                  n.Status.NodeInfo.OperatingSystem,
			ContainerRuntimeVersion: n.Status.NodeInfo.ContainerRuntimeVersion,
			KubeletVersion:          n.Status.NodeInfo.KubeletVersion,
			KubeProxyVersion:        n.Status.NodeInfo.KubeProxyVersion,
			CPUArch:                 n.Status.NodeInfo.Architecture,
			CPUCount:                n.Status.Capacity.Cpu().Value(),
			EphemeralStorage:        n.Status.Capacity.StorageEphemeral().Value(),
			MemorySize:              n.Status.Capacity.Memory().Value(),
			Pods:                    n.Status.Capacity.Pods().Value(),
			IP:                      n.Status.Addresses[0].Address,
		})
	}
	fmt.Printf("Node count is %d\n", len(nodes))
}

func TestMakePodV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	corev1client := clientset.CoreV1()

	ns, err := corev1client.Nodes().List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	var masternodename string
	for _, n := range ns.Items {
		if _, ok := n.GetObjectMeta().GetLabels()["node-role.kubernetes.io/control-plane"]; ok {
			masternodename = n.GetObjectMeta().GetName()
		}
	}

	ps, err := corev1client.Pods("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	pods := make([]kubeapi.MappingPod, 0)
	for _, n := range ps.Items {
		var nodename = n.Spec.NodeName
		var podname = n.GetObjectMeta().GetName()
		var poduid = string(n.GetObjectMeta().GetUID())

		var annotation string
		if _, ok := n.GetObjectMeta().GetAnnotations()["kubernetes.io/config.hash"]; ok {
			annotation = n.GetObjectMeta().GetAnnotations()["kubernetes.io/config.hash"]
		} else if _, ok := n.GetObjectMeta().GetAnnotations()["kubernetes.io/config.mirror"]; ok {
			annotation = n.GetObjectMeta().GetAnnotations()["kubernetes.io/config.mirror"]
		}

		var restartcount int32
		var restarttime time.Time
		for _, c := range n.Status.ContainerStatuses {
			restartcount += c.RestartCount
			if c.State.Running != nil && (restarttime.IsZero() || c.State.Running.StartedAt.Time.Unix() > restarttime.Unix()) {
				restarttime = c.State.Running.StartedAt.Time
			}
		}

		podconditions := make([]string, 0)
		for _, c := range n.Status.Conditions {
			if c.Message != "" || c.Status == corev1.ConditionFalse {
				podcondition := fmt.Sprintf("%s(%s/%s):%s", c.Type, c.Status, c.LastTransitionTime.Format(time.RFC3339), c.Message)
				podconditions = append(podconditions, podcondition)
			}
		}

		var staticpodname string
		if nodename == masternodename && len(podname) > len(nodename) && podname[len(podname)-len(nodename):] == nodename {
			staticpodname = podname[:len(podname)-len(nodename)-1]
		}

		container_status := make(map[string]kubeapi.ContainerStatus)
		for _, cs := range n.Status.ContainerStatuses {
			var container_state string
			if cs.State.Running != nil {
				container_state = "running"
			} else if cs.State.Waiting != nil {
				container_state = "waiting"
			} else if cs.State.Terminated != nil {
				container_state = "terminated"
			}

			container_status[cs.Name] = kubeapi.ContainerStatus{
				State:        container_state,
				RestartCount: cs.RestartCount,
				Ready:        cs.Ready,
				Started:      *cs.Started,
			}
		}

		containers := make([]kubeapi.MappingContainer, 0)
		for _, c := range n.Spec.Containers {
			cports := make([]string, 0)
			for _, cp := range c.Ports {
				var cpn string
				if cp.Name != "" {
					cpn = cp.Name
				}

				var cphp string
				if cp.HostPort > 0 {
					cphp = fmt.Sprintf("-%d", cp.HostPort)
				}

				var cpflag string
				if cpn != "" || cphp != "" {
					cpflag = "/"
				}

				var cpp string
				if cp.Protocol == "" {
					cpp = "TCP"
				} else {
					cpp = string(cp.Protocol)
				}

				var cphi string
				if cp.HostIP != "" {
					cphi = fmt.Sprintf(":%s", cp.HostIP)
				}

				cports = append(cports, fmt.Sprintf("%s%s%s%s%s", cpn, cphp, cpflag, cpp, cphi))
			}

			cenvs := make([]string, 0)
			for _, ce := range c.Env {
				cenvs = append(cenvs, fmt.Sprintf("%s:%s", ce.Name, ce.Value))
			}

			cvolumemounts := make([]string, 0)
			for _, cvm := range c.VolumeMounts {
				var cvmro string
				if cvm.ReadOnly {
					cvmro = "(RO)"
				}

				// SubPathExpr and SubPath are mutually exclusive.
				var cvmsp string
				if cvm.SubPath != "" {
					cvmsp = fmt.Sprintf(":%s", cvm.SubPath)
				} else if cvm.SubPathExpr != "" {
					cvmsp = fmt.Sprintf(":%s", cvm.SubPathExpr)
				}
				cvolumemounts = append(cvolumemounts, fmt.Sprintf("%s%s(%s%s)", cvm.Name, cvmro, cvm.MountPath, cvmsp))
			}
			containers = append(containers, kubeapi.MappingContainer{
				Name:             c.Name,
				UID:              poduid,
				Image:            c.Image,
				Ports:            strings.Join(cports, ","),
				Env:              strings.Join(cenvs, ","),
				LimitCpu:         c.Resources.Limits.Cpu().Value(),
				LimitMemory:      c.Resources.Limits.Memory().Value(),
				LimitStorage:     c.Resources.Limits.Storage().Value(),
				LimitEphemeral:   c.Resources.Limits.StorageEphemeral().Value(),
				RequestCpu:       c.Resources.Requests.Cpu().Value(),
				RequestMemory:    c.Resources.Requests.Memory().Value(),
				RequestStorage:   c.Resources.Requests.Storage().Value(),
				RequestEphemeral: c.Resources.Requests.StorageEphemeral().Value(),
				VolumeMounts:     strings.Join(cvolumemounts, ","),
			})
		}

		var pvc string
		for _, v := range n.Spec.Volumes {
			if v.PersistentVolumeClaim != nil {
				pvc = v.PersistentVolumeClaim.ClaimName
			}
		}

		refs := kubeapi.GetRefData(n.GetObjectMeta().GetOwnerReferences())
		pods = append(pods, kubeapi.MappingPod{
			Name:                  podname,
			UID:                   poduid,
			StartTime:             n.GetObjectMeta().GetCreationTimestamp().Time,
			Labels:                kubeapi.GetLabelSelector(n.GetObjectMeta().GetLabels()),
			Selector:              kubeapi.GetLabelSelector(n.Spec.NodeSelector),
			AnnotationUID:         annotation,
			RestartPolicy:         string(n.Spec.RestartPolicy),
			ServiceAccount:        n.Spec.ServiceAccountName,
			Status:                string(n.Status.Phase),
			HostIP:                n.Status.HostIP,
			PodIP:                 n.Status.PodIP,
			RestartCount:          restartcount,
			RestartTime:           restarttime,
			Condition:             strings.Join(podconditions, ","),
			StaticPod:             staticpodname,
			NodeName:              nodename,
			NamespaceName:         n.GetObjectMeta().GetNamespace(),
			ReferenceKind:         refs["kind"],
			ReferenceUID:          refs["UID"],
			PersistentVolumeClaim: pvc,
			Containers:            containers,
		})
	}
	fmt.Printf("Pod count is %d\n", len(pods))
}

func TestMakeServiceV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	corev1client := clientset.CoreV1()

	svcs, err := corev1client.Services("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	services := make([]kubeapi.MappingService, 0)
	for _, n := range svcs.Items {
		ports := make([]string, 0)
		for _, p := range n.Spec.Ports {
			var pt string
			if p.TargetPort.StrVal == "" || p.TargetPort.IntVal == p.Port {
				pt = fmt.Sprintf("%d", p.Port)
			} else {
				pt = fmt.Sprintf("%d:%s", p.Port, p.TargetPort.StrVal)
			}

			var npt string
			if p.NodePort != 0 {
				npt = fmt.Sprintf(":%d", p.NodePort)
			}
			port := fmt.Sprintf("%s(%s%s/%s)", p.Name, pt, npt, p.Protocol)
			ports = append(ports, port)
		}

		services = append(services, kubeapi.MappingService{
			Name:          n.GetObjectMeta().GetName(),
			UID:           string(n.GetObjectMeta().GetUID()),
			StartTime:     n.GetObjectMeta().GetCreationTimestamp().Time,
			Labels:        kubeapi.GetLabelSelector(n.GetObjectMeta().GetLabels()),
			Selector:      kubeapi.GetLabelSelector(n.Spec.Selector),
			ServiceType:   string(n.Spec.Type),
			ClusterIP:     n.Spec.ClusterIP,
			Ports:         strings.Join(ports, ","),
			NamespaceName: n.GetObjectMeta().GetNamespace(),
		})
	}
	fmt.Printf("Service count is %d\n", len(services))
}

func TestMakePersistentVolumeClaimV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	corev1client := clientset.CoreV1()

	pvcs, err := corev1client.PersistentVolumeClaims("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	persistentvolumeclaims := make([]kubeapi.MappingPvc, 0)
	for _, n := range pvcs.Items {
		persistentvolumeclaims = append(persistentvolumeclaims, kubeapi.MappingPvc{
			Name:             n.GetObjectMeta().GetName(),
			UID:              string(n.GetObjectMeta().GetUID()),
			StartTime:        n.GetObjectMeta().GetCreationTimestamp().Time,
			Labels:           kubeapi.GetLabelSelector(n.GetObjectMeta().GetLabels()),
			Selector:         kubeapi.GetSelector(n.Spec.Selector),
			RequestStorage:   n.Spec.Resources.Requests.Storage().Value(),
			AccessModes:      kubeapi.GetAccessModes(n.Spec.AccessModes),
			Status:           string(n.Status.Phase),
			StorageClassName: *n.Spec.StorageClassName,
			NamespaceName:    n.GetObjectMeta().GetNamespace(),
		})
	}
	fmt.Printf("PVC count is %d\n", len(persistentvolumeclaims))
}

func TestMakePersistentVolumeV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	corev1client := clientset.CoreV1()

	pvs, err := corev1client.PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	persistentvolumes := make([]kubeapi.MappingPv, 0)
	for _, n := range pvs.Items {
		persistentvolumes = append(persistentvolumes, kubeapi.MappingPv{
			Name:          n.GetObjectMeta().GetName(),
			UID:           string(n.GetObjectMeta().GetUID()),
			StartTime:     n.GetObjectMeta().GetCreationTimestamp().Time,
			Labels:        kubeapi.GetLabelSelector(n.GetObjectMeta().GetLabels()),
			AccessModes:   kubeapi.GetAccessModes(n.Spec.AccessModes),
			Capacity:      n.Spec.Capacity.Storage().Value(),
			ReclaimPolicy: string(n.Spec.PersistentVolumeReclaimPolicy),
			Status:        string(n.Status.Phase),
			PvcUID:        string(n.Spec.ClaimRef.UID),
		})
	}
	fmt.Printf("PV count is %d\n", len(persistentvolumes))
}

func TestMakeDeploymentV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	appsv1client := clientset.AppsV1()

	deploys, err := appsv1client.Deployments("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	deployments := make([]kubeapi.MappingDeployment, 0)
	for _, n := range deploys.Items {
		deployments = append(deployments, kubeapi.MappingDeployment{
			Name:                 n.GetObjectMeta().GetName(),
			UID:                  string(n.GetObjectMeta().GetUID()),
			StartTime:            n.GetObjectMeta().GetCreationTimestamp().Time,
			Labels:               kubeapi.GetLabelSelector(n.GetObjectMeta().GetLabels()),
			Selector:             kubeapi.GetSelector(n.Spec.Selector),
			ServiceAccount:       n.Spec.Template.Spec.ServiceAccountName,
			Replicas:             *n.Spec.Replicas,
			UpdatedReplicas:      n.Status.UpdatedReplicas,
			ReadyReplicas:        n.Status.ReadyReplicas,
			AvailableReplicas:    n.Status.AvailableReplicas,
			ObservedGeneneration: n.Status.ObservedGeneration,
			NamespaceName:        n.GetObjectMeta().GetNamespace(),
		})
	}
	fmt.Printf("Deploy count is %d\n", len(deployments))
}

func TestMakeStatefulSetV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	appsv1client := clientset.AppsV1()

	stss, err := appsv1client.StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	stateulsets := make([]kubeapi.MappingStatefulSet, 0)
	for _, n := range stss.Items {
		stateulsets = append(stateulsets, kubeapi.MappingStatefulSet{
			Name:              n.GetObjectMeta().GetName(),
			UID:               string(n.GetObjectMeta().GetUID()),
			StartTime:         n.GetObjectMeta().GetCreationTimestamp().Time,
			Labels:            kubeapi.GetLabelSelector(n.GetObjectMeta().GetLabels()),
			Selector:          kubeapi.GetSelector(n.Spec.Selector),
			ServiceAccount:    n.Spec.Template.Spec.ServiceAccountName,
			Replicas:          *n.Spec.Replicas,
			UpdatedReplicas:   n.Status.UpdatedReplicas,
			ReadyReplicas:     n.Status.ReadyReplicas,
			AvailableReplicas: n.Status.AvailableReplicas,
			NamespaceName:     n.GetObjectMeta().GetNamespace(),
		})
	}
	fmt.Printf("STS count is %d\n", len(stateulsets))
}

func TestMakeDaemonSetV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	appsv1client := clientset.AppsV1()

	dss, err := appsv1client.DaemonSets("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	daemonsets := make([]kubeapi.MappingDaemonSet, 0)
	for _, n := range dss.Items {
		daemonsets = append(daemonsets, kubeapi.MappingDaemonSet{
			Name:                   n.GetObjectMeta().GetName(),
			UID:                    string(n.GetObjectMeta().GetUID()),
			StartTime:              n.GetObjectMeta().GetCreationTimestamp().Time,
			Labels:                 kubeapi.GetLabelSelector(n.GetObjectMeta().GetLabels()),
			Selector:               kubeapi.GetSelector(n.Spec.Selector),
			ServiceAccount:         n.Spec.Template.Spec.ServiceAccountName,
			CurrentNumberScheduled: n.Status.CurrentNumberScheduled,
			DesiredNumberScheduled: n.Status.DesiredNumberScheduled,
			NumberReady:            n.Status.NumberReady,
			UpdatedNumberScheduled: n.Status.UpdatedNumberScheduled,
			NumberAvailable:        n.Status.NumberAvailable,
			NamespaceName:          n.GetObjectMeta().GetNamespace(),
		})
	}
	fmt.Printf("Daemonset count is %d\n", len(daemonsets))
}

func TestMakeReplicaSetV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	appsv1client := clientset.AppsV1()

	rss, err := appsv1client.ReplicaSets("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	replicasets := make([]kubeapi.MappingReplicaSet, 0)
	for _, n := range rss.Items {
		refs := kubeapi.GetRefData(n.GetObjectMeta().GetOwnerReferences())
		replicasets = append(replicasets, kubeapi.MappingReplicaSet{
			Name:                 n.GetObjectMeta().GetName(),
			UID:                  string(n.GetObjectMeta().GetUID()),
			StartTime:            n.GetObjectMeta().GetCreationTimestamp().Time,
			Labels:               kubeapi.GetLabelSelector(n.GetObjectMeta().GetLabels()),
			Selector:             kubeapi.GetSelector(n.Spec.Selector),
			Replicas:             *n.Spec.Replicas,
			FullyLabeledReplicas: n.Status.FullyLabeledReplicas,
			ReadyReplicas:        n.Status.ReadyReplicas,
			AvailableReplicas:    n.Status.AvailableReplicas,
			ReferenceKind:        refs["kind"],
			ReferenceUID:         refs["UID"],
			NamespaceName:        n.GetObjectMeta().GetNamespace(),
		})
	}
	fmt.Printf("Replilcaset count is %d\n", len(replicasets))
}

func TestMakeIngressV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	networkingv1client := clientset.NetworkingV1()

	ings, err := networkingv1client.Ingresses("").List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	ingresses := make([]kubeapi.MappingIngress, 0)
	for _, n := range ings.Items {
		ingresshosts := make([]kubeapi.MappingIngressHost, 0)
		if spec := n.Spec.DefaultBackend; spec != nil {
			ingress := kubeapi.MappingIngressHost{
				Hostname: "*",
				PathType: "",
				Path:     "",
				UID:      string(n.GetObjectMeta().GetUID()),
			}

			if svc := spec.Service; svc != nil {
				ingress.BackendType = "service"
				ingress.BackendName = svc.Name
				ingress.ServicePort = svc.Port.Number
			} else if rsc := spec.Resource; rsc != nil {
				ingress.BackendType = "resource"
				ingress.BackendName = rsc.Name
				ingress.ResourceAPIGroup = *rsc.APIGroup
				ingress.ResourceKind = rsc.Kind
			}

			ingresshosts = append(ingresshosts, ingress)
		} else {
			for _, r := range n.Spec.Rules {
				hostname := r.Host

				for _, p := range r.HTTP.Paths {
					ingress := kubeapi.MappingIngressHost{
						Hostname: hostname,
						PathType: string(*p.PathType),
						Path:     p.Path,
						UID:      string(n.GetObjectMeta().GetUID()),
					}

					if svc := p.Backend.Service; svc != nil {
						ingress.BackendType = "service"
						ingress.BackendName = svc.Name
						ingress.ServicePort = svc.Port.Number
					} else if rsc := p.Backend.Resource; rsc != nil {
						ingress.BackendType = "resource"
						ingress.BackendName = rsc.Name
						ingress.ResourceAPIGroup = *rsc.APIGroup
						ingress.ResourceKind = rsc.Kind
					}

					ingresshosts = append(ingresshosts, ingress)
				}
			}
		}

		ingresses = append(ingresses, kubeapi.MappingIngress{
			Name:             n.GetObjectMeta().GetName(),
			UID:              string(n.GetObjectMeta().GetUID()),
			StartTime:        n.GetObjectMeta().GetCreationTimestamp().Time,
			Labels:           kubeapi.GetLabelSelector(n.GetObjectMeta().GetLabels()),
			IngressClassName: *n.Spec.IngressClassName,
			NamespaceName:    n.GetObjectMeta().GetNamespace(),
			IngressHosts:     ingresshosts,
		})
	}
	fmt.Printf("Ingress count is %d\n", len(ingresses))
}

func TestMakeStorageClassV1Data(t *testing.T) {
	clientset := authConfigfile(t)
	storagev1client := clientset.StorageV1()

	scs, err := storagev1client.StorageClasses().List(context.TODO(), metav1.ListOptions{})
	errorCheck(err, t)

	storageclasses := make([]kubeapi.MappingStorageClass, 0)
	for _, n := range scs.Items {
		var ave bool
		if n.AllowVolumeExpansion != nil {
			ave = *n.AllowVolumeExpansion
		}

		storageclasses = append(storageclasses, kubeapi.MappingStorageClass{
			Name:                 n.GetObjectMeta().GetName(),
			UID:                  string(n.GetObjectMeta().GetUID()),
			StartTime:            n.GetObjectMeta().GetCreationTimestamp().Time,
			Labels:               kubeapi.GetLabelSelector(n.GetObjectMeta().GetLabels()),
			Provisioner:          n.Provisioner,
			ReclaimPolicy:        string(*n.ReclaimPolicy),
			VolumeBindingMode:    string(*n.VolumeBindingMode),
			AllowVolumeExpansion: ave,
		})
	}
	fmt.Printf("StorageClass count is %d\n", len(storageclasses))
}
