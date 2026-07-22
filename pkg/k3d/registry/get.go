package registry

import (
	"context"
	"slices"

	k3dNode "github.com/nikhilsbhat/terraform-provider-k3d/pkg/k3d/node"
	"github.com/rancher/k3d/v5/pkg/runtimes"
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

	filteredRegistries := make([]*k3dNode.Config, 0, len(regs))
	for _, reg := range regs {
		if slices.Contains(registry.Name, reg.Name[0]) {
			filteredRegistries = append(filteredRegistries, reg)
		}
	}

	return filteredRegistries, nil
}
