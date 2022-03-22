package registry

import (
	"context"
	k3dNode "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/node"
	"github.com/rancher/k3d/v5/pkg/runtimes"
)

type Registry interface {
	Create(context.Context, runtimes.Runtime) error
	Connect(context.Context, runtimes.Runtime) error
	Disconnect(context.Context, runtimes.Runtime) error
	Get(context.Context, runtimes.Runtime) ([]*k3dNode.Config, error)
}

// Config helps to store filtered registry data the present in selected runtime.
type Config struct {
	Name             []string          `json:"name,omitempty"`
	Image            string            `json:"image,omitempty"`
	Cluster          string            `json:"cluster,omitempty"`
	Protocol         string            `json:"protocol,omitempty"`
	Host             string            `json:"host,omitempty"`
	Port             string            `json:"port,omitempty"`
	Expose           map[string]string `json:"expose,omitempty"`
	UseProxy         bool              `json:"use_proxy,omitempty"`
	Proxy            map[string]string `json:"proxy,omitempty"`
	All              bool              `json:"all,omitempty"`
	ConnectToCluster bool              `json:"connect,omitempty"`
}
