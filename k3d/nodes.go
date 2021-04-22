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

func GetNodes(ctx context.Context, runtime runtimes.Runtime, cluster string) ([]*K3D.Node, error) {
	nodes, err := runtime.GetNodesByLabel(ctx, map[string]string{
		"k3d.cluster": cluster,
	})
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func GetNode(ctx context.Context, runtime runtimes.Runtime, name string) (*K3D.Node, error) {
	node, err := client.NodeGet(ctx, runtime, &K3D.Node{Name: name})
	if err != nil {
		return nil, err
	}
	return node, err
}

func GetNodesFromCluster(ctx context.Context, runtime runtimes.Runtime, cluster string) ([]*Node, error) {
	k3dNodes, err := GetNodes(ctx, runtime, cluster)
	if err != nil {
		return nil, err
	}
	filteredNodes := make([]*Node, 0)
	for _, node := range k3dNodes {
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
	return filteredNodes, err
}

func GetFilteredNodes(ctx context.Context, runtime runtimes.Runtime, nodes []string) ([]*Node, error) {
	k3dNodes := make([]*Node, 0)
	for _, currentNode := range nodes {
		node, err := GetNode(ctx, runtime, currentNode)
		if err != nil {
			return nil, err
		}
		k3dNodes = append(k3dNodes, &Node{
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

//func StopNodes(ctx context.Context, runtime runtimes.Runtime, cluster string, nodes []string) error {
//	nodesRaw, err := GetNodeRaw(ctx, runtime, cluster)
//	for _, node := range nodes {
//		nodeRaw, err := GetNodeRaw(ctx)
//		if err := runtime.StopNode(ctx, node); err != nil {
//			log.Fatalln(err)
//		}
//	}
//
//}
