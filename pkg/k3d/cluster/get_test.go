package cluster_test

import (
	"context"
	"testing"

	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
)

func TestMAP(t *testing.T) {
	t.Run("", func(t *testing.T) {
		nodes := []*K3D.Node{
			{Name: "node1"},
			{Name: "node2"},
			{Name: "node3"},
		}

		expected := []string{"node1", "node2", "node3"}
		names := funk.Get(nodes, "Name")
		assert.Equal(t, expected, names)
	})
}

func TestConfig_GetClusters(t *testing.T) {
	clusters := []string{"k3s-default"}
	expected := []*cluster.Config{
		{
			Name:            "k3s-default",
			Nodes:           []string{"k3d-k3s-default-serverlb", "k3d-k3s-default-server-0"},
			Network:         "k3d-k3s-default",
			Token:           "lQRTKTdNIRzjDjfVzzDS",
			ServersCount:    1,
			ServersRunning:  1,
			AgentsCount:     0,
			AgentsRunning:   0,
			ImageVolume:     "k3d-k3s-default-images",
			HasLoadBalancer: true,
			All:             false,
		},
	}

	t.Run("should be able to fetch the filtered cluster responses", func(t *testing.T) {
		cfg := cluster.Config{All: false}
		clusters, err := cfg.GetClusters(context.Background(), runtimes.SelectedRuntime, clusters)
		assert.Nil(t, err)
		assert.ElementsMatch(t, expected, clusters)
	})

	t.Run("should be able to fetch all the available clusters", func(t *testing.T) {
		cfg := cluster.Config{All: true}
		expected = append(expected, &cluster.Config{
			Name:            "test",
			Nodes:           []string{"k3d-test-serverlb", "k3d-test-server-0"},
			Network:         "k3d-test",
			Token:           "rMNfPZSdxGxGFiykcJuQ",
			ServersCount:    1,
			ServersRunning:  1,
			AgentsCount:     0,
			AgentsRunning:   0,
			ImageVolume:     "k3d-test-images",
			HasLoadBalancer: true,
			All:             false,
		})
		clusters, err := cfg.GetClusters(context.Background(), runtimes.SelectedRuntime, clusters)
		assert.Nil(t, err)
		assert.ElementsMatch(t, expected, clusters)
	})
}
