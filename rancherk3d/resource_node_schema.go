package rancherk3d

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceNodeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "name of the node",
		},
		"role": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "role of node created/retrieved",
		},
		"cluster": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "cluster to which the node belongs",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "current state of node",
		},
		"created": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "creation timestamp of node",
		},
		"memory": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "memory limit imposed on node",
		},
		"volumes": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "volumes associated with the nodes",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"networks": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "networks associated with the nodes",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"env": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "environment variables set in the node",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}
