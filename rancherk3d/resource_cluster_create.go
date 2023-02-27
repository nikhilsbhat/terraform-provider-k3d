package rancherk3d

import (
	"context"
	"github.com/rancher/k3d/v5/pkg/config/types"
	"github.com/rancher/k3d/v5/pkg/config/v1alpha4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-k3d/pkg/client"
	"github.com/nikhilsbhat/terraform-provider-k3d/pkg/utils"
	"github.com/rancher/k3d/v5/pkg/config/v1alpha2"
)

//nolint:deadcode,unused,funlen
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
				Description: "Name of the Cluster to be created",
			},
			"servers_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Count of servers",
			},
			"agents_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Count of agents in the cluster",
			},
			"image": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("K3D_IMAGE", "rancher/k3s"),
				ForceNew:    true,
				Description: "Image name to be used for creation of cluster, it would be used along with kubernetes_version",
				Computed:    true,
			},
			"network": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Network to be associated with the cluster",
				Computed:    false,
			},
			"subnetwork": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Define a subnet for the newly created container network",
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
			"volumes": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				Description: "Mount volumes into the nodes (Format: [SOURCE:]DEST[@NODEFILTER[;NODEFILTER...]]",
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
			"hostAliases": {
				Type:        schema.TypeList,
				ForceNew:    true,
				Optional:    true,
				Description: "/etc/hosts style entries to be injected into /etc/hosts in the node containers and in the NodeHosts section in CoreDNS.",
				Elem: &schema.Resource{
					Schema: resourceHostAliasesConfig(),
				},
			},
			"k3d_options": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				Description: "k3d runtime settings",
				Elem: &schema.Resource{
					Schema: resourceClusterK3dOptionsSchema(),
				},
			},
			"k3s_options": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				Description: "Options passed on to K3s itself",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"extraArgs": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    false,
							Description: "additional arguments passed to the `k3s server|agent` command; same as `--k3s-arg`",
							Elem: &schema.Resource{
								Schema: resourceClusterK3sOptionsSchema(),
							},
						},
						"nodeLabels": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    false,
							Description: "same as `--k3s-node-label 'foo=bar@agent:1'` -> this results in a Kubernetes node label",
							Elem: &schema.Resource{
								Schema: resourceClusterEnvsAndLabelsSchema(),
							},
						},
					},
				},
			},
			"kubeconfig": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Way to manage the kubeconfig generated after creating k3d clusters.",
				Default:     false,
				Computed:    false,
				Elem: &schema.Resource{
					Schema: resourceKubeconfigConfig(),
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
		},
	}
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)
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

		cfg := &v1alpha4.SimpleConfig{
			ObjectMeta: types.ObjectMeta{
				Name: utils.String(d.Get(utils.TerraformResourceName)),
			},
			Servers:      utils.Int(d.Get(utils.TerraformResourceServersCount)),
			Agents:       utils.Int(d.Get(utils.TerraformResourceAgentsCount)),
			ExposeAPI:    v1alpha4.SimpleExposureOpts{},
			Image:        utils.String(d.Get(utils.TerraformResourceImage)),
			Network:      utils.String(d.Get(utils.TerraformResourceNetwork)),
			Subnet:       utils.String(d.Get(utils.TerraformResourceSubnet)),
			ClusterToken: utils.String(d.Get(utils.TerraformResourceClusterToken)),
			Volumes: []v1alpha4.VolumeWithNodeFilters{
				{},
			},
			Ports:       nil,
			Options:     v1alpha4.SimpleConfigOptions{},
			Env:         nil,
			Registries:  v1alpha4.SimpleConfigRegistries{},
			HostAliases: nil,
		}
	}

	return nil
}

//nolint:deadcode,unused
func resourceClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Your code goes here
	return nil
}

//nolint:deadcode,unused
func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Your code goes here
	return nil
}

//nolint:deadcode,unused
func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

//nolint:deadcode,unused
func createCluster(ctx context.Context, d *schema.ResourceData, defaultConfig *client.Config) error {
	k3DSimpleConfig := &v1alpha2.SimpleConfig{}
	_ = k3DSimpleConfig

	return nil
}
