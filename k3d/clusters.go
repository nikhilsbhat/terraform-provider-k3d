package k3d

import (
	"context"
	"fmt"
	"log"

	"github.com/rancher/k3d/v4/pkg/client"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

func GetCluster(ctx context.Context, runtime runtimes.Runtime,
	cluster string) (*K3D.Cluster, error) {
	clusterConfig, err := client.ClusterGet(ctx, runtime, &K3D.Cluster{Name: cluster})
	if err != nil {
		return nil, err
	}
	return clusterConfig, nil
}

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

func GetClusters(ctx context.Context, runtime runtimes.Runtime) ([]*K3D.Cluster, error) {
	clustersList, err := client.ClusterList(ctx, runtime)
	if err != nil {
		return nil, err
	}
	return clustersList, nil
}

func StartClusters(ctx context.Context, runtime runtimes.Runtime,
	clusters []*K3D.Cluster, options K3D.ClusterStartOpts) error {
	for _, cluster := range clusters {
		if err := client.ClusterStart(ctx, runtime, cluster, options); err != nil {
			log.Fatalln(err)
		}
	}
	return nil
}

func StopClusters(ctx context.Context, runtime runtimes.Runtime,
	clusters []*K3D.Cluster) error {
	for _, cluster := range clusters {
		if err := client.ClusterStop(ctx, runtime, cluster); err != nil {
			log.Fatalln(err)
		}
	}
	return nil
}
