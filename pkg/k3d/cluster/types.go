package cluster

import (
	"context"

	"github.com/rancher/k3d/v5/pkg/runtimes"
)

type Cluster interface {
	GetClusters(context.Context, runtimes.Runtime, []string) ([]*Config, error)
	StartStopCluster(context.Context, runtimes.Runtime, []string) error
}

// Config helps storing filtered cluster data of k3d cluster.
//
//nolint:maligned
type Config struct {
	Name            string   `json:"name,omitempty" mapstructure:"name"`
	Nodes           []string `json:"nodes,omitempty" mapstructure:"nodes"`
	Network         string   `json:"network,omitempty" mapstructure:"network"`
	Token           string   `json:"cluster_token,omitempty" mapstructure:"cluster_token"`
	ServersCount    int      `json:"servers_count,omitempty" mapstructure:"servers_count"`
	ServersRunning  int      `json:"servers_running,omitempty" mapstructure:"servers_running"`
	AgentsCount     int      `json:"agents_count,omitempty" mapstructure:"agents_count"`
	AgentsRunning   int      `json:"agents_running,omitempty" mapstructure:"agents_running"`
	ImageVolume     string   `json:"image_volume,omitempty" mapstructure:"image_volume"`
	HasLoadBalancer bool     `json:"has_loadbalancer,omitempty" mapstructure:"has_loadbalancer"`
	Action          string   `json:"action,omitempty" mapstructure:"action"`
	All             bool     `json:"all,omitempty" mapstructure:"all"`
}
