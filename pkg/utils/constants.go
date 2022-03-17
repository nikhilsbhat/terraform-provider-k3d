package utils

const (
	TerraformResourceCluster          = "cluster"
	TerraformResourceImage            = "image"
	TerraformResourceImages           = "images"
	TerraformResourceClusters         = "clusters"
	TerraformResourceImagesStored     = "images_stored"
	TerraformResourceKeepTarball      = "keep_tarball"
	TerraformResourceTarballStored    = "tarball_stored"
	TerraformResourceNodes            = "nodes"
	TerraformResourceNodesList        = "node_list"
	TerraformResourceClusterList      = "clusters_list"
	TerraformResourceAll              = "all"
	TerraformResourceNotEncode        = "not_encoded"
	TerraformResourceEncode           = "encode"
	TerraformResourceKubeConfig       = "kube_config"
	TerraformResourceStatus           = "status"
	TerraformResourceStart            = "start"
	TerraformResourceStop             = "stop"
	TerraformResourceRole             = "role"
	TerraformResourceReplicas         = "replicas"
	TerraformResourceWait             = "wait"
	TerraformResourceTimeout          = "timeout"
	TerraformResourceMemory           = "memory"
	TerraformResourceCreatedAt        = "creation_time"
	TerraformResourceName             = "name"
	TerraformResourceState            = "state"
	TerraformResourceRegistries       = "registries"
	TerraformResourceRegistriesList   = "registries_list"
	TerraformResourcePorts            = "ports"
	TerraformResourceHost             = "host"
	TerraformResourceExpose           = "expose"
	TerraformResourceProtocol         = "protocol"
	TerraformResourceProxy            = "proxy"
	TerraformResourceMetadata         = "metadata"
	TerraformResourceConnect          = "connect"
	TerraformUseProxy                 = "use_proxy"
	TerraformResourceConfigFile       = "config_file"
	TerraformResourceServersCount     = "servers_count"
	TerraformResourceAgentsCount      = "agents_count"
	TerraformResourceNetwork          = "network"
	TerraformResourceClusterToken     = "cluster_token"
	TerraformResourceUpdateKubeConfig = "kubeconfig_update_default"
	TerraformResourceSwitchKubeConfig = "kubeconfig_switch_context"
	TerraformResourceLabels           = "labels"
	TerraformResourceEnv              = "env"
	TerraformResourceVolumes          = "volumes"
	TerraformResourceK3dOptions       = "k3d_options"
	TerraformResourceK3sOptions       = "k3s_options"

	TerraformK3dLabel          = "k3d.terraform"
	TerraformK3dRegistry       = "registry"
	TerraformKubernetesVersion = "kubernetes_version"
	TerraformK3dAPIVersion     = "k3d_api_version"
	TerraformK3dKind           = "kind"
	TerraformK3dRuntime        = "runtime"
	TerraformTimeOut5          = 5
	K3DRepoDEFAULT             = "rancher/k3s"
	RegistryConnectedState     = "connected"
	RegistryDisconnectedState  = "disconnected"
)
