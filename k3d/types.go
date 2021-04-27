package k3d

import (
	"context"
)

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

type K3Dimages struct {
	Images       []string        `json:"images,omitempty"`
	Cluster      string          `json:"cluster,omitempty"`
	StoreTarBall bool            `json:"keep_tarball,omitempty"`
	StoredImages StoredImages    `json:"images_stored,omitempty"`
	Context      context.Context `json:"context,omitempty"`
	Config       K3dConfig       `json:"config,omitempty"`
}

type StoredImages struct {
	Cluster string   `json:"cluster,omitempty"`
	Images  []string `json:"images,omitempty"`
}

type TarBallData struct {
	Image string `json:"image,omitempty"`
	Path  string `json:"path,omitempty"`
}

type K3DNode struct {
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

type K3DNodeStatus struct {
	Node    string `json:"node,omitempty"`
	Cluster string `json:"cluster,omitempty"`
	Role    string `json:"role,omitempty"`
	State   string `json:"state,omitempty"`
	Running bool   `json:"running,omitempty"`
}
