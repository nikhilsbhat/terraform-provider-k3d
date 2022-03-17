package node

import (
	"context"
	"time"

	"github.com/docker/go-connections/nat"
	cluster2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
)

var (
	// K3dclusterNameLabel is the label that holds cluster name in k3d node.
	K3dclusterNameLabel = "k3d.cluster"
)

// Nodes fetches details of all available nodes in the specified cluster.
func Nodes(ctx context.Context, runtime runtimes.Runtime, cluster string) ([]*K3D.Node, error) {
	nodes, err := runtime.GetNodesByLabel(ctx, map[string]string{
		"cluster": cluster,
	})
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// Node fetches details of a specified node.
func Node(ctx context.Context, runtime runtimes.Runtime, name string) (*K3D.Node, error) {
	node, err := client.NodeGet(ctx, runtime, &K3D.Node{Name: name})
	if err != nil {
		return nil, err
	}
	return node, err
}

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

// CreateNode creates list node with specified configurations.
func CreateNode(ctx context.Context, runtime runtimes.Runtime,
	node []*K3D.Node, cluster *K3D.Cluster, options K3D.NodeCreateOpts) error {
	if err := client.NodeAddToClusterMulti(ctx, runtime, node, cluster, options); err != nil {
		return err
	}
	return nil
}

// DeleteNodes deletes the specified node.
func DeleteNodes(ctx context.Context, runtime runtimes.Runtime,
	node *K3D.Node, options K3D.NodeDeleteOpts) error {
	if err := client.NodeDelete(ctx, runtime, node, options); err != nil {
		return err
	}
	return nil
}

// GetFilteredNodesFromCluster returns the fetched all nodes from a specified cluster with list of *K3Node type.
func GetFilteredNodesFromCluster(ctx context.Context, runtime runtimes.Runtime, cluster string) ([]*K3Node, error) {
	k3dNodes, err := Nodes(ctx, runtime, cluster)
	if err != nil {
		return nil, err
	}
	filteredNodes := make([]*K3Node, 0)
	for _, node := range k3dNodes {
		filteredNodes = append(filteredNodes, &K3Node{
			Name:                 node.Name,
			Role:                 string(node.Role),
			ClusterAssociated:    node.RuntimeLabels[K3dclusterNameLabel],
			State:                node.State.Status,
			Created:              node.Created,
			Volumes:              node.Volumes,
			Networks:             node.Networks,
			EnvironmentVariables: node.Env,
		})
	}
	return filteredNodes, err
}

// GetFilteredNodes returns the fetched list of specified nodes from specified cluster with list of *K3Node type.
func GetFilteredNodes(ctx context.Context, runtime runtimes.Runtime, nodes []string) ([]*K3Node, error) {
	k3dNodes := make([]*K3Node, 0)
	for _, currentNode := range nodes {
		node, err := Node(ctx, runtime, currentNode)
		if err != nil {
			return nil, err
		}
		k3dNodes = append(k3dNodes, &K3Node{
			Name:                 node.Name,
			Role:                 string(node.Role),
			ClusterAssociated:    node.RuntimeLabels[K3dclusterNameLabel],
			State:                node.State.Status,
			Created:              node.Created,
			Volumes:              node.Volumes,
			Networks:             node.Networks,
			EnvironmentVariables: node.Env,
		})
	}
	return k3dNodes, nil
}

// GetNodes returns the list of all nodes available in the specified runtime.
func GetNodes(ctx context.Context, runtime runtimes.Runtime) ([]*K3Node, error) {
	nodes, err := client.NodeList(ctx, runtime)
	if err != nil {
		return nil, err
	}
	k3dNodes := make([]*K3Node, 0)
	for _, node := range nodes {
		k3dNodes = append(k3dNodes, &K3Node{
			Name:                 node.Name,
			Role:                 string(node.Role),
			ClusterAssociated:    node.RuntimeLabels[K3dclusterNameLabel],
			State:                node.State.Status,
			Created:              node.Created,
			Memory:               node.Memory,
			Volumes:              node.Volumes,
			Networks:             node.Networks,
			EnvironmentVariables: node.Env,
			PortMapping:          getPortMaps(node.Ports),
		})
	}
	return k3dNodes, nil
}

// GetNodesByLabels gets the nodes that matches with the specified label.
func GetNodesByLabels(ctx context.Context, runtime runtimes.Runtime, label map[string]string) ([]*K3Node, error) {
	k3dNodes, err := runtime.GetNodesByLabel(ctx, label)
	if err != nil {
		return nil, err
	}
	filteredNodes := make([]*K3Node, 0)
	for _, node := range k3dNodes {
		filteredNodes = append(filteredNodes, &K3Node{
			Name:                 node.Name,
			Role:                 string(node.Role),
			ClusterAssociated:    node.RuntimeLabels[K3dclusterNameLabel],
			State:                node.State.Status,
			Created:              node.Created,
			Memory:               node.Memory,
			Volumes:              node.Volumes,
			Networks:             node.Networks,
			EnvironmentVariables: node.Env,
			Image:                node.Image,
			// dropping PortMapping as terraform schema format is yet to be figured.
			//PortMapping:          getPortMaps(node.Ports),
		})
	}
	return filteredNodes, err
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

// GetNode returns K3D.Node equivalent for an stance of K3Node.
func (c *K3Node) GetNode() *K3D.Node {
	return &K3D.Node{
		Name: c.Name,
		Role: K3D.NodeRoles[c.Role],
		RuntimeLabels: map[string]string{
			K3D.LabelRole:           c.Role,
			utils.TerraformK3dLabel: c.Created,
		},
		Image:   c.Image,
		Memory:  c.Memory,
		Restart: true,
	}
}

// CreateNodeWithTimeout creates node by setting timeouts as per input.
func CreateNodeWithTimeout(ctx context.Context, runtime runtimes.Runtime,
	cluster string, nodes []*K3Node, wait bool, timeout time.Duration) error {
	var nodeCreatOpts K3D.NodeCreateOpts
	if wait {
		nodeCreatOpts = K3D.NodeCreateOpts{Wait: wait, Timeout: timeout}
	}
	clusterFetched, err := cluster2.GetCluster(ctx, runtime, cluster)
	if err != nil {
		return err
	}
	k3dNodes := make([]*K3D.Node, 0)
	for _, node := range nodes {
		k3dNodes = append(k3dNodes, node.GetNode())
	}
	return CreateNode(ctx, runtime, k3dNodes, clusterFetched, nodeCreatOpts)
}

// DeleteNodesFromCluster deletes the specified node.
func DeleteNodesFromCluster(ctx context.Context, runtime runtimes.Runtime, node *K3D.Node) error {
	deleteOps := K3D.NodeDeleteOpts{
		SkipLBUpdate: false,
	}
	return DeleteNodes(ctx, runtime, node, deleteOps)
}

func getPortMaps(p nat.PortMap) map[string]interface{} {
	portM := make(map[string]interface{})
	for key, value := range p {
		portM[string(key)] = value
	}
	return portM
}
