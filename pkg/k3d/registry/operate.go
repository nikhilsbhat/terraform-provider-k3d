package registry

import (
	"context"

	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	k3dNode "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/node"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
)

func (registry *Config) Connect(ctx context.Context, runtime runtimes.Runtime) error {
	clusters := make([]*K3D.Cluster, 0)
	clusterCfg := cluster.Config{}

	k3dClusters, err := clusterCfg.GetClusters(ctx, runtime, []string{registry.Cluster})
	if err != nil {
		return err
	}

	for _, k3dCluster := range k3dClusters {
		clusters = append(clusters, k3dCluster.GetClusterConfig())
	}

	regs, err := k3dNode.FilteredNodes(ctx, runtime, registry.Name)
	if err != nil {
		return err
	}

	for _, reg := range regs {
		if err = client.RegistryConnectClusters(ctx, runtime, reg, clusters); err != nil {
			return err
		}
	}

	return nil
}

func (registry *Config) Disconnect(ctx context.Context, runtime runtimes.Runtime) error {
	clusterCfg := cluster.Config{}

	k3dClusters, err := clusterCfg.GetClusters(ctx, runtime, []string{registry.Cluster})
	if err != nil {
		return err
	}

	regs, err := k3dNode.FilteredNodes(ctx, runtime, registry.Name)
	if err != nil {
		return err
	}

	for _, reg := range regs {
		if err = runtime.DisconnectNodeFromNetwork(ctx, reg, k3dClusters[0].Network); err != nil {
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
