package cluster

import (
	"context"
	"fmt"

	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/config"
	"github.com/rancher/k3d/v5/pkg/config/v1alpha4"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
)

func CreateCluster(ctx context.Context, runtime runtimes.Runtime, cfg *v1alpha4.SimpleConfig) error {
	// transform simple config to cluster config

	clusterConfig, err := config.TransformSimpleToClusterConfig(ctx, runtime, *cfg)
	if err != nil {
		return err
	}

	// process cluster config
	clusterConfig, err = config.ProcessClusterConfig(*clusterConfig)
	if err != nil {
		return err
	}

	// validate cluster config
	if err = config.ValidateClusterConfig(ctx, runtimes.SelectedRuntime, *clusterConfig); err != nil {
		return err
	}

	// check if a cluster with that name exists already
	if _, err = client.ClusterGet(ctx, runtimes.SelectedRuntime, &clusterConfig.Cluster); err == nil {
		return fmt.Errorf("failed to create cluster because a cluster with that name already exists: %w", err)
	}

	// create cluster
	if err = client.ClusterRun(ctx, runtimes.SelectedRuntime, clusterConfig); err != nil {
		// rollback if creation failed
		if deleteErr := client.ClusterDelete(ctx, runtimes.SelectedRuntime, &K3D.Cluster{Name: cfg.Name},
			K3D.ClusterDeleteOpts{SkipRegistryCheck: false}); deleteErr != nil {
			return fmt.Errorf("cluster creation FAILED, also FAILED to rollback changes!: %w", deleteErr)
		}

		return err
	}

	// update default kubeconfig
	if clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig {
		if _, err := client.KubeconfigGetWrite(ctx, runtimes.SelectedRuntime,
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
