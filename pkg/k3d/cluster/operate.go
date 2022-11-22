package cluster

import (
	"context"

	"github.com/nikhilsbhat/terraform-provider-k3d/pkg/utils"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
)

func (cfg *Config) StartStopCluster(ctx context.Context, runtime runtimes.Runtime, clusters []string) error {
	fetchedClusters, err := cfg.GetClusters(ctx, runtime, clusters)
	if err != nil {
		return err
	}

	if cfg.Action == utils.TerraformResourceStart {
		for _, cluster := range fetchedClusters {
			if err = client.ClusterStart(ctx, runtime, cluster.GetClusterConfig(), K3D.ClusterStartOpts{}); err != nil {
				return err
			}
		}

		return nil
	}

	for _, cluster := range fetchedClusters {
		if err = client.ClusterStop(ctx, runtime, cluster.GetClusterConfig()); err != nil {
			return err
		}
	}

	return nil
}
