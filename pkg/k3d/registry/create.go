package registry

import (
	"context"
	"log"

	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
)

func (registry *Config) Create(ctx context.Context, runtime runtimes.Runtime) error {
	registryK3d := &K3D.Registry{}

	registryK3d.ClusterRef = registry.Cluster
	registryK3d.Protocol = registry.Protocol
	registryK3d.Host = registry.Host
	registryK3d.Image = registry.Image
	registryK3d.ExposureOpts = GetExposureOpts(registry.Expose)
	if registry.UseProxy {
		SetProxyConfig(registry.Proxy, registryK3d)
	}

	_, err := client.RegistryRun(ctx, runtime, registryK3d)
	if err != nil {
		return err
	}
	if len(registry.Cluster) != 0 {
		log.Printf("connecting the registry with cluster %v", registry.Cluster)
		return registry.Connect(ctx, runtime)
	}

	return nil
}
