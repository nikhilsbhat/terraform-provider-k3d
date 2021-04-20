package k3d

import (
	"context"
	"fmt"
	"log"

	"github.com/rancher/k3d/v4/pkg/client"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

var (
	K3dclusterNameLabel = "k3d.cluster"
)

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

func GetNodes(ctx context.Context, runtime runtimes.Runtime, cluster string) ([]*Node, error) {
	nodes, err := client.NodeList(ctx, runtime)
	if err != nil {
		return nil, err
	}
	filteredNodes := make([]*Node, 0)
	for _, node := range nodes {
		log.Print(node.Labels)
		if node.Labels["k3d.cluster"] == cluster {
			filteredNodes = append(filteredNodes, &Node{
				Name:                 node.Name,
				Role:                 string(node.Role),
				ClusterAssociated:    node.Labels[K3dclusterNameLabel],
				State:                node.State.Status,
				Created:              node.Created,
				Volumes:              node.Volumes,
				Networks:             node.Networks,
				EnvironmentVariables: node.Env,
			})
		}
	}
	if len(filteredNodes) == 0 {
		return nil, fmt.Errorf("unable to fetch nodes info as cluster %s not found", cluster)
	}
	return filteredNodes, nil
}

func GetNode(ctx context.Context, runtime runtimes.Runtime, name, cluster string) (*Node, error) {
	node, err := client.NodeGet(ctx, runtime, &K3D.Node{Name: name})
	if err != nil {
		return nil, err
	}

	return &Node{
		Name:                 node.Name,
		Role:                 string(node.Role),
		ClusterAssociated:    node.Labels[K3dclusterNameLabel],
		State:                node.State.Status,
		Created:              node.Created,
		Volumes:              node.Volumes,
		Networks:             node.Networks,
		EnvironmentVariables: node.Env,
	}, nil
}
