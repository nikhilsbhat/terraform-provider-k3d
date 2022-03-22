package node

import (
	"context"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
)

// FilteredNodes fetches details of specified list of nodes.
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

// StopNode stops the specified node.
func StopNode(ctx context.Context, runtime runtimes.Runtime, node *K3D.Node) error {
	if err := runtime.StopNode(ctx, node); err != nil {
		return err
	}
	return nil
}

// StartNode starts the specified node.
func StartNode(ctx context.Context, runtime runtimes.Runtime, node *K3D.Node) error {
	if err := runtime.StartNode(ctx, node); err != nil {
		return err
	}
	return nil
}

// StopNodes stops all specified nodes.
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

// StopNodesFromCluster stops all available nodes from a specified cluster.
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

// StartNodes starts all specified of nodes.
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

// StartNodesFromCluster starts all available nodes from a specified cluster.
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
