package image

import (
	"context"
	"fmt"
	"strings"

	terraformErrors "github.com/nikhilsbhat/terraform-provider-k3d/pkg/errors"
	cluster2 "github.com/nikhilsbhat/terraform-provider-k3d/pkg/k3d/cluster"
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

	k3dClusters, err := clusterCfg.GetClusters(ctx, runtime, []string{image.Cluster})
	if err != nil {
		return err
	}

	errors := make([]string, 0, len(k3dClusters))

	for _, k3dCluster := range k3dClusters {
		cluster := k3dCluster.GetClusterConfig()
		if err = client.ImageImportIntoClusterMulti(ctx, runtime, image.Images, cluster, loadImageOpts); err != nil {
			errors = append(errors, fmt.Sprintf("failed to import image(s) into cluster '%s': %+v", cluster.Name, err))
		}
	}

	if len(errors) != 0 {
		return fmt.Errorf("%w: \n%s", terraformErrors.ErrImportImagesFailed, strings.Join(errors, "\n"))
	}

	return nil
}
