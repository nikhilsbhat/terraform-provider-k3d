package rancherk3d

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/k3d"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/utils"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	K3D "github.com/rancher/k3d/v4/pkg/types"
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
				Description: "name the nodes to be create (index would be used to dynamically compute the names for nodes)",
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
			"metadata": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "meta data to be used for filtering registries internally",
				Elem:        &schema.Schema{Type: schema.TypeString},
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
	defaultConfig := meta.(*k3d.Config)

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

		proxy, err := utils.Map(d.Get(utils.TerraformResourceProxy))
		if err != nil {
			return diag.Errorf("errored while flattening '%s' with :%v", utils.TerraformResourceProxy, err)
		}

		expose, err := utils.Map(d.Get(utils.TerraformResourceExpose))
		if err != nil {
			return diag.Errorf("errored while flattening '%s' with :%v", utils.TerraformResourceExpose, err)
		}

		if err = d.Set(utils.TerraformResourceMetadata, getMetadata(d)); err != nil {
			return diag.Errorf("errored while setting '%s' with :%v", utils.TerraformResourceHost, err)
		}

		registry := &k3d.Registry{
			Name:     utils.String(d.Get(utils.TerraformResourceName)),
			Image:    utils.String(d.Get(utils.TerraformResourceImage)),
			Cluster:  utils.String(d.Get(utils.TerraformResourceCluster)),
			Host:     utils.String(d.Get(utils.TerraformResourceHost)),
			Protocol: utils.String(d.Get(utils.TerraformResourceProtocol)),
			Proxy:    proxy,
			UseProxy: utils.Bool(d.Get(utils.TerraformUseProxy)),
			Expose:   expose,
		}

		if err = createRegistry(ctx, defaultConfig.K3DRuntime, registry); err != nil {
			diag.Errorf("oops errored while creating registry: %v", err)
		}

		d.SetId(id)
		return resourceRegistryRead(ctx, d, meta)
	}

	return nil
}

func resourceRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.Config)

	cluster := utils.String(d.Get(utils.TerraformResourceCluster))
	metadata, err := utils.Map(d.Get(utils.TerraformResourceMetadata))
	if err != nil {
		return diag.Errorf("errored while fetching '%s'", utils.TerraformResourceMetadata)
	}
	k3dNodes, err := k3d.GetRegistry(ctx, defaultConfig.K3DRuntime, cluster, metadata["host"])
	if err != nil {
		return diag.Errorf("errored while fetching created registries: %v", k3dNodes)
	}

	flattenedRegistryNodes, err := utils.MapSlice(k3dNodes)
	if err != nil {
		return diag.Errorf("errored while flattening obtained created nodes : %v", err)
	}

	if err = d.Set(utils.TerraformResourceRegistriesList, flattenedRegistryNodes); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceRegistriesList, err)
	}

	return nil
}

func resourceRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.Config)

	id := d.Id()
	if len(id) == 0 {
		return diag.Errorf("resource with the specified ID not found")
	}

	nodes, err := getRegistriesFromState(ctx, d, defaultConfig)
	if err != nil {
		return diag.Errorf("errored while retrieving registries from state %v", err)
	}

	for _, node := range nodes {
		if err = k3d.DeletNodesFromCluster(ctx, defaultConfig.K3DRuntime, node); err != nil {
			return diag.Errorf("errored while deleting registry node %s : %v", node.Name, err)
		}
	}
	d.SetId("")
	return nil
}

func getRegistriesFromState(ctx context.Context, d *schema.ResourceData, defaultConfig *k3d.Config) ([]*K3D.Node, error) {
	registriesFromState := d.Get(utils.TerraformResourceRegistriesList)
	var nodes []*k3d.K3Node
	if err := mapstructure.Decode(registriesFromState, &nodes); err != nil {
		return nil, err
	}
	k3dNodes := make([]*K3D.Node, 0)
	for _, node := range nodes {
		nd, err := k3d.Node(ctx, defaultConfig.K3DRuntime, node.Name)
		if err != nil {
			return nil, err
		}
		k3dNodes = append(k3dNodes, nd)
	}
	return k3dNodes, nil
}

func createRegistry(ctx context.Context, runtime runtimes.Runtime, registry *k3d.Registry) error {
	k3dRegistry := &K3D.Registry{}

	k3dRegistry.ClusterRef = registry.Cluster
	k3dRegistry.Protocol = registry.Protocol
	k3dRegistry.Host = getHost(registry)
	k3dRegistry.Image = registry.Image
	k3d.GetExposureOpts(getExpose(registry.Expose), k3dRegistry)
	if registry.UseProxy {
		if !validateProxy(registry.Proxy) {
			k3d.GetProxyConfig(registry.Proxy, k3dRegistry)
		}
		return fmt.Errorf("proxy config validation failed, config cannot be empty")
	}

	if err := k3d.CreateRegistry(ctx, runtime, k3dRegistry, []string{registry.Cluster}); err != nil {
		return err
	}
	return nil
}

func validateProxy(value map[string]string) bool {
	if len(value["remoteURL"]) == 0 || len(value["username"]) == 0 || len(value["password"]) == 0 {
		return false
	}
	return true
}

func validateExpose(value map[string]string) bool {
	if len(value["hostIp"]) == 0 || len(value["hostPort"]) == 0 {
		return false
	}
	return true
}

func getExpose(expose map[string]string) map[string]string {
	if !validateExpose(expose) {
		return map[string]string{
			"hostIp":   "0.0.0.0",
			"hostPort": "5200",
		}
	}
	return expose
}

func getHost(registry *k3d.Registry) string {
	if len(registry.Host) == 0 {
		return registry.Name
	}
	return registry.Host
}

func getMetadata(d *schema.ResourceData) map[string]string {
	metadata := make(map[string]string)
	if host := utils.String(d.Get(utils.TerraformResourceHost)); len(host) == 0 {
		metadata["host"] = utils.String(d.Get(utils.TerraformResourceName))
		return metadata
	}
	metadata["host"] = utils.String(d.Get(utils.TerraformResourceHost))
	return metadata
}
