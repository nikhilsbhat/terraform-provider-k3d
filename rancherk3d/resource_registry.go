package rancherk3d

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRegistryCreate,
		ReadContext:   resourceRegistryRead,
		DeleteContext: resourceRegistryDelete,
		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("K3D_CLUSTER_NAME", nil),
				ForceNew:    true,
				Description: "cluster to be associated wth the registry",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "protocol to be used while running registry ",
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("K3D_REGISTRY_HOSTNAME", nil),
				ForceNew:    true,
				Description: "image to be used for creation of registry",
			},
			"image": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("K3D_REGISTRY_NAME", "docker.io/library/registry:2"),
				ForceNew:    true,
				Description: "image to be used for creation of registry",
			},
			"exposureOpts": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "extra options to be passed",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: false,
						},
						"binding": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostip": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: false,
									},
									"hostport": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: false,
									},
								},
							},
						},
					},
				},
			},
			"options": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "options to be passed to registry",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config-file": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    false,
							Description: "path to config file to be used for registry creation",
						},
						"proxy": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: false,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"registries": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "list of registries those were created",
				Elem:        &schema.Schema{Type: schema.TypeList},
			},
		},
	}
}

func resourceRegistryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Your code goes here
	return nil
}

func resourceRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Your code goes here
	return nil
}
