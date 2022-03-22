package registry

import (
	"context"
	k3dNode "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/node"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	"github.com/thoas/go-funk"
)

// Get fetches the information of the list of selected registries.
func (registry *Config) Get(ctx context.Context, runtime runtimes.Runtime) ([]*k3dNode.Config, error) {
	cfg := k3dNode.Config{Labels: map[string]string{"k3d.role": "registry", "k3d.cluster": registry.Cluster}}
	regs, err := cfg.GetNodesByLabels(ctx, runtime)
	if err != nil {
		return nil, err
	}

	if registry.All {
		return regs, nil
	}

	filteredRegistries := funk.Filter(regs, func(reg *k3dNode.Config) bool {
		if funk.Contains(registry.Name, reg.Name) {
			return true
		}
		return false
	}).([]*k3dNode.Config)

	return filteredRegistries, nil
}
