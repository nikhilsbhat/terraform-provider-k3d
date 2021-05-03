package k3d

import (
	"context"
	"log"

	"github.com/docker/go-connections/nat"
	"github.com/rancher/k3d/v4/pkg/client"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

func GetRegistry(ctx context.Context, runtime runtimes.Runtime, cluster string, registry string) ([]*K3DNode, error) {
	var nodes []*K3DNode
	if len(cluster) != 0 {
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

func ConnectRegistryToCluster(ctx context.Context, runtime runtimes.Runtime, clusters []string, node *K3D.Node) error {
	k3dClusters, err := GetFilteredClusters(ctx, runtime, clusters)
	if err != nil {
		return err
	}
	return client.RegistryConnectClusters(ctx, runtime, node, k3dClusters)
}

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

func GetProxyConfig(proxyCfg map[string]string, registry *K3D.Registry) {
	registry.Options.Proxy.RemoteURL = proxyCfg["remoteURL"]
	registry.Options.Proxy.RemoteURL = proxyCfg["username"]
	registry.Options.Proxy.RemoteURL = proxyCfg["password"]
}
