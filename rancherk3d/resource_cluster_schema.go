package rancherk3d

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

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
			Description: "attribute that notifies the presence of loadbalancer in the cluster",
		},
		"image_volume": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "volume to import images",
		},
	}
}

func resourceClusterRegistriesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"create": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    false,
			Description: "creates a default registry to be used with the cluster",
		},
		"use": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    false,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "some other k3d-managed registry",
		},
		"config": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    false,
			Description: "define contents of the `registries.yaml` file (or reference a file); same as `--registry-config /path/to/config.yaml`",
		},
	}
}

func resourceClusterEnvsAndLabelsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Computed:    false,
			Description: "key of key-value pair",
		},
		"value": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Computed:    false,
			Description: "value of key-value pair",
		},
		"nodeFilters": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: false,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func resourceClusterRuntimeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"gpu_request": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "GPU devices to add to the cluster node containers ('all' to pass all GPUs) [From docker].",
		},
	}
}

func resourceClusterVolumeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"source": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "Source path of volume mount",
		},
		"destination": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "Destination path for the volume",
		},
		"node_filters": {
			Type:     schema.TypeList,
			ForceNew: true,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func resourceClusterK3dOptionsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"wait": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Description: "Wait for the server(s) to be ready before returning. Use '--timeout DURATION' to not wait forever. (default true).",
		},
		"timeout": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "Rollback changes if cluster couldn't be created in specified duration.",
		},
		"no_host_ip": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Description: "Disable the automatic injection of the Host IP as 'host.k3d.internal' into the containers and CoreDNS.",
		},
		"no_image_volume": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Description: "Disable the creation of a volume for importing images.",
		},
		"no_loadbalancer": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Description: "Disable the creation of a LoadBalancer in front of the server nodes.",
		},
		"no_rollback": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Description: "Disable the automatic rollback actions, if anything goes wrong.",
		},
	}
}

func resourceClusterK3sOptionsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"k3s_agent_arg": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    false,
			Description: "Additional args passed to the k3s agent command on agent nodes",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"k3s_server_arg": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    false,
			Description: "Additional args passed to the k3s server command on server nodes",
			Elem:        &schema.Schema{Type: schema.TypeList},
		},
	}
}

func resourceClusterPortsConfig() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host": {
			Type:     schema.TypeString,
			ForceNew: true,
			Optional: true,
		},
		"host_port": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validation.IsPortNumber,
		},
		"container_port": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validation.IsPortNumber,
		},
		"protocol": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, true),
		},
		"node_filters": {
			Type:     schema.TypeList,
			ForceNew: true,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}
