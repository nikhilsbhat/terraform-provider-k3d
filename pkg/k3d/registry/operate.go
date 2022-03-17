package registry

import (
	"context"
	cluster2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	k3dNode "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/node"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
)

func (registry *Config) Connect(ctx context.Context, runtime runtimes.Runtime) error {
	k3dClusters, err := cluster2.GetFilteredClusters(ctx, runtime, []string{registry.Cluster})
	if err != nil {
		return err
	}

	regs, err := k3dNode.FilteredNodes(ctx, runtime, registry.Name)
	if err != nil {
		return err
	}

	for _, reg := range regs {
		if err = client.RegistryConnectClusters(ctx, runtime, reg, k3dClusters); err != nil {
			return err
		}
	}
	return nil
}

func (registry *Config) Disconnect(ctx context.Context, runtime runtimes.Runtime) error {
	k3dCluster, err := cluster2.GetCluster(ctx, runtime, registry.Cluster)
	if err != nil {
		return err
	}

	regs, err := k3dNode.FilteredNodes(ctx, runtime, registry.Name)
	if err != nil {
		return err
	}

	for _, reg := range regs {
		if err = runtime.DisconnectNodeFromNetwork(ctx, reg, k3dCluster.Network.Name); err != nil {
			return err
		}
	}
	return nil
}

func (registry *Config) GetRegistryStatus(registryName, state string) map[string]string {
	return map[string]string{
		"registry": registryName,
		"cluster":  registry.Cluster,
		"state":    state,
	}
}
