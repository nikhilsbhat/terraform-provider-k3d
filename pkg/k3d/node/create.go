package node

import (
	"context"
	"fmt"
	dockerunits "github.com/docker/go-units"
	cluster2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// CreateNodeWithTimeout creates node by setting timeouts as per input.
func (cfg *Config) CreateNodeWithTimeout(ctx context.Context, runtime runtimes.Runtime, nodes []*Config) error {
	var nodeCreatOpts K3D.NodeCreateOpts
	if cfg.Wait {
		nodeCreatOpts = K3D.NodeCreateOpts{Wait: cfg.Wait, Timeout: cfg.Timeout}
	}
	clusterFetched, err := cluster2.GetCluster(ctx, runtime, cfg.ClusterAssociated)
	if err != nil {
		return err
	}
	k3dNodes := make([]*K3D.Node, 0)
	for _, node := range nodes {
		k3dNodes = append(k3dNodes, node.GetNodeFromConfig())
	}
	if err = client.NodeAddToClusterMulti(ctx, runtime, k3dNodes, clusterFetched, nodeCreatOpts); err != nil {
		return err
	}

	return nil
}

// CreateNodes creates number nodes specified in 'replicas', making this startFrom if in case we support update nodes on it.
func (cfg *Config) CreateNodes(ctx context.Context, runtime runtimes.Runtime, startFrom int) error {
	nodesToCreate := make([]*Config, 0)

	if _, err := dockerunits.RAMInBytes(cfg.Memory); cfg.Memory != "" && err != nil {
		return fmt.Errorf("provided memory limit value is invalid")
	}

	for startFrom < cfg.Count {
		nodesToCreate = append(nodesToCreate, &Config{
			Name:    []string{fmt.Sprintf("%s-%d", cfg.Name[0], startFrom)},
			Role:    cfg.Role,
			Image:   cfg.Image,
			Memory:  cfg.Memory,
			Created: cfg.Created,
		})
		startFrom++
	}

	if createRrr := cfg.CreateNodeWithTimeout(ctx, runtime, nodesToCreate); createRrr != nil {
		log.Printf("creating nodes errord with: %v, cleaning up the created nodes to avoid dangling nodes", createRrr)
		for _, nodeToCreate := range nodesToCreate {
			nd := nodeToCreate.GetNodeFromConfig()
			log.Printf("cleaning up node: %s", nd.Name)
			if err := DeleteNodesFromCluster(ctx, runtime, nd); err != nil {
				log.Printf("errored while deleting node %s : %v", nd.Name, err)
			}
		}
		log.Printf("creating nodes failed")
		return fmt.Errorf("creating nodes failed with: %s", createRrr.Error())
	}
	return nil
}
