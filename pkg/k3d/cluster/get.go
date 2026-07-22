package cluster

import (
	"context"
	"slices"

	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
)

func (cfg *Config) GetClusters(ctx context.Context, runtime runtimes.Runtime, clusterList []string) ([]*Config, error) {
	clusters, err := client.ClusterList(ctx, runtime)
	if err != nil {
		return nil, err
	}

	if !cfg.All {
		clusters = filterClustersByName(clusters, clusterList)
	}

	clusterConfig := make([]*Config, 0, len(clusters))

	for _, cluster := range clusters {
		serverCount, serversRunning := cluster.ServerCountRunning()

		agentsCount, agentsRunning := cluster.AgentCountRunning()

		clusterConfig = append(clusterConfig, &Config{
			Name:            cluster.Name,
			Nodes:           nodeNames(cluster.Nodes),
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

func filterClustersByName(clusters []*K3D.Cluster, names []string) []*K3D.Cluster {
	filteredClusters := make([]*K3D.Cluster, 0, len(clusters))

	for _, cluster := range clusters {
		if slices.Contains(names, cluster.Name) {
			filteredClusters = append(filteredClusters, cluster)
		}
	}

	return filteredClusters
}

func nodeNames(nodes []*K3D.Node) []string {
	names := make([]string, 0, len(nodes))

	for _, node := range nodes {
		names = append(names, node.Name)
	}

	return names
}

func (cfg *Config) GetClusterConfig() *K3D.Cluster {
	nodes := make([]*K3D.Node, 0, len(cfg.Nodes))

	for _, node := range cfg.Nodes {
		nodes = append(nodes, &K3D.Node{Name: node})
	}

	return &K3D.Cluster{
		Name:  cfg.Name,
		Token: cfg.Token,
		Nodes: nodes,
	}
}
