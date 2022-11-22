package rancherk3d

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-k3d/pkg/client"
	k3dRegistry "github.com/nikhilsbhat/terraform-provider-k3d/pkg/k3d/registry"
	utils2 "github.com/nikhilsbhat/terraform-provider-k3d/pkg/utils"
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

	registry := &k3dRegistry.Config{
		Name:    getSlice(d.Get(utils2.TerraformResourceRegistries)),
		Cluster: utils2.String(d.Get(utils2.TerraformResourceCluster)),
		All:     utils2.Bool(d.Get(utils2.TerraformResourceAll)),
	}

	k3dNodes, err := registry.Get(ctx, defaultConfig.K3DRuntime)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while fetching registry nodes: %v", err)
	}

	if len(k3dNodes) == 0 {
		return diag.Errorf("either there are no registries in the environment or with the specified configurations")
	}

	flattenedNodes, err := utils2.MapSlice(k3dNodes)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while flattening registry nodes obtained: %v", err)
	}

	d.SetId(id)
	if err = d.Set(utils2.TerraformResourceRegistriesList, flattenedNodes); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceRegistriesList, err)
	}

	return nil
}
