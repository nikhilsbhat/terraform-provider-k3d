package node

import (
	"context"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_GetFilteredNodesFromCluster(t *testing.T) {
	t.Run("should be able to fetch the nodes that matches the labels", func(t *testing.T) {
		cfg := Config{
			ClusterAssociated: "k3s-default",
			Name:              []string{"test-node-from-terraform"},
		}

		nodes, err := cfg.GetFilteredNodesFromCluster(context.Background(), runtimes.SelectedRuntime)
		assert.NoError(t, err)
		assert.Equal(t, nodes[0].Name, []string{"test-node-from-terraform"})
		assert.Equal(t, 1, len(nodes))
	})
}

func TestConfig_GetNodeStatus(t *testing.T) {
	t.Run("should be able to get status of selected nodes", func(t *testing.T) {
		cfg := Config{
			ClusterAssociated: "k3s-default",
			Name:              []string{"test-node-from-terraform-0"},
		}
		nodes, err := cfg.GetNodeStatus(context.Background(), runtimes.SelectedRuntime)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(nodes))
	})
}
