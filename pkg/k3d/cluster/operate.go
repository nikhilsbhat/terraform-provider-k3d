package cluster

import (
	"context"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
	"log"
)

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

//func (cfg *Config) StartStopCluster(ctx context.Context, runtime runtimes.Runtime, clusters []string) error {
//	fetchedClusters, err := cfg.GetClusters(ctx, runtime, clusters)
//	if err != nil {
//		return err
//	}
//
//	if cfg.Action == utils.TerraformResourceStart {
//		for _, cluster := range fetchedClusters {
//			if err = client.ClusterStart(ctx, runtime, cluster.GetClusterConfig(), K3D.ClusterStartOpts{}); err != nil {
//				return err
//			}
//		}
//	}
//
//	for _, cluster := range fetchedClusters {
//		if err = client.ClusterStop(ctx, runtime, cluster.GetClusterConfig()); err != nil {
//			return err
//		}
//	}
//	return nil
//}

func (cfg *Config) GetClusterConfig() *K3D.Cluster {
	return &K3D.Cluster{
		Name:  cfg.Name,
		Token: cfg.Token,
	}
}
