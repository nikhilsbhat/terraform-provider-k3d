package node

import (
	"context"
	"fmt"

	"github.com/nikhilsbhat/terraform-provider-k3d/pkg/utils"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
	"github.com/thoas/go-funk"
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

func (cfg *Config) StartStopNode(ctx context.Context, runtime runtimes.Runtime) error {
	nodes, err := runtime.GetNodesByLabel(ctx, map[string]string{
		"k3d.cluster": cfg.ClusterAssociated,
		"k3d.role":    "agent",
	})
	if err != nil {
		return err
	}

	var filteredNodes []*K3D.Node

	if !cfg.All {
		filteredNodes = funk.Filter(nodes, func(node *K3D.Node) bool {
			return funk.Contains(cfg.Name, node.Name)
		}).([]*K3D.Node)
	}

	if len(filteredNodes) == 0 {
		return fmt.Errorf("nodes %v not found to start/stop them", cfg.Name)
	}

	if cfg.Action == utils.TerraformResourceStart {
		for _, filteredNode := range filteredNodes {
			if err = runtime.StartNode(ctx, filteredNode); err != nil {
				return err
			}
		}

		return nil
	}

	for _, filteredNode := range filteredNodes {
		if err = runtime.StopNode(ctx, filteredNode); err != nil {
			return err
		}
	}

	return nil
}
