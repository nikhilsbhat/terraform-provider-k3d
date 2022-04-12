package rancherk3d

import (
	"context"
	"log"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/client"
	k3dNode "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/node"
	utils2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
)

func resourceNodeAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeActionStartStop,
		ReadContext:   resourceNodeActionRead,
		DeleteContext: resourceNodeActionDelete,
		UpdateContext: resourceNodeActionUpdate,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(utils2.TerraformTimeOut5 * time.Minute),
			Update: schema.DefaultTimeout(utils2.TerraformTimeOut5 * time.Minute),
			Delete: schema.DefaultTimeout(utils2.TerraformTimeOut5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"nodes": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "list of nodes on which the action has to be applied",
			},
			"all": {
				Type:        schema.TypeBool,
				Computed:    false,
				Optional:    true,
				Description: "if enabled fetches all the nodes available in the selected cluster",
			},
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "name of the cluster of which that nodes to be acted upon",
			},
			"start": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      false,
				ForceNew:      true,
				ConflictsWith: []string{"stop"},
				Description:   "if enabled it starts a stopped nodes",
			},
			"stop": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      false,
				ForceNew:      true,
				ConflictsWith: []string{"start"},
				Description:   "if enabled it stops a running nodes",
			},
			"status": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "updated status of started/stopped nodes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node": {
							Type:        schema.TypeString,
							Required:    false,
							Computed:    true,
							Description: "node of which the current status is updated with",
						},
						"role": {
							Type:        schema.TypeString,
							Required:    false,
							Computed:    true,
							Description: "role of updated node",
						},
						"state": {
							Type:        schema.TypeString,
							Required:    false,
							Computed:    true,
							Description: "current state of the node specified",
						},
						"cluster": {
							Type:        schema.TypeString,
							Required:    true,
							Computed:    false,
							ForceNew:    true,
							Description: "name of the cluster of to which node belongs",
						},
					},
				},
			},
		},
	}
}

func resourceNodeActionStartStop(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

		actionStatus, action := getAction(utils2.Bool(d.Get(utils2.TerraformResourceStart)),
			utils2.Bool(d.Get(utils2.TerraformResourceStop)))

		if !actionStatus {
			diag.Errorf("cannot start/stop at the same time, %v", actionStatus)
		}

		cfg := k3dNode.Config{
			Name:              getSlice(d.Get(utils2.TerraformResourceNodes)),
			ClusterAssociated: utils2.String(d.Get(utils2.TerraformResourceCluster)),
			All:               utils2.Bool(d.Get(utils2.TerraformResourceAll)),
			Action:            action,
		}

		if err := cfg.StartStopNode(ctx, defaultConfig.K3DRuntime); err != nil {
			return diag.Errorf("creation failed with error: %v", err)
		}

		d.SetId(id)
		return resourceNodeActionRead(ctx, d, meta)
	}

	log.Printf("resource %s already exists", d.Id())
	return nil
}

func resourceNodeActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	cfg := k3dNode.Config{
		Name:              getSlice(d.Get(utils2.TerraformResourceNodes)),
		ClusterAssociated: utils2.String(d.Get(utils2.TerraformResourceCluster)),
		All:               utils2.Bool(d.Get(utils2.TerraformResourceAll)),
	}

	status, err := cfg.GetNodeStatus(ctx, defaultConfig.K3DRuntime)
	if err != nil {
		return diag.Errorf("errored while fetching nodes: %v", err)
	}

	flattenedNodeStatus, err := utils2.MapSlice(status)
	if err != nil {
		return diag.Errorf("errored while flattening nodes obtained: %v", err)
	}

	if err = d.Set(utils2.TerraformResourceStatus, flattenedNodeStatus); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceStatus, err)
	}
	return nil
}

func resourceNodeActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	if d.HasChange(utils2.TerraformResourceCluster) || d.HasChange(utils2.TerraformResourceNodes) ||
		d.HasChange(utils2.TerraformResourceStart) || d.HasChange(utils2.TerraformResourceStop) {
		nodes, cluster, start, stop := getUpdatedNodeActionChanges(d)

		actionStatus, action := getAction(start, stop)
		if !actionStatus {
			diag.Errorf("cannot start/stop at the same time, %v", actionStatus)
		}

		nodesConfig := k3dNode.Config{
			Name:              nodes,
			ClusterAssociated: cluster,
			Action:            action,
			All:               utils2.Bool(d.Get(utils2.TerraformResourceAll)),
		}

		if err := nodesConfig.StartStopNode(ctx, defaultConfig.K3DRuntime); err != nil {
			return diag.Errorf("creation failed with error: %v", err)
		}

		if err := d.Set(utils2.TerraformResourceCluster, cluster); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceCluster, err)
		}
		if err := d.Set(utils2.TerraformResourceNodes, nodes); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceNodes, err)
		}
		if err := d.Set(utils2.TerraformResourceAll, utils2.Bool(d.Get(utils2.TerraformResourceAll))); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceAll, err)
		}
		return resourceNodeActionRead(ctx, d, meta)
	}

	log.Printf("nothing to update so skipping")
	return nil
}

func resourceNodeActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func getUpdatedNodeActionChanges(d *schema.ResourceData) (nodes []string, cluster string, start, stop bool) {
	oldNodes, newNodes := d.GetChange(utils2.TerraformResourceNodes)
	if !cmp.Equal(oldNodes, newNodes) {
		nodes = utils2.GetSlice(newNodes.([]interface{}))
	}
	oldCluster, newCluster := d.GetChange(utils2.TerraformResourceCluster)
	if !cmp.Equal(oldCluster, newCluster) {
		cluster = utils2.String(newCluster)
	}
	oldStart, newStart := d.GetChange(utils2.TerraformResourceStart)
	if !cmp.Equal(oldStart, newStart) {
		start = utils2.Bool(newStart)
	}
	oldStop, newStop := d.GetChange(utils2.TerraformResourceStop)
	if !cmp.Equal(oldStop, newStop) {
		stop = utils2.Bool(newStop)
	}
	return
}

//func updateNodeStats(ctx context.Context, defaultConfig *client.Config, cluster, action string, nodes []string, all bool) error {
//	if action == utils2.TerraformResourceStart {
//		return startNodes(ctx, defaultConfig, cluster, nodes, all)
//	}
//	return stopNodes(ctx, defaultConfig, cluster, nodes, all)
//}

//func startNodes(ctx context.Context, defaultConfig *client.Config, cluster string, nodes []string, all bool) error {
//	if all {
//		return k3dNode.StartNodesFromCluster(ctx, defaultConfig.K3DRuntime, cluster)
//	}
//	return k3dNode.StartNodes(ctx, defaultConfig.K3DRuntime, nodes)
//}
//
//func stopNodes(ctx context.Context, defaultConfig *client.Config, cluster string, nodes []string, all bool) error {
//	if all {
//		return k3dNode.StopNodesFromCluster(ctx, defaultConfig.K3DRuntime, cluster)
//	}
//	return k3dNode.StopNodes(ctx, defaultConfig.K3DRuntime, nodes)
//}

//func getNodeStatus(ctx context.Context, defaultConfig *client.Config, cfg k3dNode.Config) ([]*k3dNode.Status, error) {
//	k3dNodes, err := cfg.GetFilteredNodesFromCluster(ctx, defaultConfig.K3DRuntime)
//	if err != nil {
//		return nil, fmt.Errorf("an error occurred while fetching nodes information : %s", err.Error())
//	}
//
//	nodeCurrentStatus := make([]*k3dNode.Status, 0)
//	for _, node := range k3dNodes {
//		nodeCurrentStatus = append(nodeCurrentStatus, &k3dNode.Status{
//			Node:    node.Name[0],
//			Cluster: node.ClusterAssociated,
//			State:   node.State,
//			Role:    node.Role,
//		})
//	}
//	return nodeCurrentStatus, nil
//}

func getAction(start, stop bool) (bool, string) {
	if start && stop {
		return false, ""
	}
	if start {
		return true, utils2.TerraformResourceStart
	}
	return true, utils2.TerraformResourceStop
}
