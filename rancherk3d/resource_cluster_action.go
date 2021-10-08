package rancherk3d

import (
	"context"
	"log"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/client"
	k3d2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	utils2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
	K3D "github.com/rancher/k3d/v4/pkg/types"
)

func resourceClusterAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterActionStartStop,
		ReadContext:   resourceClusterActionRead,
		DeleteContext: resourceClusterActionDelete,
		UpdateContext: resourceClusterActionUpdate,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(utils2.TerraformTimeOut5 * time.Minute),
			Update: schema.DefaultTimeout(utils2.TerraformTimeOut5 * time.Minute),
			Delete: schema.DefaultTimeout(utils2.TerraformTimeOut5 * time.Minute),
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
				Description:   "if enabled it stops a running cluster",
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

		clusters := getClusterSlice(d.Get(utils2.TerraformResourceClusters))
		all := utils2.Bool(d.Get(utils2.TerraformResourceAll))
		start := utils2.Bool(d.Get(utils2.TerraformResourceStart))
		stop := utils2.Bool(d.Get(utils2.TerraformResourceStop))

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
	defaultConfig := meta.(*client.Config)

	clusters := getClusterSlice(d.Get(utils2.TerraformResourceClusters))
	all := utils2.Bool(d.Get(utils2.TerraformResourceAll))
	start := utils2.Bool(d.Get(utils2.TerraformResourceStart))
	stop := utils2.Bool(d.Get(utils2.TerraformResourceStop))

	actionStatus, action := getAction(start, stop)
	if !actionStatus {
		diag.Errorf("cannot start/stop at the same time, %v", actionStatus)
	}

	clusterStatus, err := getK3dCluster(ctx, defaultConfig, clusters, all)
	if err != nil {
		return diag.Errorf("errored while fetching cluster status: %v", err)
	}

	flattenedClusterStatus, err := utils2.MapSlice(clusterStatus)
	if err != nil {
		return diag.Errorf("errored while flattening clusters obtained: %v", err)
	}

	if err = d.Set(utils2.TerraformResourceStatus, flattenedClusterStatus); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceStatus, err)
	}
	if err = d.Set(utils2.TerraformResourceState, action); err != nil {
		return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceState, err)
	}
	return nil
}

func resourceClusterActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	log.Printf("uploading newer images to k3d clusters")
	if d.HasChange(utils2.TerraformResourceClusters) || d.HasChange(utils2.TerraformResourceStart) ||
		d.HasChange(utils2.TerraformResourceStop) {
		all := utils2.Bool(d.Get(utils2.TerraformResourceAll))
		clusters, start, stop := getUpdatedClustersActionChanges(d)

		actionStatus, action := getAction(start, stop)
		if !actionStatus {
			diag.Errorf("cannot start/stop at the same time, %v", actionStatus)
		}

		if err := updateClusterStatus(ctx, defaultConfig, clusters, action, all); err != nil {
			return diag.Errorf("creation failed with error: %v", err)
		}

		if err := d.Set(utils2.TerraformResourceClusters, clusters); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceCluster, err)
		}
		if err := d.Set(utils2.TerraformResourceAll, all); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceAll, err)
		}
		return resourceClusterActionRead(ctx, d, meta)
	}

	log.Printf("nothing to update so skipping")
	return nil
}

func resourceClusterActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func updateClusterStatus(ctx context.Context, defaultConfig *client.Config, clusters []string, action string, all bool) error {
	if action == utils2.TerraformResourceStart {
		return startClusters(ctx, defaultConfig, clusters, all)
	}
	return stopClusters(ctx, defaultConfig, clusters, all)
}

func startClusters(ctx context.Context, defaultConfig *client.Config, clusters []string, all bool) error {
	var fetchedClusters []*K3D.Cluster
	if all {
		cls, err := k3d2.GetClusters(ctx, defaultConfig.K3DRuntime)
		if err != nil {
			return err
		}
		fetchedClusters = cls
	} else {
		cls, err := k3d2.GetFilteredClusters(ctx, defaultConfig.K3DRuntime, clusters)
		if err != nil {
			return err
		}
		fetchedClusters = cls
	}
	return k3d2.StartClusters(ctx, defaultConfig.K3DRuntime, fetchedClusters, K3D.ClusterStartOpts{})
}

func stopClusters(ctx context.Context, defaultConfig *client.Config, clusters []string, all bool) error {
	var fetchedClusters []*K3D.Cluster
	if all {
		cls, err := k3d2.GetClusters(ctx, defaultConfig.K3DRuntime)
		if err != nil {
			return err
		}
		fetchedClusters = cls
	} else {
		cls, err := k3d2.GetFilteredClusters(ctx, defaultConfig.K3DRuntime, clusters)
		if err != nil {
			return err
		}
		fetchedClusters = cls
	}
	return k3d2.StopClusters(ctx, defaultConfig.K3DRuntime, fetchedClusters)
}

func getUpdatedClustersActionChanges(d *schema.ResourceData) (clusters []string, start, stop bool) {
	oldClusters, newClusters := d.GetChange(utils2.TerraformResourceClusters)
	if !cmp.Equal(oldClusters, newClusters) {
		clusters = utils2.GetSlice(newClusters.([]interface{}))
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
