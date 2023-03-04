package cluster

import (
	"context"
	"fmt"

	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/config"
	"github.com/rancher/k3d/v5/pkg/config/v1alpha4"
	"github.com/rancher/k3d/v5/pkg/runtimes"
)

func CreateCluster(ctx context.Context, runtime runtimes.Runtime, cfg *v1alpha4.SimpleConfig) error {
	clusterConfig, err := config.TransformSimpleToClusterConfig(ctx, runtime, *cfg)
	if err != nil {
		return err
	}

	clusterConfig, err = config.ProcessClusterConfig(*clusterConfig)
	if err != nil {
		return err
	}

	if err = config.ValidateClusterConfig(ctx, runtimes.SelectedRuntime, *clusterConfig); err != nil {
		return err
	}

	if _, err = client.ClusterGet(ctx, runtimes.SelectedRuntime, &clusterConfig.Cluster); err == nil {
		return fmt.Errorf("failed to create cluster '%s' because a cluster with that name already exists: %w", cfg.Name, err)
	}

	if err = client.ClusterRun(ctx, runtimes.SelectedRuntime, clusterConfig); err != nil {
		return err
	}

	if clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig {
		if _, err = client.KubeconfigGetWrite(ctx, runtimes.SelectedRuntime,
			&clusterConfig.Cluster, "",
			&client.WriteKubeConfigOptions{
				UpdateExisting:       true,
				OverwriteExisting:    false,
				UpdateCurrentContext: cfg.Options.KubeconfigOptions.SwitchCurrentContext,
			}); err != nil {
			return err
		}
	}

	return nil
}

func CreateK3DCluster(ctx context.Context, runtime runtimes.Runtime, cfg *v1alpha4.SimpleConfig) error {
	return nil
}
