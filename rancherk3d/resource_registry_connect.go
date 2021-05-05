package rancherk3d

import (
	"context"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/k3d"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/utils"
	"github.com/rancher/k3d/v4/pkg/runtimes"
)

func resourceConnectRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectRegistryCluster,
		ReadContext:   resourceConnectRegistryRead,
		DeleteContext: resourceConnectRegistryDelete,
		UpdateContext: resourceConnectRegistryUpdate,
		Schema: map[string]*schema.Schema{
			"registries": {
				Type:        schema.TypeList,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "list of registries to be connected to the selected cluster",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "cluster to which registries to be associated with",
			},
			"connect": {
				Type:        schema.TypeBool,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "enable this flag if registries to be connected with specified cluster",
			},
			"status": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "updated status of registry",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"registry": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "name of the registry",
						},
						"cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "cluster to which the registry to be connected",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "updated state of registry with cluster",
						},
					},
				},
			},
		},
	}
}

func resourceConnectRegistryCluster(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

		connect := k3d.RegistryConnect{
			Registries: getSlice(d.Get(utils.TerraformResourceRegistries)),
			Cluster:    utils.String(d.Get(utils.TerraformResourceCluster)),
			Connect:    utils.Bool(d.Get(utils.TerraformResourceConnect)),
		}

		if err := connectRegistryToCluster(ctx, defaultConfig.K3DRuntime, connect); err != nil {
			return diag.Errorf("errored while connecting/disconnecting registries '%v' with cluster '%s,", connect.Registries, connect.Cluster)
		}

		d.SetId(id)
		return resourceConnectRegistryRead(ctx, d, meta)
	}
	return nil
}

func resourceConnectRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.Config)

	connect := k3d.RegistryConnect{
		Registries: getSlice(d.Get(utils.TerraformResourceRegistries)),
		Cluster:    utils.String(d.Get(utils.TerraformResourceCluster)),
		Connect:    utils.Bool(d.Get(utils.TerraformResourceConnect)),
	}

	registryStatus, err := getRegistryStatus(ctx, defaultConfig.K3DRuntime, connect)
	if err != nil {
		return diag.Errorf("errored while retrieving updated registries status '%v' from cluster '%s,", connect.Registries, connect.Cluster)
	}

	if err = d.Set(utils.TerraformResourceStatus, registryStatus); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceStatus, err)
	}
	return nil
}

func resourceConnectRegistryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.Config)

	if d.HasChange(utils.TerraformResourceRegistries) || d.HasChange(utils.TerraformResourceCluster) ||
		d.HasChange(utils.TerraformResourceConnect) || d.HasChange(utils.TerraformResourceStop) {
		connect := getUpdatedRegistriesChanges(d)
		if err := connectRegistryToCluster(ctx, defaultConfig.K3DRuntime, connect); err != nil {
			return diag.Errorf("errored while connecting/disconnecting registries '%v' with cluster '%s,", connect.Registries, connect.Cluster)
		}

		if err := d.Set(utils.TerraformResourceCluster, connect.Cluster); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceCluster, err)
		}
		if err := d.Set(utils.TerraformResourceRegistries, connect.Registries); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceRegistries, err)
		}
		if err := d.Set(utils.TerraformResourceConnect, connect.Connect); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceConnect, err)
		}

		return resourceConnectRegistryRead(ctx, d, meta)
	}
	log.Printf("nothing to update so skipping")
	return nil
}

func resourceConnectRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.Config)
	_ = defaultConfig
	// could be properly implemented once k3d supports deleting loaded images from cluster.

	id := d.Id()
	if len(id) == 0 {
		return diag.Errorf("resource with the specified ID not found")
	}
	d.SetId("")
	return nil
}

func connectRegistryToCluster(ctx context.Context, runtime runtimes.Runtime, connect k3d.RegistryConnect) error {
	log.Printf("the connect data %v", connect)
	registries, err := k3d.FilteredNodes(ctx, runtime, connect.Registries)
	if err != nil {
		return err
	}
	if connect.Connect {
		if err = k3d.ConnectRegistriesToCluster(ctx, runtime, []string{connect.Cluster}, registries); err != nil {
			return err
		}
		return nil
	}
	if err = k3d.DisconnectRegistriesFormCluster(ctx, runtime, connect.Cluster, registries); err != nil {
		return err
	}
	return nil
}

func getRegistryStatus(ctx context.Context, runtime runtimes.Runtime, connect k3d.RegistryConnect) ([]map[string]string, error) {
	updatedStatus := make([]map[string]string, 0)
	cluster, err := k3d.GetCluster(ctx, runtime, connect.Cluster)
	if err != nil {
		return nil, err
	}
	registries, err := k3d.FilteredNodes(ctx, runtime, connect.Registries)
	if err != nil {
		return nil, err
	}
	for _, registry := range registries {
		log.Printf("registry networks: %v", registry.Networks)
		log.Printf("cluster network: %s", cluster.Network.Name)
		if utils.Contains(registry.Networks, cluster.Network.Name) {
			updatedStatus = append(updatedStatus, map[string]string{
				"registry": registry.Name,
				"cluster":  cluster.Name,
				"state":    "connected",
			})
		} else {
			updatedStatus = append(updatedStatus, map[string]string{
				"registry": registry.Name,
				"cluster":  cluster.Name,
				"state":    "disconnected",
			})
		}
	}
	return updatedStatus, nil
}

func getUpdatedRegistriesChanges(d *schema.ResourceData) (registries k3d.RegistryConnect) {
	oldRegistries, newRegistries := d.GetChange(utils.TerraformResourceRegistries)
	if !cmp.Equal(oldRegistries, newRegistries) {
		registries.Registries = getSlice(newRegistries)
	}
	oldCluster, newCluster := d.GetChange(utils.TerraformResourceCluster)
	if !cmp.Equal(oldCluster, newCluster) {
		registries.Cluster = utils.String(newCluster)
	}
	oldConnect, newConnect := d.GetChange(utils.TerraformResourceConnect)
	if !cmp.Equal(oldConnect, newConnect) {
		registries.Connect = utils.Bool(newConnect)
	}
	return
}
