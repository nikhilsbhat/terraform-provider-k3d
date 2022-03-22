package node

import (
	"context"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	"time"
)

const (
	// K3dClusterNameLabel is the label that holds cluster name in k3d node.
	K3dClusterNameLabel = "k3d.cluster"
)

type K3DNode interface {
	CreateNodeWithTimeout(context.Context, runtimes.Runtime, []*Config) error
	CreateNodes(context.Context, runtimes.Runtime, int) error
	GetNodesByLabels(context.Context, runtimes.Runtime) ([]*Config, error)
}

// Config stores filtered node data of k3d cluster.
type Config struct {
	Name                 []string               `json:"name,omitempty"`
	Role                 string                 `json:"role,omitempty"`
	ClusterAssociated    string                 `json:"cluster,omitempty"`
	State                string                 `json:"state,omitempty"`
	Created              string                 `json:"created,omitempty"`
	Memory               string                 `json:"memory,omitempty"`
	Volumes              []string               `json:"volumes,omitempty"`
	Networks             []string               `json:"networks,omitempty"`
	EnvironmentVariables []string               `json:"env,omitempty"`
	Count                int                    `json:"count,omitempty"`
	Image                string                 `json:"image,omitempty"`
	PortMapping          map[string]interface{} `json:"port_mappings,omitempty"`
	Timeout              time.Duration          `json:"timeout,omitempty"`
	Wait                 bool                   `json:"wait,omitempty"`
	All                  bool                   `json:"all,omitempty"`
	Labels               map[string]string      `json:"labels,omitempty"`
}

// Status helps to store filtered node status of k3d cluster.
type Status struct {
	Node    string `json:"node,omitempty"`
	Cluster string `json:"cluster,omitempty"`
	Role    string `json:"role,omitempty"`
	State   string `json:"state,omitempty"`
	Running bool   `json:"running,omitempty"`
}

func NewConfig() *Config {
	return &Config{}
}
