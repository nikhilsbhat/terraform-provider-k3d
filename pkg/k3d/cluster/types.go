package cluster

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
