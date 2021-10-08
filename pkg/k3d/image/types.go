package image

import (
	"context"

	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/client"
)

type Registry interface {
	Upload() error
	List() error
}

// Config helps to store filtered images data that was loaded to k3d cluster.
type Config struct {
	Images       []string        `json:"images,omitempty"`
	Cluster      string          `json:"cluster,omitempty"`
	StoreTarBall bool            `json:"keep_tarball,omitempty"`
	StoredImages StoredImages    `json:"images_stored,omitempty"`
	Context      context.Context `json:"context,omitempty"`
	Config       client.Config   `json:"config,omitempty"`
}

// StoredImages holds a data of cluster to images mapping of loaded images.
type StoredImages struct {
	Cluster string   `json:"cluster,omitempty"`
	Images  []string `json:"images,omitempty"`
}

// TarBallData maps tarball stored to image.
type TarBallData struct {
	Image string `json:"image,omitempty"`
	Path  string `json:"path,omitempty"`
}
