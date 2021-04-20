package k3d

import (
	"context"
	"fmt"

	"github.com/rancher/k3d/v4/pkg/client"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

type Cluster struct {
	Name            string   `json:"name,omitempty"`
	Nodes           []string `json:"nodes,omitempty"`
	Network         string   `json:"network,omitempty"`
	Token           string   `json:"cluster_token,omitempty"`
	ServersCount    int64    `json:"servers_count,omitempty"`
	AgentsCount     int64    `json:"agents_count,omitempty"`
	AgentsRunning   int64    `json:"agents_running,omitempty"`
	ImageVolume     string   `json:"image_volume,omitempty"`
	HasLoadBalancer bool     `json:"has_loadbalancer,omitempty"`
}

func GetCluster(ctx context.Context, runtime runtimes.Runtime, cluster string) (*K3D.Cluster, error) {
	clusterConfig, err := client.ClusterGet(ctx, runtime, &K3D.Cluster{Name: cluster})
	if err != nil {
		return nil, err
	}
	return clusterConfig, nil
}

func GetFilteredClusters(ctx context.Context, runtime runtimes.Runtime, clusters []string) ([]*K3D.Cluster, error) {
	clustersList, err := client.ClusterList(ctx, runtime)
	if err != nil {
		return nil, err
	}
	var clusterConfig []*K3D.Cluster
	for _, clusterList := range clustersList {
		for _, cluster := range clusters {
			if clusterList.Name == cluster {
				clusterConfig = append(clusterConfig, clusterList)
			}
		}
	}
	if len(clusterConfig) == 0 {
		return nil, fmt.Errorf("cluster %v not found", clusters)
	}
	return clusterConfig, nil
}

func GetClusters(ctx context.Context, runtime runtimes.Runtime) ([]*K3D.Cluster, error) {
	clustersList, err := client.ClusterList(ctx, runtime)
	if err != nil {
		return nil, err
	}
	return clustersList, nil
}
