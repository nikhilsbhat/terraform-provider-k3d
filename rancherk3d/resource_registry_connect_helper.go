package rancherk3d

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/client"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	k3dRegistry "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/registry"
	utils2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
	"github.com/rancher/k3d/v5/pkg/runtimes"
)

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

func connectRegistryToCluster(ctx context.Context, runtime runtimes.Runtime, config k3dRegistry.Config) error {
	if config.ConnectToCluster {
		if err := config.Connect(ctx, runtime); err != nil {
			return err
		}
		return nil
	}
	if err := config.Disconnect(ctx, runtime); err != nil {
		return err
	}
	return nil
}

func getRegistryStatus(ctx context.Context, runtime runtimes.Runtime, config k3dRegistry.Config) ([]map[string]string, error) {
	updatedStatus := make([]map[string]string, 0)
	clusterData, err := cluster.GetCluster(ctx, runtime, config.Cluster)
	if err != nil {
		return nil, err
	}

	registries, err := config.Get(ctx, runtime)
	if err != nil {
		return nil, err
	}

	for _, registry := range registries {
		if utils2.Contains(registry.Networks, clusterData.Network.Name) {
			updatedStatus = append(updatedStatus, config.GetRegistryStatus(registry.Name[0], utils2.RegistryConnectedState))
		} else {
			updatedStatus = append(updatedStatus, config.GetRegistryStatus(registry.Name[0], utils2.RegistryDisconnectedState))
		}
	}
	return updatedStatus, nil
}

func getUpdatedRegistriesChanges(d *schema.ResourceData) (registries k3dRegistry.Config) {
	oldRegistries, newRegistries := d.GetChange(utils2.TerraformResourceRegistries)
	if !cmp.Equal(oldRegistries, newRegistries) {
		registries.Name = getSlice(newRegistries)
	}
	oldCluster, newCluster := d.GetChange(utils2.TerraformResourceCluster)
	if !cmp.Equal(oldCluster, newCluster) {
		registries.Cluster = utils2.String(newCluster)
	}
	oldConnect, newConnect := d.GetChange(utils2.TerraformResourceConnect)
	if !cmp.Equal(oldConnect, newConnect) {
		registries.ConnectToCluster = utils2.Bool(newConnect)
	}
	return
}
