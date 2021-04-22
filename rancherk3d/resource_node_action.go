package rancherk3d

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "if enabled it starts a stopped nodes",
			},
			"stop": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "if enabled it stops a running nodes",
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
						"running": {
							Type:        schema.TypeBool,
							Required:    false,
							Computed:    true,
							Description: "enabled when current state of node is running",
						},
					},
				},
			},
		},
	}
}

func resourceNodeActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceNodeActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceNodeActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceNodeActionStartStop(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}
