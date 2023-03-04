package cluster

import (
	"context"
	"strings"

	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
)

func CheckAndDeleteCluster(ctx context.Context, runtime runtimes.Runtime, cluster string) error {
	cfg := Config{Name: cluster}

	clusters, err := cfg.GetClusters(ctx, runtime, []string{cluster})
	if err != nil {
		if strings.Contains(err.Error(), "No nodes found for given cluster") {
			return nil
		}

		return err
	}

	if len(clusters) == 0 {
		return nil
	}

	return client.ClusterDelete(ctx, runtime, clusters[0].GetClusterConfig(), K3D.ClusterDeleteOpts{SkipRegistryCheck: true})
}
