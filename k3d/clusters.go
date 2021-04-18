package k3d

import (
	"context"
	"fmt"

	"github.com/rancher/k3d/v4/pkg/client"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

func GetCluster(ctx context.Context, runtime runtimes.Runtime, cluster string) (*K3D.Cluster, error) {
	clusterConfig, err := client.ClusterGet(ctx, runtime, &K3D.Cluster{Name: cluster})
	if err != nil {
		return nil, err
	}
	return clusterConfig, nil
}

func GetClusters(ctx context.Context, runtime runtimes.Runtime, clusters []string) ([]*K3D.Cluster, error) {
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
