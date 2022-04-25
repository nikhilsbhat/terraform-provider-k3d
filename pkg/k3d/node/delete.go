package node

import (
	"context"
	"fmt"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
	"github.com/thoas/go-funk"
	"strings"
)

// DeleteNodesFromCluster deletes the specified node.
func (cfg *Config) DeleteNodesFromCluster(ctx context.Context, runtime runtimes.Runtime) error {
	nodeLabel := map[string]string{
		"k3d.cluster": cfg.ClusterAssociated,
	}

	nodes, err := runtime.GetNodesByLabel(ctx, nodeLabel)
	if err != nil {
		return err
	}

	filteredNodes := funk.Filter(nodes, func(node *K3D.Node) bool {
		if funk.Contains(cfg.Name, node.Name) {
			return true
		}
		return false
	}).([]*K3D.Node)

	deleteOps := K3D.NodeDeleteOpts{
		SkipLBUpdate: false,
	}

	var errors []string
	for _, filteredNode := range filteredNodes {
		if delErr := client.NodeDelete(ctx, runtime, filteredNode, deleteOps); delErr != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) != 0 {
		return fmt.Errorf("deleting nodes failed with: %s", strings.Join(errors, "\n"))
	}
	return nil
}
