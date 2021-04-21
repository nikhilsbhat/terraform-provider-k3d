package rancherk3d

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/k3d"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/utils"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

func dataSourceClusterList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceListClusterRead,
		Schema: map[string]*schema.Schema{
			"clusters": {
				Type:        schema.TypeList,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "list of clusters of which the information to be fetched",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"all": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "if enabled fetches all clusters available",
			},
			"clusters_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of cluster of which information has been retrieved",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
		},
	}
}

func dataSourceListClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.K3dConfig)

	id := d.Id()

	if len(id) == 0 {
		newID, err := utils.GetRandomID()
		if err != nil {
			d.SetId("")
			return diag.Errorf("errored while fetching randomID %v", err)
		}
		id = newID
	}

	clusters := getClusters(d.Get(utils.TerraformResourceClusters))
	all := utils.Bool(d.Get(utils.TerraformResourceAll))

	k3dClusters, err := getK3dCluster(ctx, defaultConfig, clusters, all)
	if err != nil {
		d.SetId("")
		return diag.Errorf("errored while fetching clusters: %v", err)
	}

	flattenedClusters, err := utils.Map(k3dClusters)
	if err != nil {
		d.SetId("")
		return diag.Errorf("errored while flattening nodes obtained: %v", err)
	}
	d.SetId(id)
	if err := d.Set(utils.TerraformResourceClusterList, flattenedClusters); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceClusterList, err)
	}

	return nil
}

func getK3dCluster(ctx context.Context, defaultConfig *k3d.K3dConfig, clusters []string, all bool) ([]*k3d.Cluster, error) {
	fetchedClusters := make([]*K3D.Cluster, 0)
	if all {
		allClusters, err := k3d.GetClusters(ctx, defaultConfig.K3DRuntime)
		if err != nil {
			return nil, err
		}
		fetchedClusters = allClusters
	} else {
		allClusters, err := k3d.GetFilteredClusters(ctx, defaultConfig.K3DRuntime, clusters)
		if err != nil {
			return nil, err
		}
		fetchedClusters = allClusters
	}
	filteredClusterinfo := make([]*k3d.Cluster, 0)
	for _, cluster := range fetchedClusters {
		serverCount, _ := cluster.ServerCountRunning()
		agentsCount, agentsRunning := cluster.AgentCountRunning()
		filteredClusterinfo = append(filteredClusterinfo, &k3d.Cluster{
			Name:            cluster.Name,
			Nodes:           getNodesList(cluster.Nodes),
			Network:         cluster.Network.Name,
			Token:           cluster.Token,
			ServersCount:    int64(serverCount),
			AgentsCount:     int64(agentsCount),
			AgentsRunning:   int64(agentsRunning),
			ImageVolume:     cluster.ImageVolume,
			HasLoadBalancer: cluster.HasLoadBalancer(),
		})
	}
	return filteredClusterinfo, nil
}

func getNodesList(rawNodes []*K3D.Node) (nodes []string) {
	for _, node := range rawNodes {
		nodes = append(nodes, node.Name)
	}
	return
}
