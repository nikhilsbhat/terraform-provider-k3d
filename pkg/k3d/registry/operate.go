package registry

import (
	"context"
	cluster2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
)

// ConnectRegistryToCluster adds specified registry to a cluster.
func ConnectRegistryToCluster(ctx context.Context, runtime runtimes.Runtime, clusters []string, node *K3D.Node) error {
	k3dClusters, err := cluster2.GetFilteredClusters(ctx, runtime, clusters)
	if err != nil {
		return err
	}
	return client.RegistryConnectClusters(ctx, runtime, node, k3dClusters)
}

// DisconnectRegistryFormCluster disconnects registry from a specfied cluster.
func DisconnectRegistryFormCluster(ctx context.Context, runtime runtimes.Runtime,
	cluster string, node *K3D.Node) error {
	k3dCluster, err := cluster2.GetCluster(ctx, runtime, cluster)
	if err != nil {
		return err
	}
	return runtime.DisconnectNodeFromNetwork(ctx, node, k3dCluster.Network.Name)
}

// ConnectRegistriesToCluster connects specified registries to cluster.
func ConnectRegistriesToCluster(ctx context.Context, runtime runtimes.Runtime,
	clusters []string, nodes []*K3D.Node) error {
	k3dClusters, err := cluster2.GetFilteredClusters(ctx, runtime, clusters)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		if err = client.RegistryConnectClusters(ctx, runtime, node, k3dClusters); err != nil {
			return err
		}
	}
	return nil
}

// DisconnectRegistriesFormCluster disconnects specified registries to cluster.
func DisconnectRegistriesFormCluster(ctx context.Context, runtime runtimes.Runtime,
	cluster string, nodes []*K3D.Node) error {
	k3dCluster, err := cluster2.GetCluster(ctx, runtime, cluster)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		if err = runtime.DisconnectNodeFromNetwork(ctx, node, k3dCluster.Network.Name); err != nil {
			return err
		}
	}
	return nil
}
