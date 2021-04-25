package rancherk3d

import (
	"context"
	"log"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/k3d"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/utils"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

func resourceClusterAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterActionStartStop,
		ReadContext:   resourceClusterActionRead,
		DeleteContext: resourceClusterActionDelete,
		UpdateContext: resourceClusterActionUpdate,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"clusters": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "list of k3s clusters on which the action has to be applied",
			},
			"all": {
				Type:        schema.TypeBool,
				Computed:    false,
				Optional:    true,
				Description: "if enabled selected clusters would be started/stopped",
			},
			"start": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      false,
				ForceNew:      true,
				ConflictsWith: []string{"stop"},
				Description:   "if enabled it starts a stopped cluster",
			},
			"stop": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      false,
				ForceNew:      true,
				ConflictsWith: []string{"start"},
				Description:   "if enabled it stops a running start",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "latest state of selected clusters",
			},
			"status": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "updated status of clusters",
				Elem: &schema.Resource{
					Schema: resourceClusterSchema(),
				},
			},
		},
	}
}

func resourceClusterActionStartStop(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

		clusters := getClusterSlice(d.Get(utils.TerraformResourceClusters))
		all := utils.Bool(d.Get(utils.TerraformResourceAll))
		start := utils.Bool(d.Get(utils.TerraformResourceStart))
		stop := utils.Bool(d.Get(utils.TerraformResourceStop))

		actionStatus, action := getAction(start, stop)
		if !actionStatus {
			diag.Errorf("cannot start/stop at the same time, %v", actionStatus)
		}
		if err := updateClusterStatus(ctx, defaultConfig, clusters, action, all); err != nil {
			return diag.Errorf("creation failed with error: %v", err)
		}

		d.SetId(id)
		return resourceClusterActionRead(ctx, d, meta)
	}
	return nil
}

func resourceClusterActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.K3dConfig)

	clusters := getClusterSlice(d.Get(utils.TerraformResourceClusters))
	all := utils.Bool(d.Get(utils.TerraformResourceAll))
	start := utils.Bool(d.Get(utils.TerraformResourceStart))
	stop := utils.Bool(d.Get(utils.TerraformResourceStop))

	actionStatus, action := getAction(start, stop)
	if !actionStatus {
		diag.Errorf("cannot start/stop at the same time, %v", actionStatus)
	}

	clusterStatus, err := getK3dCluster(ctx, defaultConfig, clusters, all)
	if err != nil {
		return diag.Errorf("errored while fetching cluster status: %v", err)
	}
	flattenedClusterStatus, err := utils.Map(clusterStatus)
	log.Printf("flattenedClusterStatus, %v", flattenedClusterStatus)
	if err != nil {
		return diag.Errorf("errored while flattening clusters obtained: %v", err)
	}

	if err := d.Set(utils.TerraformResourceStatus, flattenedClusterStatus); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceStatus, err)
	}
	if err := d.Set(utils.TerraformResourceState, action); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceState, err)
	}
	return nil
}

func resourceClusterActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.K3dConfig)

	log.Printf("uploading newer images to k3d clusters")
	if d.HasChange(utils.TerraformResourceClusters) || d.HasChange(utils.TerraformResourceStart) || d.HasChange(utils.TerraformResourceStop) {

		all := utils.Bool(d.Get(utils.TerraformResourceAll))
		clusters, start, stop := getUpdatedClustersActionChanges(d)

		actionStatus, action := getAction(start, stop)
		if !actionStatus {
			diag.Errorf("cannot start/stop at the same time, %v", actionStatus)
		}

		if err := updateClusterStatus(ctx, defaultConfig, clusters, action, all); err != nil {
			return diag.Errorf("creation failed with error: %v", err)
		}

		if err := d.Set(utils.TerraformResourceClusters, clusters); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceCluster, err)
		}
		if err := d.Set(utils.TerraformResourceAll, all); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceAll, err)
		}
		return resourceClusterActionRead(ctx, d, meta)
	}

	log.Printf("nothing to update so skipping")
	return nil
}

func resourceClusterActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func updateClusterStatus(ctx context.Context, defaultConfig *k3d.K3dConfig, clusters []string, action string, all bool) error {
	if action == utils.TerraformResourceStart {
		return startClusters(ctx, defaultConfig, clusters, all)
	}
	return stopClusters(ctx, defaultConfig, clusters, all)
}

func startClusters(ctx context.Context, defaultConfig *k3d.K3dConfig, clusters []string, all bool) error {
	var fetchedClusters []*K3D.Cluster
	if all {
		cls, err := k3d.GetClusters(ctx, defaultConfig.K3DRuntime)
		if err != nil {
			return err
		}
		fetchedClusters = cls
	} else {
		cls, err := k3d.GetFilteredClusters(ctx, defaultConfig.K3DRuntime, clusters)
		if err != nil {
			return err
		}
		fetchedClusters = cls
	}
	return k3d.StartClusters(ctx, defaultConfig.K3DRuntime, fetchedClusters, K3D.ClusterStartOpts{})
}

func stopClusters(ctx context.Context, defaultConfig *k3d.K3dConfig, clusters []string, all bool) error {
	var fetchedClusters []*K3D.Cluster
	if all {
		cls, err := k3d.GetClusters(ctx, defaultConfig.K3DRuntime)
		if err != nil {
			return err
		}
		fetchedClusters = cls
	} else {
		cls, err := k3d.GetFilteredClusters(ctx, defaultConfig.K3DRuntime, clusters)
		if err != nil {
			return err
		}
		fetchedClusters = cls
	}
	return k3d.StopClusters(ctx, defaultConfig.K3DRuntime, fetchedClusters)
}

func getUpdatedClustersActionChanges(d *schema.ResourceData) (clusters []string, start, stop bool) {
	oldClusters, newClusters := d.GetChange(utils.TerraformResourceClusters)
	if !cmp.Equal(oldClusters, newClusters) {
		clusters = utils.GetSlice(newClusters.([]interface{}))
	}
	oldStart, newStart := d.GetChange(utils.TerraformResourceStart)
	if !cmp.Equal(oldStart, newStart) {
		start = utils.Bool(newStart)
	}
	oldStop, newStop := d.GetChange(utils.TerraformResourceStop)
	if !cmp.Equal(oldStop, newStop) {
		stop = utils.Bool(newStop)
	}
	return
}
