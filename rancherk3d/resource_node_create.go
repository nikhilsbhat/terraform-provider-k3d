package rancherk3d

import (
	"context"
	"fmt"
	"log"
	"time"

	dockerunits "github.com/docker/go-units"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/k3d"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/utils"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

func resourceNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeCreate,
		ReadContext:   resourceNodeRead,
		DeleteContext: resourceNodeDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "name the nodes to be create (index would be used to dynamically compute the names for nodes)",
			},
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "name of the cluster of which that nodes to be acted upon",
			},
			"image": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "image to be used for nodes creation defaults to image declared in provider",
			},
			"role": {
				Type:        schema.TypeString,
				Computed:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "cluster to which the below images are stored",
			},
			"replicas": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "number of nodes to be created",
			},
			"memory": {
				Type:        schema.TypeString,
				Computed:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "memory limit to be imposed on the node",
			},
			"wait": {
				Type:        schema.TypeBool,
				Computed:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "if enabled waits for cluster to be ready before returning",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "maximum waiting time for before canceling/returning in minutes",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "timestamp of nodes creation, this would be used to track the nodes created",
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of nodes that were created",
				Elem: &schema.Resource{
					Schema: resourceNodeSchema(),
				},
			},
		},
	}
}

func resourceNodeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.Config)

	if d.IsNewResource() {
		id := d.Id()

		if len(id) == 0 {
			newID, err := utils.GetRandomID()
			if err != nil {
				d.SetId("")
				return diag.Errorf("errored while fetching randomID %v", err)
			}
			id = newID
		}

		createInitiatedAt := time.Now()
		if err := d.Set(utils.TerraformResourceCreatedAt, createInitiatedAt.String()); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceCreatedAt, err)
		}
		node := k3d.K3Node{
			Name:              utils.String(d.Get(utils.TerraformResourceName)),
			ClusterAssociated: utils.String(d.Get(utils.TerraformResourceCluster)),
			Image:             setNodeImage(d, defaultConfig),
			Role:              utils.String(d.Get(utils.TerraformResourceRole)),
			Count:             utils.Int(d.Get(utils.TerraformResourceReplicas)),
			Memory:            utils.String(d.Get(utils.TerraformResourceMemory)),
			Created:           createInitiatedAt.String(),
		}

		timeout := utils.Int(d.Get(utils.TerraformResourceTimeout))
		wait := utils.Bool(d.Get(utils.TerraformResourceWait))

		if err := createNodes(ctx, defaultConfig.K3DRuntime, node, timeout, wait, 0); err != nil {
			if seErr := d.Set(utils.TerraformResourceCreatedAt, ""); seErr != nil {
				return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceCreatedAt, seErr)
			}
			return diag.Errorf("errored while creating nodes with: %v", err)
		}

		d.SetId(id)
		return resourceNodeRead(ctx, d, meta)
	}
	return nil
}

func resourceNodeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.Config)

	created := utils.String(d.Get(utils.TerraformResourceCreatedAt))
	k3dNodes, err := k3d.GetNodesByLabels(ctx, defaultConfig.K3DRuntime, getTerraformTimestampLabel(created))
	if err != nil {
		return diag.Errorf("errored while fetching created nodes: %v", k3dNodes)
	}

	flattenedk3dNodes, err := utils.MapSlice(k3dNodes)
	if err != nil {
		return diag.Errorf("errored while flattening obtained created nodes : %v", err)
	}

	if err = d.Set(utils.TerraformResourceNodes, flattenedk3dNodes); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceNodes, err)
	}

	return nil
}

func resourceNodeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.Config)

	id := d.Id()
	if len(id) == 0 {
		return diag.Errorf("resource with the specified ID not found")
	}

	nodes, err := getNodesFromState(ctx, d, defaultConfig)
	if err != nil {
		return diag.Errorf("errored while retrieving nodes from state %v", err)
	}

	for _, node := range nodes {
		if err = k3d.DeleteNodesFromCluster(ctx, defaultConfig.K3DRuntime, node); err != nil {
			return diag.Errorf("errored while deleting node %s : %v", node.Name, err)
		}
	}
	d.SetId("")
	return nil
}

// createNodes creates number nodes specified in 'replicas', making this startFrom if in case we support update nodes on it.
func createNodes(ctx context.Context, runtime runtimes.Runtime, node k3d.K3Node, timeout int, wait bool, startFrom int) error {
	nodeTimeout := time.Duration(timeout) * time.Minute
	nodesToCreate := make([]*k3d.K3Node, 0)

	memory := node.Memory
	if _, err := dockerunits.RAMInBytes(memory); memory != "" && err != nil {
		return fmt.Errorf("provided memory limit value is invalid")
	}

	for startFrom < node.Count {
		nodesToCreate = append(nodesToCreate, &k3d.K3Node{
			Name:    fmt.Sprintf("%s-%d", node.Name, startFrom),
			Role:    node.Role,
			Image:   node.Image,
			Memory:  node.Memory,
			Created: node.Created,
		})
		startFrom++
	}
	if err := k3d.CreateNodeWithTimeout(ctx, runtime, node.ClusterAssociated, nodesToCreate, wait, nodeTimeout); err != nil {
		log.Printf("creating nodes errord with: %v, cleaning up the created nodes to avoid dangling nodes", err)
		for _, nodeToCreate := range nodesToCreate {
			nd := nodeToCreate.GetNode()
			log.Printf("cleaning up node: %s", nd.Name)
			if err = k3d.DeleteNodesFromCluster(ctx, runtime, nd); err != nil {
				log.Printf("errored while deleting node %s : %v", nd.Name, err)
			}
		}
		log.Printf("creating nodes failed")
		return fmt.Errorf("creating nodes failed")
	}
	return nil
}

func setNodeImage(d *schema.ResourceData, defaultConfig *k3d.Config) string {
	image := utils.String(d.Get(utils.TerraformResourceImage))
	if len(image) == 0 {
		return getK3dImage(defaultConfig)
	}
	return image
}

func getNodesFromState(ctx context.Context, d *schema.ResourceData, defaultConfig *k3d.Config) ([]*K3D.Node, error) {
	nodesFromState := d.Get(utils.TerraformResourceNodes)
	var nodes []*k3d.K3Node
	if err := mapstructure.Decode(nodesFromState, &nodes); err != nil {
		return nil, err
	}
	k3dNodes := make([]*K3D.Node, 0)
	for _, node := range nodes {
		nd, err := k3d.Node(ctx, defaultConfig.K3DRuntime, node.Name)
		if err != nil {
			return nil, err
		}
		k3dNodes = append(k3dNodes, nd)
	}
	return k3dNodes, nil
}

func getTerraformTimestampLabel(time string) map[string]string {
	return map[string]string{
		utils.TerraformK3dLabel: time,
	}
}

func getK3dImage(defaultConfig *k3d.Config) string {
	return fmt.Sprintf("%s:v%s", defaultConfig.K3DRegistry, defaultConfig.KubeImageVersion)
}
