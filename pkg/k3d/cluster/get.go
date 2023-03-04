package cluster

import (
	"context"

	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
	"github.com/thoas/go-funk"
)

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
