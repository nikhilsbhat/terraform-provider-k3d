package k3d

import (
	"context"
	"fmt"

	"github.com/rancher/k3d/v4/pkg/tools"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

type K3Dimages struct {
	Images       []string     `json:"images,omitempty"`
	Cluster      string       `json:"cluster,omitempty"`
	StoreTarBall bool         `json:"keep_tarball,omitempty"`
	StoredImages StoredImages `json:"images_stored,omitempty"`
	Config       K3dConfig    `json:"config,omitempty"`
}

type StoredImages struct {
	Cluster string   `json:"cluster,omitempty"`
	Images  []string `json:"images,omitempty"`
}

type TarBallData struct {
	Image string `json:"image,omitempty"`
	Path  string `json:"path,omitempty"`
}

// StoreImages stores images in a specified clusters, also stores the tarball locally if feature is enabled.
func (client *K3Dimages) StoreImages() (*StoredImages, error) {
	loadImageOpts := K3D.ImageImportOpts{KeepTar: client.StoreTarBall}

	if err := tools.ImageImportIntoClusterMulti(context.Background(), client.Config.K3DRuntime, client.Images, Clusters(client.Cluster), loadImageOpts); err != nil {
		return nil, fmt.Errorf("failed to import image(s) into cluster '%s': %+v", client.Cluster, err)
	}

	return &StoredImages{
		Cluster: client.Cluster,
		Images:  client.Images,
	}, nil
}

func (client *StoredImages) GetImages() (images []string) {
	return client.Images
}

//func (client *StoredImages) GetClusters() (clusters []string) {
//
//}

func NewK3dImages() *K3Dimages {
	return &K3Dimages{}
}
