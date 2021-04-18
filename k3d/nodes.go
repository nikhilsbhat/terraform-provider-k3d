package k3d

import (
	"context"
	"fmt"
	"log"

	"github.com/rancher/k3d/v4/pkg/client"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

func GetNodes(ctx context.Context, runtime runtimes.Runtime, cluster string) ([]*K3D.Node, error) {
	nodes, err := client.NodeList(ctx, runtime)
	if err != nil {
		return nil, err
	}
	filteredNodes := make([]*K3D.Node, 0)
	for _, node := range nodes {
		log.Print(node.Labels)
		if node.Labels["k3d.cluster"] == cluster {
			filteredNodes = append(filteredNodes, node)
		}
	}
	if len(filteredNodes) == 0 {
		return nil, fmt.Errorf("unable to fetch nodes info as cluster %s not found", cluster)
	}
	return filteredNodes, nil
}
