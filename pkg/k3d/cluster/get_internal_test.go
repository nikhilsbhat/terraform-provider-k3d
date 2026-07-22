package cluster

import (
	"reflect"
	"testing"

	K3D "github.com/rancher/k3d/v5/pkg/types"
)

func TestFilterClustersByName(t *testing.T) {
	clusters := []*K3D.Cluster{
		{Name: "alpha"},
		{Name: "beta"},
		{Name: "gamma"},
	}

	filtered := filterClustersByName(clusters, []string{"gamma", "alpha"})

	if got, want := len(filtered), 2; got != want {
		t.Fatalf("expected %d clusters, got %d", want, got)
	}

	if got, want := []string{filtered[0].Name, filtered[1].Name}, []string{"alpha", "gamma"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("expected clusters %v, got %v", want, got)
	}
}

func TestNodeNames(t *testing.T) {
	nodes := []*K3D.Node{
		{Name: "server-0"},
		{Name: "agent-0"},
	}

	if got, want := nodeNames(nodes), []string{"server-0", "agent-0"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("expected node names %v, got %v", want, got)
	}
}

func TestConfigGetClusterConfigPreallocatesNodes(t *testing.T) {
	cfg := Config{
		Name:  "sample",
		Token: "token",
		Nodes: []string{"server-0", "agent-0"},
	}

	cluster := cfg.GetClusterConfig()

	if got, want := cluster.Name, cfg.Name; got != want {
		t.Fatalf("expected cluster name %q, got %q", want, got)
	}

	if got, want := len(cluster.Nodes), len(cfg.Nodes); got != want {
		t.Fatalf("expected %d nodes, got %d", want, got)
	}

	if got, want := []string{cluster.Nodes[0].Name, cluster.Nodes[1].Name}, cfg.Nodes; !reflect.DeepEqual(got, want) {
		t.Fatalf("expected node names %v, got %v", want, got)
	}
}
