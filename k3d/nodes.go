package k3d

import (
	"context"
	"time"

	"github.com/nikhilsbhat/terraform-provider-rancherk3d/utils"
	"github.com/rancher/k3d/v4/pkg/client"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

var (
	K3dclusterNameLabel = "k3d.cluster"
)

type K3dNode interface {
	DeleteNodes()
}

func GetNodes(ctx context.Context, runtime runtimes.Runtime, cluster string) ([]*K3D.Node, error) {
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

func CreateNode(ctx context.Context, runtime runtimes.Runtime, node []*K3D.Node, cluster *K3D.Cluster, options K3D.NodeCreateOpts) error {
	if err := client.NodeAddToClusterMulti(ctx, runtime, node, cluster, options); err != nil {
		return err
	}
	return nil
}

func DeleteNodes(ctx context.Context, runtime runtimes.Runtime, node *K3D.Node, options K3D.NodeDeleteOpts) error {
	if err := client.NodeDelete(ctx, runtime, node, options); err != nil {
		return err
	}
	return nil
}

func GetFilteredNodesFromCluster(ctx context.Context, runtime runtimes.Runtime, cluster string) ([]*K3DNode, error) {
	k3dNodes, err := GetNodes(ctx, runtime, cluster)
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

func GetNodesByLabels(ctx context.Context, runtime runtimes.Runtime, label map[string]string) ([]*K3DNode, error) {
	k3dNodes, err := runtime.GetNodesByLabel(ctx, label)
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
	nodesRaw, err := GetNodes(ctx, runtime, cluster)
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
	nodesRaw, err := GetNodes(ctx, runtime, cluster)
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

func (c *K3DNode) GetNode() *K3D.Node {
	return &K3D.Node{
		Name: c.Name,
		Role: K3D.NodeRoles[c.Role],
		Labels: map[string]string{
			K3D.LabelRole:           c.Role,
			utils.TerraformK3dLabel: c.Created,
		},
		Image:   c.Image,
		Memory:  c.Memory,
		Restart: true,
	}
}

func CreateNodeWithTimeout(ctx context.Context, runtime runtimes.Runtime, cluster string, nodes []*K3DNode, wait bool, timeout time.Duration) error {
	var nodeCreatOpts K3D.NodeCreateOpts
	if wait {
		nodeCreatOpts = K3D.NodeCreateOpts{Wait: wait, Timeout: timeout}
	}
	clusterFetched, err := GetCluster(ctx, runtime, cluster)
	if err != nil {
		return err
	}
	k3dNodes := make([]*K3D.Node, 0)
	for _, node := range nodes {
		k3dNodes = append(k3dNodes, node.GetNode())
	}
	return CreateNode(ctx, runtime, k3dNodes, clusterFetched, nodeCreatOpts)
}

func DeletNodesFromCluster(ctx context.Context, runtime runtimes.Runtime, node *K3D.Node) error {
	deleteOps := K3D.NodeDeleteOpts{
		SkipLBUpdate: false,
	}
	return DeleteNodes(ctx, runtime, node, deleteOps)
}
