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
	"github.com/rancher/k3d/v4/pkg/runtimes"
)

func resourceImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadImageLoad,
		ReadContext:   resourceLoadImageRead,
		DeleteContext: resourceLoadImageDelete,
		UpdateContext: resourceLoadImageUpdate,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(utils.TerraformTimeOut5 * time.Minute),
			Update: schema.DefaultTimeout(utils.TerraformTimeOut5 * time.Minute),
			Delete: schema.DefaultTimeout(utils.TerraformTimeOut5 * time.Minute),
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
	defaultConfig := meta.(*k3d.Config)

	if d.IsNewResource() {
		id := d.Id()

		if len(id) == 0 {
			newID, err := utils.GetChecksum(utils.String(d.Get(utils.TerraformResourceCluster)))
			if err != nil {
				return diag.Errorf("errored while fetching randomID %v", err)
			}
			id = newID
		}

		images := getImages(d.Get(utils.TerraformResourceImages))
		keepTarball := utils.Bool(d.Get(utils.TerraformResourceKeepTarball))
		cluster := utils.String(d.Get(utils.TerraformResourceCluster))
		all := utils.Bool(d.Get(utils.TerraformResourceAll))

		if err := storeImagesToClusters(ctx, defaultConfig.K3DRuntime, images, cluster, all, keepTarball); err != nil {
			return diag.Errorf("%v", err)
		}
		d.SetId(id)
		return resourceLoadImageRead(ctx, d, meta)
	}

	log.Printf("resource %s already exists", d.Id())
	return nil
}

func resourceLoadImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.Config)

	images := getImages(d.Get(utils.TerraformResourceImages))
	cluster := utils.String(d.Get(utils.TerraformResourceCluster))
	all := utils.Bool(d.Get(utils.TerraformResourceAll))

	imagesToStore, err := getImagesToBeStored(ctx, defaultConfig.K3DRuntime, images, cluster, all)
	if err != nil {
		d.SetId("")
		diag.Errorf("an error occurred while fetching images to be stored")
	}

	flattenedImagesToStore, err := utils.MapSlice(imagesToStore)
	if err != nil {
		d.SetId("")
		return diag.Errorf("errored while flattening images to store: %v", err)
	}
	if err := d.Set(utils.TerraformResourceImagesStored, flattenedImagesToStore); err != nil {
		d.SetId("")
		return diag.Errorf("oops setting 'images_stored' errored with : %v", err)
	}
	if err := d.Set(utils.TerraformResourceImages, images); err != nil {
		d.SetId("")
		return diag.Errorf("oops setting 'images' errored with : %v", err)
	}
	if err := d.Set(utils.TerraformResourceCluster, cluster); err != nil {
		d.SetId("")
		return diag.Errorf("oops setting 'cluster' errored with : %v", err)
	}
	return nil
}

func resourceLoadImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.Config)
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
	defaultConfig := meta.(*k3d.Config)

	log.Printf("uploading newer images to k3d clusters")
	if d.HasChange(utils.TerraformResourceCluster) || d.HasChange(utils.TerraformResourceImages) {
		updatedCluster, updatedImages := getUpdatedClusterAndImages(d)
		keepTarball := utils.Bool(d.Get(utils.TerraformResourceKeepTarball))
		all := utils.Bool(d.Get(utils.TerraformResourceAll))

		if err := storeImagesToClusters(ctx, defaultConfig.K3DRuntime, updatedImages, updatedCluster, all, keepTarball); err != nil {
			return diag.Errorf("%v", err)
		}

		if err := d.Set(utils.TerraformResourceCluster, updatedCluster); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceCluster, err)
		}
		if err := d.Set(utils.TerraformResourceImages, updatedImages); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", utils.TerraformResourceImages, err)
		}
		return resourceLoadImageRead(ctx, d, meta)
	}

	log.Printf("nothing to update so skipping")
	return nil
}

func storeImagesToClusters(ctx context.Context, runtime runtimes.Runtime, images []string, cluster string, all, storeTarball bool) error {
	if all {
		return k3d.StoreImagesToClusters(ctx, runtime, images, storeTarball)
	}
	return k3d.StoreImagesToCluster(ctx, runtime, images, cluster, storeTarball)
}

func getUpdatedClusterAndImages(d *schema.ResourceData) (cluster string, images []string) {
	oldCluster, newCluster := d.GetChange(utils.TerraformResourceCluster)
	if !cmp.Equal(oldCluster, newCluster) {
		cluster = utils.String(newCluster)
	}
	oldImages, newImages := d.GetChange(utils.TerraformResourceImages)
	images = getImages(oldImages)
	if !cmp.Equal(oldImages, newImages) {
		images = getImages(newImages)
	}
	return
}

func getImagesToBeStored(ctx context.Context, runtime runtimes.Runtime,
	images []string, cluster string, all bool) ([]*k3d.StoredImages, error) {
	if all {
		imagesToBeStored, err := k3d.GetImagesLoadedClusters(ctx, runtime, images)
		if err != nil {
			return nil, err
		}
		return imagesToBeStored, nil
	}
	imagesToBeStored, err := k3d.GetImagesLoadedCluster(ctx, runtime, images, cluster)
	if err != nil {
		return nil, err
	}
	return imagesToBeStored, nil
}

func getImages(images interface{}) []string {
	return utils.GetSlice(images.([]interface{}))
}

func getClusters(clusters interface{}) []string {
	return utils.GetSlice(clusters.([]interface{}))
}
