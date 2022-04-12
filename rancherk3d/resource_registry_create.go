package rancherk3d

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/client"
	k3dNode "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/node"
	k3dRegistry "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/registry"
	utils2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
)

func resourceRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRegistryCreate,
		ReadContext:   resourceRegistryRead,
		DeleteContext: resourceRegistryDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "name the registry node to be created",
			},
			"image": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("K3D_REGISTRY_NAME", "docker.io/library/registry:2"),
				ForceNew:    true,
				Description: "image to be used for creation of registry",
			},
			"cluster": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("K3D_CLUSTER_NAME", nil),
				ForceNew:    true,
				Description: "cluster to which the registry to be associated with",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "protocol to be used while running registry (defaults to http)",
				DefaultFunc: schema.EnvDefaultFunc("K3D_REGISTRY_PROTOCOL", "http"),
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("K3D_REGISTRY_HOSTNAME", nil),
				ForceNew:    true,
				Description: "host name to be assigned to the registry the would be created (defaults to name of registry)",
			},
			"config_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "config file to be used for configuring registry",
			},
			"expose": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "host to port mapping",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"use_proxy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "if enabled proxy config provided at 'proxy' would be used for configuring registry",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"proxy": {
				Type:         schema.TypeMap,
				Optional:     true,
				Computed:     false,
				ForceNew:     true,
				RequiredWith: []string{"use_proxy"},
				Description:  "proxy configurations to be used while configuring registry if enabled",
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"registries_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "list of registries those were created",
				Elem: &schema.Resource{
					Schema: resourceRegistrySchema(),
				},
			},
		},
	}
}

func resourceRegistryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	if d.IsNewResource() {
		id := d.Id()

		if len(id) == 0 {
			newID, err := utils2.GetRandomID()
			if err != nil {
				d.SetId("")
				return diag.Errorf("errored while fetching randomID %v", err)
			}
			id = newID
		}

		proxy, err := utils2.Map(d.Get(utils2.TerraformResourceProxy))
		if err != nil {
			return diag.Errorf("errored while flattening '%s' with :%v", utils2.TerraformResourceProxy, err)
		}

		expose, err := utils2.Map(d.Get(utils2.TerraformResourceExpose))
		if err != nil {
			return diag.Errorf("errored while flattening '%s' with :%v", utils2.TerraformResourceExpose, err)
		}

		//if err = d.Set(utils2.TerraformResourceMetadata, getMetadata(d)); err != nil {
		//	return diag.Errorf("errored while setting '%s' with :%v", utils2.TerraformResourceHost, err)
		//}

		registry := &k3dRegistry.Config{
			Name:     []string{utils2.String(d.Get(utils2.TerraformResourceName))},
			Image:    utils2.String(d.Get(utils2.TerraformResourceImage)),
			Cluster:  utils2.String(d.Get(utils2.TerraformResourceCluster)),
			Host:     validateAndSetHost(d),
			Protocol: utils2.String(d.Get(utils2.TerraformResourceProtocol)),
			Proxy:    validateAndSetProxy(d, proxy),
			UseProxy: utils2.Bool(d.Get(utils2.TerraformUseProxy)),
			Expose:   validateAndSetExpose(expose),
		}

		if err = registry.Create(ctx, defaultConfig.K3DRuntime); err != nil {
			return diag.Errorf("oops errored while creating registry: %v", err)
		}

		d.SetId(id)
		return resourceRegistryRead(ctx, d, meta)
	}

	return nil
}

func resourceRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	registryName := utils2.String(d.Get(utils2.TerraformResourceName))
	registry := &k3dRegistry.Config{
		Name:    []string{registryName},
		Cluster: utils2.String(d.Get(utils2.TerraformResourceCluster)),
	}

	registries, err := registry.Get(ctx, defaultConfig.K3DRuntime)
	if err != nil {
		return diag.Errorf("errored while fetching registries: '%s'", registryName)
	}

	flattenedRegistryNodes, err := utils2.MapSlice(registries)
	if err != nil {
		return diag.Errorf("errored while flattening obtained created nodes : %v", err)
	}

	if err = d.Set(utils2.TerraformResourceRegistriesList, flattenedRegistryNodes); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceRegistriesList, err)
	}

	return nil
}

func resourceRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	id := d.Id()
	if len(id) == 0 {
		return diag.Errorf("resource with the specified ID not found")
	}

	registriesFromState := d.Get(utils2.TerraformResourceRegistriesList)

	var nodes []*k3dNode.Config
	if err := mapstructure.Decode(registriesFromState, &nodes); err != nil {
		return diag.Errorf("oops decoding retrieved registries errored : %s", err.Error())
	}

	for _, node := range nodes {
		if err := node.DeleteNodesFromCluster(ctx, defaultConfig.K3DRuntime); err != nil {
			return diag.Errorf("oops errored while deleting registry node %s : %v", node.Name, err)
		}
	}

	d.SetId("")
	return nil
}
