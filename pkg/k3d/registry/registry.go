package registry

import (
	"context"
	"log"

	"github.com/docker/go-connections/nat"
	cluster2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	k3dNode "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/node"
	"github.com/rancher/k3d/v4/pkg/client"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	K3D "github.com/rancher/k3d/v4/pkg/types"
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

// CreateRegistry creates registry and connects it to specified cluster.
func CreateRegistry(ctx context.Context, runtime runtimes.Runtime, reg *K3D.Registry, clusters []string) error {
	regNode, err := client.RegistryRun(ctx, runtime, reg)
	if err != nil {
		return err
	}
	if len(clusters) != 0 {
		log.Printf("connecting the registry with cluster %v", clusters)
		return ConnectRegistryToCluster(ctx, runtimes.SelectedRuntime, clusters, regNode)
	}
	return nil
}

// GetExposureOpts fetches expose data and adds it to K3D.Config.
func GetExposureOpts(expose map[string]string, registry *K3D.Registry) {
	binding := nat.PortBinding{
		HostIP:   expose["hostIp"],
		HostPort: expose["hostPort"],
	}
	api := &K3D.ExposureOpts{}
	api.Port = nat.Port(expose["hostPort"])

	api.Binding = binding
	registry.ExposureOpts = *api
}

// GetProxyConfig fetches passed proxy config and adds it to K3D.Config.
func GetProxyConfig(proxyCfg map[string]string, registry *K3D.Registry) {
	registry.Options.Proxy.RemoteURL = proxyCfg["remoteURL"]
	registry.Options.Proxy.RemoteURL = proxyCfg["username"]
	registry.Options.Proxy.RemoteURL = proxyCfg["password"]
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
