package node

import (
	"context"
	"fmt"
	"slices"
	"strings"

	terraformErrors "github.com/nikhilsbhat/terraform-provider-k3d/pkg/errors"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
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

	filteredNodes := make([]*K3D.Node, 0, len(nodes))
	for _, node := range nodes {
		if slices.Contains(cfg.Name, node.Name) {
			filteredNodes = append(filteredNodes, node)
		}
	}

	deleteOps := K3D.NodeDeleteOpts{
		SkipLBUpdate: false,
	}

	errors := make([]string, 0, len(filteredNodes))

	for _, filteredNode := range filteredNodes {
		if delErr := client.NodeDelete(ctx, runtime, filteredNode, deleteOps); delErr != nil {
			errors = append(errors, delErr.Error())
		}
	}

	if len(errors) != 0 {
		return fmt.Errorf("%w: %s", terraformErrors.ErrDeleteNodesFailed, strings.Join(errors, "\n"))
	}

	return nil
}
