package k3d

import (
	"context"
)

// Cluster helps storing filtered cluster data of k3d cluster.
type Cluster struct {
	Name            string   `json:"name,omitempty"`
	Nodes           []string `json:"nodes,omitempty"`
	Network         string   `json:"network,omitempty"`
	Token           string   `json:"cluster_token,omitempty"`
	ServersCount    int      `json:"servers_count,omitempty"`
	ServersRunning  int      `json:"servers_running,omitempty"`
	AgentsCount     int      `json:"agents_count,omitempty"`
	AgentsRunning   int      `json:"agents_running,omitempty"`
	ImageVolume     string   `json:"image_volume,omitempty"`
	HasLoadBalancer bool     `json:"has_loadbalancer,omitempty"`
}

// Images helps storing filtered images data that was loaded to k3d cluster.
type Images struct {
	Images       []string        `json:"images,omitempty"`
	Cluster      string          `json:"cluster,omitempty"`
	StoreTarBall bool            `json:"keep_tarball,omitempty"`
	StoredImages StoredImages    `json:"images_stored,omitempty"`
	Context      context.Context `json:"context,omitempty"`
	Config       Config          `json:"config,omitempty"`
}

// StoredImages holds a data of cluster to images mapping of loaded images.
type StoredImages struct {
	Cluster string   `json:"cluster,omitempty"`
	Images  []string `json:"images,omitempty"`
}

// TarBallData maps tarball stored to image.
type TarBallData struct {
	Image string `json:"image,omitempty"`
	Path  string `json:"path,omitempty"`
}

// K3Node helps storing filtered node data of k3d cluster.
type K3Node struct {
	Name                 string                 `json:"name,omitempty"`
	Role                 string                 `json:"role,omitempty"`
	ClusterAssociated    string                 `json:"cluster,omitempty"`
	State                string                 `json:"state,omitempty"`
	Created              string                 `json:"created,omitempty"`
	Memory               string                 `json:"memory,omitempty"`
	Volumes              []string               `json:"volumes,omitempty"`
	Networks             []string               `json:"networks,omitempty"`
	EnvironmentVariables []string               `json:"env,omitempty"`
	Count                int                    `json:"count,omitempty"`
	Image                string                 `json:"image,omitempty"`
	PortMapping          map[string]interface{} `json:"port_mappings,omitempty"`
}

// NodeStatus helps storing filtered node status of k3d cluster.
type NodeStatus struct {
	Node    string `json:"node,omitempty"`
	Cluster string `json:"cluster,omitempty"`
	Role    string `json:"role,omitempty"`
	State   string `json:"state,omitempty"`
	Running bool   `json:"running,omitempty"`
}

// Registry helps storing filtered registry data the present in selected runtime.
type Registry struct {
	Name     string            `json:"name,omitempty"`
	Image    string            `json:"image,omitempty"`
	Cluster  string            `json:"cluster,omitempty"`
	Protocol string            `json:"protocol,omitempty"`
	Host     string            `json:"host,omitempty"`
	Port     string            `json:"port,omitempty"`
	Expose   map[string]string `json:"expose,omitempty"`
	UseProxy bool              `json:"use_proxy,omitempty"`
	Proxy    map[string]string `json:"proxy,omitempty"`
}

// RegistryConnect holds the status of registries connected to cluster.
type RegistryConnect struct {
	Registries []string `json:"registries,omitempty"`
	Cluster    string   `json:"cluster,omitempty"`
	Connect    bool     `json:"connect,omitempty"`
}
