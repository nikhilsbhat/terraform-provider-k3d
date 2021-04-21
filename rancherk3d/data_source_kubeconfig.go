package rancherk3d

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/k3d"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/utils"
	K3D "github.com/rancher/k3d/v4/pkg/types"
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
	defaultConfig := meta.(*k3d.K3dConfig)

	id := d.Id()
	if len(id) == 0 {
		log.Printf("fetching new ID for kubeconfig: %s", id)
		newID, err := utils.GetRandomID()
		if err != nil {
			d.SetId("")
			return diag.Errorf("errored while fetching randomID %v", err)
		}
		id = newID
	}

	clusters := getClusters(d.Get(utils.TerraformResourceClusters))
	all := utils.Bool(d.Get(utils.TerraformResourceAll))
	notEncode := utils.Bool(d.Get(utils.TerraformResourceNotEncode))

	fetchedClusters, err := getUnfilteredCluster(ctx, defaultConfig, clusters, all)
	if err != nil {
		d.SetId("")
		return diag.Errorf("oops errored while fetching clusters info: %v", err)
	}
	kubeConfig, err := k3d.GetKubeConfig(ctx, defaultConfig.K3DRuntime, fetchedClusters, notEncode)
	if err != nil {
		d.SetId("")
		return diag.Errorf("errored while fetching kube-config: %v", err)
	}

	d.SetId(id)
	if seRrr := d.Set(utils.TerraformResourceKubeConfig, kubeConfig); seRrr != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceKubeConfig, seRrr)
	}

	return nil
}

func getUnfilteredCluster(ctx context.Context, defaultConfig *k3d.K3dConfig, clusters []string, all bool) ([]*K3D.Cluster, error) {
	fetchedClusters := make([]*K3D.Cluster, 0)
	if all {
		allClusters, err := k3d.GetClusters(ctx, defaultConfig.K3DRuntime)
		if err != nil {
			return nil, err
		}
		fetchedClusters = allClusters
	} else {
		allClusters, err := k3d.GetFilteredClusters(ctx, defaultConfig.K3DRuntime, clusters)
		if err != nil {
			return nil, err
		}
		fetchedClusters = allClusters
	}

	return fetchedClusters, nil
}
