package rancherk3d

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/client"
	k3dCluster "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	k3dKube "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/config"
	utils2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
)

func dataSourceKubeConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubeConfigRead,
		Schema: map[string]*schema.Schema{
			"clusters": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "list of k3s clusters from which the kube-config has to be retrieved",
			},
			"all": {
				Type:        schema.TypeBool,
				Computed:    false,
				Optional:    true,
				Description: "set this if kube-config has to be fetched from all available k3s clusters",
			},
			"not_encoded": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Description: "set it to false if base64 encoding of kube-config is not preferred, defaults to true",
			},
			"kube_config": {
				Type:        schema.TypeMap,
				Computed:    true,
				Optional:    true,
				Sensitive:   true,
				Description: "retrieved kube-config for the specified cluster",
			},
		},
	}
}

func dataSourceKubeConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	id := d.Id()
	if len(id) == 0 {
		log.Printf("fetching new ID for kubeconfig: %s", id)
		newID, err := utils2.GetRandomID()
		if err != nil {
			d.SetId("")

			return diag.Errorf("errored while fetching randomID %v", err)
		}
		id = newID
	}

	cfg := k3dKube.Config{
		All:    utils2.Bool(d.Get(utils2.TerraformResourceAll)),
		Encode: !(utils2.Bool(d.Get(utils2.TerraformResourceNotEncode))),
	}

	clusterCfg := k3dCluster.Config{
		All: cfg.All,
	}

	clusters, err := clusterCfg.GetClusters(ctx, defaultConfig.K3DRuntime,
		utils2.GetSlice(d.Get(utils2.TerraformResourceClusters).([]interface{})))
	if err != nil {
		return diag.Errorf("fetching cluster information errored with: %v", err)
	}

	cfg.Cluster = clusters

	kubeConfig, err := cfg.GetKubeConfig(ctx, defaultConfig.K3DRuntime)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while fetching kube-config: %v", err)
	}

	d.SetId(id)
	if seRrr := d.Set(utils2.TerraformResourceKubeConfig, kubeConfig); seRrr != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceKubeConfig, seRrr)
	}

	return nil
}
