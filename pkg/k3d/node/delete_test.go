package node_test

import (
	"context"
	"testing"

	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/node"
	"github.com/rancher/k3d/v5/pkg/runtimes"
)

func TestConfig_DeleteNodesFromCluster(t *testing.T) {
	tests := []struct {
		name    string
		cfg     node.Config
		wantErr bool
	}{
		{
			name: "Should be able to delete node test-node-from-terraform",
			cfg: node.Config{
				ClusterAssociated: "k3s-default",
				Name:              []string{"test-node-from-terraform"},
			},
			wantErr: false,
		},
		{
			name: "Should be able to delete node test-node-terraform-0",
			cfg: node.Config{
				ClusterAssociated: "k3s-default",
				Name:              []string{"test-node-terraform-0"},
			},
			wantErr: false,
		},
		{
			name: "Should be able to delete node test-node-terraform-1",
			cfg: node.Config{
				ClusterAssociated: "k3s-default",
				Name:              []string{"test-node-terraform-1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.DeleteNodesFromCluster(context.TODO(), runtimes.Docker); (err != nil) != tt.wantErr {
				t.Errorf("DeleteNodesFromCluster() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
