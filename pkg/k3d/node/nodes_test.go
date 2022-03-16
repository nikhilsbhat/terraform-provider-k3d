package node

import (
	"context"
	"testing"

	"github.com/rancher/k3d/v5/pkg/runtimes"
	"github.com/stretchr/testify/assert"
)

func TestGetNodes(t *testing.T) {
	t.Run("should return the nodes of the specified cluster", func(t *testing.T) {
		nodes, err := Nodes(context.Background(), runtimes.Docker, "k3s-default")
		assert.NoError(t, err)
		assert.NotNil(t, nodes)
	})

	t.Run("should fail while fetching nodes of the specified cluster", func(t *testing.T) {
		nodes, err := Nodes(context.Background(), runtimes.Docker, "k3s-default-1")
		assert.Error(t, err, "unable to fetch nodes info as cluster k3s-default-1 not found")
		assert.Nil(t, nodes)
	})
}
