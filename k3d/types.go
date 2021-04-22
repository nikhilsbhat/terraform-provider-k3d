package k3d

import "context"

type K3Dimages struct {
	Images       []string        `json:"images,omitempty"`
	Cluster      string          `json:"cluster,omitempty"`
	StoreTarBall bool            `json:"keep_tarball,omitempty"`
	StoredImages StoredImages    `json:"images_stored,omitempty"`
	Context      context.Context `json:"context,omitempty"`
	Config       K3dConfig       `json:"config,omitempty"`
}

type StoredImages struct {
	Cluster string   `json:"cluster,omitempty"`
	Images  []string `json:"images,omitempty"`
}

type TarBallData struct {
	Image string `json:"image,omitempty"`
	Path  string `json:"path,omitempty"`
}

type Node struct {
	Name                 string   `json:"name,omitempty"`
	Role                 string   `json:"role,omitempty"`
	ClusterAssociated    string   `json:"cluster,omitempty"`
	State                string   `json:"state,omitempty"`
	Created              string   `json:"created,omitempty"`
	Volumes              []string `json:"volumes,omitempty"`
	Networks             []string `json:"networks,omitempty"`
	EnvironmentVariables []string `json:"env,omitempty"`
}
