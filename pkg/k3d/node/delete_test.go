package node

import (
	"context"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	"testing"
)

func TestConfig_DeleteNodesFromCluster(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "Should be able to delete node test-node-from-terraform",
			cfg: Config{
				ClusterAssociated: "k3s-default",
				Name:              []string{"test-node-from-terraform"},
			},
			wantErr: false,
		},
		{
			name: "Should be able to delete node test-node-terraform-0",
			cfg: Config{
				ClusterAssociated: "k3s-default",
				Name:              []string{"test-node-terraform-0"},
			},
			wantErr: false,
		},
		{
			name: "Should be able to delete node test-node-terraform-1",
			cfg: Config{
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
