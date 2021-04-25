package rancherk3d

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceClusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "cluster name that was fetched",
		},
		"nodes": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "list of nodes present in cluster",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"network": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "network associated with the cluster",
		},
		"cluster_token": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "token of the cluster",
		},
		"servers_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "count of servers",
		},
		"servers_running": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "count of servers running",
		},
		"agents_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "count of agents in the cluster",
		},
		"agents_running": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "number of agents running in the cluster",
		},
		"has_loadbalancer": {
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
			Description: "details of images and its tarball stored, if in case keep_tarball is enabled",
		},
		"image_volume": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "volume to import images",
		},
	}
}
