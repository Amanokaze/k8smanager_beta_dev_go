package kubeapi

const (
	KIND_NODE = "node"
	KIND_POD  = "pod"
	KIND_NS   = "namespace"
	KIND_SVC  = "service"
	KIND_PV   = "persistentvolume"
	KIND_PVC  = "persistentvolumeclaim"

	ENDPOINT_PREFIX_API  = "api/"
	ENDPOINT_PREFIX_APIS = "apis/"

	METRIC_CONTAINER_CPU_USAGE_SECONDS_TOTAL       = "container_cpu_usage_seconds_total"
	METRIC_CONTAINER_CPU_SYSTEM_SECONDS_TOTAL      = "container_cpu_system_seconds_total"
	METRIC_CONTAINER_CPU_USER_SECONDS_TOTAL        = "container_cpu_user_seconds_total"
	METRIC_CONTAINER_MEMORY_USAGE_BYTES            = "container_memory_usage_bytes"
	METRIC_CONTAINER_MEMORY_WORKING_SET_BYTES      = "container_memory_working_set_bytes"
	METRIC_CONTAINER_MEMORY_CACHE                  = "container_memory_cache"
	METRIC_CONTAINER_MEMORY_SWAP                   = "container_memory_swap"
	METRIC_CONTAINER_MEMORY_RSS                    = "container_memory_rss"
	METRIC_CONTAINER_FS_INODES_FREE                = "container_fs_inodes_free"
	METRIC_CONTAINER_FS_INODES_TOTAL               = "container_fs_inodes_total"
	METRIC_CONTAINER_FS_LIMIT_BYTES                = "container_fs_limit_bytes"
	METRIC_CONTAINER_FS_READS_BYTES_TOTAL          = "container_fs_reads_bytes_total"
	METRIC_CONTAINER_FS_WRITES_BYTES_TOTAL         = "container_fs_writes_bytes_total"
	METRIC_CONTAINER_FS_USAGE_BYTES                = "container_fs_usage_bytes"
	METRIC_CONTAINER_PROCESSES                     = "container_processes"
	METRIC_CONTAINER_NETWORK_RECEIVE_BYTES_TOTAL   = "container_network_receive_bytes_total"
	METRIC_CONTAINER_NETWORK_RECEIVE_ERRORS_TOTAL  = "container_network_receive_errors_total"
	METRIC_CONTAINER_NETWORK_TRANSMIT_BYTES_TOTAL  = "container_network_transmit_bytes_total"
	METRIC_CONTAINER_NETWORK_TRANSMIT_ERRORS_TOTAL = "container_network_transmit_errors_total"
)
