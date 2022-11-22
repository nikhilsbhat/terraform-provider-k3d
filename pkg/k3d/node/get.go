package node

import (
	"context"

	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
	"github.com/thoas/go-funk"
)

// GetFilteredNodesFromCluster returns the fetched all nodes from a specified cluster with list of *Config type.
func (cfg *Config) GetFilteredNodesFromCluster(ctx context.Context, runtime runtimes.Runtime) ([]*Config, error) {
	cfg.Labels = map[string]string{
		"k3d.role":    "agent",
		"k3d.cluster": cfg.ClusterAssociated,
	}
	if cfg.All {
		return cfg.GetNodesByLabels(ctx, runtime)
	}

	return cfg.GetFilteredNodes(ctx, runtime)
}

// GetFilteredNodes returns the fetched list of specified nodes from specified cluster with list of *Config type.
func (cfg *Config) GetFilteredNodes(ctx context.Context, runtime runtimes.Runtime) ([]*Config, error) {
	k3dNodes, err := cfg.GetNodesByLabels(ctx, runtime)
	if err != nil {
		return nil, err
	}

	filteredNodes := funk.Filter(k3dNodes, func(node *Config) bool {
		return funk.Contains(cfg.Name, node.Name[0])
	}).([]*Config)

	return filteredNodes, nil
}

// GetNodesByLabels gets the nodes that matches with the specified label.
func (cfg *Config) GetNodesByLabels(ctx context.Context, runtime runtimes.Runtime) ([]*Config, error) {
	k3dNodes, err := runtime.GetNodesByLabel(ctx, cfg.Labels)
	if err != nil {
		return nil, err
	}

	filteredNodes := make([]*Config, 0)
	for _, node := range k3dNodes {
		filteredNodes = append(filteredNodes, &Config{
			Name:                 []string{node.Name},
			Role:                 string(node.Role),
			ClusterAssociated:    node.RuntimeLabels[K3dClusterNameLabel],
			State:                node.State.Status,
			Created:              node.Created,
			Memory:               node.Memory,
			Volumes:              node.Volumes,
			Networks:             node.Networks,
			EnvironmentVariables: node.Env,
			Image:                node.Image,
			// dropping PortMapping as terraform schema format is yet to be figured.
			// PortMapping:          getPortMaps(node.Ports),
		})
	}

	return filteredNodes, err
}

// GetNodeStatus retrieves the latest state of the nodes.
func (cfg *Config) GetNodeStatus(ctx context.Context, runtime runtimes.Runtime) ([]*Status, error) {
	nodes, err := cfg.GetFilteredNodesFromCluster(ctx, runtime)
	if err != nil {
		return nil, err
	}

	nodeCurrentStatus := make([]*Status, 0)
	for _, node := range nodes {
		nodeCurrentStatus = append(nodeCurrentStatus, &Status{
			Node:    node.Name[0],
			Cluster: node.ClusterAssociated,
			State:   node.State,
			Role:    node.Role,
		})
	}

	return nodeCurrentStatus, nil
}

// GetNodeFromConfig returns K3D.Node equivalent for an stance of Config.
func (cfg *Config) GetNodeFromConfig() *K3D.Node {
	return &K3D.Node{
		Name: cfg.Name[0],
		Role: K3D.NodeRoles[cfg.Role],
		RuntimeLabels: map[string]string{
			K3D.LabelRole:                  cfg.Role,
			utils.TerraformK3dLabel:        cfg.Created,
			utils.TerraformCreatedK3dLabel: "true",
		},
		Image:   cfg.Image,
		Memory:  cfg.Memory,
		Restart: true,
	}
}
