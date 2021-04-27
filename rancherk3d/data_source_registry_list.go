package rancherk3d

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/k3d"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/utils"
)

func dataSourceRegistryList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRegistryListRead,
		Schema: map[string]*schema.Schema{
			"registries": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "list of registries to be retrieved from the cluster selected",
			},
			"all": {
				Type:        schema.TypeBool,
				Computed:    false,
				Optional:    true,
				Description: "if enabled fetches all the registries, if cluster is selected then all registries connected to it",
			},
			"cluster": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "name of the cluster of which registries to be retrieved",
			},
			"registries_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of registries retrieved",
				Elem: &schema.Resource{
					Schema: resourceRegistrySchema(),
				},
			},
		},
	}
}

func dataSourceRegistryListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	registries := getSlice(d.Get(utils.TerraformResourceRegistries))
	cluster := utils.String(d.Get(utils.TerraformResourceCluster))
	all := utils.Bool(d.Get(utils.TerraformResourceAll))

	k3dNodes, err := getRegistries(ctx, defaultConfig, cluster, registries, all)
	if err != nil {
		d.SetId("")
		return diag.Errorf("errored while fetching registry nodes: %v", err)
	}

	if len(k3dNodes) == 0 {
		return diag.Errorf("either there are no registries in the environment or with the specified configurations")
	}

	flattenedNodes, err := utils.Map(k3dNodes)
	if err != nil {
		d.SetId("")
		return diag.Errorf("errored while flattening registry nodes obtained: %v", err)
	}

	d.SetId(id)
	if err = d.Set(utils.TerraformResourceRegistriesList, flattenedNodes); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceRegistriesList, err)
	}
	return nil
}

func getRegistries(ctx context.Context, defaultConfig *k3d.K3dConfig, cluster string, registries []string, all bool) ([]*k3d.K3DNode, error) {
	if all {
		return k3d.GetRegistries(ctx, defaultConfig.K3DRuntime, cluster)
	}
	return k3d.GetRegistriesWithName(ctx, defaultConfig.K3DRuntime, cluster, registries)
}
