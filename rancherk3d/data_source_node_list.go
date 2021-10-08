package rancherk3d

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/client"
	k3dNode "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/node"
	utils2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
)

func dataSourceNodeList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceListNodeRead,
		Schema: map[string]*schema.Schema{
			"nodes": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "list of nodes to be listed from the cluster selected",
			},
			"all": {
				Type:        schema.TypeBool,
				Computed:    false,
				Optional:    true,
				Description: "if enabled fetches all the nodes available in the selected cluster",
			},
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "name of the cluster of which that nodes to be listed",
			},
			"node_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of nodes that were retrieved",
				Elem: &schema.Resource{
					Schema: resourceNodeSchema(),
				},
			},
		},
	}
}

func dataSourceListNodeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	id := d.Id()

	if len(id) == 0 {
		newID, err := utils2.GetRandomID()
		if err != nil {
			d.SetId("")
			return diag.Errorf("errored while fetching randomID %v", err)
		}
		id = newID
	}

	nodes := getSlice(d.Get(utils2.TerraformResourceNodes))
	cluster := utils2.String(d.Get(utils2.TerraformResourceCluster))
	all := utils2.Bool(d.Get(utils2.TerraformResourceAll))

	k3dNodes, err := getNodesFromCluster(ctx, defaultConfig, cluster, nodes, all)
	if err != nil {
		d.SetId("")
		return diag.Errorf("errored while fetching nodes: %v", err)
	}

	flattenedNodes, err := utils2.MapSlice(k3dNodes)
	if err != nil {
		d.SetId("")
		return diag.Errorf("errored while flattening nodes obtained: %v", err)
	}
	d.SetId(id)
	if err := d.Set(utils2.TerraformResourceNodesList, flattenedNodes); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceNodesList, err)
	}

	return nil
}

func getNodesFromCluster(ctx context.Context, defaultConfig *client.Config, cluster string, nodes []string, all bool) ([]*k3dNode.K3Node, error) {
	if all {
		k3dNodes, err := k3dNode.GetFilteredNodesFromCluster(ctx, defaultConfig.K3DRuntime, cluster)
		if err != nil {
			return nil, err
		}
		return k3dNodes, err
	}
	return k3dNode.GetFilteredNodes(ctx, defaultConfig.K3DRuntime, nodes)
}

func getSlice(data interface{}) []string {
	return utils2.GetSlice(data.([]interface{}))
}
