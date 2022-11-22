package registry

import (
	"context"

	k3dNode "github.com/nikhilsbhat/terraform-provider-k3d/pkg/k3d/node"
	"github.com/rancher/k3d/v5/pkg/runtimes"
)

type Registry interface {
	Create(context.Context, runtimes.Runtime) error
	Connect(context.Context, runtimes.Runtime) error
	Disconnect(context.Context, runtimes.Runtime) error
	Get(context.Context, runtimes.Runtime) ([]*k3dNode.Config, error)
}

// Config helps to store filtered registry data the present in selected runtime.
//
//nolint:maligned
type Config struct {
	Name             []string          `json:"name,omitempty" mapstructure:"name"`
	Image            string            `json:"image,omitempty"  mapstructure:"image"`
	Cluster          string            `json:"cluster,omitempty" mapstructure:"cluster"`
	Protocol         string            `json:"protocol,omitempty" mapstructure:"protocol"`
	Host             string            `json:"host,omitempty" mapstructure:"host"`
	Port             string            `json:"port,omitempty" mapstructure:"port"`
	Expose           map[string]string `json:"expose,omitempty" mapstructure:"expose"`
	UseProxy         bool              `json:"use_proxy,omitempty" mapstructure:"use_proxy"`
	Proxy            map[string]string `json:"proxy,omitempty" mapstructure:"proxy"`
	All              bool              `json:"all,omitempty" mapstructure:"all"`
	ConnectToCluster bool              `json:"connect,omitempty" mapstructure:"connect"`
}
