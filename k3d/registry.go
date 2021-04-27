package k3d

import (
	"context"

	"github.com/rancher/k3d/v4/pkg/runtimes"
)

func GetRegistry(ctx context.Context, runtime runtimes.Runtime, cluster string, registry string) ([]*K3DNode, error) {
	var nodes []*K3DNode
	if len(cluster) == 0 {
		k3dNodes, err := GetFilteredNodesFromCluster(ctx, runtime, cluster)
		if err != nil {
			return nil, err
		}
		for _, k3dNode := range k3dNodes {
			if k3dNode.Role == "registry" && k3dNode.Name == registry {
				nodes = append(nodes, k3dNode)
			}
		}
		return nodes, nil
	}
	k3dNodes, err := GetNodesByLabels(ctx, runtime, map[string]string{"k3d.role": "registry"})
	if err != nil {
		return nil, err
	}
	for _, k3dNode := range k3dNodes {
		if k3dNode.Name == registry {
			nodes = append(nodes, k3dNode)
		}
	}
	return nodes, nil
}

func GetRegistries(ctx context.Context, runtime runtimes.Runtime, cluster string) ([]*K3DNode, error) {
	if len(cluster) != 0 {
		k3dNodes, err := GetFilteredNodesFromCluster(ctx, runtime, cluster)
		if err != nil {
			return nil, err
		}
		var nodes []*K3DNode
		for _, k3dNode := range k3dNodes {
			if k3dNode.Role == "registry" {
				nodes = append(nodes, k3dNode)
			}
		}
		return nodes, nil
	}
	return GetNodesByLabels(ctx, runtime, map[string]string{"k3d.role": "registry"})
}

func GetRegistriesWithName(ctx context.Context, runtime runtimes.Runtime, cluster string, registries []string) ([]*K3DNode, error) {
	nodes := make([]*K3DNode, 0)
	for _, registry := range registries {
		regs, err := GetRegistry(ctx, runtime, cluster, registry)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, regs...)
	}
	return nodes, nil
}
