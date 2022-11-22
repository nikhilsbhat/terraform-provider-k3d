package image

import (
	"context"
	"fmt"
	"strings"

	cluster2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
)

// Upload uploads images to a specified clusters, also stores the tarball locally if feature is enabled.
func (image *Config) Upload(ctx context.Context, runtime runtimes.Runtime) error {
	loadImageOpts := K3D.ImageImportOpts{KeepTar: image.StoreTarBall}

	clusterCfg := cluster2.Config{
		All: image.All,
	}

	var clusters []*K3D.Cluster
	k3dClusters, err := clusterCfg.GetClusters(ctx, runtime, []string{image.Cluster})
	if err != nil {
		return err
	}

	for _, k3dCluster := range k3dClusters {
		clusters = append(clusters, k3dCluster.GetClusterConfig())
	}

	var errors []string
	for _, cluster := range clusters {
		if err = client.ImageImportIntoClusterMulti(ctx, runtime, image.Images, cluster, loadImageOpts); err != nil {
			errors = append(errors, fmt.Sprintf("failed to import image(s) into cluster '%s': %+v", cluster.Name, err))
		}
	}

	if len(errors) != 0 {
		return fmt.Errorf("importing images to clusters errored with: \n%s", strings.Join(errors, "\n"))
	}
	return nil
}

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
