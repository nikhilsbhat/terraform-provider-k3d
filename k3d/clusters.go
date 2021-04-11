package k3d

import K3D "github.com/rancher/k3d/v4/pkg/types"

func Clusters(cluster string) *K3D.Cluster {
	return &K3D.Cluster{
		Name: cluster,
	}
}
