package config

import k3dCluster "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"

// Config helps in storing kubeconfig information.
type Config struct {
	Cluster []*k3dCluster.Config `json:"cluster,omitempty" mapstructure:"cluster"`
	Encode  bool                 `json:"encode,omitempty" mapstructure:"encode"`
	All     bool                 `json:"all,omitempty" mapstructure:"all"`
}
