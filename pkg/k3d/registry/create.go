package registry

import (
	"context"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
	"log"
)

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

func (registry *Config) CreateRegistry(ctx context.Context, runtime runtimes.Runtime) error {
	registryK3d := &K3D.Registry{}

	registryK3d.ClusterRef = registry.Cluster
	registryK3d.Protocol = registry.Protocol
	registryK3d.Host = registry.Host
	registryK3d.Image = registry.Image
	registryK3d.ExposureOpts = GetExposureOpts(registry.Expose)
	if registry.UseProxy {
		SetProxyConfig(registry.Proxy, registryK3d)
	}

	if err := CreateRegistry(ctx, runtime, registryK3d, []string{registry.Cluster}); err != nil {
		return err
	}
	return nil
}
