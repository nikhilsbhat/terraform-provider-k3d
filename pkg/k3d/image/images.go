package image

import (
	"context"
	"fmt"

	cluster2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
)

// StoreImagesToCluster stores images in a specified clusters, also stores the tarball locally if feature is enabled.
func StoreImagesToCluster(ctx context.Context, runtime runtimes.Runtime,
	images []string, cluster string, storeTarball bool) error {
	loadImageOpts := K3D.ImageImportOpts{KeepTar: storeTarball}

	clusterCfg := cluster2.Config{}
	clusters, err := clusterCfg.GetClusters(ctx, runtime, []string{cluster})
	if err != nil {
		return err
	}
	retrievedCluster := clusters[0].GetClusterConfig()
	if err = client.ImageImportIntoClusterMulti(ctx, runtime, images, retrievedCluster, loadImageOpts); err != nil {
		return fmt.Errorf("failed to import image(s) into cluster '%s': %+v", retrievedCluster.Name, err)
	}
	return nil
}

// StoreImagesToClusters stores images to all specified clusters.
func StoreImagesToClusters(ctx context.Context, runtime runtimes.Runtime,
	images []string, storeTarball bool) error {
	loadImageOpts := K3D.ImageImportOpts{KeepTar: storeTarball}

	retrievedClusters, err := cluster2.GetClusters(ctx, runtime)
	if err != nil {
		return err
	}

	for _, retrievedCluster := range retrievedClusters {
		if err = client.ImageImportIntoClusterMulti(ctx, runtime, images, retrievedCluster, loadImageOpts); err != nil {
			return fmt.Errorf("failed to import image(s) into cluster '%s': %+v", retrievedCluster.Name, err)
		}
	}
	return nil
}

// GetImagesLoadedCluster returns list of images loaded to the cluster.
func GetImagesLoadedCluster(ctx context.Context, runtime runtimes.Runtime,
	images []string, cluster string) ([]*StoredImages, error) {

	clusterCfg := cluster2.Config{}
	clusters, err := clusterCfg.GetClusters(ctx, runtime, []string{cluster})
	if err != nil {
		return nil, err
	}
	retrievedCluster := clusters[0].GetClusterConfig()

	return []*StoredImages{{
		Cluster: retrievedCluster.Name,
		Images:  images,
	}}, nil
}

// GetImagesLoadedClusters returns list of images loaded to the clusters.
func GetImagesLoadedClusters(ctx context.Context, runtime runtimes.Runtime,
	images []string) ([]*StoredImages, error) {
	retrievedClusters, err := cluster2.GetClusters(ctx, runtime)
	if err != nil {
		return nil, err
	}

	storedImages := make([]*StoredImages, 0)
	for _, retrievedCluster := range retrievedClusters {
		storedImages = append(storedImages, &StoredImages{
			Cluster: retrievedCluster.Name,
			Images:  images,
		})
	}
	return storedImages, nil
}

func NewK3dImages() *Config {
	return &Config{}
}
