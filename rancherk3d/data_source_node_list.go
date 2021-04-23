package rancherk3d

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/k3d"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/utils"
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
				Description: "list of nodes and its details fetched from the specified cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "cluster to which the below images are stored",
						},
						"role": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "cluster to which the below images are stored",
						},
						"cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "cluster to which the below images are stored",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "cluster to which the below images are stored",
						},
						"created": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "cluster to which the below images are stored",
						},
						"volumes": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: "details of images and its tarball stored, if in case keep_tarball is enabled",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"networks": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: "details of images and its tarball stored, if in case keep_tarball is enabled",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"env": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: "details of images and its tarball stored, if in case keep_tarball is enabled",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceListNodeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.K3dConfig)

	id := d.Id()

	if len(id) == 0 {
		newID, err := utils.GetRandomID()
		if err != nil {
			d.SetId("")
			return diag.Errorf("errored while fetching randomID %v", err)
		}
		id = newID
	}

	nodes := getNodesSlice(d.Get(utils.TerraformResourceNodes))
	cluster := utils.String(d.Get(utils.TerraformResourceCluster))
	all := utils.Bool(d.Get(utils.TerraformResourceAll))

	k3dNodes, err := getNodesFromCluster(ctx, defaultConfig, cluster, nodes, all)
	if err != nil {
		d.SetId("")
		return diag.Errorf("errored while fetching nodes: %v", err)
	}

	flattenedNodes, err := utils.Map(k3dNodes)
	if err != nil {
		d.SetId("")
		return diag.Errorf("errored while flattening nodes obtained: %v", err)
	}
	d.SetId(id)
	if err := d.Set(utils.TerraformResourceNodesList, flattenedNodes); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceNodesList, err)
	}

	return nil
}

func getNodesFromCluster(ctx context.Context, defaultConfig *k3d.K3dConfig, cluster string, nodes []string, all bool) ([]*k3d.K3DNode, error) {
	if all {
		k3dNodes, err := k3d.GetFilteredNodesFromCluster(ctx, defaultConfig.K3DRuntime, cluster)
		if err != nil {
			return nil, err
		}
		return k3dNodes, err
	}
	return k3d.GetFilteredNodes(ctx, defaultConfig.K3DRuntime, nodes)
}

func getNodesSlice(nodes interface{}) []string {
	return utils.GetSlice(nodes.([]interface{}))
}
