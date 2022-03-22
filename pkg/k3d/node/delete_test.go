package node

import (
	"context"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_DeleteNodesFromCluster(t *testing.T) {
	t.Run("should be able to delete the nodes", func(t *testing.T) {
		node := Config{
			ClusterAssociated: "k3s-default",
			Name:              []string{"test-node-from-terraform-0", "test-node-terraform-0"},
		}

		err := node.DeleteNodesFromCluster(context.Background(), runtimes.SelectedRuntime)
		assert.NoError(t, err)
	})
}
