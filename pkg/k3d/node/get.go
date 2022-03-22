package node

import (
	"context"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
	"github.com/thoas/go-funk"
)

// Nodes fetches details of all available nodes in the specified cluster.
func Nodes(ctx context.Context, runtime runtimes.Runtime, cluster string) ([]*K3D.Node, error) {
	nodes, err := runtime.GetNodesByLabel(ctx, map[string]string{
		"cluster": cluster,
	})
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// Node fetches details of a specified node.
func Node(ctx context.Context, runtime runtimes.Runtime, name string) (*K3D.Node, error) {
	node, err := client.NodeGet(ctx, runtime, &K3D.Node{Name: name})
	if err != nil {
		return nil, err
	}
	return node, err
}

// GetFilteredNodesFromCluster returns the fetched all nodes from a specified cluster with list of *Config type.
func (cfg *Config) GetFilteredNodesFromCluster(ctx context.Context, runtime runtimes.Runtime) ([]*Config, error) {
	cfg.Labels = map[string]string{
		"k3d.role":    "agent",
		"k3d.cluster": cfg.ClusterAssociated,
	}
	if cfg.All {
		return cfg.GetNodesByLabels(ctx, runtime)
	}
	return cfg.GetFilteredNodes(ctx, runtime, cfg.Name)
}

// GetFilteredNodes returns the fetched list of specified nodes from specified cluster with list of *Config type.
func (cfg *Config) GetFilteredNodes(ctx context.Context, runtime runtimes.Runtime, nodes []string) ([]*Config, error) {
	k3dNodes, err := cfg.GetNodesByLabels(ctx, runtime)
	if err != nil {
		return nil, err
	}

	filteredNodes := funk.Filter(k3dNodes, func(node *Config) bool {
		if funk.Contains(nodes, node.Name[0]) {
			return false
		}
		return true
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
			//PortMapping:          getPortMaps(node.Ports),
		})
	}
	return filteredNodes, err
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
