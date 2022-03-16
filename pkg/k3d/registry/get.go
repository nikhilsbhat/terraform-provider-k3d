package registry

import (
	"context"
	k3dNode "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/node"
	"github.com/rancher/k3d/v5/pkg/runtimes"
)

// GetRegistry fetches the details of specified registry from a specified cluster.
func GetRegistry(ctx context.Context, runtime runtimes.Runtime, cluster string, registry string) ([]*k3dNode.K3Node, error) {
	var nodes []*k3dNode.K3Node
	if len(cluster) != 0 {
		k3dNodes, err := k3dNode.GetFilteredNodesFromCluster(ctx, runtime, cluster)
		if err != nil {
			return nil, err
		}
		for _, node := range k3dNodes {
			if node.Role == "registry" && node.Name == registry {
				nodes = append(nodes, node)
			}
		}
		return nodes, nil
	}
	k3dNodes, err := k3dNode.GetNodesByLabels(ctx, runtime, map[string]string{"k3d.role": "registry"})
	if err != nil {
		return nil, err
	}
	for _, node := range k3dNodes {
		if node.Name == registry {
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}

// GetRegistries fetches the details of all registries from a specified cluster.
func GetRegistries(ctx context.Context, runtime runtimes.Runtime, cluster string) ([]*k3dNode.K3Node, error) {
	if len(cluster) != 0 {
		k3dNodes, err := k3dNode.GetFilteredNodesFromCluster(ctx, runtime, cluster)
		if err != nil {
			return nil, err
		}
		var nodes []*k3dNode.K3Node
		for _, node := range k3dNodes {
			if node.Role == "registry" {
				nodes = append(nodes, node)
			}
		}
		return nodes, nil
	}
	return k3dNode.GetNodesByLabels(ctx, runtime, map[string]string{"k3d.role": "registry"})
}

// GetRegistriesWithName fetches the details of all specified registries from a specified cluster.
func GetRegistriesWithName(ctx context.Context, runtime runtimes.Runtime,
	cluster string, registries []string) ([]*k3dNode.K3Node, error) {
	nodes := make([]*k3dNode.K3Node, 0)
	for _, registry := range registries {
		regs, err := GetRegistry(ctx, runtime, cluster, registry)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, regs...)
	}
	return nodes, nil
}
