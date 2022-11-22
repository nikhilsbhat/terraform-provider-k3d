package rancherk3d

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-k3d/pkg/client"
	k3dCluster "github.com/nikhilsbhat/terraform-provider-k3d/pkg/k3d/cluster"
	utils2 "github.com/nikhilsbhat/terraform-provider-k3d/pkg/utils"
)

func dataSourceClusterList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceListClusterRead,
		Schema: map[string]*schema.Schema{
			"clusters": {
				Type:        schema.TypeList,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "list of clusters of which the information to be fetched",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"all": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "if enabled fetches all clusters available",
			},
			"clusters_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of cluster of which information has been retrieved",
				Elem: &schema.Resource{
					Schema: resourceClusterSchema(),
				},
			},
		},
	}
}

func dataSourceListClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	clusters := utils2.GetSlice(d.Get(utils2.TerraformResourceClusters).([]interface{}))
	cfg := k3dCluster.Config{
		All: utils2.Bool(d.Get(utils2.TerraformResourceAll)),
	}

	k3dClusters, err := cfg.GetClusters(ctx, defaultConfig.K3DRuntime, clusters)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while fetching clusters: %v", err)
	}

	flattenedClusters, err := utils2.MapSlice(k3dClusters)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while flattening nodes obtained: %v", err)
	}

	d.SetId(id)
	if err := d.Set(utils2.TerraformResourceClusterList, flattenedClusters); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceClusterList, err)
	}

	return nil
}
