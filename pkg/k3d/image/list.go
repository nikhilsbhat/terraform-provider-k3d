package image

import (
	"context"

	cluster2 "github.com/nikhilsbhat/terraform-provider-k3d/pkg/k3d/cluster"
	"github.com/rancher/k3d/v5/pkg/runtimes"
)

// List returns list of images loaded to the clusters.
func (image *Config) List(ctx context.Context, runtime runtimes.Runtime) ([]*StoredImages, error) {
	clusterCfg := cluster2.Config{
		All: image.All,
	}

	retrievedClusters, err := clusterCfg.GetClusters(ctx, runtime, []string{image.Cluster})
	if err != nil {
		return nil, err
	}

	storedImages := make([]*StoredImages, 0)
	for _, retrievedCluster := range retrievedClusters {
		storedImages = append(storedImages, &StoredImages{
			Cluster: retrievedCluster.Name,
			Images:  image.Images,
		})
	}

	return storedImages, nil
}
