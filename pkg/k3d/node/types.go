package node

import (
	"context"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
	"time"
)

const (
	// K3dClusterNameLabel is the label that holds cluster name in k3d node.
	K3dClusterNameLabel = "k3d.cluster"
)

type K3DNode interface {
	CreateNodeWithTimeout(context.Context, runtimes.Runtime, []*Config) error
	CreateNodes(context.Context, runtimes.Runtime, int) error
	DeleteNodesFromCluster(context.Context, runtimes.Runtime) error
	GetFilteredNodesFromCluster(context.Context, runtimes.Runtime) ([]*Config, error)
	GetFilteredNodes(context.Context, runtimes.Runtime) ([]*Config, error)
	GetNodesByLabels(context.Context, runtimes.Runtime) ([]*Config, error)
	GetNodeStatus(context.Context, runtimes.Runtime) ([]*Status, error)
	GetNodeFromConfig() *K3D.Node
	StartStopNode(context.Context, runtimes.Runtime) error
}

// Config stores filtered node data of k3d cluster.
type Config struct {
	Name                 []string               `json:"name,omitempty" mapstructure:"name"`
	Role                 string                 `json:"role,omitempty" mapstructure:"role"`
	ClusterAssociated    string                 `json:"cluster,omitempty" mapstructure:"cluster"`
	State                string                 `json:"state,omitempty" mapstructure:"state"`
	Created              string                 `json:"created,omitempty" mapstructure:"created"`
	Memory               string                 `json:"memory,omitempty" mapstructure:"memory"`
	Volumes              []string               `json:"volumes,omitempty" mapstructure:"volumes"`
	Networks             []string               `json:"networks,omitempty" mapstructure:"networks"`
	EnvironmentVariables []string               `json:"env,omitempty" mapstructure:"env"`
	Count                int                    `json:"count,omitempty" mapstructure:"count"`
	Image                string                 `json:"image,omitempty" mapstructure:"image"`
	PortMapping          map[string]interface{} `json:"port_mappings,omitempty" mapstructure:"port_mappings"`
	Timeout              time.Duration          `json:"timeout,omitempty" mapstructure:"timeout"`
	Wait                 bool                   `json:"wait,omitempty" mapstructure:"wait"`
	All                  bool                   `json:"all,omitempty" mapstructure:"all"`
	Labels               map[string]string      `json:"labels,omitempty" mapstructure:"labels"`
	Action               string                 `json:"action,omitempty" mapstructure:"action"`
}

// Status helps to store filtered node status of k3d cluster.
type Status struct {
	Node    string `json:"node,omitempty"`
	Cluster string `json:"cluster,omitempty"`
	Role    string `json:"role,omitempty"`
	State   string `json:"state,omitempty"`
	Running bool   `json:"running,omitempty"`
}
