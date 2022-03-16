package rancherk3d

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/client"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	k3d2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/config"
	utils2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
	K3D "github.com/rancher/k3d/v5/pkg/types"
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

	clusters := getClusters(d.Get(utils2.TerraformResourceClusters))
	all := utils2.Bool(d.Get(utils2.TerraformResourceAll))
	notEncode := utils2.Bool(d.Get(utils2.TerraformResourceNotEncode))

	fetchedClusters, err := getUnfilteredCluster(ctx, defaultConfig, clusters, all)
	if err != nil {
		d.SetId("")
		return diag.Errorf("oops errored while fetching clusters info: %v", err)
	}
	kubeConfig, err := k3d2.GetKubeConfig(ctx, defaultConfig.K3DRuntime, fetchedClusters, notEncode)
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

func getUnfilteredCluster(ctx context.Context, defaultConfig *client.Config, clusters []string, all bool) ([]*K3D.Cluster, error) {
	var fetchedClusters []*K3D.Cluster
	if all {
		allClusters, err := cluster.GetClusters(ctx, defaultConfig.K3DRuntime)
		if err != nil {
			return nil, err
		}
		fetchedClusters = allClusters
	} else {
		allClusters, err := cluster.GetFilteredClusters(ctx, defaultConfig.K3DRuntime, clusters)
		if err != nil {
			return nil, err
		}
		fetchedClusters = allClusters
	}

	return fetchedClusters, nil
}
