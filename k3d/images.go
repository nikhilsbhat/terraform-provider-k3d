package k3d

import (
	"fmt"
	"log"

	"github.com/rancher/k3d/v4/pkg/tools"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

// StoreImages stores images in a specified clusters, also stores the tarball locally if feature is enabled.
func (client *K3Dimages) StoreImages() (*StoredImages, error) {
	loadImageOpts := K3D.ImageImportOpts{KeepTar: client.StoreTarBall}

	cluster, err := GetCluster(client.Context, client.Config.K3DRuntime, client.Cluster)
	if err != nil {
		return nil, err
	}

	log.Printf("loading images %v to cluster %s", client.Images, cluster.Name)
	if err := tools.ImageImportIntoClusterMulti(client.Context, client.Config.K3DRuntime, client.Images, cluster, loadImageOpts); err != nil {
		return nil, fmt.Errorf("failed to import image(s) into cluster '%s': %+v", cluster.Name, err)
	}
	return &StoredImages{
		Cluster: cluster.Name,
		Images:  client.Images,
	}, nil
}

func (client *StoredImages) GetImages() (images []string) {
	return client.Images
}

func NewK3dImages() *K3Dimages {
	return &K3Dimages{}
}
