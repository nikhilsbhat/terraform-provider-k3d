package rancherk3d

import (
	"context"
	"log"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/client"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/image"
	utils2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
)

func resourceImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadImageLoad,
		ReadContext:   resourceLoadImageRead,
		DeleteContext: resourceLoadImageDelete,
		UpdateContext: resourceLoadImageUpdate,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(utils2.TerraformTimeOut5 * time.Minute),
			Update: schema.DefaultTimeout(utils2.TerraformTimeOut5 * time.Minute),
			Delete: schema.DefaultTimeout(utils2.TerraformTimeOut5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"images": {
				Type:        schema.TypeList,
				Required:    true,
				Computed:    false,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "list of images to be imported to the existing cluster",
			},
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				Description: "name of the existing cluster to which the images has to be imported to",
			},
			"all": {
				Type:        schema.TypeBool,
				Computed:    false,
				Optional:    true,
				Description: "if enabled loads images to all available clusters",
			},
			"keep_tarball": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				Description: "enable to keep the tarball of the loaded images locally",
			},
			"images_stored": {
				Type:        schema.TypeList,
				Computed:    true,
				ForceNew:    false,
				Description: "list of images loaded to the cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "cluster to which the below images are stored",
						},
						"images": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: "details of images and its tarball stored, if in case keep_tarball is enabled",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceLoadImageLoad(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	if d.IsNewResource() {
		id := d.Id()

		if len(id) == 0 {
			newID, err := utils2.GetChecksum(utils2.String(d.Get(utils2.TerraformResourceCluster)))
			if err != nil {
				return diag.Errorf("errored while fetching randomID %v", err)
			}
			id = newID
		}

		imageCfg := image.Config{
			Images:       getSlice(d.Get(utils2.TerraformResourceImages)),
			StoreTarBall: utils2.Bool(d.Get(utils2.TerraformResourceKeepTarball)),
			Cluster:      utils2.String(d.Get(utils2.TerraformResourceCluster)),
			All:          utils2.Bool(d.Get(utils2.TerraformResourceAll)),
		}

		if err := imageCfg.Upload(ctx, defaultConfig.K3DRuntime); err != nil {
			return diag.Errorf("%v", err)
		}
		d.SetId(id)

		return resourceLoadImageRead(ctx, d, meta)
	}

	log.Printf("resource %s already exists", d.Id())

	return nil
}

func resourceLoadImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	imageCfg := image.Config{
		Images:  getSlice(d.Get(utils2.TerraformResourceImages)),
		Cluster: utils2.String(d.Get(utils2.TerraformResourceCluster)),
		All:     utils2.Bool(d.Get(utils2.TerraformResourceAll)),
	}

	imagesToStore, err := imageCfg.List(ctx, defaultConfig.K3DRuntime)
	if err != nil {
		d.SetId("")

		diag.Errorf("an error occurred while fetching images to be stored")
	}

	flattenedImagesToStore, err := utils2.MapSlice(imagesToStore)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while flattening images to store: %v", err)
	}
	if err := d.Set(utils2.TerraformResourceImagesStored, flattenedImagesToStore); err != nil {
		d.SetId("")

		return diag.Errorf("oops setting 'images_stored' errored with : %v", err)
	}
	if err := d.Set(utils2.TerraformResourceImages, imageCfg.Images); err != nil {
		d.SetId("")

		return diag.Errorf("oops setting 'images' errored with : %v", err)
	}
	if err := d.Set(utils2.TerraformResourceCluster, imageCfg.Cluster); err != nil {
		d.SetId("")

		return diag.Errorf("oops setting 'cluster' errored with : %v", err)
	}

	return nil
}

func resourceLoadImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceLoadImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	log.Printf("uploading newer images to k3d clusters")
	if d.HasChange(utils2.TerraformResourceCluster) || d.HasChange(utils2.TerraformResourceImages) {
		updatedCluster, updatedImages := getUpdatedClusterAndImages(d)
		// keepTarball := utils2.Bool(d.Get(utils2.TerraformResourceKeepTarball))
		// all := utils2.Bool(d.Get(utils2.TerraformResourceAll))

		imageCfg := image.Config{
			Images:       updatedImages,
			StoreTarBall: utils2.Bool(d.Get(utils2.TerraformResourceKeepTarball)),
			Cluster:      updatedCluster,
			All:          utils2.Bool(d.Get(utils2.TerraformResourceAll)),
		}

		if err := imageCfg.Upload(ctx, defaultConfig.K3DRuntime); err != nil {
			return diag.Errorf("%v", err)
		}

		if err := d.Set(utils2.TerraformResourceCluster, updatedCluster); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceCluster, err)
		}

		if err := d.Set(utils2.TerraformResourceImages, updatedImages); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils2.TerraformResourceImages, err)
		}

		return resourceLoadImageRead(ctx, d, meta)
	}

	log.Printf("nothing to update so skipping")

	return nil
}

//nolint:nonamedreturns
func getUpdatedClusterAndImages(d *schema.ResourceData) (cluster string, images []string) {
	oldCluster, newCluster := d.GetChange(utils2.TerraformResourceCluster)
	if !cmp.Equal(oldCluster, newCluster) {
		cluster = utils2.String(newCluster)
	}
	oldImages, newImages := d.GetChange(utils2.TerraformResourceImages)
	images = getSlice(oldImages)
	if !cmp.Equal(oldImages, newImages) {
		images = getSlice(newImages)
	}

	return
}
