package utils

import (
	"testing"

	"github.com/nikhilsbhat/terraform-provider-rancherk3d/k3d"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Run("should return value in Type map", func(t *testing.T) {
		node := []*k3d.Node{
			{
				Name:              "test",
				ClusterAssociated: "test-cluster",
				Role:              "load-balancer",
				State:             "stopped",
			},
			{
				Name:              "test-1",
				ClusterAssociated: "test-cluster-1",
				Role:              "load-balancer",
				State:             "stopped",
			},
		}

		expected := []map[string]interface{}{
			{
				"cluster": "test-cluster",
				"name":    "test",
				"role":    "load-balancer",
				"state":   "stopped",
			},
			{
				"cluster": "test-cluster-1",
				"name":    "test-1",
				"role":    "load-balancer",
				"state":   "stopped",
			},
		}
		actual, err := MapSlice(node)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
