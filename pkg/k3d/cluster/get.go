package cluster

import (
	"context"
	"fmt"

	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
	"github.com/thoas/go-funk"
)

// GetFilteredClusters returns the list of *K3D.Config of specified clusters.
func GetFilteredClusters(ctx context.Context, runtime runtimes.Runtime,
	clusters []string,
) ([]*K3D.Cluster, error) {
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

func (cfg *Config) GetClusters(ctx context.Context, runtime runtimes.Runtime, clusterList []string) ([]*Config, error) {
	clusters, err := client.ClusterList(ctx, runtime)
	if err != nil {
		return nil, err
	}

	if !cfg.All {
		filteredCluster := funk.Filter(clusters, func(cluster *K3D.Cluster) bool {
			return funk.Contains(clusterList, cluster.Name)
		}).([]*K3D.Cluster)
		clusters = filteredCluster
	}

	clusterConfig := make([]*Config, 0)
	for _, cluster := range clusters {
		serverCount, serversRunning := cluster.ServerCountRunning()
		agentsCount, agentsRunning := cluster.AgentCountRunning()
		clusterConfig = append(clusterConfig, &Config{
			Name:            cluster.Name,
			Nodes:           funk.Get(cluster.Nodes, "Name").([]string),
			Network:         cluster.Network.Name,
			Token:           cluster.Token,
			ServersCount:    serverCount,
			ServersRunning:  serversRunning,
			AgentsCount:     agentsCount,
			AgentsRunning:   agentsRunning,
			ImageVolume:     cluster.ImageVolume,
			HasLoadBalancer: cluster.HasLoadBalancer(),
		})
	}

	return clusterConfig, nil
}

func (cfg *Config) GetClusterConfig() *K3D.Cluster {
	nodes := make([]*K3D.Node, 0)
	for _, node := range cfg.Nodes {
		nodes = append(nodes, &K3D.Node{Name: node})
	}

	return &K3D.Cluster{
		Name:  cfg.Name,
		Token: cfg.Token,
		Nodes: nodes,
	}
}
