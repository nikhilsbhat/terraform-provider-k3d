package rancherk3d

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/client"
	k3dNode "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/node"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
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
				Description: "name the nodes to be created (index would be used to dynamically compute the names for nodes)",
			},
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "name of the cluster to which these nodes to be connected with",
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
				Description: "role to be assigned to the node(agent)",
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
				Description: "if enabled waits for nodes to be ready before returning",
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
	defaultConfig := meta.(*client.Config)

	//nolint:nestif
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
		nodeConfig := k3dNode.Config{
			Name:              []string{utils.String(d.Get(utils.TerraformResourceName))},
			ClusterAssociated: utils.String(d.Get(utils.TerraformResourceCluster)),
			Image:             setNodeImage(d, defaultConfig),
			Role:              utils.String(d.Get(utils.TerraformResourceRole)),
			Count:             utils.Int(d.Get(utils.TerraformResourceReplicas)),
			Memory:            utils.String(d.Get(utils.TerraformResourceMemory)),
			Created:           createInitiatedAt.String(),
			Timeout:           time.Duration(utils.Int(d.Get(utils.TerraformResourceTimeout))) * time.Minute,
			Wait:              utils.Bool(d.Get(utils.TerraformResourceWait)),
		}

		if err := nodeConfig.CreateNodes(ctx, defaultConfig.K3DRuntime, 0); err != nil {
			if seErr := d.Set(utils.TerraformResourceCreatedAt, ""); seErr != nil {
				return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceCreatedAt, seErr)
			}

			return diag.Errorf("errored while creating nodes with: %v", err.Error())
		}

		d.SetId(id)

		return resourceNodeRead(ctx, d, meta)
	}

	return nil
}

func resourceNodeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	created := utils.String(d.Get(utils.TerraformResourceCreatedAt))

	cfg := k3dNode.Config{Labels: getTerraformTimestampLabel(created)}

	k3dNodes, err := cfg.GetNodesByLabels(ctx, defaultConfig.K3DRuntime)
	if err != nil {
		return diag.Errorf("errored while fetching created nodes: %v", k3dNodes)
	}

	flattenedK3dNodes, err := utils.MapSlice(k3dNodes)
	if err != nil {
		return diag.Errorf("errored while flattening obtained created nodes : %v", err)
	}

	if err = d.Set(utils.TerraformResourceNodes, flattenedK3dNodes); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceNodes, err)
	}

	return nil
}

func resourceNodeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	id := d.Id()
	if len(id) == 0 {
		return diag.Errorf("resource with the specified ID not found")
	}

	nodesFromState := d.Get(utils.TerraformResourceNodes)
	var nodes []*k3dNode.Config
	if err := mapstructure.Decode(nodesFromState, &nodes); err != nil {
		return diag.Errorf("oops reading nodes from state errored with : %s", err.Error())
	}

	for _, node := range nodes {
		if err := node.DeleteNodesFromCluster(ctx, defaultConfig.K3DRuntime); err != nil {
			return diag.Errorf("oops deleting node %s errored with : %s", node.Name[0], err.Error())
		}
	}

	d.SetId("")

	return nil
}

func setNodeImage(d *schema.ResourceData, defaultConfig *client.Config) string {
	image := utils.String(d.Get(utils.TerraformResourceImage))
	if len(image) == 0 {
		return defaultConfig.GetK3dImage()
	}

	return image
}

func getTerraformTimestampLabel(time string) map[string]string {
	return map[string]string{
		utils.TerraformK3dLabel: time,
	}
}
