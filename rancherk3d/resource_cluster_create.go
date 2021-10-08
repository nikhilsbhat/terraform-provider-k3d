package rancherk3d

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/client"
	utils2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
	"github.com/rancher/k3d/v4/pkg/config/v1alpha2"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		DeleteContext: resourceClusterDelete,
		UpdateContext: resourceClusterUpdate,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "cluster name that was fetched",
			},
			"servers_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "count of servers",
			},
			"agents_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "count of agents in the cluster",
			},
			"image": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("K3D_IMAGE", "rancher/k3s"),
				ForceNew:    true,
				Description: "image name to be used for creation of cluster, it would be used along with kubernetes_version",
				Computed:    true,
			},
			"network": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "network to be associated with the cluster",
				Computed:    false,
			},
			"cluster_token": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "superSecretToken to be used",
				Computed:    false,
			},
			"kubeconfig_update_default": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Description: "Directly update the default kubeconfig with the new cluster's context.",
				Default:     false,
				Computed:    false,
			},
			"kubeconfig_switch_context": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Description: "Directly switch the default kubeconfig's current-context to the new cluster's context",
				Default:     false,
				Computed:    false,
			},
			"labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				Description: "Labels to be added node container.",
				Elem: &schema.Resource{
					Schema: resourceClusterEnvsAndLabelsSchema(),
				},
			},
			"env": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				Description: "Environment variables to be added nodes.",
				Elem: &schema.Resource{
					Schema: resourceClusterEnvsAndLabelsSchema(),
				},
			},
			"registries": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: false,
				Elem: &schema.Resource{
					Schema: resourceClusterRegistriesSchema(),
				},
			},
			"runtime": {
				Description: "Runtime options for k3d",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: resourceClusterRuntimeSchema(),
				},
			},
			"volumes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: false,
				Elem: &schema.Resource{
					Schema: resourceClusterVolumeSchema(),
				},
			},
			"ports": {
				Type:        schema.TypeList,
				ForceNew:    true,
				Optional:    true,
				Description: "Map ports from the node containers to the host.",
				Elem: &schema.Resource{
					Schema: resourceClusterPortsConfig(),
				},
			},
			"k3d_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: false,
				Elem: &schema.Resource{
					Schema: resourceClusterK3sOptionsSchema(),
				},
			},
			"k3s_options": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				Description: "various options to be passed to k3s.",
				Elem: &schema.Resource{
					Schema: resourceClusterK3dOptionsSchema(),
				},
			},
		},
	}
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)
	if d.IsNewResource() {
		clusterName := utils2.String(d.Get(utils2.TerraformResourceName))
		_ = defaultConfig

		d.SetId(clusterName)
	}
	return nil
}

func resourceClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Your code goes here
	return nil
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Your code goes here
	return nil
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func createCluster(ctx context.Context, d *schema.ResourceData, defaultConfig *client.Config) error {
	k3DSimpleConfig := &v1alpha2.SimpleConfig{}
	_ = k3DSimpleConfig
	return nil
}
