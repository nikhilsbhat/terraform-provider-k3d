package k3d

import (
	"context"
	"fmt"
	"log"

	"github.com/rancher/k3d/v4/pkg/client"
	"github.com/rancher/k3d/v4/pkg/config"
	"github.com/rancher/k3d/v4/pkg/config/v1alpha2"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

func CreateCluster(ctx context.Context, runtime runtimes.Runtime, cfg *v1alpha2.SimpleConfig) error {
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
		return fmt.Errorf("failed to create cluster because a cluster with that name already exists: %v", err)
	}

	// create cluster
	if err = client.ClusterRun(ctx, runtimes.SelectedRuntime, clusterConfig); err != nil {
		// rollback if creation failed
		if deleteErr := client.ClusterDelete(ctx, runtimes.SelectedRuntime, &K3D.Cluster{Name: cfg.Name},
			K3D.ClusterDeleteOpts{SkipRegistryCheck: false}); deleteErr != nil {
			return fmt.Errorf("cluster creation FAILED, also FAILED to rollback changes!: %v", deleteErr)
		}
		return err
	}

	// update default kubeconfig
	if clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig {
		if _, err := client.KubeconfigGetWrite(ctx, runtimes.SelectedRuntime, &clusterConfig.Cluster, "", &client.WriteKubeConfigOptions{UpdateExisting: true, OverwriteExisting: false, UpdateCurrentContext: cfg.Options.KubeconfigOptions.SwitchCurrentContext}); err != nil {
			return err
		}
	}

	return nil
}

// GetCluster is a wrap of client.ClusterGet of k3d.
func GetCluster(ctx context.Context, runtime runtimes.Runtime,
	cluster string) (*K3D.Cluster, error) {
	clusterConfig, err := client.ClusterGet(ctx, runtime, &K3D.Cluster{Name: cluster})
	if err != nil {
		return nil, err
	}
	return clusterConfig, nil
}

// GetFilteredClusters returns the list of *K3D.Cluster of specified clusters.
func GetFilteredClusters(ctx context.Context, runtime runtimes.Runtime,
	clusters []string) ([]*K3D.Cluster, error) {
	clustersList, err := client.ClusterList(ctx, runtime)
	if err != nil {
		return nil, err
	}
	var clusterConfig []*K3D.Cluster
	for _, clusterList := range clustersList {
		for _, cluster := range clusters {
			if clusterList.Name == cluster {
				clusterConfig = append(clusterConfig, clusterList)
			}
		}
	}
	if len(clusterConfig) == 0 {
		return nil, fmt.Errorf("cluster %v not found", clusters)
	}
	return clusterConfig, nil
}

// GetClusters return the list of *K3D.Cluster of all clusters available in the specified runtime.
func GetClusters(ctx context.Context, runtime runtimes.Runtime) ([]*K3D.Cluster, error) {
	clustersList, err := client.ClusterList(ctx, runtime)
	if err != nil {
		return nil, err
	}
	return clustersList, nil
}

// StartClusters starts the specified clusters
func StartClusters(ctx context.Context, runtime runtimes.Runtime,
	clusters []*K3D.Cluster, options K3D.ClusterStartOpts) error {
	for _, cluster := range clusters {
		if err := client.ClusterStart(ctx, runtime, cluster, options); err != nil {
			log.Fatalln(err)
		}
	}
	return nil
}

// StopClusters stops the specified clusters.
func StopClusters(ctx context.Context, runtime runtimes.Runtime,
	clusters []*K3D.Cluster) error {
	for _, cluster := range clusters {
		if err := client.ClusterStop(ctx, runtime, cluster); err != nil {
			log.Fatalln(err)
		}
	}
	return nil
}
