package rancherk3d

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/k3d"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/utils"
)

func resourceNodeAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeActionStartStop,
		ReadContext:   resourceNodeActionRead,
		DeleteContext: resourceNodeActionDelete,
		UpdateContext: resourceNodeActionUpdate,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
				Description: "if enabled it stops a running nodes",
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
							Description: "node of which the current status is updated with",
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
	defaultConfig := meta.(*k3d.K3dConfig)

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

		nodes := getNodesSlice(d.Get(utils.TerraformResourceNodes))
		cluster := utils.String(d.Get(utils.TerraformResourceCluster))
		all := utils.Bool(d.Get(utils.TerraformResourceAll))
		start := utils.Bool(d.Get(utils.TerraformResourceStart))
		stop := utils.Bool(d.Get(utils.TerraformResourceStop))

		actionStatus, action := getAction(start, stop)
		if !actionStatus {
			diag.Errorf("cannot start/stop at the same time, %v", actionStatus)
		}
		if err := updateNodeStats(ctx, defaultConfig, cluster, action, nodes, all); err != nil {
			return diag.Errorf("creation failed with error: %v", err)
		}

		d.SetId(id)
		return resourceNodeActionRead(ctx, d, meta)
	}

	log.Printf("resource %s already exists", d.Id())
	return nil
}

func resourceNodeActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.K3dConfig)

	nodes := getNodesSlice(d.Get(utils.TerraformResourceNodes))
	cluster := utils.String(d.Get(utils.TerraformResourceCluster))
	all := utils.Bool(d.Get(utils.TerraformResourceAll))

	nodeStatus, err := getNodeStatus(ctx, defaultConfig, cluster, nodes, all)
	if err != nil {
		return diag.Errorf("errored while fetching nodes: %v", err)
	}
	flattenedNodeStatus, err := utils.Map(nodeStatus)
	log.Printf("flattenedNodeStatus, %v", flattenedNodeStatus)
	if err != nil {
		return diag.Errorf("errored while flattening nodes obtained: %v", err)
	}

	if err := d.Set(utils.TerraformResourceStatus, flattenedNodeStatus); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceNodesList, err)
	}
	return nil
}

func resourceNodeActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.K3dConfig)

	log.Printf("uploading newer images to k3d clusters")
	if d.HasChange(utils.TerraformResourceCluster) || d.HasChange(utils.TerraformResourceNodes) || d.HasChange(utils.TerraformResourceStart) || d.HasChange(utils.TerraformResourceStop) {

		all := utils.Bool(d.Get(utils.TerraformResourceAll))
		nodes, cluster, start, stop := getUpdatedNodeActionChanges(d)

		actionStatus, action := getAction(start, stop)
		if !actionStatus {
			diag.Errorf("cannot start/stop at the same time, %v", actionStatus)
		}

		if err := updateNodeStats(ctx, defaultConfig, cluster, action, nodes, all); err != nil {
			return diag.Errorf("creation failed with error: %v", err)
		}

		if err := d.Set(utils.TerraformResourceCluster, cluster); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceCluster, err)
		}
		if err := d.Set(utils.TerraformResourceNodes, nodes); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceCluster, err)
		}
		if err := d.Set(utils.TerraformResourceAll, all); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceAll, err)
		}
		return resourceNodeActionRead(ctx, d, meta)
	}

	log.Printf("nothing to update so skipping")
	return nil
}

func resourceNodeActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.K3dConfig)
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
	oldNodes, newNodes := d.GetChange(utils.TerraformResourceNodes)
	if !cmp.Equal(oldNodes, newNodes) {
		nodes = utils.GetSlice(newNodes.([]interface{}))
	}
	oldCluster, newCluster := d.GetChange(utils.TerraformResourceCluster)
	if !cmp.Equal(oldCluster, newCluster) {
		cluster = utils.String(newCluster)
	}
	oldStart, newStart := d.GetChange(utils.TerraformResourceNodes)
	if !cmp.Equal(oldStart, newStart) {
		start = utils.Bool(newStart)
	}
	oldStop, newStop := d.GetChange(utils.TerraformResourceCluster)
	if !cmp.Equal(oldStop, newStop) {
		stop = utils.Bool(newStop)
	}
	return
}

func updateNodeStats(ctx context.Context, defaultConfig *k3d.K3dConfig, cluster, action string, nodes []string, all bool) error {
	if action == utils.TerraformResourceStart {
		return startNodes(ctx, defaultConfig, cluster, nodes, all)
	}
	return stopNodes(ctx, defaultConfig, cluster, nodes, all)
}

func startNodes(ctx context.Context, defaultConfig *k3d.K3dConfig, cluster string, nodes []string, all bool) error {
	if all {
		return k3d.StartNodesFromCluster(ctx, defaultConfig.K3DRuntime, cluster)
	}
	return k3d.StartNodes(ctx, defaultConfig.K3DRuntime, nodes)
}

func stopNodes(ctx context.Context, defaultConfig *k3d.K3dConfig, cluster string, nodes []string, all bool) error {
	if all {
		return k3d.StopNodesFromCluster(ctx, defaultConfig.K3DRuntime, cluster)
	}
	return k3d.StopNodes(ctx, defaultConfig.K3DRuntime, nodes)
}

func getNodeStatus(ctx context.Context, defaultConfig *k3d.K3dConfig, cluster string, nodes []string, all bool) ([]*k3d.K3DNodeStatus, error) {
	k3dNodes, err := getNodesFromCluster(ctx, defaultConfig, cluster, nodes, all)
	if err != nil {
		return nil, fmt.Errorf("an error occured while fetching nodes information : %s", err.Error())
	}
	nodeCurrentStatus := make([]*k3d.K3DNodeStatus, 0)
	for _, node := range k3dNodes {
		nodeCurrentStatus = append(nodeCurrentStatus, &k3d.K3DNodeStatus{
			Node:    node.Name,
			Cluster: node.ClusterAssociated,
			State:   node.State,
			Role:    node.Role,
		})
	}
	return nodeCurrentStatus, nil
}

func getAction(start, stop bool) (bool, string) {
	if start == true && stop == true {
		return false, ""
	}
	if start == true {
		return true, "start"
	}
	return true, "stop"
}
