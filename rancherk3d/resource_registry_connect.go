package rancherk3d

import (
	"context"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/client"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/node"
	k3dRegistry "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/registry"
	utils2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
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
							Description: "cluster to which the registry is either connected/disconnected",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "updated state of registry",
						},
					},
				},
			},
		},
	}
}

func resourceConnectRegistryCluster(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

		connect := k3dRegistry.ConnectConfig{
			Registries: getSlice(d.Get(utils2.TerraformResourceRegistries)),
			Cluster:    utils2.String(d.Get(utils2.TerraformResourceCluster)),
			Connect:    utils2.Bool(d.Get(utils2.TerraformResourceConnect)),
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
	defaultConfig := meta.(*client.Config)

	connect := k3dRegistry.ConnectConfig{
		Registries: getSlice(d.Get(utils2.TerraformResourceRegistries)),
		Cluster:    utils2.String(d.Get(utils2.TerraformResourceCluster)),
		Connect:    utils2.Bool(d.Get(utils2.TerraformResourceConnect)),
	}

	registryStatus, err := getRegistryStatus(ctx, defaultConfig.K3DRuntime, connect)
	if err != nil {
		return diag.Errorf("errored while retrieving updated registries status '%v' from cluster '%s,", connect.Registries, connect.Cluster)
	}

	if err = d.Set(utils2.TerraformResourceStatus, registryStatus); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceStatus, err)
	}
	return nil
}

func resourceConnectRegistryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	if d.HasChange(utils2.TerraformResourceRegistries) || d.HasChange(utils2.TerraformResourceCluster) ||
		d.HasChange(utils2.TerraformResourceConnect) || d.HasChange(utils2.TerraformResourceStop) {
		connect := getUpdatedRegistriesChanges(d)
		if err := connectRegistryToCluster(ctx, defaultConfig.K3DRuntime, connect); err != nil {
			return diag.Errorf("errored while connecting/disconnecting registries '%v' with cluster '%s,", connect.Registries, connect.Cluster)
		}

		if err := d.Set(utils2.TerraformResourceCluster, connect.Cluster); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceCluster, err)
		}
		if err := d.Set(utils2.TerraformResourceRegistries, connect.Registries); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceRegistries, err)
		}
		if err := d.Set(utils2.TerraformResourceConnect, connect.Connect); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceConnect, err)
		}

		return resourceConnectRegistryRead(ctx, d, meta)
	}
	log.Printf("nothing to update so skipping")
	return nil
}

func resourceConnectRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)
	_ = defaultConfig
	// could be properly implemented once k3d supports deleting loaded images from cluster.

	id := d.Id()
	if len(id) == 0 {
		return diag.Errorf("resource with the specified ID not found")
	}
	d.SetId("")
	return nil
}

func connectRegistryToCluster(ctx context.Context, runtime runtimes.Runtime, connect k3dRegistry.ConnectConfig) error {
	log.Printf("the connect data %v", connect)
	registries, err := node.FilteredNodes(ctx, runtime, connect.Registries)
	if err != nil {
		return err
	}
	if connect.Connect {
		if err = k3dRegistry.ConnectRegistriesToCluster(ctx, runtime, []string{connect.Cluster}, registries); err != nil {
			return err
		}
		return nil
	}
	if err = k3dRegistry.DisconnectRegistriesFormCluster(ctx, runtime, connect.Cluster, registries); err != nil {
		return err
	}
	return nil
}

func getRegistryStatus(ctx context.Context, runtime runtimes.Runtime, connect k3dRegistry.ConnectConfig) ([]map[string]string, error) {
	updatedStatus := make([]map[string]string, 0)
	cluster, err := cluster.GetCluster(ctx, runtime, connect.Cluster)
	if err != nil {
		return nil, err
	}
	registries, err := node.FilteredNodes(ctx, runtime, connect.Registries)
	if err != nil {
		return nil, err
	}
	for _, registry := range registries {
		log.Printf("registry networks: %v", registry.Networks)
		log.Printf("cluster network: %s", cluster.Network.Name)
		if utils2.Contains(registry.Networks, cluster.Network.Name) {
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

func getUpdatedRegistriesChanges(d *schema.ResourceData) (registries k3dRegistry.ConnectConfig) {
	oldRegistries, newRegistries := d.GetChange(utils2.TerraformResourceRegistries)
	if !cmp.Equal(oldRegistries, newRegistries) {
		registries.Registries = getSlice(newRegistries)
	}
	oldCluster, newCluster := d.GetChange(utils2.TerraformResourceCluster)
	if !cmp.Equal(oldCluster, newCluster) {
		registries.Cluster = utils2.String(newCluster)
	}
	oldConnect, newConnect := d.GetChange(utils2.TerraformResourceConnect)
	if !cmp.Equal(oldConnect, newConnect) {
		registries.Connect = utils2.Bool(newConnect)
	}
	return
}
