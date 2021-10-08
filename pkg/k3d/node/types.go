package node

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

// Status helps to store filtered node status of k3d cluster.
type Status struct {
	Node    string `json:"node,omitempty"`
	Cluster string `json:"cluster,omitempty"`
	Role    string `json:"role,omitempty"`
	State   string `json:"state,omitempty"`
	Running bool   `json:"running,omitempty"`
}
