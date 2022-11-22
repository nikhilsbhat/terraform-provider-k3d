package cluster

import (
	"context"
	"testing"

	"github.com/rancher/k3d/v5/pkg/runtimes"
	"github.com/stretchr/testify/assert"
)

func TestConfig_StartStopCluster(t *testing.T) {
	t.Run("should stop selected cluster", func(t *testing.T) {
		cfg := Config{
			All:    false,
			Action: "stop",
		}

		clusters := []string{"test"}

		err := cfg.StartStopCluster(context.Background(), runtimes.SelectedRuntime, clusters)
		assert.NoError(t, err)

		updatedCluster, err := cfg.GetClusters(context.Background(), runtimes.SelectedRuntime, clusters)
		assert.NoError(t, err)
		assert.Equal(t, updatedCluster[0].ServersCount, 0)
	})
}
