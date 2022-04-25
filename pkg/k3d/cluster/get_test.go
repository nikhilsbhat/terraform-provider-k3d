package cluster

import (
	"context"
	"github.com/thoas/go-funk"
	"reflect"
	"testing"

	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestClusters(t *testing.T) {
	t.Run("should return K3D.Config successfully", func(t *testing.T) {
		expected := K3D.Cluster{
			Name:  "k3s-default",
			Nodes: []*K3D.Node{},
		}

		actual, err := GetCluster(context.Background(), runtimes.Docker, "k3s-default")
		assert.NoError(t, err)
		assert.IsType(t, actual, &expected)
	})

	t.Run("should error out while fetching cluster info as specified cluster not created", func(t *testing.T) {
		cluster, err := GetCluster(context.Background(), runtimes.Docker, "k3d-default-1")
		assert.Error(t, err, "unable to fetch nodes info as cluster k3d-default-1 not found")
		assert.Nil(t, cluster)
	})
}

func TestGetCluster(t *testing.T) {
	type args struct {
		ctx     context.Context
		runtime runtimes.Runtime
		cluster string
	}
	tests := []struct {
		name    string
		args    args
		want    *K3D.Cluster
		wantErr bool
	}{
		{
			name:    "should be able to get cluster info without error",
			args:    args{ctx: context.Background(), runtime: runtimes.Docker, cluster: "k3s-default"},
			want:    &K3D.Cluster{Name: "k3s-default"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCluster(tt.args.ctx, tt.args.runtime, tt.args.cluster) //nolint:scopelint
			if (err != nil) != tt.wantErr {                                       //nolint:scopelint
				t.Errorf("GetCluster() error = %v, wantErr %v", err, tt.wantErr) //nolint:scopelint
				return
			}
			if !reflect.DeepEqual(got.Name, tt.want.Name) { //nolint:scopelint
				t.Errorf("GetCluster() got = %v, want %v", got, tt.want) //nolint:scopelint
			}
		})
	}
}

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
	expected := []*Config{
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
		cfg := Config{All: false}
		clusters, err := cfg.GetClusters(context.Background(), runtimes.SelectedRuntime, clusters)
		assert.Nil(t, err)
		assert.ElementsMatch(t, expected, clusters)
	})

	t.Run("should be able to fetch all the available clusters", func(t *testing.T) {
		cfg := Config{All: true}
		expected = append(expected, &Config{
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
