package rancherk3d

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceRegistrySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "name of the registry",
		},
		"role": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "role of registry created/retrieved",
		},
		"image": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "image used for registry",
		},
		"cluster": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "cluster to which the registry belongs",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "current state of registry node",
		},
		"created": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "creation timestamp of registry node",
		},
		"networks": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "networks associated with the registries",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"env": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "environment variables set in the node",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"port_mappings": {
			Type:        schema.TypeMap,
			Computed:    true,
			Optional:    true,
			Description: "port mappings",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}
