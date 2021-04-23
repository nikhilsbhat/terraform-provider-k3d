package k3d

import (
	"context"

	"github.com/rancher/k3d/v4/pkg/client"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

var (
	K3dclusterNameLabel = "k3d.cluster"
)

func Nodes(ctx context.Context, runtime runtimes.Runtime, cluster string) ([]*K3D.Node, error) {
	nodes, err := runtime.GetNodesByLabel(ctx, map[string]string{
		"k3d.cluster": cluster,
	})
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func Node(ctx context.Context, runtime runtimes.Runtime, name string) (*K3D.Node, error) {
	node, err := client.NodeGet(ctx, runtime, &K3D.Node{Name: name})
	if err != nil {
		return nil, err
	}
	return node, err
}

func FilteredNodes(ctx context.Context, runtime runtimes.Runtime, nodes []string) ([]*K3D.Node, error) {
	k3dNodes, err := client.NodeList(ctx, runtime)
	if err != nil {
		return nil, err
	}
	filteredNodes := make([]*K3D.Node, 0)
	for _, k3dNode := range k3dNodes {
		for _, node := range nodes {
			if k3dNode.Name == node {
				filteredNodes = append(filteredNodes, k3dNode)
			}
		}
	}
	return filteredNodes, nil
}

func StopNode(ctx context.Context, runtime runtimes.Runtime, node *K3D.Node) error {
	if err := runtime.StopNode(ctx, node); err != nil {
		return err
	}
	return nil
}

func StartNode(ctx context.Context, runtime runtimes.Runtime, node *K3D.Node) error {
	if err := runtime.StartNode(ctx, node); err != nil {
		return err
	}
	return nil
}

func GetFilteredNodesFromCluster(ctx context.Context, runtime runtimes.Runtime, cluster string) ([]*K3DNode, error) {
	k3dNodes, err := Nodes(ctx, runtime, cluster)
	if err != nil {
		return nil, err
	}
	filteredNodes := make([]*K3DNode, 0)
	for _, node := range k3dNodes {
		filteredNodes = append(filteredNodes, &K3DNode{
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
	return filteredNodes, err
}

func GetFilteredNodes(ctx context.Context, runtime runtimes.Runtime, nodes []string) ([]*K3DNode, error) {
	k3dNodes := make([]*K3DNode, 0)
	for _, currentNode := range nodes {
		node, err := Node(ctx, runtime, currentNode)
		if err != nil {
			return nil, err
		}
		k3dNodes = append(k3dNodes, &K3DNode{
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
	return k3dNodes, nil
}

func StopNodes(ctx context.Context, runtime runtimes.Runtime, nodes []string) error {
	nodesRaw, err := FilteredNodes(ctx, runtime, nodes)
	if err != nil {
		return err
	}
	for _, node := range nodesRaw {
		if err := StopNode(ctx, runtime, node); err != nil {
			return err
		}
	}
	return nil
}

func StopNodesFromCluster(ctx context.Context, runtime runtimes.Runtime, cluster string) error {
	nodesRaw, err := Nodes(ctx, runtime, cluster)
	if err != nil {
		return err
	}
	for _, node := range nodesRaw {
		if err := StopNode(ctx, runtime, node); err != nil {
			return err
		}
	}
	return nil
}

func StartNodes(ctx context.Context, runtime runtimes.Runtime, nodes []string) error {
	nodesRaw, err := FilteredNodes(ctx, runtime, nodes)
	if err != nil {
		return err
	}
	for _, node := range nodesRaw {
		if err := StartNode(ctx, runtime, node); err != nil {
			return err
		}

	}
	return nil
}

func StartNodesFromCluster(ctx context.Context, runtime runtimes.Runtime, cluster string) error {
	nodesRaw, err := Nodes(ctx, runtime, cluster)
	if err != nil {
		return err
	}
	for _, node := range nodesRaw {
		if err := StartNode(ctx, runtime, node); err != nil {
			return err
		}
	}
	return nil
}
